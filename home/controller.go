package home

import (
	"net/http"
	a "github.com/waqqas-abdulkareem/short-url/app"
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
