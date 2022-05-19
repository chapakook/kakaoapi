package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gofiber/fiber/v2"
)

func OAuth(c *fiber.Ctx) error {
	// front-end -(query)-> back-end
	code := c.Query("code")

	// retrieve token from user code
	resp, err := http.PostForm(BASE_URL+"/oauth/token", url.Values{
		"grant_type":   []string{"authorization_code"},
		"client_id":    []string{REST_API_CLIENT_KEY},
		"redirect_uri": []string{REDIRECT_URI},
		"code":         []string{code},
	}); if err != nil {
		panic(err)
	}
	
	// decode Response JSON
	var auth KakaoAuthResult
	err = json.NewDecoder(resp.Body).Decode(&auth); if err != nil {
		panic(err)
	}

	// request access token info
	// TODO: 추후 HttpClient Wrapper 객체 작성 요망
	client := http.Client{}
	req, err := http.NewRequest("GET", BASE_API_URL + "/v1/user/access_token_info", nil)
	if (err != nil){
		panic(err)
	}
	req.Header.Add("Authorization", "Bearer "+ auth.AccessToken)
	resp, err = client.Do(req); if err != nil{
		panic(err)
	}
	
	var tokenInfo AccessTokenInfo
	err = json.NewDecoder(resp.Body).Decode(&tokenInfo); if err != nil{
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

	fmt.Println("[+] Token Info")
	fmt.Println("	ID : ", tokenInfo.ID)	
	fmt.Println("	AppIn : ", tokenInfo.AppIn)
	fmt.Println("	ExpiresIn : ", tokenInfo.ExpiresIn)


	return c.Render("one", fiber.Map{
		"Title":    "Step 1 - Get Token Success!!",
		"SubTitle": "Step 2 - Get Token Info",
		"AuthResult":    auth,
		"TokenInfo" : tokenInfo,
	})
}


func Index(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title":               "Kakao REST API example",
		"SubTitle":            "Step 1 - Get Token!",
		"REST_API_CLIENT_KEY": REST_API_CLIENT_KEY,
		"REDIRECT_URI":        REDIRECT_URI,
		"BASE_URL":            BASE_URL,
	})
}