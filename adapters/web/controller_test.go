package web

import (
	"bytes"
	"encoding/json"
	_ "fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/w-k-s/short-url/adapters/db"
	"github.com/w-k-s/short-url/domain"
	u "github.com/w-k-s/short-url/domain/urlshortener"
	"github.com/w-k-s/short-url/domain/urlshortener/usecase"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

const savedShortID = "shrt"
const savedLongURL = "http://www.example.com"
const savedShortURL = "http://small.ml/" + savedShortID

type MockShortIDGenerator struct {
	ShortID string
}

func (m MockShortIDGenerator) Generate(d usecase.ShortIDLength) string {
	return m.ShortID
}

type ControllerSuite struct {
	suite.Suite
	db                         *db.Db
	urlRepo                    *db.DefaultURLRepository
	record                     *u.URLRecord
	generator                  *MockShortIDGenerator
	shortenURLUseCase          *usecase.ShortenURLUseCase
	retrieveOriginalURLUseCase *usecase.RetrieveOriginalURLUseCase
	controller                 *Controller
}

func (suite *ControllerSuite) SetupTest() {
	suite.db = db.New("mongodb://localhost:27017/shorturl_test")

	logger := log.New(os.Stdout, "short-url: ", log.Ldate|log.Ltime)
	suite.generator = &MockShortIDGenerator{}

	suite.urlRepo = db.NewURLRepository(suite.db, logger)
	suite.shortenURLUseCase = usecase.NewShortenURLUseCase(suite.urlRepo, suite.generator, logger)
	suite.retrieveOriginalURLUseCase = usecase.NewRetrieveOriginalURLUseCase(suite.urlRepo, logger)
	suite.controller = NewController(suite.shortenURLUseCase, suite.retrieveOriginalURLUseCase, logger)

	suite.db.Instance().
		C("urls").
		RemoveAll(bson.M{})

	suite.db.Instance().
		C("visits").
		RemoveAll(bson.M{})

	suite.record = &u.URLRecord{
		savedLongURL,
		savedShortID,
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

	err := getErrOrNil(w)
	assert.Equal(suite.T(), domain.Code(usecase.ShortenURLValidation), err.Code(), "Wrong error code. Expected: %d, got: %d", usecase.ShortenURLValidation, err.Code())
}

func (suite *ControllerSuite) TestGetShortURLWhenRequestBodyDoesNotContainValidURL() {

	jsonBytes := bytes.NewBuffer([]byte("{\"longUrl\":\"hello there\"}"))
	req := httptest.NewRequest("POST", "http://small.ml/urlshortener/v", jsonBytes)
	w := httptest.NewRecorder()
	suite.controller.ShortenURL(w, req)

	err := getErrOrNil(w)
	assert.Equal(suite.T(), domain.Code(usecase.ShortenURLValidation), err.Code(), "Wrong error code. Expected: %d, got: %d", usecase.ShortenURLValidation, err.Code())

}

func (suite *ControllerSuite) TestGetShortURLWhenRequestBodyContainsRelativeURL() {

	jsonBytes := bytes.NewBuffer([]byte("{\"longUrl\":\"path/to/file\"}"))
	req := httptest.NewRequest("POST", "http://small.ml/urlshortener/v", jsonBytes)
	w := httptest.NewRecorder()
	suite.controller.ShortenURL(w, req)

	err := getErrOrNil(w)
	assert.Equal(suite.T(), domain.Code(usecase.ShortenURLValidation), err.Code(), "Wrong error code. Expected: %d, got: %d", usecase.ShortenURLValidation, err.Code())

}

func (suite *ControllerSuite) TestGetShortURLSuccessResponseWhenShortURLGenerated() {

	suite.generator.ShortID = "unique"
	jsonBytes := bytes.NewBuffer([]byte("{\"longUrl\":\"http://www.eg.com\"}"))
	req := httptest.NewRequest("POST", "http://small.ml/urlshortener/v", jsonBytes)
	w := httptest.NewRecorder()
	suite.controller.ShortenURL(w, req)

	assert.Equal(suite.T(), "application/json;charset=utf-8", w.Header()["Content-Type"][0])

	json := getJSONDictionaryOrNil(w)
	shortURL := json["shortUrl"].(string)
	assert.Contains(suite.T(), shortURL, suite.generator.ShortID, "Generated shortid '%s' not in short url '%s'", suite.generator.ShortID, shortURL)

}

func (suite *ControllerSuite) TestRedirectSuccessResponseWhenShortURLExists() {

	req := httptest.NewRequest("GET", savedShortURL, nil)
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

	err := getErrOrNil(w)
	assert.Equal(suite.T(), domain.Code(usecase.RetrieveFullURLNotFound), err.Code(), "Wrong error code. Expected: %d, got: %d", usecase.RetrieveFullURLNotFound, err.Code())

}

func (suite *ControllerSuite) TestGetLongURLRequestWhenRequestQueryDoesNotContainShortURL() {

	req := httptest.NewRequest("GET", "http://www.small.ml", nil)
	w := httptest.NewRecorder()
	suite.controller.GetLongURL(w, req)

	err := getErrOrNil(w)
	assert.Equal(suite.T(), domain.Code(usecase.RetrieveFullURLValidation), err.Code(), "Wrong error code. Expected: %d, got: %d", usecase.RetrieveFullURLValidation, err.Code())

}

func (suite *ControllerSuite) TestGetLongURLRequestWhenRequestQueryContainInvalidShortURL() {

	req := httptest.NewRequest("GET", "http://www.small.ml?shortUrlhello%20there", nil)
	w := httptest.NewRecorder()
	suite.controller.GetLongURL(w, req)

	err := getErrOrNil(w)
	assert.Equal(suite.T(), domain.Code(usecase.RetrieveFullURLValidation), err.Code(), "Wrong error code. Expected: %d, got: %d", usecase.RetrieveFullURLValidation, err.Code())

}

func (suite *ControllerSuite) TestGetLongURLRequestWhenRequestQueryContainRelativeShortURL() {

	req := httptest.NewRequest("GET", "http://www.small.ml?shortUrl=path/to/file", nil)
	w := httptest.NewRecorder()
	suite.controller.GetLongURL(w, req)

	err := getErrOrNil(w)
	assert.Equal(suite.T(), domain.Code(usecase.RetrieveFullURLValidation), err.Code(), "Wrong error code. Expected: %d, got: %d", usecase.RetrieveFullURLValidation, err.Code())

}

func (suite *ControllerSuite) TestGetLongURLRequestWhenLongURLNotFound() {

	req := httptest.NewRequest("GET", "http://www.small.ml?shortUrl=http://www.small.ml/nil", nil)
	w := httptest.NewRecorder()
	suite.controller.GetLongURL(w, req)

	err := getErrOrNil(w)
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
