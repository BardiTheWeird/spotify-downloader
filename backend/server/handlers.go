package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path"
	"spotify-downloader/clihelpers"
	"spotify-downloader/downloader"
	"spotify-downloader/models"
	"spotify-downloader/songlink"
	"spotify-downloader/spotify"
)

func (s *Server) handleSpotifyConfigure() http.HandlerFunc {
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

func (s *Server) handlePlaylist() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		id, ok := GetQueryParameter("id", r)
		if !ok {
			link, ok := GetQueryParameterOrWriteErrorResponse("link", rw, r)
			if !ok {
				return
			}
			spotifyUrl, err := url.Parse(link)
			if err != nil {
				WriteJsonResponse(rw, 400,
					models.CreateErrorPayload("link is not a valid url"))
			}
			id = path.Base(spotifyUrl.Path)
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
			WriteJsonResponse(rw, http.StatusOK, playlist)
		}
	}
}

func (s *Server) handleDownloadStart() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		id, ok := GetQueryParameterOrWriteErrorResponse("id", rw, r)
		if !ok {
			return
		}
		filepath, ok := GetQueryParameterOrWriteErrorResponse("path", rw, r)
		if !ok {
			return
		}

		youtubeLink, s2yStatus := s.SonglinkHelper.GetYoutubeLinkBySpotifyId(id)

		switch s2yStatus {
		case songlink.ErrorSendingRequest:
			rw.WriteHeader(http.StatusInternalServerError)
		case songlink.TooManyRequests:
			rw.WriteHeader(http.StatusTooManyRequests)
		case songlink.NoSongWithSuchId:
			WriteJsonResponse(rw,
				http.StatusBadRequest,
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
		}

		if s2yStatus != songlink.Found {
			return
		}

		downloadLink, exists := clihelpers.GetYoutubeDownloadLink(youtubeLink)
		if !exists {
			WriteJsonResponse(rw,
				http.StatusNotFound,
				models.CreateErrorPayload(
					"No download link found for youtube link "+youtubeLink,
				),
			)
			return
		}

		downloadStatus := s.DownloadHelper.StartDownload(
			filepath,
			downloadLink)

		switch downloadStatus {
		case downloader.DStartErrorCreatingFile:
			WriteJsonResponse(rw, 403,
				models.CreateErrorPayload(
					"could not create a file at "+filepath,
				),
			)
		case downloader.DStartErrorSendingRequest, downloader.DStartErrorReadingContentLength:
			rw.WriteHeader(http.StatusInternalServerError)
		case downloader.DStartOk:
			rw.WriteHeader(http.StatusNoContent)
		}
	}
}

func (s *Server) handleDownloadStatus() http.HandlerFunc {
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
			WriteJsonResponse(rw, http.StatusOK, downloadEntry)
		default:
			rw.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func (s *Server) handleDownloadCancel() http.HandlerFunc {
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

func (s *Server) handleFeatures() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		features := struct {
			YoutubeDl bool `json:"youtube_dl"`
			Ffmpeg    bool `json:"ffmpeg"`
		}{
			YoutubeDl: s.FeatureYoutubeDlInstalled,
			Ffmpeg:    s.FeatureFfmpegInstalled,
		}

		WriteJsonResponse(rw,
			http.StatusOK,
			features)
	}
}
