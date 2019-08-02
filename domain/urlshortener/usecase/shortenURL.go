package usecase

import (
	u "github.com/w-k-s/short-url/domain/urlshortener"
	"github.com/w-k-s/short-url/domain"
)

type ShortenURLUseCase struct{
	repo *u.URLRepository
	logger    *log.Logger
	generator ShortIDGenerator
}

func (s *ShortenURLUseCase) Execute(shortReq ShortenURLRequest) (response ShortenURLResponse, domain.Err){
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
				domain.ShortenURLShortIdInUse,
				fmt.Sprintf("Can not save shortId '%s'; possibly in-use", shortReq.ShortID),
				map[string]string{"error": err.Error()},
			)
		}
		return buildShortenedURLResponse(reqURL, longURL, newRecord), nil
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
			domain.ShortenURLFailedToSave,
			fmt.Sprintf("Failed to find a shortId after %d attempts", len(shortIDLengths)),
			map[string]string{"error": err.Error()},
		)
	}

	return buildShortenedURLResponse(shortReq, newRecord), nil
}

func buildShortenedURLResponse(shortReq ShortenURLRequest, urlRecord *URLRecord) ShortenURLResponse {
	shortURL := &url.URL{
		Scheme: shortReq.RequestURL().Scheme,
		Host:   shortReq.RequestURL().Host,
		Path:   urlRecord.ShortID
	}
	return ShortenURLResponse{
		LongURL: shortReq.LongURL().String(),
		ShortURL: shortURL.String(),
	}
}