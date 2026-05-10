package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Assistant struct {
	ID                string
	CategoryID        string
	CategoryName      *string
	Name              string
	Description       string
	Model             string
	SystemPrompt      string
	ExampleUserPrompt *string
	IsActive          bool
	CreatedAt         *time.Time
	UpdatedAt         *time.Time
}

func NewAssistant(
	categoryID string,
	name string,
	description string,
	model string,
	systemPrompt string,
	exampleUserPrompt *string,
	isActive *bool,
) (Assistant, error) {
	if !IsValidID(categoryID) {
		return Assistant{}, ErrInvalidCategoryID
	}

	normalizedName := strings.TrimSpace(name)
	if normalizedName == "" {
		return Assistant{}, ErrInvalidName
	}

	normalizedDescription := strings.TrimSpace(description)
	if normalizedDescription == "" {
		return Assistant{}, ErrInvalidDescription
	}

	normalizedModel := strings.TrimSpace(model)
	if normalizedModel == "" {
		return Assistant{}, ErrInvalidModel
	}

	normalizedSystemPrompt := strings.TrimSpace(systemPrompt)
	if normalizedSystemPrompt == "" {
		return Assistant{}, ErrInvalidSystemPrompt
	}

	active := true
	if isActive != nil {
		active = *isActive
	}

	return Assistant{
		CategoryID:        categoryID,
		Name:              normalizedName,
		Description:       normalizedDescription,
		Model:             normalizedModel,
		SystemPrompt:      normalizedSystemPrompt,
		ExampleUserPrompt: normalizeNullableString(exampleUserPrompt),
		IsActive:          active,
	}, nil
}

func ReconstituteAssistant(
	id string,
	categoryID string,
	categoryName *string,
	name string,
	description string,
	model string,
	systemPrompt string,
	exampleUserPrompt *string,
	isActive bool,
	createdAt *time.Time,
	updatedAt *time.Time,
) Assistant {
	return Assistant{
		ID:                id,
		CategoryID:        categoryID,
		CategoryName:      categoryName,
		Name:              name,
		Description:       description,
		Model:             model,
		SystemPrompt:      systemPrompt,
		ExampleUserPrompt: exampleUserPrompt,
		IsActive:          isActive,
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
	}
}

func IsValidID(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}

func normalizeNullableString(value *string) *string {
	if value == nil {
		return nil
	}

	normalized := strings.TrimSpace(*value)
	if normalized == "" {
		return nil
	}

	return &normalized
}
