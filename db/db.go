package db

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"strings"
)

type Db struct {
	session *mgo.Session
	name    string
}

func New(connString string) *Db {
	name := connString[strings.LastIndex(connString, "/")+1:]
	session, err := mgo.Dial(connString)
	if err != nil {
		panic(fmt.Sprintf("Could not connect to datastore with host %s - %v", connString, err))
	}

	return &Db{
		session,
		name,
	}
}

func (db *Db) Instance() *mgo.Database {
	return db.session.DB(db.name)
}

func (db *Db) Close() {
	db.session.Close()
}
