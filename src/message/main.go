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

	// Server static resources
	app.Static("/public", "./public")

	// auth controller
	app.Get("/", Index)
	app.Get("/oauth", OAuth)
	app.Get("/scopes", Scopes)
	app.Get("/message", Message)
	app.Get("/memo", Memo)
	app.Get("/logout", Logout)

	log.Fatal(app.Listen(PORT))
}
