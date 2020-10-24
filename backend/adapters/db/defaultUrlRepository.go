package db

import (
	"database/sql"
	"errors"
	"github.com/lib/pq"
	u "github.com/w-k-s/short-url/domain/urlshortener"
)

type DefaultURLRepository struct {
	db *sql.DB
}

func NewURLRepository(db *sql.DB) *DefaultURLRepository {
	return &DefaultURLRepository{
		db: db,
	}
}

func (ur *DefaultURLRepository) SaveRecord(record *u.URLRecord) (*u.URLRecord, error) {
	_, err := ur.db.Exec(
		`INSERT INTO url_records (long_url,short_id) VALUES ($1,$2)`,
		record.LongURL,
		record.ShortID,
	)

	return record, err
}

func (ur *DefaultURLRepository) LongURL(shortID string) (*u.URLRecord, error) {
	rows, err := ur.db.Query("SELECT * FROM url_records WHERE short_id = $1", shortID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, errors.New("Not Found")
	}

	var record u.URLRecord
	if err = rows.Scan(&record.LongURL, &record.ShortID, &record.CreateTime); err != nil {
		return nil, err
	}

	return &record, nil
}

func (ur *DefaultURLRepository) ShortURL(longURL string) (*u.URLRecord, error) {
	rows, err := ur.db.Query("SELECT * FROM url_records WHERE long_url = $1", longURL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, errors.New("Not Found")
	}

	var record u.URLRecord
	if err = rows.Scan(&record.LongURL, &record.ShortID, &record.CreateTime); err != nil {
		return nil, err
	}

	return &record, nil
}

func (ur *DefaultURLRepository) IsDup(err error) bool {
	if pqError, ok := err.(*pq.Error); ok {
		return pqError.Code.Name() == "unique_violation"
	}
	return false
}
