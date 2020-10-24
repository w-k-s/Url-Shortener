package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	u "github.com/w-k-s/short-url/domain/urlshortener"
	"os"
	"testing"
	"time"
)

const savedShortID = "shorty"
const savedLongURL = "http://www.examply.com"
const savedShortURL = "http://small.ml/" + savedShortID

type URLRepositoryTestSuite struct {
	suite.Suite
	db      *sql.DB
	urlRepo *DefaultURLRepository
	record  *u.URLRecord
}

func TestURLRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(URLRepositoryTestSuite))
}

func (suite *URLRepositoryTestSuite) SetupTest() {
	connStr := os.Getenv("TEST_DB_CONN_STRING")
	if len(connStr) == 0 {
		connStr = "postgres://localhost/url_shortener_test?sslmode=disable"
	}

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}

	suite.db = db
	suite.urlRepo = NewURLRepository(suite.db)

	suite.record = &u.URLRecord{
		LongURL:    savedLongURL,
		ShortID:    savedShortID,
		CreateTime: time.Now(),
	}
}

func (suite *URLRepositoryTestSuite) TearDownTest() {
	_, err := suite.db.Exec("DELETE FROM url_records")
	if err != nil {
		panic(err)
	}
}

func (suite *URLRepositoryTestSuite) TestSaveRecordSucccessful() {

	_, err := suite.urlRepo.SaveRecord(suite.record)

	assert.Nil(suite.T(), err, "Expected: save record. Got: %s", err)
}

func (suite *URLRepositoryTestSuite) TestDuplicateRecordFails() {
	suite.urlRepo.SaveRecord(suite.record)
	_, err := suite.urlRepo.SaveRecord(suite.record)

	assert.True(suite.T(), suite.urlRepo.IsDup(err), "Expected: duplication error. Got: %s", err)
}

func (suite *URLRepositoryTestSuite) TestFindExistingShortURL() {
	_, err := suite.urlRepo.SaveRecord(suite.record)
	if err != nil {
		panic(err)
	}

	result, err := suite.urlRepo.ShortURL(suite.record.LongURL)
	expectation := result != nil && result.ShortID == suite.record.ShortID
	assert.True(suite.T(), expectation, "Expected Matching ShortId '%s'. Got: '%v' (error: '%s')", suite.record.ShortID, result.ShortID, err)
}

func (suite *URLRepositoryTestSuite) TestFindAbsentShortURL() {

	result, err := suite.urlRepo.ShortURL("http://www.nil.com")
	assert.NotNil(suite.T(), err, "Expected err when shortId not found. Got: nil. (record: %v)", result)
}

func (suite *URLRepositoryTestSuite) TestFindExistingLongURL() {
	_, err := suite.urlRepo.SaveRecord(suite.record)
	if err != nil {
		panic(err)
	}

	result, err := suite.urlRepo.LongURL(suite.record.ShortID)
	expectation := result != nil && result.LongURL == suite.record.LongURL

	assert.True(suite.T(), expectation, "Expected Matching LongURL '%s'. Got: '%v' (error: '%s')", suite.record.LongURL, result.LongURL, err)
}

func (suite *URLRepositoryTestSuite) TestFindAbsentLongURL() {

	result, err := suite.urlRepo.LongURL("nil")
	assert.NotNil(suite.T(), err, "Expected err when longUrl not found. Got: nil. (record: %v)", result)

}
