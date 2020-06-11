package main

import (
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Hello just outputs "hello, world" using an HTML template
func Hello(c echo.Context) error {
	return c.Render(http.StatusOK, "hello", "World")
}

// Index writes the status of NoMachine
func Index(c echo.Context) error {
	status, err := getStatus()
	if err != nil {
		// TODO: log this error
		return err
	}
	return c.Render(http.StatusOK, "index", status)
}

func main() {
	t := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}

	e := echo.New()
	e.Renderer = t
	e.GET("/hello", Hello)
	e.GET("/", Index)

	// e.GET("/", func(c echo.Context) error {
	// 	return c.String(http.StatusOK, "Hello, World!")
	// })
	e.Logger.Fatal(e.Start(":1323"))
}

// Template struct
type Template struct {
	templates *template.Template
}

// Render is used by Echo handler funcs
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
