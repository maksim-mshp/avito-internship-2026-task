package v1

import "net/http"

func RegisterRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("POST /dummyLogin", h.DummyLogin)
}
