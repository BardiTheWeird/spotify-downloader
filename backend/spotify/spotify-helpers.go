package spotify

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"spotify-downloader/models"
)

var tokenString string

func Authenticate(credentialsB64 string) {
	req, _ := http.NewRequest("POST", "https://accounts.spotify.com/api/token?grant_type=client_credentials", nil)
	req.Header.Add("Authorization", "Basic "+credentialsB64)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Error doing a request %v: %s\n", req, err)
	}
	defer response.Body.Close()

	// TO DO: check for bad error statuses
	token := struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}{}
	err = json.NewDecoder(response.Body).Decode(&token)
	if err != nil {
		fmt.Printf("Error decoding client token: %s", err)
		return
	}

	tokenString = token.TokenType + " " + token.AccessToken
}

type GetPlaylistResposeStatus int

const (
	Ok GetPlaylistResposeStatus = iota
	ErrorSendingRequest
	BadOrExpiredToken
	BadOAuth
	ExceededRateLimits
	NotFound
	UnexpectedResponseStatus
)

func GetPlaylistById(id string) (models.Playlist, GetPlaylistResposeStatus) {
	req, _ := http.NewRequest("GET", "https://api.spotify.com/v1/playlists/"+id, nil)
	req.Header.Add("Authorization", tokenString)
	req.Header.Add("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(req)
	// actual error with a request or connectivity
	if err != nil {
		fmt.Printf("error sending a request to %s: %s\n", req.URL, err)
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
		// bytes, _ := ioutil.ReadAll(response.Body)
		// fmt.Println("Response from Spotify:", string(bytes))

		var playlist playlist
		// json.Unmarshal(bytes, &playlist)
		json.NewDecoder(response.Body).Decode(&playlist)
		return playlist.toModelsPlaylist(), Ok
	default:
		fmt.Printf("Response status %d was not expected\n", response.StatusCode)

		payload, _ := ioutil.ReadAll(response.Body)
		fmt.Println("payload:", string(payload))
		return models.Playlist{}, UnexpectedResponseStatus
	}
}
