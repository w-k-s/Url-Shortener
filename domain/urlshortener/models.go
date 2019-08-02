package urlshortener

import(
	"time"
	"gopkg.in/mgo.v2/bson"
)

type URLRecord struct {
	LongURL    string    `bson:"longUrl"`
	ShortID    string    `bson:"shortId"`
	CreateTime time.Time `bson:"createTime"`
}