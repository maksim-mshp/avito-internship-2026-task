package v1

import "ai-assistants-catalog/internal/categories/domain"

func mapCategory(category domain.Category) CategoryDTO {
	return CategoryDTO{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		CreatedAt:   category.CreatedAt,
	}
}

func mapCategories(categories []domain.Category) []CategoryDTO {
	result := make([]CategoryDTO, 0, len(categories))
	for _, category := range categories {
		result = append(result, mapCategory(category))
	}

	return result
}
