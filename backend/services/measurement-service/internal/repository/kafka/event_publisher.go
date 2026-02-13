// Package kafka는 Kafka 기반 EventPublisher 구현을 제공합니다.
package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/manpasik/backend/services/measurement-service/internal/service"
	"github.com/manpasik/backend/shared/events"
)

// EventPublisher는 Kafka를 사용하는 이벤트 발행기입니다.
type EventPublisher struct {
	eventBus *events.KafkaEventBus
}

// NewEventPublisher는 Kafka 기반 EventPublisher를 생성합니다.
func NewEventPublisher(eventBus *events.KafkaEventBus) *EventPublisher {
	return &EventPublisher{eventBus: eventBus}
}

// PublishMeasurementCompleted는 측정 완료 이벤트를 Kafka에 발행합니다.
func (p *EventPublisher) PublishMeasurementCompleted(ctx context.Context, event *service.MeasurementCompletedEvent) error {
	payload := map[string]interface{}{
		"session_id":    event.SessionID,
		"user_id":       event.UserID,
		"device_id":     event.DeviceID,
		"primary_value": event.PrimaryValue,
		"unit":          event.Unit,
		"completed_at":  event.CompletedAt.Format(time.RFC3339),
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("이벤트 페이로드 직렬화 실패: %w", err)
	}

	kafkaEvent := events.Event{
		Type: "measurement.completed",
		Payload: map[string]interface{}{
			"event_id":   uuid.New().String(),
			"event_type": "manpasik.measurement.completed",
			"version":    "1.0",
			"timestamp":  time.Now().UTC().Format(time.RFC3339),
			"source":     "measurement-service",
			"user_id":    event.UserID,
			"payload":    json.RawMessage(payloadBytes),
		},
	}

	return p.eventBus.Publish(ctx, kafkaEvent)
}
