package models

import "encoding/json"

type ErrorPayload struct {
	StatusCode   int    `json:"status_code,omitempty"`
	ErrorMessage string `json:"error_message"`
}

func CreateErrorPayloadWithCode(statusCode int, errorMessage string) []byte {
	response := ErrorPayload{
		StatusCode:   statusCode,
		ErrorMessage: errorMessage,
	}
	bytes, _ := json.Marshal(response)
	return bytes
}

func CreateErrorPayload(errorMessage string) []byte {
	return CreateErrorPayloadWithCode(0, errorMessage)
}
