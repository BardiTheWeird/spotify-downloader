package models

type DownloadStatus int

const (
	DownloadInProgress DownloadStatus = iota
	DownloadFinished
	DownloadFailed
)

type DownloadEntry struct {
	Filepath        string         `json:"path"`
	YoutubeLink     string         `json:"youtube_link"`
	TotalBytes      int            `json:"total_bytes"`
	DownloadedBytes int            `json:"downloaded_bytes"`
	Status          DownloadStatus `json:"status"`
}
