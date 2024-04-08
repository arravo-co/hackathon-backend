package routes_v1

import (
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

var t *Template

func init() {
	t = &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func Hello(c echo.Context) error {
	fmt.Printf(c.Request().RequestURI)
	return c.Render(http.StatusOK, "hell", "World")
}
