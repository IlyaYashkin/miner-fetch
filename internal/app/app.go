package app

import (
	"context"
	"fmt"
	"log"
	"miner-fetch/internal/config"
	"miner-fetch/internal/service"
	"miner-fetch/internal/storage"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func Start() {
	cfg := config.GetConfig()

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	services := []func(
		ctx context.Context,
		wg *sync.WaitGroup,
		storage *storage.Storage,
		cfg config.Config,
	){
		service.TelegramBotService,
		service.DeviceScannerService,
	}
	store := &storage.Storage{}

	for _, srv := range services {
		wg.Add(1)
		go srv(ctx, &wg, store, cfg)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan

	fmt.Println()
	log.Println("Shutting down...")

	cancel()

	wg.Wait()

	log.Println("All services are stopped")
}
