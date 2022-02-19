package downloader

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"spotify-downloader/models"
	"strconv"
	"sync"
	"time"
)

// filepath -> models.DownloadEntry
var DownloadEntries sync.Map

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
		fmt.Printf("Error opening file %s: %s", path, err)
		return 0, getDownloadedCantOpenFile
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		fmt.Printf("Error stating file %s: %s", path, err)
		return 0, getDownloadedCantOpenFile
	}

	return fi.Size(), getDownloadedOk
}

func PrintDownloadPercent(done chan int64, path string, total int64) {

	var stop bool = false

	for {
		select {
		case <-done:
			stop = true
		default:

			size, status := getDownloadedBytes(path)
			if status == getDownloadedCantOpenFile || status == getDownloadedCantStatFile {
				log.Fatal()
			}

			if size == 0 {
				size = 1
			}

			var percent float64 = float64(size) / float64(total) * 100

			fmt.Printf("%.0f", percent)
			fmt.Println("%")
		}

		if stop {
			break
		}

		time.Sleep(time.Second)
	}
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

		fmt.Println("content length:", size)

		resp, err := http.Get(url)
		if err != nil {
			ch <- ErrorSendingRequest
			return
		}
		defer resp.Body.Close()

		ch <- StartedDownloading

		DownloadEntries.Store(
			downloadPath,
			models.DownloadEntry{
				Filepath:    downloadPath,
				YoutubeLink: url,
				TotalBytes:  size,
				Status:      models.DownloadInProgress,
			},
		)

		n, err := io.Copy(out, resp.Body)
		fmt.Println("finished copying. n:", n, "err:", err)

		entryInterface, ok := DownloadEntries.Load(downloadPath)
		if !ok {
			entryInterface = models.DownloadEntry{
				Filepath:    downloadPath,
				YoutubeLink: url,
				TotalBytes:  size,
			}
		}

		downloadEntry := entryInterface.(models.DownloadEntry)

		if err != nil {
			downloadEntry.Status = models.DownloadFailed
		} else {
			downloadEntry.Status = models.DownloadFinished
		}

		DownloadEntries.Store(downloadPath, downloadEntry)
	}()

	return <-ch
}

func DownloadFile(url, downloadPath string) {

	file := path.Base(url)

	log.Printf("Downloading file %s from %s\n", file, url)

	var path bytes.Buffer
	path.WriteString(downloadPath)

	start := time.Now()

	out, err := os.Create(path.String())

	if err != nil {
		fmt.Println(path.String())
		panic(err)
	}

	defer out.Close()

	headResp, err := http.Head(url)

	if err != nil {
		panic(err)
	}

	defer headResp.Body.Close()

	size, err := strconv.Atoi(headResp.Header.Get("Content-Length"))

	if err != nil {
		panic(err)
	}

	done := make(chan int64)

	go PrintDownloadPercent(done, path.String(), int64(size))

	resp, err := http.Get(url)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	n, err := io.Copy(out, resp.Body)

	if err != nil {
		panic(err)
	}

	done <- n

	elapsed := time.Since(start)
	log.Printf("Download completed in %s", elapsed)
}
