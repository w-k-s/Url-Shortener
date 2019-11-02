package usecase

import (
	"fmt"
	"github.com/w-k-s/short-url/domain"
	"net/http"
	"net/url"
)

type RetrieveOriginalURLRequest struct {
	shortURL *url.URL
}

func RedirectShortURLRequest(shortURL *url.URL) RetrieveOriginalURLRequest {
	return RetrieveOriginalURLRequest{
		shortURL,
	}
}

func NewRetrieveOriginalURLRequest(req *http.Request) (RetrieveOriginalURLRequest, domain.Err) {

	shortURLReq := req.FormValue("shortUrl")

	if len(shortURLReq) == 0 {
		return RetrieveOriginalURLRequest{}, NewError(
			RetrieveFullURLValidation,
			"`shortUrl` is required",
			nil,
		)
	}

	shortURL, err := url.Parse(shortURLReq)
	if err != nil {
		return RetrieveOriginalURLRequest{}, NewError(
			RetrieveFullURLValidation,
			fmt.Sprintf("'%s' is not a valid url", shortURL),
			map[string]string{"error": err.Error()},
		)
	}

	if !shortURL.IsAbs() {
		return RetrieveOriginalURLRequest{}, NewError(
			RetrieveFullURLValidation,
			fmt.Sprintf("'%s' is a relative url. Absolute urls are expected", shortURL),
			nil,
		)
	}

	return RetrieveOriginalURLRequest{shortURL}, nil
}

func (r RetrieveOriginalURLRequest) ShortURL() *url.URL {
	return r.shortURL
}
