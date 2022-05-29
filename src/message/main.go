package main

import (
	"fmt"
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

	// Server static resources
	app.Static("/public", "./public")

	// auth controller
	app.Get("/", Index)
	app.Get("/oauth", OAuth)
	app.Get("/message", Message)
	app.Get("/logout", Logout)

	// memo
	app.Use("/talk", func(c *fiber.Ctx) error {
		fmt.Println("talk")
		return c.Next()
	})
	app.Use("/talk/memo", func(c *fiber.Ctx) error {
		fmt.Println("memo")
		return c.Next()
	})
	app.Use("/talk/meme/default", func(c *fiber.Ctx) error {
		fmt.Println("default")
		return c.Next()
	})
	app.Get("/talk/memo/default/send", Send)

	log.Fatal(app.Listen(PORT))
}
