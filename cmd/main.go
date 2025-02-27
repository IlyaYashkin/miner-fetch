package main

import (
	"miner-fetch/internal/app"
)

func main() {
	a := app.NewApp()

	a.Start()

	a.HandleShutdown()

	a.Stop()
}
