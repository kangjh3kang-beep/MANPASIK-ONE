//go:build integration

package e2e

import (
	"context"
	"testing"
	"time"

	"github.com/manpasik/backend/shared/events"
)

func TestMedicalFlowEventIntegration(t *testing.T) {
	bus := events.NewEventBus()

	// Track events received
	receivedEvents := make([]string, 0)

	bus.Subscribe(events.EventReservationCreated, func(ctx context.Context, e events.Event) error {
		receivedEvents = append(receivedEvents, e.Type)
		t.Logf("✅ Event received: %s, payload: %v", e.Type, e.Payload)
		return nil
	})

	bus.Subscribe(events.EventPrescriptionCreated, func(ctx context.Context, e events.Event) error {
		receivedEvents = append(receivedEvents, e.Type)
		t.Logf("✅ Event received: %s, payload: %v", e.Type, e.Payload)
		return nil
	})

	bus.Subscribe(events.EventPrescriptionSentToPharmacy, func(ctx context.Context, e events.Event) error {
		receivedEvents = append(receivedEvents, e.Type)
		t.Logf("✅ Event received: %s, payload: %v", e.Type, e.Payload)
		return nil
	})

	ctx := context.Background()

	// Simulate reservation created
	err := bus.Publish(ctx, events.Event{
		Type: events.EventReservationCreated,
		Payload: map[string]interface{}{
			"reservation_id": "res-001",
			"user_id":        "user-001",
			"facility_name":  "서울대병원",
			"date":           time.Now().Format("2006-01-02"),
		},
	})
	if err != nil {
		t.Fatalf("reservation event publish failed: %v", err)
	}

	// Simulate prescription created after visit
	err = bus.Publish(ctx, events.Event{
		Type: events.EventPrescriptionCreated,
		Payload: map[string]interface{}{
			"prescription_id": "presc-001",
			"user_id":         "user-001",
			"doctor_name":     "김의사",
		},
	})
	if err != nil {
		t.Fatalf("prescription event publish failed: %v", err)
	}

	// Simulate prescription sent to pharmacy
	err = bus.Publish(ctx, events.Event{
		Type: events.EventPrescriptionSentToPharmacy,
		Payload: map[string]interface{}{
			"prescription_id": "presc-001",
			"pharmacy_name":   "녹십자약국",
			"token":           "ABC123",
		},
	})
	if err != nil {
		t.Fatalf("pharmacy event publish failed: %v", err)
	}

	// Verify all events were received
	if len(receivedEvents) != 3 {
		t.Fatalf("expected 3 events, got %d", len(receivedEvents))
	}

	t.Logf("✅ Medical flow integration: %d events processed successfully", len(receivedEvents))
}
