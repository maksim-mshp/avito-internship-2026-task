package v1

import (
	"log"
	"net/http"

	"ai-assistants-catalog/internal/categories/app"
	"ai-assistants-catalog/internal/categories/app/handlers"
	corehttp "ai-assistants-catalog/internal/core/http"
)

type Handler struct {
	handlers *handlers.Handlers
}

func NewHTTPHandler(handlers *handlers.Handlers) *Handler {
	return &Handler{handlers: handlers}
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	categories, err := h.handlers.GetAll.Handle(r.Context(), app.GetAllQuery{})
	if err != nil {
		log.Printf("failed to get categories: %v", err)
		corehttp.RespondError(w, mapError(err))
		return
	}

	corehttp.Respond(w, http.StatusOK, CategoriesResponse{
		Categories: mapCategories(categories),
	})
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var request CategoryCreateRequest
	apiErr := corehttp.ParseJSONBody(r, &request)
	if apiErr != nil {
		corehttp.RespondError(w, *apiErr)
		return
	}

	category, err := h.handlers.Create.Handle(r.Context(), app.CreateCommand{
		Name:        request.Name,
		Description: request.Description,
	})
	if err != nil {
		log.Printf("failed to create category: %v", err)
		corehttp.RespondError(w, mapError(err))
		return
	}

	corehttp.Respond(w, http.StatusCreated, mapCategory(category))
}
