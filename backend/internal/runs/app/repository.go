package app

import (
	"context"

	assistantdomain "ai-assistants-catalog/internal/assistants/domain"
	"ai-assistants-catalog/internal/runs/domain"
)

type AssistantRepository interface {
	GetByID(ctx context.Context, id string, includeInactive bool, userID string) (assistantdomain.Assistant, error)
}

type RunRepository interface {
	CreatePending(ctx context.Context, run domain.Run) (domain.Run, error)
	Complete(ctx context.Context, id string, output string) (domain.Run, error)
	Fail(ctx context.Context, id string, message string) (domain.Run, error)
	ListMy(ctx context.Context, query ListMyQuery) (ListResult, error)
	ListAdmin(ctx context.Context, query ListAdminQuery) (ListResult, error)
}
