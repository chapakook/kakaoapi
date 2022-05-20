package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

func main() {
	// fiber setting html
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// server static resources
	app.Static("/public", "./public")

	// auth controller
	app.Get("/", Index)
	app.Get("/oauth", OAuth)
	app.Get("/info", Info)
	app.Get("/refresh", Refresh)
	app.Get("/logout", Logout)

	log.Fatal(app.Listen(":3000"))
}
