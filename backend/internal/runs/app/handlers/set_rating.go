package handlers

import (
	"context"
	"strings"

	"ai-assistants-catalog/internal/runs/app"
	"ai-assistants-catalog/internal/runs/domain"
)

type SetRatingHandler struct {
	repo app.RunRepository
}

func NewSetRatingHandler(repo app.RunRepository) *SetRatingHandler {
	return &SetRatingHandler{repo: repo}
}

func (h *SetRatingHandler) Handle(ctx context.Context, cmd app.SetRatingCommand) (domain.Run, error) {
	if !domain.IsValidID(cmd.ID) {
		return domain.Run{}, domain.ErrInvalidID
	}

	if !domain.IsValidID(cmd.UserID) {
		return domain.Run{}, domain.ErrInvalidUserID
	}

	if cmd.Rating == nil {
		return domain.Run{}, domain.ErrInvalidRating
	}

	rating := strings.TrimSpace(*cmd.Rating)
	if !domain.IsValidRating(rating) {
		return domain.Run{}, domain.ErrInvalidRating
	}

	return h.repo.SetRating(ctx, cmd.ID, cmd.UserID, rating)
}
