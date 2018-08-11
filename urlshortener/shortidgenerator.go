package urlshortener

import (
	"github.com/w-k-s/basenconv"
	"math/rand"
	"time"
)

type ShortIDGenerator struct {
	random *rand.Rand
}

func NewShortIDGenerator() ShortIDGenerator {
	return ShortIDGenerator{
		random: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (gen ShortIDGenerator) Generate() string {
	shortIdNum := uint64(gen.random.Intn(1<<31 - 1))
	return basenconv.FormatBase62(shortIdNum)
}
