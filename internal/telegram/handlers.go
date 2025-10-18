package telegram

import (
	"context"
	"crypto-rates/internal/database"
	"crypto-rates/internal/types"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

const (
	CommandStart     = "start"
	CommandRates     = "rates"
	CommandBTC       = "btc"
	CommandBitcoin   = "bitcoin"
	CommandETH       = "eth"
	CommandEthereum  = "ethereum"
	CommandStartAuto = "startAuto"
	CommandStopAuto  = "stopAuto"
)

func (b *Bot) handleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	command := message.Command()

	switch command {
	case CommandStart:
		b.handleStart(bot, message)
	case CommandRates:
		b.handleRates(bot, message)
	case CommandBTC, CommandBitcoin:
		b.handleSpecificCurrency(bot, message, types.BTC, "Bitcoin")
	case CommandETH, CommandEthereum:
		b.handleSpecificCurrency(bot, message, types.ETH, "Ethereum")
	case CommandStartAuto:
		b.handleStartAuto(bot, message)
	case CommandStopAuto:
		b.handleStopAuto(bot, message)
	default:
		b.sendMessage(bot, message.Chat.ID, "Unknown command: "+command)
	}
}

func (b *Bot) handleStart(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	welcomeMsg := `Welcome to Crypto Rates Bot!

Available commands:
/rates - all rates
/btc or /bitcoin - Bitcoin rate  
/eth or /ethereum - Ethereum rate
/startAuto - auto send every 10 min
/stopAuto - stop auto send`

	b.sendMessage(bot, message.Chat.ID, welcomeMsg)
}

func (b *Bot) handleRates(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rates := b.service.GetAllRateInfo(ctx)
	title := "Current rates:"

	if len(rates) == 0 {
		b.sendMessage(bot, message.Chat.ID, "No data on courses")
		return
	}

	messageText := b.formatRatesMessage(rates, title)
	b.sendMessage(bot, message.Chat.ID, messageText)
}

func (b *Bot) handleSpecificCurrency(bot *tgbotapi.BotAPI, message *tgbotapi.Message, currency types.Currency, name string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rates := b.service.GetAllRateInfo(ctx)

	if currencyInfo, exists := rates[currency]; exists {
		ratesMap := map[types.Currency]*database.RateInfo{currency: currencyInfo}
		messageText := b.formatRatesMessage(ratesMap, fmt.Sprintf("Rate %s:", name))
		b.sendMessage(bot, message.Chat.ID, messageText)
	} else {
		b.sendMessage(bot, message.Chat.ID, fmt.Sprintf("No data %s", name))
	}
}

func (b *Bot) handleStartAuto(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	b.autoSendUsers[message.Chat.ID] = true
	b.sendMessage(bot, message.Chat.ID, "Auto send enabled! Rates will be sent every 10 minutes")
}

func (b *Bot) handleStopAuto(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	delete(b.autoSendUsers, message.Chat.ID)
	b.sendMessage(bot, message.Chat.ID, "Auto send disabled")
}

func (b *Bot) handleText(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	b.sendMessage(bot, message.Chat.ID, "Send /start to see available commands")
}
