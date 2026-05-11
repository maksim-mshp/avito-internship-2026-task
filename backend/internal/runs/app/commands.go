package app

type RunAssistantCommand struct {
	AssistantID string
	UserID      string
	UserPrompt  *string
}

type SetRatingCommand struct {
	ID     string
	UserID string
	Rating *string
}
