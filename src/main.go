package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"spotify-downloader/src/models"
	"spotify-downloader/src/odeslii"
	"spotify-downloader/src/spotify"
)

var (
	appConfig AppConfig
)

type AppConfig struct {
	ClientId     string
	ClientSecret string
}

func B64Strict(s string) string {
	return base64.RawStdEncoding.Strict().EncodeToString([]byte(s))
}

func (config AppConfig) GetB64() string {
	return B64Strict(config.ClientId + ":" + config.ClientSecret)
}

func getEnvOrDefault(envKey, def string) string {
	val, err := os.LookupEnv(envKey)
	if !err {
		fmt.Printf("%s is not set, using %s\n", envKey, def)
		val = def
	}
	return val
}

func configureApp() {
	appConfig = AppConfig{
		ClientId:     getEnvOrDefault("CLIENT_ID", "00000000000000000000000000000000"),
		ClientSecret: getEnvOrDefault("CLIENT_SECRET", "00000000000000000000000000000000"),
	}
}

func SetContentTypeToJson(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
}

func main() {
	configureApp()
	spotify.Authenticate(appConfig.GetB64())

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("there might be a list of all the endpoints here sometime in the future"))
	})

	// 500
	// 401 => not authorized (maybe?)
	// 404 => no playlist with such id
	// 429 => too many requests
	// 200 + playlist payload
	http.HandleFunc("/playlist/", func(rw http.ResponseWriter, r *http.Request) {
		SetContentTypeToJson(rw, r)
		id := r.URL.Path[len("/playlist/"):]

		spotifyPlaylist, status := spotify.GetPlaylistById(id)
		if status == spotify.BadOrExpiredToken {
			spotifyPlaylist, status = spotify.GetPlaylistById(id)
		}

		switch status {
		case spotify.ErrorSendingRequest, spotify.UnexpectedResponseStatus:
			rw.WriteHeader(http.StatusInternalServerError)
		case spotify.BadOrExpiredToken, spotify.BadOAuth:
			rw.WriteHeader(http.StatusUnauthorized)
		case spotify.ExceededRateLimits:
			rw.WriteHeader(http.StatusTooManyRequests)
		case spotify.NotFound:
			rw.WriteHeader(http.StatusNotFound)
		case spotify.Ok:
			playlist := models.FromSpotifyPlaylist(spotifyPlaylist)
			bytes, _ := json.Marshal(playlist)
			SetContentTypeToJson(rw, r)
			rw.Write(bytes)
		}
	})

	// 500
	// 404 => (no such id / no yt link) -> error payload
	// 200 + songToDownload payload
	http.HandleFunc("/s2y/", func(rw http.ResponseWriter, r *http.Request) {
		SetContentTypeToJson(rw, r)

		spotifyId := r.URL.Path[len("/s2y/"):]

		songToDownload, statusCode := odeslii.GetYoutubeLinkBySpotifyId(spotifyId)

		switch statusCode {
		case odeslii.ErrorSendingRequest:
			rw.WriteHeader(http.StatusInternalServerError)
		case odeslii.NoSongWithSuchId:
			rw.WriteHeader(http.StatusNotFound)
			rw.Write(models.CreateErrorPayload(400, fmt.Sprintf("No entry for song with id %s", spotifyId)))
		case odeslii.NoYoutubeLinkForSong:
			rw.WriteHeader(http.StatusNotFound)
			rw.Write(models.CreateErrorPayload(404, fmt.Sprintf("No YouTube link for song with id %s", spotifyId)))
		case odeslii.Found:
			rw.WriteHeader(http.StatusOK)
			bytes, _ := json.Marshal(songToDownload)
			rw.Write(bytes)
		}
	})

	log.Println("Starting a server at :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
