package domain

import (
	"strings"
	"time"
)

type Category struct {
	ID          string
	Name        string
	Description *string
	CreatedAt   *time.Time
}

func NewCategory(name string, description *string) (Category, error) {
	normalizedName := strings.TrimSpace(name)
	if normalizedName == "" {
		return Category{}, ErrInvalidName
	}

	return Category{
		Name:        normalizedName,
		Description: normalizeDescription(description),
	}, nil
}

func ReconstituteCategory(id string, name string, description *string, createdAt *time.Time) Category {
	return Category{
		ID:          id,
		Name:        name,
		Description: description,
		CreatedAt:   createdAt,
	}
}

func normalizeDescription(description *string) *string {
	if description == nil {
		return nil
	}

	value := strings.TrimSpace(*description)
	if value == "" {
		return nil
	}

	return &value
}
