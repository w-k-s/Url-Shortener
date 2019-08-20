package main

import (
	d "github.com/w-k-s/short-url/adapters/db"
	"github.com/w-k-s/short-url/adapters/logging"
	"github.com/w-k-s/short-url/adapters/web"
	"github.com/w-k-s/short-url/domain/urlshortener/usecase"
	"log"
	"net/http"
	"os"
)

var app *web.App
var db *d.Db

func init() {
	app = web.Init()
	db = d.New(os.Getenv("DB_CONN_STRING"), false)
}

func main() {
	defer db.Close()

	configureURLController()
	configureLoggingMiddleware()

	log.Panic(app.ListenAndServe())
}

func configureURLController() {
	urlRepo := d.NewURLRepository(db, app.Logger())

	shortenURLUseCase := usecase.NewShortenURLUseCase(urlRepo, usecase.DefaultShortIDGenerator{}, app.Logger())
	retrieveOriginalURLUseCase := usecase.NewRetrieveOriginalURLUseCase(urlRepo, app.Logger())

	urlController := web.NewController(shortenURLUseCase, retrieveOriginalURLUseCase, app.Logger())

	app.Post("/urlshortener/v1/url", urlController.ShortenURL)
	app.Get("/urlshortener/v1/url", urlController.GetLongURL)
	app.Get("/{shortUrl}", urlController.RedirectToLongURL)
}

func configureLoggingMiddleware() {
	logRepository := logging.NewLogRepository(app.Logger(), db)
	app.Middleware(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			sw := &logging.StatusWriter{ResponseWriter: w}

			record := logRepository.LogRequest(r)

			next.ServeHTTP(sw, r)

			logRepository.LogResponse(sw, record)
		})
	})
}
