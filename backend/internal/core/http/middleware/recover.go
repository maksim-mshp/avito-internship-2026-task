package middleware

import (
	"log"
	"net/http"

	corehttp "ai-assistants-catalog/internal/core/http"
)

func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic recovered: %v", err)
				corehttp.Respond(w, http.StatusInternalServerError, corehttp.ErrorResponse{
					Error: "INTERNAL_ERROR",
				})
			}
		}()

		next.ServeHTTP(w, r)
	})
}
