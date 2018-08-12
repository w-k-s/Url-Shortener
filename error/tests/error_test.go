package error

import (
	err "github.com/w-k-s/short-url/error"
	"reflect"
	"testing"
)

const code err.Code = 10000
const domain string = "domain"
const message string = "message"

func TestNewError(t *testing.T) {

	fields := map[string]string{"key": "value"}

	err := err.NewError(
		code,
		domain,
		message,
		fields,
	)

	if err.Code() != code {
		t.Errorf("error.Code, got: %d, want: %d.", err.Code(), code)
	}

	if err.Domain() != domain {
		t.Errorf("error.Domain, got: %s, want: %s.", err.Domain(), domain)
	}

	if err.Error() != message {
		t.Errorf("error.Message, got: %s, want: %s.", err.Error(), message)
	}

	if !reflect.DeepEqual(fields, err.Fields()) {
		t.Errorf("error.Fields, got: %v, want: %v.", err.Fields(), fields)
	}
}

func TestTypeAssertionToError(t *testing.T) {

	fields := map[string]string{"key": "value"}

	var anError interface{} = err.NewError(
		code,
		domain,
		message,
		fields,
	)

	if _, ok := anError.(error); !ok {
		t.Errorf("err.Error doesn't comply to `error` interface")
	}
}

func TestTypeAssertionToErr(t *testing.T) {

	fields := map[string]string{"key": "value"}

	var anError interface{} = err.NewError(
		code,
		domain,
		message,
		fields,
	)

	if _, ok := anError.(err.Err); !ok {
		t.Errorf("err.Error doesn't comply to `err.Err` interface")
	}

}
