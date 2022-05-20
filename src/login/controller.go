package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/gofiber/fiber/v2"
)

func Index(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title":               "Kakao REST API example",
		"SubTitle":            "Step 1 - Get Token!",
		"REST_API_CLIENT_KEY": REST_API_CLIENT_KEY,
		"REDIRECT_URI":        REDIRECT_URI,
		"BASE_URL":            BASE_URL,
	})
}

func OAuth(c *fiber.Ctx) error {
	// front-end -(query)-> back-end
	code := c.Query("code")

	// retrieve token from user code
	resp, err := http.PostForm(BASE_URL+"/oauth/token", url.Values{
		"grant_type":   []string{"authorization_code"},
		"client_id":    []string{REST_API_CLIENT_KEY},
		"redirect_uri": []string{REDIRECT_URI},
		"code":         []string{code},
	})
	if err != nil {
		panic(err)
	}

	// decode Response JSON
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

	return c.Render("oauth", fiber.Map{
		"Title":        "Step 1 - Get Token Success!!",
		"SubTitle":     "Step 2 - Get Token Info",
		"AuthResult":   auth,
		"TokenType":    auth.TokenType,
		"AccessToken":  auth.AccessToken,
		"ExpiresIn":    auth.ExpiresIn,
		"RefreshToken": auth.RefreshToken,
		"Scope":        auth.Scope,
	})
}

func Info(c *fiber.Ctx) error {
	client := http.Client{}
	req, err := http.NewRequest("GET", BASE_API_URL+"/v1/user/access_token_info", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", "Bearer "+c.Query("accesstoken"))
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	var tokenInfo AccessTokenInfo
	err = json.NewDecoder(resp.Body).Decode(&tokenInfo)
	if err != nil {
		panic(err)
	}

	// print info's in console
	fmt.Println("[+] Token Info")
	fmt.Println("	ID : ", tokenInfo.ID)
	fmt.Println("	AppIn : ", tokenInfo.AppIn)
	fmt.Println("	ExpiresIn : ", tokenInfo.ExpiresIn)

	return c.Render("info", fiber.Map{
		"Title":        "Step 2 - Get Token Success!!",
		"SubTitle":     "Step 3 - Get Refresh Token",
		"ID":           tokenInfo.ID,
		"AppIn":        tokenInfo.AppIn,
		"ExpiresIn":    tokenInfo.ExpiresIn,
		"AccessToken":  c.Query("accesstoken"),
		"RefreshToken": c.Query("refreshtoken"),
	})
}

func Refresh(c *fiber.Ctx) error {
	// retrieve token from user code
	resp, err := http.PostForm(BASE_URL+"/oauth/token", url.Values{
		"grant_type":    []string{"refresh_token"},
		"client_id":     []string{REST_API_CLIENT_KEY},
		"refresh_token": []string{c.Query("refreshtoken")},
	})
	if err != nil {
		panic(err)
	}

	// decode Response JSON
	var auth KakaoAuthResult
	err = json.NewDecoder(resp.Body).Decode(&auth)
	if err != nil {
		panic(err)
	}

	// print info's in console
	fmt.Println("[+] Refresh Auth Info")
	fmt.Println("	TokenType : ", auth.TokenType)
	fmt.Println("	AccessToken : ", auth.AccessToken)
	fmt.Println("	ExpiresIn : ", auth.ExpiresIn)
	fmt.Println("	RefreshToken : ", auth.RefreshToken)
	fmt.Println("	RefreshTokenExpiresIn : ", auth.RefreshTokenExpiresIn)

	return c.Render("refresh", fiber.Map{
		"Title":                 "Step 3 - Refresh Token Success!!",
		"SubTitle":              "Step 4 - Logout",
		"TokenType":             auth.TokenType,
		"AccessToken":           auth.AccessToken,
		"ExpiresIn":             auth.ExpiresIn,
		"RefreshToken":          auth.RefreshToken,
		"RefreshTokenExpiresIn": auth.RefreshTokenExpiresIn,
		"BASE_URL":              BASE_URL,
		"REST_API_CLIENT_KEY":   REST_API_CLIENT_KEY,
		"LOGOUT_REDIRECT_URI":   LOGOUT_REDIRECT_URI,
	})
}

func Logout(c *fiber.Ctx) error {
	fmt.Println(c.Query("accesstoken"))
	client := http.Client{}
	req, err := http.NewRequest("POST", BASE_API_URL+"/v1/user/logout", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", "Bearer "+c.Query("accesstoken"))
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	var id LogoutInfo
	err = json.NewDecoder(resp.Body).Decode(&id)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[+] Logout")
	fmt.Println("ID: ", id)

	return c.SendString("success")
}
