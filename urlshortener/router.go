package urlshortener

import (
	"github.com/gorilla/mux"
	a "github.com/w-k-s/short-url/app"
)

func Configure(app *a.App, r *mux.Router) {

	c := NewController(NewService(
		NewURLRepository(app.Db),
		app.Logger,
		DefaultShortIdGenerator{},
	))

	r.HandleFunc("/{shortUrl}", c.RedirectToLongURL).
		Methods("GET")

	s := r.PathPrefix("/urlshortener/v1").
		Subrouter()

	s.HandleFunc("/url", c.ShortenURL).
		Methods("POST")

	s.HandleFunc("/url", c.GetLongURL).
		Methods("GET")
}
