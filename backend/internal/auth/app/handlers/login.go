package handlers

import (
	"context"

	"ai-assistants-catalog/internal/auth/app"
	"ai-assistants-catalog/internal/auth/domain"
	"ai-assistants-catalog/internal/core/security"

	"golang.org/x/crypto/bcrypt"
)

type LoginResult struct {
	Token string
	User  domain.User
}

type LoginHandler struct {
	jwtToken string
	repo     app.Repository
}

func NewLoginHandler(jwtToken string, repo app.Repository) *LoginHandler {
	return &LoginHandler{
		jwtToken: jwtToken,
		repo:     repo,
	}
}

func (h *LoginHandler) Handle(ctx context.Context, cmd app.LoginCommand) (LoginResult, error) {
	email, err := domain.NormalizeEmail(cmd.Email)
	if err != nil {
		return LoginResult{}, err
	}

	if err = domain.ValidatePassword(cmd.Password); err != nil {
		return LoginResult{}, err
	}

	authUser, err := h.repo.GetAuthUserByEmail(ctx, email)
	if err != nil {
		return LoginResult{}, err
	}

	if authUser.PasswordHash == "" {
		return LoginResult{}, domain.ErrInvalidCredentials
	}

	if err = bcrypt.CompareHashAndPassword([]byte(authUser.PasswordHash), []byte(cmd.Password)); err != nil {
		return LoginResult{}, domain.ErrInvalidCredentials
	}

	token, err := security.GenerateJWT(h.jwtToken, authUser.ID, authUser.Role.String())
	if err != nil {
		return LoginResult{}, err
	}

	return LoginResult{
		Token: token,
		User:  authUser.User,
	}, nil
}
