package downloader

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"spotify-downloader/models"
	"strconv"
)

type DownloadHelper struct {
	// filepath -> models.DownloadEntry
	DownloadEntries downloadEntrySyncMap
}

type DownloadStartStatus int

const (
	StartedDownloading DownloadStartStatus = iota
	ErrorCreatingFile
	ErrorSendingRequest
	ErrorReadingContentLength
)

type readerWithCancellationFunc func(p []byte) (n int, err error)

func (rf readerWithCancellationFunc) Read(p []byte) (n int, err error) {
	return rf(p)
}

func (d *DownloadHelper) StartDownload(downloadPath, url string) DownloadStartStatus {
	ch := make(chan DownloadStartStatus)

	go func() {
		cancelled := false
		func() {
			out, err := os.Create(downloadPath)
			if err != nil {
				ch <- ErrorCreatingFile
				return
			}
			defer out.Close()

			headResp, err := http.Head(url)
			if err != nil {
				ch <- ErrorSendingRequest
				return
			}
			defer headResp.Body.Close()

			size, err := strconv.Atoi(headResp.Header.Get("Content-Length"))
			if err != nil {
				ch <- ErrorReadingContentLength
				return
			}

			resp, err := http.Get(url)
			if err != nil {
				ch <- ErrorSendingRequest
				return
			}
			defer resp.Body.Close()

			log.Println("download at", downloadPath, "started")
			ch <- StartedDownloading

			ctx, fn := context.WithCancel(context.Background())

			d.DownloadEntries.Store(
				downloadPath,
				models.DownloadEntry{
					Filepath:         downloadPath,
					YoutubeLink:      url,
					TotalBytes:       size,
					Status:           models.DownloadInProgress,
					CancellationFunc: fn,
				},
			)

			_, err = io.Copy(out, readerWithCancellationFunc(func(p []byte) (n int, err error) {
				select {
				case <-ctx.Done():
					cancelled = true
					return 0, io.EOF
				default:
					return resp.Body.Read(p)
				}
			}))

			entry, ok := d.DownloadEntries.Load(downloadPath)
			if !ok {
				entry = models.DownloadEntry{
					Filepath:    downloadPath,
					YoutubeLink: url,
					TotalBytes:  size,
				}
			}

			switch {
			case cancelled:
				log.Println("download at", downloadPath, "was cancelled")
				entry.Status = models.DownloadedCancelled
			case err != nil:
				log.Println("download at", downloadPath, "failed:", err)
				entry.Status = models.DownloadFailed
			default:
				log.Println("finished downloading at", downloadPath)
				entry.Status = models.DownloadFinished
			}

			d.DownloadEntries.Store(downloadPath, entry)
		}()

		if cancelled {
			os.Remove(downloadPath)
		}
	}()

	return <-ch
}

type GetDownloadStatusStatus int

const (
	GetDownloadStatusOk GetDownloadStatusStatus = iota
	GetDownloadStatusNotFound
	GetDownloadStatusError
)

func (d *DownloadHelper) GetDownloadStatus(path string) (models.DownloadEntry, GetDownloadStatusStatus) {
	entry, ok := d.DownloadEntries.Load(path)
	if !ok {
		return models.DownloadEntry{}, GetDownloadStatusNotFound
	}

	file, err := os.Open(path)
	if err != nil {
		log.Printf("Error opening file %s: %s", path, err)
		return models.DownloadEntry{}, GetDownloadStatusError
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		log.Printf("Error stating file %s: %s", path, err)
		return models.DownloadEntry{}, GetDownloadStatusError
	}

	entry.DownloadedBytes = int(fi.Size())
	return entry, GetDownloadStatusOk
}

type CancelDownloadStatus int

const (
	CancelDownloadStatusOk CancelDownloadStatus = iota
	CancelDownloadStatusNotFound
	CancelDownloadStatusNotInProgress
)

func (d *DownloadHelper) CancelDownload(path string) CancelDownloadStatus {
	entry, ok := d.DownloadEntries.Load(path)
	if !ok {
		return CancelDownloadStatusNotFound
	}
	switch entry.Status {
	case models.DownloadInProgress:
		entry.CancellationFunc()
		return CancelDownloadStatusOk
	default:
		return CancelDownloadStatusNotInProgress
	}
}
