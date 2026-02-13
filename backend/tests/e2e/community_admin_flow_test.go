//go:build integration

package e2e

import (
	"context"
	"testing"

	"github.com/manpasik/backend/shared/events"
)

// TestCommunityFlow_PostToNotification verifies community post creation triggers
// a notification to followers.
func TestCommunityFlow_PostToNotification(t *testing.T) {
	bus := events.NewEventBus()
	ctx := context.Background()

	steps := make([]string, 0, 3)

	// Step 1: Post created → comment created (simulating immediate self-comment)
	bus.Subscribe(events.EventCommunityPostCreated, func(ctx context.Context, e events.Event) error {
		steps = append(steps, "post_created")
		t.Logf("✅ Event: %s, author=%v title=%v", e.Type, e.Payload["author_id"], e.Payload["title"])

		return bus.Publish(ctx, events.Event{
			Type: events.EventCommunityCommentCreated,
			Payload: map[string]interface{}{
				"post_id":   e.Payload["post_id"],
				"author_id": "user-commenter-001",
				"body":      "좋은 글이네요!",
			},
		})
	})

	// Step 2: Comment created → notification to post author
	bus.Subscribe(events.EventCommunityCommentCreated, func(ctx context.Context, e events.Event) error {
		steps = append(steps, "comment_created")
		t.Logf("✅ Event: %s, post=%v commenter=%v", e.Type, e.Payload["post_id"], e.Payload["author_id"])

		return bus.Publish(ctx, events.Event{
			Type: events.EventNotificationSent,
			Payload: map[string]interface{}{
				"type":    "community_comment",
				"post_id": e.Payload["post_id"],
				"message": "새 댓글이 달렸습니다.",
			},
		})
	})

	// Step 3: Notification sent (terminal)
	bus.Subscribe(events.EventNotificationSent, func(_ context.Context, e events.Event) error {
		steps = append(steps, "notification_sent")
		t.Logf("✅ Event: %s, type=%v", e.Type, e.Payload["type"])
		return nil
	})

	// Kick off
	err := bus.Publish(ctx, events.Event{
		Type: events.EventCommunityPostCreated,
		Payload: map[string]interface{}{
			"post_id":   "post-001",
			"author_id": "user-community-test",
			"title":     "첫 번째 게시글",
		},
	})
	if err != nil {
		t.Fatalf("publish CommunityPostCreated failed: %v", err)
	}

	expected := []string{"post_created", "comment_created", "notification_sent"}
	if len(steps) != len(expected) {
		t.Fatalf("expected %d steps, got %d: %v", len(expected), len(steps), steps)
	}
	for i, want := range expected {
		if steps[i] != want {
			t.Errorf("step %d: got %q, want %q", i, steps[i], want)
		}
	}

	t.Logf("✅ Community flow: post → comment → notification (%d steps)", len(steps))
}

// TestAdminFlow_ActionToNotification verifies that admin actions trigger
// appropriate audit events and notifications.
func TestAdminFlow_ActionToNotification(t *testing.T) {
	bus := events.NewEventBus()
	ctx := context.Background()

	received := make(map[string]bool)

	// Admin action → notification to affected user
	bus.Subscribe(events.EventAdminActionPerformed, func(ctx context.Context, e events.Event) error {
		received["admin_action"] = true
		t.Logf("✅ Admin action: action=%v target=%v", e.Payload["action"], e.Payload["target_user_id"])

		return bus.Publish(ctx, events.Event{
			Type: events.EventNotificationSent,
			Payload: map[string]interface{}{
				"type":    "admin_action",
				"user_id": e.Payload["target_user_id"],
				"message": "관리자에 의해 계정이 정지되었습니다.",
			},
		})
	})

	bus.Subscribe(events.EventNotificationSent, func(_ context.Context, e events.Event) error {
		received["notification_sent"] = true
		t.Logf("✅ Notification: type=%v user=%v", e.Payload["type"], e.Payload["user_id"])
		return nil
	})

	err := bus.Publish(ctx, events.Event{
		Type: events.EventAdminActionPerformed,
		Payload: map[string]interface{}{
			"admin_id":       "admin-001",
			"action":         "suspend_user",
			"target_user_id": "user-bad-actor",
			"reason":         "커뮤니티 가이드라인 위반",
		},
	})
	if err != nil {
		t.Fatalf("publish AdminActionPerformed failed: %v", err)
	}

	if !received["admin_action"] {
		t.Error("admin_action event not received")
	}
	if !received["notification_sent"] {
		t.Error("notification_sent event not received")
	}

	t.Logf("✅ Admin flow: admin action → notification completed")
}

// TestConsentFlow_GrantAndRevoke verifies consent grant/revoke with family data sharing.
func TestConsentFlow_GrantAndRevoke(t *testing.T) {
	bus := events.NewEventBus()
	ctx := context.Background()

	steps := make([]string, 0, 3)

	// Consent granted → share family data
	bus.Subscribe(events.EventConsentGranted, func(ctx context.Context, e events.Event) error {
		steps = append(steps, "consent_granted")
		t.Logf("✅ Consent granted: user=%v scope=%v", e.Payload["user_id"], e.Payload["scope"])

		return bus.Publish(ctx, events.Event{
			Type: events.EventFamilyDataShared,
			Payload: map[string]interface{}{
				"user_id":   e.Payload["user_id"],
				"family_id": "family-001",
				"scope":     e.Payload["scope"],
			},
		})
	})

	bus.Subscribe(events.EventFamilyDataShared, func(_ context.Context, e events.Event) error {
		steps = append(steps, "family_data_shared")
		t.Logf("✅ Family data shared: family=%v scope=%v", e.Payload["family_id"], e.Payload["scope"])
		return nil
	})

	// Consent revoked (separate flow)
	bus.Subscribe(events.EventConsentRevoked, func(_ context.Context, e events.Event) error {
		steps = append(steps, "consent_revoked")
		t.Logf("✅ Consent revoked: user=%v", e.Payload["user_id"])
		return nil
	})

	// Flow 1: Grant consent → share data
	err := bus.Publish(ctx, events.Event{
		Type: events.EventConsentGranted,
		Payload: map[string]interface{}{
			"user_id": "user-consent-test",
			"scope":   "measurement_history",
		},
	})
	if err != nil {
		t.Fatalf("publish ConsentGranted failed: %v", err)
	}

	// Flow 2: Revoke consent
	err = bus.Publish(ctx, events.Event{
		Type: events.EventConsentRevoked,
		Payload: map[string]interface{}{
			"user_id": "user-consent-test",
			"scope":   "measurement_history",
		},
	})
	if err != nil {
		t.Fatalf("publish ConsentRevoked failed: %v", err)
	}

	expected := []string{"consent_granted", "family_data_shared", "consent_revoked"}
	if len(steps) != len(expected) {
		t.Fatalf("expected %d steps, got %d: %v", len(expected), len(steps), steps)
	}
	for i, want := range expected {
		if steps[i] != want {
			t.Errorf("step %d: got %q, want %q", i, steps[i], want)
		}
	}

	t.Logf("✅ Consent flow: grant → share → revoke completed (%d steps)", len(steps))
}
