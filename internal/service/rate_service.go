package service

import (
	"crypto-rates/internal/api"
	"crypto-rates/internal/database"
	"fmt"
	"log"
)

type RateService struct {
	binanceClient *api.BinanceClient
	rateRepo      *database.RateRepository
}

func NewRateService(binanceClient *api.BinanceClient, rateRepo *database.RateRepository) *RateService {
	return &RateService{
		binanceClient: binanceClient,
		rateRepo:      rateRepo,
	}
}

func (s *RateService) FetchAndStoreRates() error {
	currencies := []string{"BTC", "ETH"}
	for _, currency := range currencies {
		price, err := s.binanceClient.GetRate(currency)
		if err != nil {
			return fmt.Errorf("ошибка получения %s: %w", currency, err)
		}
		err = s.rateRepo.SaveRate(currency, price)
		if err != nil {
			return fmt.Errorf("ошибка сохранения %s: %w", currency, err)
		}
		log.Printf("Успешно обновлен %s: $%.2f", currency, price)
	}
	log.Println("Все курсы успешно обновлены!")
	return nil
}
