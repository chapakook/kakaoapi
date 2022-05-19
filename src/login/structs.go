package main

type KakaoAuthResult struct {
	AccessToken           string `json:"access_token"`
	TokenType             string `json:"token_type"`
	RefreshToken          string `json:"refresh_token"`
	ExpiresIn             uint64 `json:"expires_in"`
	RefreshTokenExpiresIn uint64 `json:"refresh_token_expires_in"`
	Scope                 string `json:"scope"`
}

type AccessTokenInfo struct {
	ID         			  uint64 `json:"id"`
	ExpiresIn  			  int    `json:"expires_in"`
	AppIn      			  int    `json:"app_id"`
}
