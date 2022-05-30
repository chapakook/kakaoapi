package main

type KakaoAuthResult struct {
	AccessToken           string `json:"access_token"`
	TokenType             string `json:"token_type"`
	RefreshToken          string `json:"refresh_token"`
	ExpiresIn             int    `json:"expires_in"`
	RefreshTokenExpiresIn uint64 `json:"refresh_token_expires_in"`
	Scope                 string `json:"scope"`
}

type MemoSendResult struct {
	ResultCode int `json:"result_code"`
}

type FailureInfo struct {
	Code          int      `json:"code"`
	Msg           string   `json:"msg"`
	ReceiverUuids []string `json:"receiver_uuids"`
}
type FriendsSendResult struct {
	SuccessfulReceiverUuids []string      `json:"successful_receiver_uuids"`
	FailureInfos            []FailureInfo `json:"failure_info"`
}
type Link struct {
	WebUrl string `json:"web_url"`
	//MobileUrl string `json:"mobile_web_url"`
	//Android   string `json:"android_execution_params"`
	//IOS       string `json:"ios_execution_params"`
}

type TextMessageTemplate struct {
	ObjectType string `json:"object_type"`
	Text       string `json:"text"`
	Link       Link   `json:"link"`
}

type Scope struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	Type        string `json:"type"`
	Using       bool   `json:"using"`
	Agreed      bool   `json:"agreed"`
	Revocable   bool   `json:"revocable"`
}

type ScopeResult struct {
	ID     uint64  `json:"id"`
	Scopes []Scope `json:"scopes"`
}

type LogoutInfo struct {
	ID uint64 `json:"id"`
}
