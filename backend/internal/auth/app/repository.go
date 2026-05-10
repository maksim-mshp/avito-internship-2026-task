package app

import (
	"context"

	"ai-assistants-catalog/internal/auth/domain"
)

type Repository interface {
	GetByRole(ctx context.Context, role domain.Role) (domain.User, error)
	CreateUser(ctx context.Context, user domain.User, passwordHash string) (domain.User, error)
	GetAuthUserByEmail(ctx context.Context, email string) (domain.AuthUser, error)
}
