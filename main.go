package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"spotify-downloader/downloader"
	"spotify-downloader/models"
	"spotify-downloader/odeslii"
	"spotify-downloader/spotify"
	"strings"
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

func SetContentTypeToJson(rw http.ResponseWriter) {
	rw.Header().Add("Content-Type", "application/json")
}

func WriteJsonResponse(rw http.ResponseWriter, statusCode int, payload []byte) {
	SetContentTypeToJson(rw)
	rw.WriteHeader(statusCode)
	rw.Write(payload)
}

func GetQueryParameterOrWriteErrorResponse(parameter string, rw http.ResponseWriter, r *http.Request) (string, bool) {
	val := r.URL.Query().Get(parameter)
	present := true
	if len(val) == 0 {
		WriteJsonResponse(
			rw,
			400,
			models.CreateErrorPayload(
				400,
				fmt.Sprintf("'%s' query parameter is missing", parameter),
			),
		)
		present = false
	}
	return val, present
}

func RunCliCommand(name string, params ...string) (string, string, error) {
	// get the download link with the power of youtube-dl
	cmd := exec.Command(name, params...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func GetYoutubeDownloadLink(youtubeLink string) (string, bool) {
	link, _, err := RunCliCommand("youtube-dl", "-x", "-g", youtubeLink)
	exists := true
	if err != nil {
		exists = false
		fmt.Println("error querying youtube-dl:", err)
	}

	return strings.TrimSpace(link), exists
}

func main() {
	configureApp()
	spotify.Authenticate(appConfig.GetB64())

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("there might be a list of all the endpoints here sometime in the future"))
	})

	// /playlist?id={spotify_playlist_id}
	// 200 + playlist payload
	// 400 => "id" is empty
	// 401 => not authorized (maybe?)
	// 404 => no playlist with such id
	// 429 => too many requests
	// 500
	http.HandleFunc("/playlist", func(rw http.ResponseWriter, r *http.Request) {
		id, ok := GetQueryParameterOrWriteErrorResponse("id", rw, r)
		if !ok {
			return
		}

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
			WriteJsonResponse(rw, http.StatusOK, bytes)
		}
	})

	// /s2y?id={spotify_song_id}
	// 200 + songToDownload payload
	// 400 => 'id' is empty
	// 404 => (no such id / no yt link) + error payload:
	//     400 => no entry for song with {id}
	//     404 => no YouTube link for song with {id}
	// 500
	http.HandleFunc("/s2y", func(rw http.ResponseWriter, r *http.Request) {
		SetContentTypeToJson(rw)
		id, ok := GetQueryParameterOrWriteErrorResponse("id", rw, r)
		if !ok {
			return
		}

		songToDownload, statusCode := odeslii.GetYoutubeLinkBySpotifyId(id)

		switch statusCode {
		case odeslii.ErrorSendingRequest:
			rw.WriteHeader(http.StatusInternalServerError)
		case odeslii.NoSongWithSuchId:
			rw.WriteHeader(http.StatusNotFound)
			rw.Write(models.CreateErrorPayload(400, fmt.Sprintf("No entry for song with id %s", id)))
		case odeslii.NoYoutubeLinkForSong:
			rw.WriteHeader(http.StatusNotFound)
			rw.Write(models.CreateErrorPayload(404, fmt.Sprintf("No YouTube link for song with id %s", id)))
		case odeslii.Found:
			rw.WriteHeader(http.StatusOK)
			bytes, _ := json.Marshal(songToDownload)
			rw.Write(bytes)
		}
	})

	// 400 + error payload
	//     400 => error decoding body
	//     403 => error creating a file
	// 404 => youtube-dl couldn't find a download link
	// 405 => only allows POST
	// 500 => youtube-dl execution error
	http.HandleFunc("/start-download", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var downloadRequest models.DownloadRequest
		err := json.NewDecoder(r.Body).Decode(&downloadRequest)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			SetContentTypeToJson(rw)
			rw.Write(models.CreateErrorPayload(400, "error decoding the body"))
			return
		}

		if len(downloadRequest.Filepath) == 0 {
			WriteJsonResponse(
				rw,
				http.StatusBadRequest,
				models.CreateErrorPayload(
					400,
					"filepath is empty",
				),
			)
			return
		}

		if len(downloadRequest.YoutubeLink) == 0 {
			WriteJsonResponse(
				rw,
				http.StatusBadRequest,
				models.CreateErrorPayload(
					400,
					"youtube_link is empty",
				),
			)
			return
		}

		// get the download link with the power of youtube-dl
		downloadLink, exists := GetYoutubeDownloadLink(downloadRequest.YoutubeLink)
		if !exists {
			rw.WriteHeader(http.StatusNotFound)
			return
		}

		status := downloader.StartDownload(
			downloadRequest.Filepath,
			downloadLink)

		switch status {
		case downloader.ErrorCreatingFile:
			rw.WriteHeader(http.StatusBadRequest)
			SetContentTypeToJson(rw)
			rw.Write(models.CreateErrorPayload(
				403,
				fmt.Sprint("could not create a file at", downloadRequest.Filepath)))
		case downloader.ErrorSendingRequest:
			rw.WriteHeader(http.StatusInternalServerError)
		case downloader.ErrorReadingContentLength:
			rw.WriteHeader(http.StatusBadRequest)
			SetContentTypeToJson(rw)
			rw.Write(models.CreateErrorPayload(
				400,
				"error reading content-length at the download link"))
		case downloader.StartedDownloading:
			rw.WriteHeader(http.StatusOK)
		}
	})

	// /download-status?path={path}
	// 200
	// 400 + payload error => path not provided
	// 404 => no downloads at path
	// 500 => can't stat file OR unhandled GetDownloadStatusStatus
	http.HandleFunc("/download-status", func(rw http.ResponseWriter, r *http.Request) {
		path, ok := GetQueryParameterOrWriteErrorResponse("path", rw, r)
		if !ok {
			return
		}

		downloadEntry, responseStatus := downloader.GetDownloadStatus(path)
		switch responseStatus {
		case downloader.GetDownloadStatusNotFound:
			rw.WriteHeader(http.StatusNotFound)
		case downloader.GetDownloadStatusGetDownloadedError:
			rw.WriteHeader(http.StatusInternalServerError)
		case downloader.GetDownloadStatusOk:
			bytes, err := json.Marshal(downloadEntry)
			if err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}

			rw.WriteHeader(http.StatusOK)
			SetContentTypeToJson(rw)
			rw.Write(bytes)
		default:
			rw.WriteHeader(http.StatusInternalServerError)
		}
	})

	log.Println("Starting a server at :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
