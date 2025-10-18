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
			zap.L().Info("Tables or indexes already exist (this is normal)")
			return nil
		}
		return err
	}
	zap.L().Info("Migrations applied successfully")
	return nil
}
