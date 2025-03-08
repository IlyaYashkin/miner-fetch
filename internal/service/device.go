package service

import (
	"fmt"
	"miner-fetch/internal/device"
	"miner-fetch/internal/util"
	"sync"
)

type CommandNotFoundError struct {
	msg string
}

func (e *CommandNotFoundError) Error() string {
	return e.msg
}

type DeviceService struct {
	mu              sync.Mutex
	devices         map[string]device.Device
	previousDevices map[string]device.Device
	commandMap      map[string]func() (string, error)
}

func NewDeviceService() *DeviceService {
	d := &DeviceService{}
	d.commandMap = map[string]func() (string, error){
		"temperature": d.GetTemperature,
		"status":      d.GetStatus,
		"ips":         d.GetIps,
	}

	return d
}

func (d *DeviceService) SetDevices(devices map[string]device.Device) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.previousDevices = d.devices

	d.devices = devices
}

func (d *DeviceService) GetDevices() map[string]device.Device {
	d.mu.Lock()
	defer d.mu.Unlock()

	return d.devices
}

func (d *DeviceService) GetOfflineDevices() map[string]device.Device {
	d.mu.Lock()
	defer d.mu.Unlock()

	diff := make(map[string]device.Device)

	for k, pd := range d.previousDevices {
		_, ok := d.devices[k]
		if !ok {
			diff[k] = pd
		}
	}

	return diff
}

func (d *DeviceService) GetNewDevices() map[string]device.Device {
	d.mu.Lock()
	defer d.mu.Unlock()

	diff := make(map[string]device.Device)

	for k, pd := range d.devices {
		_, ok := d.previousDevices[k]
		if !ok {
			diff[k] = pd
		}
	}

	return diff
}

func (d *DeviceService) ExecuteCommand(command string) (string, error) {
	commandFunc, ok := d.commandMap[command]
	if !ok {
		return "", &CommandNotFoundError{msg: fmt.Sprintf("command \"%s\" not found", command)}
	}

	message, err := commandFunc()
	return message, err
}

func (d *DeviceService) GetTemperature() (string, error) {
	var message string

	for _, d := range d.GetDevices() {
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

func (d *DeviceService) GetStatus() (string, error) {
	var message string

	for _, d := range d.GetDevices() {
		versionCommand := device.VersionCommand{}
		err := d.SendCommand(&versionCommand)
		if err != nil {
			return "", err
		}

		poolsCommand := device.PoolsCommand{}
		err = d.SendCommand(&poolsCommand)
		if err != nil {
			return "", err
		}

		statsCommand := device.StatsCommand{}
		err = d.SendCommand(&statsCommand)
		if err != nil {
			return "", err
		}

		message += fmt.Sprintf(
			"%s [%s] %s\nHashrate real: %.0f MH/s\nHashrate AVG: %.0f MH/s\n",
			versionCommand.Response.Version[0].Type,
			poolsCommand.Response.Pools[0].User,
			util.FormatDuration(statsCommand.Response.Stats[1].Elapsed),
			statsCommand.Response.Stats[1].Ghs5s,
			statsCommand.Response.Stats[1].GhsAv,
		)

		message += "\n"
	}

	return message, nil
}

func (d *DeviceService) GetIps() (string, error) {
	var message string

	for _, d := range d.GetDevices() {
		versionCommand := device.VersionCommand{}
		err := d.SendCommand(&versionCommand)
		if err != nil {
			return "", err
		}

		poolsCommand := device.PoolsCommand{}
		err = d.SendCommand(&poolsCommand)
		if err != nil {
			return "", err
		}

		message += fmt.Sprintf(
			"%s [%s] %s\n",
			versionCommand.Response.Version[0].Type,
			poolsCommand.Response.Pools[0].User,
			d.IP,
		)
	}

	return message, nil
}
