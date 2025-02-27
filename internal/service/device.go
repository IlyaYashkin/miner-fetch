package service

import (
	"fmt"
	"miner-fetch/internal/device"
	"sync"
)

type Device struct {
	mu         sync.Mutex
	devices    []device.Device
	commandMap map[string]func() (string, error)
}

func NewDevice() *Device {
	d := &Device{}
	d.commandMap = map[string]func() (string, error){
		"info": d.GetDevicesInfo,
	}

	return d
}

func (d *Device) SetDevices(devices []device.Device) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.devices = devices
}

func (d *Device) GetDevices() []device.Device {
	d.mu.Lock()
	defer d.mu.Unlock()

	return d.devices
}

func (d *Device) ExecuteCommand(command string) (string, error) {
	commandFunc, ok := d.commandMap[command]
	if !ok {
		return "", fmt.Errorf("command %s not found", command)
	}

	message, err := commandFunc()
	return message, err
}

func (d *Device) GetDevicesInfo() (string, error) {
	devices := d.GetDevices()

	var message string

	for _, d := range devices {
		versionCommand := device.VersionCommand{}
		err := d.SendCommand(&versionCommand)
		if err != nil {
			return "", err
		}

		statsCommand := device.StatsCommand{}
		err = d.SendCommand(&statsCommand)
		if err != nil {
			return "", err
		}

		poolsCommand := device.PoolsCommand{}
		err = d.SendCommand(&poolsCommand)
		if err != nil {
			return "", err
		}

		message += fmt.Sprintf(
			"%s [%s]\n",
			versionCommand.Response.Version[0].Type,
			poolsCommand.Response.Pools[0].User,
		)

		message += fmt.Sprintf(
			"Temp 1 — %d %d\n",
			statsCommand.Response.Stats[1].Temp1,
			statsCommand.Response.Stats[1].Temp21,
		)

		message += fmt.Sprintf(
			"Temp 2 — %d %d\n",
			statsCommand.Response.Stats[1].Temp2,
			statsCommand.Response.Stats[1].Temp22,
		)

		message += fmt.Sprintf(
			"Temp 3 — %d %d\n",
			statsCommand.Response.Stats[1].Temp3,
			statsCommand.Response.Stats[1].Temp23,
		)

		message += "\n"
	}

	return message, nil
}
