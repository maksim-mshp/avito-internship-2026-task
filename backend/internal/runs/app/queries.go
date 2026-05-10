package app

import "ai-assistants-catalog/internal/runs/domain"

type ListMyQuery struct {
	UserID   string
	Status   *string
	Page     int
	PageSize int
}

type ListAdminQuery struct {
	AssistantID *string
	Status      *string
	Page        int
	PageSize    int
}

type ListResult struct {
	Runs     []domain.Run
	Page     int
	PageSize int
	Total    int
}
