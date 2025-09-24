package main

import (
	"crypto-rates/internal/api"
	"crypto-rates/internal/database"
	"crypto-rates/internal/service"
	"log"
	"time"
)

func main() {
	dbConfig := database.Config{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "password",
		DBName:   "crypto_rates",
	}
	db, err := database.NewPostgresConnection(dbConfig)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	err = database.Migrate(db, "migrations")
	if err != nil {
		log.Printf("Предупреждение: %v", err)
	}

	binanceClient := api.NewBinanceClient()
	rateRepo := database.NewRateRepository(db)
	rateService := service.NewRateService(binanceClient, rateRepo)
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	log.Println("Запуск первого обновления...")
	err = rateService.FetchAndStoreRates()
	if err != nil {
		log.Printf("Ошибка первого обновления: %v", err)
	}
	log.Println("Сервис запущен. Обновление каждые 5 минут...")
	for range ticker.C {
		log.Println("Запуск периодического обновления...")
		err = rateService.FetchAndStoreRates()
		if err != nil {
			log.Printf("Ошибка периодического обновления: %v", err)
		}
	}
}
