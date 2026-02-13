//go:build integration

package e2e

import (
	"context"
	"testing"

	"github.com/manpasik/backend/shared/events"
)

// TestCommerceFlow_SubscriptionToPayment verifies the subscription → shop → payment
// event chain using the in-process EventBus.
func TestCommerceFlow_SubscriptionToPayment(t *testing.T) {
	bus := events.NewEventBus()
	ctx := context.Background()

	received := make(map[string]bool)

	// When a subscription is created, simulate a shop order.
	bus.Subscribe(events.EventSubscriptionCreated, func(ctx context.Context, e events.Event) error {
		received["subscription_created"] = true
		t.Logf("✅ Event received: %s, payload: %v", e.Type, e.Payload)

		// Trigger shop order as a downstream effect
		return bus.Publish(ctx, events.Event{
			Type: events.EventShopOrderCreated,
			Payload: map[string]interface{}{
				"subscription_id": e.Payload["subscription_id"],
				"user_id":         e.Payload["user_id"],
				"plan":            e.Payload["plan"],
			},
		})
	})

	// When a shop order is created, simulate payment completion.
	bus.Subscribe(events.EventShopOrderCreated, func(ctx context.Context, e events.Event) error {
		received["shop_order_created"] = true
		t.Logf("✅ Event received: %s, payload: %v", e.Type, e.Payload)

		return bus.Publish(ctx, events.Event{
			Type: events.EventPaymentCompleted,
			Payload: map[string]interface{}{
				"subscription_id": e.Payload["subscription_id"],
				"amount":          "29900",
			},
		})
	})

	// Terminal event — just record it.
	bus.Subscribe(events.EventPaymentCompleted, func(_ context.Context, e events.Event) error {
		received["payment_completed"] = true
		t.Logf("✅ Event received: %s, payload: %v", e.Type, e.Payload)
		return nil
	})

	// Kick off the flow
	err := bus.Publish(ctx, events.Event{
		Type: events.EventSubscriptionCreated,
		Payload: map[string]interface{}{
			"subscription_id": "sub-001",
			"user_id":         "user-commerce-test",
			"plan":            "premium",
		},
	})
	if err != nil {
		t.Fatalf("publish SubscriptionCreated failed: %v", err)
	}

	// Verify all three events were processed
	if !received["subscription_created"] {
		t.Error("subscription_created event not received")
	}
	if !received["shop_order_created"] {
		t.Error("shop_order_created event not received")
	}
	if !received["payment_completed"] {
		t.Error("payment_completed event not received")
	}

	t.Logf("✅ Commerce flow: subscription → shop → payment completed successfully")
}

// TestCommerceFlow_SubscriptionCancellation verifies the cancellation flow.
func TestCommerceFlow_SubscriptionCancellation(t *testing.T) {
	bus := events.NewEventBus()
	ctx := context.Background()

	received := make(map[string]bool)

	bus.Subscribe(events.EventSubscriptionCancelled, func(_ context.Context, e events.Event) error {
		received["subscription_cancelled"] = true
		t.Logf("✅ Event received: %s, payload: %v", e.Type, e.Payload)
		return nil
	})

	err := bus.Publish(ctx, events.Event{
		Type: events.EventSubscriptionCancelled,
		Payload: map[string]interface{}{
			"subscription_id": "sub-002",
			"user_id":         "user-cancel-test",
			"reason":          "user_requested",
		},
	})
	if err != nil {
		t.Fatalf("publish SubscriptionCancelled failed: %v", err)
	}

	if !received["subscription_cancelled"] {
		t.Error("subscription_cancelled event not received")
	}

	t.Logf("✅ Subscription cancellation flow completed")
}
