package home

import (
	"fmt"
	a "github.com/w-k-s/short-url/app"
	err "github.com/w-k-s/short-url/error"
	"net/http"
)

const (
	//Shortening URL
	HomeRenderFailed err.Code = 20100
)

func domain(e err.Code) string {
	switch e {
	//Shortening URL
	case HomeRenderFailed:
		return "home.renderFailed"
	default:
		return fmt.Sprintf("Unknown Domain (%d)", e)
	}
}

type Controller struct {
	app *a.App
}

func NewController(app *a.App) *Controller {
	return &Controller{
		app,
	}
}

func (c *Controller) Index(w http.ResponseWriter, req *http.Request) {

	data := struct {
		Host string
	}{
		req.Host,
	}

	_err_ := c.app.Templates.ExecuteTemplate(w, "home.gohtml", &data)
	if _err_ != nil {
		err.SendError(w, http.StatusInternalServerError, err.NewError(
			HomeRenderFailed,
			domain(HomeRenderFailed),
			"Failed to render template",
			map[string]string{"error": _err_.Error()},
		))
	}
}
