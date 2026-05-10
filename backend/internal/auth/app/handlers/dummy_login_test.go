package handlers

import (
	"context"
	"errors"
	"testing"

	"ai-assistants-catalog/internal/auth/app"
	"ai-assistants-catalog/internal/auth/domain"
	"ai-assistants-catalog/internal/core/security"
)

func TestDummyLoginHandlerHandle(t *testing.T) {
	handler := NewDummyLoginHandler("secret")

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
