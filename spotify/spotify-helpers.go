package spotify

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var token ClientToken

func Authenticate(credentialsB64 string) {
	req, _ := http.NewRequest("POST", "https://accounts.spotify.com/api/token?grant_type=client_credentials", nil)
	req.Header.Add("Authorization", "Basic "+credentialsB64)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Error doing a request %v: %s\n", req, err)
	}

	// TO DO: check for bad error statuses

	err = json.NewDecoder(response.Body).Decode(&token)
	if err != nil {
		fmt.Printf("Error decoding client token: %s", err)
	}

	defer response.Body.Close()
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

func GetPlaylistById(id string) (Playlist, GetPlaylistResposeStatus) {
	req, _ := http.NewRequest("GET", "https://api.spotify.com/v1/playlists/"+id, nil)
	req.Header.Add("Authorization", token.TokenType+" "+token.AccessToken)
	req.Header.Add("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(req)
	// actual error with a request or connectivity
	if err != nil {
		fmt.Printf("error sending a request to %s: %s\n", req.URL, err)
		return Playlist{}, ErrorSendingRequest
	}

	switch response.StatusCode {
	case 401:
		return Playlist{}, BadOrExpiredToken
	case 403:
		return Playlist{}, BadOAuth
	case 404:
		return Playlist{}, NotFound
	case 429:
		return Playlist{}, ExceededRateLimits
	case 200:
		// bytes, _ := ioutil.ReadAll(response.Body)
		// fmt.Println("Response from Spotify:", string(bytes))

		var playlist Playlist
		// json.Unmarshal(bytes, &playlist)
		json.NewDecoder(response.Body).Decode(&playlist)
		return playlist, Ok
	default:
		fmt.Printf("Response status %d was not expected\n", response.StatusCode)
		return Playlist{}, UnexpectedResponseStatus
	}
}
