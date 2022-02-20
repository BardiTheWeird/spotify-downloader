package server

import "github.com/go-chi/chi/v5"

func (s *Server) ConfigureRoutes() {
	r := chi.NewRouter()
	r.Mount("/api/v1", s.apiRouter())
	s.router = r
}

func (s *Server) apiRouter() *chi.Mux {
	r := chi.NewRouter()
	// r.Get("/", s.handleRoot())
	r.Get("/playlist", s.handlePlaylist())
	r.Get("/s2y", s.handleS2Y())
	r.Route("/download", func(r chi.Router) {
		r.Post("/start", s.handleDownloadStart())
		r.Get("/status", s.handleDownloadStatus())
		r.Post("/cancel", s.handleDownloadCancel())
	})
	return r
}
