package storage

import "AlexSarva/gofermart/models"

type Repo interface {
	Ping() bool
	NewUser(user *models.User) error
	GetUser(username string) (*models.User, error)
	CheckOrder(orderNum int) (*models.Order, error)
	NewOrder(order *models.Order) error
}
