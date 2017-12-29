package urlshortener

import (
	"github.com/gorilla/mux"
	a "github.com/waqqas-abdulkareem/short-url/app"
)

func Configure(app *a.App, r *mux.Router) {
	
	c := NewController(app)

	s := r.PathPrefix("/urlshortener/v1").Subrouter()

	s.HandleFunc("/url", c.ShortenUrl).
		Methods("POST")
}
