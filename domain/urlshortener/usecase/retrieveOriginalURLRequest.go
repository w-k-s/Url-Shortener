package usecase

import(
	"net/http"
	"net/url"
	"github.com/w-k-s/short-url/domain"
)

type RetrieveOriginalURLRequest struct {
	shortURL   *url.URL
}

func NewRetrieveOriginalURLRequest(req *http.Request) (RetrieveOriginalURLRequest, domain.Err) {

	shortURLReq := req.FormValue("shortUrl")

	if len(shortURLReq) == 0 {
		return RetrieveOriginalURLRequest{}, NewError(
			domain.RetrieveFullURLValidation,
			"`shortUrl` is required",
			nil,
		)
	}

	shortURL, err := url.Parse(shortURLReq)
	if err != nil {
		return RetrieveOriginalURLRequest{}, NewError(
			domain.RetrieveFullURLValidation,
			fmt.Sprintf("'%s' is not a valid url", shortURL),
			map[string]string{"error": err.Error()},
		)
	}

	if !shortURL.IsAbs() {
		return RetrieveOriginalURLRequest{}, NewError(
			domain.RetrieveFullURLValidation,
			fmt.Sprintf("'%s' is a relative url. Absolute urls are expected", shortURL),
			nil,
		)
	}

	return RetrieveOriginalURLRequest{shortURL}, nil
}