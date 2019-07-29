package urlshortener

import (
	database "github.com/w-k-s/short-url/db"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

const collNameURLs = "urls"

const fieldShortID = "shortId"
const fieldLongURL = "longUrl"

type URLRecord struct {
	LongURL    string    `bson:"longUrl"`
	ShortID    string    `bson:"shortId"`
	CreateTime time.Time `bson:"createTime"`
}

type VisitTrack struct {
	ShortID    string    `bson:"shortId"`
	IPAddress  string    `bson:"visitIp"`
	CreateTime time.Time `bson:"createTime"`
}

type URLRepository struct {
	db     *database.Db
	logger *log.Logger
}

func NewURLRepository(db *database.Db, logger *log.Logger) *URLRepository {
	return &URLRepository{
		db:     db,
		logger: logger,
	}
}

func (ur *URLRepository) urlCollection() *mgo.Collection {
	return ur.db.Instance().C(collNameURLs)
}

func (ur *URLRepository) updateIndexes() error {
	index := mgo.Index{
		Key:        []string{fieldShortID},
		Unique:     true,  //only allow unique url-ids
		DropDups:   false, //raise error if url-id is not unique
		Background: false, //other connections cant use collection while index is under construction
		Sparse:     true,  //if document is missing url-id, do not index it
	}

	return ur.urlCollection().EnsureIndex(index)
}

func (ur *URLRepository) SaveRecord(record *URLRecord) (*URLRecord, error) {
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

func (ur *URLRepository) LongURL(shortID string) (*URLRecord, error) {
	var record URLRecord
	err := ur.urlCollection().
		Find(bson.M{fieldShortID: shortID}).
		One(&record)

	if err != nil {
		panicIfConnectionError(err)
		return nil, err
	}

	return &record, nil
}

func (ur *URLRepository) ShortURL(longURL string) (*URLRecord, error) {
	var record URLRecord
	err := ur.urlCollection().
		Find(bson.M{fieldLongURL: longURL}).
		One(&record)

	if err != nil {
		panicIfConnectionError(err)
		return nil, err
	}

	return &record, nil
}

func (ur *URLRepository) logLastError(err error) {
	if lastError, ok := err.(*mgo.LastError); ok {
		ur.logger.Printf("Last Error. Code: %d, Message: %s (rows affected: %d)\n", lastError.Code, lastError.Err, lastError.N)
	}
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
