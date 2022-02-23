package server

import (
	"log"
	"net/http"
	"spotify-downloader/models"
	"time"
)

func LogEndpoint() func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			timeStart := time.Now()
			h.ServeHTTP(rw, r)
			log.Printf("%s %s from %s %s",
				r.Method,
				r.URL.Path,
				r.RemoteAddr,
				time.Since(timeStart).String(),
			)
		})
	}
}

func IsFeatureEnabled(feature *bool, featureName string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			if !*feature {
				WriteJsonResponse(rw,
					http.StatusServiceUnavailable,
					models.CreateErrorPayload(
						featureName+" is not available",
					))
				return
			}
			h.ServeHTTP(rw, r)
		})
	}
}
