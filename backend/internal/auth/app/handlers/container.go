package handlers

import "ai-assistants-catalog/internal/auth/app"

type Handlers struct {
	DummyLogin *DummyLoginHandler
}

func BuildHandlers(jwtToken string, repo app.Repository) *Handlers {
	return &Handlers{
		DummyLogin: NewDummyLoginHandler(jwtToken, repo),
	}
}
