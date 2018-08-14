package db

import (
	"github.com/stretchr/testify/assert"
	"github.com/w-k-s/short-url/db"
	"testing"
)

const DB_NAME string = "shorturl_test"
const CONN_STRING string = "mongodb://localhost:27017/" + DB_NAME

func TestNew(t *testing.T) {
	db := db.New(CONN_STRING)
	actualName := db.Instance().Name
	assert.Equal(t, actualName, DB_NAME, "Db name parsed incorrectly, got: %s, want: %s", actualName, DB_NAME)
}
