package domain

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

const code Code = 10000
const domain string = "domain"
const message string = "message"

func TestNewError(t *testing.T) {

	fields := map[string]string{"key": "value"}

	err := NewError(
		code,
		domain,
		message,
		fields,
	)

	assert.Equal(t, err.Code(), code, "error.Code, got: %d, want: %d.", err.Code(), code)
	assert.Equal(t, err.Domain(), domain, "error.Domain, got: %s, want: %s.", err.Domain(), domain)
	assert.Equal(t, err.Error(), message, "error.Message, got: %s, want: %s.", err.Error(), message)
	assert.True(t, reflect.DeepEqual(fields, err.Fields()), "error.Fields, got: %v, want: %v.", err.Fields(), fields)
}
