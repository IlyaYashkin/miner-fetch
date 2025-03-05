package runner

import (
	"context"
	"github.com/go-telegram/bot"
	"miner-fetch/internal/handler/telegram"
	"os"
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
			bot.WithMiddlewares(handler.AuthMiddleware, handler.ChatIdSaveMiddleware),
			bot.WithMessageTextHandler("/start", bot.MatchTypeExact, handler.Start),
			bot.WithDefaultHandler(handler.Default),
		}

		b, err := bot.New(t.cfg.TgAPIKey, opts...)
		if err != nil {
			t.s.Logger.Log(err)
			os.Exit(1)
		}

		err = t.s.TelegramSender.Init(b)
		if err != nil {
			t.s.Logger.Log(err)
		}

		b.Start(t.ctx)

		t.stopCh <- true
	}()
}

func (t *TelegramBot) GetName() string {
	return "TelegramBot"
}
