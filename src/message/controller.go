package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Index(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title":               "Kakao REST API example - Message",
		"SubTitle":            "Step 1 - Go login",
		"BASE_URL":            BASE_URL,
		"REST_API_CLIENT_KEY": REST_API_CLIENT_KEY,
		"REDIRECT_URI":        REDIRECT_URI,
	})
}

func OAuth(c *fiber.Ctx) error {
	code := c.Query("code")

	resp, err := http.PostForm(BASE_URL+"/oauth/token", url.Values{
		"grant_type":   []string{"authorization_code"},
		"client_id":    []string{REST_API_CLIENT_KEY},
		"redirect_uri": []string{REDIRECT_URI},
		"code":         []string{code},
	})
	if err != nil {
		panic(err)
	}

	var auth KakaoAuthResult
	err = json.NewDecoder(resp.Body).Decode(&auth)
	if err != nil {
		panic(err)
	}

	// print info's in console
	fmt.Println("[+] Auth Info")
	fmt.Println("	AccessToken : ", auth.AccessToken)
	fmt.Println("	ExpiresIn : ", auth.ExpiresIn)
	fmt.Println("	RefreshToken : ", auth.RefreshToken)
	fmt.Println("	RefreshTokenExpiresIn : ", auth.RefreshTokenExpiresIn)
	fmt.Println("	Scope : ", auth.Scope)
	fmt.Println("	TokenType : ", auth.TokenType)

	cookie := new(fiber.Cookie)
	cookie.Name = "accesstoken"
	cookie.Value = auth.AccessToken
	cookie.Expires = time.Now().Add(24 * time.Hour)
	c.Cookie(cookie)

	cookie.Name = "refreshtoken"
	cookie.Value = auth.RefreshToken
	cookie.Expires = time.Now().Add(24 * time.Hour)
	c.Cookie(cookie)

	return c.Redirect("http://localhost:3000/message")
}

func Talk(c *fiber.Ctx) error {
	return c.Render("message", fiber.Map{
		"Title":    "Sucess Login",
		"SubTitle": "Send To Me",
	})
}

func Memo(c *fiber.Ctx) error {
	return c.SendString("")
}

func Logout(c *fiber.Ctx) error {
	client := http.Client{}

	// Logout
	req, err := http.NewRequest("POST", BASE_API_URL+"/v1/user/logout", nil)
	if err != nil {
		panic(err)
	}

	req.Header.Add("Authorization", "Bearer "+c.Cookies("accesstoken"))
	_, err = client.Do(req)
	if err != nil {
		panic(err)
	}

	// Unlink
	req, err = http.NewRequest("POST", BASE_API_URL+"/v1/user/unlink", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", "Bearer "+c.Cookies("accesstoken"))
	_, err = client.Do(req)
	if err != nil {
		panic(err)
	}

	fmt.Println("Suceess")
	return c.Redirect("http://localhost:3000")
}
