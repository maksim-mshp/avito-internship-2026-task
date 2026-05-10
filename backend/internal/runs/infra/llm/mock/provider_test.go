package mock

import (
	"context"
	"errors"
	"strings"
	"testing"

	"ai-assistants-catalog/internal/runs/app"
)

func TestProviderGenerate(t *testing.T) {
	provider := NewProvider()

	output, err := provider.Generate(context.Background(), app.LLMRequest{
		Model:        "mock-smart",
		SystemPrompt: "system",
		UserPrompt:   " hello ",
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !strings.Contains(output, "mock-smart") || !strings.Contains(output, "hello") || !strings.Contains(output, "system") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestProviderGenerateError(t *testing.T) {
	provider := NewProvider()

	_, err := provider.Generate(context.Background(), app.LLMRequest{
		Model:      "mock-smart",
		UserPrompt: "mock_error",
	})
	if !errors.Is(err, ErrMockProvider) {
		t.Fatalf("expected ErrMockProvider, got %v", err)
	}
}
