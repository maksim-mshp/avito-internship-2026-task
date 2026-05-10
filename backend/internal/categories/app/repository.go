package app

import (
	"context"

	"ai-assistants-catalog/internal/categories/domain"
)

type Repository interface {
	GetAll(ctx context.Context) ([]domain.Category, error)
	Create(ctx context.Context, category domain.Category) (domain.Category, error)
}
