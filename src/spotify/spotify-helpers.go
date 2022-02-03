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

func GetPlaylistById(id string) (Playlist, int, error) {
	req, _ := http.NewRequest("GET", "https://api.spotify.com/v1/playlists/"+id, nil)
	req.Header.Add("Authorization", token.TokenType+" "+token.AccessToken)
	req.Header.Add("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return Playlist{}, -1, err
	}

	if response.StatusCode != 200 {
		return Playlist{}, response.StatusCode, fmt.Errorf("%d %s StatusCode when queying a playlist %s", response.StatusCode, response.Status, id)
	}

	var playlist Playlist
	json.NewDecoder(response.Body).Decode(&playlist)

	return playlist, 0, nil
}
