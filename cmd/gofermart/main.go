package main

import (
	"AlexSarva/gofermart/internal/app"
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
	log.Printf("ServerAddress: %v", cfg.ServerAddress)
	// Перезаписываем из параметров запуска
	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "host:port to listen on")
	flag.StringVar(&cfg.Database, "d", cfg.Database, "database config")
	flag.StringVar(&cfg.Database, "r", cfg.Database, "address of the accrual system")
	flag.Parse()
	log.Printf("ServerAddress: %v", cfg.ServerAddress)
	DB, dbErr := app.NewStorage(cfg.Database)
	if dbErr != nil {
		log.Fatal(dbErr)
	}
	ping := DB.Repo.Ping()
	log.Println(ping)
	MainApp := server.NewServer(&cfg, DB)
	if err := MainApp.Run(); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
