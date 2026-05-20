package notifier_test

import (
	"context"
	"testing"

	"github.com/manpasik/backend/services/emergency-service/internal/notifier"
)

func TestLogNotifier_Dispatch119(t *testing.T) {
	n := notifier.NewLogNotifier()
	ctx := context.Background()

	alert := &notifier.EmergencyAlert{
		EmergencyID: "em-001",
		UserID:      "user-1",
		Type:        "cardiac",
		Severity:    4,
		Location:    "서울 강남구",
		Description: "흉통 발생",
	}

	result, err := n.Dispatch119(ctx, alert)
	if err != nil {
		t.Fatalf("LogNotifier.Dispatch119 에러: %v", err)
	}
	if result == nil {
		t.Fatal("result가 nil입니다")
	}
	if !result.Success {
		t.Error("Success가 true여야 합니다")
	}
	if result.DispatchID != "LOG-em-001" {
		t.Errorf("DispatchID = %q, want %q", result.DispatchID, "LOG-em-001")
	}
}

func TestLogNotifier_NotifyContacts(t *testing.T) {
	n := notifier.NewLogNotifier()
	ctx := context.Background()

	alert := &notifier.EmergencyAlert{EmergencyID: "em-002"}
	err := n.NotifyContacts(ctx, alert, []string{"010-1234-5678"})
	if err != nil {
		t.Fatalf("NotifyContacts 에러: %v", err)
	}
}

func TestLogNotifier_ProviderName(t *testing.T) {
	n := notifier.NewLogNotifier()
	if n.ProviderName() != "log" {
		t.Errorf("ProviderName = %q, want %q", n.ProviderName(), "log")
	}
}
