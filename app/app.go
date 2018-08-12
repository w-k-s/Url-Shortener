package app

import (
	"github.com/w-k-s/short-url/db"
	"log"
	"os"
)

type App struct {
	Logger *log.Logger
	Db     *db.Db
}

func Init(connString string) *App {
	logger := log.New(os.Stdout, "short-url: ", log.Ldate|log.Ltime)

	db := db.New(connString)

	app := &App{
		logger,
		db,
	}

	return app
}
