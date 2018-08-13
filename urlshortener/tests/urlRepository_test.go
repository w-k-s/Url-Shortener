package tests

// import (
// 	_ "fmt"
// 	database "github.com/w-k-s/short-url/db"
// 	repo "github.com/w-k-s/short-url/urlshortener"
// 	"gopkg.in/mgo.v2"
// 	"gopkg.in/mgo.v2/bson"
// 	"os"
// 	"testing"
// 	"time"
// )

// func TestMain(m *testing.M) {
// 	setup()
// 	retCode := m.Run()
// 	tearDown()
// 	os.Exit(retCode)
// }

// var db *database.Db
// var urlRepo *repo.URLRepository
// var record *repo.URLRecord

// func setup() {
// 	db = database.New("mongodb://localhost:27017/shorturl_test")
// 	urlRepo = repo.NewURLRepository(db)

// 	record = &repo.URLRecord{
// 		"http://www.example.com",
// 		"shrt",
// 		time.Now(),
// 	}
// }

// func tearDown() {
// 	db.Instance().
// 		C("urls").
// 		RemoveAll(bson.M{})

// 	defer db.Close()
// }

// func TestSaveRecordSucccessful(t *testing.T) {

// 	if _, err := urlRepo.SaveRecord(record); err != nil {
// 		t.Errorf("Expected: save record. Got: %s", err)
// 	}

// }

// func TestDuplicateRecordFails(t *testing.T) {

// 	urlRepo.SaveRecord(record)

// 	if _, err := urlRepo.SaveRecord(record); !mgo.IsDup(err) {
// 		t.Errorf("Expected: duplication error. Got: %s", err)
// 	}
// }

// func TestFindExistingShortURL(t *testing.T) {

// 	if result, err := urlRepo.ShortURL(record.LongUrl); result == nil || result.ShortId != record.ShortId {
// 		t.Errorf("Expected Matching ShortId '%s'. Got: '%v' (error: '%s')", record.ShortId, result, err)
// 	}

// }

// func TestFindAbsentShortURL(t *testing.T) {

// 	if result, err := urlRepo.ShortURL("http://www.nil.com"); err == nil {
// 		t.Errorf("Expected err when shortId not found. Got: nil. (record: %v)", result)
// 	}

// }

// func TestFindExistingLongURL(t *testing.T) {

// 	if result, err := urlRepo.LongURL(record.ShortId); result == nil || result.LongUrl != record.LongUrl {
// 		t.Errorf("Expected Matching LongUrl '%s'. Got: '%v' (error: '%s')", record.LongUrl, result, err)
// 	}

// }

// func TestFindAbsentLongURL(t *testing.T) {

// 	if result, err := urlRepo.LongURL("nil"); err == nil {
// 		t.Errorf("Expected err when longUrl not found. Got: nil. (record: %v)", result)
// 	}

// }
