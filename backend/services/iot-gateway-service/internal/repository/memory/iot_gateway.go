// Package memory는 인메모리 IoT 저장소입니다 (개발/테스트용).
package memory

import (
	"context"
	"sync"

	"github.com/manpasik/backend/services/iot-gateway-service/internal/service"
)

// IoTRepository는 인메모리 IoT 저장소입니다.
type IoTRepository struct {
	mu       sync.RWMutex
	devices  map[string]*service.IoTDevice  // key: DeviceID
	commands map[string]*service.IoTCommand // key: Command ID
	data     []*service.IoTData
}

// NewIoTRepository는 인메모리 IoTRepository를 생성합니다.
func NewIoTRepository() *IoTRepository {
	return &IoTRepository{
		devices:  make(map[string]*service.IoTDevice),
		commands: make(map[string]*service.IoTCommand),
		data:     make([]*service.IoTData, 0),
	}
}

// RegisterDevice는 디바이스를 등록합니다.
func (r *IoTRepository) RegisterDevice(_ context.Context, device *service.IoTDevice) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	cp := *device
	if cp.Metadata != nil {
		cp.Metadata = make(map[string]string, len(device.Metadata))
		for k, v := range device.Metadata {
			cp.Metadata[k] = v
		}
	}
	r.devices[device.DeviceID] = &cp
	return nil
}

// GetDevice는 디바이스 ID로 조회합니다.
func (r *IoTRepository) GetDevice(_ context.Context, deviceID string) (*service.IoTDevice, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	d, ok := r.devices[deviceID]
	if !ok {
		return nil, nil
	}
	cp := *d
	if d.Metadata != nil {
		cp.Metadata = make(map[string]string, len(d.Metadata))
		for k, v := range d.Metadata {
			cp.Metadata[k] = v
		}
	}
	return &cp, nil
}

// SendCommand는 명령을 저장합니다.
func (r *IoTRepository) SendCommand(_ context.Context, cmd *service.IoTCommand) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	cp := *cmd
	r.commands[cmd.ID] = &cp
	return nil
}

// GetCommand는 명령 ID로 조회합니다.
func (r *IoTRepository) GetCommand(_ context.Context, commandID string) (*service.IoTCommand, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	c, ok := r.commands[commandID]
	if !ok {
		return nil, nil
	}
	cp := *c
	return &cp, nil
}

// ReceiveData는 수신 데이터를 저장합니다.
func (r *IoTRepository) ReceiveData(_ context.Context, data *service.IoTData) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	cp := *data
	r.data = append(r.data, &cp)
	return nil
}

// ListData는 디바이스의 수신 데이터 목록을 조회합니다.
func (r *IoTRepository) ListData(_ context.Context, deviceID string, limit, offset int) ([]*service.IoTData, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var filtered []*service.IoTData
	for _, d := range r.data {
		if d.DeviceID == deviceID {
			cp := *d
			filtered = append(filtered, &cp)
		}
	}

	total := len(filtered)
	if offset >= total {
		return nil, nil
	}
	end := offset + limit
	if end > total {
		end = total
	}

	return filtered[offset:end], nil
}
