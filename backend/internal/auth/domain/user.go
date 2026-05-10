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
