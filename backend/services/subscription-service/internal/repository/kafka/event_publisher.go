// Package kafka는 Kafka 기반 EventPublisher 구현을 제공합니다.
package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/manpasik/backend/services/subscription-service/internal/service"
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

// PublishSubscriptionChanged는 구독 변경 이벤트를 Kafka에 발행합니다.
func (p *EventPublisher) PublishSubscriptionChanged(ctx context.Context, event *service.SubscriptionChangedEvent) error {
	payload := map[string]interface{}{
		"subscription_id":      event.SubscriptionID,
		"previous_tier":        event.PreviousTier,
		"new_tier":             event.NewTier,
		"change_type":          event.ChangeType,
		"effective_at":         event.EffectiveAt.Format(time.RFC3339),
		"max_devices":          event.MaxDevices,
		"max_family_members":   event.MaxFamilyMembers,
		"ai_coaching_enabled":  event.AICoachingEnabled,
		"telemedicine_enabled": event.TelemedicineEnabled,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("이벤트 페이로드 직렬화 실패: %w", err)
	}

	kafkaEvent := events.Event{
		Type: "subscription.changed",
		Payload: map[string]interface{}{
			"event_id":   uuid.New().String(),
			"event_type": "manpasik.subscription.changed",
			"version":    "1.0",
			"timestamp":  time.Now().UTC().Format(time.RFC3339),
			"source":     "subscription-service",
			"user_id":    event.UserID,
			"payload":    json.RawMessage(payloadBytes),
		},
	}

	return p.eventBus.Publish(ctx, kafkaEvent)
}
