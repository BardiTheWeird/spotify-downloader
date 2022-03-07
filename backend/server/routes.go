package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func (s *Server) ConfigureRoutes() {
	r := chi.NewRouter()
	r.Use(LogEndpoint())
	r.Mount("/api/v1", s.apiRouter())
	s.router = r
}

func (s *Server) apiRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

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
	r.Get("/features", s.handleFeatures())
	return r
}
