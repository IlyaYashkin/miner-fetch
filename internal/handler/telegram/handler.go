package telegram

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"miner-fetch/internal/config"
	"miner-fetch/internal/service"
)

type Handler struct {
	s   *service.Service
	cfg config.Config
}

func NewHandler(s *service.Service, cfg config.Config) *Handler {
	return &Handler{
		s:   s,
		cfg: cfg,
	}
}

func (h *Handler) Start(ctx context.Context, b *bot.Bot, update *models.Update) {
	commands := []models.BotCommand{
		{Command: "start", Description: "Запустить бота"},
		{Command: "info", Description: "Получить информацию"},
	}
	_, err := b.SetMyCommands(ctx, &bot.SetMyCommandsParams{Commands: commands})
	if err != nil {
		h.s.Logger.Log(err)
	}
}

func (h *Handler) Info(ctx context.Context, b *bot.Bot, update *models.Update) {
	if h.cfg.IsScanner {
		text, err := h.s.Device.GetDevicesInfo()
		if err != nil {
			h.s.Logger.Log(err)
		}

		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   text,
		})
		if err != nil {
			h.s.Logger.Log(err)
		}
	}

	h.s.Polling.Send("info")
}
