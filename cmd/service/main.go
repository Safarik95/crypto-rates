// @title Crypto Rates API
// @version 1.0
// @description Сервис для отслеживания курсов криптовалют

// @host localhost:8080
// @BasePath /
package main

import (
	"context"
	"crypto-rates/internal/api"
	"crypto-rates/internal/api/rest"
	"crypto-rates/internal/config"
	"crypto-rates/internal/database"
	"crypto-rates/internal/logger"
	"crypto-rates/internal/service"
	"crypto-rates/internal/telegram"
	"fmt"
	"go.uber.org/zap"
	"time"
)

func main() {
	log, err := logger.New()
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	zap.ReplaceGlobals(log)

	cfg, err := config.Load()
	if err != nil {
		zap.L().Fatal("Ошибка загрузки конфигурации", zap.Error(err))
	}

	if err := waitForDB(cfg); err != nil {
		zap.L().Fatal("База данных недоступна", zap.Error(err))
	}

	db, err := database.NewPostgresConnection(cfg)
	if err != nil {
		zap.L().Fatal("Ошибка подключения к БД", zap.Error(err))
	}
	defer db.Close()

	if err := database.Migrate(db); err != nil {
		zap.L().Error("Ошибка применения миграций", zap.Error(err))
	}

	binanceClient := api.NewBinanceClient(cfg.BinanceAPIURL)
	rateRepo := database.NewRateRepository(db)
	rateService := service.NewRateService(binanceClient, rateRepo)

	restServer := rest.NewServer(cfg.APIPort, rateService)
	go func() {
		zap.L().Info("Запуск REST API сервера...")
		if err := restServer.Start(); err != nil {
			zap.L().Error("Ошибка REST API сервера", zap.Error(err))
		}
	}()
	defer restServer.Stop(context.Background())

	if cfg.TelegramBotToken != "" {
		telegramBot := telegram.NewBot(cfg.TelegramBotToken, rateService)

		go func() {
			zap.L().Info("Запуск Telegram бота...")
			if err := telegramBot.Start(); err != nil {
				zap.L().Error("Ошибка Telegram бота", zap.Error(err))
			}
		}()
		defer telegramBot.Stop()

		zap.L().Info("Telegram бот запущен")
	} else {
		zap.L().Warn("TELEGRAM_BOT_TOKEN не установлен, бот не запущен")
	}

	interval := time.Duration(cfg.UpdateIntervalMinutes) * time.Minute
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	zap.L().Info("Сервис запущен",
		zap.String("интервал", interval.String()),
		zap.String("API", cfg.BinanceAPIURL),
	)

	if err := rateService.FetchAndStoreRates(); err != nil {
		zap.L().Error("Ошибка первого обновления", zap.Error(err))
	}

	showRateInfo(rateService)

	for range ticker.C {
		zap.L().Debug("Запуск обновления курсов")

		if err := rateService.FetchAndStoreRates(); err != nil {
			zap.L().Error("Ошибка обновления курсов", zap.Error(err))
		}

		showRateInfo(rateService)
	}
}

func waitForDB(cfg *config.Config) error {
	zap.L().Info("Ожидание подключения к базе данных...",
		zap.String("хост", cfg.DBHost),
		zap.Int("порт", cfg.DBPort),
	)

	timeout := time.Second * 30
	start := time.Now()

	for {
		db, err := database.NewPostgresConnection(cfg)
		if err == nil {
			db.Close()
			zap.L().Info("База данных доступна",
				zap.Duration("время_ожидания", time.Since(start)),
			)
			return nil
		}

		if time.Since(start) > timeout {
			return fmt.Errorf("таймаут подключения к БД: %v", err)
		}

		zap.L().Debug("База данных еще не доступна, повторная попытка...",
			zap.Error(err),
		)
		time.Sleep(2 * time.Second)
	}
}

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

		fmt.Printf("%s: $%.2f (24ч: $%.2f - $%.2f) %s\n",
			currency, info.CurrentPrice, info.MinPrice24h, info.MaxPrice24h, info.Change1h)
	}
	fmt.Println()
}
