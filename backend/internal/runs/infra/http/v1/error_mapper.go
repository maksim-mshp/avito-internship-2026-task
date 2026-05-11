package v1

import (
	"errors"
	"net/http"

	corehttp "ai-assistants-catalog/internal/core/http"
	"ai-assistants-catalog/internal/runs/domain"
)

func mapError(err error) corehttp.APIError {
	if errors.Is(err, domain.ErrNotFound) {
		return corehttp.APIError{
			StatusCode: http.StatusNotFound,
			Code:       "RUN_NOT_FOUND",
			Message:    "run not found",
		}
	}

	if errors.Is(err, domain.ErrAssistantNotFound) {
		return corehttp.APIError{
			StatusCode: http.StatusNotFound,
			Code:       "ASSISTANT_NOT_FOUND",
			Message:    "assistant not found",
		}
	}

	if errors.Is(err, domain.ErrAssistantInactive) {
		return corehttp.APIError{
			StatusCode: http.StatusConflict,
			Code:       "ASSISTANT_INACTIVE",
			Message:    "assistant inactive",
		}
	}

	if errors.Is(err, domain.ErrProviderFailed) {
		return corehttp.APIError{
			StatusCode: http.StatusBadGateway,
			Code:       "LLM_PROVIDER_ERROR",
			Message:    "llm provider error",
		}
	}

	if errors.Is(err, domain.ErrInvalidID) ||
		errors.Is(err, domain.ErrInvalidAssistantID) ||
		errors.Is(err, domain.ErrInvalidUserID) ||
		errors.Is(err, domain.ErrInvalidUserPrompt) ||
		errors.Is(err, domain.ErrInvalidStatus) ||
		errors.Is(err, domain.ErrInvalidRating) ||
		errors.Is(err, domain.ErrInvalidPagination) {
		return corehttp.ErrInvalidRequest
	}

	return corehttp.ErrInternal
}
