package tests

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	database "github.com/w-k-s/short-url/db"
	repo "github.com/w-k-s/short-url/logging"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http/httptest"
	"os"
	"testing"
)

type LogRepositoryTestSuite struct {
	suite.Suite
	db      *database.Db
	logRepo *repo.LogRepository
}

func TestLogRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(LogRepositoryTestSuite))
}

func (suite *LogRepositoryTestSuite) SetupTest() {
	logger := log.New(os.Stdout, "short-url: ", log.Ldate|log.Ltime)
	suite.db = database.New("mongodb://localhost:27017/shorturl_test")
	suite.logRepo = repo.NewLogRepository(logger, suite.db)
}

func (suite *LogRepositoryTestSuite) TearDownTest() {

	suite.db.Instance().
		C("logs").
		RemoveAll(bson.M{})

	defer suite.db.Close()

}

func (suite *LogRepositoryTestSuite) TestSaveRecordSucccessful() {

	stringBody := "{\"longUrl\":\"http://www.eg.com\"}"
	jsonBytes := bytes.NewBuffer([]byte(stringBody))
	req := httptest.NewRequest("POST", "http://small.ml/urlshortener/v", jsonBytes)

	logRecord := suite.logRepo.LogRequest(req)

	assert.NotNil(suite.T(), logRecord.Time)
	assert.NotNil(suite.T(), logRecord.IPAddress)
	assert.Equal(suite.T(), "POST", logRecord.Method)
	assert.Equal(suite.T(), req.RequestURI, logRecord.URI)
	assert.Equal(suite.T(), stringBody, logRecord.Body)
	assert.Equal(suite.T(), 0, logRecord.Status)

	sw := repo.StatusWriter{ResponseWriter: httptest.NewRecorder()}
	sw.WriteHeader(200)
	err := suite.logRepo.LogResponse(&sw, logRecord)

	assert.Equal(suite.T(), 200, logRecord.Status)
	assert.Nil(suite.T(), err, "Expected: save record. Got: %s", err)
}
