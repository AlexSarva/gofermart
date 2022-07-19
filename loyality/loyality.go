package loyality

import (
	"AlexSarva/gofermart/internal/app"
	"AlexSarva/gofermart/models"
	"log"
	"time"
)

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
	}
}

func GetProcessedInfo(client *ProcessingClient, ordersCh chan string, procesedCh chan models.ProcessingOrder) {
	for order := range ordersCh {
		log.Printf("Отправляем %s", order)
		if orderInfo, orderInfoErr := client.GetOrder(order); orderInfoErr == nil {
			log.Printf("Получаем %s", orderInfo)
			procesedCh <- orderInfo
		}

		procesedCh <- models.ProcessingOrder{
			OrderNum: order,
			Status:   "PROCESSED",
		}
	}
}

func ApplyLoyality(database app.Database, procesedCh chan models.ProcessingOrder) {
	for order := range procesedCh {
		log.Printf("%+v\n", order)
		database.Repo.UpdateOrder(order)
	}
}
