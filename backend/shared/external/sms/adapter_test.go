package sms_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/manpasik/backend/shared/external/sms"
)

// TestValidateMessage_Valid는 정상 메시지 검증을 확인합니다.
func TestValidateMessage_Valid(t *testing.T) {
	msg := &sms.Message{
		From: "+821012345678",
		To:   "+821087654321",
		Body: "안녕하세요",
	}
	if err := sms.ValidateMessage(msg); err != nil {
		t.Errorf("ValidateMessage 실패: %v", err)
	}
}

// TestValidateMessage_InvalidPhone은 잘못된 전화번호 거부를 검증합니다.
func TestValidateMessage_InvalidPhone(t *testing.T) {
	cases := []*sms.Message{
		{To: "010-1234-5678", Body: "x"}, // 하이픈
		{To: "01012345678", Body: "x"},   // + 없음
		{To: "+", Body: "x"},              // 잘못된 형식
	}
	for i, c := range cases {
		if err := sms.ValidateMessage(c); err == nil {
			t.Errorf("case %d: 잘못된 전화번호 통과: %s", i, c.To)
		}
	}
}

// TestValidateMessage_BodyLength는 본문 길이 검증을 확인합니다.
func TestValidateMessage_BodyLength(t *testing.T) {
	// SMS 90자 초과 거부
	long := strings.Repeat("가", 91)
	msg := &sms.Message{To: "+821012345678", Body: long, IsLMS: false}
	if err := sms.ValidateMessage(msg); err == nil {
		t.Error("SMS 90자 초과가 허용됨")
	}

	// LMS 2000자 초과 거부
	veryLong := strings.Repeat("가", 2001)
	msg2 := &sms.Message{To: "+821012345678", Body: veryLong, IsLMS: true}
	if err := sms.ValidateMessage(msg2); err == nil {
		t.Error("LMS 2000자 초과가 허용됨")
	}
}

// TestValidateMessage_EmptyBody는 빈 본문 거부를 검증합니다.
func TestValidateMessage_EmptyBody(t *testing.T) {
	msg := &sms.Message{To: "+821012345678", Body: ""}
	if err := sms.ValidateMessage(msg); err == nil {
		t.Error("빈 본문이 허용됨")
	}
}

// TestNoopAdapter_Send는 Noop 어댑터 발송을 검증합니다.
func TestNoopAdapter_Send(t *testing.T) {
	a := sms.NewNoopAdapter()
	msg := &sms.Message{To: "+821012345678", Body: "테스트"}

	result, err := a.Send(context.Background(), msg)
	if err != nil {
		t.Fatalf("Send 실패: %v", err)
	}
	if result.Provider != "noop" {
		t.Errorf("Provider = %q, want noop", result.Provider)
	}
	if result.Status != "sent" {
		t.Errorf("Status = %q, want sent", result.Status)
	}
	if a.Count() != 1 {
		t.Errorf("Count = %d, want 1", a.Count())
	}
}

// TestNoopAdapter_ValidationFailure는 검증 실패 시 발송 거부를 검증합니다.
func TestNoopAdapter_ValidationFailure(t *testing.T) {
	a := sms.NewNoopAdapter()
	msg := &sms.Message{To: "invalid", Body: "x"}

	_, err := a.Send(context.Background(), msg)
	if err == nil {
		t.Error("잘못된 메시지가 통과됨")
	}
	if a.Count() != 0 {
		t.Errorf("Count = %d, want 0 (검증 실패는 저장 안됨)", a.Count())
	}
}

// TestNoopAdapter_Clear는 메시지 클리어를 검증합니다.
func TestNoopAdapter_Clear(t *testing.T) {
	a := sms.NewNoopAdapter()
	_, _ = a.Send(context.Background(), &sms.Message{To: "+821012345678", Body: "x"})
	_, _ = a.Send(context.Background(), &sms.Message{To: "+821012345678", Body: "y"})

	a.Clear()
	if a.Count() != 0 {
		t.Errorf("Clear 후 Count = %d, want 0", a.Count())
	}
}

// TestTwilioAdapter_Send_Success는 Twilio 정상 발송을 검증합니다.
func TestTwilioAdapter_Send_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Basic auth 검증
		user, pass, ok := r.BasicAuth()
		if !ok || user != "AC123" || pass != "secret" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"sid": "SM123"}`))
	}))
	defer server.Close()

	a := sms.NewTwilioAdapter("AC123", "secret", "+15551234567")
	a.SetEndpoint(server.URL)

	msg := &sms.Message{To: "+821012345678", Body: "Hi"}
	result, err := a.Send(context.Background(), msg)
	if err != nil {
		t.Fatalf("Send 실패: %v", err)
	}
	if result.Provider != "twilio" {
		t.Errorf("Provider = %q, want twilio", result.Provider)
	}
	if result.Status != "sent" {
		t.Errorf("Status = %q, want sent", result.Status)
	}
}

// TestTwilioAdapter_Send_Unauthorized는 Twilio 인증 실패 처리를 검증합니다.
func TestTwilioAdapter_Send_Unauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	}))
	defer server.Close()

	a := sms.NewTwilioAdapter("AC123", "wrong", "+15551234567")
	a.SetEndpoint(server.URL)

	_, err := a.Send(context.Background(), &sms.Message{To: "+821012345678", Body: "x"})
	if err == nil {
		t.Error("401 응답에 에러가 발생하지 않음")
	}
}

// TestTwilioAdapter_HealthCheck는 헬스체크를 검증합니다.
func TestTwilioAdapter_HealthCheck(t *testing.T) {
	a := sms.NewTwilioAdapter("AC123", "secret", "+15551234567")
	if err := a.HealthCheck(context.Background()); err != nil {
		t.Errorf("HealthCheck 실패: %v", err)
	}

	// 자격증명 미설정
	a2 := sms.NewTwilioAdapter("", "", "")
	if err := a2.HealthCheck(context.Background()); err == nil {
		t.Error("자격증명 미설정에 헬스체크 통과")
	}
}

// TestTossAdapter_Send_Success는 Toss 정상 발송을 검증합니다.
func TestTossAdapter_Send_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-key" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message_id": "ts-001", "status": "sent", "cost": 22.5}`))
	}))
	defer server.Close()

	a := sms.NewTossAdapter("test-key", "+821012345678")
	a.SetEndpoint(server.URL)

	msg := &sms.Message{To: "+821087654321", Body: "안녕"}
	result, err := a.Send(context.Background(), msg)
	if err != nil {
		t.Fatalf("Send 실패: %v", err)
	}
	if result.Provider != "toss" {
		t.Errorf("Provider = %q", result.Provider)
	}
	if result.MessageID != "ts-001" {
		t.Errorf("MessageID = %q, want ts-001", result.MessageID)
	}
	if result.Cost != 22.5 {
		t.Errorf("Cost = %f, want 22.5", result.Cost)
	}
}

