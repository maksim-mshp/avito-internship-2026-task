package middleware

import (
	"net/http"

	"ai-assistants-catalog/internal/core/security"
)

func RequireUserMiddleware() func(http.Handler) http.Handler {
	return RequireRoleMiddleware(security.RoleUser)
}
