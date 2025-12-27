package models

import "github.com/jackc/pgx/v5/pgtype"

type LoginRequest struct {
	Username pgtype.Text `json:"username"`
	Password pgtype.Text `json:"password"`
}

type TokenResponse struct {
	Token     string `json:"token"`
	TokenType string `json:"token_type"`
	ExpiresIn int    `json:"expires_in"`
}
