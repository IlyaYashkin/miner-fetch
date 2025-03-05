package telegram

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (h *Handler) AuthMiddleware(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		fromUsername := update.Message.From.Username

		for _, username := range h.cfg.TgAdminUsernames {
			if fromUsername == username {
				next(ctx, b, update)
				return
			}
		}
	}
}

func (h *Handler) ChatIdSaveMiddleware(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		err := h.s.TelegramSender.SaveChatId(update.Message.From.Username, update.Message.Chat.ID)
		if err != nil {
			h.s.Logger.Log(err)
		}

		next(ctx, b, update)
		return
	}
}
