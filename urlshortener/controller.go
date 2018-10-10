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

type urlResponse struct {
	LongURL  string `json:"longUrl"`
	ShortURL string `json:"shortUrl"`
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

	longURL, err := parseShortenURLRequest(req)
	if err != nil {
		SendError(w, err)
		return
	}

	shortURL, err := c.service.ShortenURL(reqURL, longURL)
	if err != nil {
		SendError(w, err)
		return
	}

	encoder := json.NewEncoder(w)
	_err_ := encoder.Encode(&urlResponse{
		LongURL:  longURL.String(),
		ShortURL: shortURL.String(),
	})

	if _err_ != nil {
		SendError(w, NewError(
			ShortenURLEncoding,
			"Error encoding response",
			map[string]string{"error": _err_.Error()},
		))
		return
	}
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

	encoder := json.NewEncoder(w)
	_err_ := encoder.Encode(&urlResponse{
		LongURL:  longURL.String(),
		ShortURL: shortURL.String(),
	})

	if _err_ != nil {
		SendError(w, NewError(
			RetrieveFullURLEncoding,
			"Error encoding response",
			map[string]string{"error": _err_.Error()},
		))
		return
	}
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
