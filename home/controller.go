package home

import (
	a "github.com/w-k-s/short-url/app"
	"html/template"
	"net/http"
)

type Controller struct {
	tpl *template.Template
}

func NewController(app *a.App) *Controller {
	tpl := template.Must(template.ParseGlob("home/templates/*.html"))
	return &Controller{
		tpl,
	}
}

func (c *Controller) Index(w http.ResponseWriter, req *http.Request) {
	err := c.tpl.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
