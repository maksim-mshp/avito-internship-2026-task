package handlers

import (
	"context"

	"ai-assistants-catalog/internal/categories/app"
	"ai-assistants-catalog/internal/categories/domain"
)

type CreateHandler struct {
	repo app.Repository
}

func NewCreateHandler(repo app.Repository) *CreateHandler {
	return &CreateHandler{repo: repo}
}

func (h *CreateHandler) Handle(ctx context.Context, cmd app.CreateCommand) (domain.Category, error) {
	var name string
	if cmd.Name != nil {
		name = *cmd.Name
	}

	category, err := domain.NewCategory(name, cmd.Description)
	if err != nil {
		return domain.Category{}, err
	}

	return h.repo.Create(ctx, category)
}
