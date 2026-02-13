package memory

import (
	"context"
	"sync"

	"github.com/manpasik/backend/services/measurement-service/internal/service"
)

// EventPublisher는 인메모리 이벤트 발행기입니다 (개발용, 실제는 Kafka).
type EventPublisher struct {
	mu     sync.Mutex
	events []*service.MeasurementCompletedEvent
}

// NewEventPublisher는 인메모리 EventPublisher를 생성합니다.
func NewEventPublisher() *EventPublisher {
	return &EventPublisher{}
}

func (p *EventPublisher) PublishMeasurementCompleted(_ context.Context, event *service.MeasurementCompletedEvent) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.events = append(p.events, event)
	return nil
}
