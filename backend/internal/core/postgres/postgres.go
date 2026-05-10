package postgres

import (
	"context"
	"fmt"
	"net/url"

	"ai-assistants-catalog/internal/core/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func MakeConnectionString(cfg config.Database) string {
	dsn := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.User, cfg.Password),
		Host:   fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Path:   cfg.Database,
	}

	return dsn.String()
}

func NewPostgres(cfg config.Database) (*pgxpool.Pool, error) {
	dsn := MakeConnectionString(cfg)
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	if err = pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	return pool, nil
}
