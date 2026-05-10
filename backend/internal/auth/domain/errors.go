package domain

import "errors"

var ErrInvalidRole = errors.New("invalid role")

var ErrUserNotFound = errors.New("user not found")

var ErrInvalidEmail = errors.New("invalid email")

var ErrInvalidPassword = errors.New("invalid password")

var ErrEmailTaken = errors.New("email taken")

var ErrInvalidCredentials = errors.New("invalid credentials")
