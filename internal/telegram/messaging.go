package telegram

import (
	"crypto-rates/internal/database"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"strings"
)

func (b *Bot) formatRatesMessage(rates map[string]*database.RateInfo, title string) string {
	var sb strings.Builder
	for currency, info := range rates {
		trend := "➡️"
		if strings.HasPrefix(info.Change1h, "+") {
			trend = "📈"
		} else if strings.HasPrefix(info.Change1h, "-") {
			trend = "📉"
		}
		sb.WriteString(fmt.Sprintf(
			"*%s*: $%.2f\n"+
				"24ч: $%.2f - $%.2f\n"+
				"%s Изменение за час: %s\n\n",
			currency,
			info.CurrentPrice,
			info.MinPrice24h,
			info.MaxPrice24h,
			trend,
			info.Change1h,
		))
	}
	sb.WriteString("Обновляется каждые 5 минут")
	return sb.String()
}

func (b *Bot) sendMessage(bot *tgbotapi.BotAPI, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	_, err := bot.Send(msg)
	if err != nil {
		zap.L().Error("Ошибка отправки сообщения",
			zap.Int64("chat_id", chatID),
			zap.Error(err))
	}
}
