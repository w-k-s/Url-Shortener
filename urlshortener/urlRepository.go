package urlshortener

import (
	database "github.com/w-k-s/short-url/db"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

const collNameURLs = "urls"
const colNameVisits = "visits"

const fieldShortId = "shortId"
const fieldLongURL = "longUrl"

type URLRecord struct {
	LongURL    string    `bson:"longUrl"`
	ShortId    string    `bson:"shortId"`
	CreateTime time.Time `bson:"createTime"`
}

type VisitTrack struct {
	ShortId    string    `bson:"shortId"`
	IpAddress  string    `bson:"visitIp"`
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

func (ur *URLRepository) visitTrackCollection() *mgo.Collection {
	return ur.db.Instance().C(colNameVisits)
}

func (ur *URLRepository) updateIndexes() error {
	index := mgo.Index{
		Key:        []string{fieldShortId},
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
		ur.logLastError(err)
		return nil, err
	}

	err = ur.updateIndexes()
	if err != nil {
		return nil, err
	}

	return record, nil
}

func (ur *URLRepository) LongURL(shortId string) (*URLRecord, error) {
	var record URLRecord
	err := ur.urlCollection().
		Find(bson.M{fieldShortId: shortId}).
		One(&record)

	if err != nil {
		return nil, err
	}

	return &record, nil
}

func (ur *URLRepository) TrackVisit(visitTrack *VisitTrack) error {
	err := ur.visitTrackCollection().
		Insert(visitTrack)

	if err != nil {
		ur.logLastError(err)
	}

	return err
}

func (ur *URLRepository) ShortURL(longURL string) (*URLRecord, error) {
	var record URLRecord
	err := ur.urlCollection().
		Find(bson.M{fieldLongURL: longURL}).
		One(&record)

	if err != nil {
		return nil, err
	}

	return &record, nil
}

func (ur *URLRepository) logLastError(err error) {
	if lastError, ok := err.(*mgo.LastError); ok {
		ur.logger.Printf("Last Error. Code: %d, Message: %s (rows affected: %d)\n", lastError.Code, lastError.Err, lastError.N)
	}
}
