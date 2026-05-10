package middleware

import (
	"log"
	"net/http"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (w *responseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *responseWriter) StatusCode() int {
	if w.status == 0 {
		return http.StatusOK
	}

	return w.status
}

func (w *responseWriter) Flush() {
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := &responseWriter{
			ResponseWriter: w,
		}

		next.ServeHTTP(rw, r)

		log.Printf("%s %s | %d | %s | %s",
			r.Method, r.URL.Path, rw.StatusCode(), r.RemoteAddr, r.UserAgent(),
		)
	})
}
