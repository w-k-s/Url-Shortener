package urlshortener

import (
	"fmt"
	err "github.com/w-k-s/short-url/error"
	"net/http"
)

const (
	//Shortening URL
	ShortenURLDecoding           err.Code = 10200
	ShortenURLValidation                  = 10300
	ShortenURLFailedToSave                = 10400
	ShortenURLSaveVisitTimeError          = 10401
	ShortenURLUndocumented                = 10999

	//Retrieving Long Url
	RetrieveFullURLDecoding     = 11200
	RetrieveFullURLValidation   = 11300
	RetrieveFullURLNotFound     = 11400
	RetrieveFullURLParsing      = 11500
	RetrieveFullURLUndocumented = 11999

	//Redirectign to Long Url
	RedirectionFullURLNotFound = 12100
	RedirectionUndocumented    = 12999

	//URLResponse
	URLResponseEncoding = 13000
)

func domain(e err.Code) string {
	switch e {
	//Shortening URL
	case ShortenURLDecoding:
		return "shortenUrl.decoding"
	case ShortenURLValidation:
		return "shortenUrl.validation"
	case ShortenURLFailedToSave:
		return "shortenUrl.failedToSave"
	case ShortenURLUndocumented:
		return "shortenUrl.undocumented"

	//Retrieving Long Url
	case RetrieveFullURLDecoding:
		return "retrieveFullURL.decoding"
	case RetrieveFullURLValidation:
		return "retrieveFullURL.validation"
	case RetrieveFullURLNotFound:
		return "retrieveFullURL.urlNotFound"
	case RetrieveFullURLParsing:
		return "retrieveFullURL.urlParsing"
	case RetrieveFullURLUndocumented:
		return "retrieveFullURL.undocumented"

	//Redirectign to Long Url
	case RedirectionFullURLNotFound:
		return "redirection.urlNotFound"
	case RedirectionUndocumented:
		return "redirection.undocumented"

	//Encoding URLResponse
	case URLResponseEncoding:
		return "urlResponse.encoding"

	default:
		return fmt.Sprintf("Unknown Domain (%d)", e)
	}
}

func NewError(code err.Code, message string, fields map[string]string) *err.Error {
	return err.NewError(
		code,
		domain(code),
		message,
		fields,
	)
}

func SendError(w http.ResponseWriter, e err.Err) {
	err.SendError(w, httpStatusCode(e.Code()), e)
}

func httpStatusCode(e err.Code) int {
	switch e {
	case ShortenURLValidation:
		fallthrough
	case RetrieveFullURLValidation:
		return http.StatusBadRequest
	case RetrieveFullURLNotFound:
		fallthrough
	case RedirectionFullURLNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
