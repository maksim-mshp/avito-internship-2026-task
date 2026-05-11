package domain

import (
	"errors"
	"testing"
)

const testCategoryID = "4b042f5e-3887-46bf-9c54-6fae1d664c49"

func TestNewAssistant(t *testing.T) {
	examplePrompt := "  ingredients  "

	assistant, err := NewAssistant(
		testCategoryID,
		"  Cook  ",
		"  Helps with recipes  ",
		"  mock-smart  ",
		"  System prompt  ",
		&examplePrompt,
		[]string{"  recipes  ", "food", "recipes"},
		nil,
	)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if assistant.Name != "Cook" {
		t.Fatalf("unexpected name: %s", assistant.Name)
	}
	if assistant.Description != "Helps with recipes" {
		t.Fatalf("unexpected description: %s", assistant.Description)
	}
	if assistant.Model != "mock-smart" {
		t.Fatalf("unexpected model: %s", assistant.Model)
	}
	if assistant.SystemPrompt != "System prompt" {
		t.Fatalf("unexpected system prompt: %s", assistant.SystemPrompt)
	}
	if assistant.ExampleUserPrompt == nil || *assistant.ExampleUserPrompt != "ingredients" {
		t.Fatalf("unexpected example prompt: %v", assistant.ExampleUserPrompt)
	}
	if len(assistant.Tags) != 2 || assistant.Tags[0] != "recipes" || assistant.Tags[1] != "food" {
		t.Fatalf("unexpected tags: %v", assistant.Tags)
	}
	if assistant.IsFavorite {
		t.Fatalf("expected new assistant to be not favorite")
	}
	if !assistant.IsActive {
		t.Fatalf("expected active assistant by default")
	}
}

func TestNewAssistantValidation(t *testing.T) {
	tests := []struct {
		name          string
		categoryID    string
		assistantName string
		description   string
		model         string
		systemPrompt  string
		expected      error
	}{
		{name: "invalid category", categoryID: "bad", assistantName: "Name", description: "Description", model: "model", systemPrompt: "prompt", expected: ErrInvalidCategoryID},
		{name: "empty name", categoryID: testCategoryID, assistantName: " ", description: "Description", model: "model", systemPrompt: "prompt", expected: ErrInvalidName},
		{name: "empty description", categoryID: testCategoryID, assistantName: "Name", description: " ", model: "model", systemPrompt: "prompt", expected: ErrInvalidDescription},
		{name: "empty model", categoryID: testCategoryID, assistantName: "Name", description: "Description", model: " ", systemPrompt: "prompt", expected: ErrInvalidModel},
		{name: "empty system prompt", categoryID: testCategoryID, assistantName: "Name", description: "Description", model: "model", systemPrompt: " ", expected: ErrInvalidSystemPrompt},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewAssistant(tt.categoryID, tt.assistantName, tt.description, tt.model, tt.systemPrompt, nil, nil, nil)
			if !errors.Is(err, tt.expected) {
				t.Fatalf("expected %v, got %v", tt.expected, err)
			}
		})
	}
}
