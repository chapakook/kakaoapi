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

	token      kakaoTokenResult
	token_info kakaoTokenInfo
)

type kakaoTokenResult struct {
	AccessToken           string `json:"access_token"`
	TokenType             string `json:"token_type"`
	RefreshToken          string `json:"refresh_token"`
	ExpiresIn             uint64 `json:"expires_in"`
	RefreshTokenExpiresIn uint64 `json:"refresh_token_expires_in"`
	Scope                 string `json:"scope"`
}

type kakaoTokenInfo struct {
	ID         uint64 `json:"id"`
	Expires_in int    `json:"expires_in"`
	App_in     int    `json:"app_id"`
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
		// front-end -(query)-> back-end
		code := c.Query("code")
		fmt.Println("code: ", code)

		// retrieve token from user code
		resp, rErr := http.PostForm(BASE_URL+"/oauth/token", url.Values{
			"grant_type":   []string{"authorization_code"},
			"client_id":    []string{REST_API_CLIENT_KEY},
			"redirect_uri": []string{REDIRECT_URI},
			"code":         []string{code},
		})
		if rErr != nil {
			panic("retrieve token failed")
		}
		dErr := json.NewDecoder(resp.Body).Decode(&token)
		if dErr != nil {
			panic("resp decode failed")
		}
		fmt.Println("token: ", token)

		// check token info (https://kapi.kakao.com/v1/user/access_token_info)
		// req, rErr := http.NewRequest("GET", "https://kapi.kakao.com/v1/user/access_token_info", nil)
		// if rErr != nil {
		// 	log.Fatal("make check token request failed")
		// 	return c.SendString("make check token request failed")
		// }
		// req.Header.Add("Authorization", "Bearer "+token.AccessToken)
		// client := &http.Client{}
		// cResp, cErr := client.Do(req)
		// if cErr != nil {
		// 	log.Fatal("check token info failed")
		// 	return c.SendString("check token info failed")
		// }
		// defer cResp.Body.Close()
		// fmt.Println(token_info)

		return c.Render("one", fiber.Map{
			"Title":    "Step 1 - Get Token Success!!",
			"SubTitle": "Step 2 - Get Token Info",
			"Token":    token,
		})
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title":               "Kakao REST API example",
			"SubTitle":            "Step 1 - Get Token!",
			"REST_API_CLIENT_KEY": REST_API_CLIENT_KEY,
			"REDIRECT_URI":        REDIRECT_URI,
			"BASE_URL":            BASE_URL,
		})
	})

	log.Fatal(app.Listen(":3000"))
}
