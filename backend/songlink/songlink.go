package songlink

import (
	"encoding/json"
	"log"
	"net/http"
	"spotify-downloader/models"

	"go.uber.org/ratelimit"
)

type SonglinkHelper struct {
	Endpoint  string
	RateLimit ratelimit.Limiter
}

func (s *SonglinkHelper) SetDefaults() {
	s.Endpoint = "https://api.song.link/v1-alpha.1/links"
	s.RateLimit = ratelimit.New(10)
}

type QueryResponseStatus int

const (
	Found QueryResponseStatus = iota
	ErrorSendingRequest
	NoSongWithSuchId
	NoYoutubeLinkForSong
	TooManyRequests
)

func (s *SonglinkHelper) GetYoutubeLinkBySpotifyId(spotifyId string) (models.DownloadLink, QueryResponseStatus) {
	s.RateLimit.Take()

	type SonglinkResponse struct {
		LinksByPlatform struct {
			Youtube struct {
				Url string
			}
		}
	}

	req, _ := http.NewRequest("GET", s.Endpoint, nil)
	query := req.URL.Query()
	query.Add("platform", "spotify")
	query.Add("type", "song")
	query.Add("id", spotifyId)
	req.URL.RawQuery = query.Encode()

	response, err := http.DefaultClient.Do(req)
	// actual ERRORS with a request or connectivity
	if err != nil {
		log.Printf("error sending a request to %s: %s\n", req.URL, err)
		return models.DownloadLink{}, ErrorSendingRequest
	}

	// too many requests
	if response.StatusCode == 429 {
		return models.DownloadLink{}, TooManyRequests
	}

	// no spotify song with such id exists
	if response.StatusCode == 404 {
		return models.DownloadLink{}, NoSongWithSuchId
	}

	body := SonglinkResponse{}
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
