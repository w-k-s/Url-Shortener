package usecase

import(
	"net/http"
	"net/url"
	"github.com/w-k-s/short-url/domain"
)

type ShortenURLRequest struct {
	longURL   string `json:"longUrl"`
	shortID   string `json:"ShortId"`
	parsedURL *url.URL
	requestURL *url.URL
}

func NewShortenURLRequest(req *http.Request) (ShortenURLRequest, domain.Err) {
	scheme := req.URL.Scheme
	if len(scheme) == 0 {
		scheme = "https"
	}

	requestURL := &url.URL{
		Scheme: scheme,
		Host:   req.Host,
	}

	decoder := json.NewDecoder(req.Body)

	var shortenReq ShortenURLRequest
	err := decoder.Decode(&shortenReq)
	if err != nil {
		return ShortenURLRequest{}, domain.NewError(
			ShortenURLDecoding,
			"JSON Body must include `longUrl`",
			map[string]string{"error": err.Error()},
		)
	}

	rawURL, err := url.Parse(shortenReq.LongURL)
	if err != nil {
		return ShortenURLRequest{}, domain.NewError(
			ShortenURLValidation,
			fmt.Sprintf("'%s' is not a valid url", shortenReq.LongURL),
			map[string]string{"error": err.Error()},
		)
	}

	if !rawURL.IsAbs() {
		return ShortenURLRequest{}, domain.NewError(
			ShortenURLValidation,
			fmt.Sprintf("'%s' is a relative url. Absolute urls are expected", shortenReq.LongURL),
			nil,
		)
	}

	return ShortenURLRequest{
		longURL:   shortenReq.LongURL,
		shortID:   shortenReq.ShortID,
		parsedURL: rawURL,
		requestURL: requestURL,
	}, nil
}

func (s ShortenURLRequest) UserDidSpecifyShortId() bool {
	return len(s.ShortID) > 0
}

func (s ShortenURLRequest) LongURL() string {
	return s.longURL
}

func (s ShortenURLRequest) ShortID() string{
	return s.shortID
}

func (s ShortenURLRequest) ParsedURL() *url.URL{
	return s.parsedURL
}

func (a ShortenURLRequest) RequestURL() *url.URL{
	return s.requestURL
}