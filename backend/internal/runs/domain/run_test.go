package domain

import (
	"errors"
	"testing"
)

const (
	testAssistantID = "8a471523-4b7a-4f07-a0c8-1f901962a3da"
	testUserID      = "44c75af3-eca3-4867-85fc-b8245eaafa3a"
)

func TestNewPendingRun(t *testing.T) {
	run, err := NewPendingRun(testAssistantID, testUserID, " mock-smart ", " hello ")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if run.Model != "mock-smart" {
		t.Fatalf("unexpected model: %s", run.Model)
	}
	if run.UserPrompt != "hello" {
		t.Fatalf("unexpected user prompt: %s", run.UserPrompt)
	}
	if run.Status != StatusPending {
		t.Fatalf("unexpected status: %s", run.Status)
	}
}

func TestNewPendingRunValidation(t *testing.T) {
	tests := []struct {
		name        string
		assistantID string
		userID      string
		model       string
		userPrompt  string
		expected    error
	}{
		{name: "invalid assistant", assistantID: "bad", userID: testUserID, model: "model", userPrompt: "prompt", expected: ErrInvalidAssistantID},
		{name: "invalid user", assistantID: testAssistantID, userID: "bad", model: "model", userPrompt: "prompt", expected: ErrInvalidUserID},
		{name: "empty model", assistantID: testAssistantID, userID: testUserID, model: " ", userPrompt: "prompt", expected: ErrProviderFailed},
		{name: "empty prompt", assistantID: testAssistantID, userID: testUserID, model: "model", userPrompt: " ", expected: ErrInvalidUserPrompt},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewPendingRun(tt.assistantID, tt.userID, tt.model, tt.userPrompt)
			if !errors.Is(err, tt.expected) {
				t.Fatalf("expected %v, got %v", tt.expected, err)
			}
		})
	}
}

func TestIsValidStatus(t *testing.T) {
	if !IsValidStatus(StatusPending) || !IsValidStatus(StatusSuccess) || !IsValidStatus(StatusFailed) {
		t.Fatalf("expected known statuses to be valid")
	}
	if IsValidStatus("unknown") {
		t.Fatalf("expected unknown status to be invalid")
	}
}
