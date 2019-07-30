package app

import (
	"github.com/w-k-s/short-url/db"
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

type App struct {
	logger     *log.Logger
	db         *db.Db
	server 	   *http.Server
	router     *mux.Router
	production bool
	listeningAndServing bool
}

func Init() *App {
	production := os.Getenv("PROD") == "1"

	dbConnString := os.Getenv("MONGO_ADDRESS")
	if len(dbConnString) == 0 {
		dbConnString = "mongodb://localhost:27017/shorturl"
	}

	db := db.New(connString)

	address := os.Getenv("ADDRESS")
	if len(address) == 0 {
		address = ":8080"
	}
	
	router := mux.NewRouter()

	server := createServer(, address)

	logger := log.New(os.Stdout, "short-url: ", log.Ldate|log.Ltime)


	app := &App{
		logger,
		db,
		server,
		router,
		production,
		false,
	}

	log.Printf("Address: '%s'", address)
	log.Printf("Connection String: %s", dbConnString)
	log.Printf("Production: %v", production)
	log.Printf("Init Complete. Running on %s", address)

	return app
}

func (a *App) ListenAndServe(errchan chan error){
	go func(c chan error) {
		
		a.listeningAndServing = true
		err := a.server.ListenAndServe()
		a.listeningAndServing = false
		
		if err != nil {
			errchan <- err
		}
	}(errchan)
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

func createServer(h http.Handler, address string) *http.Server {
	return &http.Server{
		Handler: h,
		Addr:    address,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
}
