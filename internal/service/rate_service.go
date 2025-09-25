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
	}

	zap.L().Info("Все курсы обновлены")
	return nil
}

// GetRateInfo просто вызывает метод репозитория
func (s *RateService) GetRateInfo(currency string) (*database.RateInfo, error) {
	return s.rateRepo.GetRateInfo(currency)
}

// GetAllRateInfo возвращает информацию для всех валют
func (s *RateService) GetAllRateInfo() map[string]*database.RateInfo {
	currencies := []string{"BTC", "ETH"}
	results := make(map[string]*database.RateInfo)

	for _, currency := range currencies {
		info, err := s.rateRepo.GetRateInfo(currency)
		if err != nil {
			zap.L().Error("Ошибка получения информации",
				zap.String("валюта", currency),
				zap.Error(err),
			)
			continue
		}
		results[currency] = info
	}

	return results
}
