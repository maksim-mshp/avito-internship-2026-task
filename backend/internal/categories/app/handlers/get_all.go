package handlers

import (
	"context"

	"ai-assistants-catalog/internal/categories/app"
	"ai-assistants-catalog/internal/categories/domain"
)

type GetAllHandler struct {
	repo app.Repository
}

func NewGetAllHandler(repo app.Repository) *GetAllHandler {
	return &GetAllHandler{repo: repo}
}

func (h *GetAllHandler) Handle(ctx context.Context, _ app.GetAllQuery) ([]domain.Category, error) {
	return h.repo.GetAll(ctx)
}
