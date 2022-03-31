package server

import (
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

	clihelpers.CliHelper

	router *chi.Mux
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(rw, r)
}

func (s *Server) ConfigureDefaults() {
	s.ConfigureRoutes()

	s.SonglinkHelper.SetDefaults()
	s.FeaturesSetDefaults()

	s.DownloadHelper.CliHelper = &s.CliHelper
}
