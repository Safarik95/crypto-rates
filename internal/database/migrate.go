package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func Migrate(db *sql.DB, migrationsDir string) error {
	files, err := filepath.Glob(filepath.Join(migrationsDir, "*.up.sql"))
	if err != nil {
		return fmt.Errorf("ошибка чтения директории миграций: %w", err)
	}
	for _, file := range files {
		log.Printf("Применяем миграцию: %s", file)
		sqlContent, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("ошибка чтения файла %s: %w", file, err)
		}
		_, err = db.Exec(string(sqlContent))
		if err != nil {
			if strings.Contains(err.Error(), "already exists") {
				log.Printf("Таблица уже существует: %s", file)
				continue
			}
			return fmt.Errorf("ошибка выполнения миграции %s: %w", file, err)
		}
		log.Printf("Миграция успешно применена: %s", file)
	}
	log.Println("Миграции завершены")
	return nil
}
