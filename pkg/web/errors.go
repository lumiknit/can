package web

import (
	"fmt"
	"net/http"
)

// WebError represents an HTTP error with status code and message
type WebError struct {
	Code    int
	Message string
}

func (e *WebError) Error() string {
	return fmt.Sprintf("HTTP %d: %s", e.Code, e.Message)
}

// WriteError writes an HTTP error response
func (e *WebError) WriteError(w http.ResponseWriter) {
	http.Error(w, e.Message, e.Code)
}

// Common HTTP errors
var (
	ErrBadRequest          = &WebError{Code: http.StatusBadRequest, Message: "Bad Request"}
	ErrUnauthorized        = &WebError{Code: http.StatusUnauthorized, Message: "Unauthorized"}
	ErrNotFound            = &WebError{Code: http.StatusNotFound, Message: "Not Found"}
	ErrInternalServerError = &WebError{Code: http.StatusInternalServerError, Message: "Internal Server Error"}
	ErrTooManyRequests     = &WebError{Code: http.StatusTooManyRequests, Message: "Too Many Requests"}
)

// NewWebError creates a new WebError with custom message
func NewWebError(code int, message string) *WebError {
	return &WebError{Code: code, Message: message}
}
