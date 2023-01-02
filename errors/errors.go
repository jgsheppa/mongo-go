package errors

import "net/http"

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
