package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/w-k-s/short-url/app"
	"github.com/w-k-s/short-url/home"
	"github.com/w-k-s/short-url/urlshortener"
	"golang.org/x/crypto/acme/autocert"
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
var autocertManager *autocert.Manager

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
	app := app.Init(dbConnString)
	defer app.Db.Close()

	httpsRerouter := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://"+r.Host+r.URL.String(), http.StatusMovedPermanently)
	})

	mainRouter := mux.NewRouter()

	urlshortener.Configure(app, mainRouter)
	home.Configure(app, mainRouter)

	mainRouter.
		PathPrefix(STATIC_DIR).
		Handler(http.StripPrefix(STATIC_DIR, http.FileServer(http.Dir("."+STATIC_DIR))))

	var httpRouter http.Handler
	if production {
		initAutocertManager()
		httpRouter = httpsRerouter
	} else {
		httpRouter = mainRouter
	}

	errchan := make(chan error, 1)
	listenAndServeHTTPServer(httpRouter, errchan)
	listenAndServeHTTPSServer(mainRouter, errchan)

	err := <-errchan
	log.Fatalf("Shutting down server with error: %s", err)
}

//HTTPS Server

func initAutocertManager() {
	hostPolicy := func(ctx context.Context, host string) error {
		if isHostAllowed(host) {
			return nil
		}

		err := fmt.Errorf("acme/autocert: host '%v' not allowed. Allowed domains: '%v'", host, domains)
		log.Println(err)
		return err
	}

	autocertManager = &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: hostPolicy,
		Cache:      autocert.DirCache(certDir),
	}
}

func listenAndServeHTTPSServer(h http.Handler, errchan chan error) {

	httpsServer := createServer(h, addressTLS)

	go func(c chan error) {
		if production {

			httpsServer.TLSConfig = &tls.Config{GetCertificate: autocertManager.GetCertificate}

			err := httpsServer.ListenAndServeTLS("", "")
			if err != nil {
				c <- fmt.Errorf("Error while configuring HTTPS Server (production mode): %v", err)
			}

		} else {

			err := httpsServer.ListenAndServeTLS("server.crt", "server.key")
			if err != nil {
				c <- fmt.Errorf("Error while configuring HTTPS Server (development mode): %v", err)
			}

		}

	}(errchan)
}

func isHostAllowed(host string) bool {
	ok, _ := domains[host]
	return ok
}

//HTTP Server

func listenAndServeHTTPServer(h http.Handler, errchan chan error) {

	httpServer := createServer(h, address)
	go func(c chan error) {
		if autocertManager != nil {
			// allow autocert handle Let's Encrypt auth callbacks over HTTP.
			// it'll pass all other urls to our hanlder
			httpServer.Handler = autocertManager.HTTPHandler(httpServer.Handler)
		}
		err := httpServer.ListenAndServe()
		if err != nil {
			c <- fmt.Errorf("Error while configuring HTTP Server: %v", err)
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
