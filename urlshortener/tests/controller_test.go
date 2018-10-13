package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/w-k-s/short-url/db"
	err "github.com/w-k-s/short-url/error"
	u "github.com/w-k-s/short-url/urlshortener"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

type ControllerSuite struct {
	suite.Suite
	db         *db.Db
	urlRepo    *u.URLRepository
	record     *u.URLRecord
	generator  *MockShortIDGenerator
	service    *u.Service
	controller *u.Controller
}

func (suite *ControllerSuite) SetupTest() {
	suite.db = db.New("mongodb://localhost:27017/shorturl_test")
	suite.urlRepo = u.NewURLRepository(suite.db)

	logger := log.New(os.Stdout, "short-url: ", log.Ldate|log.Ltime)
	suite.generator = &MockShortIDGenerator{}

	suite.service = u.NewService(suite.urlRepo, logger, suite.generator)
	suite.controller = u.NewController(suite.service)

	suite.db.Instance().
		C("urls").
		RemoveAll(bson.M{})

	suite.record = &u.URLRecord{
		SAVED_LONG_URL,
		SAVED_SHORT_ID,
		time.Now(),
	}

	_, err := suite.urlRepo.SaveRecord(suite.record)
	if err != nil {
		panic(err)
	}
}

func (suite *ControllerSuite) TearDownTest() {
	defer suite.db.Close()
}

func TestControllerSuite(t *testing.T) {
	suite.Run(t, new(ControllerSuite))
}

func (suite *ControllerSuite) TestGetShortURLWhenRequestBodyDoesNotContainLongURL() {

	jsonBytes := bytes.NewBuffer([]byte("{}"))
	req := httptest.NewRequest("POST", "http://small.ml/urlshortener/v", jsonBytes)
	w := httptest.NewRecorder()
	suite.controller.ShortenURL(w, req)

	error := getErrOrNil(w)
	assert.Equal(suite.T(), err.Code(u.ShortenURLValidation), error.Code(), "Wrong error code. Expected: %d, got: %d", u.ShortenURLValidation, error.Code())
}

func (suite *ControllerSuite) TestGetShortURLWhenRequestBodyDoesNotContainValidURL() {

	jsonBytes := bytes.NewBuffer([]byte("{\"longUrl\":\"hello there\"}"))
	req := httptest.NewRequest("POST", "http://small.ml/urlshortener/v", jsonBytes)
	w := httptest.NewRecorder()
	suite.controller.ShortenURL(w, req)

	error := getErrOrNil(w)
	assert.Equal(suite.T(), err.Code(u.ShortenURLValidation), error.Code(), "Wrong error code. Expected: %d, got: %d", u.ShortenURLValidation, error.Code())

}

func (suite *ControllerSuite) TestGetShortURLWhenRequestBodyContainsRelativeURL() {

	jsonBytes := bytes.NewBuffer([]byte("{\"longUrl\":\"path/to/file\"}"))
	req := httptest.NewRequest("POST", "http://small.ml/urlshortener/v", jsonBytes)
	w := httptest.NewRecorder()
	suite.controller.ShortenURL(w, req)

	error := getErrOrNil(w)
	assert.Equal(suite.T(), err.Code(u.ShortenURLValidation), error.Code(), "Wrong error code. Expected: %d, got: %d", u.ShortenURLValidation, error.Code())

}

func (suite *ControllerSuite) TestGetShortURLSuccessResponseWhenShortURLGenerated() {

	suite.generator.ShortId = "unique"
	jsonBytes := bytes.NewBuffer([]byte("{\"longUrl\":\"http://www.eg.com\"}"))
	req := httptest.NewRequest("POST", "http://small.ml/urlshortener/v", jsonBytes)
	w := httptest.NewRecorder()
	suite.controller.ShortenURL(w, req)

	assert.Equal(suite.T(), "application/json;charset=utf-8", w.Header()["Content-Type"][0])
	assert.Contains(suite.T(), w.Header(), "Etag")

	json := getJSONDictionaryOrNil(w)
	shortUrl := json["shortUrl"].(string)
	assert.Contains(suite.T(), shortUrl, suite.generator.ShortId, "Generated shortid '%s' not in short url '%s'", suite.generator.ShortId, shortUrl)

}

