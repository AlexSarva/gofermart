package loyality

import (
	"AlexSarva/gofermart/internal/app"
	"AlexSarva/gofermart/models"
	"time"
)

// GetOrdersToProcessing get order from database and put it into the orders channel
func GetOrdersToProcessing(database app.Database, ordersCh chan string) {
	for {
		loyal, _ := database.Repo.GetOrdersForProcessing()
		if len(loyal) > 0 {
			for _, order := range loyal {
				ordersCh <- order
			}
		} else {
			time.Sleep(time.Second * 10)
		}
		time.Sleep(time.Second * 2)
	}
}

// GetProcessedInfo get order from orders channel and check it in loyalty system
// If it processed - put it in processed channel
func GetProcessedInfo(client *ProcessingClient, ordersCh chan string, procesedCh chan models.ProcessingOrder) {
	for order := range ordersCh {
		if orderInfo, orderInfoErr := client.GetOrder(order); orderInfoErr == nil {
			procesedCh <- orderInfo
		}
	}
}

// ApplyLoyalty apply processed order in database
func ApplyLoyalty(database app.Database, procesedCh chan models.ProcessingOrder) {
	for order := range procesedCh {
		database.Repo.UpdateOrder(order)
	}
}
