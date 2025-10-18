package telegram

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"time"
)

func (b *Bot) autoSendWorker(bot *tgbotapi.BotAPI) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	lastSent := make(map[int64]time.Time, len(b.autoSendUsers))

	for {
		select {
		case <-ticker.C:
			now := time.Now()
			for userID, autoSendEnabled := range b.autoSendUsers {
				if !autoSendEnabled {
					continue
				}

				if lastSend, exists := lastSent[userID]; exists {
					nextSend := lastSend.Add(10 * time.Minute)
					if now.Before(nextSend) {
						continue
					}
				}

				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				rates := b.service.GetAllRateInfo(ctx)
				cancel()

				message := b.formatRatesMessage(rates, "Automatic update:")

				msg := tgbotapi.NewMessage(userID, message)
				_, err := bot.Send(msg)
				if err != nil {
					zap.L().Error("Error sending auto message",
						zap.Int64("user_id", userID),
						zap.Error(err),
					)
					delete(b.autoSendUsers, userID)
					continue
				}

				lastSent[userID] = now

				zap.L().Debug("Auto-sending completed",
					zap.Int64("user_id", userID),
				)
			}

		case <-b.stopCh:
			return
		}
	}
}
