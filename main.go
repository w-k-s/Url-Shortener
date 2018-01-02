package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/w-k-s/short-url/app"
	"github.com/w-k-s/short-url/home"
	"github.com/w-k-s/short-url/urlshortener"
	"log"
	"net/http"
)

// global flags
var port int
var dbConnString string

func init() {
	flag.IntVar(&port, "port", 8080, "Specify the port to listen to.")
	flag.StringVar(&dbConnString, "dbdsn", "mongodb://localhost/shorturl", "Specifies the MongoDB connection string")

	flag.Parse()

	log.Printf("Port: %d", 8080)
	log.Printf("Connection String: %s", dbConnString)
	log.Println("Init Complete")
}

func main() {

	app := app.Init(dbConnString)
	defer app.Session.Close()

	r := mux.NewRouter()

	home.Configure(app, r)
	urlshortener.Configure(app, r)

	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}
