package fcm_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/manpasik/backend/shared/external/fcm"
)

// TestValidateMessage_Token은 토큰 메시지 검증을 확인합니다.
func TestValidateMessage_Token(t *testing.T) {
	msg := &fcm.Message{
		Token: "device-token-123",
		Notification: &fcm.Notification{Title: "T", Body: "B"},
	}
	if err := fcm.ValidateMessage(msg); err != nil {
		t.Errorf("ValidateMessage 실패: %v", err)
	}
}

// TestValidateMessage_Topic는 토픽 메시지 검증을 확인합니다.
func TestValidateMessage_Topic(t *testing.T) {
	msg := &fcm.Message{
		Topic: "all-users",
		Data:  map[string]string{"key": "value"},
	}
	if err := fcm.ValidateMessage(msg); err != nil {
		t.Errorf("ValidateMessage 실패: %v", err)
	}
}

// TestValidateMessage_NoTarget는 타겟 누락 거부를 검증합니다.
func TestValidateMessage_NoTarget(t *testing.T) {
	msg := &fcm.Message{
		Notification: &fcm.Notification{Title: "T", Body: "B"},
	}
	if err := fcm.ValidateMessage(msg); err == nil {
		t.Error("타겟 없이 통과됨")
	}
}

// TestValidateMessage_MultipleTargets는 다중 타겟 거부를 검증합니다.
func TestValidateMessage_MultipleTargets(t *testing.T) {
	msg := &fcm.Message{
		Token: "t",
		Topic: "x",
		Notification: &fcm.Notification{Title: "T", Body: "B"},
	}
	if err := fcm.ValidateMessage(msg); err == nil {
		t.Error("token+topic 동시 지정이 통과됨")
	}
}

// TestValidateMessage_NoBody는 본문 누락 거부를 검증합니다.
func TestValidateMessage_NoBody(t *testing.T) {
	msg := &fcm.Message{Token: "t"}
	if err := fcm.ValidateMessage(msg); err == nil {
		t.Error("notification/data 없이 통과됨")
	}
}

// TestValidateMessage_EmptyNotification는 빈 알림 거부를 검증합니다.
func TestValidateMessage_EmptyNotification(t *testing.T) {
	msg := &fcm.Message{
		Token:        "t",
		Notification: &fcm.Notification{},
	}
	if err := fcm.ValidateMessage(msg); err == nil {
		t.Error("빈 알림이 통과됨")
	}
}

// TestNoopAdapter_Send는 Noop 발송을 검증합니다.
func TestNoopAdapter_Send(t *testing.T) {
	a := fcm.NewNoopAdapter()
	msg := &fcm.Message{
		Token: "t",
		Notification: &fcm.Notification{Title: "Hi", Body: "Hello"},
	}
	result, err := a.Send(context.Background(), msg)
	if err != nil {
		t.Fatalf("Send 실패: %v", err)
	}
	if result.Provider != "noop" {
		t.Errorf("Provider = %q", result.Provider)
	}
	if a.Count() != 1 {
		t.Errorf("Count = %d, want 1", a.Count())
	}
}

// TestNoopAdapter_SendMulticast는 다중 발송을 검증합니다.
func TestNoopAdapter_SendMulticast(t *testing.T) {
	a := fcm.NewNoopAdapter()
	msg := &fcm.Message{
		Notification: &fcm.Notification{Title: "T", Body: "B"},
	}
	tokens := []string{"t1", "t2", "t3"}

	result, err := a.SendMulticast(context.Background(), msg, tokens)
	if err != nil {
		t.Fatalf("SendMulticast 실패: %v", err)
	}
	if result.SuccessCount != 3 {
		t.Errorf("SuccessCount = %d, want 3", result.SuccessCount)
	}
	if a.Count() != 3 {
		t.Errorf("Count = %d, want 3", a.Count())
	}
}

