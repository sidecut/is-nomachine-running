package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault("port", 80)
	viper.SetDefault("sslport", 443)
}

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

func main() {
	log.Println(debug.ReadBuildInfo())
	log.Println("###")

	t := &Template{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}

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

	viper.AutomaticEnv()
	viper.SetEnvPrefix("isno")
	port := viper.GetInt("port")
	sslport := viper.GetInt("sslport")

	// Start port 443
	go func(c *echo.Echo) {
		e.Logger.Fatal(e.StartAutoTLS(fmt.Sprintf(":%v", sslport)))
	}(e)

	// Start port 80
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", port)))
}

// Template struct
type Template struct {
	templates *template.Template
}

// Render is used by Echo handler funcs
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
