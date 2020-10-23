package web

import (
	"github.com/gorilla/mux"
	"github.com/w-k-s/short-url/log"
	"net/http"
	"time"
)

type Routable interface {
	Route(*mux.Router)
}

type App struct {
	server *http.Server
	router *mux.Router
}

func Init(listenAddress string) *App {

	router := mux.NewRouter()

	server := createServer(router, listenAddress)

	app := &App{
		server,
		router,
	}

	return app
}

func (a *App) ListenAndServe() error {
	log.Printf("Listening on address: %s", a.server.Addr)
	return a.server.ListenAndServe()
}

func (a *App) Register(routable Routable) {
	routable.Route(a.router)
}

func createServer(h http.Handler, address string) *http.Server {
	return &http.Server{
		Handler:      h,
		Addr:         address,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
}
