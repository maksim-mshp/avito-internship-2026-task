package v1

import (
	"log"
	"net/http"

	"ai-assistants-catalog/internal/auth/app"
	"ai-assistants-catalog/internal/auth/app/handlers"
	corehttp "ai-assistants-catalog/internal/core/http"
)

type Handler struct {
	handlers *handlers.Handlers
}

func NewHTTPHandler(handlers *handlers.Handlers) *Handler {
	return &Handler{handlers: handlers}
}

func (h *Handler) DummyLogin(w http.ResponseWriter, r *http.Request) {
	var request DummyLoginRequest
	apiErr := corehttp.ParseJSONBody(r, &request)
	if apiErr != nil {
		corehttp.RespondError(w, *apiErr)
		return
	}

	result, err := h.handlers.DummyLogin.Handle(r.Context(), app.DummyLoginCommand{
		Role: request.Role,
	})
	if err != nil {
		log.Printf("failed to build dummy login token: %v", err)
		corehttp.RespondError(w, mapError(err))
		return
	}

	corehttp.Respond(w, http.StatusOK, TokenResponse{
		Token: result.Token,
		User:  mapUser(result.User),
	})
}
