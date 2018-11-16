package main

import (
	"github.com/gorilla/mux"
	"github.com/w-k-s/short-url/app"
	"github.com/w-k-s/short-url/logging"
	"github.com/w-k-s/short-url/urlshortener"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var dbConnString string
var address string
var addressTLS string
var production bool
var certDir string
var domains map[string]bool

const STATIC_DIR string = "/public/"

func init() {
	production = os.Getenv("PROD") == "1"

	address = os.Getenv("ADDRESS")
	if len(address) == 0 {
		address = ":8080"
	}

	addressTLS = os.Getenv("ADDRESS_TLS")
	if len(addressTLS) == 0 {
		addressTLS = ":4430"
	}

	dbConnString = os.Getenv("MONGO_ADDRESS")
	if len(dbConnString) == 0 {
		dbConnString = "mongodb://localhost:27017/shorturl"
	}

	certDir = os.Getenv("CERT_DIR")
	if len(certDir) == 0 {
		certDir = "."
	}

	commaSeperatedDomains := os.Getenv("DOMAINS")
	domains = make(map[string]bool)
	for _, domain := range strings.Split(commaSeperatedDomains, ",") {
		domains[strings.Trim(domain, " ")] = true
	}

	log.Printf("Address: '%s'", address)
	log.Printf("AddressTLS: '%s'", addressTLS)
	log.Printf("Connection String: %s", dbConnString)
	log.Printf("Production: %v", production)
	log.Printf("CertDir: %s", certDir)
	log.Printf("Domains: %v", domains)
	log.Printf("Init Complete. Running on %s and %s", address, addressTLS)
}

func main() {
	logger := log.New(os.Stdout, "short-url: ", log.Ldate|log.Ltime)
	app := app.Init(logger, dbConnString)
	defer app.Db.Close()

	logRepository := logging.NewLogRepository(logger, app.Db)

	httpRouter := mux.NewRouter()
	httpRouter.Use(loggingMiddleware(logRepository))

	urlshortener.Configure(app, httpRouter)

	errchan := make(chan error, 1)
	listenAndServeHTTPServer(httpRouter, errchan)
	log.Fatalf("Error while configuring HTTP Server: %v", <-errchan)
}

func loggingMiddleware(logRepository *logging.LogRepository) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			sw := logging.StatusWriter{ResponseWriter: w}

			record := logRepository.LogRequest(r)

			next.ServeHTTP(sw, r)

			logRepository.LogResponse(sw, record)
		})
	}
}

//HTTP Server

func listenAndServeHTTPServer(h http.Handler, errchan chan error) {

	httpServer := createServer(h, address)
	go func(c chan error) {
		err := httpServer.ListenAndServe()
		if err != nil {
			errchan <- err
		}
	}(errchan)
}

//Utils

func createServer(h http.Handler, addr string) *http.Server {
	return &http.Server{
		Handler: h,
		Addr:    addr,
		// set timeouts so that a slow or malicious client doesn't
		// hold resources forever
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
}
