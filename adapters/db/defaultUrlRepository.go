package db

import (
	u "github.com/w-k-s/short-url/domain/urlshortener"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

const collNameURLs = "urls"

const fieldShortID = "shortId"
const fieldLongURL = "longUrl"

type DefaultURLRepository struct {
	db     *Db
	logger *log.Logger
}

func NewURLRepository(db *Db, logger *log.Logger) *DefaultURLRepository {
	return &DefaultURLRepository{
		db:     db,
		logger: logger,
	}
}

func (ur *DefaultURLRepository) urlCollection() *mgo.Collection {
	return ur.db.Instance().C(collNameURLs)
}

func (ur *DefaultURLRepository) updateIndexes() error {
	index := mgo.Index{
		Key:        []string{fieldShortID},
		Unique:     true,  //only allow unique url-ids
		DropDups:   false, //raise error if url-id is not unique
		Background: false, //other connections cant use collection while index is under construction
		Sparse:     true,  //if document is missing url-id, do not index it
	}

	return ur.urlCollection().EnsureIndex(index)
}

func (ur *DefaultURLRepository) SaveRecord(record *u.URLRecord) (*u.URLRecord, error) {
	err := ur.urlCollection().
		Insert(record)

	if err != nil {
		panicIfConnectionError(err)
		return nil, err
	}

	err = ur.updateIndexes()
	if err != nil {
		log.Panic(err)
	}

	return record, nil
}

func (ur *DefaultURLRepository) LongURL(shortID string) (*u.URLRecord, error) {
	var record u.URLRecord
	err := ur.urlCollection().
		Find(bson.M{fieldShortID: shortID}).
		One(&record)

	if err != nil {
		panicIfConnectionError(err)
		return nil, err
	}

	return &record, nil
}

func (ur *DefaultURLRepository) ShortURL(longURL string) (*u.URLRecord, error) {
	var record u.URLRecord
	err := ur.urlCollection().
		Find(bson.M{fieldLongURL: longURL}).
		One(&record)

	if err != nil {
		panicIfConnectionError(err)
		return nil, err
	}

	return &record, nil
}

func isConnectionError(err error) bool {
	otherError := mgo.IsDup(err) ||
		err == mgo.ErrNotFound
	return !otherError
}

func panicIfConnectionError(err error) {
	if isConnectionError(err) {
		log.Fatal(err)
	}
}
