package urlshortener

type URLRepository interface{
	SaveRecord(record *URLRecord) (*URLRecord, error)
	LongURL(shortID string) (*URLRecord, error)
	ShortURL(longURL string) (*URLRecord, error)
}