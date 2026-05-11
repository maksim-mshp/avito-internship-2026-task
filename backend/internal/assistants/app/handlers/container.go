package handlers

import "ai-assistants-catalog/internal/assistants/app"

type Handlers struct {
	List           *ListHandler
	GetByID        *GetByIDHandler
	Create         *CreateHandler
	Update         *UpdateHandler
	AddFavorite    *AddFavoriteHandler
	RemoveFavorite *RemoveFavoriteHandler
}

func BuildHandlers(repo app.Repository) *Handlers {
	return &Handlers{
		List:           NewListHandler(repo),
		GetByID:        NewGetByIDHandler(repo),
		Create:         NewCreateHandler(repo),
		Update:         NewUpdateHandler(repo),
		AddFavorite:    NewAddFavoriteHandler(repo),
		RemoveFavorite: NewRemoveFavoriteHandler(repo),
	}
}
