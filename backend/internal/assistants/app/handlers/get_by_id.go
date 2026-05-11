package handlers

import (
	"context"

	"ai-assistants-catalog/internal/assistants/app"
	"ai-assistants-catalog/internal/assistants/domain"
)

type GetByIDHandler struct {
	repo app.Repository
}

func NewGetByIDHandler(repo app.Repository) *GetByIDHandler {
	return &GetByIDHandler{repo: repo}
}

func (h *GetByIDHandler) Handle(ctx context.Context, query app.GetByIDQuery) (domain.Assistant, error) {
	if !domain.IsValidID(query.ID) {
		return domain.Assistant{}, domain.ErrInvalidID
	}

	return h.repo.GetByID(ctx, query.ID, query.IncludeInactive, query.UserID)
}