func (suite *ControllerSuite) TestRedirectSuccessResponseWhenShortURLExists() {

	req := httptest.NewRequest("GET", SAVED_SHORT_URL, nil)
	w := httptest.NewRecorder()
	suite.controller.RedirectToLongURL(w, req)
	resp := w.Result()

	assert.Equal(suite.T(), resp.StatusCode, http.StatusSeeOther)
}

func (suite *ControllerSuite) TestRedirectFailureResponseWhenShortURLDoesNotExist() {

	req := httptest.NewRequest("GET", "http://www.small.ml/nil", nil)
	w := httptest.NewRecorder()
	suite.controller.RedirectToLongURL(w, req)
	resp := w.Result()

	assert.Equal(suite.T(), resp.StatusCode, http.StatusNotFound)

	error := getErrOrNil(w)
	assert.Equal(suite.T(), err.Code(u.RetrieveFullURLNotFound), error.Code(), "Wrong error code. Expected: %d, got: %d", u.RetrieveFullURLNotFound, error.Code())

}

func (suite *ControllerSuite) TestGetLongURLRequestWhenRequestQueryDoesNotContainShortURL() {

	req := httptest.NewRequest("GET", "http://www.small.ml", nil)
	w := httptest.NewRecorder()
	suite.controller.GetLongURL(w, req)

	error := getErrOrNil(w)
	assert.Equal(suite.T(), err.Code(u.RetrieveFullURLValidation), error.Code(), "Wrong error code. Expected: %d, got: %d", u.RetrieveFullURLValidation, error.Code())

}

func (suite *ControllerSuite) TestGetLongURLRequestWhenRequestQueryContainInvalidShortURL() {

	req := httptest.NewRequest("GET", "http://www.small.ml?shortUrlhello%20there", nil)
	w := httptest.NewRecorder()
	suite.controller.GetLongURL(w, req)

	error := getErrOrNil(w)
	assert.Equal(suite.T(), err.Code(u.RetrieveFullURLValidation), error.Code(), "Wrong error code. Expected: %d, got: %d", u.RetrieveFullURLValidation, error.Code())

}

func (suite *ControllerSuite) TestGetLongURLRequestWhenRequestQueryContainRelativeShortURL() {

	req := httptest.NewRequest("GET", "http://www.small.ml?shortUrl=path/to/file", nil)
	w := httptest.NewRecorder()
	suite.controller.GetLongURL(w, req)

	error := getErrOrNil(w)
	assert.Equal(suite.T(), err.Code(u.RetrieveFullURLValidation), error.Code(), "Wrong error code. Expected: %d, got: %d", u.RetrieveFullURLValidation, error.Code())

}

func (suite *ControllerSuite) TestGetLongURLRequestWhenLongURLNotFound() {

	req := httptest.NewRequest("GET", "http://www.small.ml?shortUrl=http://www.small.ml/nil", nil)
	w := httptest.NewRecorder()
	suite.controller.GetLongURL(w, req)

	error := getErrOrNil(w)
	assert.Equal(suite.T(), err.Code(u.RetrieveFullURLNotFound), error.Code(), "Wrong error code. Expected: %d, got: %d", u.RetrieveFullURLNotFound, error)

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

func getErrOrNil(w *httptest.ResponseRecorder) err.Err {
	JSONDictionary := getJSONDictionaryOrNil(w)
	if JSONDictionary == nil {
		return nil
	}

	var ok bool
	var code float64
	var domain string
	var message string
	var fields map[string]string

	fmt.Printf("JSONDictionary -> %v\n", JSONDictionary)

	if code, ok = JSONDictionary["code"].(float64); !ok {
		fmt.Printf("getErrOrNil -> couldnt map `code\n`")
	}

	if domain, ok = JSONDictionary["domain"].(string); !ok {
		fmt.Printf("getErrOrNil -> couldnt map `domain`\n")
	}

	if message, ok = JSONDictionary["message"].(string); !ok {
		fmt.Printf("getErrOrNil -> couldnt map `message`\n")
	}

	if fields, ok = JSONDictionary["fields"].(map[string]string); !ok {
		fmt.Printf("getErrOrNil -> couldnt map `fields`\n")
	}

	return err.NewError(err.Code(int(code)), domain, message, fields)
}
