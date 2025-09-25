package main

import (
	"crypto-rates/internal/api"
	"crypto-rates/internal/config"
	"crypto-rates/internal/database"
	"crypto-rates/internal/logger"
	"crypto-rates/internal/service"
	"fmt"
	"go.uber.org/zap"
	"time"
)

func main() {
	// Инициализируем логгер
	log, err := logger.New()
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	// Устанавливаем глобальный логгер
	zap.ReplaceGlobals(log)

	// Загружаем конфигурацию
	cfg, err := config.Load()
	if err != nil {
		zap.L().Fatal("Ошибка загрузки конфигурации", zap.Error(err))
	}

	// Подключаемся к БД
	db, err := database.NewPostgresConnection(cfg)
	if err != nil {
		zap.L().Fatal("Ошибка подключения к БД", zap.Error(err))
	}
	defer db.Close()

	// Применяем миграции
	if err := database.Migrate(db); err != nil {
		zap.L().Error("Ошибка применения миграций", zap.Error(err))
	}

	// Инициализируем сервисы
	binanceClient := api.NewBinanceClient(cfg.BinanceAPIURL)
	rateRepo := database.NewRateRepository(db)
	rateService := service.NewRateService(binanceClient, rateRepo)

	// Настраиваем интервал обновления
	interval := time.Duration(cfg.UpdateIntervalMinutes) * time.Minute
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	zap.L().Info("Сервис запущен",
		zap.String("интервал", interval.String()),
		zap.String("API", cfg.BinanceAPIURL),
	)

	// Первый запуск
	if err := rateService.FetchAndStoreRates(); err != nil {
		zap.L().Error("Ошибка первого обновления", zap.Error(err))
	}

	// Показываем информацию о курсах после первого обновления
	showRateInfo(rateService)

	// Основной цикл
	for range ticker.C {
		zap.L().Debug("Запуск обновления курсов")

		if err := rateService.FetchAndStoreRates(); err != nil {
			zap.L().Error("Ошибка обновления курсов", zap.Error(err))
		}

		// Показываем информацию каждые 5 циклов (чтобы не засорять логи)
		showRateInfo(rateService)
	}
}

// showRateInfo показывает информацию о курсах
func showRateInfo(rateService *service.RateService) {
	zap.L().Info("=== ИНФОРМАЦИЯ О КУРСАХ ===")

	rateInfo := rateService.GetAllRateInfo()
	for currency, info := range rateInfo {
		zap.L().Info("Курс",
			zap.String("валюта", currency),
			zap.Float64("текущая_цена", info.CurrentPrice),
			zap.Float64("мин_24ч", info.MinPrice24h),
			zap.Float64("макс_24ч", info.MaxPrice24h),
			zap.String("изменение_за_час", info.Change1h),
		)

		// Красивый вывод в консоль
		fmt.Printf("%s: $%.2f (24ч: $%.2f - $%.2f) %s\n",
			currency, info.CurrentPrice, info.MinPrice24h, info.MaxPrice24h, info.Change1h)
	}
	fmt.Println() // Пустая строка для разделения
}
