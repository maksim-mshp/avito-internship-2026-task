package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"ai-assistants-catalog/internal/core"
	"ai-assistants-catalog/internal/core/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	app, err := core.Start(cfg)
	if err != nil {
		log.Fatalf("failed to start app: %v", err)
	}

	serverErr := make(chan error, 1)

	go func() {
		log.Printf("server started on :%d", cfg.Port)
		serverErr <- app.Server.ListenAndServe()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err = <-serverErr:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("server error: %v", err)
		}
	case <-stop:
		log.Println("shutdown signal received")
	}

	if err = app.Stop(context.Background()); err != nil {
		log.Printf("failed to stop app: %v", err)
	}

	log.Println("application stopped")
}
