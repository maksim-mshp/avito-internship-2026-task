package handlers

import (
	"context"
	"errors"
	"testing"

	"ai-assistants-catalog/internal/categories/app"
	"ai-assistants-catalog/internal/categories/domain"
)

type fakeRepository struct {
	created      domain.Category
	createCalled bool
}

func (r *fakeRepository) GetAll(_ context.Context) ([]domain.Category, error) {
	return []domain.Category{r.created}, nil
}

func (r *fakeRepository) Create(_ context.Context, category domain.Category) (domain.Category, error) {
	r.createCalled = true
	r.created = category
	r.created.ID = "category-id"
	return r.created, nil
}

func TestCreateHandlerHandle(t *testing.T) {
	repo := &fakeRepository{}
	handler := NewCreateHandler(repo)
	name := "  Food  "

	category, err := handler.Handle(context.Background(), app.CreateCommand{Name: &name})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if !repo.createCalled {
		t.Fatalf("expected repository call")
	}
	if category.Name != "Food" {
		t.Fatalf("unexpected category name: %s", category.Name)
	}
}

func TestCreateHandlerInvalidName(t *testing.T) {
	repo := &fakeRepository{}
	handler := NewCreateHandler(repo)

	_, err := handler.Handle(context.Background(), app.CreateCommand{})
	if !errors.Is(err, domain.ErrInvalidName) {
		t.Fatalf("expected ErrInvalidName, got %v", err)
	}
	if repo.createCalled {
		t.Fatalf("expected no repository call")
	}
}

func TestGetAllHandlerHandle(t *testing.T) {
	repo := &fakeRepository{created: domain.ReconstituteCategory("id", "Food", nil, nil)}
	handler := NewGetAllHandler(repo)

	categories, err := handler.Handle(context.Background(), app.GetAllQuery{})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if len(categories) != 1 || categories[0].Name != "Food" {
		t.Fatalf("unexpected categories: %v", categories)
	}
}
