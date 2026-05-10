package middleware

import (
	"net/http"

	"ai-assistants-catalog/internal/core/security"
)

func RequireAdminMiddleware() func(http.Handler) http.Handler {
	return RequireRoleMiddleware(security.RoleAdmin)
}
