package handlers

import (
	"context"
	"errors"
	"testing"

	"ai-assistants-catalog/internal/assistants/app"
	"ai-assistants-catalog/internal/assistants/domain"
)

const (
	testAssistantID = "8a471523-4b7a-4f07-a0c8-1f901962a3da"
	testCategoryID  = "4b042f5e-3887-46bf-9c54-6fae1d664c49"
)

type fakeRepository struct {
	assistant     domain.Assistant
	listQuery     app.ListQuery
	createCalled  bool
	updateCalled  bool
	getByIDCalled bool
	listCalled    bool
	addCalled     bool
	removeCalled  bool
}

func (r *fakeRepository) List(_ context.Context, query app.ListQuery) (app.ListResult, error) {
	r.listCalled = true
	r.listQuery = query
	return app.ListResult{
		Assistants: []domain.Assistant{r.assistant},
		Page:       query.Page,
		PageSize:   query.PageSize,
		Total:      1,
	}, nil
}

func (r *fakeRepository) GetByID(_ context.Context, _ string, _ bool, _ string) (domain.Assistant, error) {
	r.getByIDCalled = true
	return r.assistant, nil
}

func (r *fakeRepository) Create(_ context.Context, assistant domain.Assistant) (domain.Assistant, error) {
	r.createCalled = true
	r.assistant = assistant
	r.assistant.ID = testAssistantID
	return r.assistant, nil
}

func (r *fakeRepository) Update(_ context.Context, assistant domain.Assistant) (domain.Assistant, error) {
	r.updateCalled = true
	r.assistant = assistant
	return r.assistant, nil
}

func (r *fakeRepository) AddFavorite(_ context.Context, _ string, _ string, _ bool) error {
	r.addCalled = true
	return nil
}

func (r *fakeRepository) RemoveFavorite(_ context.Context, _ string, _ string) error {
	r.removeCalled = true
	return nil
}

func TestCreateHandlerHandle(t *testing.T) {
	repo := &fakeRepository{}
	handler := NewCreateHandler(repo)
	command := validCreateCommand()

	assistant, err := handler.Handle(context.Background(), command)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if !repo.createCalled {
		t.Fatalf("expected repository call")
	}
	if assistant.ID != testAssistantID {
		t.Fatalf("unexpected assistant id: %s", assistant.ID)
	}
}

func TestCreateHandlerInvalidCommand(t *testing.T) {
	repo := &fakeRepository{}
	handler := NewCreateHandler(repo)

	_, err := handler.Handle(context.Background(), app.CreateCommand{})
	if !errors.Is(err, domain.ErrInvalidCategoryID) {
		t.Fatalf("expected ErrInvalidCategoryID, got %v", err)
	}
	if repo.createCalled {
		t.Fatalf("expected no repository call")
	}
}

func TestUpdateHandlerHandle(t *testing.T) {
	repo := &fakeRepository{}
	handler := NewUpdateHandler(repo)
	command := validUpdateCommand()

	assistant, err := handler.Handle(context.Background(), command)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if !repo.updateCalled {
		t.Fatalf("expected repository call")
	}
	if assistant.ID != testAssistantID {
		t.Fatalf("unexpected assistant id: %s", assistant.ID)
	}
}

func TestUpdateHandlerRequiresActiveState(t *testing.T) {
	repo := &fakeRepository{}
	handler := NewUpdateHandler(repo)
	command := validUpdateCommand()
	command.IsActive = nil

	_, err := handler.Handle(context.Background(), command)
	if !errors.Is(err, domain.ErrInvalidActiveState) {
		t.Fatalf("expected ErrInvalidActiveState, got %v", err)
	}
	if repo.updateCalled {
		t.Fatalf("expected no repository call")
	}
}

func TestListHandlerHandle(t *testing.T) {
	repo := &fakeRepository{assistant: validAssistant()}
	handler := NewListHandler(repo)
	search := "  cook  "

	result, err := handler.Handle(context.Background(), app.ListQuery{
		CategoryID: &[]string{testCategoryID}[0],
		Search:     &search,
		Page:       1,
		PageSize:   10,
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if !repo.listCalled {
		t.Fatalf("expected repository call")
	}
	if repo.listQuery.Search == nil || *repo.listQuery.Search != "cook" {
		t.Fatalf("unexpected search value: %v", repo.listQuery.Search)
	}
	if result.Total != 1 {
		t.Fatalf("unexpected total: %d", result.Total)
	}
}

func TestListHandlerInvalidPagination(t *testing.T) {
	repo := &fakeRepository{}
	handler := NewListHandler(repo)

	_, err := handler.Handle(context.Background(), app.ListQuery{Page: 0, PageSize: 10})
	if !errors.Is(err, domain.ErrInvalidPagination) {
		t.Fatalf("expected ErrInvalidPagination, got %v", err)
	}
	if repo.listCalled {
		t.Fatalf("expected no repository call")
	}
}

func TestGetByIDHandlerInvalidID(t *testing.T) {
	repo := &fakeRepository{}
	handler := NewGetByIDHandler(repo)

	_, err := handler.Handle(context.Background(), app.GetByIDQuery{ID: "bad"})
	if !errors.Is(err, domain.ErrInvalidID) {
		t.Fatalf("expected ErrInvalidID, got %v", err)
	}
	if repo.getByIDCalled {
		t.Fatalf("expected no repository call")
	}
}

func TestAddFavoriteHandlerHandle(t *testing.T) {
	repo := &fakeRepository{}
	handler := NewAddFavoriteHandler(repo)

	err := handler.Handle(context.Background(), app.FavoriteCommand{
		UserID:      "44c75af3-eca3-4867-85fc-b8245eaafa3a",
		AssistantID: testAssistantID,
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !repo.addCalled {
		t.Fatalf("expected repository call")
	}
}

func TestRemoveFavoriteHandlerInvalidUserID(t *testing.T) {
	repo := &fakeRepository{}
	handler := NewRemoveFavoriteHandler(repo)

	err := handler.Handle(context.Background(), app.FavoriteCommand{
		UserID:      "bad",
		AssistantID: testAssistantID,
	})
	if !errors.Is(err, domain.ErrInvalidID) {
		t.Fatalf("expected ErrInvalidID, got %v", err)
	}
	if repo.removeCalled {
		t.Fatalf("expected no repository call")
	}
}

func validCreateCommand() app.CreateCommand {
	categoryID := testCategoryID
	name := "Assistant"
	description := "Description"
	model := "mock-smart"
	systemPrompt := "System"
	isActive := true

	return app.CreateCommand{
		CategoryID:   &categoryID,
		Name:         &name,
		Description:  &description,
		Model:        &model,
		SystemPrompt: &systemPrompt,
		Tags:         []string{"support", "text"},
		IsActive:     &isActive,
	}
}

func validUpdateCommand() app.UpdateCommand {
	command := validCreateCommand()

	return app.UpdateCommand{
		ID:           testAssistantID,
		CategoryID:   command.CategoryID,
		Name:         command.Name,
		Description:  command.Description,
		Model:        command.Model,
		SystemPrompt: command.SystemPrompt,
		Tags:         command.Tags,
		IsActive:     command.IsActive,
	}
}

func validAssistant() domain.Assistant {
	return domain.ReconstituteAssistant(
		testAssistantID,
		testCategoryID,
		nil,
		"Assistant",
		"Description",
		"mock-smart",
		"System",
		nil,
		[]string{"support"},
		false,
		true,
		nil,
		nil,
	)
}
