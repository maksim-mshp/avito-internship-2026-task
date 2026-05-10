package core

import (
	"context"
	"net/http"
	"time"

	assistantsApp "ai-assistants-catalog/internal/assistants/app/handlers"
	assistantsV1HTTP "ai-assistants-catalog/internal/assistants/infra/http/v1"
	assistantsPostgres "ai-assistants-catalog/internal/assistants/infra/postgres"
	authApp "ai-assistants-catalog/internal/auth/app/handlers"
	authV1HTTP "ai-assistants-catalog/internal/auth/infra/http/v1"
	categoriesApp "ai-assistants-catalog/internal/categories/app/handlers"
	categoriesV1HTTP "ai-assistants-catalog/internal/categories/infra/http/v1"
	categoriesPostgres "ai-assistants-catalog/internal/categories/infra/postgres"
	"ai-assistants-catalog/internal/core/config"
	corehttp "ai-assistants-catalog/internal/core/http"
	"ai-assistants-catalog/internal/core/http/middleware"
	"ai-assistants-catalog/internal/core/postgres"
	runsApp "ai-assistants-catalog/internal/runs/app/handlers"
	runsV1HTTP "ai-assistants-catalog/internal/runs/infra/http/v1"
	runsMockLLM "ai-assistants-catalog/internal/runs/infra/llm/mock"
	runsPostgres "ai-assistants-catalog/internal/runs/infra/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	Config   *config.Config
	Database *pgxpool.Pool
	Server   *http.Server
}

func Start(cfg *config.Config) (*App, error) {
	db, err := postgres.NewPostgres(cfg.Database)
	if err != nil {
		return nil, err
	}

	needCloseDB := true
	defer func() {
		if needCloseDB {
			db.Close()
		}
	}()

	if err = postgres.RunMigrations(cfg.Database); err != nil {
		return nil, err
	}

	mux := http.NewServeMux()

	corehttp.RegisterRoutes(mux)
	authMW := middleware.RequireAuthMiddleware(cfg.JWTToken)
	adminMW := middleware.RequireAdminMiddleware()

	authHandlers := authApp.BuildHandlers(cfg.JWTToken)
	authHTTPHandler := authV1HTTP.NewHTTPHandler(authHandlers)
	authV1HTTP.RegisterRoutes(mux, authHTTPHandler)

	categoriesRepo := categoriesPostgres.NewRepository(db)
	categoriesHandlers := categoriesApp.BuildHandlers(categoriesRepo)
	categoriesHTTPHandler := categoriesV1HTTP.NewHTTPHandler(categoriesHandlers)
	categoriesV1HTTP.RegisterRoutes(mux, categoriesHTTPHandler, authMW, adminMW)

	assistantsRepo := assistantsPostgres.NewRepository(db)
	assistantsHandlers := assistantsApp.BuildHandlers(assistantsRepo)
	assistantsHTTPHandler := assistantsV1HTTP.NewHTTPHandler(assistantsHandlers)
	assistantsV1HTTP.RegisterRoutes(mux, assistantsHTTPHandler, authMW, adminMW)

	runsRepo := runsPostgres.NewRepository(db)
	llmProvider := runsMockLLM.NewProvider()
	runsHandlers := runsApp.BuildHandlers(runsRepo, assistantsRepo, llmProvider)
	runsHTTPHandler := runsV1HTTP.NewHTTPHandler(runsHandlers)
	runsV1HTTP.RegisterRoutes(mux, runsHTTPHandler, authMW, adminMW)

	handler := middleware.RecoverMiddleware(mux)
	handler = middleware.LoggingMiddleware(handler)

	server, err := corehttp.NewServer(cfg.Port, handler)
	if err != nil {
		return nil, err
	}

	needCloseDB = false

	return &App{
		Config:   cfg,
		Database: db,
		Server:   server,
	}, nil
}

func (a *App) Stop(ctx context.Context) error {
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := a.Server.Shutdown(shutdownCtx)
	a.Database.Close()

	return err
}
