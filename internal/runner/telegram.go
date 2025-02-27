package runner

import (
	"context"
	"github.com/go-telegram/bot"
	"log"
	"miner-fetch/internal/handler/telegram"
)

type TelegramBot struct {
	CommonRunner
}

func NewTelegramBot(runner CommonRunner) *TelegramBot {
	ctxc, cancel := context.WithCancel(runner.ctx)
	runner.ctx = ctxc
	runner.cancel = cancel

	return &TelegramBot{runner}
}

func (t *TelegramBot) Start() {
	go func() {
		handler := telegram.NewHandler(t.s, t.cfg)

		opts := []bot.Option{
			bot.WithMiddlewares(handler.AuthMiddleware),
			bot.WithMessageTextHandler("/start", bot.MatchTypeExact, handler.Start),
			bot.WithMessageTextHandler("/info", bot.MatchTypeExact, handler.Info),
		}

		b, err := bot.New(t.cfg.TgAPIKey, opts...)

		if err != nil {
			log.Fatalln(err)
		}

		b.Start(t.ctx)

		t.stopCh <- true
	}()
}

func (t *TelegramBot) GetName() string {
	return "TelegramBot"
}
