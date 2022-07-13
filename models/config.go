package models

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	Database      string `env:"DATABASE_DSN"`
	AccrualSystem string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}
