package database

import (
	"context"
	"crypto-rates/internal/config"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"time"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func NewPostgresConnection(cfg *config.Config) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("error opening connection: %w", err)
	}

	db.SetConnMaxLifetime(time.Hour)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	ctx, cancel := context.WithTimeout(context.Background(), config.DatabasePingTimeout)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("DB ping error: %w", err)
	}
	return db, nil
}
