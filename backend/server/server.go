package server

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	cliHelpers "spotify-downloader/cliHelpers"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	SpotifyClientId     string
	SpotifyClientSecret string

	FeatureYoutubeDlInstalled bool
	FeatureFfmpegInstalled    bool

	router *chi.Mux
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(rw, r)
}

func (s *Server) GetB64() string {
	return base64.RawStdEncoding.Strict().
		EncodeToString([]byte(s.SpotifyClientId + ":" + s.SpotifyClientSecret))
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

	s.SpotifyClientId = getEnvOrDefault("CLIENT_ID", "00000000000000000000000000000000")
	s.SpotifyClientSecret = getEnvOrDefault("CLIENT_SECRET", "00000000000000000000000000000000")

	_, _, err := cliHelpers.RunCliCommand("youtube-dl", "--version")
	if err == nil {
		s.FeatureYoutubeDlInstalled = true
		fmt.Println("youtube-dl detected")
	} else {
		fmt.Println("youtube-dl could not be detected. Downloads will be unavailable")
	}
	_, _, err = cliHelpers.RunCliCommand("ffmpeg", "-version")
	if err == nil {
		s.FeatureFfmpegInstalled = true
		fmt.Println("ffmpeg detected")
	} else {
		fmt.Println("ffmpeg could not be detected. Conversion from mp4 will not be available")
	}
}
