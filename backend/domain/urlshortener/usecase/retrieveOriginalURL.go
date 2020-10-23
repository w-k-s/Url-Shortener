package usecase

import (
	"fmt"
	"github.com/w-k-s/short-url/domain"
	u "github.com/w-k-s/short-url/domain/urlshortener"
	"net/url"
)

type RetrieveOriginalURLUseCase struct {
	repo u.URLRepository
}

func NewRetrieveOriginalURLUseCase(repo u.URLRepository) *RetrieveOriginalURLUseCase {
	return &RetrieveOriginalURLUseCase{
		repo,
	}
}

func (s *RetrieveOriginalURLUseCase) Execute(retrieveRequest RetrieveOriginalURLRequest) (RetrieveOriginalURLResponse, domain.Err) {

	var shortID string
	path := retrieveRequest.ShortURL().Path
	if len(path) == 0 {
		return RetrieveOriginalURLResponse{}, NewError(
			RetrieveFullURLValidation,
			fmt.Sprintf("The URL '%s' does not have a path.", retrieveRequest.ShortURL().String()),
			nil,
		)
	}

	if path[0] == '/' {
		shortID = path[1:]
	}

	record, err := s.repo.LongURL(shortID)
	if err != nil {
		return RetrieveOriginalURLResponse{}, NewError(
			RetrieveFullURLNotFound,
			fmt.Sprintf("No URL for %s", shortID),
			map[string]string{"error": err.Error()},
		)
	}

	longURL, err := url.Parse(record.LongURL)
	if err != nil {
		return RetrieveOriginalURLResponse{}, NewError(
			RetrieveFullURLParsing,
			fmt.Sprintf("Failed to parse %s", record.LongURL),
			map[string]string{"error": err.Error()},
		)
	}

	return RetrieveOriginalURLResponse{
		LongURL:  longURL.String(),
		ShortURL: retrieveRequest.ShortURL().String(),
	}, nil
}
