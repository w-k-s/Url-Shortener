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

func init() {

	address = os.Getenv("ADDRESS")
	if len(address) == 0 {
		address = ":8080"
	}

	addressTLS = os.Getenv("ADDRESS_TLS")
	if len(addressTLS) == 0 {
		addressTLS = ":443"
	}

	dbConnString = os.Getenv("MONGO_ADDRESS")
	if len(dbConnString) == 0 {
		dbConnString = "mongodb://localhost:27017/shorturl"
	}

	production = os.Getenv("PROD") == "1"

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
	log.Printf("Production: %s", production)
	log.Printf("CertDir: %s", certDir)
	log.Printf("Domains: %s", domains)
	log.Println("Init Complete")
}

func main() {

	app := app.Init(dbConnString)
	defer app.Session.Close()

	r := mux.NewRouter()

	urlshortener.Configure(app, r)
	home.Configure(app, r)

	if production {
		initAutocertManager()
	}

	errchan := make(chan error, 1)
	listenAndServeHTTPServer(r, errchan)
	listenAndServeHTTPSServer(r, errchan)

	err := <-errchan
	log.Fatalf("Shutting down server with error: %s", err.Error())
}

func initAutocertManager() {
	hostPolicy := func(ctx context.Context, host string) error {
		if isHostAllowed(host) {
			return nil
		}

		log.Fatalf("acme/autocert: only '%s' host is allowed", domains)
		return fmt.Errorf("acme/autocert: only '%s' host is allowed", domains)
	}

	autocertManager = &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: hostPolicy,
		Cache:      autocert.DirCache(certDir),
	}
}

//HTTPS Server

func listenAndServeHTTPSServer(r *mux.Router, errchan chan error) {
	httpsServer := createServer(r, addressTLS)
	go func(c chan error) {
		if production {

			httpsServer.TLSConfig = &tls.Config{GetCertificate: autocertManager.GetCertificate}

			err := httpsServer.ListenAndServeTLS("", "")
			if err != nil {
				c <- fmt.Errorf("Error while configuring HTTPS Server (production mode): %s", err)
			}
		} else {
			err := httpsServer.ListenAndServeTLS("server.crt", "server.key")
			if err != nil {
				c <- fmt.Errorf("Error while configuring HTTPS Server (development mode): %s", err)
			}
		}
	}(errchan)
}

func isHostAllowed(host string) bool {
	ok, _ := domains[host]
	return ok
}

//HTTP Server

func listenAndServeHTTPServer(r *mux.Router, errchan chan error) {
	httpServer := createServer(r, address)
	go func(c chan error) {
		if autocertManager != nil {
			// allow autocert handle Let's Encrypt auth callbacks over HTTP.
			// it'll pass all other urls to our hanlder
			httpServer.Handler = autocertManager.HTTPHandler(httpServer.Handler)
		}
		err := httpServer.ListenAndServe()
		if err != nil {
			c <- fmt.Errorf("Error while configuring HTTP Server: %s", err)
		}
	}(errchan)
}

//Utils

func createServer(r *mux.Router, addr string) *http.Server {
	return &http.Server{
		Handler: r,
		Addr:    addr,
		// set timeouts so that a slow or malicious client doesn't
		// hold resources forever
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
}
