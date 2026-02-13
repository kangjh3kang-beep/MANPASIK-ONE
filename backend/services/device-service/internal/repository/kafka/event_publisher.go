// Package kafka는 Kafka 기반 EventPublisher 구현을 제공합니다.
package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/manpasik/backend/services/device-service/internal/service"
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

// PublishDeviceRegistered는 디바이스 등록 이벤트를 Kafka에 발행합니다.
func (p *EventPublisher) PublishDeviceRegistered(ctx context.Context, event *service.DeviceRegisteredEvent) error {
	payload := map[string]interface{}{
		"device_id":        event.DeviceID,
		"serial_number":    event.SerialNumber,
		"firmware_version": event.FirmwareVersion,
		"registered_at":    event.RegisteredAt.Format(time.RFC3339),
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("이벤트 페이로드 직렬화 실패: %w", err)
	}

	kafkaEvent := events.Event{
		Type: "device.registered",
		Payload: map[string]interface{}{
			"event_id":   uuid.New().String(),
			"event_type": "manpasik.device.registered",
			"version":    "1.0",
			"timestamp":  time.Now().UTC().Format(time.RFC3339),
			"source":     "device-service",
			"user_id":    event.UserID,
			"payload":    json.RawMessage(payloadBytes),
		},
	}

	return p.eventBus.Publish(ctx, kafkaEvent)
}

// PublishDeviceStatusChanged는 디바이스 상태 변경 이벤트를 Kafka에 발행합니다.
func (p *EventPublisher) PublishDeviceStatusChanged(ctx context.Context, event *service.DeviceStatusChangedEvent) error {
	payload := map[string]interface{}{
		"device_id":        event.DeviceID,
		"serial_number":    event.SerialNumber,
		"previous_status":  event.PreviousStatus,
		"new_status":       event.NewStatus,
		"battery_percent":  event.BatteryPercent,
		"firmware_version": event.FirmwareVersion,
		"last_seen":        event.LastSeen.Format(time.RFC3339),
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("이벤트 페이로드 직렬화 실패: %w", err)
	}

	kafkaEvent := events.Event{
		Type: "device.status.changed",
		Payload: map[string]interface{}{
			"event_id":   uuid.New().String(),
			"event_type": "manpasik.device.status.changed",
			"version":    "1.0",
			"timestamp":  time.Now().UTC().Format(time.RFC3339),
			"source":     "device-service",
			"user_id":    event.UserID,
			"payload":    json.RawMessage(payloadBytes),
		},
	}

	return p.eventBus.Publish(ctx, kafkaEvent)
}
