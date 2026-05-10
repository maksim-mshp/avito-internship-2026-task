package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Run struct {
	ID            string
	AssistantID   string
	AssistantName *string
	CategoryID    *string
	CategoryName  *string
	UserID        string
	Model         string
	UserPrompt    string
	Output        *string
	Status        string
	Error         *string
	CreatedAt     *time.Time
}

func NewPendingRun(assistantID string, userID string, model string, userPrompt string) (Run, error) {
	if !IsValidID(assistantID) {
		return Run{}, ErrInvalidAssistantID
	}

	if !IsValidID(userID) {
		return Run{}, ErrInvalidUserID
	}

	normalizedModel := strings.TrimSpace(model)
	if normalizedModel == "" {
		return Run{}, ErrProviderFailed
	}

	normalizedPrompt := strings.TrimSpace(userPrompt)
	if normalizedPrompt == "" {
		return Run{}, ErrInvalidUserPrompt
	}

	return Run{
		AssistantID: assistantID,
		UserID:      userID,
		Model:       normalizedModel,
		UserPrompt:  normalizedPrompt,
		Status:      StatusPending,
	}, nil
}

func ReconstituteRun(
	id string,
	assistantID string,
	assistantName *string,
	categoryID *string,
	categoryName *string,
	userID string,
	model string,
	userPrompt string,
	output *string,
	status string,
	runError *string,
	createdAt *time.Time,
) Run {
	return Run{
		ID:            id,
		AssistantID:   assistantID,
		AssistantName: assistantName,
		CategoryID:    categoryID,
		CategoryName:  categoryName,
		UserID:        userID,
		Model:         model,
		UserPrompt:    userPrompt,
		Output:        output,
		Status:        status,
		Error:         runError,
		CreatedAt:     createdAt,
	}
}

func IsValidID(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}
