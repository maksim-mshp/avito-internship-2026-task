package domain

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

const (
	adminUserID = "ca2a62f3-c998-4050-96c1-0c0f62cf6568"
	userUserID  = "44c75af3-eca3-4867-85fc-b8245eaafa3a"
)

func NewRole(value string) (Role, error) {
	role := Role(value)
	switch role {
	case RoleAdmin, RoleUser:
		return role, nil
	default:
		return "", ErrInvalidRole
	}
}

func (r Role) String() string {
	return string(r)
}

func (r Role) UserID() string {
	if r == RoleAdmin {
		return adminUserID
	}

	return userUserID
}

func (r Role) Email() string {
	if r == RoleAdmin {
		return "admin@example.com"
	}

	return "user@example.com"
}
