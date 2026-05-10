package domain

import (
	"net/mail"
	"strings"
	"time"
)

type User struct {
	ID        string
	Email     string
	Role      Role
	CreatedAt *time.Time
}

type AuthUser struct {
	User
	PasswordHash string
}

func NewUser(email string, role Role) User {
	return User{
		Email: email,
		Role:  role,
	}
}

func NewDummyUser(role Role) User {
	return User{
		ID:    role.UserID(),
		Email: role.Email(),
		Role:  role,
	}
}

func ReconstituteUser(id string, email string, role Role, createdAt *time.Time) User {
	return User{
		ID:        id,
		Email:     email,
		Role:      role,
		CreatedAt: createdAt,
	}
}

func NormalizeEmail(email string) (string, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" {
		return "", ErrInvalidEmail
	}

	parsed, err := mail.ParseAddress(email)
	if err != nil || parsed.Address != email {
		return "", ErrInvalidEmail
	}

	return email, nil
}

func ValidatePassword(password string) error {
	if strings.TrimSpace(password) == "" {
		return ErrInvalidPassword
	}

	return nil
}
