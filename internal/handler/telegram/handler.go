package telegram

import (
	"context"
	"errors"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"miner-fetch/internal/config"
	"miner-fetch/internal/service"
	"strings"
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

func (h *Handler) Start(ctx context.Context, b *bot.Bot, _ *models.Update) {
	commands := []models.BotCommand{
		{Command: "start", Description: "Запустить бота"},
		{Command: "status", Description: "Статус"},
		{Command: "temperature", Description: "Температура"},
		{Command: "ips", Description: "IP-адреса"},
	}
	_, err := b.SetMyCommands(ctx, &bot.SetMyCommandsParams{Commands: commands})
	if err != nil {
		h.s.Logger.Log(err)
	}
}

func (h *Handler) Default(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if !strings.HasPrefix(update.Message.Text, "/") {
		return
	}

	command := strings.TrimLeft(update.Message.Text, "/")

	if h.cfg.IsScanner {
		message, err := h.s.Device.ExecuteCommand(command)
		target := &service.CommandNotFoundError{}
		if err != nil && !errors.As(err, &target) {
			h.s.Logger.Log(err)
			return
		} else if errors.As(err, &target) {
			message = err.Error()
		}

		err = h.s.TelegramSender.SendMessage(ctx, update.Message.Chat.ID, h.cfg.NodeName, message)
		if err != nil {
			h.s.Logger.Log(err)
		}
	}

	h.s.Polling.Send(service.Payload{
		Command: command,
		ChatID:  update.Message.Chat.ID,
	})
}
