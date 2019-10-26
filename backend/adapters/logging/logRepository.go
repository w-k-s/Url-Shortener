package logging

import (
	"bytes"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type logRecord struct {
	Time      time.Time `bson:"createTime"`
	Method    string    `bson:"method"`
	URI       string    `bson:"requestUri"`
	IPAddress string    `bson:"ipAddress"`
	Status    int       `bson:"status"`
	Body      string    `bson:"body"`
}

func (lr logRecord) String() string {
	return fmt.Sprintf("%s: %s %s %s - %d",
		lr.IPAddress,
		lr.Method,
		lr.URI,
		lr.Body,
		lr.Status,
	)
}

type LogRepository struct {
	db     *sql.DB
	logger *log.Logger
}

func NewLogRepository(db *sql.DB, logger *log.Logger) *LogRepository {
	return &LogRepository{
		db:     db,
		logger: logger,
	}
}

func (lr *LogRepository) LogRequest(r *http.Request) *logRecord {
	return &logRecord{
		Time:      time.Now(),
		Method:    r.Method,
		URI:       r.RequestURI,
		IPAddress: r.Header.Get("X-Forwarded-For"),
		Body:      readRequestBody(r),
	}
}

func (lr *LogRepository) LogResponse(sw *StatusWriter, record *logRecord) error {
	record.Status = sw.Status()
	lr.logger.Println(record.String())

	_, err := lr.db.Exec(
		`INSERT INTO logs (method,uri,ip_address,status,body,create_time) VALUES ($1,$2,$3,$4,$5,$6)`,
		record.Method,
		record.URI,
		record.IPAddress,
		record.Status,
		record.Body,
		record.Time,
	)

	return err
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
