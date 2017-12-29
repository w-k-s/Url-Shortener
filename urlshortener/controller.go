package urlshortener

import (
	"encoding/json"
	"fmt"
	a "github.com/waqqas-abdulkareem/short-url/app"
	"io"
	"net/http"
)

type Controller struct {
	service * Service
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	url, err := shortenReq.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	url, err = c.service.ShortenUrl(req.Host,url)
	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	io.WriteString(w, fmt.Sprintf("The url is %s", url))
}
