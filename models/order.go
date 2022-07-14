package models

import (
	"github.com/google/uuid"
	"time"
)

type Order struct {
	UserID   uuid.UUID   `json:"user_id,omitempty" db:"user_id"`
	OrderNum string      `json:"number" db:"order_num"`
	Status   string      `json:"status" db:"status"`
	Accrual  interface{} `json:"accrual,omitempty" db:"accrual"`
	Created  time.Time   `json:"uploaded_at" db:"created"`
}

type OrderDB struct {
	OrderNum string      `json:"number" db:"order_num"`
	Status   string      `json:"status" db:"status"`
	Accrual  interface{} `json:"accrual,omitempty" db:"accrual"`
	Created  time.Time   `json:"uploaded_at" db:"created"`
}

type Balance struct {
	Current  interface{} `json:"current" db:"current"`
	Withdraw interface{} `json:"withdraw" db:"withdraw"`
}
