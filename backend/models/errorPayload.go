package models

type ErrorPayload struct {
	StatusCode   int    `json:"status_code,omitempty"`
	ErrorMessage string `json:"error_message"`
}

func CreateErrorPayloadWithCode(statusCode int, errorMessage string) ErrorPayload {
	// response := ErrorPayload{
	// 	StatusCode:   statusCode,
	// 	ErrorMessage: errorMessage,
	// }
	// bytes, _ := json.Marshal(response)
	// return bytes

	return ErrorPayload{
		StatusCode:   statusCode,
		ErrorMessage: errorMessage,
	}
}

func CreateErrorPayload(errorMessage string) ErrorPayload {
	return CreateErrorPayloadWithCode(0, errorMessage)
}
