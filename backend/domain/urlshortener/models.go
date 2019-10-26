package urlshortener

import (
	"time"
)

type URLRecord struct {
	LongURL    string    `bson:"longUrl"`
	ShortID    string    `bson:"shortId"`
	CreateTime time.Time `bson:"createTime"`
}

type URLRepository interface {
	SaveRecord(record *URLRecord) (*URLRecord, error)
	LongURL(shortID string) (*URLRecord, error)
	ShortURL(longURL string) (*URLRecord, error)
}
