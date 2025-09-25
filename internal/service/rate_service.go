package service

import (
	"crypto-rates/internal/api"
	"crypto-rates/internal/database"
	"fmt"
	"go.uber.org/zap"
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

		zap.L().Info("Курс обновлен",
			zap.String("валюта", currency),
			zap.Float64("цена", price),
		)
	}

	zap.L().Info("Все курсы обновлены")
	return nil
}
