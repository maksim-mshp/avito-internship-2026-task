package handlers

import (
	"context"
	"errors"
	"testing"

	"ai-assistants-catalog/internal/runs/app"
	"ai-assistants-catalog/internal/runs/domain"
)

func TestListMyHandlerValidation(t *testing.T) {
	repo := &fakeRunRepository{}
	handler := NewListMyHandler(repo)
	status := "bad"

	_, err := handler.Handle(context.Background(), app.ListMyQuery{
		UserID:   testUserID,
		Status:   &status,
		Page:     1,
		PageSize: 10,
	})
	if !errors.Is(err, domain.ErrInvalidStatus) {
		t.Fatalf("expected ErrInvalidStatus, got %v", err)
	}
}

func TestListAdminHandlerValidation(t *testing.T) {
	repo := &fakeRunRepository{}
	handler := NewListAdminHandler(repo)
	assistantID := "bad"

	_, err := handler.Handle(context.Background(), app.ListAdminQuery{
		AssistantID: &assistantID,
		Page:        1,
		PageSize:    10,
	})
	if !errors.Is(err, domain.ErrInvalidAssistantID) {
		t.Fatalf("expected ErrInvalidAssistantID, got %v", err)
	}
}

func TestListMyHandlerSuccess(t *testing.T) {
	repo := &fakeRunRepository{}
	handler := NewListMyHandler(repo)

	result, err := handler.Handle(context.Background(), app.ListMyQuery{
		UserID:   testUserID,
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if result.Total != 0 {
		t.Fatalf("unexpected total: %d", result.Total)
	}
}
