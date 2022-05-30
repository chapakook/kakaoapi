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

	return c.Redirect("http://localhost:3000/scopes")
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

	return c.Redirect("http://localhost:3000/message")
}

func Message(c *fiber.Ctx) error {
	return c.Render("message", fiber.Map{
		"Title":    "Sucess Login",
		"SubTitle": "Step 1 - Send To Me",
	})
}

func Memo(c *fiber.Ctx) error {
	// make text template
	text := c.Query("text")
	template := TextMessageTemplate{
		ObjectType: "text",
		Text:       text,
		Link: Link{
			WebUrl: "https://developers.kakao.com/docs/latest/ko/message/common",
		},
	}
	out, err := json.Marshal(template)
	CheckErr(err)
	params := url.Values{}
	params.Add("template_object", string(out))

	req, err := http.NewRequest("POST", BASE_API_URL+"/v2/api/talk/memo/default/send", bytes.NewBufferString(params.Encode()))
	CheckErr(err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Bearer "+c.Cookies("accesstoken"))

	client := http.Client{}
	resp, err := client.Do(req)
	CheckErr(err)
	CheckStatus(resp)

	var memo MemoSendResult
	err = json.NewDecoder(resp.Body).Decode(&memo)
	CheckErr(err)

	// print info's in console
	fmt.Println("[+] Memo Default Send Result")
	fmt.Println("	Result Code : ", memo.ResultCode)

	return c.Render("memo", fiber.Map{
		"Title":      "Sucess Send To me",
		"ResultCode": memo.ResultCode,
		"Text":       text,
	})
}

func Logout(c *fiber.Ctx) error {
	client := http.Client{}

	req, err := http.NewRequest("POST", BASE_API_URL+"/v1/user/logout", nil)
	CheckErr(err)

	req.Header.Add("Authorization", "Bearer "+c.Cookies("accesstoken"))
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

	return c.Redirect("http://localhost:3000")
}
