package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
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
	t := &Template{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}

	e := echo.New()
	e.Renderer = t
	e.Use(middleware.Gzip())
	e.Use(middleware.Logger())

	corsConfig := middleware.CORSConfig{AllowOrigins: []string{"*"}}
	e.Use(middleware.CORSWithConfig(corsConfig))

	e.GET("/hello", Hello)
	e.GET("/api", statusAPI)
	e.Static("/", "dist")

	viper.AutomaticEnv()
	viper.SetEnvPrefix("isno")
	port := viper.GetInt("port")
	sslport := viper.GetInt("sslport")

	oldLevel := e.Logger.Level()
	e.Logger.SetLevel(log.INFO)
	e.Logger.Infoj(log.JSON{"port": port, "sslport": sslport})
	e.Logger.SetLevel(oldLevel)

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
