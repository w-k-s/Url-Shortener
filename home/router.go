package home

import (
	"github.com/gorilla/mux"
	a "github.com/w-k-s/short-url/app"
)

func Configure(app *a.App, r *mux.Router) {

	c := NewController(app)

	r.HandleFunc("/", c.Index).
		Methods("GET")
}
