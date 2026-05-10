package handlers

import (
	"context"
	"errors"
	"testing"
	"time"

	"ai-assistants-catalog/internal/auth/app"
	"ai-assistants-catalog/internal/auth/domain"
	"ai-assistants-catalog/internal/core/security"
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

func TestDummyLoginHandlerHandle(t *testing.T) {
	handler := NewDummyLoginHandler("secret", fakeRepository{})

	result, err := handler.Handle(context.Background(), app.DummyLoginCommand{Role: "admin"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	claims, err := security.ParseJWT("secret", result.Token)
	if err != nil {
		t.Fatalf("expected valid token, got %v", err)
	}
	if claims.Role != "admin" {
		t.Fatalf("expected admin role, got %s", claims.Role)
	}
	if claims.UserID != "ca2a62f3-c998-4050-96c1-0c0f62cf6568" {
		t.Fatalf("unexpected admin user id: %s", claims.UserID)
	}
	if result.User.Email != "admin@example.com" {
		t.Fatalf("unexpected admin email: %s", result.User.Email)
	}

	_, err = handler.Handle(context.Background(), app.DummyLoginCommand{Role: "guest"})
	if !errors.Is(err, domain.ErrInvalidRole) {
		t.Fatalf("expected ErrInvalidRole, got %v", err)
	}
}
