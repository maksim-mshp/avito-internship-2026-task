package app

import (
	"context"

	"ai-assistants-catalog/internal/assistants/domain"
)

type Repository interface {
	List(ctx context.Context, query ListQuery) (ListResult, error)
	GetByID(ctx context.Context, id string, includeInactive bool, userID string) (domain.Assistant, error)
	Create(ctx context.Context, assistant domain.Assistant) (domain.Assistant, error)
	Update(ctx context.Context, assistant domain.Assistant) (domain.Assistant, error)
	AddFavorite(ctx context.Context, userID string, assistantID string, includeInactive bool) error
	RemoveFavorite(ctx context.Context, userID string, assistantID string) error
}
