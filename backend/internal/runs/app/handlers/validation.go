package handlers

import "ai-assistants-catalog/internal/runs/domain"

func validateStatus(status *string) error {
	if status != nil && !domain.IsValidStatus(*status) {
		return domain.ErrInvalidStatus
	}

	return nil
}

func validatePagination(page int, pageSize int) error {
	if page <= 0 || pageSize <= 0 || pageSize > 100 {
		return domain.ErrInvalidPagination
	}

	return nil
}
