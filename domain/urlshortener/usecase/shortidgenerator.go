package urlshortener

import (
	"github.com/w-k-s/basenconv"
	"math"
	"math/rand"
	"time"
)

// Gaussian Distribution is used to generate random numbers
// Random numbers are biased to around 10000 which generates a short id in base64
// This is to keep shortIds far apart and usually short.
//
// If Deviation is greater than or equal to bias, the bias number is more liely to occur
// If deviation is less than half the bias, the bias number is unlikely to occur
type ShortIDLength int

const bias int = 10000
const (
	VERY_SHORT ShortIDLength = 2 * 10000 // bias * 3
	SHORT      ShortIDLength = 10000     // bias
	MEDIUM     ShortIDLength = 10000 / 4 // bias/4
	VERY_LONG  ShortIDLength = 1
)

type ShortIDGenerator interface {
	Generate(idLength ShortIDLength) string
}

type DefaultShortIDGenerator struct{}

func (gen DefaultShortIDGenerator) Generate(idLength ShortIDLength) string {
	biasedRandom := uint64(randBias(0, 1<<31-1, bias, float64(idLength)))
	return basenconv.FormatBase62(biasedRandom)
}

func randBias(min, max, bias int, deviation float64) int {

	influence := rand.Intn(101)
	x := randInRange(min, max)

	//this is the part that moves the random number closer to the bell
	//not sure how it works.
	if x > bias {
		return x + int(math.Floor(gauss(influence, deviation)*float64(bias-x)))
	}

	return x - int(math.Floor(gauss(influence, deviation)*float64(x-bias)))
}

func gauss(_x int, deviation float64) float64 {
	a := 1.0       //height of the bell curve, arbitrary value
	b := 50.0      //position of the center of the bell, arbitrary value
	c := deviation //width of the bell; wider bell means value near bias is less frequent
	x := float64(_x)

	exp := (-1 * (x - b) * (x - b)) / (2 * c * c)
	return a * math.Exp(exp)
}

func randInRange(min, max int) int {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	return random.Intn(max-min) + min
}
