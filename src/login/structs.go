package main

type KakaoAuthResult struct {
	AccessToken           string `json:"access_token"`
	TokenType             string `json:"token_type"`
	RefreshToken          string `json:"refresh_token"`
	ExpiresIn             uint64 `json:"expires_in"`
	RefreshTokenExpiresIn uint64 `json:"refresh_token_expires_in"`
}

type AccessTokenInfo struct {
	ID        uint64 `json:"id"`
	ExpiresIn int    `json:"expires_in"`
	AppIn     int    `json:"app_id"`
}

type LogoutInfo struct {
	ID uint64 `json:"id"`
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

type RevokeResult struct {
	TargetIDType string  `json:"target_id_type"`
	TargetID     uint64  `json:"target_id"`
	Scopes       []Scope `json:"scopes"`
}
