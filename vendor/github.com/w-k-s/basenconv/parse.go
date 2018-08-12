package basenconv

import (
	"bytes"
	"math"
)

// ParseInt interprets a string value in base 62
// and returns the corresponding uint64
func ParseBase62(value string) uint64 {
	return ParseUint(value, Base62Alphabet)
}

// ParseInt interprets a string value in base 16
// and returns the corresponding uint64
func ParseHex(value string) uint64 {
	return ParseUint(value, Base16Alphabet)
}

// ParseInt interprets a string value in base 8
// and returns the corresponding uint64
func ParseOctal(value string) uint64 {
	return ParseUint(value, Base8Alphabet)
}

// ParseInt interprets a string value in base 2
// and returns the corresponding uint64
func ParseBinary(value string) uint64 {
	return ParseUint(value, Base2Alphabet)
}

// ParseInt interprets a string s in base n and
// returns the corresponding value uint64.
// base is determined using the lenght of the given alphabet
func ParseUint(value string, alphabet string) uint64 {
	alphabytes := []byte(alphabet)
	base := float64(len(alphabet))

	var decoded float64 = 0
	for i, c := range reverse(value) {
		repr := float64(bytes.IndexRune(alphabytes, c))
		decoded += repr * math.Pow(base, float64(i))
	}
	return uint64(decoded)
}
