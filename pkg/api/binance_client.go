package api

import (
	"context"
	"crypto-rates/internal/types"
)

type BinanceClientInterface interface {
	GetRate(ctx context.Context, currency types.Currency) (float64, error)
}
