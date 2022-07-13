package models

import "github.com/google/uuid"

type Status struct {
	Result string `json:"result"`
}

type User struct {
	ID       uuid.UUID `json:"id" db:"id"`
	Username string    `json:"login" db:"username"`
	Password string    `json:"password" db:"passwd"`
	Cookie   string    `json:"cookie" db:"cookie"`
}
