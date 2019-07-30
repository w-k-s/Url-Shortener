package main

import (
	a "github.com/w-k-s/short-url/app"
	"github.com/w-k-s/short-url/logging"
	u "github.com/w-k-s/short-url/urlshortener"
	"log"
	"net/http"
)

var app *a.App

func init() {
	app = a.Init()
}

func main() {
	defer app.Close()

	configureURLController()
	configureLoggingMiddleware()

	log.Panic(app.ListenAndServe())
}

func configureURLController() {
	urlRepo := u.NewURLRepository(app.Db(), app.Logger())
	urlService := u.NewService(urlRepo, app.Logger(), u.DefaultShortIDGenerator{})
	urlController := u.NewController(urlService)

	app.Post("/urlshortener/v1/url", urlController.ShortenURL)
	app.Get("/urlshortener/v1/url", urlController.GetLongURL)
	app.Get("/{shortUrl}", urlController.RedirectToLongURL)
}

func configureLoggingMiddleware() {
	logRepository := logging.NewLogRepository(app.Logger(), app.Db())
	app.Middleware(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			sw := &logging.StatusWriter{ResponseWriter: w}

			record := logRepository.LogRequest(r)

			next.ServeHTTP(sw, r)

			logRepository.LogResponse(sw, record)
		})
	})
}
