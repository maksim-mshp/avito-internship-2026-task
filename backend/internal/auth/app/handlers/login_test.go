package handlers

import (
	"context"
	"errors"
	"testing"

	"ai-assistants-catalog/internal/auth/app"
	"ai-assistants-catalog/internal/auth/domain"
	"ai-assistants-catalog/internal/core/security"
)

func TestLoginHandlerHandle(t *testing.T) {
	handler := NewLoginHandler("secret", fakeRepository{})

	result, err := handler.Handle(context.Background(), app.LoginCommand{
		Email:    "USER@example.com",
		Password: "password",
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if result.User.Email != "user@example.com" {
		t.Fatalf("unexpected email: %s", result.User.Email)
	}

	claims, err := security.ParseJWT("secret", result.Token)
	if err != nil {
		t.Fatalf("expected valid token, got %v", err)
	}
	if claims.Role != "user" {
		t.Fatalf("expected user role, got %s", claims.Role)
	}
}

func TestLoginHandlerHandleInvalidPassword(t *testing.T) {
	handler := NewLoginHandler("secret", fakeRepository{})

	_, err := handler.Handle(context.Background(), app.LoginCommand{
		Email:    "user@example.com",
		Password: "wrong",
	})
	if !errors.Is(err, domain.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}
