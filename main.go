package main

import (
	d "github.com/w-k-s/short-url/adapters/db"
	"github.com/w-k-s/short-url/adapters/logging"
	"github.com/w-k-s/short-url/adapters/web"
	"github.com/w-k-s/short-url/domain/urlshortener/usecase"
	"log"
	"net/http"
	"os"
	"database/sql"
	_ "github.com/lib/pq"
)

var app *web.App
var db *sql.DB

func init() {
	connStr := os.Getenv("DB_CONN_STRING")
	if len(connStr) == 0 {
		connStr = "postgres://localhost/url_shortener_test?sslmode=disable"
	}
	
	var err error
	db, err = sql.Open("postgres", connStr)

	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}

	app = web.Init()
}

func main() {
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
