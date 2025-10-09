package database

import (
	"database/sql"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
	"strings"
)

func Migrate(db *sql.DB) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	// Применяем миграции
	err := goose.Up(db, "migrations")
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			zap.L().Info("Таблицы или индексы уже существуют (это нормально)")
			return nil
		}
		return err
	}
	zap.L().Info("Миграции применены успешно")
	return nil
}
