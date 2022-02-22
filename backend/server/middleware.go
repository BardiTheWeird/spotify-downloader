package server

import (
	"fmt"
	"net/http"
	"spotify-downloader/models"
)

func (s *Server) IsFeatureEnabled(feature *bool, featureName string, h http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if !*feature {
			WriteJsonResponse(rw,
				http.StatusServiceUnavailable,
				models.CreateErrorPayload(
					0,
					fmt.Sprint(featureName, " is not available"),
				))
			return
		}
		h(rw, r)
	}
}
