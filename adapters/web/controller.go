package web

import (
	"encoding/json"
	"github.com/w-k-s/short-url/domain/urlshortener/usecase"
	"log"
	"net/http"
)

type Controller struct {
	shortenURLUseCase          *usecase.ShortenURLUseCase
	retrieveOriginalURLUseCase *usecase.RetrieveOriginalURLUseCase
	logger                     *log.Logger
}

func NewController(shortenURLUseCase *usecase.ShortenURLUseCase,
	retrieveOriginalURLUseCase *usecase.RetrieveOriginalURLUseCase,
	logger *log.Logger) *Controller {
	return &Controller{
		shortenURLUseCase,
		retrieveOriginalURLUseCase,
		logger,
	}
}

//--Shorten URL

func (c *Controller) ShortenURL(w http.ResponseWriter, req *http.Request) {

	shortenRequest, err := usecase.NewShortenURLRequest(req)
	if err != nil {
		SendError(w, err)
		return
	}

	shortenResponse, err := c.shortenURLUseCase.Execute(shortenRequest)
	if err != nil {
		SendError(w, err)
		return
	}

	sendResponse(w, http.StatusOK, shortenResponse)
}

//--Shorten URL

func (c *Controller) GetLongURL(w http.ResponseWriter, req *http.Request) {

	retrieveRequest, err := usecase.NewRetrieveOriginalURLRequest(req)
	if err != nil {
		SendError(w, err)
		return
	}

	retrieveResponse, err := c.retrieveOriginalURLUseCase.Execute(retrieveRequest)
	if err != nil {
		SendError(w, err)
		return
	}

	sendResponse(w, http.StatusOK, retrieveResponse)
}

//--Redirect

func (c *Controller) RedirectToLongURL(w http.ResponseWriter, req *http.Request) {

	redirectRequest := usecase.RedirectShortURLRequest(req.URL)

	redirectResponse, err := c.retrieveOriginalURLUseCase.Execute(redirectRequest)
	if err != nil {
		SendError(w, err)
		return
	}

	c.logger.Printf("redirecting to %s\n", redirectResponse.LongURL)
	http.Redirect(w, req, redirectResponse.LongURL, http.StatusSeeOther)
}

//--URLResponse

func sendResponse(w http.ResponseWriter, status int, body interface{}) {

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(status)
	encoder := json.NewEncoder(w)
	err := encoder.Encode(body)

	if err != nil {
		SendEncodingError(w, body, err)
		return
	}
}
