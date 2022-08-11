package storage

import (
	"AlexSarva/gofermart/models"

	"github.com/google/uuid"
)

// Repo primary interface for all types of databases
type Repo interface {
	Ping() bool
	NewUser(user *models.User) error
	GetUser(username string) (*models.User, error)
	CheckOrder(orderNum string) (*models.Order, error)
	NewOrder(order *models.Order) error
	GetOrders(userID uuid.UUID) ([]*models.OrderDB, error)
	GetBalance(userID uuid.UUID) (*models.Balance, error)
	NewWithdraw(withdraw *models.Withdraw) error
	GetAllWithdraw(userID uuid.UUID) ([]*models.WithdrawBD, error)
	GetOrdersForProcessing() ([]string, error)
	UpdateOrder(order models.ProcessingOrder)
}
