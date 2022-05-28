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

	// Groups
	talk := app.Group("/talk", Talk)
	memo := app.Group("/memo", Memo)

	// Server static resources
	app.Static("/public", "./public")
	talk.Static("/public", "./public")
	memo.Static("/public", "./public")

	// auth controller
	app.Get("/", Index)
	app.Get("/oauth", OAuth)
	app.Get("/logout", Logout)

	log.Fatal(app.Listen(PORT))
}
