package v1

import (
	"errors"

	"ai-assistants-catalog/internal/categories/domain"
	corehttp "ai-assistants-catalog/internal/core/http"
)

func mapError(err error) corehttp.APIError {
	if errors.Is(err, domain.ErrInvalidName) {
		return corehttp.ErrInvalidRequest
	}

	return corehttp.ErrInternal
}
