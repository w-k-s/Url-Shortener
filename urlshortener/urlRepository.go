package urlshortener

import (
	database "github.com/w-k-s/short-url/db"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const collNameUrls = "urls"

const fieldShortId = "shortId"
const fieldLongUrl = "longUrl"

type URLRecord struct {
	LongUrl    string    `bson:"longUrl"`
	ShortId    string    `bson:"shortId"`
	CreateTime time.Time `bson:"createTime"`
}

type URLRepository struct {
	db *database.Db
}

func NewURLRepository(db *database.Db) *URLRepository {
	return &URLRepository{
		db: db,
	}
}

func (ur *URLRepository) urlCollection() *mgo.Collection {
	return ur.db.Instance().C(collNameUrls)
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

func (ur *URLRepository) ShortURL(longUrl string) (*URLRecord, error) {
	var record URLRecord
	err := ur.urlCollection().
		Find(bson.M{fieldLongUrl: longUrl}).
		One(&record)

	if err != nil {
		return nil, err
	}

	return &record, nil
}