// TestV1Adapter_Send_Success는 v1 정상 발송을 검증합니다.
func TestV1Adapter_Send_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"name": "projects/test/messages/123"}`))
	}))
	defer server.Close()

	a := fcm.NewV1Adapter("test-project", "")
	a.SetEndpoint(server.URL)
	a.SetTokenProvider(&fcm.StaticTokenProvider{Token: "test-token"})

	msg := &fcm.Message{
		Token: "device-1",
		Notification: &fcm.Notification{Title: "T", Body: "B"},
	}
	result, err := a.Send(context.Background(), msg)
	if err != nil {
		t.Fatalf("Send 실패: %v", err)
	}
	if result.MessageID != "projects/test/messages/123" {
		t.Errorf("MessageID = %q", result.MessageID)
	}
}

// TestV1Adapter_Send_InvalidToken은 무효 토큰 식별을 검증합니다.
func TestV1Adapter_Send_InvalidToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": {"code": 400, "message": "invalid", "status": "UNREGISTERED"}}`))
	}))
	defer server.Close()

	a := fcm.NewV1Adapter("p", "")
	a.SetEndpoint(server.URL)
	a.SetTokenProvider(&fcm.StaticTokenProvider{Token: "tok"})

	msg := &fcm.Message{
		Token:        "expired-token",
		Notification: &fcm.Notification{Title: "T", Body: "B"},
	}
	result, _ := a.Send(context.Background(), msg)
	if result == nil {
		t.Fatal("result nil")
	}
	if result.Status != "invalid_token" {
		t.Errorf("Status = %q, want invalid_token", result.Status)
	}
}

// TestV1Adapter_NoTokenProvider는 토큰 제공자 누락 거부를 검증합니다.
func TestV1Adapter_NoTokenProvider(t *testing.T) {
	a := fcm.NewV1Adapter("p", "")
	msg := &fcm.Message{
		Token:        "t",
		Notification: &fcm.Notification{Title: "T", Body: "B"},
	}
	_, err := a.Send(context.Background(), msg)
	if err == nil {
		t.Error("토큰 제공자 없이 발송 허용됨")
	}
}

// TestV1Adapter_DryRun은 dry_run 모드를 검증합니다.
func TestV1Adapter_DryRun(t *testing.T) {
	gotValidateOnly := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 본문에 validate_only 포함 여부 확인 (대략)
		buf := make([]byte, 1024)
		n, _ := r.Body.Read(buf)
		if n > 0 && contains(string(buf[:n]), "validate_only") {
			gotValidateOnly = true
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"name": "x"}`))
	}))
	defer server.Close()

	a := fcm.NewV1Adapter("p", "")
	a.SetEndpoint(server.URL)
	a.SetTokenProvider(&fcm.StaticTokenProvider{Token: "tok"})

	msg := &fcm.Message{
		Token:        "t",
		Notification: &fcm.Notification{Title: "T", Body: "B"},
		DryRun:       true,
	}
	_, _ = a.Send(context.Background(), msg)
	if !gotValidateOnly {
		t.Error("dry_run 시 validate_only가 요청에 포함되지 않음")
	}
}

// TestV1Adapter_AndroidConfig는 안드로이드 설정 빌드를 검증합니다.
func TestV1Adapter_AndroidConfig(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"name": "x"}`))
	}))
	defer server.Close()

	a := fcm.NewV1Adapter("p", "")
	a.SetEndpoint(server.URL)
	a.SetTokenProvider(&fcm.StaticTokenProvider{Token: "tok"})

	msg := &fcm.Message{
		Token: "t",
		Notification: &fcm.Notification{Title: "T", Body: "B"},
		Android: &fcm.Android{
			Priority:    "high",
			TTL:         time.Hour,
			ChannelID:   "health",
			Sound:       "alert.mp3",
			CollapseKey: "group-1",
		},
	}
	_, err := a.Send(context.Background(), msg)
	if err != nil {
		t.Errorf("Android 설정 발송 실패: %v", err)
	}
}

