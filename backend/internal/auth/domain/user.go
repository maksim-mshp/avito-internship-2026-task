package domain

import "time"

type User struct {
	ID        string
	Email     string
	Role      Role
	CreatedAt *time.Time
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
