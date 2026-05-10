package v1

import "time"

type DummyLoginRequest struct {
	Role string `json:"role"`
}

type UserDTO struct {
	ID        string     `json:"id"`
	Email     string     `json:"email"`
	Role      string     `json:"role"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
}

type TokenResponse struct {
	Token string  `json:"token"`
	User  UserDTO `json:"user"`
}
