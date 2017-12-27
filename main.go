package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

type ShortenUrlRequest struct {
	LongUrl string `json:"longUrl"`
}

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/", Index).
		Methods("GET")

	shortenerRouter := r.PathPrefix("/urlshortener/v1").Subrouter()

	shortenerRouter.HandleFunc("/url", ShortenUrl).
		Methods("POST")

	http.ListenAndServe(":8080", r)
}

func Index(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Hello World")
}

func ShortenUrl(w http.ResponseWriter, req *http.Request) {

	var shortenReq ShortenUrlRequest
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&shortenReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	io.WriteString(w, fmt.Sprintf("The url is %s", shortenReq.LongUrl))
}
