package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

var (
	apiKey       string
	baseUrl      string
	redirectUrl  string
	clientID     string
	clientSecret string
	OAuthConf    *oauth2.Config
)

func main() {
	initEnv()
	initOAuth()

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

	app.Get("/oauth/callback", func(c *fiber.Ctx) error {
		return c.Render("kakao_login", nil)
	})

	api := app.Group("/api")

	v1 := api.Group("/v1")
	v1.Get("/login", func(c *fiber.Ctx) error {
		kakaoLogin()
		return c.SendString("Get /api/v1/login")
	})

	log.Fatal(app.Listen(":3000"))
}

func initOAuth() {
	OAuthConf = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectUrl,
	}
}

func kakaoLogin() {
	getOAuthorize()
	getToken()
}

func getToken() error {
	fmt.Println("getToken")
	return nil
}

func getOAuthorize() error {
	url := fmt.Sprintf("%s/oauth/authorize?client_id=%s&redirect_uri=%s&response_type=code", baseUrl, apiKey, redirectUrl)
	//url := fmt.Sprintf("%s/oauth/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=account_emai", baseUrl, apiKey, redirectUrl)
	res, err := http.Get(url)
	//fmt.Println(url)
	//fmt.Println(res)
	if err != nil {
		return err
	}

	data, dErr := ioutil.ReadAll(res.Body)
	if dErr != nil {
		return dErr
	}

	file, cErr := os.Create("views/kakao_login.html")
	if cErr != nil {
		return cErr
	}
	_, wErr := file.Write([]byte(data))
	if wErr != nil {
		return wErr
	}
	return nil
}

func initEnv() {
	eErr := godotenv.Load()
	checkErr(eErr, "Error loading .env file")

	apiKey = os.Getenv("REST_API_KEY")
	baseUrl = os.Getenv("BASE_URL")
	redirectUrl = os.Getenv("REDIRECT_URI")
	clientID = os.Getenv("Client_ID")
	clientSecret = os.Getenv("Client_Secret")
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatal(msg)
	}
}
