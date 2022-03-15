package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html"
)

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

	app.Get("/hello", Hello)
	app.Get("/api", statusAPI)
	app.Static("/", "dist")

	// Start port 80
	log.Fatal(app.Listen(":80"))
}
