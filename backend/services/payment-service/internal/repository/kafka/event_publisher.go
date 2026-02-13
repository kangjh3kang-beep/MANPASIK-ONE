// Package kafka는 Kafka 기반 EventPublisher 구현을 제공합니다.
package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/manpasik/backend/services/payment-service/internal/service"
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

// PublishPaymentCompleted는 결제 완료 이벤트를 Kafka에 발행합니다.
func (p *EventPublisher) PublishPaymentCompleted(ctx context.Context, event *service.PaymentCompletedEvent) error {
	payload := map[string]interface{}{
		"payment_id":        event.PaymentID,
		"order_id":          event.OrderID,
		"amount":            event.AmountKRW,
		"currency":          "KRW",
		"payment_method":    event.PaymentMethod,
		"pg_provider":       event.PgProvider,
		"pg_transaction_id": event.PgTransactionID,
		"payment_type":      event.PaymentType,
		"subscription_id":   event.SubscriptionID,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("이벤트 페이로드 직렬화 실패: %w", err)
	}

	kafkaEvent := events.Event{
		Type: "payment.completed",
		Payload: map[string]interface{}{
			"event_id":   uuid.New().String(),
			"event_type": "manpasik.payment.completed",
			"version":    "1.0",
			"timestamp":  time.Now().UTC().Format(time.RFC3339),
			"source":     "payment-service",
			"user_id":    event.UserID,
			"payload":    json.RawMessage(payloadBytes),
		},
	}

	return p.eventBus.Publish(ctx, kafkaEvent)
}

// PublishPaymentFailed는 결제 실패 이벤트를 Kafka에 발행합니다.
func (p *EventPublisher) PublishPaymentFailed(ctx context.Context, event *service.PaymentFailedEvent) error {
	payload := map[string]interface{}{
		"payment_id":    event.PaymentID,
		"order_id":      event.OrderID,
		"amount":        event.AmountKRW,
		"currency":      "KRW",
		"error_code":    event.ErrorCode,
		"error_message": event.ErrorMessage,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("이벤트 페이로드 직렬화 실패: %w", err)
	}

	kafkaEvent := events.Event{
		Type: "payment.failed",
		Payload: map[string]interface{}{
			"event_id":   uuid.New().String(),
			"event_type": "manpasik.payment.failed",
			"version":    "1.0",
			"timestamp":  time.Now().UTC().Format(time.RFC3339),
			"source":     "payment-service",
			"user_id":    event.UserID,
			"payload":    json.RawMessage(payloadBytes),
		},
	}

	return p.eventBus.Publish(ctx, kafkaEvent)
}

// PublishPaymentRefunded는 환불 완료 이벤트를 Kafka에 발행합니다.
func (p *EventPublisher) PublishPaymentRefunded(ctx context.Context, event *service.PaymentRefundedEvent) error {
	payload := map[string]interface{}{
		"payment_id":      event.PaymentID,
		"refund_id":       event.RefundID,
		"refund_amount":   event.RefundAmountKRW,
		"currency":        "KRW",
		"reason":          event.Reason,
		"is_full_refund":  event.IsFullRefund,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("이벤트 페이로드 직렬화 실패: %w", err)
	}

	kafkaEvent := events.Event{
		Type: "payment.refunded",
		Payload: map[string]interface{}{
			"event_id":   uuid.New().String(),
			"event_type": "manpasik.payment.refunded",
			"version":    "1.0",
			"timestamp":  time.Now().UTC().Format(time.RFC3339),
			"source":     "payment-service",
			"user_id":    event.UserID,
			"payload":    json.RawMessage(payloadBytes),
		},
	}

	return p.eventBus.Publish(ctx, kafkaEvent)
}
