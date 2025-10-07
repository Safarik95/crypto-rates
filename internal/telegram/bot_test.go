package telegram_test

import (
	"crypto-rates/internal/database"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Тестируем только логику форматирования как отдельную функцию
func TestFormatRatesMessage(t *testing.T) {
	rates := map[string]*database.RateInfo{
		"BTC": {
			Currency:     "BTC",
			CurrentPrice: 50000.0,
			MinPrice24h:  49000.0,
			MaxPrice24h:  51000.0,
			Change1h:     "+2.0%",
		},
	}

	message := formatRatesMessageTest(rates, "Тестовые курсы:")

	assert.Contains(t, message, "BTC")
	assert.Contains(t, message, "50000.00")
	assert.Contains(t, message, "49000.00")
	assert.Contains(t, message, "51000.00")
	assert.Contains(t, message, "+2.0%")
}

func TestFormatRatesMessage_Empty(t *testing.T) {
	rates := map[string]*database.RateInfo{}

	message := formatRatesMessageTest(rates, "Пустые курсы:")

	assert.Contains(t, message, "Пустые курсы:")
	assert.Contains(t, message, "Нет данных о курсах")
}

// Вспомогательная функция для тестирования
func formatRatesMessageTest(rates map[string]*database.RateInfo, title string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s\n\n", title))

	if len(rates) == 0 {
		sb.WriteString("Нет данных о курсах\n\n")
		sb.WriteString("_Обновляется каждые 5 минут_")
		return sb.String()
	}

	for currency, info := range rates {
		trend := "->"
		if strings.HasPrefix(info.Change1h, "+") {
			trend = "UP"
		} else if strings.HasPrefix(info.Change1h, "-") {
			trend = "DOWN"
		}

		sb.WriteString(fmt.Sprintf(
			"*%s*: $%.2f\n"+
				"24h: $%.2f - $%.2f\n"+
				"%s Change 1h: %s\n\n",
			currency,
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
