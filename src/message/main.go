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

	message := app.Group("/message", Message)

	// default
	basic := message.Group("default", Basic)
	basic.Get("/sendtome", SendToMe)
	basic.Get("/sendtofriends", SendToFriends)

	// default User
	basicuser := message.Group("/scrap", BasicUser)
	basicuser.Get("/sendtome", SendToMe)
	basicuser.Get("/sendtofriends", SendToFriends)

	// Scrap
	scrap := message.Group("/scrap", Scrap)
	scrap.Get("/sendtome", SendToMe)
	scrap.Get("/sendtofriends", SendToFriends)

	// Scrap User
	scrapuser := message.Group("/scrapuser", ScrapUser)
	scrapuser.Get("/sendtome", SendToMe)
	scrapuser.Get("/sendtofriends", SendToFriends)

	log.Fatal(app.Listen(PORT))
}
