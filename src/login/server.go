package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/joho/godotenv"
)

var (
	apiKey      string
	baseUrl     string
	redirectUrl string
)

func main() {
	initEnv()
	getAuthorize()
	getToken()

	// Load html
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Load static file like CSS, Images, JS ...
	app.Static("/public", "./public")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title":    "Kakao login REST API",
			"SubTitle": "Start kakao REST API call!!",
		})
	})

	log.Fatal(app.Listen(":3000"))
}

func callKakaoLogin() {
	fmt.Println("Called")
}

func getToken() {
}

func getAuthorize() {
	url := fmt.Sprintf("%s/oauth/authorize?client_id=%s&redirect_uri=%s&response_type=code", baseUrl, apiKey, redirectUrl)
	//url := fmt.Sprintf("%s/oauth/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=account_emai", baseUrl, apiKey, redirectUrl)
	res, err := http.Get(url)
	checkErr(err, "Error Token check code plz!!!")
	fmt.Println(res)
}

func initEnv() {
	eErr := godotenv.Load()
	checkErr(eErr, "Error loading .env file")
	apiKey = os.Getenv("REST_API_KEY")
	baseUrl = os.Getenv("BASE_URL")
	redirectUrl = os.Getenv("REDIRECT_URI")
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatal(msg)
	}
}
