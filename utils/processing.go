package utils

import (
	"AlexSarva/gofermart/internal/app"
	"AlexSarva/gofermart/models"
	"log"
)

func InsertOrderToDB(database app.Database, insertOrdersCh chan models.Order) {
	for order := range insertOrdersCh {
		insertErr := database.Repo.NewOrder(&order)
		if insertErr != nil {
			log.Println(insertErr)
			return
		}
	}
}
