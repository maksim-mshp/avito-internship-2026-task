package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"time"
)

type APIError struct {
	StatusCode int    `json:"-"`
	Code       string `json:"-"`
	Message    string `json:"-"`
}

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

func Respond(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if data == nil {
		return
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("failed to encode json response: %v", err)
	}
}

func RespondError(w http.ResponseWriter, apiErr APIError) {
	Respond(w, apiErr.StatusCode, ErrorResponse{
		Error: ErrorBody{
			Code:    apiErr.Code,
			Message: apiErr.Message,
		},
	})
}

func ParseJSONBody(r *http.Request, data any) *APIError {
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("failed to close body reader: %v", err)
		}
	}()

	if contentType := r.Header.Get("Content-Type"); contentType != "" {
		mediaType, _, err := mime.ParseMediaType(contentType)
		if err != nil || mediaType != "application/json" {
			return &ErrInvalidRequest
		}
	}

	const maxBodyBytes = 1024 * 1024
	r.Body = http.MaxBytesReader(nil, r.Body, maxBodyBytes)

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(data); err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			return &ErrInvalidRequest
		}

		return &ErrInvalidRequest
	}

	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return &ErrInvalidRequest
	}

	return nil
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
