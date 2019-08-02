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
	generator ShortIDGenerator
}

func NewService(repo *URLRepository, logger *log.Logger, generator ShortIDGenerator) *Service {
	return &Service{
		repo,
		logger,
		generator,
	}
}

func (s *Service) ShortenURL(reqURL *url.URL, shortReq shortenURLRequest) (*url.URL, err.Err) {

	longURL := shortReq.parsedURL
	existingRecord, _ := s.repo.ShortURL(longURL.String())

	if existingRecord != nil {
		s.logger.Printf("Record found. Long Url: %s, shortURL: %s", longURL, existingRecord.ShortID)
		return buildShortenedURL(reqURL, existingRecord), nil
	}

	if shortReq.UserDidSpecifyShortId() {
		newRecord, err := s.repo.SaveRecord(&URLRecord{
			LongURL:    longURL.String(),
			ShortID:    shortReq.ShortID,
			CreateTime: time.Now(),
		})
		if err != nil {
			return nil, NewError(
				ShortenURLShortIdInUse,
				fmt.Sprintf("Can not save shortId '%s'; possibly in-use", shortReq.ShortID),
				map[string]string{"error": err.Error()},
			)
		}
		return buildShortenedURL(reqURL, newRecord), nil
	}

	shortIDLengths := []ShortIDLength{VERY_SHORT, SHORT, MEDIUM, VERY_LONG}
	inserted := false
	var newRecord *URLRecord
	var err error

	for try := 0; !inserted && try < len(shortIDLengths); try++ {
		shortID := s.generator.Generate(shortIDLengths[try])
		newRecord, err = s.repo.SaveRecord(&URLRecord{
			LongURL:    longURL.String(),
			ShortID:    shortID,
			CreateTime: time.Now(),
		})

		s.logger.Printf("longURL '%s' (Attempt %d): Using shortId '%s'.\n\t-- Error: %v\n\n", longURL, try, shortID, err)
		inserted = err == nil
	}

	if !inserted {
		return nil, NewError(
			ShortenURLFailedToSave,
			fmt.Sprintf("Failed to find a shortId after %d attempts", len(shortIDLengths)),
			map[string]string{"error": err.Error()},
		)
	}

	return buildShortenedURL(reqURL, newRecord), nil
}

func buildShortenedURL(reqURL *url.URL, urlRecord *URLRecord) *url.URL {
	return &url.URL{
		Scheme: reqURL.Scheme,
		Host:   reqURL.Host,
		Path:   urlRecord.ShortID,
	}
}

func (s *Service) GetLongURL(shortURL *url.URL) (*url.URL, err.Err) {

	var shortID string
	path := shortURL.Path
	if len(path) == 0 {
		return nil, NewError(
			RetrieveFullURLValidation,
			fmt.Sprintf("The URL '%s' does not have a path.", shortURL),
			nil,
		)
	}

	if path[0] == '/' {
		shortID = path[1:]
	}

	record, err := s.repo.LongURL(shortID)
	if err != nil {
		return nil, NewError(
			RetrieveFullURLNotFound,
			fmt.Sprintf("No URL for %s", shortID),
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
