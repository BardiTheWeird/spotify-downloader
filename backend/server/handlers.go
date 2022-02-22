package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	cliHelpers "spotify-downloader/cliHelpers"
	"spotify-downloader/downloader"
	"spotify-downloader/models"
	"spotify-downloader/songlink"
	"spotify-downloader/spotify"
)

// "/playlist"
func (s *Server) handlePlaylist() func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		id, ok := GetQueryParameterOrWriteErrorResponse("id", rw, r)
		if !ok {
			return
		}

		playlist, status := s.SpotifyHelper.GetPlaylistById(id)
		if status == spotify.BadOrExpiredToken {
			s.SpotifyHelper.Authenticate()
			playlist, status = s.SpotifyHelper.GetPlaylistById(id)
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
		case songlink.NoSongWithSuchId:
			WriteJsonResponse(rw,
				http.StatusNotFound,
				models.CreateErrorPayload(
					404,
					fmt.Sprintf("No entry for song with id %s", id),
				),
			)
		case songlink.NoYoutubeLinkForSong:
			WriteJsonResponse(rw,
				http.StatusNotFound,
				models.CreateErrorPayload(
					404,
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
		if !s.FeatureYoutubeDlInstalled {
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

		downloadLink, exists := cliHelpers.GetYoutubeDownloadLink(youtubeLink)
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
	}
}

// "/download/status"
func (s *Server) handleDownloadStatus() func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
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
	}
}

// "/download/cancel"
func (s *Server) handleDownloadCancel() func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
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
	}
}
