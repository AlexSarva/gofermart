package models

import (
	"github.com/google/uuid"
	"time"
)

type Order struct {
	UserID   uuid.UUID   `json:"user_id,omitempty" db:"user_id"`
	OrderNum int         `json:"number" db:"order_num"`
	Status   string      `json:"status" db:"status"`
	Accrual  interface{} `json:"accrual,omitempty" db:"accrual"`
	Created  time.Time   `json:"uploaded_at" db:"created"`
}

type OrderDB struct {
	OrderNum int         `json:"number" db:"order_num"`
	Status   string      `json:"status" db:"status"`
	Accrual  interface{} `json:"accrual,omitempty" db:"accrual"`
	Created  time.Time   `json:"uploaded_at" db:"created"`
}
