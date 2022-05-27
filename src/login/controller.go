package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
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
	client := http.Client{}
	req, err := http.NewRequest("GET", BASE_API_URL+"/v1/user/access_token_info", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", "Bearer "+c.Cookies("accesstoken"))
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
	if err != nil {
		panic(err)
	}

	// decode Response JSON
	var auth KakaoAuthResult
	err = json.NewDecoder(resp.Body).Decode(&auth)
	if err != nil {
		panic(err)
	}

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
	client := http.Client{}
	req, err := http.NewRequest("GET", BASE_API_URL+"/v2/user/scopes", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", "Bearer "+c.Cookies("accesstoken"))
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	var scopes ScopeResult
	err = json.NewDecoder(resp.Body).Decode(&scopes)
	if err != nil {
		panic(err)
	}

	// print info's in console
	fmt.Println("[+] Scope Result")
	fmt.Println("	ID : ", scopes.ID)

	id := ""
	for i, s := range scopes.Scopes {
		fmt.Println("	Scopes - ", i)
		fmt.Println("	  - Agreed : ", s.Agreed)
		fmt.Println("	  - DisplayName : ", s.DisplayName)
		fmt.Println("	  - ID : ", s.ID)
		fmt.Println("	  - Revocable : ", s.Revocable)
		fmt.Println("	  - Type : ", s.Type)
		fmt.Println("	  - Using : ", s.Using)
		id += s.ID + "|"
	}

	return c.Render("scopes", fiber.Map{
		"Title":    "Get Scopes Success!!",
		"SubTitle": "Step 5 - Revoke",
		"Scopes":   scopes.Scopes,
		"ID":       id,
	})
}

func Revoke(c *fiber.Ctx) error {
	// id parsing
	tmp := c.Query("id")
	var id []string
	for tmp != "" {
		n := strings.Index(tmp, "|")
		id = append(id, tmp[:n-1])
		arr := tmp[n+1:]
		tmp = arr
	}

	params := url.Values{}
	for _, i := range id {
		params.Add("scopes", i)
	}

	client := http.Client{}
	req, err := http.NewRequest("POST", BASE_API_URL+"/v2/user/revoke/scopes", bytes.NewBufferString(params.Encode()))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", "Bearer "+c.Cookies("accesstoken"))

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	var revoke RevokeResult
	err = json.NewDecoder(resp.Body).Decode(&revoke)
	if err != nil {
		panic(err)
	}

	// print info's in console
	fmt.Println("[+] Revoke Result")
	fmt.Println("	TargetID : ", revoke.TargetID)
	fmt.Println("	TargetIDType : ", revoke.TargetIDType)

	for i, s := range revoke.Scopes {
		fmt.Println("	Scopes - ", i)
		fmt.Println("	  - Agreed : ", s.Agreed)
		fmt.Println("	  - DisplayName : ", s.DisplayName)
		fmt.Println("	  - ID : ", s.ID)
		fmt.Println("	  - Revocable : ", s.Revocable)
		fmt.Println("	  - Type : ", s.Type)
		fmt.Println("	  - Using : ", s.Using)
	}

	return c.Render("revoke", fiber.Map{
		"Title":        "Revoke Success!!",
		"SubTitle":     "Step 6 - Logout",
		"TargetID":     revoke.TargetID,
		"TargetIDType": revoke.TargetIDType,
		"Scopes":       revoke.Scopes,
	})
}

func Logout(c *fiber.Ctx) error {
	client := http.Client{}

	// logout
	req, err := http.NewRequest("POST", BASE_API_URL+"/v1/user/logout", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", "Bearer "+c.Cookies("accesstoken"))
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	var logoutinfo LogoutInfo
	err = json.NewDecoder(resp.Body).Decode(&logoutinfo)
	if err != nil {
		panic(err)
	}

	// print info's in console
	fmt.Println("[+] Logout")
	fmt.Println("	ID : ", logoutinfo.ID)

	// unlink
	req, err = http.NewRequest("POST", BASE_API_URL+"/v1/user/unlink", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", "Bearer "+c.Cookies("accesstoken"))
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}

	err = json.NewDecoder(resp.Body).Decode(&logoutinfo)
	if err != nil {
		panic(err)
	}

	// cookies clear
	c.ClearCookie()

	// print info's in console
	fmt.Println("[+] Unlink")
	fmt.Println("	ID : ", logoutinfo.ID)

	return c.Redirect("http://localhost:3000/end?ID=" + fmt.Sprint(logoutinfo.ID))
}

func End(c *fiber.Ctx) error {
	return c.Render("end", fiber.Map{
		"Title": "Logout Success",
		"ID":    c.Query("ID"),
	})
}
