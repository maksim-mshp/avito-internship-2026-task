package handlers

import (
	"context"

	"ai-assistants-catalog/internal/assistants/app"
	"ai-assistants-catalog/internal/assistants/domain"
)

type UpdateHandler struct {
	repo app.Repository
}

func NewUpdateHandler(repo app.Repository) *UpdateHandler {
	return &UpdateHandler{repo: repo}
}

func (h *UpdateHandler) Handle(ctx context.Context, cmd app.UpdateCommand) (domain.Assistant, error) {
	if !domain.IsValidID(cmd.ID) {
		return domain.Assistant{}, domain.ErrInvalidID
	}

	if cmd.IsActive == nil {
		return domain.Assistant{}, domain.ErrInvalidActiveState
	}

	assistant, err := domain.NewAssistant(
		valueOrEmpty(cmd.CategoryID),
		valueOrEmpty(cmd.Name),
		valueOrEmpty(cmd.Description),
		valueOrEmpty(cmd.Model),
		valueOrEmpty(cmd.SystemPrompt),
		cmd.ExampleUserPrompt,
		cmd.Tags,
		cmd.IsActive,
	)
	if err != nil {
		return domain.Assistant{}, err
	}

	assistant.ID = cmd.ID

	return h.repo.Update(ctx, assistant)
}
