package v1

import "net/http"

func RegisterRoutes(
	mux *http.ServeMux,
	h *Handler,
	authMW func(http.Handler) http.Handler,
	adminMW func(http.Handler) http.Handler,
) {
	mux.Handle("POST /assistants/{assistantId}/run", authMW(http.HandlerFunc(h.RunAssistant)))
	mux.Handle("GET /runs/my", authMW(http.HandlerFunc(h.ListMy)))
	mux.Handle("GET /admin/runs", authMW(adminMW(http.HandlerFunc(h.ListAdmin))))
}
