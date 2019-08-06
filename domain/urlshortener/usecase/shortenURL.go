package usecase

import (
	"fmt"
	"github.com/w-k-s/short-url/domain"
	u "github.com/w-k-s/short-url/domain/urlshortener"
	"log"
	"net/url"
	"time"
)

type ShortenURLUseCase struct {
	repo      u.URLRepository
	generator ShortIDGenerator
	logger    *log.Logger
}

func NewShortenURLUseCase(repo u.URLRepository, generator ShortIDGenerator, logger *log.Logger) *ShortenURLUseCase {
	return &ShortenURLUseCase{
		repo,
		generator,
		logger,
	}
}

func (s *ShortenURLUseCase) Execute(shortReq ShortenURLRequest) (ShortenURLResponse, domain.Err) {
	longURL := shortReq.parsedURL
	existingRecord, _ := s.repo.ShortURL(longURL.String())

	if existingRecord != nil {
		s.logger.Printf("Record found. Long Url: %s, shortURL: %s", longURL, existingRecord.ShortID)
		return buildShortenedURLResponse(shortReq, existingRecord), nil
	}

	if shortReq.UserDidSpecifyShortId() {
		newRecord, err := s.repo.SaveRecord(&u.URLRecord{
			LongURL:    longURL.String(),
			ShortID:    shortReq.ShortID,
			CreateTime: time.Now(),
		})
		if err != nil {
			return ShortenURLResponse{}, NewError(
				ShortenURLShortIDInUse,
				fmt.Sprintf("Can not save shortId '%s'; possibly in-use", shortReq.ShortID),
				map[string]string{"error": err.Error()},
			)
		}
		return buildShortenedURLResponse(shortReq, newRecord), nil
	}

	shortIDLengths := []ShortIDLength{VeryShort, Short, Medium, VeryLong}
	inserted := false
	var newRecord *u.URLRecord
	var err error

	for try := 0; !inserted && try < len(shortIDLengths); try++ {
		shortID := s.generator.Generate(shortIDLengths[try])
		newRecord, err = s.repo.SaveRecord(&u.URLRecord{
			LongURL:    longURL.String(),
			ShortID:    shortID,
			CreateTime: time.Now(),
		})

		s.logger.Printf("longURL '%s' (Attempt %d): Using shortId '%s'.\n\t-- Error: %v\n\n", longURL, try, shortID, err)
		inserted = err == nil
	}

	if !inserted {
		return ShortenURLResponse{}, NewError(
			ShortenURLFailedToSave,
			fmt.Sprintf("Failed to find a shortId after %d attempts", len(shortIDLengths)),
			map[string]string{"error": err.Error()},
		)
	}

	return buildShortenedURLResponse(shortReq, newRecord), nil
}

func buildShortenedURLResponse(shortReq ShortenURLRequest, urlRecord *u.URLRecord) ShortenURLResponse {
	shortURL := &url.URL{
		Scheme: shortReq.RequestURL().Scheme,
		Host:   shortReq.RequestURL().Host,
		Path:   urlRecord.ShortID,
	}
	return ShortenURLResponse{
		LongURL:  shortReq.LongURL,
		ShortURL: shortURL.String(),
	}
}
