package runner

import (
	"context"
	"fmt"
	"miner-fetch/internal/device"
	"time"
)

var period = 10 * time.Second

type Scanner interface {
	Scan(ctx context.Context) (map[string]device.Device, error)
}

type DeviceScanner struct {
	CommonRunner
	scanner Scanner
}

func NewDeviceScanner(runner CommonRunner) *DeviceScanner {
	ctxc, cancel := context.WithCancel(runner.ctx)
	runner.ctx = ctxc
	runner.cancel = cancel

	return &DeviceScanner{
		CommonRunner: runner,
		scanner:      device.GetRustScanScanner(),
	}
}

func (d *DeviceScanner) Start() {
	go func() {
		timer := time.NewTimer(0)

		defer timer.Stop()
	L:
		for {
			select {
			case <-timer.C:
				timer.Reset(period)

				devices, err := d.Scan()
				if err != nil {
					d.s.Logger.Log(err)
					continue
				}

				d.s.Device.SetDevices(devices)

				d.NotifyUsers()

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

func (d *DeviceScanner) Scan() (map[string]device.Device, error) {
	devices, err := d.scanner.Scan(d.ctx)
	if err != nil {
		return nil, err
	}

	for k, dev := range devices {
		versionCommand := device.VersionCommand{}
		err := dev.SendCommand(&versionCommand)

		poolsCommand := device.PoolsCommand{}
		err = dev.SendCommand(&poolsCommand)
		if err != nil {
			dev.Info = "unsupported device"
		}

		info := fmt.Sprintf(
			"%s \\[%s] %s",
			versionCommand.Response.Version[0].Type,
			poolsCommand.Response.Pools[0].User,
			dev.IP,
		)

		dev.Info = info

		devices[k] = dev
	}

	return devices, nil
}

func (d *DeviceScanner) NotifyUsers() {
	offlineDev := d.s.Device.GetOfflineDevices()
	newDev := d.s.Device.GetNewDevices()

	var message string

	if len(offlineDev) > 0 {
		message += " ❗️*НЕ Работают* ❗️\n\n"
		for _, dev := range offlineDev {
			message += dev.Info + "\n"
		}
	}

	if len(newDev) > 0 {
		message += " ❗️*Новые* ❗️\n\n"
		for _, dev := range newDev {
			message += dev.Info + "\n"
		}
	}

	if message != "" {
		for _, chatId := range d.s.TelegramSender.GetChatIds() {
			err := d.s.TelegramSender.SendMessage(d.ctx, chatId, d.cfg.NodeName, message)
			if err != nil {
				d.s.Logger.Log(err)
			}
		}
	}
}
