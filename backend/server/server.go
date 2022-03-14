package server

import (
	"log"
	"net/http"
	"spotify-downloader/clihelpers"
	"spotify-downloader/downloader"
	"spotify-downloader/songlink"
	"spotify-downloader/spotify"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	SpotifyHelper  spotify.SpotifyHelper
	SonglinkHelper songlink.SonglinkHelper
	downloader.DownloadHelper

	SettingsFileLocation string

	FeatureYoutubeDlInstalled bool
	FeatureFfmpegInstalled    bool

	router *chi.Mux
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(rw, r)
}

func (s *Server) DiscoverFeatures() {
	_, _, err := clihelpers.RunCliCommand("youtube-dl", "--version")
	if err == nil {
		s.FeatureYoutubeDlInstalled = true
		log.Println("youtube-dl detected")
	} else {
		log.Println("youtube-dl could not be detected. Downloads will be unavailable")
	}
	_, _, err = clihelpers.RunCliCommand("ffmpeg", "-version")
	if err == nil {
		s.FeatureFfmpegInstalled = true
		log.Println("ffmpeg detected")
	} else {
		log.Println("ffmpeg could not be detected. Conversion from mp4 will not be available")
	}
}
