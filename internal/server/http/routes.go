package http

import (
	"github.com/B-Dmitriy/expenses/pgk/web"
	"net/http"
)

func initRoutes(serv *http.ServeMux) *http.ServeMux {
	serv.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		web.WriteOK(w, "ping")
	})
	serv.HandleFunc("GET /json", func(w http.ResponseWriter, r *http.Request) {
		web.WriteOK(w, "{\"json\": \"data\"}")
	})
	serv.HandleFunc("GET /nil", func(w http.ResponseWriter, r *http.Request) {
		web.WriteOK(w, nil)
	})
	// /query?first=1&second=none&third=test
	serv.HandleFunc("GET /query", func(w http.ResponseWriter, r *http.Request) {
		params, _ := web.ParseQueryParams(r, "first", "third")

		web.WriteOK(w, params)
	})
	serv.HandleFunc("GET /user/{id}", func(w http.ResponseWriter, r *http.Request) {

		id, _ := web.ParseIDFromURL(r, "id")

		resp := struct {
			ID int `json:"id"`
		}{
			ID: id,
		}

		web.WriteOK(w, resp)
	})
	return serv
}
