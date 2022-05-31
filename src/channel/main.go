package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

func main() {
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/public", "./public")

	app.Get("/", Index)
	app.Get("/oauth", OAuth)
	app.Get("/channel", Channel)
	app.Get("/check", Check)

	log.Fatal(app.Listen(PORT))
}
