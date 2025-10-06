package telegram

import (
	"crypto-rates/internal/service"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type Bot struct {
	token         string
	service       *service.RateService
	stopCh        chan struct{}
	autoSendUsers map[int64]bool
}

func NewBot(token string, rateService *service.RateService) *Bot {
	return &Bot{
		token:         token,
		service:       rateService,
		stopCh:        make(chan struct{}),
		autoSendUsers: make(map[int64]bool),
	}
}

func (b *Bot) Start() error {
	bot, err := tgbotapi.NewBotAPI(b.token)
	if err != nil {
		return fmt.Errorf("ошибка создания бота: %w", err)
	}
	bot.Debug = false
	zap.L().Info("Авторизован", zap.String("аккаунт", bot.Self.UserName))
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	go b.autoSendWorker(bot)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		zap.L().Debug("Получено сообщение",
			zap.Int64("user_id", update.Message.From.ID),
			zap.String("text", update.Message.Text))
		b.handleMessage(bot, update.Message)
	}
	return nil
}

func (b *Bot) Stop() {
	close(b.stopCh)
}

func (b *Bot) handleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	switch {
	case message.IsCommand():
		b.handleCommand(bot, message)
	default:
		b.handleText(bot, message)
	}
}
