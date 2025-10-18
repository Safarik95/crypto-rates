package database

import (
	"context"
	"crypto-rates/internal/config"
	"crypto-rates/internal/types"
	"database/sql"
	"fmt"
)

type RateRepository struct {
	db *sql.DB
}

func NewRateRepository(db *sql.DB) *RateRepository {
	return &RateRepository{db: db}
}

func (r *RateRepository) SaveRate(ctx context.Context, currency types.Currency, price float64) error {
	query := `INSERT INTO rates (currency_code, price) VALUES ($1, $2)`

	_, err := r.db.ExecContext(ctx, query, currency.String(), price)
	if err != nil {
		return fmt.Errorf("failed to save rate: %w", err)
	}

	return nil
}

type RateInfo struct {
	Currency     string
	CurrentPrice float64
	MinPrice24h  float64
	MaxPrice24h  float64
	Change1h     string
}

func (r *RateRepository) GetCurrentPrice(ctx context.Context, currency types.Currency) (float64, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DatabaseQueryTimeout)
	defer cancel()
	query := `SELECT price FROM rates WHERE currency_code = $1 ORDER BY timestamp DESC LIMIT 1`

	var price float64
	err := r.db.QueryRowContext(ctx, query, currency.String()).Scan(&price)
	if err != nil {
		return 0, fmt.Errorf("failed to get %s price: %w", currency, err)
	}

	return price, nil
}

func (r *RateRepository) GetSimpleDailyStats(ctx context.Context, currency types.Currency) (minPrice, maxPrice float64, err error) {
	query := `
		SELECT 
			COALESCE(MIN(price), 0),
			COALESCE(MAX(price), 0)
		FROM rates 
		WHERE currency_code = $1 
		AND timestamp >= NOW() - INTERVAL '24 hours'
	`

	err = r.db.QueryRowContext(ctx, query, currency.String()).Scan(&minPrice, &maxPrice)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get daily stats: %w", err)
	}
	return minPrice, maxPrice, nil
}

func (r *RateRepository) GetHourlyChangePercent(ctx context.Context, currency types.Currency) string {
	var currentPrice float64
	err := r.db.QueryRowContext(ctx, "SELECT price FROM rates WHERE currency_code = $1 ORDER BY timestamp DESC LIMIT 1", currency.String()).Scan(&currentPrice)
	if err != nil {
		return "0%"
	}

	var hourAgoPrice float64
	err = r.db.QueryRowContext(ctx, `
		SELECT price 
		FROM rates 
		WHERE currency_code = $1 
		AND timestamp <= NOW() - INTERVAL '1 hour' 
		ORDER BY timestamp DESC 
		LIMIT 1`, currency.String()).Scan(&hourAgoPrice)

	if err != nil {
		return "0%"
	}

	if hourAgoPrice == 0 {
		return "0%"
	}

	change := ((currentPrice - hourAgoPrice) / hourAgoPrice) * 100

	if change > 0 {
		return fmt.Sprintf("+%.2f%%", change)
	}
	return fmt.Sprintf("%.2f%%", change)
}

func (r *RateRepository) GetRateInfo(ctx context.Context, currency types.Currency) (*RateInfo, error) {
	currentPrice, err := r.GetCurrentPrice(ctx, currency)
	if err != nil {
		return nil, err
	}

	minPrice, maxPrice, err := r.GetSimpleDailyStats(ctx, currency)
	if err != nil {
		minPrice = currentPrice
		maxPrice = currentPrice
	}

	change1h := r.GetHourlyChangePercent(ctx, currency)

	return &RateInfo{
		Currency:     currency.String(),
		CurrentPrice: currentPrice,
		MinPrice24h:  minPrice,
		MaxPrice24h:  maxPrice,
		Change1h:     change1h,
	}, nil
}
