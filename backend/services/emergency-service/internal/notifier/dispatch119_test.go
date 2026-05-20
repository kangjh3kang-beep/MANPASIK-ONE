package notifier_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/manpasik/backend/services/emergency-service/internal/notifier"
)

func TestDispatch119Notifier_Dispatch_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Error("Authorization 헤더 누락")
		}
		if r.URL.Path != "/v1/dispatch" {
			t.Errorf("Path = %q, want %q", r.URL.Path, "/v1/dispatch")
		}
		resp := map[string]interface{}{
			"success":       true,
			"dispatch_id":   "D-2026-0001",
			"estimated_eta": 8,
			"message":       "접수 완료",
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	n := notifier.NewDispatch119Notifier(notifier.Dispatch119Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
	})
	ctx := context.Background()

	alert := &notifier.EmergencyAlert{
		EmergencyID: "em-001",
		UserID:      "user-1",
		Type:        "cardiac",
		Severity:    4,
		Location:    "서울 강남구 삼성동",
		Description: "흉통이 심합니다",
		MedicalInfo: "고혈압 약 복용 중",
	}

	result, err := n.Dispatch119(ctx, alert)
	if err != nil {
		t.Fatalf("Dispatch119 실패: %v", err)
	}
	if !result.Success {
		t.Error("Success가 true여야 합니다")
	}
	if result.DispatchID != "D-2026-0001" {
		t.Errorf("DispatchID = %q, want %q", result.DispatchID, "D-2026-0001")
	}
	if result.EstimatedETA != 8 {
		t.Errorf("EstimatedETA = %d, want 8", result.EstimatedETA)
	}
}

func TestDispatch119Notifier_Dispatch_EmptyAPIKey(t *testing.T) {
	n := notifier.NewDispatch119Notifier(notifier.Dispatch119Config{
		APIKey: "",
	})
	ctx := context.Background()

	_, err := n.Dispatch119(ctx, &notifier.EmergencyAlert{})
	if err == nil {
		t.Fatal("빈 API 키에 에러가 반환되어야 합니다")
	}
}

func TestDispatch119Notifier_Dispatch_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()

	n := notifier.NewDispatch119Notifier(notifier.Dispatch119Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
	})
	ctx := context.Background()

	_, err := n.Dispatch119(ctx, &notifier.EmergencyAlert{Type: "cardiac"})
	if err == nil {
		t.Fatal("HTTP 503에 에러가 반환되어야 합니다")
	}
}

func TestDispatch119Notifier_NotifyContacts_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/notify" {
			t.Errorf("Path = %q, want %q", r.URL.Path, "/v1/notify")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n := notifier.NewDispatch119Notifier(notifier.Dispatch119Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
	})
	ctx := context.Background()

	alert := &notifier.EmergencyAlert{
		EmergencyID: "em-002",
		UserID:      "user-1",
		Location:    "서울",
	}

	err := n.NotifyContacts(ctx, alert, []string{"010-1234-5678", "010-8765-4321"})
	if err != nil {
		t.Fatalf("NotifyContacts 실패: %v", err)
	}
}

func TestDispatch119Notifier_NotifyContacts_EmptyAPIKey(t *testing.T) {
	n := notifier.NewDispatch119Notifier(notifier.Dispatch119Config{
		APIKey: "",
	})
	ctx := context.Background()

	err := n.NotifyContacts(ctx, &notifier.EmergencyAlert{}, []string{"010-0000-0000"})
	if err == nil {
		t.Fatal("빈 API 키에 에러가 반환되어야 합니다")
	}
}

func TestDispatch119Notifier_ProviderName(t *testing.T) {
	n := notifier.NewDispatch119Notifier(notifier.Dispatch119Config{
		APIKey: "test",
	})
	if n.ProviderName() != "dispatch119" {
		t.Errorf("ProviderName = %q, want %q", n.ProviderName(), "dispatch119")
	}
}
