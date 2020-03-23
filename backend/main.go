package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	d "github.com/w-k-s/short-url/adapters/db"
	"github.com/w-k-s/short-url/adapters/logging"
	"github.com/w-k-s/short-url/adapters/web"
	"github.com/w-k-s/short-url/domain/urlshortener/usecase"
	"log"
	"net/http"
	"net/url"
	"os"
)

var app *web.App
var db *sql.DB
var baseURL *url.URL

func init() {
	var err error

	baseURL, err = url.Parse(os.Getenv("BASE_URL"))
	if err != nil {
		log.Fatalf("Failed to parse env variable 'BASE_URL': '%s'", os.Getenv("BASE_URL"))
	}
	if len(baseURL.Scheme) == 0 {
		log.Fatalf("Failed to determine scheme from BASE_URL %q", baseURL)
	}
	if len(baseURL.Host) == 0 {
		log.Fatalf("Failed to determine host from BASE_URL %q", baseURL)
	}

	connStr := os.Getenv("DB_CONN_STRING")
	if len(connStr) == 0 {
		log.Fatal("Connection String is required")
	}

	db, err = sql.Open("postgres", connStr)

	if err != nil {
		log.Fatalf("Failed to open db with connection string %q: %s", connStr, err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to ping db with connection string %q: %s", connStr, err)
	}

	app = web.Init()
}

func main() {
	configureHealthCheck()
	configureURLController()
	configureLoggingMiddleware()

	log.Panic(app.ListenAndServe())
}

func configureHealthCheck() {
	app.Get("/health", func(w http.ResponseWriter, req *http.Request) {
		if err := db.Ping(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}

func configureURLController() {
	urlRepo := d.NewURLRepository(db, app.Logger())

	shortenURLUseCase := usecase.NewShortenURLUseCase(urlRepo, baseURL, usecase.DefaultShortIDGenerator{}, app.Logger())
	retrieveOriginalURLUseCase := usecase.NewRetrieveOriginalURLUseCase(urlRepo, app.Logger())

	urlController := web.NewController(shortenURLUseCase, retrieveOriginalURLUseCase, app.Logger())

	app.Post("/urlshortener/v1/url", urlController.ShortenURL)
	app.Get("/urlshortener/v1/url", urlController.GetLongURL)
	app.Get("/{shortUrl}", urlController.RedirectToLongURL)
}

func configureLoggingMiddleware() {
	logRepository := logging.NewLogRepository(db, app.Logger())
	app.Middleware(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			sw := &logging.StatusWriter{ResponseWriter: w}

			record := logRepository.LogRequest(r)

			next.ServeHTTP(sw, r)

			logRepository.LogResponse(sw, record)
		})
	})
}