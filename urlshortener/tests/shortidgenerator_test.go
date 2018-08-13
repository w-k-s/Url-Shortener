package tests

import (
	u "github.com/w-k-s/short-url/urlshortener"
	"sort"
	"testing"
)

func TestGenerator(t *testing.T) {

	for i := 0; i < 10; i++ {

		gen := u.DefaultShortIDGenerator{}
		shortIds := []string{
			gen.Generate(u.VERY_SHORT),
			gen.Generate(u.SHORT),
			gen.Generate(u.MEDIUM),
			gen.Generate(u.VERY_LONG),
		}

		compareLengths := func(i, j int) bool {
			left := shortIds[i]
			right := shortIds[j]
			return len(right) > len(left)
		}

		sorted := sort.SliceIsSorted(shortIds, compareLengths)

		if !sorted {
			t.Errorf("Expected VERY_SHORT, SHORT, MEDIUM, VERY_LONG. Got %v", shortIds)
		}
	}
}
