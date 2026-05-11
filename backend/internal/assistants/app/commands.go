package app

type CreateCommand struct {
	CategoryID        *string
	Name              *string
	Description       *string
	Model             *string
	SystemPrompt      *string
	ExampleUserPrompt *string
	Tags              []string
	IsActive          *bool
}

type UpdateCommand struct {
	ID                string
	CategoryID        *string
	Name              *string
	Description       *string
	Model             *string
	SystemPrompt      *string
	ExampleUserPrompt *string
	Tags              []string
	IsActive          *bool
}

type FavoriteCommand struct {
	UserID          string
	AssistantID     string
	IncludeInactive bool
}
