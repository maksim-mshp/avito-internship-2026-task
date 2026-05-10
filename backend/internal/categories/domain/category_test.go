package domain

import (
	"errors"
	"testing"
)

func TestNewCategory(t *testing.T) {
	description := "  helpful category  "

	category, err := NewCategory("  Food  ", &description)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if category.Name != "Food" {
		t.Fatalf("unexpected name: %s", category.Name)
	}
	if category.Description == nil || *category.Description != "helpful category" {
		t.Fatalf("unexpected description: %v", category.Description)
	}
}

func TestNewCategoryEmptyName(t *testing.T) {
	_, err := NewCategory("   ", nil)
	if !errors.Is(err, ErrInvalidName) {
		t.Fatalf("expected ErrInvalidName, got %v", err)
	}
}

func TestNewCategoryEmptyDescription(t *testing.T) {
	description := "   "

	category, err := NewCategory("Food", &description)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if category.Description != nil {
		t.Fatalf("expected nil description, got %v", category.Description)
	}
}
