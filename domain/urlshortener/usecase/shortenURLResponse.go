package usecase

type ShortenURLResponse struct {
	LongURL  string `json:"longUrl"`
	ShortURL string `json:"shortUrl"`
}