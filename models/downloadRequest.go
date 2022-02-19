package models

type DownloadRequest struct {
	YoutubeLink string `json:"youtube_link"`
	Filepath    string `json:"filepath"`
}
