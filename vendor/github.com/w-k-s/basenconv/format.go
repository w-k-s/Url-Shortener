package basenconv

import (
	"bytes"
)

const Base62Alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
const Base16Alphabet = "0123456789ABCDEF"
const Base8Alphabet = "01234567"
const Base2Alphabet = "01"

// FormatBase62 returns the string representation of value in the base 62,
func FormatBase62(value uint64) string {
	return FormatUint(value, Base62Alphabet)
}

// FormatBase62 returns the string representation of value in the base 16,
func FormatHex(value uint64) string {
	return FormatUint(value, Base16Alphabet)
}

// FormatBase62 returns the string representation of value in the base 8
func FormatOctal(value uint64) string {
	return FormatUint(value, Base8Alphabet)
}

// FormatBase62 returns the string representation of value in the base 2
func FormatBinary(value uint64) string {
	return FormatUint(value, Base2Alphabet)
}

// FormatUint returns the string representation of i in base n,
// The base is determined using the length of the given alphabet
func FormatUint(value uint64, alphabet string) string {
	base := uint64(len(alphabet))
	quotient := value / base
	remainder := value % base

	var buffer bytes.Buffer
	buffer.WriteByte(alphabet[remainder])

	for quotient >= base {
		oldQuotient := quotient
		quotient = oldQuotient / base
		remainder = oldQuotient % base

		buffer.WriteByte(alphabet[remainder])
	}

	if quotient > 0 {
		buffer.WriteByte(alphabet[quotient])
	}

	return reverse(buffer.String())
}
