package home

import (
	a "github.com/w-k-s/short-url/app"
	"net/http"
)

type Controller struct{}

func NewController(app *a.App) *Controller {
	return &Controller{}
}

func (c *Controller) Index(w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req, "/public/docs", http.StatusSeeOther)
}
