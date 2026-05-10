package v1

import "ai-assistants-catalog/internal/auth/domain"

func mapUser(user domain.User) UserDTO {
	return UserDTO{
		ID:        user.ID,
		Email:     user.Email,
		Role:      user.Role.String(),
		CreatedAt: user.CreatedAt,
	}
}
