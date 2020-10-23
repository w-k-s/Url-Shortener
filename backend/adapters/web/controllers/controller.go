package controllers

import (
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/w-k-s/short-url/adapters/logging"
	"github.com/w-k-s/short-url/adapters/web"
	"github.com/w-k-s/short-url/domain/urlshortener/usecase"
	"github.com/w-k-s/short-url/log"
	"net/http"
)

// Shorten URL

type ShortenURLHandler http.HandlerFunc

func (h ShortenURLHandler) Route(r *mux.Router) {
	r.HandleFunc("/urlshortener/v1/url", h).
		Methods("POST")
}

func GetShortenURLHandler(useCase *usecase.ShortenURLUseCase, responseFmt web.ResponseFmt) ShortenURLHandler {
	return func(w http.ResponseWriter, req *http.Request) {
		shortenRequest, err := usecase.NewShortenURLRequest(req)
		if err != nil {
			responseFmt.Error(w, err)
			return
		}

		shortenResponse, err := useCase.Execute(shortenRequest)
		if err != nil {
			responseFmt.Error(w, err)
			return
		}

		responseFmt.Print(w, http.StatusOK, shortenResponse)
	}
}

// Get Original URL

type RetrieveOriginalURLHandler http.HandlerFunc

func (h RetrieveOriginalURLHandler) Route(r *mux.Router) {
	r.HandleFunc("/urlshortener/v1/url", h).
		Methods("GET")
}

func GetRetrieveOriginalURLHandler(useCase *usecase.RetrieveOriginalURLUseCase, responseFmt web.ResponseFmt) RetrieveOriginalURLHandler {
	return func(w http.ResponseWriter, req *http.Request) {
		retrieveRequest, err := usecase.NewRetrieveOriginalURLRequest(req)
		if err != nil {
			responseFmt.Error(w, err)
			return
		}

		retrieveResponse, err := useCase.Execute(retrieveRequest)
		if err != nil {
			responseFmt.Error(w, err)
			return
		}

		responseFmt.Print(w, http.StatusOK, retrieveResponse)
	}
}

//--Redirect

type RedirectToOriginalURLHandler http.HandlerFunc

func (h RedirectToOriginalURLHandler) Route(r *mux.Router) {
	r.HandleFunc("/{shortUrl}", h).
		Methods("GET")
}

func GetRedirectToOriginalURLHandler(useCase *usecase.RetrieveOriginalURLUseCase, responseFmt web.ResponseFmt) RedirectToOriginalURLHandler {
	return func(w http.ResponseWriter, req *http.Request) {
		redirectRequest := usecase.RedirectShortURLRequest(req.URL)

		redirectResponse, err := useCase.Execute(redirectRequest)
		if err != nil {
			responseFmt.Error(w, err)
			return
		}

		log.Printf("redirecting to %s\n", redirectResponse.LongURL)
		http.Redirect(w, req, redirectResponse.LongURL, http.StatusSeeOther)
	}
}

// Health Check

type HealthCheckHandler http.HandlerFunc

func (h HealthCheckHandler) Route(r *mux.Router) {
	r.HandleFunc("/health", h).
		Methods("GET")
}

func GetHealthCheckHandler(db *sql.DB) HealthCheckHandler {
	return func(w http.ResponseWriter, req *http.Request) {
		log.Printf("Do health check")

		if err := db.Ping(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

// Middleware

type LogRequestMiddleware mux.MiddlewareFunc

func (m LogRequestMiddleware) Route(r *mux.Router) {
	r.Use(mux.MiddlewareFunc(m))
}

func GetLogRequestMiddleware(logRepository *logging.LogRepository) LogRequestMiddleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			sw := &logging.StatusWriter{ResponseWriter: w}

			record := logRepository.LogRequest(r)

			next.ServeHTTP(sw, r)

			logRepository.LogResponse(sw, record)
		})
	}
}
