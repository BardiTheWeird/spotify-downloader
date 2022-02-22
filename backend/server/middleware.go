package server

import (
	"fmt"
	"log"
	"net/http"
	"spotify-downloader/models"
	"time"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc
type MiddlewareArr struct {
	Arr []Middleware
}

func Chain(ms ...Middleware) MiddlewareArr {
	return MiddlewareArr{ms}
}

func (m MiddlewareArr) ThenChain(ms ...Middleware) MiddlewareArr {
	return Chain(append(m.Arr, ms...)...)
}

func (m MiddlewareArr) Then(h http.HandlerFunc) http.HandlerFunc {
	for i := range m.Arr {
		h = m.Arr[len(m.Arr)-1-i](h)
	}
	return h
}

func IsFeatureEnabled(feature *bool, featureName string) Middleware {
	return func(h http.HandlerFunc) http.HandlerFunc {
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
}

func LogEndpoint() Middleware {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) {
			timeStart := time.Now()
			h(rw, r)
			log.Printf("%s %s from %s %s",
				r.Method,
				r.URL.Path,
				r.RemoteAddr,
				time.Since(timeStart).String(),
			)
		}
	}
}
