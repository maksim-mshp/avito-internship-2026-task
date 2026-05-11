package handlers

import "ai-assistants-catalog/internal/runs/app"

type Handlers struct {
	RunAssistant *RunAssistantHandler
	SetRating    *SetRatingHandler
	ListMy       *ListMyHandler
	ListAdmin    *ListAdminHandler
}

func BuildHandlers(
	runs app.RunRepository,
	assistants app.AssistantRepository,
	provider app.LLMProvider,
) *Handlers {
	return &Handlers{
		RunAssistant: NewRunAssistantHandler(runs, assistants, provider),
		SetRating:    NewSetRatingHandler(runs),
		ListMy:       NewListMyHandler(runs),
		ListAdmin:    NewListAdminHandler(runs),
	}
}
