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

func (s *Service) ShortenUrl(reqUrl *url.URL, longUrl *url.URL) (*url.URL, err.Err) {

	existingRecord, _ := s.repo.ShortURL(longUrl.String())

	if existingRecord != nil {
		s.logger.Printf("Record found. Long Url: %s, shortUrl: %s", longUrl, existingRecord.ShortId)
		return buildShortenedUrl(reqUrl, existingRecord), nil
	}

	deviations := []Deviation{VERY_SHORT, SHORT, MEDIUM, VERY_LONG}
	inserted := false
	var newRecord *URLRecord
	var err error

	for try := 0; !inserted && try < len(deviations); try++ {
		shortId := s.generator.Generate(deviations[try])
		newRecord, err = s.repo.SaveRecord(&URLRecord{
			LongUrl:    longUrl.String(),
			ShortId:    shortId,
			CreateTime: time.Now(),
		})

		s.logger.Printf("longUrl '%s' (Attempt %d): Using shortId '%s'. Error: %s", longUrl, try, shortId, err)
		inserted = err == nil
	}

	if !inserted {
		return nil, NewError(
			ShortenURLFailedToSave,
			fmt.Sprintf("Failed to find a shortId after %d attempts", len(deviations)),
			map[string]string{"error": err.Error()},
		)
	}

	return buildShortenedUrl(reqUrl, newRecord), nil
}

func buildShortenedUrl(reqUrl *url.URL, urlRecord *URLRecord) *url.URL {
	return &url.URL{
		Scheme: reqUrl.Scheme,
		Host:   reqUrl.Host,
		Path:   urlRecord.ShortId,
	}
}

func (s *Service) GetLongUrl(shortUrl *url.URL) (*url.URL, err.Err) {

	var shortId string
	path := shortUrl.Path
	if len(path) == 0 {
		return nil, NewError(
			RetrieveFullURLValidation,
			fmt.Sprintf("The URL '%s' does not have a path.", shortUrl),
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

	longUrl, err := url.Parse(record.LongUrl)
	if err != nil {
		return nil, NewError(
			RetrieveFullURLParsing,
			fmt.Sprintf("Failed to parse %s", record.LongUrl),
			map[string]string{"error": err.Error()},
		)
	}

	return longUrl, nil
}
