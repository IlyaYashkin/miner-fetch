package storage

import (
	"miner-fetch/internal/device"
	"sync"
)

type Storage struct {
	deviceMutex sync.Mutex
	devices     []device.Device
}

func (s *Storage) SetDevices(d []device.Device) {
	s.deviceMutex.Lock()
	defer s.deviceMutex.Unlock()

	s.devices = d
}

func (s *Storage) GetDevices() []device.Device {
	s.deviceMutex.Lock()
	defer s.deviceMutex.Unlock()

	return s.devices
}
