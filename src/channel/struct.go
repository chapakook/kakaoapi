package main

import "time"

type KakaoAuthResult struct {
	AccessToken           string `json:"access_token"`
	TokenType             string `json:"token_type"`
	RefreshToken          string `json:"refresh_token"`
	ExpiresIn             int    `json:"expires_in"`
	RefreshTokenExpiresIn uint64 `json:"refresh_token_expires_in"`
	Scope                 string `json:"scope"`
}

type Channels struct {
	ChannelUuid     string    `json:"channel_uuid"`
	ChannelPublicID string    `json:"channel_public_id"`
	Relation        string    `json:"relation"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type key struct {
	UserID   uint64     `json:"user_id"`
	Channels []Channels `json:"channels"`
}
