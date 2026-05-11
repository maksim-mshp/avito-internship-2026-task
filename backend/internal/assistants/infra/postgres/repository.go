package postgres

import (
	"context"
	"errors"
	"time"

	"ai-assistants-catalog/internal/assistants/app"
	"ai-assistants-catalog/internal/assistants/domain"
	corepostgres "ai-assistants-catalog/internal/core/postgres"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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

func (r *Repository) List(ctx context.Context, query app.ListQuery) (app.ListResult, error) {
	offset := (query.Page - 1) * query.PageSize

	rows, err := r.db.Query(
		ctx,
		listQuery,
		nullableString(query.CategoryID),
		nullableString(query.Search),
		nullableString(query.Tag),
		nullableUserID(query.UserID),
		query.IncludeInactive,
		query.FavoritesOnly,
		query.PageSize,
		offset,
	)
	if err != nil {
		return app.ListResult{}, err
	}

	defer rows.Close()

	assistants := make([]domain.Assistant, 0)
	total := 0

	for rows.Next() {
		assistant, rowTotal, scanErr := scanAssistantWithTotal(rows)
		if scanErr != nil {
			return app.ListResult{}, scanErr
		}

		assistants = append(assistants, assistant)
		total = rowTotal
	}

	if err = rows.Err(); err != nil {
		return app.ListResult{}, err
	}

	return app.ListResult{
		Assistants: assistants,
		Page:       query.Page,
		PageSize:   query.PageSize,
		Total:      total,
	}, nil
}

func (r *Repository) GetByID(ctx context.Context, id string, includeInactive bool, userID string) (domain.Assistant, error) {
	assistant, err := scanAssistant(r.db.QueryRow(ctx, getByIDQuery, id, includeInactive, nullableUserID(userID)))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Assistant{}, domain.ErrNotFound
		}

		return domain.Assistant{}, err
	}

	return assistant, nil
}

func (r *Repository) Create(ctx context.Context, assistant domain.Assistant) (domain.Assistant, error) {
	created, err := scanAssistant(r.db.QueryRow(
		ctx,
		createQuery,
		assistant.CategoryID,
		assistant.Name,
		assistant.Description,
		assistant.Model,
		assistant.SystemPrompt,
		assistant.ExampleUserPrompt,
		assistant.Tags,
		assistant.IsActive,
	))
	if err != nil {
		return domain.Assistant{}, mapDatabaseError(err)
	}

	return created, nil
}

func (r *Repository) Update(ctx context.Context, assistant domain.Assistant) (domain.Assistant, error) {
	updated, err := scanAssistant(r.db.QueryRow(
		ctx,
		updateQuery,
		assistant.ID,
		assistant.CategoryID,
		assistant.Name,
		assistant.Description,
		assistant.Model,
		assistant.SystemPrompt,
		assistant.ExampleUserPrompt,
		assistant.Tags,
		assistant.IsActive,
	))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Assistant{}, domain.ErrNotFound
		}

		return domain.Assistant{}, mapDatabaseError(err)
	}

	return updated, nil
}

func (r *Repository) AddFavorite(ctx context.Context, userID string, assistantID string, includeInactive bool) error {
	var id string
	err := r.db.QueryRow(ctx, addFavoriteQuery, userID, assistantID, includeInactive).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrNotFound
		}

		return err
	}

	return nil
}

func (r *Repository) RemoveFavorite(ctx context.Context, userID string, assistantID string) error {
	_, err := r.db.Exec(ctx, removeFavoriteQuery, userID, assistantID)
	return err
}

func scanAssistant(row rowScanner) (domain.Assistant, error) {
	var (
		id                string
		categoryID        string
		categoryName      *string
		name              string
		description       string
		model             string
		systemPrompt      string
		exampleUserPrompt *string
		tags              []string
		isFavorite        bool
		isActive          bool
		createdAt         time.Time
		updatedAt         time.Time
	)

	if err := row.Scan(
		&id,
		&categoryID,
		&categoryName,
		&name,
		&description,
		&model,
		&systemPrompt,
		&exampleUserPrompt,
		&tags,
		&isFavorite,
		&isActive,
		&createdAt,
		&updatedAt,
	); err != nil {
		return domain.Assistant{}, err
	}

	return domain.ReconstituteAssistant(
		id,
		categoryID,
		categoryName,
		name,
		description,
		model,
		systemPrompt,
		exampleUserPrompt,
		tags,
		isFavorite,
		isActive,
		&createdAt,
		&updatedAt,
	), nil
}

func scanAssistantWithTotal(row rowScanner) (domain.Assistant, int, error) {
	var (
		id                string
		categoryID        string
		categoryName      *string
		name              string
		description       string
		model             string
		systemPrompt      string
		exampleUserPrompt *string
		tags              []string
		isFavorite        bool
		isActive          bool
		createdAt         time.Time
		updatedAt         time.Time
		total             int
	)

	if err := row.Scan(
		&id,
		&categoryID,
		&categoryName,
		&name,
		&description,
		&model,
		&systemPrompt,
		&exampleUserPrompt,
		&tags,
		&isFavorite,
		&isActive,
		&createdAt,
		&updatedAt,
		&total,
	); err != nil {
		return domain.Assistant{}, 0, err
	}

	assistant := domain.ReconstituteAssistant(
		id,
		categoryID,
		categoryName,
		name,
		description,
		model,
		systemPrompt,
		exampleUserPrompt,
		tags,
		isFavorite,
		isActive,
		&createdAt,
		&updatedAt,
	)

	return assistant, total, nil
}

func nullableString(value *string) any {
	if value == nil {
		return nil
	}

	return *value
}

func nullableUserID(value string) any {
	if value == "" {
		return nil
	}

	return value
}

func mapDatabaseError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23503" {
		return domain.ErrCategoryNotFound
	}

	return err
}
