package home

import (
	"html/template"
	"net/http"
)

type Controller struct {
	tpl *template.Template
}

func NewController(tpl *template.Template) *Controller {
	return &Controller{
		tpl,
	}
}

func (c *Controller) Index(w http.ResponseWriter, req *http.Request) {
	err := c.tpl.ExecuteTemplate(w, "home.gohtml", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
