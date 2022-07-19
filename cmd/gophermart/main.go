package main

import (
	"AlexSarva/gofermart/internal/app"
	"AlexSarva/gofermart/loyality"
	"AlexSarva/gofermart/models"
	"AlexSarva/gofermart/server"
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
)

func main() {
	var cfg models.Config
	// Приоритет будет у ФЛАГОВ
	// Загружаем конфиг из переменных окружения
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	// Перезаписываем из параметров запуска
	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "host:port to listen on")
	flag.StringVar(&cfg.Database, "d", cfg.Database, "database config")
	flag.StringVar(&cfg.AccrualSystem, "r", cfg.AccrualSystem, "address of the accrual system")
	flag.Parse()
	log.Printf("%+v\n", cfg)
	log.Printf("ServerAddress: %v", cfg.ServerAddress)
	DB, dbErr := app.NewStorage(cfg.Database)
	client := loyality.NewProcessingClient(cfg.AccrualSystem, "/api/orders")
	ordersToProcessingCh := make(chan string)
	ordersProcessedCh := make(chan models.ProcessingOrder)
	go loyality.GetOrdersToProcessing(*DB, ordersToProcessingCh)
	go loyality.GetProcessedInfo(client, ordersToProcessingCh, ordersProcessedCh)
	go loyality.ApplyLoyality(*DB, ordersProcessedCh)
	if dbErr != nil {
		log.Fatal(dbErr)
	}
	ping := DB.Repo.Ping()
	log.Println(ping)
	MainApp := server.NewServer(&cfg, DB)
	if runErr := MainApp.Run(); runErr != nil {
		log.Printf("%s", runErr.Error())
	}
}
