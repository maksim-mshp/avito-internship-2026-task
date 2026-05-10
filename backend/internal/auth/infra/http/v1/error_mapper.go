package v1

import (
	"errors"
	"net/http"

	"ai-assistants-catalog/internal/auth/domain"
	corehttp "ai-assistants-catalog/internal/core/http"
)

func mapError(err error) corehttp.APIError {
	if errors.Is(err, domain.ErrInvalidRole) {
		return corehttp.ErrInvalidRequest
	}

	if errors.Is(err, domain.ErrInvalidEmail) || errors.Is(err, domain.ErrInvalidPassword) {
		return corehttp.ErrInvalidRequest
	}

	if errors.Is(err, domain.ErrEmailTaken) {
		return corehttp.APIError{
			StatusCode: http.StatusConflict,
			Code:       "EMAIL_TAKEN",
			Message:    "email taken",
		}
	}

	if errors.Is(err, domain.ErrInvalidCredentials) {
		return corehttp.APIError{
			StatusCode: http.StatusUnauthorized,
			Code:       "INVALID_CREDENTIALS",
			Message:    "invalid credentials",
		}
	}

	return corehttp.ErrInternal
}
