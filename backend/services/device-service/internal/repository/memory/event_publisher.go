package memory

import (
	"context"

	"github.com/manpasik/backend/services/device-service/internal/service"
)

// EventPublisher는 인메모리 이벤트 발행기입니다 (개발용, 실제는 Kafka).
type EventPublisher struct{}

// NewEventPublisher는 인메모리 EventPublisher를 생성합니다.
func NewEventPublisher() *EventPublisher {
	return &EventPublisher{}
}

func (p *EventPublisher) PublishDeviceRegistered(_ context.Context, _ *service.DeviceRegisteredEvent) error {
	return nil
}

func (p *EventPublisher) PublishDeviceStatusChanged(_ context.Context, _ *service.DeviceStatusChangedEvent) error {
	return nil
}
