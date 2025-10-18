package database

import (
	"context"
	"crypto-rates/internal/types"
)

type RateRepositoryInterface interface {
	SaveRate(ctx context.Context, currency types.Currency, price float64) error
	GetCurrentPrice(ctx context.Context, currency types.Currency) (float64, error)
	GetSimpleDailyStats(ctx context.Context, currency types.Currency) (minPrice, maxPrice float64, err error)
	GetHourlyChangePercent(ctx context.Context, currency types.Currency) string
	GetRateInfo(ctx context.Context, currency types.Currency) (*RateInfo, error)
}
