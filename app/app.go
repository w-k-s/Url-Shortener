package app

import (
	"github.com/w-k-s/short-url/db"
	"log"
)

type App struct {
	Logger *log.Logger
	Db     *db.Db
}

func Init(logger *log.Logger, connString string) *App {
	db := db.New(connString)

	app := &App{
		logger,
		db,
	}

	return app
}
