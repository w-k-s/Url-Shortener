package home

import (
	a "github.com/w-k-s/short-url/app"
	"net/http"
)

type Controller struct {
	app *a.App
}

func NewController(app *a.App) *Controller {
	return &Controller{
		app,
	}
}

func (c *Controller) Index(w http.ResponseWriter, req *http.Request) {

	err := c.app.Templates.ExecuteTemplate(w, "home.gohtml", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
