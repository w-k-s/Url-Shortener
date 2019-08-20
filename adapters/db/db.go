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

func New(connString string, safeMode bool) *Db {
	if len(connString) == 0 {
		panic("Blank database connection string")
	}

	name := connString[strings.LastIndex(connString, "/")+1:]
	session, err := mgo.Dial(connString)
	if safeMode {
		session.SetSafe(&mgo.Safe{
			W:     1,
			FSync: true,
		})
	}
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
