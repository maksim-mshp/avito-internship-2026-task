package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	categoryHandlers "ai-assistants-catalog/internal/categories/app/handlers"
	"ai-assistants-catalog/internal/categories/domain"
)

type fakeRepository struct {
	categories []domain.Category
}

func (r *fakeRepository) GetAll(_ context.Context) ([]domain.Category, error) {
	return r.categories, nil
}

func (r *fakeRepository) Create(_ context.Context, category domain.Category) (domain.Category, error) {
	category.ID = "category-id"
	r.categories = append(r.categories, category)
	return category, nil
}

func TestGetAll(t *testing.T) {
	repo := &fakeRepository{
		categories: []domain.Category{domain.ReconstituteCategory("category-id", "Food", nil, nil)},
	}
	handler := NewHTTPHandler(categoryHandlers.BuildHandlers(repo))
	req := httptest.NewRequest(http.MethodGet, "/categories", nil)
	recorder := httptest.NewRecorder()

	handler.GetAll(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}

	var response CategoriesResponse
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(response.Categories) != 1 || response.Categories[0].Name != "Food" {
		t.Fatalf("unexpected categories: %v", response.Categories)
	}
}

func TestCreate(t *testing.T) {
	repo := &fakeRepository{}
	handler := NewHTTPHandler(categoryHandlers.BuildHandlers(repo))
	req := httptest.NewRequest(http.MethodPost, "/categories", strings.NewReader(`{"name":"Food"}`))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	handler.Create(recorder, req)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", recorder.Code)
	}

	var response CategoryDTO
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if response.Name != "Food" {
		t.Fatalf("unexpected category: %v", response)
	}
}

func TestCreateInvalidBody(t *testing.T) {
	repo := &fakeRepository{}
	handler := NewHTTPHandler(categoryHandlers.BuildHandlers(repo))
	req := httptest.NewRequest(http.MethodPost, "/categories", strings.NewReader(`{"name":`))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	handler.Create(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}
}
