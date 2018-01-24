package main

import (
	"github.com/gorilla/mux"
	"github.com/w-k-s/short-url/app"
	"github.com/w-k-s/short-url/home"
	"github.com/w-k-s/short-url/urlshortener"
	"log"
	"net/http"
	"os"
)

var dbConnString string
var address string

func init() {

	address = os.Getenv("ADDRESS")
	if len(address) == 0 {
		address = ":8080"
	}

	dbConnString = os.Getenv("MONGO_ADDRESS")
	if len(dbConnString) == 0 {
		dbConnString = "mongodb://localhost:27017/shorturl"
	}

	log.Printf("Address: '%s'", address)
	log.Printf("Connection String: %s", dbConnString)
	log.Println("Init Complete")
}

func main() {

	app := app.Init(dbConnString)
	defer app.Session.Close()

	r := mux.NewRouter()

	//r.PathPrefix("/.well-known/").Handler(http.FileServer(http.Dir("./static/")))

	urlshortener.Configure(app, r)
	home.Configure(app, r)

	err := http.ListenAndServeTLS(address, "server.crt", "server.key", r)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
