package handlers

import (
	"context"

	"ai-assistants-catalog/internal/runs/app"
	"ai-assistants-catalog/internal/runs/domain"
)

type ListAdminHandler struct {
	runs app.RunRepository
}

func NewListAdminHandler(runs app.RunRepository) *ListAdminHandler {
	return &ListAdminHandler{runs: runs}
}

func (h *ListAdminHandler) Handle(ctx context.Context, query app.ListAdminQuery) (app.ListResult, error) {
	if query.AssistantID != nil && !domain.IsValidID(*query.AssistantID) {
		return app.ListResult{}, domain.ErrInvalidAssistantID
	}

	if err := validateStatus(query.Status); err != nil {
		return app.ListResult{}, err
	}

	if err := validatePagination(query.Page, query.PageSize); err != nil {
		return app.ListResult{}, err
	}

	return h.runs.ListAdmin(ctx, query)
}
