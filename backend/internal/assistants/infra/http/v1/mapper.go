package v1

import (
	"ai-assistants-catalog/internal/assistants/app"
	"ai-assistants-catalog/internal/assistants/domain"
)

func mapAssistant(assistant domain.Assistant) AssistantDTO {
	systemPrompt := assistant.SystemPrompt

	return AssistantDTO{
		ID:                assistant.ID,
		CategoryID:        assistant.CategoryID,
		CategoryName:      assistant.CategoryName,
		Name:              assistant.Name,
		Description:       assistant.Description,
		Model:             assistant.Model,
		SystemPrompt:      &systemPrompt,
		ExampleUserPrompt: assistant.ExampleUserPrompt,
		IsActive:          assistant.IsActive,
		CreatedAt:         assistant.CreatedAt,
		UpdatedAt:         assistant.UpdatedAt,
	}
}

func mapAssistants(assistants []domain.Assistant) []AssistantDTO {
	result := make([]AssistantDTO, 0, len(assistants))
	for _, assistant := range assistants {
		result = append(result, mapAssistant(assistant))
	}

	return result
}

func mapListResult(result app.ListResult) AssistantsResponse {
	return AssistantsResponse{
		Assistants: mapAssistants(result.Assistants),
		Pagination: PaginationDTO{
			Page:     result.Page,
			PageSize: result.PageSize,
			Total:    result.Total,
		},
	}
}
