package downloader

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"spotify-downloader/models"
	"strconv"
	"sync"
)

// filepath -> models.DownloadEntry
var DownloadEntries sync.Map

type CancelDownloadStatus int

const (
	CancelDownloadStatusOk CancelDownloadStatus = iota
	CancelDownloadStatusNotFound
	CancelDownloadStatusNotInProgress
)

func CancelDownload(path string) CancelDownloadStatus {
	entry, ok := DownloadEntries.Load(path)
	if !ok {
		return CancelDownloadStatusNotFound
	}
	switch entry := entry.(models.DownloadEntry); entry.Status {
	case models.DownloadInProgress:
		entry.CancellationFunc()
		return CancelDownloadStatusOk
	default:
		return CancelDownloadStatusNotInProgress
	}
}

type GetDownloadStatusStatus int

const (
	GetDownloadStatusOk GetDownloadStatusStatus = iota
	GetDownloadStatusNotFound
	GetDownloadStatusGetDownloadedError
	GetDownloadStatusUnknownError
)

func GetDownloadStatus(path string) (models.DownloadEntry, GetDownloadStatusStatus) {
	entryInterface, ok := DownloadEntries.Load(path)
	if !ok {
		return models.DownloadEntry{}, GetDownloadStatusNotFound
	}

	switch downloadedBytes, status := getDownloadedBytes(path); status {
	case getDownloadedCantOpenFile:
		return models.DownloadEntry{}, GetDownloadStatusNotFound
	case getDownloadedCantStatFile:
		return models.DownloadEntry{}, GetDownloadStatusGetDownloadedError
	case getDownloadedOk:
		entry := entryInterface.(models.DownloadEntry)
		entry.DownloadedBytes = int(downloadedBytes)
		return entry, GetDownloadStatusOk
	default:
		return models.DownloadEntry{}, GetDownloadStatusUnknownError
	}
}

type getDownloadedBytesResponseStatus int

const (
	getDownloadedOk getDownloadedBytesResponseStatus = iota
	getDownloadedCantOpenFile
	getDownloadedCantStatFile
)

func getDownloadedBytes(path string) (int64, getDownloadedBytesResponseStatus) {
	file, err := os.Open(path)
	if err != nil {
		log.Printf("Error opening file %s: %s", path, err)
		return 0, getDownloadedCantOpenFile
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		log.Printf("Error stating file %s: %s", path, err)
		return 0, getDownloadedCantOpenFile
	}

	return fi.Size(), getDownloadedOk
}

type readerWithCancellationFunc func(p []byte) (n int, err error)

func (rf readerWithCancellationFunc) Read(p []byte) (n int, err error) {
	return rf(p)
}

type DownloadStartStatus int

const (
	StartedDownloading DownloadStartStatus = iota
	ErrorCreatingFile
	ErrorSendingRequest
	ErrorReadingContentLength
)

func StartDownload(downloadPath, url string) DownloadStartStatus {
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

			DownloadEntries.Store(
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

			entryInterface, ok := DownloadEntries.Load(downloadPath)
			if !ok {
				entryInterface = models.DownloadEntry{
					Filepath:    downloadPath,
					YoutubeLink: url,
					TotalBytes:  size,
				}
			}

			downloadEntry := entryInterface.(models.DownloadEntry)
			switch {
			case cancelled:
				log.Println("download at", downloadPath, "was cancelled")
				downloadEntry.Status = models.DownloadedCancelled
			case err != nil:
				log.Println("download at", downloadPath, "failed:", err)
				downloadEntry.Status = models.DownloadFailed
			default:
				log.Println("finished downloading at", downloadPath)
				downloadEntry.Status = models.DownloadFinished
			}

			DownloadEntries.Store(downloadPath, downloadEntry)
		}()

		if cancelled {
			os.Remove(downloadPath)
		}
	}()

	return <-ch
}
