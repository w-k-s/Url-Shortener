package db

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const dbName string = "shorturl_test"
const connString string = "mongodb://localhost:27017/" + dbName

func TestNew(t *testing.T) {
	db := New(connString)
	actualName := Instance().Name
	assert.Equal(t, actualName, dbName, "Db name parsed incorrectly, got: %s, want: %s", actualName, dbName)
}
