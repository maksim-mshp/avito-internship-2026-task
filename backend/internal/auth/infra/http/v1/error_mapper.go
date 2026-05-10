package v1

import (
	"errors"

	"ai-assistants-catalog/internal/auth/domain"
	corehttp "ai-assistants-catalog/internal/core/http"
)

func mapError(err error) corehttp.APIError {
	if errors.Is(err, domain.ErrInvalidRole) {
		return corehttp.ErrInvalidRequest
	}

	return corehttp.ErrInternal
}
