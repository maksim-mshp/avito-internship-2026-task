package app

import "context"

type LLMRequest struct {
	Model        string
	SystemPrompt string
	UserPrompt   string
}

type LLMProvider interface {
	Generate(ctx context.Context, request LLMRequest) (string, error)
}
