package memory

import (
	"context"
	"sync"

	"github.com/manpasik/backend/services/device-service/internal/service"
)

// DeviceEventRepository는 인메모리 디바이스 이벤트 저장소입니다.
type DeviceEventRepository struct {
	mu     sync.Mutex
	events []*service.DeviceEvent
}

// NewDeviceEventRepository는 인메모리 DeviceEventRepository를 생성합니다.
func NewDeviceEventRepository() *DeviceEventRepository {
	return &DeviceEventRepository{}
}

func (r *DeviceEventRepository) LogEvent(_ context.Context, event *service.DeviceEvent) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.events = append(r.events, event)
	return nil
}
