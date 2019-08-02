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
	ShortID string
}

func (m MockShortIDGenerator) Generate(d ShortIDLength) string {
	return m.ShortID
}

//-----

const savedShortID = "shrt"
const savedLongURL = "http://www.example.com"
const savedShortURL = "http://small.ml/" + savedShortID

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
		savedLongURL,
		savedShortID,
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
	testURL, _ := url.Parse(savedLongURL)
	shortURL, _ := suite.service.ShortenURL(hostURL, shortenURLRequest{
		parsedURL: testURL,
	})
	expectation := "http://www.small.ml/" + savedShortID

	assert.Equal(suite.T(), shortURL.String(), expectation, "ShortenURL generates wrong url. Expected '%s'. Got: %s", expectation, shortURL.String())
}

func (suite *ServiceSuite) TestShortURLCreatedWhenRecordDoesNotExist() {

	hostURL, _ := url.Parse("http://www.small.ml")
	testURL, _ := url.Parse("http://www.1.com")
	suite.generator.ShortID = "alpha"
	shortURL, _ := suite.service.ShortenURL(hostURL, shortenURLRequest{
		parsedURL: testURL,
	})
	expectation := "http://www.small.ml/" + suite.generator.ShortID

	assert.Equal(suite.T(), shortURL.String(), expectation, "ShortenURL generates wrong url. Expected '%s'. Got: %s", expectation, shortURL.String())
}

func (suite *ServiceSuite) TestCustomShortIDIsUsedWhenNotInUse() {

	hostURL, _ := url.Parse("http://www.small.ml")
	testURL, _ := url.Parse("http://www.google.com")

	shortURL, _ := suite.service.ShortenURL(hostURL, shortenURLRequest{
		parsedURL: testURL,
		ShortID:   "custom",
	})
	expectation := "http://www.small.ml/custom"

	assert.Equal(suite.T(), shortURL.String(), expectation, "ShortenURL generates wrong url. Expected '%s'. Got: %s", expectation, shortURL.String())

}

func (suite *ServiceSuite) TestInUseShortIDReturnedEvenThoughCustomShortIDProvided() {

	hostURL, _ := url.Parse("http://www.small.ml")
	testURL, _ := url.Parse(savedLongURL)

	shortURL, _ := suite.service.ShortenURL(hostURL, shortenURLRequest{
		parsedURL: testURL,
		ShortID:   "custom",
	})
	expectation := "http://www.small.ml/" + savedShortID

	assert.Equal(suite.T(), shortURL.String(), expectation, "ShortenURL generates wrong url. Expected '%s'. Got: %s", expectation, shortURL.String())

}

func (suite *ServiceSuite) TestErrorWhenCustomShortIDIsNotUnique() {

	hostURL, _ := url.Parse("http://www.small.ml")
	testURL, _ := url.Parse("http://www.google.com")

	_, err := suite.service.ShortenURL(hostURL, shortenURLRequest{
		parsedURL: testURL,
		ShortID:   savedShortID,
	})
	expectation := ShortenURLShortIdInUse

	assert.Equal(suite.T(), int(err.Code()), expectation, "ShortenURL generates wrong error code. Expected '%v'. Got: %v", expectation, err.Code())

}

func (suite *ServiceSuite) TestShortUrlErrorWhenShortIDNotUnique() {

	hostURL, _ := url.Parse("http://www.small.ml")
	testURL, _ := url.Parse("http://www.2.com")
	suite.generator.ShortID = savedShortID
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

	testURL, _ := url.Parse(savedShortURL)
	url, _ := suite.service.GetLongURL(testURL)

	assert.Equal(suite.T(), url.String(), savedLongURL, "GetLongURL returned wrong original url. Expected %s, Got: %s", savedLongURL, url.String())
}

func (suite *ServiceSuite) TestGetLongURLErrorWhenRecordExistsButInvalidURL() {

	record := &URLRecord{
		LongURL:    "",
		ShortID:    "wrong",
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
