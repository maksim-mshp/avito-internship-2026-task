package v1

import "time"

type AssistantDTO struct {
	ID                string     `json:"id"`
	CategoryID        string     `json:"categoryId"`
	CategoryName      *string    `json:"categoryName,omitempty"`
	Name              string     `json:"name"`
	Description       string     `json:"description"`
	Model             string     `json:"model"`
	SystemPrompt      *string    `json:"systemPrompt,omitempty"`
	ExampleUserPrompt *string    `json:"exampleUserPrompt,omitempty"`
	Tags              []string   `json:"tags"`
	IsFavorite        bool       `json:"isFavorite"`
	IsActive          bool       `json:"isActive"`
	CreatedAt         *time.Time `json:"createdAt,omitempty"`
	UpdatedAt         *time.Time `json:"updatedAt,omitempty"`
}

type AssistantCreateRequest struct {
	CategoryID        *string  `json:"categoryId"`
	Name              *string  `json:"name"`
	Description       *string  `json:"description"`
	Model             *string  `json:"model"`
	SystemPrompt      *string  `json:"systemPrompt"`
	ExampleUserPrompt *string  `json:"exampleUserPrompt"`
	Tags              []string `json:"tags"`
	IsActive          *bool    `json:"isActive"`
}

type AssistantUpdateRequest struct {
	CategoryID        *string  `json:"categoryId"`
	Name              *string  `json:"name"`
	Description       *string  `json:"description"`
	Model             *string  `json:"model"`
	SystemPrompt      *string  `json:"systemPrompt"`
	ExampleUserPrompt *string  `json:"exampleUserPrompt"`
	Tags              []string `json:"tags"`
	IsActive          *bool    `json:"isActive"`
}

type AssistantsResponse struct {
	Assistants []AssistantDTO `json:"assistants"`
	Pagination PaginationDTO  `json:"pagination"`
}

type PaginationDTO struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
	Total    int `json:"total"`
}
