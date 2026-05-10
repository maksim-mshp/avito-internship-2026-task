package mock

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"ai-assistants-catalog/internal/runs/app"
)

var ErrMockProvider = errors.New("mock provider error")

type Provider struct{}

func NewProvider() *Provider {
	return &Provider{}
}

func (p *Provider) Generate(ctx context.Context, request app.LLMRequest) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	userPrompt := strings.TrimSpace(request.UserPrompt)
	if strings.Contains(strings.ToLower(userPrompt), "mock_error") {
		return "", ErrMockProvider
	}

	return fmt.Sprintf(
		"Mock response from %s: user prompt %q was processed with system prompt %q.",
		request.Model,
		userPrompt,
		strings.TrimSpace(request.SystemPrompt),
	), nil
}
