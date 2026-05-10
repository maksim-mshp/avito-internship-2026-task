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

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var request RegisterRequest
	apiErr := corehttp.ParseJSONBody(r, &request)
	if apiErr != nil {
		corehttp.RespondError(w, *apiErr)
		return
	}
	if request.Email == nil || request.Password == nil {
		corehttp.RespondError(w, corehttp.ErrInvalidRequest)
		return
	}

	role := ""
	if request.Role != nil {
		role = *request.Role
	}

	user, err := h.handlers.Register.Handle(r.Context(), app.RegisterCommand{
		Email:    *request.Email,
		Password: *request.Password,
		Role:     role,
	})
	if err != nil {
		log.Printf("failed to register user: %v", err)
		corehttp.RespondError(w, mapError(err))
		return
	}

	corehttp.Respond(w, http.StatusCreated, RegisterResponse{
		User: mapUser(user),
	})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var request LoginRequest
	apiErr := corehttp.ParseJSONBody(r, &request)
	if apiErr != nil {
		corehttp.RespondError(w, *apiErr)
		return
	}
	if request.Email == nil || request.Password == nil {
		corehttp.RespondError(w, corehttp.ErrInvalidRequest)
		return
	}

	result, err := h.handlers.Login.Handle(r.Context(), app.LoginCommand{
		Email:    *request.Email,
		Password: *request.Password,
	})
	if err != nil {
		log.Printf("failed to login user: %v", err)
		corehttp.RespondError(w, mapError(err))
		return
	}

	corehttp.Respond(w, http.StatusOK, TokenResponse{
		Token: result.Token,
		User:  mapUser(result.User),
	})
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
