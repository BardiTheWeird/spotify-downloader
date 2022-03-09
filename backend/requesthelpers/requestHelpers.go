package requesthelpers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func SetContentTypeToJson(rw http.ResponseWriter) {
	rw.Header().Add("Content-Type", "application/json")
}

func WriteJsonResponse(rw http.ResponseWriter, statusCode int, payload interface{}) {
	bytes, _ := json.Marshal(payload)
	SetContentTypeToJson(rw)
	rw.WriteHeader(statusCode)
	rw.Write(bytes)
}

func GetQueryParameter(parameter string, r *http.Request) (string, bool) {
	val := r.URL.Query().Get(parameter)
	return val, len(val) != 0
}

func GetQueryParameterOrWriteErrorResponse(parameter string, rw http.ResponseWriter, r *http.Request) (string, bool) {
	val, present := GetQueryParameter(parameter, r)
	if !present {
		WriteJsonResponse(
			rw,
			400,
			CreateErrorPayload(
				fmt.Sprintf("'%s' query parameter is missing", parameter)),
		)
	}
	return val, present
}
