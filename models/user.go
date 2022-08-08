package models

import (
	"time"

	"github.com/google/uuid"
)

type Status struct {
	Result string `json:"result"`
}

type User struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Username  string    `json:"login" db:"username"`
	Password  string    `json:"password" db:"passwd"`
	Cookie    string    `json:"cookie" db:"cookie"`
	CookieExp time.Time `json:"cookie_expires" db:"cookie_expires"`
}

type TestUser struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Username  string    `json:"login" db:"username"`
	Password  string    `json:"password" db:"passwd"`
	Cookie    string    `json:"cookie" db:"cookie"`
	CookieExp time.Time `json:"cookie_expires" db:"cookie_expires"`
	Token     string    `json:"token"`
}

type Token struct {
	TokenType   string `json:"token_type"`
	AuthToken   string `json:"auth_token"`
	GeneratedAt string `json:"generated_at"`
	ExpiresAt   string `json:"expires_at"`
}
