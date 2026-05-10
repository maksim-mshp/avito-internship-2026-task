package handlers

import (
	"context"
	"errors"
	"testing"

	"ai-assistants-catalog/internal/auth/app"
	"ai-assistants-catalog/internal/auth/domain"
)

func TestRegisterHandlerHandle(t *testing.T) {
	handler := NewRegisterHandler(fakeRepository{})

	user, err := handler.Handle(context.Background(), app.RegisterCommand{
		Email:    " USER@Example.COM ",
		Password: "password",
		Role:     "user",
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if user.Email != "user@example.com" {
		t.Fatalf("unexpected email: %s", user.Email)
	}
	if user.Role != domain.RoleUser {
		t.Fatalf("unexpected role: %s", user.Role)
	}
}

func TestRegisterHandlerHandleInvalidEmail(t *testing.T) {
	handler := NewRegisterHandler(fakeRepository{})

	_, err := handler.Handle(context.Background(), app.RegisterCommand{
		Email:    "bad",
		Password: "password",
		Role:     "user",
	})
	if !errors.Is(err, domain.ErrInvalidEmail) {
		t.Fatalf("expected ErrInvalidEmail, got %v", err)
	}
}
