package app

import (
	"github.com/w-k-s/short-url/db"
	"log"
)

type App struct {
	logger     *log.Logger
	db         *db.Db
	production bool
}

func Init(logger *log.Logger, connString string, production bool) *App {
	db := db.New(connString)

	app := &App{
		logger,
		db,
		production,
	}

	return app
}

func (a *App) Db() *db.Db {
	return a.db
}

func (a *App) Logger() *log.Logger {
	return a.logger
}

func (a *App) IsProd() bool {
	return a.production
}

func (a *App) Close() {
	a.db.Close()
}
