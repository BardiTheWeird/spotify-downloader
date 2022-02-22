package server

import (
	"github.com/go-chi/chi/v5"
)

func (s *Server) ConfigureRoutes() {
	r := chi.NewRouter()
	r.Mount("/api/v1", s.apiRouter())
	s.router = r
}

func (s *Server) apiRouter() *chi.Mux {
	r := chi.NewRouter()
	commonMiddleware := Chain(LogEndpoint())
	// r.Get("/", s.handleRoot())
	r.Get("/playlist", commonMiddleware.Then(s.handlePlaylist()))
	r.Get("/s2y", commonMiddleware.Then(s.handleS2Y()))
	r.Route("/download", func(r chi.Router) {
		r.Post("/start", commonMiddleware.ThenChain(
			IsFeatureEnabled(&s.FeatureYoutubeDlInstalled, "youtube-dl")).
			Then(s.handleDownloadStart()))
		r.Get("/status", commonMiddleware.Then(s.handleDownloadStatus()))
		r.Post("/cancel", commonMiddleware.Then(s.handleDownloadCancel()))
	})
	return r
}
