package spotify

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"spotify-downloader/models"
)

type SpotifyHelper struct {
}

type GetPlaylistResponseStatus int

const (
	Ok GetPlaylistResponseStatus = iota
	ErrorSendingRequest
	BadOrExpiredToken
	BadClientCredentials
	BadOAuth
	ExceededRateLimits
	NotFound
	UnexpectedResponseStatus
)

func (s *SpotifyHelper) GetPlaylistById(id, linkType, accessToken string) (models.Playlist, GetPlaylistResponseStatus) {
	req, _ := http.NewRequest("GET", "https://api.spotify.com/v1/"+linkType+"/"+id, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", accessToken)
	http.DefaultClient.Do(req)

	response, err := http.DefaultClient.Do(req)
	// actual error with a request or connectivity
	if err != nil {
		log.Printf("error sending a request to %s: %s\n", req.URL, err)
		return models.Playlist{}, ErrorSendingRequest
	}

	switch response.StatusCode {
	case 400, 401:
		return models.Playlist{}, BadOrExpiredToken
	case 403:
		return models.Playlist{}, BadOAuth
	case 404:
		return models.Playlist{}, NotFound
	case 429:
		return models.Playlist{}, ExceededRateLimits
	case 200:
		var modelsPlaylist models.Playlist
		switch linkType {
		case "playlists":
			var playlist playlistTracks
			json.NewDecoder(response.Body).Decode(&playlist)
			modelsPlaylist = toModelsPlaylist(playlist.toTracks())
		case "albums":
			var album albumTracks
			json.NewDecoder(response.Body).Decode(&album)
			modelsPlaylist = toModelsPlaylist(album.toTracks())
		case "tracks":
			var trackVariable track
			json.NewDecoder(response.Body).Decode(&trackVariable)
			modelsPlaylist = toModelsPlaylist([]track{trackVariable})
		default:
			return models.Playlist{}, NotFound
		}
		return modelsPlaylist, Ok

	default:
		log.Printf("Response status %d was not expected\n", response.StatusCode)

		payload, _ := ioutil.ReadAll(response.Body)
		log.Println("payload:", string(payload))
		return models.Playlist{}, UnexpectedResponseStatus
	}
}
