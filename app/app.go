package app

import (
	"gopkg.in/mgo.v2"
	"html/template"
	"log"
	"os"
)

type App struct {
	Templates *template.Template
	Session   *mgo.Session
	Logger    *log.Logger
}

func Init(dataSourceName string) *App {
	tpl := template.Must(template.ParseGlob("./templates/*"))

	session, err := mgo.Dial(dataSourceName)
	if err != nil {
		panic(err)
	}

	logger := log.New(os.Stdout, "short-url: ", log.Ldate|log.Ltime)

	app := &App{
		tpl,
		session,
		logger,
	}

	app.ensureIndexes()

	return app
}
