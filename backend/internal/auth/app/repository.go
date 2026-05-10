package app

import (
	"context"

	"ai-assistants-catalog/internal/auth/domain"
)

type Repository interface {
	GetByRole(ctx context.Context, role domain.Role) (domain.User, error)
}
