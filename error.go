package lade

import (
	"net/http"
)

type APIError struct {
	Status  int    `json:"-"`
	Type    string `json:"error"`
	Message string `json:"error_description"`
}

func (e *APIError) Error() string {
	message := e.Message
	if message == "" {
		message = http.StatusText(e.Status)
	}
	return message
}

var ErrNotFound = &APIError{
	Type:    "not_found",
	Message: "Resource not found",
	Status:  http.StatusNotFound,
}

var ErrServerError = &APIError{
	Type:    "server_error",
	Message: "Unexpected server error",
	Status:  http.StatusInternalServerError,
}
