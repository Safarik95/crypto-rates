package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"time"
)

const (
	HTTPClientTimeout     = 10 * time.Second
	DatabasePingTimeout   = 5 * time.Second
	ServerReadTimeout     = 10 * time.Second
	ServerWriteTimeout    = 10 * time.Second
	ServerIdleTimeout     = 60 * time.Second
	BinanceRequestTimeout = 10 * time.Second
	DatabaseQueryTimeout  = 5 * time.Second
	InitialRatesTimeout   = 30 * time.Second
)

type Config struct {
	DBHost                string
	DBPort                int
	DBUser                string
	DBPassword            string
	DBName                string
	BinanceAPIURL         string
	UpdateIntervalMinutes int
	LogLevel              string
	TelegramBotToken      string
	APIPort               string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnvAsInt("DB_PORT", 5432),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "crypto_rates"),

		BinanceAPIURL:         getEnv("BINANCE_API_URL", "https://api.binance.com/api/v3"),
		UpdateIntervalMinutes: getEnvAsInt("UPDATE_INTERVAL_MINUTES", 5),
		LogLevel:              getEnv("LOG_LEVEL", "info"),
		TelegramBotToken:      getEnv("TELEGRAM_BOT_TOKEN", ""),
		APIPort:               getEnv("API_PORT", "8080"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
