package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"time"
)

func (b *Bot) autoSendWorker(bot *tgbotapi.BotAPI) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	lastSent := make(map[int64]time.Time)
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

				rates := b.service.GetAllRateInfo()
				message := b.formatRatesMessage(rates, "Автоматическое обновление:")
				msg := tgbotapi.NewMessage(userID, message)
				_, err := bot.Send(msg)
				if err != nil {
					zap.L().Error("Ошибка отправки автосообщения",
						zap.Int64("user_id", userID),
						zap.Error(err),
					)
					delete(b.autoSendUsers, userID)
					continue
				}
				lastSent[userID] = now
				zap.L().Debug("Автоотправка выполнена",
					zap.Int64("user_id", userID),
				)
			}
		case <-b.stopCh:
			return
		}
	}
}
