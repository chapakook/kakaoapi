package main

import (
	"bytes"
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

	return c.Redirect("http://localhost:3000/channel")
}

func Channel(c *fiber.Ctx) error {
	return c.Render("channel", fiber.Map{
		"Title":    "Success Login",
		"Subtitle": "Step 1 - Check Kakao Talk Channel Relationship",
	})
}

func Check(c *fiber.Ctx) error {
	params := url.Values{}
	params.Add("channel_public_ids", CHANNEL_PUBLIC_ID)

	req, err := http.NewRequest("GET", BASE_API_URL+"/v1/api/talk/channels", bytes.NewBufferString(params.Encode()))
	CheckErr(err)
	req.Header.Add("Authorization", "Bearer "+c.Cookies("accesstoken"))

	client := http.Client{}
	resp, err := client.Do(req)
	CheckErr(err)
	CheckStatus(resp)

	var key key
	err = json.NewDecoder(resp.Body).Decode(&key)
	CheckErr(err)

	// print info's in console
	fmt.Println("[+] Channel Relationship Info")
	fmt.Println("	UserID : ", key.UserID)
	for i, c := range key.Channels {
		fmt.Println("	Channels - ", i)
		fmt.Println("	  - Channel uuid : ", c.ChannelUuid)
		fmt.Println("	  - Channel Public ID : ", c.ChannelPublicID)
		fmt.Println("	  - Relation : ", c.Relation)
		fmt.Println("	  - Updated At : ", c.UpdatedAt)
	}

	return c.SendString("check")
}
