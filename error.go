package lade

import (
	"net/http"
)

type APIError struct {
	Type    string `json:"error"`
	Message string `json:"error_description"`
	Status  int    `json:"-"`
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
