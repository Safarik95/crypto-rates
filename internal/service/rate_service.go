package service

import (
	"crypto-rates/internal/api"
	"crypto-rates/internal/database"
	"fmt"
	"go.uber.org/zap"
)

type RateService struct {
	binanceClient api.BinanceClientInterface
	rateRepo      database.RateRepositoryInterface
}

func NewRateService(binanceClient api.BinanceClientInterface, rateRepo database.RateRepositoryInterface) *RateService {
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

// GetRateInfo возвращает информацию для ОДНОЙ валюты
func (s *RateService) GetRateInfo(currency string) (*database.RateInfo, error) {
	zap.L().Debug("GetRateInfo вызван для валюты", zap.String("валюта", currency))

	// Получаем данные только для запрошенной валюты
	rateInfo, err := s.rateRepo.GetRateInfo(currency)
	if err != nil {
		zap.L().Error("Ошибка в GetRateInfo", zap.String("валюта", currency), zap.Error(err))
		return nil, fmt.Errorf("ошибка получения информации для %s: %w", currency, err)
	}

	zap.L().Debug("GetRateInfo успешно завершен",
		zap.String("валюта", rateInfo.Currency),
		zap.Float64("цена", rateInfo.CurrentPrice),
	)

	return rateInfo, nil
}

// GetAllRateInfo возвращает информацию для ВСЕХ валют
func (s *RateService) GetAllRateInfo() map[string]*database.RateInfo {
	currencies := []string{"BTC", "ETH"}
	results := make(map[string]*database.RateInfo)

	for _, currency := range currencies {
		zap.L().Debug("Получение информации для валюты", zap.String("валюта", currency))

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

	zap.L().Debug("GetAllRateInfo завершен", zap.Int("количество_валют", len(results)))
	return results
}
