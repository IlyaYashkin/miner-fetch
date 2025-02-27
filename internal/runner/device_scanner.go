package runner

import (
	"context"
	"log"
	"miner-fetch/internal/device"
	"time"
)

var period = 5 * time.Minute

type DeviceScanner struct {
	CommonRunner
}

func NewDeviceScanner(runner CommonRunner) *DeviceScanner {
	ctxc, cancel := context.WithCancel(runner.ctx)
	runner.ctx = ctxc
	runner.cancel = cancel

	return &DeviceScanner{runner}
}

func (d *DeviceScanner) Start() {
	go func() {
		ticker := time.NewTicker(period)

		defer ticker.Stop()

		scanner := device.GetRustScanScanner()

		log.Println("scanning devices...")
		devices, err := scanner.Scan(d.ctx)
		if err != nil {
			d.s.Logger.Log(err)
		} else {
			d.s.Device.SetDevices(devices)
			log.Println("scanning completed")
		}

	L:
		for {
			select {
			case <-ticker.C:
				devices, err := scanner.Scan(d.ctx)
				if err != nil {
					d.s.Logger.Log(err)
				}
				d.s.Device.SetDevices(devices)
			case <-d.ctx.Done():
				d.stopCh <- true
				break L
			}
		}
	}()
}

func (d *DeviceScanner) GetName() string {
	return "DeviceScanner"
}
