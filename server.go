package main

import (
	"html/template"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Hello just outputs "Hello, world" using an HTML template
func Hello(c *gin.Context) {
	c.HTML(http.StatusOK, "hello", "World")
}

func statusAPI(c *gin.Context) {
	status, err := getStatus()
	if err != nil {
		// TODO: log this error
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, status)
}

func main() {
	r := gin.Default()
	// r.Use(middleware.Gzip())

	// corsConfig := middleware.CORSConfig{AllowOrigins: []string{"*"}}
	// r.Use(middleware.CORSWithConfig(corsConfig))

	r.GET("/hello", Hello)
	r.GET("/api", statusAPI)
	r.Static("/", "dist")

	// Start port 80
	r.Run()
}

// Template struct
type Template struct {
	templates *template.Template
}

// Render is used by Echo handler funcs
func (t *Template) Render(w io.Writer, name string, data interface{}, c *gin.Context) {
	t.templates.ExecuteTemplate(w, name, data)
}
