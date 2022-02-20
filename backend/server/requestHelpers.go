package server

import (
	"fmt"
	"net/http"
	"spotify-downloader/models"
)

func SetContentTypeToJson(rw http.ResponseWriter) {
	rw.Header().Add("Content-Type", "application/json")
}

func WriteJsonResponse(rw http.ResponseWriter, statusCode int, payload []byte) {
	SetContentTypeToJson(rw)
	rw.WriteHeader(statusCode)
	rw.Write(payload)
}

func GetQueryParameterOrWriteErrorResponse(parameter string, rw http.ResponseWriter, r *http.Request) (string, bool) {
	val := r.URL.Query().Get(parameter)
	present := true
	if len(val) == 0 {
		WriteJsonResponse(
			rw,
			400,
			models.CreateErrorPayload(
				400,
				fmt.Sprintf("'%s' query parameter is missing", parameter),
			),
		)
		present = false
	}
	return val, present
}
