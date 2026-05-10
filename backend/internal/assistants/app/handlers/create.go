package handlers

import (
	"context"

	"ai-assistants-catalog/internal/assistants/app"
	"ai-assistants-catalog/internal/assistants/domain"
)

type CreateHandler struct {
	repo app.Repository
}

func NewCreateHandler(repo app.Repository) *CreateHandler {
	return &CreateHandler{repo: repo}
}

func (h *CreateHandler) Handle(ctx context.Context, cmd app.CreateCommand) (domain.Assistant, error) {
	categoryID := valueOrEmpty(cmd.CategoryID)
	name := valueOrEmpty(cmd.Name)
	description := valueOrEmpty(cmd.Description)
	model := valueOrEmpty(cmd.Model)
	systemPrompt := valueOrEmpty(cmd.SystemPrompt)

	assistant, err := domain.NewAssistant(
		categoryID,
		name,
		description,
		model,
		systemPrompt,
		cmd.ExampleUserPrompt,
		cmd.IsActive,
	)
	if err != nil {
		return domain.Assistant{}, err
	}

	return h.repo.Create(ctx, assistant)
}

func valueOrEmpty(value *string) string {
	if value == nil {
		return ""
	}

	return *value
}
