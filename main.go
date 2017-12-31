package main

import (
	"github.com/gorilla/mux"
	"github.com/w-k-s/short-url/home"
	"github.com/w-k-s/short-url/urlshortener"
	"github.com/w-k-s/short-url/app"
	"net/http"
)

func main() {

	app := app.Init()
	defer app.Session.Close()

	r := mux.NewRouter()

	home.Configure(app, r)
	urlshortener.Configure(app, r)

	http.ListenAndServe(":8080", r)
}