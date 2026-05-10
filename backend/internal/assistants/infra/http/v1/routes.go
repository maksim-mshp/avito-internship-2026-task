package v1

import "net/http"

func RegisterRoutes(
	mux *http.ServeMux,
	h *Handler,
	authMW func(http.Handler) http.Handler,
	adminMW func(http.Handler) http.Handler,
) {
	mux.Handle("GET /assistants", authMW(http.HandlerFunc(h.List)))
	mux.Handle("POST /assistants", authMW(adminMW(http.HandlerFunc(h.Create))))
	mux.Handle("GET /assistants/{assistantId}", authMW(http.HandlerFunc(h.GetByID)))
	mux.Handle("PUT /assistants/{assistantId}", authMW(adminMW(http.HandlerFunc(h.Update))))
}
