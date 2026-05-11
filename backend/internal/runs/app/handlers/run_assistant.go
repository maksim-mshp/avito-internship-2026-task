package handlers

import (
	"context"
	"errors"
	"strings"
	"time"

	assistantdomain "ai-assistants-catalog/internal/assistants/domain"
	"ai-assistants-catalog/internal/runs/app"
	"ai-assistants-catalog/internal/runs/domain"
)

type RunAssistantHandler struct {
	runs       app.RunRepository
	assistants app.AssistantRepository
	provider   app.LLMProvider
}

func NewRunAssistantHandler(
	runs app.RunRepository,
	assistants app.AssistantRepository,
	provider app.LLMProvider,
) *RunAssistantHandler {
	return &RunAssistantHandler{
		runs:       runs,
		assistants: assistants,
		provider:   provider,
	}
}

func (h *RunAssistantHandler) Handle(ctx context.Context, cmd app.RunAssistantCommand) (domain.Run, error) {
	if !domain.IsValidID(cmd.AssistantID) {
		return domain.Run{}, domain.ErrInvalidAssistantID
	}

	if !domain.IsValidID(cmd.UserID) {
		return domain.Run{}, domain.ErrInvalidUserID
	}

	userPrompt := valueOrEmpty(cmd.UserPrompt)
	if strings.TrimSpace(userPrompt) == "" {
		return domain.Run{}, domain.ErrInvalidUserPrompt
	}

	assistant, err := h.assistants.GetByID(ctx, cmd.AssistantID, true, "")
	if err != nil {
		if errors.Is(err, assistantdomain.ErrNotFound) {
			return domain.Run{}, domain.ErrAssistantNotFound
		}

		return domain.Run{}, err
	}

	if !assistant.IsActive {
		return domain.Run{}, domain.ErrAssistantInactive
	}

	run, err := domain.NewPendingRun(assistant.ID, cmd.UserID, assistant.Model, userPrompt)
	if err != nil {
		return domain.Run{}, err
	}

	created, err := h.runs.CreatePending(ctx, run)
	if err != nil {
		return domain.Run{}, err
	}

	providerCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	output, err := h.provider.Generate(providerCtx, app.LLMRequest{
		Model:        assistant.Model,
		SystemPrompt: assistant.SystemPrompt,
		UserPrompt:   run.UserPrompt,
	})
	if err != nil {
		_, failErr := h.runs.Fail(ctx, created.ID, err.Error())
		if failErr != nil {
			return domain.Run{}, failErr
		}

		return domain.Run{}, domain.ErrProviderFailed
	}

	output = strings.TrimSpace(output)
	if output == "" {
		_, failErr := h.runs.Fail(ctx, created.ID, "empty provider response")
		if failErr != nil {
			return domain.Run{}, failErr
		}

		return domain.Run{}, domain.ErrProviderFailed
	}

	return h.runs.Complete(ctx, created.ID, output)
}

func valueOrEmpty(value *string) string {
	if value == nil {
		return ""
	}

	return *value
}
