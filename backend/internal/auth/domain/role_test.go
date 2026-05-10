package domain

import (
	"errors"
	"testing"
)

func TestNewRole(t *testing.T) {
	admin, err := NewRole("admin")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if admin.UserID() != "ca2a62f3-c998-4050-96c1-0c0f62cf6568" {
		t.Fatalf("unexpected admin id: %s", admin.UserID())
	}
	if admin.Email() != "admin@example.com" {
		t.Fatalf("unexpected admin email: %s", admin.Email())
	}

	user, err := NewRole("user")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if user.UserID() != "44c75af3-eca3-4867-85fc-b8245eaafa3a" {
		t.Fatalf("unexpected user id: %s", user.UserID())
	}
}

func TestNewRoleInvalid(t *testing.T) {
	_, err := NewRole("guest")
	if !errors.Is(err, ErrInvalidRole) {
		t.Fatalf("expected ErrInvalidRole, got %v", err)
	}
}
