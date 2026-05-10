package middleware

import (
	"net/http"
	"strings"

	corehttp "ai-assistants-catalog/internal/core/http"
	"ai-assistants-catalog/internal/core/security"
)

func RequireAuthMiddleware(secretKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, ok := parseBearerToken(r.Header.Get("Authorization"))
			if !ok {
				corehttp.RespondError(w, corehttp.ErrUnauthorized)
				return
			}

			claims, err := security.ParseJWT(secretKey, token)
			if err != nil {
				corehttp.RespondError(w, corehttp.ErrUnauthorized)
				return
			}

			next.ServeHTTP(w, r.WithContext(security.WithClaims(r.Context(), claims)))
		})
	}
}

func parseBearerToken(value string) (string, bool) {
	if value == "" {
		return "", false
	}

	token, ok := strings.CutPrefix(value, "Bearer ")
	if !ok || token == "" {
		return "", false
	}

	return token, true
}
