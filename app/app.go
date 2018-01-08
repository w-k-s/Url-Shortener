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

func Init(host string) *App {
	logger := log.New(os.Stdout, "short-url: ", log.Ldate|log.Ltime)

	tpl := template.Must(template.ParseGlob("./templates/*"))

	session, err := mgo.Dial(host)
	if err != nil {
		logger.Panicf("Could not connect to datastore with host %s - %v", host, err)
	}

	app := &App{
		tpl,
		session,
		logger,
	}

	app.ensureIndexes()

	return app
}
