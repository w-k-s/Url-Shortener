package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/w-k-s/short-url/domain"
	u "github.com/w-k-s/short-url/domain/urlshortener"
	"github.com/w-k-s/short-url/domain/urlshortener/usecase"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"
)

const savedShortID = "shrt"
const savedLongURL = "http://www.example.com"
const savedShortURL = "http://small.ml/" + savedShortID

//-- MockShortIDGenerator

type MockShortIDGenerator struct {
	ShortID string
}

func (m MockShortIDGenerator) Generate(d usecase.ShortIDLength) string {
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

type ControllerSuite struct {
	suite.Suite
	urlRepo                    *MockURLRepository
	record                     *u.URLRecord
	generator                  *MockShortIDGenerator
	shortenURLUseCase          *usecase.ShortenURLUseCase
	retrieveOriginalURLUseCase *usecase.RetrieveOriginalURLUseCase
	controller                 *Controller
}

func (suite *ControllerSuite) SetupTest() {
	logger := log.New(os.Stdout, "short-url: ", log.Ldate|log.Ltime)

	baseURL, _ := url.Parse("https://small.ml")

	suite.generator = &MockShortIDGenerator{}

	suite.urlRepo = &MockURLRepository{}
	suite.shortenURLUseCase = usecase.NewShortenURLUseCase(suite.urlRepo, baseURL, suite.generator, logger)
	suite.retrieveOriginalURLUseCase = usecase.NewRetrieveOriginalURLUseCase(suite.urlRepo, logger)
	suite.controller = NewController(suite.shortenURLUseCase, suite.retrieveOriginalURLUseCase, logger)

	suite.record = &u.URLRecord{
		LongURL:    savedLongURL,
		ShortID:    savedShortID,
		CreateTime: time.Now(),
	}

	_, err := suite.urlRepo.SaveRecord(suite.record)
	if err != nil {
		panic(fmt.Sprintf("Setup Test: %s", err.Error()))
	}
}

func TestControllerSuite(t *testing.T) {
	suite.Run(t, new(ControllerSuite))
}

func (suite *ControllerSuite) TestGivenEmptyBody_WhenShorteningURL_ThenReturnsError() {
	//Given
	jsonBytes := bytes.NewBuffer([]byte("{}"))

	//When
	req := httptest.NewRequest("POST", "http://small.ml/urlshortener/v", jsonBytes)
	w := httptest.NewRecorder()
	suite.controller.ShortenURL(w, req)

	//Then
	err := getErrOrNil(w)
	assert.NotNil(suite.T(), err, "ShortURL: Expected error; got nil")
	assert.Equal(suite.T(), domain.Code(usecase.ShortenURLValidation), err.Code(), "Wrong error code. Expected: %d, got: %d", usecase.ShortenURLValidation, err.Code())
}

func (suite *ControllerSuite) TestGivenInvalidLongURL_WhenShorteningURL_ThenReturnError() {
	//Given
	jsonBytes := bytes.NewBuffer([]byte("{\"longUrl\":\"hello there\"}"))

	//When
	req := httptest.NewRequest("POST", "http://small.ml/urlshortener/v", jsonBytes)
	w := httptest.NewRecorder()
	suite.controller.ShortenURL(w, req)

	//Then
	err := getErrOrNil(w)
	assert.NotNil(suite.T(), err, "ShortURL: Expected error; got nil")
	assert.Equal(suite.T(), domain.Code(usecase.ShortenURLValidation), err.Code(), "Wrong error code. Expected: %d, got: %d", usecase.ShortenURLValidation, err.Code())

}

func (suite *ControllerSuite) TestGivenRelativeLongURL_WhenShorteningURL_ThenReturnError() {
	//Given
	jsonBytes := bytes.NewBuffer([]byte("{\"longUrl\":\"path/to/file\"}"))

	//When
	req := httptest.NewRequest("POST", "http://small.ml/urlshortener/v", jsonBytes)
	w := httptest.NewRecorder()
	suite.controller.ShortenURL(w, req)

	//Then
	err := getErrOrNil(w)
	assert.NotNil(suite.T(), err, "ShortURL: Expected error; got nil")
	assert.Equal(suite.T(), domain.Code(usecase.ShortenURLValidation), err.Code(), "Wrong error code. Expected: %d, got: %d", usecase.ShortenURLValidation, err.Code())

}

func (suite *ControllerSuite) TestGivenLongURL_WhenShorteningURL_() {

	//Given
	jsonBytes := bytes.NewBuffer([]byte("{\"longUrl\":\"http://www.eg.com\"}"))
	suite.generator.ShortID = "unique"
	suite.urlRepo.ShortURLRecordResult = &u.URLRecord{
		LongURL:    "http://www.eg.com",
		ShortID:    "unique",
		CreateTime: time.Now(),
	}

	//When
	req := httptest.NewRequest("POST", "http://small.ml/urlshortener/v", jsonBytes)
	w := httptest.NewRecorder()
	suite.controller.ShortenURL(w, req)

	//Then
	assert.Equal(suite.T(), "application/json;charset=utf-8", w.Header()["Content-Type"][0])

	json := getJSONDictionaryOrNil(w)
	shortURL := json["shortUrl"].(string)
	assert.Contains(suite.T(), shortURL, suite.generator.ShortID, "Generated shortid '%s' not in short url '%s'", suite.generator.ShortID, shortURL)

}

func (suite *ControllerSuite) TestGivenShortURLExists_WhenRedirecting_ThenSeeOtherResponse() {

	//Given
	suite.urlRepo.LongURLRecordResult = suite.record

	//When
	req := httptest.NewRequest("GET", savedShortURL, nil)
	w := httptest.NewRecorder()
	suite.controller.RedirectToLongURL(w, req)

	//Then
	resp := w.Result()
	assert.Equal(suite.T(), resp.StatusCode, http.StatusSeeOther)
}

func (suite *ControllerSuite) TestGivenShortURLDoesNotExist_WhenRedirecting_ThenNotFoundResponse() {
	//Given
	suite.urlRepo.ReturnError = true
	suite.urlRepo.LongURLRecordError = errors.New("Not Found")

	//When
	req := httptest.NewRequest("GET", "http://www.small.ml/nil", nil)
	w := httptest.NewRecorder()
	suite.controller.RedirectToLongURL(w, req)

	//Then
	resp := w.Result()
	assert.Equal(suite.T(), resp.StatusCode, http.StatusNotFound)

	err := getErrOrNil(w)
	assert.NotNil(suite.T(), err, "ShortURL: Expected error; got nil")
	assert.Equal(suite.T(), domain.Code(usecase.RetrieveFullURLNotFound), err.Code(), "Wrong error code. Expected: %d, got: %d", usecase.RetrieveFullURLNotFound, err.Code())

}

func (suite *ControllerSuite) TestGivenNoShortURL_WhenGetLongURLRequest_ThenReturnRetrieveFullURLValidationError() {

	req := httptest.NewRequest("GET", "http://www.small.ml", nil)
	w := httptest.NewRecorder()
	suite.controller.GetLongURL(w, req)

	err := getErrOrNil(w)
	assert.NotNil(suite.T(), err, "ShortURL: Expected error; got nil")
	assert.Equal(suite.T(), domain.Code(usecase.RetrieveFullURLValidation), err.Code(), "Wrong error code. Expected: %d, got: %d", usecase.RetrieveFullURLValidation, err.Code())

}

func (suite *ControllerSuite) TestGivenInvalidShortURL_WhenGetLongURLRequest_ThenReturnRetrieveFullURLValidationError() {

	//When
	req := httptest.NewRequest("GET", "http://www.small.ml?shortUrlhello%20there", nil)
	w := httptest.NewRecorder()
	suite.controller.GetLongURL(w, req)

	//Then
	err := getErrOrNil(w)
	assert.NotNil(suite.T(), err, "ShortURL: Expected error; got nil")
	assert.Equal(suite.T(), domain.Code(usecase.RetrieveFullURLValidation), err.Code(), "Wrong error code. Expected: %d, got: %d", usecase.RetrieveFullURLValidation, err.Code())

}

func (suite *ControllerSuite) TestGivenRelativeShortURL_WhenGetLongURLRequest_ThenReturnRetrieveFullURLValidationError() {

	//When
	req := httptest.NewRequest("GET", "http://www.small.ml?shortUrl=path/to/file", nil)
	w := httptest.NewRecorder()
	suite.controller.GetLongURL(w, req)

	//Then
	err := getErrOrNil(w)
	assert.NotNil(suite.T(), err, "ShortURL: Expected error; got nil")
	assert.Equal(suite.T(), domain.Code(usecase.RetrieveFullURLValidation), err.Code(), "Wrong error code. Expected: %d, got: %d", usecase.RetrieveFullURLValidation, err.Code())
}

func (suite *ControllerSuite) TestGivenShortURLDoesNotExist_WhenGetLongURLRequest_ThenRetrieveFullURLNotFoundError() {
	//Given
	suite.urlRepo.ReturnError = true
	suite.urlRepo.LongURLRecordError = errors.New("Not Found")

	//When
	req := httptest.NewRequest("GET", "http://www.small.ml?shortUrl=http://www.small.ml/nil", nil)
	w := httptest.NewRecorder()
	suite.controller.GetLongURL(w, req)

	//Then
	err := getErrOrNil(w)
	assert.NotNil(suite.T(), err, "ShortURL: Expected error; got nil")
	assert.Equal(suite.T(), domain.Code(usecase.RetrieveFullURLNotFound), err.Code(), "Wrong error code. Expected: %d, got: %d", usecase.RetrieveFullURLNotFound, err.Code())
}

func getJSONDictionaryOrNil(w *httptest.ResponseRecorder) map[string]interface{} {
	var JSONDictionary map[string]interface{}

	resp := w.Result()
	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(&JSONDictionary)
	if err != nil {
		return nil
	}
	return JSONDictionary
}

func getErrOrNil(w *httptest.ResponseRecorder) domain.Err {
	JSONDictionary := getJSONDictionaryOrNil(w)

	if JSONDictionary == nil {
		return nil
	}

	code := int(JSONDictionary["code"].(float64))
	domainString := JSONDictionary["domain"].(string)
	message := JSONDictionary["message"].(string)
	fields, _ := JSONDictionary["fields"].(map[string]string)

	return domain.NewError(domain.Code(code), domainString, message, fields)
}
