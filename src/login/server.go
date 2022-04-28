package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

var (
	apiKey string
)

func main() {
	initEnv()
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("hello world")
	})
	log.Fatal(app.Listen(":3000"))
}

func initEnv() {
	eErr := godotenv.Load()
	checkErr(eErr, "Error loading .env file")
	apiKey = os.Getenv("REST_API_KEY")
	fmt.Println(apiKey)
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatal(msg)
	}
}
