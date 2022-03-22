package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/pflag"
)

var port *int

// Hello just outputs "Hello, world" using an HTML template
func Hello(c echo.Context) error {
	return c.Render(http.StatusOK, "hello", "World")
}

func statusAPI(c echo.Context) error {
	status, err := getStatus()
	if err != nil {
		// TODO: log this error
		return err
	}
	return c.JSON(http.StatusOK, status)
}

func init() {
	port = pflag.IntP("port", "p", 3123, "Specify port to listen on")
}

func main() {
	t := &Template{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}

	pflag.Parse()

	e := echo.New()
	e.Renderer = t
	e.Use(middleware.Gzip())
	e.Use(middleware.Logger())
	e.Use(middleware.RequestID())

	corsConfig := middleware.CORSConfig{AllowOrigins: []string{"*"}}
	e.Use(middleware.CORSWithConfig(corsConfig))

	e.GET("/hello", Hello)
	e.GET("/api", statusAPI)
	e.Static("/", "dist")

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", *port)))
}

// Template struct
type Template struct {
	templates *template.Template
}

// Render is used by Echo handler funcs
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
