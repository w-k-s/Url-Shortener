package urlshortener

import (
	"errors"
	"github.com/w-k-s/basenconv"
	a "github.com/w-k-s/short-url/app"
	"github.com/w-k-s/short-url/db"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"math/rand"
	"net/url"
	"time"
)

type urlRecord struct {
	LongUrl    string    `json:"longUrl" bson:"longUrl"`
	ShortId    string    `json:"-" bson:"shortId"`
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

func (s *Service) ShortenUrl(reqUrl *url.URL, longUrl *url.URL) (*url.URL, error) {

	var urlRecords []urlRecord
	err := s.urlsColl().Find(bson.M{db.UrlsFieldLongUrl: longUrl.String()}).All(&urlRecords)

	if err != nil {
		return nil, err
	}

	if len(urlRecords) == 1 {
		s.app.Logger.Println("Record found")
		record := urlRecords[0]
		return s.buildShortenedUrl(reqUrl, record), nil
	}

	maxTries := 3
	inserted := false
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	urlRec := urlRecord{
		LongUrl:    longUrl.String(),
		CreateTime: time.Now(),
	}

	for try := 0; try < maxTries; try++ {

		shortIdNum := s.generateShortIdNumber(try, random)
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

	return s.buildShortenedUrl(reqUrl, urlRec), nil
}

func (s *Service) generateShortIdNumber(try int, random *rand.Rand) uint64 {
	//31 should be extracted as a configuration, probably
	//still, not the best solution, sometimes the shortId will be short, othertimes long
	return uint64(random.Intn(1<<31 - 1))
}

func (s *Service) buildShortenedUrl(reqUrl *url.URL, urlRecord urlRecord) *url.URL {
	s.app.Logger.Println("reqUrl.Host = ", reqUrl.String())

	return &url.URL{
		Scheme: reqUrl.Scheme,
		Host:   reqUrl.Host,
		Path:   urlRecord.ShortId,
	}
}

func (s *Service) GetLongUrl(shortUrl *url.URL) (*url.URL, bool, error) {

	path := shortUrl.Path
	if len(path) == 0 {
		return nil, false, errors.New("expected url to have a path")
	}

	if path[0] == '/' {
		path = path[1:]
	}

	var urlRecords []urlRecord
	err := s.urlsColl().Find(bson.M{db.UrlsFieldShortId: path}).
		All(&urlRecords)

	if err != nil {
		return nil, false, err
	}

	if len(urlRecords) == 0 {
		return nil, false, nil
	}

	longUrl, err := url.Parse(urlRecords[0].LongUrl)
	if err != nil {
		return nil, false, err
	}

	return longUrl, true, nil
}
