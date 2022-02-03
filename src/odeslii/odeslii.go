package odeslii

import (
	"encoding/json"
	"fmt"
	"net/http"
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

func GetYoutubeLinkBySpotifyId(id string) (string, bool, error) {
	req, _ := http.NewRequest("GET", endpoint, nil)
	query := req.URL.Query()
	query.Add("platform", "spotify")
	query.Add("type", "song")
	query.Add("id", id)
	req.URL.RawQuery = query.Encode()

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		// actual ERRORS with a request or connectivity
		return "", false, fmt.Errorf("error sending a request to %s: %s", req.URL, err)
	}

	if response.StatusCode >= 400 {
		// bad request formatting or too many requests
		return "", false, fmt.Errorf("%d %s when sending a request to %s", response.StatusCode, response.Status, req.URL)
	}

	body := OdesliiResponse{}
	json.NewDecoder(response.Body).Decode(&body)

	youtubeLink := body.LinksByPlatform.Youtube.Url
	if len(youtubeLink) > 0 {
		// there is a link, so here's a link
		return youtubeLink, true, nil
	} else {
		// there's no link
		return youtubeLink, false, nil
	}
}
