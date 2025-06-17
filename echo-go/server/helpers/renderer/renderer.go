package renderer

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

type Renderer interface {
	Render(w io.Writer, name string, data interface{}, c echo.Context) error
}

type renderer struct {
	template *template.Template
	debug    bool
	location string
}

func NewRenderer(location string, debug bool) Renderer {
	t := &renderer{
		debug:    debug,
		location: location,
	}
	t.reloadTemplates()
	return t
}

func (r *renderer) reloadTemplates() {
	r.template = template.Must(template.ParseGlob(r.location))
}

func (r *renderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if r.debug {
		r.reloadTemplates()
	}
	return r.template.ExecuteTemplate(w, name, data)
}
