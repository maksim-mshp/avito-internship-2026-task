package middleware

import (
	"net/http"

	corehttp "ai-assistants-catalog/internal/core/http"
	"ai-assistants-catalog/internal/core/security"
)

func RequireRoleMiddleware(roles ...string) func(http.Handler) http.Handler {
	allowed := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		allowed[role] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := security.ClaimsFromContext(r.Context())
			if !ok {
				corehttp.RespondError(w, corehttp.ErrUnauthorized)
				return
			}

			if _, ok = allowed[claims.Role]; !ok {
				corehttp.RespondError(w, corehttp.ErrForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
