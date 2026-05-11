package http

import "net/http"

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /_info", info)
}

// @Summary Состояние сервиса
// @Tags System
// @Success 200 "OK"
// @Router /_info [GET]
func info(w http.ResponseWriter, r *http.Request) {
	Respond(w, http.StatusOK, nil)
}
