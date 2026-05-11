package v1

import (
	"log"
	"net/http"
	"strconv"

	"ai-assistants-catalog/internal/assistants/app"
	"ai-assistants-catalog/internal/assistants/app/handlers"
	corehttp "ai-assistants-catalog/internal/core/http"
	"ai-assistants-catalog/internal/core/security"
)

type Handler struct {
	handlers *handlers.Handlers
}

func NewHTTPHandler(handlers *handlers.Handlers) *Handler {
	return &Handler{handlers: handlers}
}

// @Summary Список ассистентов
// @Tags Assistants
// @Security BearerAuth
// @Param categoryId query string false "ID категории"
// @Param q query string false "Поиск"
// @Param tag query string false "Тег"
// @Param includeInactive query bool false "Показывать неактивных"
// @Param page query int false "Страница"
// @Param pageSize query int false "Размер страницы"
// @Success 200 {object} AssistantsResponse "Список ассистентов"
// @Failure 400 {object} corehttp.ErrorResponse "Некорректный запрос"
// @Failure 401 {object} corehttp.ErrorResponse "Нет авторизации"
// @Failure 500 {object} corehttp.ErrorResponse "Внутренняя ошибка"
// @Router /assistants [GET]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	query, apiErr := parseListQuery(r)
	if apiErr != nil {
		corehttp.RespondError(w, *apiErr)
		return
	}

	result, err := h.handlers.List.Handle(r.Context(), query)
	if err != nil {
		log.Printf("failed to list assistants: %v", err)
		corehttp.RespondError(w, mapError(err))
		return
	}

	corehttp.Respond(w, http.StatusOK, mapListResult(result))
}

// @Summary Получить ассистента
// @Tags Assistants
// @Security BearerAuth
// @Param assistantId path string true "ID ассистента"
// @Success 200 {object} AssistantDTO "Ассистент"
// @Failure 401 {object} corehttp.ErrorResponse "Нет авторизации"
// @Failure 404 {object} corehttp.ErrorResponse "Ассистент не найден"
// @Failure 500 {object} corehttp.ErrorResponse "Внутренняя ошибка"
// @Router /assistants/{assistantId} [GET]
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	claims, ok := security.ClaimsFromContext(r.Context())
	if !ok {
		corehttp.RespondError(w, corehttp.ErrUnauthorized)
		return
	}

	assistant, err := h.handlers.GetByID.Handle(r.Context(), app.GetByIDQuery{
		ID:              r.PathValue("assistantId"),
		IncludeInactive: claims.Role == security.RoleAdmin,
	})
	if err != nil {
		log.Printf("failed to get assistant: %v", err)
		corehttp.RespondError(w, mapError(err))
		return
	}

	corehttp.Respond(w, http.StatusOK, mapAssistant(assistant))
}

// @Summary Создать ассистента
// @Tags Assistants
// @Security BearerAuth
// @Param request body AssistantCreateRequest true "AssistantCreateRequest"
// @Success 201 {object} AssistantDTO "Ассистент создан"
// @Failure 400 {object} corehttp.ErrorResponse "Некорректный запрос"
// @Failure 401 {object} corehttp.ErrorResponse "Нет авторизации"
// @Failure 403 {object} corehttp.ErrorResponse "Недостаточно прав"
// @Failure 500 {object} corehttp.ErrorResponse "Внутренняя ошибка"
// @Router /assistants [POST]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var request AssistantCreateRequest
	apiErr := corehttp.ParseJSONBody(r, &request)
	if apiErr != nil {
		corehttp.RespondError(w, *apiErr)
		return
	}

	assistant, err := h.handlers.Create.Handle(r.Context(), app.CreateCommand{
		CategoryID:        request.CategoryID,
		Name:              request.Name,
		Description:       request.Description,
		Model:             request.Model,
		SystemPrompt:      request.SystemPrompt,
		ExampleUserPrompt: request.ExampleUserPrompt,
		Tags:              request.Tags,
		IsActive:          request.IsActive,
	})
	if err != nil {
		log.Printf("failed to create assistant: %v", err)
		corehttp.RespondError(w, mapError(err))
		return
	}

	corehttp.Respond(w, http.StatusCreated, mapAssistant(assistant))
}

// @Summary Обновить ассистента
// @Tags Assistants
// @Security BearerAuth
// @Param assistantId path string true "ID ассистента"
// @Param request body AssistantUpdateRequest true "AssistantUpdateRequest"
// @Success 200 {object} AssistantDTO "Ассистент обновлен"
// @Failure 400 {object} corehttp.ErrorResponse "Некорректный запрос"
// @Failure 401 {object} corehttp.ErrorResponse "Нет авторизации"
// @Failure 403 {object} corehttp.ErrorResponse "Недостаточно прав"
// @Failure 404 {object} corehttp.ErrorResponse "Ассистент не найден"
// @Failure 500 {object} corehttp.ErrorResponse "Внутренняя ошибка"
// @Router /assistants/{assistantId} [PUT]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	var request AssistantUpdateRequest
	apiErr := corehttp.ParseJSONBody(r, &request)
	if apiErr != nil {
		corehttp.RespondError(w, *apiErr)
		return
	}

	assistant, err := h.handlers.Update.Handle(r.Context(), app.UpdateCommand{
		ID:                r.PathValue("assistantId"),
		CategoryID:        request.CategoryID,
		Name:              request.Name,
		Description:       request.Description,
		Model:             request.Model,
		SystemPrompt:      request.SystemPrompt,
		ExampleUserPrompt: request.ExampleUserPrompt,
		Tags:              request.Tags,
		IsActive:          request.IsActive,
	})
	if err != nil {
		log.Printf("failed to update assistant: %v", err)
		corehttp.RespondError(w, mapError(err))
		return
	}

	corehttp.Respond(w, http.StatusOK, mapAssistant(assistant))
}

func parseListQuery(r *http.Request) (app.ListQuery, *corehttp.APIError) {
	values := r.URL.Query()

	page, err := parseIntQuery(values.Get("page"), 1)
	if err != nil {
		return app.ListQuery{}, &corehttp.ErrInvalidRequest
	}

	pageSize, err := parseIntQuery(values.Get("pageSize"), 10)
	if err != nil {
		return app.ListQuery{}, &corehttp.ErrInvalidRequest
	}

	includeInactive, err := parseBoolQuery(values.Get("includeInactive"), false)
	if err != nil {
		return app.ListQuery{}, &corehttp.ErrInvalidRequest
	}

	claims, ok := security.ClaimsFromContext(r.Context())
	if !ok {
		return app.ListQuery{}, &corehttp.ErrUnauthorized
	}

	if claims.Role != security.RoleAdmin {
		includeInactive = false
	}

	return app.ListQuery{
		CategoryID:      optionalQuery(values.Get("categoryId")),
		Search:          optionalQuery(values.Get("q")),
		Tag:             optionalQuery(values.Get("tag")),
		IncludeInactive: includeInactive,
		Page:            page,
		PageSize:        pageSize,
	}, nil
}

func parseIntQuery(value string, defaultValue int) (int, error) {
	if value == "" {
		return defaultValue, nil
	}

	return strconv.Atoi(value)
}

func parseBoolQuery(value string, defaultValue bool) (bool, error) {
	if value == "" {
		return defaultValue, nil
	}

	return strconv.ParseBool(value)
}

func optionalQuery(value string) *string {
	if value == "" {
		return nil
	}

	return &value
}
