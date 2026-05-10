package core

import (
	"context"
	"net/http"
	"time"

	authApp "ai-assistants-catalog/internal/auth/app/handlers"
	authV1HTTP "ai-assistants-catalog/internal/auth/infra/http/v1"
	"ai-assistants-catalog/internal/core/config"
	corehttp "ai-assistants-catalog/internal/core/http"
	"ai-assistants-catalog/internal/core/http/middleware"
)

type App struct {
	Config *config.Config
	Server *http.Server
}

func Start(cfg *config.Config) (*App, error) {
	mux := http.NewServeMux()

	corehttp.RegisterRoutes(mux)

	authHandlers := authApp.BuildHandlers(cfg.JWTToken)
	authHTTPHandler := authV1HTTP.NewHTTPHandler(authHandlers)
	authV1HTTP.RegisterRoutes(mux, authHTTPHandler)

	handler := middleware.RecoverMiddleware(mux)
	handler = middleware.LoggingMiddleware(handler)

	server, err := corehttp.NewServer(cfg.Port, handler)
	if err != nil {
		return nil, err
	}

	return &App{
		Config: cfg,
		Server: server,
	}, nil
}

func (a *App) Stop(ctx context.Context) error {
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return a.Server.Shutdown(shutdownCtx)
}
