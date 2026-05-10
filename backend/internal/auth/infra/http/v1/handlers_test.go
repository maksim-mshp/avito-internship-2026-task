package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	authHandlers "ai-assistants-catalog/internal/auth/app/handlers"
	"ai-assistants-catalog/internal/auth/domain"
	corehttp "ai-assistants-catalog/internal/core/http"
)

type fakeRepository struct{}

const testPasswordHash = "$2a$10$hfcLxclybGOmg8CIJeuZi.3qIHtMKLdw2bHne4ZRDbARnai0Sgh96"

func (r fakeRepository) GetByRole(_ context.Context, role domain.Role) (domain.User, error) {
	createdAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	return domain.ReconstituteUser(role.UserID(), role.Email(), role, &createdAt), nil
}

func (r fakeRepository) CreateUser(_ context.Context, user domain.User, _ string) (domain.User, error) {
	createdAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	return domain.ReconstituteUser("0fa9ed62-6138-4794-b567-3813a8d1a4fb", user.Email, user.Role, &createdAt), nil
}

func (r fakeRepository) GetAuthUserByEmail(_ context.Context, email string) (domain.AuthUser, error) {
	createdAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	return domain.AuthUser{
		User:         domain.ReconstituteUser("0fa9ed62-6138-4794-b567-3813a8d1a4fb", email, domain.RoleUser, &createdAt),
		PasswordHash: testPasswordHash,
	}, nil
}

func TestRegister(t *testing.T) {
	handler := NewHTTPHandler(authHandlers.BuildHandlers("secret", fakeRepository{}))

	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(`{"email":"user@example.com","password":"password","role":"user"}`))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	handler.Register(recorder, req)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", recorder.Code)
	}

	var response RegisterResponse
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if response.User.Email != "user@example.com" {
		t.Fatalf("unexpected email: %s", response.User.Email)
	}
}

func TestLogin(t *testing.T) {
	handler := NewHTTPHandler(authHandlers.BuildHandlers("secret", fakeRepository{}))

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"email":"user@example.com","password":"password"}`))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	handler.Login(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}

	var response TokenResponse
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if response.Token == "" {
		t.Fatalf("expected token")
	}
	if response.User.Email != "user@example.com" {
		t.Fatalf("unexpected email: %s", response.User.Email)
	}
}

func TestDummyLogin(t *testing.T) {
	handler := NewHTTPHandler(authHandlers.BuildHandlers("secret", fakeRepository{}))

	req := httptest.NewRequest(http.MethodPost, "/dummyLogin", strings.NewReader(`{"role":"user"}`))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	handler.DummyLogin(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}

	var response TokenResponse
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if response.Token == "" {
		t.Fatalf("expected token")
	}
	if response.User.ID != "44c75af3-eca3-4867-85fc-b8245eaafa3a" {
		t.Fatalf("unexpected user id: %s", response.User.ID)
	}
	if response.User.Role != "user" {
		t.Fatalf("unexpected role: %s", response.User.Role)
	}
}

func TestDummyLoginInvalidRole(t *testing.T) {
	handler := NewHTTPHandler(authHandlers.BuildHandlers("secret", fakeRepository{}))

	req := httptest.NewRequest(http.MethodPost, "/dummyLogin", strings.NewReader(`{"role":"guest"}`))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	handler.DummyLogin(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}

	var response corehttp.ErrorResponse
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if response.Error.Code != "INVALID_REQUEST" {
		t.Fatalf("unexpected error code: %s", response.Error.Code)
	}
}
