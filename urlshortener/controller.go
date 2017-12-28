package urlshortener

import (
	"encoding/json"
	"io"
	"net/http"
	"fmt"
)

type Controller struct {
}

func NewController() *Controller {
	return &Controller{}
}

func (c *Controller) ShortenUrl(w http.ResponseWriter, req *http.Request) {

	var shortenReq ShortenUrlRequest
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&shortenReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	io.WriteString(w, fmt.Sprintf("The url is %s", shortenReq.LongUrl))
}
