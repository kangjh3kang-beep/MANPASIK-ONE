// Package memory는 device-service의 인메모리 저장소 구현입니다.
package memory

import (
	"context"
	"sync"
	"time"

	"github.com/manpasik/backend/services/device-service/internal/service"
)

// DeviceRepository는 인메모리 디바이스 저장소입니다.
type DeviceRepository struct {
	mu      sync.RWMutex
	devices map[string]*service.Device // key: device ID
}

// NewDeviceRepository는 인메모리 DeviceRepository를 생성합니다.
func NewDeviceRepository() *DeviceRepository {
	return &DeviceRepository{
		devices: make(map[string]*service.Device),
	}
}

func (r *DeviceRepository) Create(_ context.Context, device *service.Device) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.devices[device.ID] = device
	return nil
}

func (r *DeviceRepository) GetByID(_ context.Context, deviceID string) (*service.Device, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	d, ok := r.devices[deviceID]
	if !ok {
		return nil, nil
	}
	cp := *d
	return &cp, nil
}

func (r *DeviceRepository) ListByUser(_ context.Context, userID string) ([]*service.Device, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*service.Device
	for _, d := range r.devices {
		if d.UserID == userID {
			cp := *d
			result = append(result, &cp)
		}
	}
	return result, nil
}

func (r *DeviceRepository) UpdateStatus(_ context.Context, deviceID string, status service.DeviceStatus, battery int, lastSeen time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	d, ok := r.devices[deviceID]
	if !ok {
		return nil
	}
	d.Status = status
	d.BatteryPercent = battery
	d.LastSeen = lastSeen
	return nil
}

func (r *DeviceRepository) CountByUser(_ context.Context, userID string) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	count := 0
	for _, d := range r.devices {
		if d.UserID == userID {
			count++
		}
	}
	return count, nil
}

func (r *DeviceRepository) Delete(_ context.Context, deviceID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.devices, deviceID)
	return nil
}
