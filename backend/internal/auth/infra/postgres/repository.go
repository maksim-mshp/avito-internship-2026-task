package postgres

import (
	"context"
	"errors"
	"time"

	"ai-assistants-catalog/internal/auth/domain"
	corepostgres "ai-assistants-catalog/internal/core/postgres"

	"github.com/jackc/pgx/v5"
)

type Repository struct {
	db corepostgres.DBTX
}

func NewRepository(db corepostgres.DBTX) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetByRole(ctx context.Context, role domain.Role) (domain.User, error) {
	var (
		id        string
		email     string
		roleValue string
		createdAt time.Time
	)

	err := r.db.QueryRow(ctx, getByRoleQuery, role.UserID(), role.String()).Scan(&id, &email, &roleValue, &createdAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.ErrUserNotFound
		}

		return domain.User{}, err
	}

	foundRole, err := domain.NewRole(roleValue)
	if err != nil {
		return domain.User{}, err
	}

	return domain.ReconstituteUser(id, email, foundRole, &createdAt), nil
}
