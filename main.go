package main

import (
	
	"github.com/w-k-s/short-url/app"
	"github.com/w-k-s/short-url/logging"
	"github.com/w-k-s/short-url/urlshortener"
	"log"
	"net/http"
	"os"
	"time"
)

var app app.App

func init() {
	app = app.Init()
}

func main() {
	defer app.Close()

	logRepository := logging.NewLogRepository(app.Logger(), app.Db())

	app.Middleware(loggingMiddleware(logRepository))

	urlshortener.Configure(app, httpRouter)

	errchan := make(chan error, 1)
	app.ListenAndServe(errchan)
	log.Fatalf("Error while configuring HTTP Server: %v", <-errchan)
}

func loggingMiddleware(logRepository *logging.LogRepository) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			sw := &logging.StatusWriter{ResponseWriter: w}

			record := logRepository.LogRequest(r)

			next.ServeHTTP(sw, r)

			logRepository.LogResponse(sw, record)
		})
	}
}
