package usecase

import (
	"fmt"
	"github.com/w-k-s/short-url/domain"
)

const (
	//Shortening URL
	ShortenURLDecoding        domain.Code = 10200
	ShortenURLValidation                  = 10300
	ShortenURLFailedToSave                = 10400
	ShortenURLTrackVisitError             = 10401
	ShortenURLShortIdInUse                = 10402
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

func domainString(e domain.Code) string {
	switch e {
	//Shortening URL
	case ShortenURLDecoding:
		return "shortenUrl.decoding"
	case ShortenURLValidation:
		return "shortenUrl.validation"
	case ShortenURLFailedToSave:
		return "shortenUrl.failedToSave"
	case ShortenURLShortIdInUse:
		return "shortenUrl.shortIdInUse"
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

func NewError(code domain.Code, message string, fields map[string]string) *domain.Error {
	return domain.NewError(
		code,
		domainString(code),
		message,
		fields,
	)
}
