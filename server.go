package main

import (
	"html/template"

	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
)

// Hello just outputs "Hello, world" using an HTML template
func Hello(ctx iris.Context) {
	ctx.ViewData("name", "World")
	ctx.View("hello.html")
}

func statusAPI(ctx iris.Context) {
	status, err := getStatus()
	if err != nil {
		// TODO: log this error
		ctx.StopWithProblem(500, iris.NewProblem().DetailErr(err))
		return
	}
	ctx.JSON(status)
}

func main() {
	app := iris.New()

	tmpl := iris.HTML("./views", ".html")
	app.RegisterView(tmpl)

	app.Use(iris.Compression)

	// corsConfig := middleware.CORSConfig{AllowOrigins: []string{"*"}}
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})
	app.Use(crs)

	app.Get("/hello", Hello)
	app.Get("/api", statusAPI)
	app.HandleDir("/", iris.Dir("./dist"))

	// TODO: TLS and logging
	// e.Logger.Fatal(e.StartAutoTLS(":443"))
	app.Listen(":80")
}

// Template struct
type Template struct {
	templates *template.Template
}
