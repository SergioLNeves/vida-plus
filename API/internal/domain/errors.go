package domain

import (
	"errors"
	"fmt"
	"net/http"
)

type APIError struct {
	Type    string `json:"type,omitempty"`
	Title   string `json:"title,omitempty"`
	Status  int    `json:"status,omitempty"`
	Details any    `json:"details,omitempty"`
}

func (e *APIError) Error() string {
	if e.Details != nil {
		return fmt.Sprint(e.Details)
	}
	return e.Title
}

var (
	ErrDocumentNotFound = errors.New("document not found")
)

func NewAPIError(statusCode int, details any) *APIError {
	return &APIError{
		Type:    getTypeByStatusCode(statusCode),
		Title:   http.StatusText(statusCode),
		Status:  statusCode,
		Details: details,
	}
}

// NewInternalError creates an internal server error
func NewInternalError(message string) *APIError {
	return NewAPIError(http.StatusInternalServerError, message)
}

// NewConflictError creates a conflict error
func NewConflictError(message string) *APIError {
	return NewAPIError(http.StatusConflict, message)
}

// NewNotFoundError creates a not found error
func NewNotFoundError(message string) *APIError {
	return NewAPIError(http.StatusNotFound, message)
}

// NewBadRequestError creates a bad request error
func NewBadRequestError(message string) *APIError {
	return NewAPIError(http.StatusBadRequest, message)
}

// NewUnauthorizedError creates an unauthorized error
func NewUnauthorizedError(message string) *APIError {
	return NewAPIError(http.StatusUnauthorized, message)
}

func getTypeByStatusCode(statusCode int) string {
	switch statusCode {
	case http.StatusBadRequest:
		return "https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/400"
	case http.StatusUnauthorized:
		return "https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/401"
	case http.StatusForbidden:
		return "https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/403"
	case http.StatusNotFound:
		return "https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/404"
	case http.StatusConflict:
		return "https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/409"
	case http.StatusInternalServerError:
		return "https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/500"
	}
	return ""
}
