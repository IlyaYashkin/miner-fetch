package telegram

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/pkg/errors"
	"log"
	"miner-fetch/internal/config"
	"miner-fetch/internal/usecase/telegram"
)

type Handler struct {
	Usecase *telegram.Usecase
	Config  config.Config
}

func (h *Handler) Start(ctx context.Context, b *bot.Bot, update *models.Update) {
	commands := []models.BotCommand{
		{Command: "start", Description: "Запустить бота"},
		{Command: "info", Description: "Получить информацию"},
	}
	_, err := b.SetMyCommands(ctx, &bot.SetMyCommandsParams{Commands: commands})
	if err != nil {
		log.Println(errors.Wrap(err, "telegram bot"))
	}
}

func (h *Handler) Info(ctx context.Context, b *bot.Bot, update *models.Update) {
	text := h.Usecase.GetDevicesInfo()

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   text,
	})
	if err != nil {
		log.Println(errors.Wrap(err, "telegram bot"))
	}
}
