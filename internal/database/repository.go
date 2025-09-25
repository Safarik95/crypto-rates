package database

import (
	"database/sql"
	"fmt"
	"go.uber.org/zap"
)

type RateRepository struct {
	db *sql.DB
}

func NewRateRepository(db *sql.DB) *RateRepository {
	return &RateRepository{db: db}
}

// SaveRate сохраняет курс валюты в базу данных
func (r *RateRepository) SaveRate(currency string, price float64) error {
	query := `INSERT INTO rates (currency_code, price) VALUES ($1, $2)`

	_, err := r.db.Exec(query, currency, price)
	if err != nil {
		return fmt.Errorf("ошибка сохранения курса: %w", err)
	}

	return nil
}

// RateInfo упрощенная структура для хранения информации
type RateInfo struct {
	Currency     string
	CurrentPrice float64
	MinPrice24h  float64
	MaxPrice24h  float64
	Change1h     string
}

// GetCurrentPrice возвращает текущую цену
func (r *RateRepository) GetCurrentPrice(currency string) (float64, error) {
	query := `SELECT price FROM rates WHERE currency_code = $1 ORDER BY timestamp DESC LIMIT 1`

	var price float64
	err := r.db.QueryRow(query, currency).Scan(&price)
	if err != nil {
		return 0, fmt.Errorf("ошибка получения цены %s: %w", currency, err)
	}

	return price, nil
}

// GetSimpleDailyStats простой метод для мин/макс за 24 часа
func (r *RateRepository) GetSimpleDailyStats(currency string) (minPrice, maxPrice float64, err error) {
	query := `
		SELECT 
			COALESCE(MIN(price), 0),
			COALESCE(MAX(price), 0)
		FROM rates 
		WHERE currency_code = $1 
		AND timestamp >= NOW() - INTERVAL '24 hours'
	`

	err = r.db.QueryRow(query, currency).Scan(&minPrice, &maxPrice)
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка получения статистики: %w", err)
	}
	return
}

// GetHourlyChangePercent возвращает изменение в процентах за час
func (r *RateRepository) GetHourlyChangePercent(currency string) string {
	query := `
		WITH current_price AS (
			SELECT price, timestamp
			FROM rates 
			WHERE currency_code = $1 
			ORDER BY timestamp DESC 
			LIMIT 1
		),
		hour_ago_price AS (
			SELECT price
			FROM rates 
			WHERE currency_code = $1 
			AND timestamp <= (SELECT timestamp FROM current_price) - INTERVAL '1 hour'
			ORDER BY timestamp DESC 
			LIMIT 1
		)
		SELECT 
			current_price.price as current,
			hour_ago_price.price as hour_ago
		FROM current_price, hour_ago_price
	`

	var currentPrice, hourAgoPrice float64
	err := r.db.QueryRow(query, currency).Scan(&currentPrice, &hourAgoPrice)
	if err != nil {
		zap.L().Debug("Не удалось получить данные за час", zap.String("валюта", currency), zap.Error(err))
		return "0%"
	}

	// Если нет данных за прошлый час
	if hourAgoPrice == 0 {
		return "0%"
	}

	change := ((currentPrice - hourAgoPrice) / hourAgoPrice) * 100

	if change > 0 {
		return fmt.Sprintf("+%.2f%%", change)
	}
	return fmt.Sprintf("%.2f%%", change)
}

// GetRateInfo возвращает всю информацию для валюты
func (r *RateRepository) GetRateInfo(currency string) (*RateInfo, error) {
	// Текущая цена
	currentPrice, err := r.GetCurrentPrice(currency)
	if err != nil {
		return nil, err
	}

	// Мин/макс за 24 часа
	minPrice, maxPrice, err := r.GetSimpleDailyStats(currency)
	if err != nil {
		// Fallback - используем текущую цену
		minPrice = currentPrice
		maxPrice = currentPrice
	}

	// Изменение за час
	change1h := r.GetHourlyChangePercent(currency)

	return &RateInfo{
		Currency:     currency,
		CurrentPrice: currentPrice,
		MinPrice24h:  minPrice,
		MaxPrice24h:  maxPrice,
		Change1h:     change1h,
	}, nil
}
