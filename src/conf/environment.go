package conf

import (
	"os"
)

type Config struct {
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
	APISecret  string
}

func LoadConfig() *Config {
	apiSecret := os.Getenv("API_SECRET")
	if apiSecret == "" {
		apiSecret = "supersecretkey" // Default for development
	}

	return &Config{
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBName:     os.Getenv("DB_NAME"),
		APISecret:  apiSecret,
	}
}
