package email_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/manpasik/backend/shared/external/email"
)

// TestValidateAddress_Valid는 정상 이메일 검증을 확인합니다.
func TestValidateAddress_Valid(t *testing.T) {
	addrs := []string{
		"user@example.com",
		"first.last+tag@sub.domain.co.kr",
		"name123@test-domain.org",
	}
	for _, e := range addrs {
		if err := email.ValidateAddress(&email.Address{Email: e}); err != nil {
			t.Errorf("정상 주소 %q 거부됨: %v", e, err)
		}
	}
}

// TestValidateAddress_Invalid는 잘못된 이메일 거부를 검증합니다.
func TestValidateAddress_Invalid(t *testing.T) {
	addrs := []string{
		"notanemail",
		"@nodomain.com",
		"missing@dot",
		"spaces in@email.com",
	}
	for _, e := range addrs {
		if err := email.ValidateAddress(&email.Address{Email: e}); err == nil {
			t.Errorf("잘못된 주소 %q 통과됨", e)
		}
	}
}

// TestValidateMessage_Valid는 정상 메시지 검증을 확인합니다.
func TestValidateMessage_Valid(t *testing.T) {
	msg := &email.Message{
		From:     &email.Address{Email: "from@example.com"},
		To:       []*email.Address{{Email: "to@example.com"}},
		Subject:  "Test",
		TextBody: "Hello",
	}
	if err := email.ValidateMessage(msg); err != nil {
		t.Errorf("ValidateMessage 실패: %v", err)
	}
}

// TestValidateMessage_NoFrom는 발신자 누락 거부를 검증합니다.
func TestValidateMessage_NoFrom(t *testing.T) {
	msg := &email.Message{
		To:       []*email.Address{{Email: "to@example.com"}},
		Subject:  "Test",
		TextBody: "Hello",
	}
	if err := email.ValidateMessage(msg); err == nil {
		t.Error("From 누락이 허용됨")
	}
}

// TestValidateMessage_NoTo는 수신자 누락 거부를 검증합니다.
func TestValidateMessage_NoTo(t *testing.T) {
	msg := &email.Message{
		From:     &email.Address{Email: "from@example.com"},
		Subject:  "Test",
		TextBody: "Hello",
	}
	if err := email.ValidateMessage(msg); err == nil {
		t.Error("To 누락이 허용됨")
	}
}

// TestValidateMessage_NoBody는 본문 누락 거부를 검증합니다.
func TestValidateMessage_NoBody(t *testing.T) {
	msg := &email.Message{
		From:    &email.Address{Email: "f@x.com"},
		To:      []*email.Address{{Email: "t@x.com"}},
		Subject: "Test",
	}
	if err := email.ValidateMessage(msg); err == nil {
		t.Error("본문 누락이 허용됨")
	}
}

// TestAddress_String은 RFC 2822 형식 변환을 검증합니다.
func TestAddress_String(t *testing.T) {
	a1 := &email.Address{Email: "user@example.com"}
	if a1.String() != "user@example.com" {
		t.Errorf("addr1 = %q", a1.String())
	}

	a2 := &email.Address{Email: "user@example.com", Name: "홍길동"}
	expected := "홍길동 <user@example.com>"
	if a2.String() != expected {
		t.Errorf("addr2 = %q, want %q", a2.String(), expected)
	}
}

// TestRenderTemplate_Welcome은 환영 템플릿 렌더링을 검증합니다.
func TestRenderTemplate_Welcome(t *testing.T) {
	html, err := email.RenderTemplate("welcome", map[string]string{
		"name": "홍길동",
	})
	if err != nil {
		t.Fatalf("Render 실패: %v", err)
	}
	if !strings.Contains(html, "홍길동") {
		t.Errorf("이름이 치환되지 않음: %s", html)
	}
	if strings.Contains(html, "{{name}}") {
		t.Error("템플릿 변수가 치환되지 않음")
	}
}

// TestRenderTemplate_Alert는 알림 템플릿 렌더링을 검증합니다.
func TestRenderTemplate_Alert(t *testing.T) {
	html, err := email.RenderTemplate("alert", map[string]string{
		"message":     "혈당 수치가 높습니다",
		"measured_at": "2026-04-30 14:00",
	})
	if err != nil {
		t.Fatalf("Render 실패: %v", err)
	}
	if !strings.Contains(html, "혈당 수치가 높습니다") {
		t.Error("메시지가 치환되지 않음")
	}
}

// TestRenderTemplate_NotFound는 미지의 템플릿 거부를 검증합니다.
func TestRenderTemplate_NotFound(t *testing.T) {
	_, err := email.RenderTemplate("xyz", map[string]string{})
	if err == nil {
		t.Error("미지 템플릿이 통과됨")
	}
}

