package events

import (
	"context"
	"sync"
)

// EventPublisher is the interface for publishing and subscribing to events.
// Both the in-memory EventBus and KafkaEventBus satisfy this interface.
type EventPublisher interface {
	Subscribe(eventType string, handler Handler)
	Publish(ctx context.Context, event Event) error
}

// Event represents a domain event
type Event struct {
	Type    string
	Payload map[string]interface{}
}

// Handler is a function that handles an event
type Handler func(ctx context.Context, event Event) error

// EventBus provides publish-subscribe messaging
type EventBus struct {
	mu       sync.RWMutex
	handlers map[string][]Handler
}

// NewEventBus creates a new event bus
func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[string][]Handler),
	}
}

// Subscribe registers a handler for an event type
func (b *EventBus) Subscribe(eventType string, handler Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[eventType] = append(b.handlers[eventType], handler)
}

// Publish sends an event to all subscribed handlers
func (b *EventBus) Publish(ctx context.Context, event Event) error {
	b.mu.RLock()
	handlers := b.handlers[event.Type]
	b.mu.RUnlock()

	for _, h := range handlers {
		if err := h(ctx, event); err != nil {
			return err
		}
	}
	return nil
}

// Event types constants
const (
	// Reservation
	EventReservationCreated   = "reservation.created"
	EventReservationCancelled = "reservation.cancelled"

	// Prescription
	EventPrescriptionCreated        = "prescription.created"
	EventPrescriptionSentToPharmacy = "prescription.sent_to_pharmacy"
	EventPrescriptionDispensed      = "prescription.dispensed"

	// Measurement & Calibration
	EventMeasurementCompleted   = "measurement.completed"
	EventCalibrationCompleted   = "calibration.completed"
	EventCartridgeReplaced      = "cartridge.replaced"

	// AI / Inference
	EventAIAnalysisCompleted = "ai.analysis_completed"

	// Coaching
	EventCoachingGoalSet     = "coaching.goal_set"
	EventCoachingTipDelivered = "coaching.tip_delivered"

	// Health
	EventHealthAlertTriggered = "health_alert.triggered"

	// Commerce
	EventSubscriptionCreated   = "subscription.created"
	EventSubscriptionCancelled = "subscription.cancelled"
	EventPaymentCompleted      = "payment.completed"
	EventShopOrderCreated      = "shop.order_created"

	// Consent / Family
	EventConsentGranted   = "consent.granted"
	EventConsentRevoked   = "consent.revoked"
	EventFamilyDataShared = "family.data_shared"

	// Community
	EventCommunityPostCreated    = "community.post_created"
	EventCommunityCommentCreated = "community.comment_created"

	// Admin
	EventAdminActionPerformed = "admin.action_performed"

	// Notification
	EventNotificationSent = "notification.sent"
)
