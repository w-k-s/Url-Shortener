package app


import(
	"gopkg.in/mgo.v2"
	"github.com/waqqas-abdulkareem/short-url/db"
)

func (a *App) ensureIndexes(){
	index := mgo.Index{
		Key: []string{db.DocNameShortId},
		Unique: true,//only allow unique url-ids
		DropDups: false,//raise error if url-id is not unique
		Background: false,//other connections cant use collection while index is under construction
		Sparse: true,//if document is missing url-id, do not index it
	}

	err := a.UrlsColl().EnsureIndex(index)
	if err != nil{
		panic(err)
	}
}

func (a *App) DB() *mgo.Database{
	return a.Session.DB(db.Name)
}

func (a *App) Coll(name string) *mgo.Collection{
	return a.DB().C(name)
}

func (a *App) UrlsColl() *mgo.Collection{
	return a.Coll(db.CollNameUrls)
}