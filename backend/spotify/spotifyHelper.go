package spotify

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"spotify-downloader/models"
	"time"
)

type SpotifyHelper struct {
	ClientId     string
	ClientSecret string

	Token struct {
		Value     string
		ExpiresIn time.Duration
		Timestamp time.Time
	}
}

func (s *SpotifyHelper) UseClientAuthentication(r *http.Request) bool {
	if len(s.Token.Value) == 0 ||
		time.Now().After(s.Token.Timestamp.Add(time.Second*s.Token.ExpiresIn)) {

		if !s.UpdateClientToken() {
			return false
		}
		log.Println("Spotify authentication token was refreshed")
	}
	r.Header.Add("Authorization", s.Token.Value)
	return true
}

func (s *SpotifyHelper) UpdateClientToken() bool {
	credentialsB64 := base64.RawStdEncoding.Strict().
		EncodeToString([]byte(s.ClientId + ":" + s.ClientSecret))
	req, _ := http.NewRequest("POST", "https://accounts.spotify.com/api/token?grant_type=client_credentials", nil)
	req.Header.Add("Authorization", "Basic "+credentialsB64)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error doing a request %v: %s\n", req, err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		return false
	}

	token := struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}{}
	err = json.NewDecoder(response.Body).Decode(&token)
	if err != nil {
		log.Panicln("Error decoding client token:", err)
	}

	s.Token.Value = token.TokenType + " " + token.AccessToken
	s.Token.ExpiresIn = time.Second * time.Duration(token.ExpiresIn)
	s.Token.Timestamp = time.Now()
	log.Println("spotify client token refreshed")
	return true
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

func (s *SpotifyHelper) GetPlaylistById(id string) (models.Playlist, GetPlaylistResponseStatus) {
	req, _ := http.NewRequest("GET", "https://api.spotify.com/v1/playlists/"+id, nil)
	req.Header.Add("Content-Type", "application/json")
	if !s.UseClientAuthentication(req) {
		return models.Playlist{}, BadClientCredentials
	}
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
		var playlist playlist
		json.NewDecoder(response.Body).Decode(&playlist)
		return playlist.toModelsPlaylist(), Ok
	default:
		log.Printf("Response status %d was not expected\n", response.StatusCode)

		payload, _ := ioutil.ReadAll(response.Body)
		log.Println("payload:", string(payload))
		return models.Playlist{}, UnexpectedResponseStatus
	}
}
