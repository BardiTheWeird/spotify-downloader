package server

import (
	"github.com/go-chi/chi/v5"
)

func (s *Server) ConfigureRoutes() {
	r := chi.NewRouter()
	r.Use(LogEndpoint())
	r.Mount("/api/v1", s.apiRouter())
	s.router = r
}

func (s *Server) apiRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Route("/spotify", func(r chi.Router) {
		r.Get("/playlist", s.handlePlaylist())
		r.Post("/configure", s.handleSpotifyConfigure())
	})
	r.Get("/s2y", s.handleS2Y())
	r.Route("/download", func(r chi.Router) {
		r.With(IsFeatureEnabled(&s.FeatureYoutubeDlInstalled, "youtube-dl")).
			Post("/start", s.handleDownloadStart())
		r.Get("/status", s.handleDownloadStatus())
		r.Post("/cancel", s.handleDownloadCancel())
	})
	return r
}
