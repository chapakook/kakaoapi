package main

type KakaoAuthResult struct {
	AccessToken           string `json:"access_token"`
	TokenType             string `json:"token_type"`
	RefreshToken          string `json:"refresh_token"`
	ExpiresIn             uint64 `json:"expires_in"`
	RefreshTokenExpiresIn uint64 `json:"refresh_token_expires_in"`
	Scope                 string `json:"scope"`
}

type KakaoSendToMeResult struct {
	ResultCode int `json:"result_code"`
}

type Link struct {
	WebUrl    string `json:"web_url"`
	MobileUrl string `json:"mobile_web_url"`
	Android   string `json:"android_execution_params"`
	iOS       string `json:"ios_execution_params"`
}

type TextMessageTemplate struct {
	ObjectType string `json:"object_type"`
	Text       string `json:"text"`
	Link       Link   `json:"link"`
}
