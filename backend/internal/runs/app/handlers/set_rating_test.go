package handlers

import (
	"context"
	"errors"
	"testing"

	"ai-assistants-catalog/internal/runs/app"
	"ai-assistants-catalog/internal/runs/domain"
)

func TestSetRatingHandlerSuccess(t *testing.T) {
	repo := &fakeRunRepository{}
	handler := NewSetRatingHandler(repo)
	rating := domain.RatingLike

	run, err := handler.Handle(context.Background(), app.SetRatingCommand{
		ID:     testRunID,
		UserID: testUserID,
		Rating: &rating,
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !repo.ratingCall {
		t.Fatalf("expected repository call")
	}
	if run.Rating == nil || *run.Rating != domain.RatingLike {
		t.Fatalf("unexpected rating: %v", run.Rating)
	}
}

func TestSetRatingHandlerInvalidRating(t *testing.T) {
	repo := &fakeRunRepository{}
	handler := NewSetRatingHandler(repo)
	rating := "bad"

	_, err := handler.Handle(context.Background(), app.SetRatingCommand{
		ID:     testRunID,
		UserID: testUserID,
		Rating: &rating,
	})
	if !errors.Is(err, domain.ErrInvalidRating) {
		t.Fatalf("expected ErrInvalidRating, got %v", err)
	}
	if repo.ratingCall {
		t.Fatalf("expected no repository call")
	}
}
