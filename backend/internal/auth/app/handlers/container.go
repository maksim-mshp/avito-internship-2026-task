package handlers

import "ai-assistants-catalog/internal/auth/app"

type Handlers struct {
	DummyLogin *DummyLoginHandler
	Register   *RegisterHandler
	Login      *LoginHandler
}

func BuildHandlers(jwtToken string, repo app.Repository) *Handlers {
	return &Handlers{
		DummyLogin: NewDummyLoginHandler(jwtToken, repo),
		Register:   NewRegisterHandler(repo),
		Login:      NewLoginHandler(jwtToken, repo),
	}
}
