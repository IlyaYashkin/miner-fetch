package app

import (
	"context"
	"fmt"
	"log"
	"miner-fetch/internal/config"
	"miner-fetch/internal/runner"
	"miner-fetch/internal/service"
	"os"
	"os/signal"
	"syscall"
)

type Runner interface {
	Start()
	Stop()
	GetName() string
}

type App struct {
	runners []Runner
}

func NewApp() *App {
	app := &App{}

	cfg := config.GetConfig()

	ctx := context.Background()

	s := &service.Service{
		Device:  service.NewDevice(),
		Polling: service.NewPolling(),
		Logger:  service.NewLogger(),
	}

	commonRunner := runner.NewCommonRunner(ctx, s, cfg)

	if cfg.IsScanner {
		app.runners = append(app.runners, runner.NewDeviceScanner(commonRunner))
	}

	if cfg.Mode == "parent" {
		app.runners = append(app.runners, runner.NewTelegramBot(commonRunner))
		app.runners = append(app.runners, runner.NewHttpServer(commonRunner))
	} else if cfg.Mode == "child" {
		app.runners = append(app.runners, runner.NewPoller(commonRunner))
	}

	app.runners = append(app.runners, runner.NewLogger(commonRunner))

	return app
}

func (a *App) Start() {
	for _, r := range a.runners {
		r.Start()
	}
}

func (a *App) HandleShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
}

func (a *App) Stop() {
	fmt.Println()
	log.Printf("Shutting down...")

	for _, r := range a.runners {
		log.Printf("Stopping service '%s'...\n", r.GetName())
		r.Stop()
	}

	log.Println("All runners are stopped")
}
