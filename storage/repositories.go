package storage

import "AlexSarva/gofermart/models"

//var ErrDuplicatePK = errors.New("duplicate PK")

type Repo interface {
	Ping() bool
	NewUser(user *models.User) error
}
