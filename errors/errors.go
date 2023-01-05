package errors

import (
	"errors"
	"net/http"
)

var ErrNoToken = errors.New("no token found or token is invalid")

type ErrorResponse struct {
	Message      string
	Error        bool
	ErrorMessage error
	StatusCode   int
}

func NotFound(err error) ErrorResponse {
	return ErrorResponse{
		Message:      "Document not found",
		Error:        true,
		ErrorMessage: err,
		StatusCode:   http.StatusNotFound,
	}
}

func InternalError(message string, err error) ErrorResponse {
	return ErrorResponse{
		Message:      message,
		Error:        true,
		ErrorMessage: err,
		StatusCode:   http.StatusInternalServerError,
	}
}

func Unauthorized(err error) ErrorResponse {
	return ErrorResponse{
		Message:      "Unauthorized",
		Error:        true,
		ErrorMessage: err,
		StatusCode:   http.StatusUnauthorized,
	}
}
