package app

import (
	"AlexSarva/gofermart/storage"
	"AlexSarva/gofermart/storage/storagepg"
	"errors"
	"fmt"
)

// Database interface for different types of databases
type Database struct {
	Repo storage.Repo
}

// NewStorage generate new instance of database
func NewStorage(database string) (*Database, error) {
	if len(database) > 0 {
		Storage := storagepg.NewPostgresDBConnection(database)
		fmt.Println("Using PostgreSQL Database")
		return &Database{
			Repo: Storage,
		}, nil
	}

	return &Database{}, errors.New("u must use database config")

}
