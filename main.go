package main

import (
	"github.com/gorilla/mux"
	"github.com/waqqas-abdulkareem/short-url/home"
	"github.com/waqqas-abdulkareem/short-url/urlshortener"
	"net/http"
)

func main() {

	r := mux.NewRouter()

	home.ConfigureRoutes(r)
	urlshortener.ConfigureRoutes(r)

	http.ListenAndServe(":8080", r)
}
