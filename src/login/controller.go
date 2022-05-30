package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Index(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title":               "Kakao REST API Example - Login",
		"SubTitle":            "Step 1 - Get Token!",
		"REST_API_CLIENT_KEY": REST_API_CLIENT_KEY,
		"REDIRECT_URI":        REDIRECT_URI,
		"BASE_URL":            BASE_URL,
	})
}

func OAuth(c *fiber.Ctx) error {
	code := c.Query("code")

	// retrieve token from user code
	resp, err := http.PostForm(BASE_URL+"/oauth/token", url.Values{
		"grant_type":   []string{"authorization_code"},
		"client_id":    []string{REST_API_CLIENT_KEY},
		"redirect_uri": []string{REDIRECT_URI},
		"code":         []string{code},
	})
	CheckErr(err)
	CheckStatus(resp)

	var auth KakaoAuthResult
	err = json.NewDecoder(resp.Body).Decode(&auth)
	CheckErr(err)

	// add cookies
	cookie := new(fiber.Cookie)
	cookie.Name = "accesstoken"
	cookie.Value = auth.AccessToken
	cookie.Expires = time.Now().Add(time.Duration(auth.ExpiresIn) * time.Second)
	c.Cookie(cookie)

	cookie.Name = "refreshtoken"
	cookie.Value = auth.RefreshToken
	cookie.Expires = time.Now().Add(time.Duration(auth.RefreshTokenExpiresIn) * time.Second)
	c.Cookie(cookie)

	// print info's in console
	fmt.Println("[+] Auth Info")
	fmt.Println("	AccessToken : ", auth.AccessToken)
	fmt.Println("	ExpiresIn : ", auth.ExpiresIn)
	fmt.Println("	RefreshToken : ", auth.RefreshToken)
	fmt.Println("	RefreshTokenExpiresIn : ", auth.RefreshTokenExpiresIn)
	fmt.Println("	TokenType : ", auth.TokenType)

	return c.Render("oauth", fiber.Map{
		"Title":        "Get Token Success!!",
		"SubTitle":     "Step 2 - Get Token Info",
		"AuthResult":   auth,
		"TokenType":    auth.TokenType,
		"AccessToken":  auth.AccessToken,
		"ExpiresIn":    auth.ExpiresIn,
		"RefreshToken": auth.RefreshToken,
	})
}

func Info(c *fiber.Ctx) error {
	req, err := http.NewRequest("GET", BASE_API_URL+"/v1/user/access_token_info", nil)
	CheckErr(err)
	req.Header.Add("Authorization", "Bearer "+c.Cookies("accesstoken"))

	client := http.Client{}
	resp, err := client.Do(req)
	CheckErr(err)
	CheckStatus(resp)

	var tokenInfo AccessTokenInfo
	err = json.NewDecoder(resp.Body).Decode(&tokenInfo)
	CheckErr(err)

	// print info's in console
	fmt.Println("[+] Token Info")
	fmt.Println("	ID : ", tokenInfo.ID)
	fmt.Println("	AppIn : ", tokenInfo.AppIn)
	fmt.Println("	ExpiresIn : ", tokenInfo.ExpiresIn)

	return c.Render("info", fiber.Map{
		"Title":     "Info Token Success!!",
		"SubTitle":  "Step 3 - Refresh Token",
		"ID":        tokenInfo.ID,
		"AppIn":     tokenInfo.AppIn,
		"ExpiresIn": tokenInfo.ExpiresIn,
	})
}

