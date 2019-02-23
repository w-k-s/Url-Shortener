package logging

import (
	"bytes"
	"fmt"
	database "github.com/w-k-s/short-url/db"
	"gopkg.in/mgo.v2"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const collNameLogs = "logs"

const fieldShortId = "shortId"
const fieldLongURL = "longUrl"

type logRecord struct {
	Time      time.Time `bson:"createTime"`
	Method    string    `bson:"method"`
	URI       string    `bson:"requestUri"`
	IpAddress string    `bson:"ipAddress"`
	Status    int       `bson:"status"`
	Body      string    `bson:"body"`
}

func (lr logRecord) String() string {
	return fmt.Sprintf("%s: %s %s %s - %d",
		lr.IpAddress,
		lr.Method,
		lr.URI,
		lr.Body,
		lr.Status,
	)
}

type LogRepository struct {
	db     *database.Db
	logger *log.Logger
}

func NewLogRepository(logger *log.Logger, db *database.Db) *LogRepository {

	createLogsCollectionIfNotExists(db.Instance())

	return &LogRepository{
		db:     db,
		logger: logger,
	}
}

func createLogsCollectionIfNotExists(db *mgo.Database) error {
	if exists := logCollectionExists(db); !exists {
		return createLogsCollection(db)
	}
	return nil
}

func logCollectionExists(db *mgo.Database) bool {

	names, err := db.CollectionNames()
	if err != nil {
		panic(err.Error)
	}

	for _, name := range names {
		if name == collNameLogs {
			return true
		}
	}

	return false
}

func createLogsCollection(db *mgo.Database) error {

	coll := &mgo.Collection{
		Database: db,
		Name:     collNameLogs,
		FullName: fmt.Sprintf("%s.%s", db.Name, collNameLogs),
	}

	return coll.Create(&mgo.CollectionInfo{
		Capped:   true,
		MaxBytes: 5000000, //5MB
	})
}

func (lr *LogRepository) logsCollection() *mgo.Collection {
	return lr.db.Instance().C(collNameLogs)
}

func (lr *LogRepository) LogRequest(r *http.Request) *logRecord {

	return &logRecord{
		Time:      time.Now(),
		Method:    r.Method,
		URI:       r.RequestURI,
		IpAddress: r.RemoteAddr,
		Body:      readRequestBody(r),
	}
}

func (lr *LogRepository) LogResponse(sw *StatusWriter, record *logRecord) error {
	record.Status = sw.Status()
	lr.logger.Println(record.String())
	return lr.logsCollection().Insert(record)
}

func readRequestBody(r *http.Request) string {

	var bodyBytes []byte
	var err error

	if r.Body != nil {
		bodyBytes, err = ioutil.ReadAll(r.Body)
		if err != nil {
			return ""
		}
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	return string(bodyBytes)
}
