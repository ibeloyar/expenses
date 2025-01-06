package http

import (
	"net/http"

	"github.com/ibeloyar/expenses/pgk/web"
)

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Authorization, Origin")
		if r.Method == "OPTIONS" {
			CorsOptionHandlerFunc(w, r)
		}
		next.ServeHTTP(w, r)
	})
}

func CorsOptionHandlerFunc(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write(nil) // Write code 200
	if err != nil {
		web.WriteServerError(w)
	}
	return
}
