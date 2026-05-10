package postgres

import (
	"context"
	"time"

	"ai-assistants-catalog/internal/core/postgres"
	"ai-assistants-catalog/internal/runs/app"
	"ai-assistants-catalog/internal/runs/domain"
)

type Repository struct {
	db postgres.DBTX
}

type rowScanner interface {
	Scan(dest ...any) error
}

func NewRepository(db postgres.DBTX) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreatePending(ctx context.Context, run domain.Run) (domain.Run, error) {
	return scanRun(r.db.QueryRow(
		ctx,
		createPendingQuery,
		run.AssistantID,
		run.UserID,
		run.Model,
		run.UserPrompt,
	))
}

func (r *Repository) Complete(ctx context.Context, id string, output string) (domain.Run, error) {
	return scanRun(r.db.QueryRow(ctx, completeQuery, id, output))
}

func (r *Repository) Fail(ctx context.Context, id string, message string) (domain.Run, error) {
	return scanRun(r.db.QueryRow(ctx, failQuery, id, message))
}

func (r *Repository) ListMy(ctx context.Context, query app.ListMyQuery) (app.ListResult, error) {
	offset := (query.Page - 1) * query.PageSize
	return r.list(ctx, listMyQuery, query.UserID, nullableString(query.Status), query.PageSize, offset, query.Page, query.PageSize)
}

func (r *Repository) ListAdmin(ctx context.Context, query app.ListAdminQuery) (app.ListResult, error) {
	offset := (query.Page - 1) * query.PageSize
	return r.list(ctx, listAdminQuery, nullableString(query.AssistantID), nullableString(query.Status), query.PageSize, offset, query.Page, query.PageSize)
}

func (r *Repository) list(
	ctx context.Context,
	sql string,
	filterA any,
	filterB any,
	limit int,
	offset int,
	page int,
	pageSize int,
) (app.ListResult, error) {
	rows, err := r.db.Query(ctx, sql, filterA, filterB, limit, offset)
	if err != nil {
		return app.ListResult{}, err
	}

	defer rows.Close()

	runs := make([]domain.Run, 0)
	total := 0

	for rows.Next() {
		run, rowTotal, scanErr := scanRunWithTotal(rows)
		if scanErr != nil {
			return app.ListResult{}, scanErr
		}

		runs = append(runs, run)
		total = rowTotal
	}

	if err = rows.Err(); err != nil {
		return app.ListResult{}, err
	}

	return app.ListResult{
		Runs:     runs,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}, nil
}

func scanRun(row rowScanner) (domain.Run, error) {
	var (
		id            string
		assistantID   string
		assistantName *string
		categoryID    *string
		categoryName  *string
		userID        string
		model         string
		userPrompt    string
		output        *string
		status        string
		runError      *string
		createdAt     time.Time
	)

	if err := row.Scan(
		&id,
		&assistantID,
		&assistantName,
		&categoryID,
		&categoryName,
		&userID,
		&model,
		&userPrompt,
		&output,
		&status,
		&runError,
		&createdAt,
	); err != nil {
		return domain.Run{}, err
	}

	return domain.ReconstituteRun(
		id,
		assistantID,
		assistantName,
		categoryID,
		categoryName,
		userID,
		model,
		userPrompt,
		output,
		status,
		runError,
		&createdAt,
	), nil
}

func scanRunWithTotal(row rowScanner) (domain.Run, int, error) {
	var (
		id            string
		assistantID   string
		assistantName *string
		categoryID    *string
		categoryName  *string
		userID        string
		model         string
		userPrompt    string
		output        *string
		status        string
		runError      *string
		createdAt     time.Time
		total         int
	)

	if err := row.Scan(
		&id,
		&assistantID,
		&assistantName,
		&categoryID,
		&categoryName,
		&userID,
		&model,
		&userPrompt,
		&output,
		&status,
		&runError,
		&createdAt,
		&total,
	); err != nil {
		return domain.Run{}, 0, err
	}

	run := domain.ReconstituteRun(
		id,
		assistantID,
		assistantName,
		categoryID,
		categoryName,
		userID,
		model,
		userPrompt,
		output,
		status,
		runError,
		&createdAt,
	)

	return run, total, nil
}

func nullableString(value *string) any {
	if value == nil {
		return nil
	}

	return *value
}
