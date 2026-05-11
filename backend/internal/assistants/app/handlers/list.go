package handlers

import (
	"context"
	"strings"

	"ai-assistants-catalog/internal/assistants/app"
	"ai-assistants-catalog/internal/assistants/domain"
)

type ListHandler struct {
	repo app.Repository
}

func NewListHandler(repo app.Repository) *ListHandler {
	return &ListHandler{repo: repo}
}

func (h *ListHandler) Handle(ctx context.Context, query app.ListQuery) (app.ListResult, error) {
	if query.CategoryID != nil && !domain.IsValidID(*query.CategoryID) {
		return app.ListResult{}, domain.ErrInvalidCategoryID
	}

	if query.Page <= 0 || query.PageSize <= 0 || query.PageSize > 100 {
		return app.ListResult{}, domain.ErrInvalidPagination
	}

	if query.Search != nil {
		search := strings.TrimSpace(*query.Search)
		if search == "" {
			query.Search = nil
		} else {
			query.Search = &search
		}
	}

	if query.Tag != nil {
		tag := strings.TrimSpace(*query.Tag)
		if tag == "" {
			query.Tag = nil
		} else {
			query.Tag = &tag
		}
	}

	return h.repo.List(ctx, query)
}
