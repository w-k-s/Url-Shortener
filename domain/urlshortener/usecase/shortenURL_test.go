package usecase

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	u "github.com/w-k-s/short-url/domain/urlshortener"
	"log"
	"net/url"
	"os"
	"testing"
	"time"
)

//-- MockShortIDGenerator

type MockShortIDGenerator struct {
	ShortID string
}

func (m MockShortIDGenerator) Generate(d ShortIDLength) string {
	return m.ShortID
}

//-- MockURLRepository

type MockURLRepository struct {
	ReturnError bool

	SaveURLRecordResult *u.URLRecord
	SaveURLRecordError  error

	LongURLRecordResult *u.URLRecord
	LongURLRecordError  error

	ShortURLRecordResult *u.URLRecord
	ShortURLRecordError  error
}

func (m MockURLRepository) SaveRecord(record *u.URLRecord) (*u.URLRecord, error) {
	if m.ReturnError {
		return nil, m.SaveURLRecordError
	}
	return m.SaveURLRecordResult, nil
}

func (m MockURLRepository) LongURL(shortID string) (*u.URLRecord, error) {
	if m.ReturnError {
		return nil, m.LongURLRecordError
	}
	return m.LongURLRecordResult, nil
}

func (m MockURLRepository) ShortURL(longURL string) (*u.URLRecord, error) {
	if m.ReturnError {
		return nil, m.ShortURLRecordError
	}
	return m.ShortURLRecordResult, nil
}

//-----

const savedShortID = "shrt"
const savedLongURL = "http://www.example.com"
const savedShortURL = "http://small.ml/" + savedShortID

type ShortenURLUseCaseTestSuite struct {
	suite.Suite
	urlRepo   *MockURLRepository
	record    *u.URLRecord
	generator *MockShortIDGenerator
	useCase   *ShortenURLUseCase
}

func (suite *ShortenURLUseCaseTestSuite) SetupTest() {
	suite.record = &u.URLRecord{
		LongURL:    savedLongURL,
		ShortID:    savedShortID,
		CreateTime: time.Now(),
	}

	logger := log.New(os.Stdout, "short-url: ", log.Ldate|log.Ltime)
	suite.generator = &MockShortIDGenerator{}
	suite.urlRepo = &MockURLRepository{}
	suite.useCase = NewShortenURLUseCase(suite.urlRepo, suite.generator, logger)
}

func TestShortenURLUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(ShortenURLUseCaseTestSuite))
}

func (suite *ShortenURLUseCaseTestSuite) TestGivenRecordExists_WhenShorteningURL_ThenExistingRecordReturnedTestShortURLReturnedWhenRecordExists() {

	//Given
	hostURL, _ := url.Parse("http://www.small.ml")
	testURL, _ := url.Parse(savedLongURL)
	suite.urlRepo.ShortURLRecordResult = suite.record

	//When
	response, _ := suite.useCase.Execute(ShortenURLRequest{
		LongURL:    "http://www.1.ml",
		parsedURL:  testURL,
		requestURL: hostURL,
	})

	//Then
	expectation := "http://www.small.ml/" + savedShortID
	assert.Equal(suite.T(), expectation, response.ShortURL, "ShortenURL generates wrong url. Expected '%s'. Got: %s", expectation, response.ShortURL)
}

func (suite *ShortenURLUseCaseTestSuite) TestGivenRecordDoesNotExists_WhenShorteningURL_ThenRecordCreated() {

	//Given
	suite.generator.ShortID = "alpha"
	hostURL, _ := url.Parse("http://www.small.ml")
	testURL, _ := url.Parse("http://www.1.com")
	suite.urlRepo.SaveURLRecordResult = &u.URLRecord{
		LongURL:    savedLongURL,
		ShortID:    suite.generator.ShortID,
		CreateTime: time.Now(),
	}

	//When
	response, _ := suite.useCase.Execute(ShortenURLRequest{
		LongURL:    "http://www.1.ml",
		parsedURL:  testURL,
		requestURL: hostURL,
	})

	//Then
	expectation := "http://www.small.ml/" + suite.generator.ShortID
	assert.Equal(suite.T(), expectation, response.ShortURL, "ShortenURL generates wrong url. Expected '%s'. Got: %s", expectation, response.ShortURL)
}

