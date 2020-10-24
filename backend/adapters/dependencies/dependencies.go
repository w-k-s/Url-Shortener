package dependencies

import (
	"database/sql"
	_ "github.com/lib/pq"
	persistence "github.com/w-k-s/short-url/adapters/db"
	"github.com/w-k-s/short-url/adapters/logging"
	"github.com/w-k-s/short-url/adapters/web"
	"github.com/w-k-s/short-url/config"
	"github.com/w-k-s/short-url/domain/urlshortener"
	"github.com/w-k-s/short-url/domain/urlshortener/usecase"
	"log"
	"net/url"
)

var Db *sql.DB
var urlRepo urlshortener.URLRepository
var baseURL *url.URL
var ShortenURLUseCase *usecase.ShortenURLUseCase
var RetrieveOriginalURLUseCase *usecase.RetrieveOriginalURLUseCase
var LogRepository *logging.LogRepository
var JsonFmt web.JsonFmt

func Init() {
	initDB()
	initURLRepository()
	initShortenURLUseCase()
	initRetrieveOriginalUseCase()
	initLogRepository()
	initJsonFmt()
}

func initDB() {
	var err error
	Db, err = sql.Open("postgres", config.Settings.DatabaseConnectionString)

	if err != nil && Db.Ping() != nil {
		log.Fatalf("Failed to ping db with connection string %q: %s", config.Settings.DatabaseConnectionString, err)
	}
}

func initURLRepository() {
	urlRepo = persistence.NewURLRepository(Db)
}

func initShortenURLUseCase() {
	ShortenURLUseCase = usecase.NewShortenURLUseCase(urlRepo, config.Settings.GetBaseURL(), usecase.DefaultShortIDGenerator{})
}

func initRetrieveOriginalUseCase() {
	RetrieveOriginalURLUseCase = usecase.NewRetrieveOriginalURLUseCase(urlRepo)
}

func initLogRepository() {
	LogRepository = logging.NewLogRepository(Db)
}

func initJsonFmt() {
	JsonFmt = web.NewJsonFmtWithHeaders(map[string]string{
		"Access-Control-Allow-Origin": config.Settings.AccessControlAllowOriginHeader,
	})
}
