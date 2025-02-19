package service

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/pkg/errors"
	"log"
	"miner-fetch/internal/config"
	"miner-fetch/internal/handler/telegram"
	"miner-fetch/internal/storage"
	telegramUsecase "miner-fetch/internal/usecase/telegram"
	"sync"
)

func TelegramBotService(
	ctx context.Context,
	wg *sync.WaitGroup,
	storage *storage.Storage,
	cfg config.Config,
) {
	defer wg.Done()

	usecase := &telegramUsecase.Usecase{Storage: storage}
	handler := &telegram.Handler{Usecase: usecase, Config: cfg}

	opts := []bot.Option{
		bot.WithMiddlewares(handler.AuthMiddleware),
		bot.WithMessageTextHandler("/start", bot.MatchTypeExact, handler.Start),
		bot.WithMessageTextHandler("/info", bot.MatchTypeExact, handler.Info),
	}

	b, err := bot.New(cfg.TelegramAPIKey, opts...)

	if err != nil {
		log.Fatalln(errors.Wrap(err, "telegram bot"))
	}

	b.Start(ctx)

	log.Println("Service 'TelegramBot' has stopped")
}
