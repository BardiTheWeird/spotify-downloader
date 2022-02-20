package models

type SongToDownload struct {
	SpotifyId   string `json:"spotify_id"`
	YoutubeLink string `json:"youtube_url"`
}
