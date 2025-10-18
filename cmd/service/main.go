// @title Crypto Rates API
// @version 1.0
// @description Сервис для отслеживания курсов криптовалют

// @host localhost:8080
// @BasePath /
package main

import (
	"context"
	"crypto-rates/internal/api/rest"
	"crypto-rates/internal/config"
	"crypto-rates/internal/database"
	"crypto-rates/internal/service"
	"crypto-rates/internal/telegram"
	"crypto-rates/pkg/api"
	"crypto-rates/pkg/logger"
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
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
		zap.L().Fatal("Failed to load configuration", zap.Error(err))
	}

	if err := waitForDB(cfg); err != nil {
		zap.L().Fatal("Database unavailable", zap.Error(err))
	}

	db, err := database.NewPostgresConnection(cfg)
	if err != nil {
		zap.L().Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	if err := database.Migrate(db); err != nil {
		zap.L().Error("Migration error", zap.Error(err))
	}

	binanceClient := api.NewBinanceClient(cfg.BinanceAPIURL, config.HTTPClientTimeout, config.BinanceRequestTimeout)
	rateRepo := database.NewRateRepository(db)
	rateService := service.NewRateService(binanceClient, rateRepo)
	restServer := rest.NewServer(cfg.APIPort, rateService)

	var telegramBot *telegram.Bot

	g, ctx := errgroup.WithContext(context.Background())

	// REST API server
	g.Go(func() error {
		zap.L().Info("Starting REST API server...")
		return restServer.Start()
	})

	// Telegram bot
	if cfg.TelegramBotToken != "" {
		telegramBot = telegram.NewBot(cfg.TelegramBotToken, rateService)
		g.Go(func() error {
			zap.L().Info("Starting Telegram bot...")
			return telegramBot.Start()
		})
		zap.L().Info("Telegram bot started")
	} else {
		zap.L().Warn("TELEGRAM_BOT_TOKEN not set, bot not started")
	}

	g.Go(func() error {
		interval := time.Duration(cfg.UpdateIntervalMinutes) * time.Minute
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		initialCtx, initialCancel := context.WithTimeout(ctx, config.InitialRatesTimeout)
		defer initialCancel()
		if err := rateService.FetchAndStoreRates(initialCtx); err != nil {
			zap.L().Error("First update error", zap.Error(err))
		}
		showRateInfo(rateService)

		for {
			select {
			case <-ticker.C:
				if err := rateService.FetchAndStoreRates(ctx); err != nil {
					zap.L().Error("Update rates error", zap.Error(err))
				}
				showRateInfo(rateService)
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	})

	defer restServer.Stop(context.Background())
	if telegramBot != nil {
		defer telegramBot.Stop()
	}

	zap.L().Info("Service started",
		zap.Int("interval", cfg.UpdateIntervalMinutes),
		zap.String("API", cfg.BinanceAPIURL),
	)

	// Wait for all goroutines
	if err := g.Wait(); err != nil {
		zap.L().Error("Error in one of goroutines", zap.Error(err))
	}
}

func waitForDB(cfg *config.Config) error {
	zap.L().Info("Waiting for database connection...",
		zap.String("host", cfg.DBHost),
		zap.Int("port", cfg.DBPort),
	)

	timeout := 30 * time.Second
	start := time.Now()

	for {
		db, err := database.NewPostgresConnection(cfg)
		if err == nil {
			db.Close()
			zap.L().Info("Database available",
				zap.Duration("wait_time", time.Since(start)),
			)
			return nil
		}

		if time.Since(start) > timeout {
			return fmt.Errorf("database connection timeout: %v", err)
		}

		zap.L().Debug("Database not available yet, retrying...",
			zap.Error(err),
		)
		time.Sleep(2 * time.Second)
	}
}

func showRateInfo(rateService *service.RateService) {
	zap.L().Info("=== RATES INFORMATION ===")
	ctx := context.Background()
	rateInfo := rateService.GetAllRateInfo(ctx)
	for currency, info := range rateInfo {
		zap.L().Info("Rate",
			zap.String("currency", currency.String()),
			zap.Float64("current_price", info.CurrentPrice),
			zap.Float64("min_24h", info.MinPrice24h),
			zap.Float64("max_24h", info.MaxPrice24h),
			zap.String("change_1h", info.Change1h),
		)
		fmt.Printf("%s: $%.2f (24h: $%.2f - $%.2f) %s\n",
			currency, info.CurrentPrice, info.MinPrice24h, info.MaxPrice24h, info.Change1h)
	}
	fmt.Println()
}
