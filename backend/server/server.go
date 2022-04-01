package server

import (
	"flag"
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
	s.ConfigureFromCli()
	s.FeaturesSetDefaults()

	s.DownloadHelper.CliHelper = &s.CliHelper
}

func (s *Server) ConfigureFromCli() {
	flag.StringVar(&s.Features.Ffmpeg.Path, "ffmpeg-path", "ffmpeg", "configure path to ffmpeg")
	flag.StringVar(&s.Features.YoutubeDl.Path, "youtube-dl-path", "youtube-dl", "configure path to youtube-dl")
	flag.Parse()
}
