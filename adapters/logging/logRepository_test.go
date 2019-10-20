package logging

import (
	"bytes"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"log"
	"net/http/httptest"
	"os"
	"testing"
)

type LogRepositoryTestSuite struct {
	suite.Suite
	db      *sql.DB
	logRepo *LogRepository
}

func TestLogRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(LogRepositoryTestSuite))
}

func (suite *LogRepositoryTestSuite) SetupTest() {
	connStr := "postgres://localhost/url_shortener_test?sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}

	logger := log.New(os.Stdout, "short-url: ", log.Ldate|log.Ltime)
	suite.db = db
	suite.logRepo = NewLogRepository(suite.db, logger)
}

func (suite *LogRepositoryTestSuite) TearDownTest() {
	_, err := suite.db.Exec("DELETE FROM logs")
	if err != nil {
		panic(err)
	}
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

	sw := StatusWriter{ResponseWriter: httptest.NewRecorder()}
	sw.WriteHeader(200)
	err := suite.logRepo.LogResponse(&sw, logRecord)

	assert.Equal(suite.T(), 200, logRecord.Status)
	assert.Nil(suite.T(), err, "Expected: save record. Got: %s", err)
}
