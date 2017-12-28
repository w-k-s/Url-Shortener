package home

import (
	"github.com/gorilla/mux"
	"html/template"
)

var c *Controller

func init() {
	tpl := template.Must(template.ParseGlob("./templates/*"))
	c = NewController(tpl)
}

func ConfigureRoutes(r *mux.Router) {
	r.HandleFunc("/", c.Index).
		Methods("GET")
}
