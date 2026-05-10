package v1

import "time"

type RunDTO struct {
	ID            string     `json:"id"`
	AssistantID   string     `json:"assistantId"`
	AssistantName *string    `json:"assistantName,omitempty"`
	CategoryID    *string    `json:"categoryId,omitempty"`
	CategoryName  *string    `json:"categoryName,omitempty"`
	UserID        string     `json:"userId"`
	Model         string     `json:"model"`
	UserPrompt    string     `json:"userPrompt"`
	Output        *string    `json:"output,omitempty"`
	Status        string     `json:"status"`
	Error         *string    `json:"error,omitempty"`
	CreatedAt     *time.Time `json:"createdAt,omitempty"`
}

type RunCreateRequest struct {
	UserPrompt *string `json:"userPrompt"`
}

type RunsResponse struct {
	Runs       []RunDTO      `json:"runs"`
	Pagination PaginationDTO `json:"pagination"`
}

type PaginationDTO struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
	Total    int `json:"total"`
}
