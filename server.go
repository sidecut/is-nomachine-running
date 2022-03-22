package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault("port", 80)
	viper.SetDefault("sslport", 443)
}

// Hello just outputs "Hello, world" using an HTML template
func Hello(c *fiber.Ctx) error {
	return c.Render("hello", fiber.Map{"Hello": "World"})
}

func statusAPI(c *fiber.Ctx) error {
	status, err := getStatus()
	if err != nil {
		// TODO: log this error
		return err
	}
	return c.Format(status)
}

func main() {
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{Views: engine})
	app.Use(compress.New())

	app.Use(cors.New())

	viper.AutomaticEnv()
	viper.SetEnvPrefix("isno")
	port := viper.GetInt("port")
	// sslport := viper.GetInt("sslport")

	app.Use(logger.New())

	app.Get("/hello", Hello)
	app.Get("/api", statusAPI)
	app.Static("/", "dist")

	// Start port 443
	// go func(c *viper.Viper) {
	// 	logger.Fatal(app.StartAutoTLS(fmt.Sprintf(":%v", sslport)))
	// }(app)

	// Start port 80
	log.Fatal(app.Listen(fmt.Sprintf(":%v", port)))
}

// // Template struct
// type Template struct {
// 	templates *template.Template
// }

// // Render is used by Echo handler funcs
// func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
// 	return t.templates.ExecuteTemplate(w, name, data)
// }
