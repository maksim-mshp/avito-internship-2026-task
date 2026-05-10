package handlers

import (
	"context"

	"ai-assistants-catalog/internal/runs/app"
	"ai-assistants-catalog/internal/runs/domain"
)

type ListMyHandler struct {
	runs app.RunRepository
}

func NewListMyHandler(runs app.RunRepository) *ListMyHandler {
	return &ListMyHandler{runs: runs}
}

func (h *ListMyHandler) Handle(ctx context.Context, query app.ListMyQuery) (app.ListResult, error) {
	if !domain.IsValidID(query.UserID) {
		return app.ListResult{}, domain.ErrInvalidUserID
	}

	if err := validateStatus(query.Status); err != nil {
		return app.ListResult{}, err
	}

	if err := validatePagination(query.Page, query.PageSize); err != nil {
		return app.ListResult{}, err
	}

	return h.runs.ListMy(ctx, query)
}
