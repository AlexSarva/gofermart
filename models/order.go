package models

import (
	"github.com/google/uuid"
	"time"
)

type Order struct {
	UserID   uuid.UUID `json:"user_id" db:"user_id"`
	OrderNum int       `json:"number" db:"order_num"`
	Status   string    `json:"status,omitempty" db:"status"`
	Accrual  int       `json:"accrual,omitempty" db:"accrual"`
	Created  time.Time `json:"uploaded_at" db:"created"`
}
