package domain

import (
	"errors"
	"testing"
)

func TestNormalizeEmail(t *testing.T) {
	email, err := NormalizeEmail(" USER@example.COM ")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if email != "user@example.com" {
		t.Fatalf("unexpected email: %s", email)
	}
}

func TestNormalizeEmailInvalid(t *testing.T) {
	_, err := NormalizeEmail("bad")
	if !errors.Is(err, ErrInvalidEmail) {
		t.Fatalf("expected ErrInvalidEmail, got %v", err)
	}
}

func TestValidatePassword(t *testing.T) {
	if err := ValidatePassword("password"); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	err := ValidatePassword(" ")
	if !errors.Is(err, ErrInvalidPassword) {
		t.Fatalf("expected ErrInvalidPassword, got %v", err)
	}
}
