package main

import (
	dep "github.com/w-k-s/short-url/adapters/dependencies"
	"github.com/w-k-s/short-url/adapters/web"
	"github.com/w-k-s/short-url/adapters/web/controllers"
	"github.com/w-k-s/short-url/config"
	"github.com/w-k-s/short-url/log"
)

var app *web.App

func init() {
	config.Init()
	dep.Init()
	log.Init()
}

func main() {

	app = web.Init(config.Settings.ListenAddress)

	app.Register(controllers.GetHealthCheckHandler(dep.Db))
	app.Register(controllers.GetShortenURLHandler(dep.ShortenURLUseCase, dep.JsonFmt))
	app.Register(controllers.GetRetrieveOriginalURLHandler(dep.RetrieveOriginalURLUseCase, dep.JsonFmt))
	//app.Register(controllers.GetRedirectToOriginalURLHandler(dep.RetrieveOriginalURLUseCase, dep.JsonFmt))
	//app.Register(controllers.GetLogRequestMiddleware(dep.LogRepository))

	log.Panic(app.ListenAndServe())
}