func Refresh(c *fiber.Ctx) error {
	// retrieve token from user code
	resp, err := http.PostForm(BASE_URL+"/oauth/token", url.Values{
		"grant_type":    []string{"refresh_token"},
		"client_id":     []string{REST_API_CLIENT_KEY},
		"refresh_token": []string{c.Cookies("refreshtoken")},
	})
	CheckErr(err)
	CheckStatus(resp)

	var auth KakaoAuthResult
	err = json.NewDecoder(resp.Body).Decode(&auth)
	CheckErr(err)

	if auth.RefreshToken == "" {
		auth.RefreshToken = c.Cookies("refreshtoken")
	}

	// add cookies
	cookie := new(fiber.Cookie)
	cookie.Name = "accesstoken"
	cookie.Value = auth.AccessToken
	cookie.Expires = time.Now().Add(24 * time.Hour)
	c.Cookie(cookie)

	cookie.Name = "refreshtoken"
	cookie.Value = auth.RefreshToken
	cookie.Expires = time.Now().Add(24 * time.Hour)
	c.Cookie(cookie)

	// print info's in console
	fmt.Println("[+] Refresh Auth Info")
	fmt.Println("	TokenType : ", auth.TokenType)
	fmt.Println("	AccessToken : ", auth.AccessToken)
	fmt.Println("	ExpiresIn : ", auth.ExpiresIn)
	fmt.Println("	RefreshToken : ", auth.RefreshToken)
	fmt.Println("	RefreshTokenExpiresIn : ", auth.RefreshTokenExpiresIn)

	return c.Render("refresh", fiber.Map{
		"Title":                 "Refresh Token Success!!",
		"SubTitle":              "Step 4 - Get Scopes",
		"TokenType":             auth.TokenType,
		"AccessToken":           auth.AccessToken,
		"ExpiresIn":             auth.ExpiresIn,
		"RefreshToken":          auth.RefreshToken,
		"RefreshTokenExpiresIn": auth.RefreshTokenExpiresIn,
	})
}

func Scopes(c *fiber.Ctx) error {
	req, err := http.NewRequest("GET", BASE_API_URL+"/v2/user/scopes", nil)
	CheckErr(err)
	req.Header.Add("Authorization", "Bearer "+c.Cookies("accesstoken"))

	client := http.Client{}
	resp, err := client.Do(req)
	CheckErr(err)
	CheckStatus(resp)

	var scopes ScopeResult
	err = json.NewDecoder(resp.Body).Decode(&scopes)
	CheckErr(err)

	// print info's in console
	fmt.Println("[+] Scope Result")
	fmt.Println("	ID : ", scopes.ID)

	for i, s := range scopes.Scopes {
		fmt.Println("	Scopes - ", i)
		fmt.Println("	  - Agreed : ", s.Agreed)
		fmt.Println("	  - DisplayName : ", s.DisplayName)
		fmt.Println("	  - ID : ", s.ID)
		fmt.Println("	  - Revocable : ", s.Revocable)
		fmt.Println("	  - Type : ", s.Type)
		fmt.Println("	  - Using : ", s.Using)
	}

	return c.Render("scopes", fiber.Map{
		"Title":    "Get Scopes Success!!",
		"SubTitle": "Step 5 - Logout",
		"Scopes":   scopes.Scopes,
	})
}

func Logout(c *fiber.Ctx) error {
	req, err := http.NewRequest("POST", BASE_API_URL+"/v1/user/logout", nil)
	CheckErr(err)
	req.Header.Add("Authorization", "Bearer "+c.Cookies("accesstoken"))

	client := http.Client{}
	resp, err := client.Do(req)
	CheckErr(err)
	CheckStatus(resp)

	var logoutinfo LogoutInfo
	err = json.NewDecoder(resp.Body).Decode(&logoutinfo)
	CheckErr(err)

	// print info's in console
	fmt.Println("[+] Logout")
	fmt.Println("	ID : ", logoutinfo.ID)

	// cookies clear
	c.ClearCookie()

	return c.Redirect("http://localhost:3000/end?id=" + strconv.FormatUint(logoutinfo.ID, 10))
}

func End(c *fiber.Ctx) error {
	return c.Render("end", fiber.Map{
		"Title": "Logout Success",
		"ID":    c.Query("id"),
	})
}
