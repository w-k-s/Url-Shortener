package urlshortener

import (
	"github.com/gorilla/mux"
)

var c *Controller

func init() {
	c = NewController()
}

func ConfigureRoutes(r *mux.Router) {
	s := r.PathPrefix("/urlshortener/v1").Subrouter()

	s.HandleFunc("/url", c.ShortenUrl).
		Methods("POST")
}
