package usecase

import (
	"encoding/json"
	"fmt"
	"github.com/w-k-s/short-url/domain"
	"net/http"
	"net/url"
)

type ShortenURLRequest struct {
	LongURL   string `json:"longUrl"`
	ShortID   string `json:"ShortId"`
	parsedURL *url.URL
}

func NewShortenURLRequest(req *http.Request) (ShortenURLRequest, domain.Err) {

	decoder := json.NewDecoder(req.Body)

	var shortenReq ShortenURLRequest
	err := decoder.Decode(&shortenReq)
	if err != nil {
		return ShortenURLRequest{}, NewError(
			ShortenURLDecoding,
			"JSON Body must include `longUrl`",
			map[string]string{"error": err.Error()},
		)
	}

	rawURL, err := url.Parse(shortenReq.LongURL)
	if err != nil {
		return ShortenURLRequest{}, NewError(
			ShortenURLValidation,
			fmt.Sprintf("'%s' is not a valid url", shortenReq.LongURL),
			map[string]string{"error": err.Error()},
		)
	}

	if !rawURL.IsAbs() {
		return ShortenURLRequest{}, NewError(
			ShortenURLValidation,
			fmt.Sprintf("'%s' is a relative url. Absolute urls are expected", shortenReq.LongURL),
			nil,
		)
	}

	return ShortenURLRequest{
		LongURL:   shortenReq.LongURL,
		ShortID:   shortenReq.ShortID,
		parsedURL: rawURL,
	}, nil
}

func (s ShortenURLRequest) UserDidSpecifyShortId() bool {
	return len(s.ShortID) > 0
}

func (s ShortenURLRequest) ParsedURL() *url.URL {
	return s.parsedURL
}
