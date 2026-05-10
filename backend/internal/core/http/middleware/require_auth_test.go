package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"ai-assistants-catalog/internal/core/security"
)

func TestRequireAuthMiddleware(t *testing.T) {
	t.Run("missing token", func(t *testing.T) {
		middleware := RequireAuthMiddleware("secret")
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

	t.Run("invalid token", func(t *testing.T) {
		middleware := RequireAuthMiddleware("secret")
		handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatalf("handler should not be called")
		}))

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer bad-token")
		recorder := httptest.NewRecorder()

		handler.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusUnauthorized {
			t.Fatalf("expected status 401, got %d", recorder.Code)
		}
	})

	t.Run("valid token", func(t *testing.T) {
		token, err := security.GenerateJWT("secret", "user-id", security.RoleUser)
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}

		middleware := RequireAuthMiddleware("secret")
		handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := security.ClaimsFromContext(r.Context())
			if !ok {
				t.Fatalf("expected claims in context")
			}
			if claims.UserID != "user-id" || claims.Role != security.RoleUser {
				t.Fatalf("unexpected claims: %#v", claims)
			}

			w.WriteHeader(http.StatusNoContent)
		}))

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		recorder := httptest.NewRecorder()

		handler.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusNoContent {
			t.Fatalf("expected status 204, got %d", recorder.Code)
		}
	})
}
