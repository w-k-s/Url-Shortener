package app

import (
	"github.com/w-k-s/short-url/db"
	"html/template"
	"log"
	"os"
)

type App struct {
	Templates *template.Template
	Logger    *log.Logger
	Db        *db.Db
}

func Init(host string, dbName string) *App {
	logger := log.New(os.Stdout, "short-url: ", log.Ldate|log.Ltime)

	tpl := template.Must(template.ParseGlob("./templates/*"))

	db := db.New(host, dbName)

	app := &App{
		tpl,
		logger,
		db,
	}

	return app
}
