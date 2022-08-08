package models

import (
	"time"

	"github.com/google/uuid"
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
	Current  float32 `json:"current" db:"current"`
	Withdraw float32 `json:"withdrawn" db:"withdraw"`
}

type Withdraw struct {
	UserID   uuid.UUID `json:"user_id,omitempty" db:"user_id"`
	OrderNum string    `json:"order" db:"order_num"`
	Withdraw float32   `json:"sum" db:"withdraw"`
	Created  time.Time `json:"uploaded_at" db:"created"`
}

type WithdrawBD struct {
	OrderNum string    `json:"order" db:"order_num"`
	Withdraw float32   `json:"sum" db:"withdraw"`
	Created  time.Time `json:"processed_at" db:"created"`
}

type ProcessingOrder struct {
	OrderNum string      `json:"order"`
	Status   string      `json:"status"`
	Accrual  interface{} `json:"accrual,omitempty"`
}

type TestType struct {
	Name       string `json:"name"`
	Diameter   string `json:"diameter"`
	Population string `json:"population"`
}

type MyChans struct {
	InsertOrdersCh chan Order
}
