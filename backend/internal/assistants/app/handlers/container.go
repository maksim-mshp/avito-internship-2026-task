package handlers

import "ai-assistants-catalog/internal/assistants/app"

type Handlers struct {
	List    *ListHandler
	GetByID *GetByIDHandler
	Create  *CreateHandler
	Update  *UpdateHandler
}

func BuildHandlers(repo app.Repository) *Handlers {
	return &Handlers{
		List:    NewListHandler(repo),
		GetByID: NewGetByIDHandler(repo),
		Create:  NewCreateHandler(repo),
		Update:  NewUpdateHandler(repo),
	}
}