func (suite *ShortenURLUseCaseTestSuite) TestGivenShortIDProvided_WhenShortIDNotInUse_ThenProvidedShortIDUsed() {

	//Given
	suite.generator.ShortID = "NotUsed"
	hostURL, _ := url.Parse("http://www.small.ml")
	testURL, _ := url.Parse("http://www.1.com")
	suite.urlRepo.SaveURLRecordResult = &u.URLRecord{
		LongURL:    savedLongURL,
		ShortID:    "Used",
		CreateTime: time.Now(),
	}

	//When
	response, _ := suite.useCase.Execute(ShortenURLRequest{
		LongURL:    "http://www.1.ml",
		ShortID:    "Used",
		parsedURL:  testURL,
		requestURL: hostURL,
	})

	//Then
	expectation := "http://www.small.ml/Used"
	assert.Equal(suite.T(), expectation, response.ShortURL, "ShortenURL generates wrong url. Expected '%s'. Got: %s", expectation, response.ShortURL)

}

func (suite *ShortenURLUseCaseTestSuite) TestGivenShortIDProvided_WhenShortIDInUse_ThenProvidedShortIDNotUsed() {

	//Given
	suite.generator.ShortID = "NotUsed"
	hostURL, _ := url.Parse("http://www.small.ml")
	testURL, _ := url.Parse("http://www.1.com")
	suite.urlRepo.ShortURLRecordResult = suite.record

	//When
	response, _ := suite.useCase.Execute(ShortenURLRequest{
		LongURL:    "http://www.1.ml",
		ShortID:    "NotUsed",
		parsedURL:  testURL,
		requestURL: hostURL,
	})

	//Then
	expectation := "http://www.small.ml/" + savedShortID
	assert.Equal(suite.T(), expectation, response.ShortURL, "ShortenURL generates wrong url. Expected '%s'. Got: %s", expectation, response.ShortURL)

}

func (suite *ShortenURLUseCaseTestSuite) TestGivenShortIDProvided_WhenShortIDNotUnique_ThenErrorReturned() {

	//Given
	suite.generator.ShortID = "NotUsed"
	hostURL, _ := url.Parse("http://www.small.ml")
	testURL, _ := url.Parse("http://www.1.com")
	suite.urlRepo.ReturnError = true
	suite.urlRepo.SaveURLRecordError = errors.New("short id exists")

	//When
	_, err := suite.useCase.Execute(ShortenURLRequest{
		LongURL:    "http://www.1.ml",
		ShortID:    savedShortID,
		parsedURL:  testURL,
		requestURL: hostURL,
	})

	//Then
	expectation := ShortenURLShortIDInUse
	assert.NotNil(suite.T(), err, "ShortenURL: Expected Error, got nil")
	assert.Equal(suite.T(), expectation, int(err.Code()), "ShortenURL generates wrong error code. Expected '%v'. Got: %v", expectation, err.Code())
}

func (suite *ShortenURLUseCaseTestSuite) TestGivenShortIDIsGenerated_WhenShortIDNotUnique_ThenReturnError() {

	//Given
	suite.generator.ShortID = savedShortID
	suite.urlRepo.ReturnError = true
	suite.urlRepo.SaveURLRecordError = errors.New("short id exists")

	hostURL, _ := url.Parse("http://www.small.ml")
	testURL, _ := url.Parse("http://www.2.com")

	//When
	_, err := suite.useCase.Execute(ShortenURLRequest{
		LongURL:    "http://www.2.com",
		parsedURL:  testURL,
		requestURL: hostURL,
	})

	expectation := ShortenURLFailedToSave
	assert.NotNil(suite.T(), err, "ShortenURL: Expected Error, got nil")
	assert.Equal(suite.T(), expectation, int(err.Code()), "ShortenURL wrong error code. Expected '%d'. Got: %d", expectation, err)
}
