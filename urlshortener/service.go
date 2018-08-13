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

func NewService(repo *URLRepository, logger *log.Logger) *Service {
	return &Service{
		repo,
		logger,
		NewShortIDGenerator(),
	}
}

func (s *Service) ShortenUrl(reqUrl *url.URL, longUrl *url.URL) (*url.URL, err.Err) {

	record, _ := s.repo.ShortURL(longUrl.String())

	if record != nil {
		s.logger.Printf("Record found. Long Url: %s, shortUrl: %s", longUrl, record.ShortId)
		return buildShortenedUrl(reqUrl, record), nil
	}

	deviations := []Deviation{VERY_SHORT, SHORT, MEDIUM, VERY_LONG}
	inserted := false
	var err error

	for try, deviation := range deviations {
		record, err = s.repo.SaveRecord(&URLRecord{
			LongUrl:    longUrl.String(),
			ShortId:    s.generator.Generate(deviation),
			CreateTime: time.Now(),
		})

		s.logger.Printf("longUrl '%s' (Attempt %d): Using shortId '%d'. Error: %s", longUrl, try, record.ShortId, err)
		inserted = err == nil
	}

	if !inserted {
		return nil, NewError(
			ShortenURLFailedToSave,
			fmt.Sprintf("Failed to find a shortId after %d attempts", len(deviations)),
			map[string]string{"error": err.Error()},
		)
	}

	return buildShortenedUrl(reqUrl, record), nil
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
