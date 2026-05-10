package v1

import "time"

type CategoryDTO struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description *string    `json:"description,omitempty"`
	CreatedAt   *time.Time `json:"createdAt,omitempty"`
}

type CategoryCreateRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type CategoriesResponse struct {
	Categories []CategoryDTO `json:"categories"`
}
