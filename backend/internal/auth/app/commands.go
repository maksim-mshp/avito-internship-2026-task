package app

type DummyLoginCommand struct {
	Role string
}

type RegisterCommand struct {
	Email    string
	Password string
	Role     string
}

type LoginCommand struct {
	Email    string
	Password string
}
