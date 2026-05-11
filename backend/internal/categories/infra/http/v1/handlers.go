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

// @Summary Список категорий
// @Tags Categories
// @Security BearerAuth
// @Success 200 {object} CategoriesResponse "Список категорий"
// @Failure 401 {object} corehttp.ErrorResponse "Нет авторизации"
// @Failure 500 {object} corehttp.ErrorResponse "Внутренняя ошибка"
// @Router /categories [GET]
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

// @Summary Создать категорию
// @Tags Categories
// @Security BearerAuth
// @Param request body CategoryCreateRequest true "CategoryCreateRequest"
// @Success 201 {object} CategoryDTO "Категория создана"
// @Failure 400 {object} corehttp.ErrorResponse "Некорректный запрос"
// @Failure 401 {object} corehttp.ErrorResponse "Нет авторизации"
// @Failure 403 {object} corehttp.ErrorResponse "Недостаточно прав"
// @Failure 500 {object} corehttp.ErrorResponse "Внутренняя ошибка"
// @Router /categories [POST]
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
