package urlshortener

import (
	"fmt"
	err "github.com/w-k-s/short-url/error"
	"log"
	"net/url"
	"time"
)

type Service struct {
	repo      *URLRepository
	logger    *log.Logger
	generator ShortIdGenerator
}

func NewService(repo *URLRepository, logger *log.Logger, generator ShortIdGenerator) *Service {
	return &Service{
		repo,
		logger,
		generator,
	}
}

func (s *Service) ShortenURL(reqURL *url.URL, longURL *url.URL) (*url.URL, err.Err) {

	existingRecord, _ := s.repo.ShortURL(longURL.String())

	if existingRecord != nil {
		s.logger.Printf("Record found. Long Url: %s, shortURL: %s", longURL, existingRecord.ShortId)
		return buildShortenedURL(reqURL, existingRecord), nil
	}

	deviations := []Deviation{VERY_SHORT, SHORT, MEDIUM, VERY_LONG}
	inserted := false
	var newRecord *URLRecord
	var err error

	for try := 0; !inserted && try < len(deviations); try++ {
		shortId := s.generator.Generate(deviations[try])
		newRecord, err = s.repo.SaveRecord(&URLRecord{
			LongURL:    longURL.String(),
			ShortId:    shortId,
			CreateTime: time.Now(),
		})

		s.logger.Printf("longURL '%s' (Attempt %d): Using shortId '%s'.\n\t-- Error: %s\n\n", longURL, try, shortId, err)
		inserted = err == nil
	}

	if !inserted {
		return nil, NewError(
			ShortenURLFailedToSave,
			fmt.Sprintf("Failed to find a shortId after %d attempts", len(deviations)),
			map[string]string{"error": err.Error()},
		)
	}

	return buildShortenedURL(reqURL, newRecord), nil
}

func buildShortenedURL(reqURL *url.URL, urlRecord *URLRecord) *url.URL {
	return &url.URL{
		Scheme: reqURL.Scheme,
		Host:   reqURL.Host,
		Path:   urlRecord.ShortId,
	}
}

func (s *Service) GetLongURL(shortURL *url.URL) (*url.URL, err.Err) {

	var shortId string
	path := shortURL.Path
	if len(path) == 0 {
		return nil, NewError(
			RetrieveFullURLValidation,
			fmt.Sprintf("The URL '%s' does not have a path.", shortURL),
			nil,
		)
	}

	if path[0] == '/' {
		shortId = path[1:]
	}

	record, err := s.repo.LongURL(shortId)
	if err != nil {
		return nil, NewError(
			RetrieveFullURLNotFound,
			fmt.Sprintf("No URL for %s", shortId),
			map[string]string{"error": err.Error()},
		)
	}

	longURL, err := url.Parse(record.LongURL)
	if err != nil {
		return nil, NewError(
			RetrieveFullURLParsing,
			fmt.Sprintf("Failed to parse %s", record.LongURL),
			map[string]string{"error": err.Error()},
		)
	}

	return longURL, nil
}
