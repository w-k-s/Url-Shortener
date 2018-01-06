package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/w-k-s/short-url/app"
	"github.com/w-k-s/short-url/home"
	"github.com/w-k-s/short-url/urlshortener"
	"log"
	"os"
	"net/http"
	"strings"
)

// global flags
var port int

func init() {
	flag.IntVar(&port, "port", 8080, "Specify the port to listen to.")
	
	flag.Parse()

	dbIp :=os.Getenv("MONGO_PORT")
	dbIp = strings.Replace(dbIp,"tcp","mongodb",1)
	dbConnString = fmt.Sprintf("%s/shorturl",dbIp)

	log.Printf("Port: %d", port)
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
