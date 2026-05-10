package handlers

import (
	"context"

	"ai-assistants-catalog/internal/auth/app"
	"ai-assistants-catalog/internal/auth/domain"

	"golang.org/x/crypto/bcrypt"
)

type RegisterHandler struct {
	repo app.Repository
}

func NewRegisterHandler(repo app.Repository) *RegisterHandler {
	return &RegisterHandler{repo: repo}
}

func (h *RegisterHandler) Handle(ctx context.Context, cmd app.RegisterCommand) (domain.User, error) {
	email, err := domain.NormalizeEmail(cmd.Email)
	if err != nil {
		return domain.User{}, err
	}

	if err = domain.ValidatePassword(cmd.Password); err != nil {
		return domain.User{}, err
	}

	roleValue := cmd.Role
	if roleValue == "" {
		roleValue = domain.RoleUser.String()
	}

	role, err := domain.NewRole(roleValue)
	if err != nil {
		return domain.User{}, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(cmd.Password), bcrypt.DefaultCost)
	if err != nil {
		return domain.User{}, err
	}

	return h.repo.CreateUser(ctx, domain.NewUser(email, role), string(hash))
}
