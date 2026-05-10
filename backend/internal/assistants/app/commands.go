package app

type CreateCommand struct {
	CategoryID        *string
	Name              *string
	Description       *string
	Model             *string
	SystemPrompt      *string
	ExampleUserPrompt *string
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
	IsActive          *bool
}
