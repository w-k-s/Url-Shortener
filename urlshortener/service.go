package urlshortener

import (
	"errors"
	a "github.com/w-k-s/short-url/app"
	"github.com/w-k-s/short-url/db"
	"github.com/w-k-s/basenconv"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"math/rand"
	"net/url"
	"time"
)

type urlRecord struct {
	LongUrl string `json:"longUrl" bson:"longUrl"`
	ShortId string `json:"-" bson:"shortId"`
	CreateTime time.Time `json:"-" bson:"createTime"`
}

type Service struct {
	app *a.App
}

func NewService(app *a.App) *Service {
	return &Service{
		app,
	}
}

func (s *Service) urlsColl() *mgo.Collection {
	return s.app.UrlsColl()
}

func (s *Service) ShortenUrl(host string, longUrl *url.URL) (*url.URL, error) {

	var urlRecords []urlRecord
	err := s.urlsColl().Find(bson.M{db.UrlsFieldLongUrl: longUrl.String()}).All(&urlRecords)

	if err != nil {
		return nil, err
	}

	if len(urlRecords) == 1 {
		s.app.Logger.Println("Record found")
		record := urlRecords[0]
		return s.buildShortenedUrl(longUrl, host, record), nil
	}

	maxTries := 3
	inserted := false
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	urlRec := urlRecord{
		LongUrl: longUrl.String(),
		CreateTime: time.Now(),
	}

	for try := 0; try < maxTries; try++ {

		shortIdNum := s.generateShortIdNumber(try,random)
		urlRec.ShortId = basenconv.FormatBase62(shortIdNum)

		err = s.urlsColl().Insert(urlRec)
		if err == nil {
			inserted = true
			break
		}
		if mgo.IsDup(err) {
			s.app.Logger.Println("Duplication Error")
			continue
		} else {
			s.app.Logger.Println("Insert error", err.Error())
			return nil, err
		}
	}

	if !inserted {
		return nil, errors.New("Could not save url after several attempts")
	}

	return s.buildShortenedUrl(longUrl, host, urlRec), nil
}

func (s *Service) generateShortIdNumber(try int,random *rand.Rand) uint64{
	//31 should be extracted as a configuration, probably
	//still, not the best solution, sometimes the shortId will be short, othertimes long
	return uint64(random.Intn(1<<31 - 1))
}

func (s *Service) buildShortenedUrl(original *url.URL, host string, urlRecord urlRecord) *url.URL {

	shortUrl, _ := url.Parse(original.String())
	shortUrl.Host = host
	shortUrl.Path = urlRecord.ShortId

	return shortUrl
}
