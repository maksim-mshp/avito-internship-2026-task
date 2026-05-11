package handlers

import (
	"context"

	"ai-assistants-catalog/internal/assistants/app"
	"ai-assistants-catalog/internal/assistants/domain"
)

type AddFavoriteHandler struct {
	repo app.Repository
}

func NewAddFavoriteHandler(repo app.Repository) *AddFavoriteHandler {
	return &AddFavoriteHandler{repo: repo}
}

func (h *AddFavoriteHandler) Handle(ctx context.Context, cmd app.FavoriteCommand) error {
	if !domain.IsValidID(cmd.UserID) {
		return domain.ErrInvalidID
	}

	if !domain.IsValidID(cmd.AssistantID) {
		return domain.ErrInvalidID
	}

	return h.repo.AddFavorite(ctx, cmd.UserID, cmd.AssistantID, cmd.IncludeInactive)
}

type RemoveFavoriteHandler struct {
	repo app.Repository
}

func NewRemoveFavoriteHandler(repo app.Repository) *RemoveFavoriteHandler {
	return &RemoveFavoriteHandler{repo: repo}
}

func (h *RemoveFavoriteHandler) Handle(ctx context.Context, cmd app.FavoriteCommand) error {
	if !domain.IsValidID(cmd.UserID) {
		return domain.ErrInvalidID
	}

	if !domain.IsValidID(cmd.AssistantID) {
		return domain.ErrInvalidID
	}

	return h.repo.RemoveFavorite(ctx, cmd.UserID, cmd.AssistantID)
}
