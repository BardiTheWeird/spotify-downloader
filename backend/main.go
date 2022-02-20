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
	appConfig                 AppConfig
	featureYoutubeDlInstalled bool
	featureFfmpegInstalled    bool
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

	_, _, err := RunCliCommand("youtube-dl", "--version")
	if err == nil {
		featureYoutubeDlInstalled = true
		fmt.Println("youtube-dl detected")
	} else {
		fmt.Println("youtube-dl could not be detected. Downloads will be unavailable")
	}
	_, _, err = RunCliCommand("ffmpeg", "-version")
	if err == nil {
		featureFfmpegInstalled = true
		fmt.Println("ffmpeg detected")
	} else {
		fmt.Println("ffmpeg could not be detected. Conversion from mp4 will not be available")
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
		rw.Write([]byte("there might be a frontend here sometime in the future"))
	})

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

	http.HandleFunc("/s2y", func(rw http.ResponseWriter, r *http.Request) {
		SetContentTypeToJson(rw)
		id, ok := GetQueryParameterOrWriteErrorResponse("id", rw, r)
		if !ok {
			return
		}

		downloadLink, statusCode := odeslii.GetYoutubeLinkBySpotifyId(id)

		switch statusCode {
		case odeslii.ErrorSendingRequest:
			rw.WriteHeader(http.StatusInternalServerError)
		case odeslii.NoSongWithSuchId:
			WriteJsonResponse(rw,
				http.StatusNotFound,
				models.CreateErrorPayload(
					404,
					fmt.Sprintf("No entry for song with id %s", id),
				),
			)
		case odeslii.NoYoutubeLinkForSong:
			WriteJsonResponse(rw,
				http.StatusNotFound,
				models.CreateErrorPayload(
					404,
					fmt.Sprintf("No YouTube link for song with id %s", id),
				),
			)
		case odeslii.Found:
			bytes, _ := json.Marshal(downloadLink)
			WriteJsonResponse(rw,
				http.StatusOK,
				bytes,
			)
		}
	})

	http.HandleFunc("/start-download", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if !featureYoutubeDlInstalled {
			rw.WriteHeader(http.StatusServiceUnavailable)
			rw.Write([]byte("youtube-dl is not installed, thus downloads are unavailable"))
		}

		filepath, ok := GetQueryParameterOrWriteErrorResponse("path", rw, r)
		if !ok {
			return
		}
		youtubeLink, ok := GetQueryParameterOrWriteErrorResponse("link", rw, r)
		if !ok {
			return
		}

		downloadLink, exists := GetYoutubeDownloadLink(youtubeLink)
		if !exists {
			rw.WriteHeader(http.StatusNotFound)
			return
		}

		status := downloader.StartDownload(
			filepath,
			downloadLink)

		switch status {
		case downloader.ErrorCreatingFile:
			WriteJsonResponse(rw,
				http.StatusBadRequest,
				models.CreateErrorPayload(
					403,
					"could not create a file at "+filepath,
				),
			)
		case downloader.ErrorSendingRequest:
			rw.WriteHeader(http.StatusInternalServerError)
		case downloader.ErrorReadingContentLength:
			WriteJsonResponse(rw,
				http.StatusBadRequest,
				models.CreateErrorPayload(
					400,
					"error reading content-length at the download link",
				),
			)
		case downloader.StartedDownloading:
			rw.WriteHeader(http.StatusNoContent)
		}
	})

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
			bytes, _ := json.Marshal(downloadEntry)
			WriteJsonResponse(rw,
				http.StatusOK,
				bytes,
			)
		default:
			rw.WriteHeader(http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/cancel-download", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		path, ok := GetQueryParameterOrWriteErrorResponse("path", rw, r)
		if !ok {
			return
		}

		switch status := downloader.CancelDownload(path); status {
		case downloader.CancelDownloadStatusNotFound:
			rw.WriteHeader(http.StatusNotFound)
		case downloader.CancelDownloadStatusNotInProgress:
			rw.WriteHeader(http.StatusConflict)
		case downloader.CancelDownloadStatusOk:
			rw.WriteHeader(http.StatusNoContent)
		default:
			rw.WriteHeader(http.StatusInternalServerError)
		}
	})

	log.Println("Starting a server at :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
