package app

import "ai-assistants-catalog/internal/assistants/domain"

type GetByIDQuery struct {
	ID              string
	UserID          string
	IncludeInactive bool
}

type ListQuery struct {
	CategoryID      *string
	Search          *string
	Tag             *string
	UserID          string
	FavoritesOnly   bool
	IncludeInactive bool
	Page            int
	PageSize        int
}

type ListResult struct {
	Assistants []domain.Assistant
	Page       int
	PageSize   int
	Total      int
}
