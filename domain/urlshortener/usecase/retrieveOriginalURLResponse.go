package usecase

type RetrieveOriginalURLResponse struct {
	LongURL  string `json:"longUrl"`
	ShortURL string `json:"shortUrl"`
}
