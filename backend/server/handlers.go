package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"spotify-downloader/clihelpers"
	"spotify-downloader/downloader"
	"spotify-downloader/models"
	"spotify-downloader/songlink"
	"spotify-downloader/spotify"
)

// OPTIONS /
func (s *Server) handleOptions() func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		rw.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Length, Authorization")
	}
}

// "/playlist"
func (s *Server) handlePlaylist() func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		id, ok := GetQueryParameterOrWriteErrorResponse("id", rw, r)
		if !ok {
			return
		}

		playlist, status := s.SpotifyHelper.GetPlaylistById(id)

		switch status {
		case spotify.ErrorSendingRequest, spotify.UnexpectedResponseStatus:
			rw.WriteHeader(http.StatusInternalServerError)
		case spotify.BadOrExpiredToken, spotify.BadOAuth:
			rw.WriteHeader(http.StatusUnauthorized)
		case spotify.ExceededRateLimits:
			rw.WriteHeader(http.StatusTooManyRequests)
		case spotify.NotFound:
			rw.WriteHeader(http.StatusNotFound)
		case spotify.BadClientCredentials:
			WriteJsonResponse(rw,
				http.StatusHTTPVersionNotSupported,
				models.CreateErrorPayload("bad spotify client id or key"))
		case spotify.Ok:
			bytes, _ := json.Marshal(playlist)
			WriteJsonResponse(rw, http.StatusOK, bytes)
		}
	}
}

// "/s2y"
func (s *Server) handleS2Y() func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		SetContentTypeToJson(rw)
		id, ok := GetQueryParameterOrWriteErrorResponse("id", rw, r)
		if !ok {
			return
		}

		downloadLink, statusCode := s.SonglinkHelper.GetYoutubeLinkBySpotifyId(id)

		switch statusCode {
		case songlink.ErrorSendingRequest:
			rw.WriteHeader(http.StatusInternalServerError)
		case songlink.TooManyRequests:
			rw.WriteHeader(http.StatusTooManyRequests)
		case songlink.NoSongWithSuchId:
			WriteJsonResponse(rw,
				http.StatusNotFound,
				models.CreateErrorPayload(
					fmt.Sprintf("No entry for song with id %s", id),
				),
			)
		case songlink.NoYoutubeLinkForSong:
			WriteJsonResponse(rw,
				http.StatusNotFound,
				models.CreateErrorPayload(
					fmt.Sprintf("No YouTube link for song with id %s", id),
				),
			)
		case songlink.Found:
			bytes, _ := json.Marshal(downloadLink)
			WriteJsonResponse(rw,
				http.StatusOK,
				bytes,
			)
		}
	}
}

// "/download/start"
func (s *Server) handleDownloadStart() func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		// if !s.FeatureYoutubeDlInstalled {
		// 	rw.WriteHeader(http.StatusServiceUnavailable)
		// 	rw.Write([]byte("youtube-dl is not installed, thus downloads are unavailable"))
		// }

		filepath, ok := GetQueryParameterOrWriteErrorResponse("path", rw, r)
		if !ok {
			return
		}
		youtubeLink, ok := GetQueryParameterOrWriteErrorResponse("link", rw, r)
		if !ok {
			return
		}

		downloadLink, exists := clihelpers.GetYoutubeDownloadLink(youtubeLink)
		if !exists {
			rw.WriteHeader(http.StatusNotFound)
			return
		}

		status := s.DownloadHelper.StartDownload(
			filepath,
			downloadLink)

		switch status {
		case downloader.DStartErrorCreatingFile:
			WriteJsonResponse(rw,
				http.StatusBadRequest,
				models.CreateErrorPayloadWithCode(
					403,
					"could not create a file at "+filepath,
				),
			)
		case downloader.DStartErrorSendingRequest:
			rw.WriteHeader(http.StatusInternalServerError)
		case downloader.DStartErrorReadingContentLength:
			WriteJsonResponse(rw,
				http.StatusBadRequest,
				models.CreateErrorPayloadWithCode(
					400,
					"error reading content-length at the download link",
				),
			)
		case downloader.DStartOk:
			rw.WriteHeader(http.StatusNoContent)
		}
	}
}

// "/download/status"
func (s *Server) handleDownloadStatus() func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		path, ok := GetQueryParameterOrWriteErrorResponse("path", rw, r)
		if !ok {
			return
		}

		downloadEntry, responseStatus := s.DownloadHelper.GetDownloadStatus(path)
		switch responseStatus {
		case downloader.DStatusFound:
			rw.WriteHeader(http.StatusNotFound)
		case downloader.DStatusOk:
			bytes, _ := json.Marshal(downloadEntry)
			WriteJsonResponse(rw, http.StatusOK, bytes)
		default:
			rw.WriteHeader(http.StatusInternalServerError)
		}
	}
}

// "/download/cancel"
func (s *Server) handleDownloadCancel() func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		path, ok := GetQueryParameterOrWriteErrorResponse("path", rw, r)
		if !ok {
			return
		}

		switch status := s.DownloadHelper.CancelDownload(path); status {
		case downloader.DCancelNotFound:
			rw.WriteHeader(http.StatusNotFound)
		case downloader.DCancelNotInProgress:
			rw.WriteHeader(http.StatusConflict)
		case downloader.DCancelOk:
			rw.WriteHeader(http.StatusNoContent)
		default:
			rw.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func (s *Server) handleSpotifyConfigure() func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		configuration := struct {
			ClientId     string `json:"client_id"`
			ClientSecret string `json:"client_secret"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&configuration)
		if err != nil {
			log.Panicln("error decoding client configuration")
		}

		if len(configuration.ClientId) == 0 {
			WriteJsonResponse(rw, 400,
				models.CreateErrorPayload("client_id is empty"))
			return
		}
		if len(configuration.ClientSecret) == 0 {
			WriteJsonResponse(rw, 400,
				models.CreateErrorPayload("client_secret is empty"))
			return
		}

		oldClientId := s.SpotifyHelper.ClientId
		oldClientSecret := s.SpotifyHelper.ClientSecret
		s.SpotifyHelper.ClientId = configuration.ClientId
		s.SpotifyHelper.ClientSecret = configuration.ClientSecret

		if !s.SpotifyHelper.UpdateClientToken() {
			WriteJsonResponse(rw, 400,
				models.CreateErrorPayloadWithCode(401,
					"can't authenticate using new credentials"))

			s.SpotifyHelper.ClientId = oldClientId
			s.SpotifyHelper.ClientSecret = oldClientSecret
			s.SpotifyHelper.UpdateClientToken()
			return
		}
		rw.WriteHeader(http.StatusNoContent)
	}
}
