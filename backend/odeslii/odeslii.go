package odeslii

import (
	"encoding/json"
	"fmt"
	"net/http"
	"spotify-downloader/models"
)

var endpoint string = "https://api.song.link/v1-alpha.1/links"

type YoutubeLinks struct {
	Url string
}

type LinksByPlatform struct {
	Youtube YoutubeLinks
}

type OdesliiResponse struct {
	LinksByPlatform LinksByPlatform
}

type QueryResponseStatus int

const (
	Found QueryResponseStatus = iota
	ErrorSendingRequest
	NoSongWithSuchId
	NoYoutubeLinkForSong
)

func GetYoutubeLinkBySpotifyId(spotifyId string) (models.DownloadLink, QueryResponseStatus) {
	req, _ := http.NewRequest("GET", endpoint, nil)
	query := req.URL.Query()
	query.Add("platform", "spotify")
	query.Add("type", "song")
	query.Add("id", spotifyId)
	req.URL.RawQuery = query.Encode()

	response, err := http.DefaultClient.Do(req)
	// actual ERRORS with a request or connectivity
	if err != nil {
		fmt.Printf("error sending a request to %s: %s\n", req.URL, err)
		return models.DownloadLink{}, ErrorSendingRequest
	}

	// no spotify song with such id exists
	if response.StatusCode == 404 {
		return models.DownloadLink{}, NoSongWithSuchId
	}

	body := OdesliiResponse{}
	json.NewDecoder(response.Body).Decode(&body)

	youtubeLink := body.LinksByPlatform.Youtube.Url
	// this song couldn't be found on YouTube
	if len(youtubeLink) == 0 {
		return models.DownloadLink{}, NoYoutubeLinkForSong

	}
	// actually found a song
	return models.DownloadLink{
			SpotifyId: spotifyId,
			Link:      youtubeLink,
		},
		Found
}
