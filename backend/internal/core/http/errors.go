package http

import "net/http"

var ErrInvalidRequest = APIError{
	StatusCode: http.StatusBadRequest,
	Code:       "INVALID_REQUEST",
	Message:    "invalid request",
}

var ErrUnauthorized = APIError{
	StatusCode: http.StatusUnauthorized,
	Code:       "UNAUTHORIZED",
	Message:    "unauthorized",
}

var ErrForbidden = APIError{
	StatusCode: http.StatusForbidden,
	Code:       "FORBIDDEN",
	Message:    "forbidden",
}

var ErrNotFound = APIError{
	StatusCode: http.StatusNotFound,
	Code:       "NOT_FOUND",
	Message:    "not found",
}

var ErrConflict = APIError{
	StatusCode: http.StatusConflict,
	Code:       "CONFLICT",
	Message:    "conflict",
}

var ErrBadGateway = APIError{
	StatusCode: http.StatusBadGateway,
	Code:       "BAD_GATEWAY",
	Message:    "bad gateway",
}

var ErrInternal = APIError{
	StatusCode: http.StatusInternalServerError,
	Code:       "INTERNAL_ERROR",
	Message:    "internal server error",
}
