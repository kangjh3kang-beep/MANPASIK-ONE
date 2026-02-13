//go:build integration

package e2e

import (
	"context"
	"testing"

	"github.com/manpasik/backend/shared/events"
)

// TestAIHardwareFlow_MeasurementToCoaching verifies the full hardware pipeline:
// measurement.completed → ai.analysis_completed → coaching.tip_delivered
func TestAIHardwareFlow_MeasurementToCoaching(t *testing.T) {
	bus := events.NewEventBus()
	ctx := context.Background()

	steps := make([]string, 0, 3)

	// Step 1: Measurement completes → trigger AI analysis
	bus.Subscribe(events.EventMeasurementCompleted, func(ctx context.Context, e events.Event) error {
		steps = append(steps, "measurement_completed")
		t.Logf("✅ Step 1 — %s: session=%v device=%v", e.Type, e.Payload["session_id"], e.Payload["device_id"])

		return bus.Publish(ctx, events.Event{
			Type: events.EventAIAnalysisCompleted,
			Payload: map[string]interface{}{
				"session_id": e.Payload["session_id"],
				"result":     "normal",
				"confidence": 0.95,
			},
		})
	})

	// Step 2: AI analysis completes → deliver coaching tip
	bus.Subscribe(events.EventAIAnalysisCompleted, func(ctx context.Context, e events.Event) error {
		steps = append(steps, "ai_analysis_completed")
		t.Logf("✅ Step 2 — %s: result=%v confidence=%v", e.Type, e.Payload["result"], e.Payload["confidence"])

		return bus.Publish(ctx, events.Event{
			Type: events.EventCoachingTipDelivered,
			Payload: map[string]interface{}{
				"session_id": e.Payload["session_id"],
				"tip":        "혈당이 정상 범위입니다. 현재 식단을 유지하세요.",
			},
		})
	})

	// Step 3: Coaching tip delivered (terminal)
	bus.Subscribe(events.EventCoachingTipDelivered, func(_ context.Context, e events.Event) error {
		steps = append(steps, "coaching_tip_delivered")
		t.Logf("✅ Step 3 — %s: tip=%v", e.Type, e.Payload["tip"])
		return nil
	})

	// Kick off
	err := bus.Publish(ctx, events.Event{
		Type: events.EventMeasurementCompleted,
		Payload: map[string]interface{}{
			"session_id": "sess-ai-test",
			"device_id":  "dev-001",
			"user_id":    "user-ai-test",
		},
	})
	if err != nil {
		t.Fatalf("publish MeasurementCompleted failed: %v", err)
	}

	// Verify step ordering
	expected := []string{"measurement_completed", "ai_analysis_completed", "coaching_tip_delivered"}
	if len(steps) != len(expected) {
		t.Fatalf("expected %d steps, got %d: %v", len(expected), len(steps), steps)
	}
	for i, want := range expected {
		if steps[i] != want {
			t.Errorf("step %d: got %q, want %q", i, steps[i], want)
		}
	}

	t.Logf("✅ AI/Hardware flow: measurement → analysis → coaching completed (%d steps)", len(steps))
}

// TestAIHardwareFlow_CalibrationAndCartridge verifies calibration + cartridge events.
func TestAIHardwareFlow_CalibrationAndCartridge(t *testing.T) {
	bus := events.NewEventBus()
	ctx := context.Background()

	received := make(map[string]bool)

	bus.Subscribe(events.EventCartridgeReplaced, func(ctx context.Context, e events.Event) error {
		received["cartridge_replaced"] = true
		t.Logf("✅ Event: %s, device=%v", e.Type, e.Payload["device_id"])

		// Cartridge replacement triggers calibration
		return bus.Publish(ctx, events.Event{
			Type: events.EventCalibrationCompleted,
			Payload: map[string]interface{}{
				"device_id":    e.Payload["device_id"],
				"cartridge_id": e.Payload["cartridge_id"],
				"status":       "passed",
			},
		})
	})

	bus.Subscribe(events.EventCalibrationCompleted, func(_ context.Context, e events.Event) error {
		received["calibration_completed"] = true
		t.Logf("✅ Event: %s, status=%v", e.Type, e.Payload["status"])
		return nil
	})

	err := bus.Publish(ctx, events.Event{
		Type: events.EventCartridgeReplaced,
		Payload: map[string]interface{}{
			"device_id":    "dev-002",
			"cartridge_id": "cart-new-001",
		},
	})
	if err != nil {
		t.Fatalf("publish CartridgeReplaced failed: %v", err)
	}

	if !received["cartridge_replaced"] {
		t.Error("cartridge_replaced event not received")
	}
	if !received["calibration_completed"] {
		t.Error("calibration_completed event not received")
	}

	t.Logf("✅ Calibration flow: cartridge replace → calibration completed")
}

// TestAIHardwareFlow_HealthAlert verifies that abnormal AI results trigger a health alert.
func TestAIHardwareFlow_HealthAlert(t *testing.T) {
	bus := events.NewEventBus()
	ctx := context.Background()

	steps := make([]string, 0, 3)

	bus.Subscribe(events.EventAIAnalysisCompleted, func(ctx context.Context, e events.Event) error {
		steps = append(steps, "ai_analysis_completed")

		// Abnormal result triggers health alert
		if e.Payload["result"] == "abnormal" {
			return bus.Publish(ctx, events.Event{
				Type: events.EventHealthAlertTriggered,
				Payload: map[string]interface{}{
					"session_id": e.Payload["session_id"],
					"severity":   "high",
					"message":    "혈당 수치 이상 감지",
				},
			})
		}
		return nil
	})

	bus.Subscribe(events.EventHealthAlertTriggered, func(ctx context.Context, e events.Event) error {
		steps = append(steps, "health_alert_triggered")
		t.Logf("✅ Health alert: severity=%v message=%v", e.Payload["severity"], e.Payload["message"])

		// Alert triggers notification
		return bus.Publish(ctx, events.Event{
			Type: events.EventNotificationSent,
			Payload: map[string]interface{}{
				"type":    "health_alert",
				"message": e.Payload["message"],
			},
		})
	})

	bus.Subscribe(events.EventNotificationSent, func(_ context.Context, e events.Event) error {
		steps = append(steps, "notification_sent")
		t.Logf("✅ Notification sent: type=%v", e.Payload["type"])
		return nil
	})

	err := bus.Publish(ctx, events.Event{
		Type: events.EventAIAnalysisCompleted,
		Payload: map[string]interface{}{
			"session_id": "sess-alert-test",
			"result":     "abnormal",
			"confidence": 0.88,
		},
	})
	if err != nil {
		t.Fatalf("publish AIAnalysisCompleted failed: %v", err)
	}

	expected := []string{"ai_analysis_completed", "health_alert_triggered", "notification_sent"}
	if len(steps) != len(expected) {
		t.Fatalf("expected %d steps, got %d: %v", len(expected), len(steps), steps)
	}
	for i, want := range expected {
		if steps[i] != want {
			t.Errorf("step %d: got %q, want %q", i, steps[i], want)
		}
	}

	t.Logf("✅ Health alert flow: abnormal analysis → alert → notification (%d steps)", len(steps))
}
