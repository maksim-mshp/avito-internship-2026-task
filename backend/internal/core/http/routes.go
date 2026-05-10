package http

import "net/http"

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /_info", info)
}

func info(w http.ResponseWriter, r *http.Request) {
	Respond(w, http.StatusOK, nil)
}
