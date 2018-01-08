package urlshortener

import (
	"encoding/json"
	a "github.com/w-k-s/short-url/app"
	"net/http"
	"net/url"
)

type Controller struct {
	service *Service
}

func NewController(app *a.App) *Controller {
	return &Controller{
		NewService(app),
	}
}

func (c *Controller) ShortenUrl(w http.ResponseWriter, req *http.Request) {

	var shortenReq ShortenUrlRequest
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&shortenReq)
	if err != nil {
		a.EncodeNewErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	longUrl, err := shortenReq.Validate()
	if err != nil {
		a.EncodeNewErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortUrl, err := c.service.ShortenUrl(req.Host, longUrl)
	if err != nil {
		a.EncodeNewErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	encoder := json.NewEncoder(w)
	err = encoder.Encode(&UrlResponse{
		LongUrl:  longUrl.String(),
		ShortUrl: shortUrl.String(),
	})

	if err != nil {
		a.EncodeNewErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *Controller) GetLongUrl(w http.ResponseWriter, req *http.Request) {

	shortUrlReq := req.FormValue("shortUrl")

	if len(shortUrlReq) == 0 {
		a.EncodeNewErrorJSON(w, "Missing parameter: shortUrl", http.StatusBadRequest)
		return
	}

	shortUrl, err := url.Parse(shortUrlReq)
	if err != nil || !shortUrl.IsAbs() {
		a.EncodeNewErrorJSON(w, "url invalid or not absolute", http.StatusBadRequest)
		return
	}

	longUrl, found, err := c.service.GetLongUrl(shortUrl)
	if err != nil {
		a.EncodeNewErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !found {
		a.EncodeNewErrorJSON(w, "Not Found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	encoder := json.NewEncoder(w)
	err = encoder.Encode(&UrlResponse{
		LongUrl:  longUrl.String(),
		ShortUrl: shortUrl.String(),
	})

	if err != nil {
		a.EncodeNewErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *Controller) RedirectToLongUrl(w http.ResponseWriter, req *http.Request) {

	longUrl, found, err := c.service.GetLongUrl(req.URL)

	if found {
		http.Redirect(w, req, longUrl.String(), http.StatusSeeOther)
		return
	} else {
		a.EncodeNewErrorJSON(w, "Not Found", http.StatusNotFound)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
