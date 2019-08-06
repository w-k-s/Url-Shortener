package usecase

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	u "github.com/w-k-s/short-url/domain/urlshortener"
	"log"
	"net/url"
	"os"
	"testing"
	"time"
)

type RetrieveOriginalURLUseCaseTestSuite struct {
	suite.Suite
	urlRepo *MockURLRepository
	record  *u.URLRecord
	useCase *RetrieveOriginalURLUseCase
}

func (suite *RetrieveOriginalURLUseCaseTestSuite) SetupTest() {
	suite.record = &u.URLRecord{
		savedLongURL,
		savedShortID,
		time.Now(),
	}

	logger := log.New(os.Stdout, "short-url: ", log.Ldate|log.Ltime)
	suite.urlRepo = &MockURLRepository{}
	suite.useCase = NewRetrieveOriginalURLUseCase(suite.urlRepo, logger)
}

func TestRetrieveOriginalURLUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(RetrieveOriginalURLUseCaseTestSuite))
}

func (suite *RetrieveOriginalURLUseCaseTestSuite) GivenShortURL_WhenShortlURLHasNoPath_ThenReturnError() {

	//Given
	testURL, _ := url.Parse("http://www.small.ml")

	//When
	_, err := suite.useCase.Execute(RetrieveOriginalURLRequest{
		shortURL: testURL,
	})

	//Then
	expectation := RetrieveFullURLValidation
	assert.Equal(suite.T(), expectation, int(err.Code()), "GetLongURL wrong error code. Expected '%d'. Got: %d", expectation, int(err.Code()))
}

func (suite *RetrieveOriginalURLUseCaseTestSuite) GivenShortURL_WhenRecordDoesNotExist_ThenReturnError() {

	//Given
	testURL, _ := url.Parse("http://www.small.ml/nil")

	//When
	_, err := suite.useCase.Execute(RetrieveOriginalURLRequest{
		shortURL: testURL,
	})

	//Then
	expectation := RetrieveFullURLNotFound
	assert.Equal(suite.T(), expectation, int(err.Code()), "GetLongURL wrong error code. Expected '%d'. Got: %d", expectation, int(err.Code()))

}

func (suite *RetrieveOriginalURLUseCaseTestSuite) GivenShortURL_WhenRecordExists_ThenReturnOriginalURL() {

	//Given
	testURL, _ := url.Parse(savedShortURL)
	suite.urlRepo.LongURLRecordResult = suite.record

	//When
	resp, _ := suite.useCase.Execute(RetrieveOriginalURLRequest{
		shortURL: testURL,
	})

	assert.Equal(suite.T(), savedLongURL, resp.LongURL, "GetLongURL returned wrong original url. Expected %s, Got: %s", savedLongURL, resp.LongURL)
}

func (suite *RetrieveOriginalURLUseCaseTestSuite) GivenRecordExists_WhenRecordHasInvalidURL_ThenReturnError() {

	//Given
	testURL, _ := url.Parse("http://small.ml/wrong")
	suite.urlRepo.LongURLRecordResult = &u.URLRecord{
		LongURL:    "@",
		ShortID:    "wrong",
		CreateTime: time.Now(),
	}

	//When
	_, err := suite.useCase.Execute(RetrieveOriginalURLRequest{
		shortURL: testURL,
	})

	//Then
	expectation := RetrieveFullURLParsing
	assert.Equal(suite.T(), expectation, int(err.Code()), "GetLongURL wrong error code. Expected '%v'. Got: %v", expectation, int(err.Code()))
}
