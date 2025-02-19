package service

import (
	"context"
	"log"
	"miner-fetch/internal/config"
	"miner-fetch/internal/device"
	"miner-fetch/internal/storage"
	"sync"
	"time"
)

func DeviceScannerService(
	ctx context.Context,
	wg *sync.WaitGroup,
	storage *storage.Storage,
	cfg config.Config,
) {
	ticker := time.NewTicker(5 * time.Minute)

	defer ticker.Stop()
	defer wg.Done()

	scanner := device.GetRustScanScanner()

	log.Println("scanning devices...")
	devices, err := scanner.Scan(ctx)
	if err != nil {
		log.Println(err)
	} else {
		storage.SetDevices(devices)
		log.Println("scanning completed")
	}

L:
	for {
		select {
		case <-ticker.C:
			devices, err := scanner.Scan(ctx)
			if err != nil {
				log.Println(err)
			}
			storage.SetDevices(devices)
		case <-ctx.Done():
			break L
		}
	}

	log.Println("Service 'DeviceScanner' has stopped")
}
