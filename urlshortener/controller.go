package urlshortener

import (
	"encoding/json"
	"fmt"
	app "github.com/w-k-s/short-url/app"
	err "github.com/w-k-s/short-url/error"
	"net/http"
	"net/url"
)

type Controller struct {
	service *Service
}

func NewController(app *app.App) *Controller {
	return &Controller{
		NewService(
			NewURLRepository(app.Db),
			app.Logger,
		),
	}
}

//--URLResponse

type urlResponse struct {
	LongUrl  string `json:"longUrl"`
	ShortUrl string `json:"shortUrl"`
}

//--Shorten URL

type shortenUrlRequest struct {
	LongUrl string `json:"longUrl"`
}

func parseShortenUrlRequest(req *http.Request) (*url.URL, err.Err) {

	var shortenReq shortenUrlRequest
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&shortenReq)
	if err != nil {
		return nil, NewError(
			ShortenURLDecoding,
			"JSON Body must include `longUrl`",
			map[string]string{"error": err.Error()},
		)
	}

	rawUrl, err := url.Parse(shortenReq.LongUrl)
	if err != nil {
		return nil, NewError(
			ShortenURLValidation,
			fmt.Sprintf("'%s' is not a valid url", shortenReq.LongUrl),
			map[string]string{"error": err.Error()},
		)
	}

	if !rawUrl.IsAbs() {
		return nil, NewError(
			ShortenURLValidation,
			fmt.Sprintf("'%s' is a relative url. Absolute urls are expected", shortenReq.LongUrl),
			nil,
		)
	}

	return rawUrl, nil
}

func (c *Controller) ShortenUrl(w http.ResponseWriter, req *http.Request) {

	scheme := req.URL.Scheme
	if len(scheme) == 0 {
		scheme = "https"
	}

	reqUrl := &url.URL{
		Scheme: scheme,
		Host:   req.Host,
	}

	longUrl, err := parseShortenUrlRequest(req)
	if err != nil {
		SendError(w, err)
		return
	}

	shortUrl, err := c.service.ShortenUrl(reqUrl, longUrl)
	if err != nil {
		SendError(w, err)
		return
	}

	encoder := json.NewEncoder(w)
	_err_ := encoder.Encode(&urlResponse{
		LongUrl:  longUrl.String(),
		ShortUrl: shortUrl.String(),
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

	shortUrlReq := req.FormValue("shortUrl")

	if len(shortUrlReq) == 0 {
		return nil, NewError(
			RetrieveFullURLValidation,
			"`shortUrl` is required",
			nil,
		)
	}

	shortUrl, err := url.Parse(shortUrlReq)
	if err != nil {
		return nil, NewError(
			RetrieveFullURLValidation,
			fmt.Sprintf("'%s' is not a valid url", shortUrl),
			map[string]string{"error": err.Error()},
		)
	}

	if !shortUrl.IsAbs() {
		return nil, NewError(
			RetrieveFullURLValidation,
			fmt.Sprintf("'%s' is a relative url. Absolute urls are expected", shortUrl),
			nil,
		)
	}

	return shortUrl, nil
}

func (c *Controller) GetLongUrl(w http.ResponseWriter, req *http.Request) {

	shortUrl, err := parseRetrieveFullURLRequest(req)
	if err != nil {
		SendError(w, err)
		return
	}

	longUrl, err := c.service.GetLongUrl(shortUrl)
	if err != nil {
		SendError(w, err)
		return
	}

	encoder := json.NewEncoder(w)
	_err_ := encoder.Encode(&urlResponse{
		LongUrl:  longUrl.String(),
		ShortUrl: shortUrl.String(),
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

func (c *Controller) RedirectToLongUrl(w http.ResponseWriter, req *http.Request) {

	longUrl, err := c.service.GetLongUrl(req.URL)

	if err != nil {
		SendError(w, err)
		return
	}

	http.Redirect(w, req, longUrl.String(), http.StatusSeeOther)
}
