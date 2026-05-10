package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func Respond(w http.ResponseWriter, statusCode int, data any) {
	if data == nil {
		w.WriteHeader(statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("failed to encode json response: %v", err)
	}
}

func NewServer(port int, handler http.Handler) (*http.Server, error) {
	if port <= 0 {
		return nil, fmt.Errorf("port must be positive")
	}

	return &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}, nil
}
