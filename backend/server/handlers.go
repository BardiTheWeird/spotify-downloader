package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"spotify-downloader/clihelpers"
	"spotify-downloader/downloader"
	"spotify-downloader/models"
	"spotify-downloader/requesthelpers"
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
			requesthelpers.WriteJsonResponse(rw, 400,
				requesthelpers.CreateErrorPayload("client_id is empty"))
			return
		}
		if len(configuration.ClientSecret) == 0 {
			requesthelpers.WriteJsonResponse(rw, 400,
				requesthelpers.CreateErrorPayload("client_secret is empty"))
			return
		}

		oldClientId := s.SpotifyHelper.ClientId
		oldClientSecret := s.SpotifyHelper.ClientSecret
		s.SpotifyHelper.ClientId = configuration.ClientId
		s.SpotifyHelper.ClientSecret = configuration.ClientSecret

		if !s.SpotifyHelper.UpdateClientToken() {
			requesthelpers.WriteJsonResponse(rw, 400,
				requesthelpers.CreateErrorPayloadWithCode(401,
					"can't authenticate using new credentials"))

			s.SpotifyHelper.ClientId = oldClientId
			s.SpotifyHelper.ClientSecret = oldClientSecret
			s.SpotifyHelper.UpdateClientToken()
			return
		}
		s.UpdateSettingsFile()
		rw.WriteHeader(http.StatusNoContent)
	}
}

func (s *Server) handlePlaylist() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		id, ok := requesthelpers.GetQueryParameter("id", r)
		if !ok {
			link, ok := requesthelpers.GetQueryParameterOrWriteErrorResponse("link", rw, r)
			if !ok {
				return
			}
			spotifyUrl, err := url.Parse(link)
			if err != nil {
				requesthelpers.WriteJsonResponse(rw, 400,
					requesthelpers.CreateErrorPayload("link is not a valid url"))
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
			requesthelpers.WriteJsonResponse(rw,
				http.StatusHTTPVersionNotSupported,
				requesthelpers.CreateErrorPayload("bad spotify client id or key"))
		case spotify.Ok:
			requesthelpers.WriteJsonResponse(rw, http.StatusOK, playlist)
		}
	}
}

func (s *Server) handleDownloadStart() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var downloadRequest struct {
			Id       string `json:"id" validate:"required"`
			Folder   string `json:"folder" validate:"required,dir"`
			Filename string `json:"filename" validate:"required"`

			Title  string `json:"title"`
			Artist string `json:"artist"`
			Album  string `json:"album"`
			Image  string `json:"image"`
		}
		json.NewDecoder(r.Body).Decode(&downloadRequest)

		if !models.ValidateAndWriteResponse(rw, downloadRequest) {
			return
		}

		youtubeLink, s2yStatus := s.SonglinkHelper.GetYoutubeLinkBySpotifyId(downloadRequest.Id)

		switch s2yStatus {
		case songlink.ErrorSendingRequest:
			rw.WriteHeader(http.StatusInternalServerError)
		case songlink.TooManyRequests:
			rw.WriteHeader(http.StatusTooManyRequests)
		case songlink.NoSongWithSuchId:
			requesthelpers.WriteJsonResponse(rw,
				http.StatusBadRequest,
				requesthelpers.CreateErrorPayload(
					fmt.Sprintf("No entry for song with id %s", downloadRequest.Id),
				),
			)
		case songlink.NoYoutubeLinkForSong:
			requesthelpers.WriteJsonResponse(rw,
				http.StatusNotFound,
				requesthelpers.CreateErrorPayload(
					fmt.Sprintf("No YouTube link for song with id %s", downloadRequest.Id),
				),
			)
		}

		if s2yStatus != songlink.Found {
			return
		}

		downloadLink, exists := clihelpers.GetYoutubeDownloadLink(youtubeLink)
		if !exists {
			requesthelpers.WriteJsonResponse(rw,
				http.StatusNotFound,
				requesthelpers.CreateErrorPayload(
					"No download link found for youtube link "+youtubeLink,
				),
			)
			return
		}

		downloadStatus := s.DownloadHelper.StartDownload(
			downloadRequest.Folder,
			downloadRequest.Filename,
			downloadLink,
			clihelpers.FfmpegMetadata{
				Title:  downloadRequest.Title,
				Artist: downloadRequest.Artist,
				Album:  downloadRequest.Album,
			},
			s.FeatureFfmpegInstalled)

		switch downloadStatus {
		case downloader.DStartErrorCreatingFile:
			requesthelpers.WriteJsonResponse(rw, 403,
				requesthelpers.CreateErrorPayload(
					"could not create a file at "+filepath.Join(downloadRequest.Folder, downloadRequest.Filename),
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
		folder, ok := requesthelpers.GetQueryParameterOrWriteErrorResponse("folder", rw, r)
		if !ok {
			return
		}
		filename, ok := requesthelpers.GetQueryParameterOrWriteErrorResponse("filename", rw, r)
		if !ok {
			return
		}

		downloadEntry, responseStatus := s.DownloadHelper.GetDownloadStatus(folder, filename)
		switch responseStatus {
		case downloader.DStatusNotFound:
			rw.WriteHeader(http.StatusNotFound)
		case downloader.DStatusOk:
			requesthelpers.WriteJsonResponse(rw, http.StatusOK, downloadEntry)
		default:
			rw.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func (s *Server) handleDownloadCancel() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		folder, ok := requesthelpers.GetQueryParameterOrWriteErrorResponse("folder", rw, r)
		if !ok {
			return
		}
		filename, ok := requesthelpers.GetQueryParameterOrWriteErrorResponse("filename", rw, r)
		if !ok {
			return
		}

		switch status := s.DownloadHelper.CancelDownload(folder, filename); status {
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

		requesthelpers.WriteJsonResponse(rw,
			http.StatusOK,
			features)
	}
}
