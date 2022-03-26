package server

import (
	"encoding/json"
	"fmt"
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
	"strings"
	"time"
)

func (s *Server) handlePlaylist() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		linkType := "playlist"
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
			splitPath := strings.Split(spotifyUrl.Path, "/")
			if len(splitPath) < 3 {
				requesthelpers.WriteJsonResponse(rw, 400,
					requesthelpers.CreateErrorPayload("Invalid Spotify link"))
				return
			}
			// the path is supposed to be /{playlist/album}/{id}
			linkType = splitPath[1]
			fmt.Println("link type:", linkType)
			id = path.Base(spotifyUrl.Path)
		}

		playlist, status := s.SpotifyHelper.GetPlaylistById(
			id,
			linkType+"s",
			r.Header.Get("Authorization"))

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

		var downloadStatus downloader.DownloadStartStatus
		for i := 0; i < 3; i++ {
			downloadLink, exists := s.GetYoutubeDownloadLink(youtubeLink)
			if !exists {
				requesthelpers.WriteJsonResponse(rw,
					http.StatusNotFound,
					requesthelpers.CreateErrorPayload(
						"No download link found for "+youtubeLink,
					),
				)
				return
			}

			downloadStatus = s.DownloadHelper.StartDownload(
				downloadRequest.Id,
				downloadRequest.Folder,
				downloadRequest.Filename,
				downloadLink,
				clihelpers.FfmpegMetadata{
					Title:  downloadRequest.Title,
					Artist: downloadRequest.Artist,
					Album:  downloadRequest.Album,
					Image:  downloadRequest.Image,
				})

			if downloadStatus != downloader.DStartErrorInvalidUrl {
				break
			}

			if i < 2 {
				time.Sleep(time.Second * 3)
			}
		}

		switch downloadStatus {
		case downloader.DStartErrorCreatingFile:
			requesthelpers.WriteJsonResponse(rw, 403,
				requesthelpers.CreateErrorPayload(
					"could not create a file at "+filepath.Join(downloadRequest.Folder, downloadRequest.Filename),
				),
			)
		case downloader.DStartErrorInvalidUrl:
			rw.WriteHeader(http.StatusRequestTimeout)
		case downloader.DStartErrorSendingRequest, downloader.DStartErrorReadingContentLength:
			rw.WriteHeader(http.StatusInternalServerError)
		case downloader.DStartOk:
			rw.WriteHeader(http.StatusNoContent)
		}
	}
}

func (s *Server) handleDownloadStatus() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		trackId, ok := requesthelpers.GetQueryParameterOrWriteErrorResponse("id", rw, r)
		if !ok {
			return
		}

		downloadEntry, responseStatus := s.DownloadHelper.GetDownloadStatus(trackId)
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
		trackId, ok := requesthelpers.GetQueryParameterOrWriteErrorResponse("id", rw, r)
		if !ok {
			return
		}

		switch status := s.DownloadHelper.CancelDownload(trackId); status {
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

func (s *Server) handleConfigureFfmpeg() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		path, ok := requesthelpers.GetQueryParameterOrWriteErrorResponse("path", rw, r)
		if !ok {
			return
		}

		s.FfmpegPath = path
		s.FeatureFfmpegInstalled = s.DiscoverFeature(path, "-version")

		if s.FeatureFfmpegInstalled {
			rw.WriteHeader(http.StatusNoContent)
		} else {
			rw.WriteHeader(http.StatusNotFound)
		}
	}
}

func (s *Server) handleConfigureYoutubeDl() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		path, ok := requesthelpers.GetQueryParameterOrWriteErrorResponse("path", rw, r)
		if !ok {
			return
		}

		s.YoutubeDlPath = path
		s.FeatureYoutubeDlInstalled = s.DiscoverFeature(path, "-version")

		if s.FeatureYoutubeDlInstalled {
			rw.WriteHeader(http.StatusNoContent)
		} else {
			rw.WriteHeader(http.StatusNotFound)
		}
	}
}
