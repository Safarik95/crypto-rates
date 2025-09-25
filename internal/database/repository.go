package database

import (
	"database/sql"
	"fmt"
)

type RateRepository struct {
	db *sql.DB
}

func NewRateRepository(db *sql.DB) *RateRepository {
	return &RateRepository{db: db}
}

func (r *RateRepository) SaveRate(currency string, price float64) error {
	query := `INSERT INTO rates (currency_code, price) VALUES ($1, $2)`
	_, err := r.db.Exec(query, currency, price)
	if err != nil {
		return fmt.Errorf("ошибка сохранения курса: %w", err)
	}
	return nil
}
