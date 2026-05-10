package v1

import (
	"errors"
	"net/http"

	"ai-assistants-catalog/internal/assistants/domain"
	corehttp "ai-assistants-catalog/internal/core/http"
)

func mapError(err error) corehttp.APIError {
	if errors.Is(err, domain.ErrNotFound) {
		return corehttp.APIError{
			StatusCode: http.StatusNotFound,
			Code:       "ASSISTANT_NOT_FOUND",
			Message:    "assistant not found",
		}
	}

	if errors.Is(err, domain.ErrCategoryNotFound) {
		return corehttp.APIError{
			StatusCode: http.StatusBadRequest,
			Code:       "CATEGORY_NOT_FOUND",
			Message:    "category not found",
		}
	}

	if errors.Is(err, domain.ErrInvalidID) ||
		errors.Is(err, domain.ErrInvalidCategoryID) ||
		errors.Is(err, domain.ErrInvalidName) ||
		errors.Is(err, domain.ErrInvalidDescription) ||
		errors.Is(err, domain.ErrInvalidModel) ||
		errors.Is(err, domain.ErrInvalidSystemPrompt) ||
		errors.Is(err, domain.ErrInvalidActiveState) ||
		errors.Is(err, domain.ErrInvalidPagination) {
		return corehttp.ErrInvalidRequest
	}

	return corehttp.ErrInternal
}
