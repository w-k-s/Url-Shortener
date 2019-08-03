package web

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"time"
)

type MiddlewareFunc mux.MiddlewareFunc

type App struct {
	logger     *log.Logger
	server     *http.Server
	router     *mux.Router
	production bool
}

func Init() *App {
	production := os.Getenv("PROD") == "1"

	address := os.Getenv("ADDRESS")
	if len(address) == 0 {
		address = ":8080"
	}

	router := mux.NewRouter()

	server := createServer(router, address)

	logger := log.New(os.Stdout, "", log.Llongfile|log.Ldate|log.Ltime|log.LUTC)

	app := &App{
		logger,
		server,
		router,
		production,
	}

	logger.Printf("Address: '%s'", address)
	logger.Printf("Production: %v", production)
	logger.Print("Init Complete.")

	return app
}

func (a *App) ListenAndServe() error {
	a.logger.Printf("Listening on address: %s", a.server.Addr)
	err := a.server.ListenAndServe()

	return err
}

func (a *App) Logger() *log.Logger {
	return a.logger
}

func (a *App) IsProd() bool {
	return a.production
}

func (a *App) Middleware(middlewareFunc MiddlewareFunc) {
	a.router.Use(mux.MiddlewareFunc(middlewareFunc))
}

func (a *App) Get(path string, f func(http.ResponseWriter, *http.Request)) {
	a.logRegisteredRoute("GET", path)
	a.router.HandleFunc(path, f).Methods("GET")
}

func (a *App) Post(path string, f func(http.ResponseWriter, *http.Request)) {
	a.logRegisteredRoute("POST", path)
	a.router.HandleFunc(path, f).Methods("POST")
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

func (a *App) logRegisteredRoute(method string, path string) {
	a.logger.Printf("Adding Route: '%s %s'", method, path)
}
