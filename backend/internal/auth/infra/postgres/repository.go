package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"ai-assistants-catalog/internal/auth/domain"
	corepostgres "ai-assistants-catalog/internal/core/postgres"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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

func (r *Repository) CreateUser(ctx context.Context, user domain.User, passwordHash string) (domain.User, error) {
	var (
		id        string
		email     string
		roleValue string
		createdAt time.Time
	)

	err := r.db.QueryRow(ctx, createUserQuery, user.Email, user.Role.String(), passwordHash).
		Scan(&id, &email, &roleValue, &createdAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.User{}, domain.ErrEmailTaken
		}

		return domain.User{}, fmt.Errorf("create user: %w", err)
	}

	role, err := domain.NewRole(roleValue)
	if err != nil {
		return domain.User{}, err
	}

	return domain.ReconstituteUser(id, email, role, &createdAt), nil
}

func (r *Repository) GetAuthUserByEmail(ctx context.Context, email string) (domain.AuthUser, error) {
	var (
		id           string
		emailValue   string
		roleValue    string
		passwordHash string
		createdAt    time.Time
	)

	err := r.db.QueryRow(ctx, getAuthUserByEmailQuery, email).Scan(
		&id,
		&emailValue,
		&roleValue,
		&createdAt,
		&passwordHash,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.AuthUser{}, domain.ErrInvalidCredentials
		}

		return domain.AuthUser{}, fmt.Errorf("get auth user by email: %w", err)
	}

	role, err := domain.NewRole(roleValue)
	if err != nil {
		return domain.AuthUser{}, err
	}

	return domain.AuthUser{
		User:         domain.ReconstituteUser(id, emailValue, role, &createdAt),
		PasswordHash: passwordHash,
	}, nil
}
