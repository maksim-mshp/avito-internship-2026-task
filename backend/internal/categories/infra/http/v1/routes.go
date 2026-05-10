package v1

import "net/http"

func RegisterRoutes(
	mux *http.ServeMux,
	h *Handler,
	authMW func(http.Handler) http.Handler,
	adminMW func(http.Handler) http.Handler,
) {
	mux.Handle("GET /categories", authMW(http.HandlerFunc(h.GetAll)))
	mux.Handle("POST /categories", authMW(adminMW(http.HandlerFunc(h.Create))))
}
