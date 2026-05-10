package v1

import (
	"errors"

	corehttp "ai-assistants-catalog/internal/core/http"
	"ai-assistants-catalog/internal/runs/domain"
)

func mapError(err error) corehttp.APIError {
	if errors.Is(err, domain.ErrAssistantNotFound) {
		return corehttp.ErrNotFound
	}

	if errors.Is(err, domain.ErrAssistantInactive) {
		return corehttp.ErrConflict
	}

	if errors.Is(err, domain.ErrProviderFailed) {
		return corehttp.ErrBadGateway
	}

	if errors.Is(err, domain.ErrInvalidID) ||
		errors.Is(err, domain.ErrInvalidAssistantID) ||
		errors.Is(err, domain.ErrInvalidUserID) ||
		errors.Is(err, domain.ErrInvalidUserPrompt) ||
		errors.Is(err, domain.ErrInvalidStatus) ||
		errors.Is(err, domain.ErrInvalidPagination) {
		return corehttp.ErrInvalidRequest
	}

	return corehttp.ErrInternal
}
