package v1

import (
	"log"
	"net/http"
	"strconv"

	corehttp "ai-assistants-catalog/internal/core/http"
	"ai-assistants-catalog/internal/core/security"
	"ai-assistants-catalog/internal/runs/app"
	"ai-assistants-catalog/internal/runs/app/handlers"
)

type Handler struct {
	handlers *handlers.Handlers
}

func NewHTTPHandler(handlers *handlers.Handlers) *Handler {
	return &Handler{handlers: handlers}
}

// @Summary Запустить ассистента
// @Tags Runs
// @Security BearerAuth
// @Param assistantId path string true "ID ассистента"
// @Param request body RunCreateRequest true "RunCreateRequest"
// @Success 201 {object} RunDTO "Запуск создан"
// @Failure 400 {object} corehttp.ErrorResponse "Некорректный запрос"
// @Failure 401 {object} corehttp.ErrorResponse "Нет авторизации"
// @Failure 404 {object} corehttp.ErrorResponse "Ассистент не найден"
// @Failure 409 {object} corehttp.ErrorResponse "Ассистент выключен"
// @Failure 500 {object} corehttp.ErrorResponse "Внутренняя ошибка"
// @Failure 502 {object} corehttp.ErrorResponse "Ошибка LLM-провайдера"
// @Router /assistants/{assistantId}/run [POST]
func (h *Handler) RunAssistant(w http.ResponseWriter, r *http.Request) {
	claims, ok := security.ClaimsFromContext(r.Context())
	if !ok {
		corehttp.RespondError(w, corehttp.ErrUnauthorized)
		return
	}

	var request RunCreateRequest
	apiErr := corehttp.ParseJSONBody(r, &request)
	if apiErr != nil {
		corehttp.RespondError(w, *apiErr)
		return
	}

	run, err := h.handlers.RunAssistant.Handle(r.Context(), app.RunAssistantCommand{
		AssistantID: r.PathValue("assistantId"),
		UserID:      claims.UserID,
		UserPrompt:  request.UserPrompt,
	})
	if err != nil {
		log.Printf("failed to run assistant: %v", err)
		corehttp.RespondError(w, mapError(err))
		return
	}

	corehttp.Respond(w, http.StatusCreated, mapRun(run))
}

// @Summary История запусков пользователя
// @Tags Runs
// @Security BearerAuth
// @Param status query string false "Статус"
// @Param page query int false "Страница"
// @Param pageSize query int false "Размер страницы"
// @Success 200 {object} RunsResponse "История запусков"
// @Failure 400 {object} corehttp.ErrorResponse "Некорректный запрос"
// @Failure 401 {object} corehttp.ErrorResponse "Нет авторизации"
// @Failure 500 {object} corehttp.ErrorResponse "Внутренняя ошибка"
// @Router /runs/my [GET]
func (h *Handler) ListMy(w http.ResponseWriter, r *http.Request) {
	claims, ok := security.ClaimsFromContext(r.Context())
	if !ok {
		corehttp.RespondError(w, corehttp.ErrUnauthorized)
		return
	}

	query, apiErr := parseListQuery(r)
	if apiErr != nil {
		corehttp.RespondError(w, *apiErr)
		return
	}

	result, err := h.handlers.ListMy.Handle(r.Context(), app.ListMyQuery{
		UserID:   claims.UserID,
		Status:   query.Status,
		Page:     query.Page,
		PageSize: query.PageSize,
	})
	if err != nil {
		log.Printf("failed to list user runs: %v", err)
		corehttp.RespondError(w, mapError(err))
		return
	}

	corehttp.Respond(w, http.StatusOK, mapListResult(result))
}

// @Summary Все запуски ассистентов
// @Tags Runs
// @Security BearerAuth
// @Param assistantId query string false "ID ассистента"
// @Param status query string false "Статус"
// @Param page query int false "Страница"
// @Param pageSize query int false "Размер страницы"
// @Success 200 {object} RunsResponse "Список запусков"
// @Failure 400 {object} corehttp.ErrorResponse "Некорректный запрос"
// @Failure 401 {object} corehttp.ErrorResponse "Нет авторизации"
// @Failure 403 {object} corehttp.ErrorResponse "Недостаточно прав"
// @Failure 500 {object} corehttp.ErrorResponse "Внутренняя ошибка"
// @Router /admin/runs [GET]
func (h *Handler) ListAdmin(w http.ResponseWriter, r *http.Request) {
	query, apiErr := parseListQuery(r)
	if apiErr != nil {
		corehttp.RespondError(w, *apiErr)
		return
	}

	result, err := h.handlers.ListAdmin.Handle(r.Context(), app.ListAdminQuery{
		AssistantID: optionalQuery(r.URL.Query().Get("assistantId")),
		Status:      query.Status,
		Page:        query.Page,
		PageSize:    query.PageSize,
	})
	if err != nil {
		log.Printf("failed to list admin runs: %v", err)
		corehttp.RespondError(w, mapError(err))
		return
	}

	corehttp.Respond(w, http.StatusOK, mapListResult(result))
}

type listQuery struct {
	Status   *string
	Page     int
	PageSize int
}

func parseListQuery(r *http.Request) (listQuery, *corehttp.APIError) {
	values := r.URL.Query()

	page, err := parseIntQuery(values.Get("page"), 1)
	if err != nil {
		return listQuery{}, &corehttp.ErrInvalidRequest
	}

	pageSize, err := parseIntQuery(values.Get("pageSize"), 20)
	if err != nil {
		return listQuery{}, &corehttp.ErrInvalidRequest
	}

	return listQuery{
		Status:   optionalQuery(values.Get("status")),
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func parseIntQuery(value string, defaultValue int) (int, error) {
	if value == "" {
		return defaultValue, nil
	}

	return strconv.Atoi(value)
}

func optionalQuery(value string) *string {
	if value == "" {
		return nil
	}

	return &value
}
