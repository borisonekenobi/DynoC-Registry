package models

import "github.com/jackc/pgx/v5/pgtype"

type RegisterRequest struct {
	Username pgtype.Text `json:"username"`
	Email    pgtype.Text `json:"email"`
	Password pgtype.Text `json:"password"`
}

type RegisterResponse struct {
	UserID    pgtype.UUID        `json:"user_id"`
	Username  pgtype.Text        `json:"username"`
	Email     pgtype.Text        `json:"email"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
}

type AccountResponse struct {
	UserID    pgtype.UUID        `json:"user_id"`
	Username  pgtype.Text        `json:"username,omitempty"`
	Email     pgtype.Text        `json:"email,omitempty"`
	CreatedAt pgtype.Timestamptz `json:"created_at,omitempty"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at,omitempty"`
}

type UpdateAccountRequest struct {
	Username *pgtype.Text `json:"username,omitempty"`
	Email    *pgtype.Text `json:"email,omitempty"`
	Password *pgtype.Text `json:"password,omitempty"`
}
