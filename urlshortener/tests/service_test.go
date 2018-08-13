package tests

import (
	_ "fmt"
	database "github.com/w-k-s/short-url/db"
	"github.com/w-k-s/short-url/error"
	u "github.com/w-k-s/short-url/urlshortener"
	_ "gopkg.in/mgo.v2"
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

var db *database.Db
var urlRepo *u.URLRepository
var record *u.URLRecord
var generator *MockShortIDGenerator
var service *u.Service

func TestMain(m *testing.M) {
	setup()
	retCode := m.Run()
	tearDown()
	os.Exit(retCode)
}

func setup() {
	db = database.New("mongodb://localhost:27017/shorturl_test")
	urlRepo = u.NewURLRepository(db)

	record = &u.URLRecord{
		SAVED_LONG_URL,
		SAVED_SHORT_ID,
		time.Now(),
	}

	db.Instance().
		C("urls").
		RemoveAll(bson.M{})

	_, err := urlRepo.SaveRecord(record)
	if err != nil {
		panic(err)
	}

	logger := log.New(os.Stdout, "short-url: ", log.Ldate|log.Ltime)
	generator = &MockShortIDGenerator{}

	service = u.NewService(urlRepo, logger, generator)
}

func tearDown() {
	defer db.Close()
}

func TestShortUrlReturnedWhenRecordExists(t *testing.T) {

	hostUrl, _ := url.Parse("http://www.small.ml")
	testUrl, _ := url.Parse(SAVED_LONG_URL)
	shortUrl, _ := service.ShortenUrl(hostUrl, testUrl)
	expectation := "http://www.small.ml/" + SAVED_SHORT_ID

	if shortUrl.String() != expectation {
		t.Errorf("ShortenURL generates wrong url. Expected '%s'. Got: %s", expectation, shortUrl.String())
	}

}

func TestShortUrlCreatedWhenRecordDoesNotExist(t *testing.T) {

	hostUrl, _ := url.Parse("http://www.small.ml")
	testUrl, _ := url.Parse("http://www.1.com")
	generator.ShortId = "alpha"
	shortUrl, _ := service.ShortenUrl(hostUrl, testUrl)
	expectation := "http://www.small.ml/" + generator.ShortId

	if shortUrl.String() != expectation {
		t.Errorf("ShortenURL generates wrong url. Expected '%s'. Got: %s", expectation, shortUrl.String())
	}

}

func TestShortUrlErrorWhenShortIDNotUnique(t *testing.T) {

	hostUrl, _ := url.Parse("http://www.small.ml")
	testUrl, _ := url.Parse("http://www.2.com")
	generator.ShortId = SAVED_SHORT_ID
	_, err := service.ShortenUrl(hostUrl, testUrl)
	expectation := error.Code(u.ShortenURLFailedToSave)

	if err == nil || err.Code() != expectation {
		t.Errorf("ShortenURL wrong error code. Expected '%d'. Got: %d", expectation, err)
	}

}

// func TestGetLongURLErrorWhenShortURLHasNoPath(t *testing.T){

// }

// func TestGetLongURLErrorWhenRecordDoesNotExist(t *testing.T){

// }

// func TestGetLongURLWhenRecordExists(t *testing.T){

// }

// func TestGetLongURLErrorWhenRecordExistsButInvalidURL(t *testing.T){

// }
