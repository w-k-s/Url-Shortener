package web

import (
	"encoding/json"
	"fmt"
	err "github.com/w-k-s/short-url/error"
	"github.com/w-k-s/short-url/domain/urlshortener/usecase"
	"net/http"
	"net/url"
)

type Controller struct {
	shortenURLUseCase *useCase.shortenURLUseCase
}

func NewController(service *Service) *Controller {
	return &Controller{
		shortenURLUseCase,
	}
}

//--URLResponse


func sendResponse(w http.ResponseWriter, int status, body interface{}) {

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(status)
	encoder := json.NewEncoder(w)
	err := encoder.Encode(body)

	if err != nil {
		SendError(w, NewError(
			domain.URLResponseEncoding,
			"Error encoding response",
			map[string]string{"error": err.Error()},
		))
		return
	}
}

//--Shorten URL

func (c *Controller) ShortenURL(w http.ResponseWriter, req *http.Request) {

	shortenRequest, appErr := usecase.NewShortenURLRequest(req)
	if appErr != nil {
		SendError(w, appErr)
		return
	}

	shortenResponse, appErr := c.shortenURLUseCase.Execute(shortenRequest)
	if appErr != nil {
		SendError(w, appErr)
		return
	}

	sendURLResponse(w, http.StatusOK, shortenResponse)
}

//--Shorten URL

func parseRetrieveFullURLRequest(req *http.Request) (*url.URL, err.Err) {

	shortURLReq := req.FormValue("shortUrl")

	if len(shortURLReq) == 0 {
		return nil, NewError(
			RetrieveFullURLValidation,
			"`shortUrl` is required",
			nil,
		)
	}

	shortURL, err := url.Parse(shortURLReq)
	if err != nil {
		return nil, NewError(
			RetrieveFullURLValidation,
			fmt.Sprintf("'%s' is not a valid url", shortURL),
			map[string]string{"error": err.Error()},
		)
	}

	if !shortURL.IsAbs() {
		return nil, NewError(
			RetrieveFullURLValidation,
			fmt.Sprintf("'%s' is a relative url. Absolute urls are expected", shortURL),
			nil,
		)
	}

	return shortURL, nil
}

func (c *Controller) GetLongURL(w http.ResponseWriter, req *http.Request) {

	shortURL, err := parseRetrieveFullURLRequest(req)
	if err != nil {
		SendError(w, err)
		return
	}

	longURL, err := c.service.GetLongURL(shortURL)
	if err != nil {
		SendError(w, err)
		return
	}

	sendURLResponse(w, req, &urlResponse{
		LongURL:  longURL.String(),
		ShortURL: shortURL.String(),
	})
}

//--Redirect

func (c *Controller) RedirectToLongURL(w http.ResponseWriter, req *http.Request) {

	longURL, err := c.service.GetLongURL(req.URL)

	fmt.Printf("redirecting to %s\n", longURL)

	if err != nil {
		SendError(w, err)
		return
	}

	http.Redirect(w, req, longURL.String(), http.StatusSeeOther)
}
