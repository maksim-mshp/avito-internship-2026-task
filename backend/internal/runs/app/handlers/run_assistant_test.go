package handlers

import (
	"context"
	"errors"
	"testing"

	assistantdomain "ai-assistants-catalog/internal/assistants/domain"
	"ai-assistants-catalog/internal/runs/app"
	"ai-assistants-catalog/internal/runs/domain"
)

const (
	testAssistantID = "8a471523-4b7a-4f07-a0c8-1f901962a3da"
	testCategoryID  = "4b042f5e-3887-46bf-9c54-6fae1d664c49"
	testUserID      = "44c75af3-eca3-4867-85fc-b8245eaafa3a"
	testRunID       = "10d7eb1f-6506-4e57-95e3-80c53040fd75"
)

type fakeAssistantRepository struct {
	assistant assistantdomain.Assistant
	err       error
}

func (r fakeAssistantRepository) GetByID(_ context.Context, _ string, _ bool, _ string) (assistantdomain.Assistant, error) {
	if r.err != nil {
		return assistantdomain.Assistant{}, r.err
	}

	return r.assistant, nil
}

type fakeRunRepository struct {
	created      domain.Run
	completed    domain.Run
	failed       domain.Run
	createCall   bool
	completeCall bool
	failCall     bool
}

func (r *fakeRunRepository) CreatePending(_ context.Context, run domain.Run) (domain.Run, error) {
	r.createCall = true
	r.created = domain.ReconstituteRun(
		testRunID,
		run.AssistantID,
		nil,
		nil,
		nil,
		run.UserID,
		run.Model,
		run.UserPrompt,
		nil,
		domain.StatusPending,
		nil,
		nil,
	)

	return r.created, nil
}

func (r *fakeRunRepository) Complete(_ context.Context, id string, output string) (domain.Run, error) {
	r.completeCall = true
	r.completed = domain.ReconstituteRun(
		id,
		r.created.AssistantID,
		nil,
		nil,
		nil,
		r.created.UserID,
		r.created.Model,
		r.created.UserPrompt,
		&output,
		domain.StatusSuccess,
		nil,
		nil,
	)

	return r.completed, nil
}

func (r *fakeRunRepository) Fail(_ context.Context, id string, message string) (domain.Run, error) {
	r.failCall = true
	r.failed = domain.ReconstituteRun(
		id,
		r.created.AssistantID,
		nil,
		nil,
		nil,
		r.created.UserID,
		r.created.Model,
		r.created.UserPrompt,
		nil,
		domain.StatusFailed,
		&message,
		nil,
	)

	return r.failed, nil
}

func (r *fakeRunRepository) ListMy(_ context.Context, _ app.ListMyQuery) (app.ListResult, error) {
	return app.ListResult{}, nil
}

func (r *fakeRunRepository) ListAdmin(_ context.Context, _ app.ListAdminQuery) (app.ListResult, error) {
	return app.ListResult{}, nil
}

type fakeProvider struct {
	output  string
	err     error
	request app.LLMRequest
	called  bool
}

func (p *fakeProvider) Generate(_ context.Context, request app.LLMRequest) (string, error) {
	p.called = true
	p.request = request
	if p.err != nil {
		return "", p.err
	}

	return p.output, nil
}

func TestRunAssistantHandlerSuccess(t *testing.T) {
	runs := &fakeRunRepository{}
	provider := &fakeProvider{output: "answer"}
	handler := NewRunAssistantHandler(runs, fakeAssistantRepository{assistant: activeAssistant()}, provider)
	userPrompt := "hello"

	run, err := handler.Handle(context.Background(), app.RunAssistantCommand{
		AssistantID: testAssistantID,
		UserID:      testUserID,
		UserPrompt:  &userPrompt,
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if run.Status != domain.StatusSuccess {
		t.Fatalf("expected success status, got %s", run.Status)
	}
	if run.Output == nil || *run.Output != "answer" {
		t.Fatalf("unexpected output: %v", run.Output)
	}
	if !runs.createCall || !runs.completeCall || runs.failCall {
		t.Fatalf("unexpected repository calls: create=%v complete=%v fail=%v", runs.createCall, runs.completeCall, runs.failCall)
	}
	if !provider.called {
		t.Fatalf("expected provider call")
	}
	if provider.request.SystemPrompt != "system" {
		t.Fatalf("unexpected system prompt: %s", provider.request.SystemPrompt)
	}
}

func TestRunAssistantHandlerInactiveAssistant(t *testing.T) {
	runs := &fakeRunRepository{}
	provider := &fakeProvider{output: "answer"}
	assistant := activeAssistant()
	assistant.IsActive = false
	handler := NewRunAssistantHandler(runs, fakeAssistantRepository{assistant: assistant}, provider)
	userPrompt := "hello"

	_, err := handler.Handle(context.Background(), app.RunAssistantCommand{
		AssistantID: testAssistantID,
		UserID:      testUserID,
		UserPrompt:  &userPrompt,
	})
	if !errors.Is(err, domain.ErrAssistantInactive) {
		t.Fatalf("expected ErrAssistantInactive, got %v", err)
	}
	if runs.createCall || provider.called {
		t.Fatalf("expected no run creation and no provider call")
	}
}

func TestRunAssistantHandlerProviderFailure(t *testing.T) {
	runs := &fakeRunRepository{}
	provider := &fakeProvider{err: errors.New("provider error")}
	handler := NewRunAssistantHandler(runs, fakeAssistantRepository{assistant: activeAssistant()}, provider)
	userPrompt := "hello"

	_, err := handler.Handle(context.Background(), app.RunAssistantCommand{
		AssistantID: testAssistantID,
		UserID:      testUserID,
		UserPrompt:  &userPrompt,
	})
	if !errors.Is(err, domain.ErrProviderFailed) {
		t.Fatalf("expected ErrProviderFailed, got %v", err)
	}
	if !runs.createCall || !runs.failCall || runs.completeCall {
		t.Fatalf("unexpected repository calls: create=%v complete=%v fail=%v", runs.createCall, runs.completeCall, runs.failCall)
	}
	if runs.failed.Error == nil || *runs.failed.Error != "provider error" {
		t.Fatalf("unexpected failed run error: %v", runs.failed.Error)
	}
}

func activeAssistant() assistantdomain.Assistant {
	return assistantdomain.ReconstituteAssistant(
		testAssistantID,
		testCategoryID,
		nil,
		"assistant",
		"description",
		"mock-smart",
		"system",
		nil,
		[]string{"support"},
		false,
		true,
		nil,
		nil,
	)
}
