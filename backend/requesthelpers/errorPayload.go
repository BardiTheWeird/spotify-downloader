package requesthelpers

type ErrorPayload struct {
	StatusCode   int    `json:"status_code,omitempty"`
	ErrorMessage string `json:"error_message"`
}

func CreateErrorPayloadWithCode(statusCode int, errorMessage string) ErrorPayload {
	return ErrorPayload{
		StatusCode:   statusCode,
		ErrorMessage: errorMessage,
	}
}

func CreateErrorPayload(errorMessage string) ErrorPayload {
	return CreateErrorPayloadWithCode(0, errorMessage)
}
