package loyality

import (
	"AlexSarva/gofermart/internal/app"
	"AlexSarva/gofermart/models"
	"AlexSarva/gofermart/storage/storagepg"
	"log"
)

func GetOrdersToProcessing(database app.Database, ordersCh chan string) {
	for {
		loyal, loyalErr := database.Repo.GetOrdersForProcessing()
		if loyalErr != nil {
			if loyalErr == storagepg.ErrNoValues {
				log.Println(loyalErr)
			}
			log.Println(loyalErr)
		}
		for _, order := range loyal {
			ordersCh <- order
		}
	}

}

func GetProcessedInfo(client *ProcessingClient, ordersCh chan string, procesedCh chan models.ProcessingOrder) {
	for order := range ordersCh {
		if orderInfo, orderInfoErr := client.GetOrder(order); orderInfoErr == nil {
			procesedCh <- orderInfo
		}
	}
}

func ApplyLoyality(database app.Database, procesedCh chan models.ProcessingOrder) {
	for order := range procesedCh {
		database.Repo.UpdateOrder(order)
	}
}
