package usecase

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func TestGenerator(t *testing.T) {

	for i := 0; i < 10; i++ {

		gen := DefaultShortIDGenerator{}
		shortIDs := []string{
			gen.Generate(VeryShort),
			gen.Generate(Short),
			gen.Generate(Medium),
			gen.Generate(VeryLong),
		}

		compareLengths := func(i, j int) bool {
			left := shortIDs[i]
			right := shortIDs[j]
			return len(right) > len(left)
		}

		sorted := sort.SliceIsSorted(shortIDs, compareLengths)

		assert.True(t, sorted, "Expected VERY_SHORT, SHORT, MEDIUM, VERY_LONG. Got %v", shortIDs)
	}
}
