package service

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/basel-ax/luckysix/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// TelegramService handles sending notifications to Telegram.
type TelegramService struct {
	bot    *tgbotapi.BotAPI
	chatID int64
}

// NewTelegramService creates a new Telegram service.
func NewTelegramService(cfg *config.Config) (*TelegramService, error) {
	if cfg.PARAM.TgBotApi == "" {
		return nil, errors.New("telegram bot token is not configured")
	}
	if cfg.PARAM.TgChatId == "" {
		return nil, errors.New("telegram chat id is not configured")
	}

	bot, err := tgbotapi.NewBotAPI(cfg.PARAM.TgBotApi)
	if err != nil {
		return nil, fmt.Errorf("failed to create telegram bot: %w", err)
	}

	chatID, err := strconv.ParseInt(cfg.PARAM.TgChatId, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse telegram chat id: %w", err)
	}

	return &TelegramService{
		bot:    bot,
		chatID: chatID,
	}, nil
}

// SendNotification sends a message to the configured Telegram chat.
func (s *TelegramService) SendNotification(message string) error {
	msg := tgbotapi.NewMessage(s.chatID, message)
	msg.ParseMode = "Markdown"
	msg.DisableWebPagePreview = false

	_, err := s.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send telegram message: %w", err)
	}
	return nil
}
