package telegram

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (h *Handler) AuthMiddleware(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		fromUsername := update.Message.From.Username

		for _, username := range h.Config.TelegramAdminUsernames {
			if fromUsername == username {
				next(ctx, b, update)
				return
			}
		}
	}
}
