package urlshortener

import (
	"encoding/json"
	"fmt"
	err "github.com/w-k-s/short-url/error"
	"net/http"
	"net/url"
)

type Controller struct {
	service *Service
}

func NewController(service *Service) *Controller {
	return &Controller{
		service,
	}
}

//--URLResponse

const urlResponseCacheControlMaxAge = 172800 // 2 days

type urlResponse struct {
	LongURL  string `json:"longUrl"`
	ShortURL string `json:"shortUrl"`
}

func sendURLResponse(w http.ResponseWriter, req *http.Request, urlResponse *urlResponse) {

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	err := encoder.Encode(urlResponse)

	if err != nil {
		SendError(w, NewError(
			URLResponseEncoding,
			"Error encoding response",
			map[string]string{"error": err.Error()},
		))
		return
	}
}

//--Shorten URL

type shortenURLRequest struct {
	LongURL string `json:"longUrl"`
}

func parseShortenURLRequest(req *http.Request) (*url.URL, err.Err) {

	var shortenReq shortenURLRequest
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&shortenReq)
	if err != nil {
		return nil, NewError(
			ShortenURLDecoding,
			"JSON Body must include `longUrl`",
			map[string]string{"error": err.Error()},
		)
	}

	rawURL, err := url.Parse(shortenReq.LongURL)
	if err != nil {
		return nil, NewError(
			ShortenURLValidation,
			fmt.Sprintf("'%s' is not a valid url", shortenReq.LongURL),
			map[string]string{"error": err.Error()},
		)
	}

	if !rawURL.IsAbs() {
		return nil, NewError(
			ShortenURLValidation,
			fmt.Sprintf("'%s' is a relative url. Absolute urls are expected", shortenReq.LongURL),
			nil,
		)
	}

	return rawURL, nil
}

func (c *Controller) ShortenURL(w http.ResponseWriter, req *http.Request) {

	scheme := req.URL.Scheme
	if len(scheme) == 0 {
		scheme = "https"
	}

	reqURL := &url.URL{
		Scheme: scheme,
		Host:   req.Host,
	}

	longURL, appErr := parseShortenURLRequest(req)
	if appErr != nil {
		SendError(w, appErr)
		return
	}

	shortURL, appErr := c.service.ShortenURL(reqURL, longURL)
	if appErr != nil {
		SendError(w, appErr)
		return
	}

	sendURLResponse(w, req, &urlResponse{
		LongURL:  longURL.String(),
		ShortURL: shortURL.String(),
	})
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
