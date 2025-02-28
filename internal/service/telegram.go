package service

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
)

type TelegramSender struct {
	b *bot.Bot
}

func NewTelegramSender() *TelegramSender {
	return &TelegramSender{}
}

func (t *TelegramSender) SetBot(b *bot.Bot) {
	t.b = b
}

func (t *TelegramSender) SendMessage(ctx context.Context, chatID int64, nodeName string, message string) error {
	params := bot.SendMessageParams{
		ChatID: chatID,
		Text:   fmt.Sprintf("%s\n\n%s", nodeName, message),
	}

	_, err := t.b.SendMessage(ctx, &params)
	if err != nil {
		return err
	}

	return nil
}
