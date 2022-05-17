package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/joho/godotenv"
)

var (
	BASE_URL            string
	REST_API_CLIENT_KEY string
	REDIRECT_URI        string

	token kakaoTokenResult
)

type kakaoTokenResult struct {
	AccessToken           string `json:"access_toekn"`
	TokenType             string `json:"token_type"`
	RefreshToken          string `json:"refresh_token"`
	ExpiresIn             uint64 `json:"expires_in"`
	RefreshTokenExpiresIn uint64 `json:"refresh_token_expires_in"`
}

func main() {
	// load dotenv
	err := godotenv.Load()
	if err != nil {
		panic("dotenv load failed")
	}

	// set up enviromental variable
	BASE_URL = os.Getenv("BASE_URL")
	REST_API_CLIENT_KEY = os.Getenv("REST_API_CLIENT_KEY")
	REDIRECT_URI = os.Getenv("REDIRECT_URI")

	// fiber setting html
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// server static resources
	app.Static("/public", "./public")

	// auth controller
	app.Get("/oauth", func(c *fiber.Ctx) error {
		code := c.Query("code")
		fmt.Println("code: ", code)

		// retrieve token from user code
		resp, err := http.PostForm(BASE_URL+"/oauth/token", url.Values{
			"grant_type":   []string{"authorization_code"},
			"client_id":    []string{os.Getenv("REST_API_CLIENT_KEY")},
			"redirect_uri": []string{os.Getenv("BASE_URL")},
			"code":         []string{code},
		})
		if err != nil {
			panic("retrieve token failed")
		}

		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(&token); err == nil {
			fmt.Println(token)
			return c.SendString("Get oauth")
		}

		return c.SendString("Get oauth")
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title":               "카카오 REST API 예제",
			"SubTitle":            "예제입니다",
			"REST_API_CLIENT_KEY": REST_API_CLIENT_KEY,
			"REDIRECT_URI":        REDIRECT_URI,
			"BASE_URL":            BASE_URL,
		})
	})

	log.Fatal(app.Listen(":3000"))
}
