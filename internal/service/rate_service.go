package service

import (
	"context"
	"crypto-rates/internal/database"
	"crypto-rates/internal/types"
	"crypto-rates/pkg/api"
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

func (s *RateService) FetchAndStoreRates(ctx context.Context) error {
	currencies := []types.Currency{types.BTC, types.ETH}

	for _, currency := range currencies {
		price, err := s.binanceClient.GetRate(ctx, currency)
		if err != nil {
			return fmt.Errorf("failed to get %s rate: %w", currency, err)
		}

		err = s.rateRepo.SaveRate(ctx, currency, price)
		if err != nil {
			return fmt.Errorf("failed to save %s rate: %w", currency, err)
		}

		zap.L().Info("Rate updated",
			zap.String("currency", currency.String()),
			zap.Float64("price", price),
		)
	}

	zap.L().Info("All rates updated")
	return nil
}

func (s *RateService) GetRateInfo(ctx context.Context, currency types.Currency) (*database.RateInfo, error) {
	rateInfo, err := s.rateRepo.GetRateInfo(ctx, currency)
	if err != nil {
		zap.L().Error("Error getting rate info",
			zap.String("currency", currency.String()),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get info for %s: %w", currency, err)
	}

	return rateInfo, nil
}

func (s *RateService) GetAllRateInfo(ctx context.Context) map[types.Currency]*database.RateInfo {
	currencies := []types.Currency{types.BTC, types.ETH}
	results := make(map[types.Currency]*database.RateInfo, len(currencies))

	for _, currency := range currencies {
		info, err := s.rateRepo.GetRateInfo(ctx, currency)
		if err != nil {
			zap.L().Error("Error getting rate info",
				zap.String("currency", currency.String()),
				zap.Error(err),
			)
			continue
		}
		results[currency] = info
	}

	return results
}
