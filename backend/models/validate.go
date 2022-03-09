package models

import (
	"fmt"
	"net/http"
	"spotify-downloader/requesthelpers"
	"strings"

	"gopkg.in/go-playground/validator.v9"
)

func Validate(entity interface{}) (string, bool) {
	v := validator.New()
	err := v.Struct(entity)
	if err != nil {
		validationErrors := make([]string, 0)
		for _, e := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, fmt.Sprint(e))
		}
		return strings.Join(validationErrors, "; "), false
	}
	return "", true
}

func ValidateAndWriteResponse(rw http.ResponseWriter, entity interface{}) bool {
	errors, ok := Validate(entity)
	if !ok {
		requesthelpers.WriteJsonResponse(rw, 400,
			requesthelpers.CreateErrorPayload(errors))
	}
	return ok
}
