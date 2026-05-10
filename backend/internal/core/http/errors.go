package http

import "net/http"

var ErrInvalidRequest = APIError{
	StatusCode: http.StatusBadRequest,
	Code:       "INVALID_REQUEST",
	Message:    "invalid request",
}

var ErrInternal = APIError{
	StatusCode: http.StatusInternalServerError,
	Code:       "INTERNAL_ERROR",
	Message:    "internal server error",
}
