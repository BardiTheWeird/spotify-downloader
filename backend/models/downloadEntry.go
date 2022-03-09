package models

import "context"

type DownloadStatus int

const (
	DownloadInProgress DownloadStatus = iota
	DownloadConvertationInProgress
	DownloadFinished
	DownloadErrorConverting
	DownloadFailed
	DownloadedCancelled
)

type DownloadEntry struct {
	TotalBytes      int            `json:"total_bytes"`
	DownloadedBytes int            `json:"downloaded_bytes,omitempty"`
	Status          DownloadStatus `json:"status"`

	CancellationFunc context.CancelFunc `json:"-"`
}
