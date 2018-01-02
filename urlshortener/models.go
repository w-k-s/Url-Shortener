package urlshortener

import (
	"errors"
	"net/url"
)

type ShortenUrlRequest struct {
	LongUrl string `json:"longUrl"`
}

func (r ShortenUrlRequest) Validate() (*url.URL, error) {
	rawUrl, err := url.Parse(r.LongUrl)
	if err != nil {
		return nil, err
	}

	if !rawUrl.IsAbs() {
		return nil, errors.New("url must be absolute")
	}

	return rawUrl, nil
}

type UrlResponse struct {
	LongUrl  string `json:"longUrl"`
	ShortUrl string `json:"shortUrl"`
}
