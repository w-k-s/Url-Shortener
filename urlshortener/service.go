package urlshortener

import(
	"gopkg.in/mgo.v2"
	"net/url"
	"gopkg.in/mgo.v2/bson"
	a "github.com/waqqas-abdulkareem/short-url/app"
	"github.com/waqqas-abdulkareem/short-url/db"
	"errors"
	"math/rand"
	"time"
	"strconv"
)

type urlRecord struct{
	LongUrl string `json:"longUrl" bson:"longUrl"`

	//should be uint64 but bson does not support such large numbers, so use string instead
	ShortId string `json:"-" bson:"shortId"`
}

type Service struct{
	app *a.App
}

func NewService(app *a.App) *Service{
	return &Service{
		app,
	}
}

func (s *Service) urlsColl() *mgo.Collection{
	return s.app.UrlsColl()
}

func (s *Service) ShortenUrl(host string, longUrl *url.URL) (*url.URL,error){

	var urlRecords []urlRecord
	err := s.urlsColl().Find(bson.M{db.DocNameLongUrl: longUrl.String()}).All(&urlRecords)
	
	if err != nil{
		return nil,err
	}

	if len(urlRecords) == 1 {
		s.app.Logger.Println("Record found")
		record := urlRecords[0]
		return s.buildShortenedUrl(longUrl,host,record.ShortId),nil
	}

	maxTries := 3
	inserted := false
	src := rand.NewSource(time.Now().UnixNano())
	random := rand.New(src)

	urlRec := urlRecord{
		LongUrl: longUrl.String(),
		ShortId: strconv.FormatUint(random.Uint64(),10),
	}

	for try := 0; try < maxTries; try++{
		
		err = s.urlsColl().Insert(urlRec)
		if err == nil{
			inserted = true
			break
		}
		if mgo.IsDup(err){
			s.app.Logger.Println("Duplication Error")
			urlRec.ShortId = strconv.FormatUint(random.Uint64(),10)
			continue
		}else{
			s.app.Logger.Println("Insert error",err.Error())
			return nil,err
		}
	}

	if !inserted{
		return nil,errors.New("Could not save url after several attempts")
	}

	return s.buildShortenedUrl(longUrl,host,urlRec.ShortId),nil
}

func (s *Service) buildShortenedUrl(original *url.URL, host string, shortId string) *url.URL{
	
	s.app.Logger.Printf("Short Id: %v\n",shortId)

	shortUrl,_ := url.Parse(original.String())
	shortUrl.Host = host
	shortUrlVals := shortUrl.Query()
	shortUrlVals.Add("id",shortId)
	shortUrl.RawQuery = shortUrlVals.Encode()

	return shortUrl
}