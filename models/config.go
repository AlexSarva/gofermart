package models

type Config struct {
	ServerAddress string `env:"RUN_ADDRESS" envDefault:"localhost:8080"`
	Database      string `env:"DATABASE_URI"`
	AccrualSystem string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}
