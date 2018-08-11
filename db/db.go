package db

import (
	"fmt"
	"gopkg.in/mgo.v2"
)

type Db struct {
	session *mgo.Session
	name    string
}

func New(host string, name string) *Db {
	session, err := mgo.Dial(host)
	if err != nil {
		panic(fmt.Sprintf("Could not connect to datastore with host %s - %v", host, err))
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
