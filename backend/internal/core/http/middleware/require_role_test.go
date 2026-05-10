package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"ai-assistants-catalog/internal/core/security"
)

func TestRequireRoleMiddleware(t *testing.T) {
	t.Run("missing claims", func(t *testing.T) {
		middleware := RequireRoleMiddleware(security.RoleAdmin)
		handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatalf("handler should not be called")
		}))

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		recorder := httptest.NewRecorder()

		handler.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusUnauthorized {
			t.Fatalf("expected status 401, got %d", recorder.Code)
		}
	})

	t.Run("forbidden role", func(t *testing.T) {
		middleware := RequireRoleMiddleware(security.RoleAdmin)
		handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatalf("handler should not be called")
		}))

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		ctx := security.WithClaims(req.Context(), security.Claims{
			UserID: "user-id",
			Role:   security.RoleUser,
		})
		recorder := httptest.NewRecorder()

		handler.ServeHTTP(recorder, req.WithContext(ctx))

		if recorder.Code != http.StatusForbidden {
			t.Fatalf("expected status 403, got %d", recorder.Code)
		}
	})

	t.Run("allowed role", func(t *testing.T) {
		middleware := RequireRoleMiddleware(security.RoleAdmin, security.RoleUser)
		handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		}))

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		ctx := security.WithClaims(req.Context(), security.Claims{
			UserID: "user-id",
			Role:   security.RoleUser,
		})
		recorder := httptest.NewRecorder()

		handler.ServeHTTP(recorder, req.WithContext(ctx))

		if recorder.Code != http.StatusNoContent {
			t.Fatalf("expected status 204, got %d", recorder.Code)
		}
	})
}
