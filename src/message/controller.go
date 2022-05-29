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
	CheckErr(err)

	var auth KakaoAuthResult
	err = json.NewDecoder(resp.Body).Decode(&auth)
	CheckErr(err)

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
	cookie.Expires = time.Now().Add(time.Duration(auth.ExpiresIn) * time.Second)
	c.Cookie(cookie)

	cookie.Name = "refreshtoken"
	cookie.Value = auth.RefreshToken
	cookie.Expires = time.Now().Add(time.Duration(auth.RefreshTokenExpiresIn) * time.Second)
	c.Cookie(cookie)

	return c.Redirect("http://localhost:3000/message")
}

func Message(c *fiber.Ctx) error {
	return c.Render("message", fiber.Map{
		"Title":    "Sucess Login",
		"SubTitle": "Send To Me",
	})
}

func Send(c *fiber.Ctx) error {
	client := http.Client{}
	req, err := http.NewRequest("POST", BASE_API_URL+"/v2/api/talk/memo/default/send", nil)
	CheckErr(err)
	resp, err := client.Do(req)
	CheckErr(err)

	var sendResult SendResult
	err = json.NewDecoder(resp.Body).Decode(&sendResult)
	CheckErr(err)

	// print info's in console
	fmt.Println("[+] Default Send Result")
	fmt.Println("	Result Code : ", sendResult.ResultCode)

	return c.SendString("send")
}

func Logout(c *fiber.Ctx) error {
	client := http.Client{}

	// Logout
	req, err := http.NewRequest("POST", BASE_API_URL+"/v1/user/logout", nil)
	CheckErr(err)

	req.Header.Add("Authorization", "Bearer "+c.Cookies("accesstoken"))
	_, err = client.Do(req)
	CheckErr(err)

	// Unlink
	req, err = http.NewRequest("POST", BASE_API_URL+"/v1/user/unlink", nil)
	CheckErr(err)
	req.Header.Add("Authorization", "Bearer "+c.Cookies("accesstoken"))
	_, err = client.Do(req)
	CheckErr(err)

	return c.Redirect("http://localhost:3000")
}
