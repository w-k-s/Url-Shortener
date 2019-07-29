package urlshortener

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	database "github.com/w-k-s/short-url/db"
	"github.com/w-k-s/short-url/error"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/url"
	"os"
	"testing"
	"time"
)

//-- MockShortIDGenerator

type MockShortIDGenerator struct {
	ShortId string
}

func (m MockShortIDGenerator) Generate(d ShortIdLength) string {
	return m.ShortId
}

//-----

const SAVED_SHORT_ID = "shrt"
const SAVED_LONG_URL = "http://www.example.com"
const SAVED_SHORT_URL = "http://small.ml/" + SAVED_SHORT_ID

type ServiceSuite struct {
	suite.Suite
	db        *database.Db
	urlRepo   *URLRepository
	record    *URLRecord
	generator *MockShortIDGenerator
	service   *Service
}

func (suite *ServiceSuite) SetupTest() {
	suite.db = database.New("mongodb://localhost:27017/shorturl_test")

	suite.record = &URLRecord{
		SAVED_LONG_URL,
		SAVED_SHORT_ID,
		time.Now(),
	}

	suite.db.Instance().
		C("urls").
		RemoveAll(bson.M{})

	logger := log.New(os.Stdout, "short-url: ", log.Ldate|log.Ltime)
	suite.generator = &MockShortIDGenerator{}
	suite.urlRepo = NewURLRepository(suite.db, logger)
	suite.service = NewService(suite.urlRepo, logger, suite.generator)

	_, err := suite.urlRepo.SaveRecord(suite.record)
	if err != nil {
		panic(err)
	}
}

func (suite *ServiceSuite) TearDownTest() {
	defer suite.db.Close()
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}

func (suite *ServiceSuite) TestShortURLReturnedWhenRecordExists() {

	hostURL, _ := url.Parse("http://www.small.ml")
	testURL, _ := url.Parse(SAVED_LONG_URL)
	shortURL, _ := suite.service.ShortenURL(hostURL, shortenURLRequest{
		parsedURL: testURL,
	})
	expectation := "http://www.small.ml/" + SAVED_SHORT_ID

	assert.Equal(suite.T(), shortURL.String(), expectation, "ShortenURL generates wrong url. Expected '%s'. Got: %s", expectation, shortURL.String())
}

func (suite *ServiceSuite) TestShortURLCreatedWhenRecordDoesNotExist() {

	hostURL, _ := url.Parse("http://www.small.ml")
	testURL, _ := url.Parse("http://www.1.com")
	suite.generator.ShortId = "alpha"
	shortURL, _ := suite.service.ShortenURL(hostURL, shortenURLRequest{
		parsedURL: testURL,
	})
	expectation := "http://www.small.ml/" + suite.generator.ShortId

	assert.Equal(suite.T(), shortURL.String(), expectation, "ShortenURL generates wrong url. Expected '%s'. Got: %s", expectation, shortURL.String())
}

func (suite *ServiceSuite) TestCustomShortIdIsUsedWhenNotInUse() {

	hostURL, _ := url.Parse("http://www.small.ml")
	testURL, _ := url.Parse("http://www.google.com")

	shortURL, _ := suite.service.ShortenURL(hostURL, shortenURLRequest{
		parsedURL: testURL,
		ShortId:   "custom",
	})
	expectation := "http://www.small.ml/custom"

	assert.Equal(suite.T(), shortURL.String(), expectation, "ShortenURL generates wrong url. Expected '%s'. Got: %s", expectation, shortURL.String())

}

func (suite *ServiceSuite) TestInUseShortIdReturnedEvenThoughCustomShortIdProvided() {

	hostURL, _ := url.Parse("http://www.small.ml")
	testURL, _ := url.Parse(SAVED_LONG_URL)

	shortURL, _ := suite.service.ShortenURL(hostURL, shortenURLRequest{
		parsedURL: testURL,
		ShortId:   "custom",
	})
	expectation := "http://www.small.ml/" + SAVED_SHORT_ID

	assert.Equal(suite.T(), shortURL.String(), expectation, "ShortenURL generates wrong url. Expected '%s'. Got: %s", expectation, shortURL.String())

}

func (suite *ServiceSuite) TestErrorWhenCustomShortIdIsNotUnique() {

	hostURL, _ := url.Parse("http://www.small.ml")
	testURL, _ := url.Parse("http://www.google.com")

	_, err := suite.service.ShortenURL(hostURL, shortenURLRequest{
		parsedURL: testURL,
		ShortId:   SAVED_SHORT_ID,
	})
	expectation := ShortenURLShortIdInUse

	assert.Equal(suite.T(), int(err.Code()), expectation, "ShortenURL generates wrong error code. Expected '%v'. Got: %v", expectation, err.Code())

}

func (suite *ServiceSuite) TestShortUrlErrorWhenShortIDNotUnique() {

	hostURL, _ := url.Parse("http://www.small.ml")
	testURL, _ := url.Parse("http://www.2.com")
	suite.generator.ShortId = SAVED_SHORT_ID
	_, err := suite.service.ShortenURL(hostURL, shortenURLRequest{
		parsedURL: testURL,
	})
	expectation := error.Code(ShortenURLFailedToSave)

	assert.True(suite.T(), err != nil && err.Code() == expectation, "ShortenURL wrong error code. Expected '%d'. Got: %d", expectation, err)
}

func (suite *ServiceSuite) TestGetLongURLErrorWhenShortURLHasNoPath() {

	testURL, _ := url.Parse("http://www.small.ml")
	_, err := suite.service.GetLongURL(testURL)
	expectation := error.Code(RetrieveFullURLValidation)

	assert.True(suite.T(), err != nil && err.Code() == expectation, "GetLongURL wrong error code. Expected '%d'. Got: %d", expectation, err)
}

func (suite *ServiceSuite) TestGetLongURLErrorWhenRecordDoesNotExist() {

	testURL, _ := url.Parse("http://www.small.ml/nil")
	_, err := suite.service.GetLongURL(testURL)
	expectation := error.Code(RetrieveFullURLNotFound)

	assert.True(suite.T(), err != nil && err.Code() == expectation, "GetLongURL wrong error code. Expected '%d'. Got: %d", expectation, err)

}

func (suite *ServiceSuite) TestGetLongURLWhenRecordExists() {

	testURL, _ := url.Parse(SAVED_SHORT_URL)
	url, _ := suite.service.GetLongURL(testURL)

	assert.Equal(suite.T(), url.String(), SAVED_LONG_URL, "GetLongURL returned wrong original url. Expected %s, Got: %s", SAVED_LONG_URL, url.String())
}

func (suite *ServiceSuite) TestGetLongURLErrorWhenRecordExistsButInvalidURL() {

	record := &URLRecord{
		LongURL:    "",
		ShortId:    "wrong",
		CreateTime: time.Now(),
	}

	_, err := suite.urlRepo.SaveRecord(record)
	if err != nil {
		panic(err)
	}

	testURL, _ := url.Parse("http://small.ml/wrong")
	r, _ := suite.service.GetLongURL(testURL)
	assert.Equal(suite.T(), r.String(), "", "GetLongURL wrong error code. Expected '%v'. Got: %v", "", r.String())
}