// TestV1Adapter_IOSConfig는 iOS APNS 설정 빌드를 검증합니다.
func TestV1Adapter_IOSConfig(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"name": "x"}`))
	}))
	defer server.Close()

	a := fcm.NewV1Adapter("p", "")
	a.SetEndpoint(server.URL)
	a.SetTokenProvider(&fcm.StaticTokenProvider{Token: "tok"})

	msg := &fcm.Message{
		Token: "t",
		Notification: &fcm.Notification{Title: "T", Body: "B"},
		IOS: &fcm.IOS{
			BadgeCount:       5,
			Sound:            "default",
			Category:         "HEALTH_ALERT",
			ContentAvailable: true,
			MutableContent:   true,
		},
	}
	_, err := a.Send(context.Background(), msg)
	if err != nil {
		t.Errorf("iOS 설정 발송 실패: %v", err)
	}
}

// TestV1Adapter_HealthCheck는 헬스체크를 검증합니다.
func TestV1Adapter_HealthCheck(t *testing.T) {
	a := fcm.NewV1Adapter("project", "")
	a.SetTokenProvider(&fcm.StaticTokenProvider{Token: "tok"})

	if err := a.HealthCheck(context.Background()); err != nil {
		t.Errorf("HealthCheck 실패: %v", err)
	}

	a2 := fcm.NewV1Adapter("", "")
	if err := a2.HealthCheck(context.Background()); err == nil {
		t.Error("project_id 없이 헬스체크 통과")
	}
}

// TestNewFromEnv_Default는 기본 Noop 폴백을 검증합니다.
func TestNewFromEnv_Default(t *testing.T) {
	t.Setenv("FCM_ENABLED", "")
	a := fcm.NewFromEnv()
	if a.Provider() != "noop" {
		t.Errorf("Provider = %q, want noop", a.Provider())
	}
}

// TestStaticTokenProvider_AccessToken은 정적 토큰을 검증합니다.
func TestStaticTokenProvider_AccessToken(t *testing.T) {
	p := &fcm.StaticTokenProvider{Token: "abc"}
	tok, err := p.AccessToken(context.Background())
	if err != nil {
		t.Fatalf("AccessToken 실패: %v", err)
	}
	if tok != "abc" {
		t.Errorf("token = %q", tok)
	}

	p2 := &fcm.StaticTokenProvider{}
	if _, err := p2.AccessToken(context.Background()); err == nil {
		t.Error("빈 토큰 통과됨")
	}
}

// TestV1Adapter_SendMulticast는 다중 발송 + 무효 토큰 수집을 검증합니다.
func TestV1Adapter_SendMulticast(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// "bad-token" 메시지에 대해서만 실패 응답
		buf := make([]byte, 2048)
		n, _ := r.Body.Read(buf)
		body := string(buf[:n])
		if contains(body, "bad-token") {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": {"code": 400, "status": "UNREGISTERED"}}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"name": "ok"}`))
	}))
	defer server.Close()

	a := fcm.NewV1Adapter("p", "")
	a.SetEndpoint(server.URL)
	a.SetTokenProvider(&fcm.StaticTokenProvider{Token: "tok"})

	msg := &fcm.Message{
		Notification: &fcm.Notification{Title: "T", Body: "B"},
	}
	tokens := []string{"good-1", "bad-token", "good-2"}

	result, err := a.SendMulticast(context.Background(), msg, tokens)
	if err != nil {
		t.Fatalf("SendMulticast 실패: %v", err)
	}
	if result.SuccessCount != 2 {
		t.Errorf("SuccessCount = %d, want 2", result.SuccessCount)
	}
	if result.FailureCount != 1 {
		t.Errorf("FailureCount = %d, want 1", result.FailureCount)
	}
	if len(result.InvalidTokens) != 1 || result.InvalidTokens[0] != "bad-token" {
		t.Errorf("InvalidTokens = %v, want [bad-token]", result.InvalidTokens)
	}
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