// TestTossAdapter_LMS는 LMS(장문) 발송을 검증합니다.
func TestTossAdapter_LMS(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message_id": "lms-001", "status": "sent", "cost": 33.0}`))
	}))
	defer server.Close()

	a := sms.NewTossAdapter("k", "+821012345678")
	a.SetEndpoint(server.URL)

	long := strings.Repeat("긴 메시지 ", 50)
	msg := &sms.Message{To: "+821087654321", Body: long, IsLMS: true, Subject: "공지"}
	result, err := a.Send(context.Background(), msg)
	if err != nil {
		t.Fatalf("Send 실패: %v", err)
	}
	if result.Status != "sent" {
		t.Errorf("Status = %q", result.Status)
	}
}

// TestTossAdapter_HealthCheck는 헬스체크를 검증합니다.
func TestTossAdapter_HealthCheck(t *testing.T) {
	a := sms.NewTossAdapter("k", "")
	if err := a.HealthCheck(context.Background()); err != nil {
		t.Errorf("HealthCheck 실패: %v", err)
	}

	a2 := sms.NewTossAdapter("", "")
	if err := a2.HealthCheck(context.Background()); err == nil {
		t.Error("apiKey 없이 헬스체크 통과")
	}
}

// TestNewFromEnv_Default는 환경변수 미설정 시 Noop 폴백을 검증합니다.
func TestNewFromEnv_Default(t *testing.T) {
	t.Setenv("SMS_PROVIDER", "")
	a := sms.NewFromEnv()
	if a.Provider() != "noop" {
		t.Errorf("Provider = %q, want noop", a.Provider())
	}
}

// TestNewFromEnv_Twilio는 SMS_PROVIDER=twilio 스위칭을 검증합니다.
func TestNewFromEnv_Twilio(t *testing.T) {
	t.Setenv("SMS_PROVIDER", "twilio")
	t.Setenv("TWILIO_ACCOUNT_SID", "AC123")
	t.Setenv("TWILIO_AUTH_TOKEN", "secret")
	t.Setenv("TWILIO_FROM", "+15551234567")

	a := sms.NewFromEnv()
	if a.Provider() != "twilio" {
		t.Errorf("Provider = %q, want twilio", a.Provider())
	}
}

// TestNewFromEnv_Toss는 SMS_PROVIDER=toss 스위칭을 검증합니다.
func TestNewFromEnv_Toss(t *testing.T) {
	t.Setenv("SMS_PROVIDER", "toss")
	t.Setenv("TOSS_SMS_API_KEY", "key")

	a := sms.NewFromEnv()
	if a.Provider() != "toss" {
		t.Errorf("Provider = %q, want toss", a.Provider())
	}
}

// TestE164Validation는 다양한 E.164 형식을 검증합니다.
func TestE164Validation(t *testing.T) {
	valid := []string{"+821012345678", "+15551234567", "+447911123456"}
	invalid := []string{"+0", "+0123456", "010-1234-5678", "1234567"}

	for _, p := range valid {
		msg := &sms.Message{To: p, Body: "x"}
		if err := sms.ValidateMessage(msg); err != nil {
			t.Errorf("정상 번호 %q 거부됨: %v", p, err)
		}
	}
	for _, p := range invalid {
		msg := &sms.Message{To: p, Body: "x"}
		if err := sms.ValidateMessage(msg); err == nil {
			t.Errorf("잘못된 번호 %q 통과됨", p)
		}
	}
}

// dummyDoer는 항상 에러를 반환하는 HTTP 클라이언트입니다 (테스트용).
type dummyDoer struct {
	err error
}

func (d *dummyDoer) Do(_ *http.Request) (*http.Response, error) {
	return nil, d.err
}

// TestTwilioAdapter_NetworkError는 네트워크 에러 처리를 검증합니다.
func TestTwilioAdapter_NetworkError(t *testing.T) {
	a := sms.NewTwilioAdapter("AC", "secret", "+15551234567")
	a.SetHTTPClient(&dummyDoer{err: errors.New("network down")})

	_, err := a.Send(context.Background(), &sms.Message{To: "+821012345678", Body: "x"})
	if err == nil {
		t.Error("네트워크 에러가 전파되지 않음")
	}
}