// TestNoopAdapter_Send는 Noop 발송을 검증합니다.
func TestNoopAdapter_Send(t *testing.T) {
	a := email.NewNoopAdapter()
	msg := &email.Message{
		From:     &email.Address{Email: "f@x.com"},
		To:       []*email.Address{{Email: "t@x.com"}},
		Subject:  "Test",
		TextBody: "Hello",
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

// TestSendGridAdapter_Send_Success는 SendGrid 정상 발송을 검증합니다.
func TestSendGridAdapter_Send_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-key" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		w.Header().Set("X-Message-Id", "sg-test-001")
		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	a := email.NewSendGridAdapter("test-key", "from@x.com", "From")
	a.SetEndpoint(server.URL)

	msg := &email.Message{
		From:     &email.Address{Email: "from@x.com"},
		To:       []*email.Address{{Email: "to@x.com"}},
		Subject:  "Test",
		HTMLBody: "<p>Hello</p>",
	}

	result, err := a.Send(context.Background(), msg)
	if err != nil {
		t.Fatalf("Send 실패: %v", err)
	}
	if result.MessageID != "sg-test-001" {
		t.Errorf("MessageID = %q, want sg-test-001", result.MessageID)
	}
}

// TestSendGridAdapter_Send_NoAPIKey는 API 키 누락 시 거부를 검증합니다.
func TestSendGridAdapter_Send_NoAPIKey(t *testing.T) {
	a := email.NewSendGridAdapter("", "from@x.com", "")
	msg := &email.Message{
		From:     &email.Address{Email: "f@x.com"},
		To:       []*email.Address{{Email: "t@x.com"}},
		Subject:  "Test",
		TextBody: "Hello",
	}
	_, err := a.Send(context.Background(), msg)
	if err == nil {
		t.Error("API 키 없이 발송 허용됨")
	}
}

// TestSendGridAdapter_HealthCheck는 헬스체크를 검증합니다.
func TestSendGridAdapter_HealthCheck(t *testing.T) {
	a := email.NewSendGridAdapter("key", "f@x.com", "")
	if err := a.HealthCheck(context.Background()); err != nil {
		t.Errorf("HealthCheck 실패: %v", err)
	}

	a2 := email.NewSendGridAdapter("", "", "")
	if err := a2.HealthCheck(context.Background()); err == nil {
		t.Error("API 키 없이 헬스체크 통과")
	}
}

// TestSESAdapter_Send_Success는 SES 발송을 검증합니다.
func TestSESAdapter_Send_Success(t *testing.T) {
	a := email.NewSESAdapter("us-east-1", "from@x.com")
	msg := &email.Message{
		From:     &email.Address{Email: "from@x.com"},
		To:       []*email.Address{{Email: "to@x.com"}},
		Subject:  "Test",
		TextBody: "Hello",
	}
	result, err := a.Send(context.Background(), msg)
	if err != nil {
		t.Fatalf("Send 실패: %v", err)
	}
	if result.Provider != "ses" {
		t.Errorf("Provider = %q", result.Provider)
	}
}

// TestSESAdapter_NoRegion는 region 누락 거부를 검증합니다.
func TestSESAdapter_NoRegion(t *testing.T) {
	a := email.NewSESAdapter("", "from@x.com")
	if err := a.HealthCheck(context.Background()); err == nil {
		t.Error("region 없이 헬스체크 통과")
	}
}

// TestNewFromEnv_Default는 기본 Noop 폴백을 검증합니다.
func TestNewFromEnv_Default(t *testing.T) {
	t.Setenv("EMAIL_PROVIDER", "")
	a := email.NewFromEnv()
	if a.Provider() != "noop" {
		t.Errorf("Provider = %q, want noop", a.Provider())
	}
}

// TestNewFromEnv_SendGrid는 sendgrid 스위칭을 검증합니다.
func TestNewFromEnv_SendGrid(t *testing.T) {
	t.Setenv("EMAIL_PROVIDER", "sendgrid")
	t.Setenv("SENDGRID_API_KEY", "key")
	a := email.NewFromEnv()
	if a.Provider() != "sendgrid" {
		t.Errorf("Provider = %q, want sendgrid", a.Provider())
	}
}

// TestNewFromEnv_SES는 ses 스위칭을 검증합니다.
func TestNewFromEnv_SES(t *testing.T) {
	t.Setenv("EMAIL_PROVIDER", "ses")
	t.Setenv("AWS_REGION", "us-east-1")
	a := email.NewFromEnv()
	if a.Provider() != "ses" {
		t.Errorf("Provider = %q, want ses", a.Provider())
	}
}

// TestSendGridAdapter_TemplateSupport는 템플릿 발송을 검증합니다.
func TestSendGridAdapter_TemplateSupport(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	a := email.NewSendGridAdapter("k", "from@x.com", "")
	a.SetEndpoint(server.URL)

	msg := &email.Message{
		From:         &email.Address{Email: "f@x.com"},
		To:           []*email.Address{{Email: "t@x.com"}},
		Subject:      "Welcome",
		TemplateID:   "d-12345",
		TemplateData: map[string]string{"name": "홍길동"},
	}

	_, err := a.Send(context.Background(), msg)
	if err != nil {
		t.Fatalf("템플릿 발송 실패: %v", err)
	}
}

// TestSendGridAdapter_CcBcc는 Cc/Bcc 발송을 검증합니다.
func TestSendGridAdapter_CcBcc(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	a := email.NewSendGridAdapter("k", "f@x.com", "")
	a.SetEndpoint(server.URL)

	msg := &email.Message{
		From:     &email.Address{Email: "f@x.com"},
		To:       []*email.Address{{Email: "t@x.com"}},
		Cc:       []*email.Address{{Email: "cc@x.com"}},
		Bcc:      []*email.Address{{Email: "bcc@x.com"}},
		Subject:  "Test",
		TextBody: "x",
	}

	_, err := a.Send(context.Background(), msg)
	if err != nil {
		t.Fatalf("Cc/Bcc 발송 실패: %v", err)
	}
}
