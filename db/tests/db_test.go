package db

import (
	"github.com/stretchr/testify/assert"
	"github.com/w-k-s/short-url/db"
	"testing"
)

const dbName string = "shorturl_test"
const connString string = "mongodb://localhost:27017/" + dbName

func TestNew(t *testing.T) {
	db := db.New(connString)
	actualName := db.Instance().Name
	assert.Equal(t, actualName, dbName, "Db name parsed incorrectly, got: %s, want: %s", actualName, dbName)
}
