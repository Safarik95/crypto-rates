package telegram

import (
	"crypto-rates/internal/database"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	command := message.Command()

	switch command {
	case "start":
		b.handleStart(bot, message)
	case "rates":
		b.handleRates(bot, message)
	case "btc", "bitcoin": // ← отдельная команда для BTC
		b.handleBTC(bot, message)
	case "eth", "ethereum": // ← отдельная команда для ETH
		b.handleETH(bot, message)
	case "startAuto":
		b.handleStartAuto(bot, message)
	case "stopAuto":
		b.handleStopAuto(bot, message)
	default:
		b.sendMessage(bot, message.Chat.ID, "Неизвестная команда: "+command)
	}
}

func (b *Bot) handleStart(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	welcomeMsg := `Добро пожаловать в Crypto Rates Bot!

Доступные команды:
/rates - все курсы
/btc или /bitcoin - курс Bitcoin  
/eth или /ethereum - курс Ethereum
/startAuto - автоотправка каждые 10 мин
/stopAuto - остановить автоотправку`

	b.sendMessage(bot, message.Chat.ID, welcomeMsg)
}

func (b *Bot) handleRates(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	rates := b.service.GetAllRateInfo()
	title := "Текущие курсы:"

	if len(rates) == 0 {
		b.sendMessage(bot, message.Chat.ID, "Нет данных о курсах")
		return
	}

	messageText := b.formatRatesMessage(rates, title)
	b.sendMessage(bot, message.Chat.ID, messageText)
}

func (b *Bot) handleBTC(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	rates := b.service.GetAllRateInfo()

	if btcInfo, exists := rates["BTC"]; exists {
		rates = map[string]*database.RateInfo{"BTC": btcInfo}
		messageText := b.formatRatesMessage(rates, "Курс Bitcoin:")
		b.sendMessage(bot, message.Chat.ID, messageText)
	} else {
		b.sendMessage(bot, message.Chat.ID, "Нет данных по Bitcoin")
	}
}

func (b *Bot) handleETH(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	rates := b.service.GetAllRateInfo()

	if ethInfo, exists := rates["ETH"]; exists {
		rates = map[string]*database.RateInfo{"ETH": ethInfo}
		messageText := b.formatRatesMessage(rates, "Курс Ethereum:")
		b.sendMessage(bot, message.Chat.ID, messageText)
	} else {
		b.sendMessage(bot, message.Chat.ID, "Нет данных по Ethereum")
	}
}

func (b *Bot) handleStartAuto(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	b.autoSendUsers[message.Chat.ID] = true
	b.sendMessage(bot, message.Chat.ID, "Автоотправка включена! Курсы будут приходить каждые 10 минут")
}

func (b *Bot) handleStopAuto(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	delete(b.autoSendUsers, message.Chat.ID)
	b.sendMessage(bot, message.Chat.ID, "Автоотправка отключена")
}

func (b *Bot) handleText(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	b.sendMessage(bot, message.Chat.ID, "Отправьте /start для просмотра доступных команд")
}
