package memory

import (
	"context"

	"github.com/manpasik/backend/services/payment-service/internal/service"
)

// EventPublisher는 인메모리 이벤트 발행기입니다 (개발용, 실제는 Kafka).
type EventPublisher struct{}

// NewEventPublisher는 인메모리 EventPublisher를 생성합니다.
func NewEventPublisher() *EventPublisher {
	return &EventPublisher{}
}

func (p *EventPublisher) PublishPaymentCompleted(_ context.Context, _ *service.PaymentCompletedEvent) error {
	return nil
}

func (p *EventPublisher) PublishPaymentFailed(_ context.Context, _ *service.PaymentFailedEvent) error {
	return nil
}

func (p *EventPublisher) PublishPaymentRefunded(_ context.Context, _ *service.PaymentRefundedEvent) error {
	return nil
}
