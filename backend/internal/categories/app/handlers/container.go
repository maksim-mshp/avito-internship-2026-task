package handlers

import "ai-assistants-catalog/internal/categories/app"

type Handlers struct {
	GetAll *GetAllHandler
	Create *CreateHandler
}

func BuildHandlers(repo app.Repository) *Handlers {
	return &Handlers{
		GetAll: NewGetAllHandler(repo),
		Create: NewCreateHandler(repo),
	}
}
