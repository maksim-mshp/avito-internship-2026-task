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
