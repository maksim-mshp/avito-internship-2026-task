package handlers

import (
	"context"

	"ai-assistants-catalog/internal/auth/app"
	"ai-assistants-catalog/internal/auth/domain"
	"ai-assistants-catalog/internal/core/security"
)

type DummyLoginResult struct {
	Token string
	User  domain.User
}

type DummyLoginHandler struct {
	jwtToken string
}

func NewDummyLoginHandler(jwtToken string) *DummyLoginHandler {
	return &DummyLoginHandler{jwtToken: jwtToken}
}

func (h *DummyLoginHandler) Handle(ctx context.Context, cmd app.DummyLoginCommand) (DummyLoginResult, error) {
	role, err := domain.NewRole(cmd.Role)
	if err != nil {
		return DummyLoginResult{}, err
	}

	user := domain.NewDummyUser(role)
	token, err := security.GenerateJWT(h.jwtToken, user.ID, user.Role.String())
	if err != nil {
		return DummyLoginResult{}, err
	}

	return DummyLoginResult{
		Token: token,
		User:  user,
	}, nil
}
