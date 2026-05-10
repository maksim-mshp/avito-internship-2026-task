package postgres

import (
	"context"
	"time"

	"ai-assistants-catalog/internal/categories/domain"
	corepostgres "ai-assistants-catalog/internal/core/postgres"

	"github.com/jackc/pgx/v5"
)

type Repository struct {
	db corepostgres.DBTX
}

type rowScanner interface {
	Scan(dest ...any) error
}

func NewRepository(db corepostgres.DBTX) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAll(ctx context.Context) ([]domain.Category, error) {
	rows, err := r.db.Query(ctx, getAllQuery)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, scanCollectableCategory)
}

func (r *Repository) Create(ctx context.Context, category domain.Category) (domain.Category, error) {
	row := r.db.QueryRow(ctx, createQuery, category.Name, category.Description)
	return scanCategory(row)
}

func scanCollectableCategory(row pgx.CollectableRow) (domain.Category, error) {
	return scanCategory(row)
}

func scanCategory(row rowScanner) (domain.Category, error) {
	var (
		id          string
		name        string
		description *string
		createdAt   time.Time
	)

	if err := row.Scan(&id, &name, &description, &createdAt); err != nil {
		return domain.Category{}, err
	}

	return domain.ReconstituteCategory(id, name, description, &createdAt), nil
}
