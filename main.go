package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/w-k-s/short-url/app"
	"github.com/w-k-s/short-url/home"
	"github.com/w-k-s/short-url/urlshortener"
	"log"
	"net/http"
	"os"
	"strings"
)

var dbConnString string
var address string

func init() {

	address = os.Getenv("ADDRESS")
	if len(address) == 0 {
		address = ":8080"
	}

	// MONGO_PORT is an env variable created by docker
	// when the web app container is linked to a container named 'mongo'
	// MONGO_PORT is the ip address of the container
	dbAddress := os.Getenv("MONGO_PORT")
	if len(dbAddress) != 0 {
		dbAddress := strings.Replace(dbAddress, "tcp", "mongodb", 1)
		dbConnString = fmt.Sprintf("%s/shorturl", dbAddress)
	} else {
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

	urlshortener.Configure(app, r)
	home.Configure(app, r)

	http.ListenAndServe(address, r)
}
