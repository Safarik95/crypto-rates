package telegram

import (
	"crypto-rates/internal/database"
	"crypto-rates/internal/types"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"strings"
)

func (b *Bot) formatRatesMessage(rates map[types.Currency]*database.RateInfo, title string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s\n\n", title))

	if len(rates) == 0 {
		sb.WriteString("No rate data\n\n")
		sb.WriteString("_Updated every 5 minutes_")
		return sb.String()
	}

	for currency, info := range rates {
		trend := "➡️"
		if strings.HasPrefix(info.Change1h, "+") {
			trend = "📈"
		} else if strings.HasPrefix(info.Change1h, "-") {
			trend = "📉"
		}

		sb.WriteString(fmt.Sprintf(
			"*%s*: $%.2f\n"+
				"24h: $%.2f - $%.2f\n"+
				"%s Change 1h: %s\n\n",
			currency.String(),
			info.CurrentPrice,
			info.MinPrice24h,
			info.MaxPrice24h,
			trend,
			info.Change1h,
		))
	}

	sb.WriteString("_Updated every 5 minutes_")
	return sb.String()
}

func (b *Bot) sendMessage(bot *tgbotapi.BotAPI, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown

	_, err := bot.Send(msg)
	if err != nil {
		zap.L().Error("Error sending message",
			zap.Int64("chat_id", chatID),
			zap.Error(err),
		)
	}
}
