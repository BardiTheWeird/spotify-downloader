package downloader

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"spotify-downloader/clihelpers"
	"spotify-downloader/models"
	"strconv"
)

type DownloadHelper struct {
	// trackId -> models.DownloadEntry
	downloadEntries downloadEntrySyncMap
	*clihelpers.CliHelper
}

type DownloadStartStatus int

const (
	DStartOk DownloadStartStatus = iota
	DStartErrorInvalidUrl
	DStartErrorCreatingFile
	DStartErrorSendingRequest
	DStartErrorReadingContentLength
)

type readerWithCancellationFunc func(p []byte) (n int, err error)

func (rf readerWithCancellationFunc) Read(p []byte) (n int, err error) {
	return rf(p)
}

func (d *DownloadHelper) StartDownload(trackId, downloadFolder, filename, url string, metadata clihelpers.FfmpegMetadata) DownloadStartStatus {
	ch := make(chan DownloadStartStatus)
	filepathNoExt := filepath.Join(downloadFolder, filename)
	filepathTmp := filepathNoExt + ".tmp"

	downloadTrack := func(out *os.File) (bool, error) {
		cancelled := false

		var size int
		resp, err := http.Get(url)
		if err != nil {
			ch <- DStartErrorSendingRequest
			return false, err
		}
		defer resp.Body.Close()

		size, err = strconv.Atoi(resp.Header.Get("Content-Length"))
		if err != nil || size == 0 {
			ch <- DStartErrorInvalidUrl
			return false, fmt.Errorf("invalid url")
		}

		log.Println("download at", filepathNoExt, "started")
		ch <- DStartOk

		ctx, fn := context.WithCancel(context.Background())

		d.downloadEntries.Store(
			trackId,
			models.DownloadEntry{
				TotalBytes:       size,
				Status:           models.DownloadInProgress,
				FilepathNoExt:    filepathNoExt,
				CancellationFunc: fn,
			},
		)

		readerWithCancellation := readerWithCancellationFunc(func(p []byte) (n int, err error) {
			select {
			case <-ctx.Done():
				cancelled = true
				return 0, io.EOF
			default:
				return resp.Body.Read(p)
			}
		})

		bytesWritten, err := io.Copy(out, readerWithCancellation)

		fmt.Println("download at", filepathNoExt, "was finished:\n",
			"\tbytesWritten:", bytesWritten,
			"\n\terr:", err)

		return cancelled, err
	}

	go func() {
		cancelled := false
		func() {
			out, err := os.Create(filepathTmp)
			if err != nil {
				ch <- DStartErrorCreatingFile
				return
			}
			defer out.Close()

			cancelled, err = downloadTrack(out)

			out.Close()

			entry, ok := d.downloadEntries.Load(trackId)
			if !ok {
				entry = models.DownloadEntry{}
			}

			switch {
			case cancelled:
				log.Println("download at", filepathNoExt, "was cancelled")
				entry.Status = models.DownloadedCancelled
			case err != nil:
				log.Println("download at", filepathNoExt, "failed:", err)
				entry.Status = models.DownloadFailed
				os.Remove(filepathTmp)
			default:
				log.Println("download at", filepathNoExt, "was finished")

				if d.Features.Ffmpeg.Installed {
					entry.Status = models.DownloadConvertationInProgress
					d.downloadEntries.Store(trackId, entry)

					err := d.FfmpegConvert(filepathTmp, filepathNoExt+".mp3", metadata)
					if err != nil {
						entry.Status = models.DownloadErrorConverting
					} else {
						entry.Status = models.DownloadFinished
					}
					os.Remove(filepathTmp)
				} else {
					os.Rename(filepathTmp, filepathNoExt+".mp4")
				}
			}

			d.downloadEntries.Store(trackId, entry)
		}()

		if cancelled {
			os.Remove(filepathTmp)
		}
	}()

	return <-ch
}

type GetDownloadStatusStatus int

const (
	DStatusOk GetDownloadStatusStatus = iota
	DStatusNotFound
	DStatusError
)

func (d *DownloadHelper) GetDownloadStatus(trackId string) (models.DownloadEntry, GetDownloadStatusStatus) {
	entry, ok := d.downloadEntries.Load(trackId)
	if !ok {
		return models.DownloadEntry{}, DStatusNotFound
	}

	if entry.Status == models.DownloadInProgress {
		filepathTmp := entry.FilepathNoExt + ".tmp"
		file, err := os.Open(filepathTmp)
		if err != nil {
			log.Printf("Error opening file %s: %s", filepathTmp, err)
			return models.DownloadEntry{}, DStatusError
		}
		defer file.Close()

		fi, err := file.Stat()
		if err != nil {
			log.Printf("Error stating file %s: %s", filepathTmp, err)
			return models.DownloadEntry{}, DStatusError
		}

		entry.DownloadedBytes = int(fi.Size())
	}

	return entry, DStatusOk
}

type CancelDownloadStatus int

const (
	DCancelOk CancelDownloadStatus = iota
	DCancelNotFound
	DCancelNotInProgress
)

func (d *DownloadHelper) CancelDownload(trackId string) CancelDownloadStatus {
	entry, ok := d.downloadEntries.Load(trackId)
	if !ok {
		return DCancelNotFound
	}
	switch entry.Status {
	case models.DownloadInProgress:
		entry.CancellationFunc()
		return DCancelOk
	default:
		return DCancelNotInProgress
	}
}
