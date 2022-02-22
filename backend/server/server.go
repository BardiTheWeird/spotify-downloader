package server

import (
	"fmt"
	"net/http"
	"os"
	"spotify-downloader/clihelpers"
	"spotify-downloader/songlink"
	"spotify-downloader/spotify"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	SpotifyHelper  spotify.SpotifyHelper
	SonglinkHelper songlink.SonglinkHelper

	FeatureYoutubeDlInstalled bool
	FeatureFfmpegInstalled    bool

	router *chi.Mux
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(rw, r)
}

func (s *Server) ConfigureFromEnv() {
	getEnvOrDefault := func(envKey, def string) string {
		val, err := os.LookupEnv(envKey)
		if !err {
			fmt.Printf("%s is not set, using %s\n", envKey, def)
			val = def
		}
		return val
	}

	s.SpotifyHelper.ClientId = getEnvOrDefault("CLIENT_ID", "00000000000000000000000000000000")
	s.SpotifyHelper.ClientSecret = getEnvOrDefault("CLIENT_SECRET", "00000000000000000000000000000000")
}

func (s *Server) DiscoverFeatures() {
	_, _, err := clihelpers.RunCliCommand("youtube-dl", "--version")
	if err == nil {
		s.FeatureYoutubeDlInstalled = true
		fmt.Println("youtube-dl detected")
	} else {
		fmt.Println("youtube-dl could not be detected. Downloads will be unavailable")
	}
	_, _, err = clihelpers.RunCliCommand("ffmpeg", "-version")
	if err == nil {
		s.FeatureFfmpegInstalled = true
		fmt.Println("ffmpeg detected")
	} else {
		fmt.Println("ffmpeg could not be detected. Conversion from mp4 will not be available")
	}
}
