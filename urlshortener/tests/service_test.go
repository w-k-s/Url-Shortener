package tests

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	database "github.com/w-k-s/short-url/db"
	"github.com/w-k-s/short-url/error"
	u "github.com/w-k-s/short-url/urlshortener"
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

func (m MockShortIDGenerator) Generate(d u.Deviation) string {
	return m.ShortId
}

//-----

const SAVED_SHORT_ID = "shrt"
const SAVED_LONG_URL = "http://www.example.com"
const SAVED_SHORT_URL = "http://small.ml/" + SAVED_SHORT_ID

type ServiceSuite struct {
	suite.Suite
	db        *database.Db
	urlRepo   *u.URLRepository
	record    *u.URLRecord
	generator *MockShortIDGenerator
	service   *u.Service
}

func (suite *ServiceSuite) SetupTest() {
	suite.db = database.New("mongodb://localhost:27017/shorturl_test")

	suite.record = &u.URLRecord{
		SAVED_LONG_URL,
		SAVED_SHORT_ID,
		time.Now(),
	}

	suite.db.Instance().
		C("urls").
		RemoveAll(bson.M{})

	logger := log.New(os.Stdout, "short-url: ", log.Ldate|log.Ltime)
	suite.generator = &MockShortIDGenerator{}
	suite.urlRepo = u.NewURLRepository(suite.db, logger)
	suite.service = u.NewService(suite.urlRepo, logger, suite.generator)

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

func (suite *ServiceSuite) TestShortUrlReturnedWhenRecordExists() {

	hostUrl, _ := url.Parse("http://www.small.ml")
	testUrl, _ := url.Parse(SAVED_LONG_URL)
	shortUrl, _ := suite.service.ShortenURL(hostUrl, testUrl)
	expectation := "http://www.small.ml/" + SAVED_SHORT_ID

	assert.Equal(suite.T(), shortUrl.String(), expectation, "ShortenURL generates wrong url. Expected '%s'. Got: %s", expectation, shortUrl.String())
}

func (suite *ServiceSuite) TestShortUrlCreatedWhenRecordDoesNotExist() {

	hostUrl, _ := url.Parse("http://www.small.ml")
	testUrl, _ := url.Parse("http://www.1.com")
	suite.generator.ShortId = "alpha"
	shortUrl, _ := suite.service.ShortenURL(hostUrl, testUrl)
	expectation := "http://www.small.ml/" + suite.generator.ShortId

	assert.Equal(suite.T(), shortUrl.String(), expectation, "ShortenURL generates wrong url. Expected '%s'. Got: %s", expectation, shortUrl.String())
}

func (suite *ServiceSuite) TestShortUrlErrorWhenShortIDNotUnique() {

	hostUrl, _ := url.Parse("http://www.small.ml")
	testUrl, _ := url.Parse("http://www.2.com")
	suite.generator.ShortId = SAVED_SHORT_ID
	_, err := suite.service.ShortenURL(hostUrl, testUrl)
	expectation := error.Code(u.ShortenURLFailedToSave)

	assert.True(suite.T(), err != nil && err.Code() == expectation, "ShortenURL wrong error code. Expected '%d'. Got: %d", expectation, err)
}

func (suite *ServiceSuite) TestGetLongURLErrorWhenShortURLHasNoPath() {

	testUrl, _ := url.Parse("http://www.small.ml")
	_, err := suite.service.GetLongURL(testUrl)
	expectation := error.Code(u.RetrieveFullURLValidation)

	assert.True(suite.T(), err != nil && err.Code() == expectation, "GetLongURL wrong error code. Expected '%d'. Got: %d", expectation, err)
}

func (suite *ServiceSuite) TestGetLongURLErrorWhenRecordDoesNotExist() {

	testUrl, _ := url.Parse("http://www.small.ml/nil")
	_, err := suite.service.GetLongURL(testUrl)
	expectation := error.Code(u.RetrieveFullURLNotFound)

	assert.True(suite.T(), err != nil && err.Code() == expectation, "GetLongURL wrong error code. Expected '%d'. Got: %d", expectation, err)

}

func (suite *ServiceSuite) TestGetLongURLWhenRecordExists() {

	testUrl, _ := url.Parse(SAVED_SHORT_URL)
	url, _ := suite.service.GetLongURL(testUrl)

	assert.Equal(suite.T(), url.String(), SAVED_LONG_URL, "GetLongURL returned wrong original url. Expected %s, Got: %s", SAVED_LONG_URL, url.String())
}

func (suite *ServiceSuite) TestGetLongURLErrorWhenRecordExistsButInvalidURL() {

	record := &u.URLRecord{
		LongURL:    "",
		ShortId:    "wrong",
		CreateTime: time.Now(),
	}

	_, err := suite.urlRepo.SaveRecord(record)
	if err != nil {
		panic(err)
	}

	testUrl, _ := url.Parse("http://small.ml/wrong")
	r, _ := suite.service.GetLongURL(testUrl)
	assert.Equal(suite.T(), r.String(), "", "GetLongURL wrong error code. Expected '%v'. Got: %v", "", r.String())
}
