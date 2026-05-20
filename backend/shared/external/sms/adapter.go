// Package sms는 SMS 발송을 위한 통합 어댑터 인터페이스를 제공합니다.
//
// 지원 프로바이더: Twilio, Toss SMS
// 환경변수 SMS_PROVIDER로 스위칭 (기본: noop)
package sms

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

// ============================================================================
// 도메인 모델
// ============================================================================

// Message는 SMS 발송 메시지입니다.
type Message struct {
	From      string // 발신번호 (E.164 형식)
	To        string // 수신번호 (E.164 형식)
	Body      string // 본문 (90자 이하 권장)
	IsLMS     bool   // LMS (장문) 여부 (한국 한정)
	Subject   string // LMS 제목 (한국 한정)
	RequestID string // 추적용 ID
}

// SendResult는 발송 결과입니다.
type SendResult struct {
	MessageID  string
	Provider   string
	SentAt     time.Time
	Cost       float64 // 원/달러
	Currency   string
	Status     string // "queued", "sent", "delivered", "failed"
	ErrorCode  string
	ErrorMsg   string
}

// ============================================================================
// 어댑터 인터페이스
// ============================================================================

// Adapter는 SMS 프로바이더 어댑터 인터페이스입니다.
type Adapter interface {
	Send(ctx context.Context, msg *Message) (*SendResult, error)
	Provider() string
	HealthCheck(ctx context.Context) error
}

// ============================================================================
// 검증
// ============================================================================

// e164Regex는 E.164 형식 전화번호 검증입니다.
var e164Regex = regexp.MustCompile(`^\+[1-9]\d{1,14}$`)

// ValidateMessage는 메시지를 검증합니다.
func ValidateMessage(msg *Message) error {
	if msg == nil {
		return errors.New("message is nil")
	}
	if !e164Regex.MatchString(msg.To) {
		return fmt.Errorf("invalid To: %q (must be E.164)", msg.To)
	}
	if msg.From != "" && !e164Regex.MatchString(msg.From) {
		return fmt.Errorf("invalid From: %q (must be E.164)", msg.From)
	}
	body := strings.TrimSpace(msg.Body)
	if body == "" {
		return errors.New("body is empty")
	}
	if msg.IsLMS && len(body) > 2000 {
		return errors.New("LMS body exceeds 2000 characters")
	}
	if !msg.IsLMS && len([]rune(body)) > 90 {
		return errors.New("SMS body exceeds 90 characters (use LMS for longer)")
	}
	return nil
}

// ============================================================================
// 팩토리 (환경변수 스위칭)
// ============================================================================

// NewFromEnv는 환경변수 SMS_PROVIDER에 따라 어댑터를 생성합니다.
//
// SMS_PROVIDER:
//   - "twilio": Twilio 어댑터 (TWILIO_ACCOUNT_SID, TWILIO_AUTH_TOKEN)
//   - "toss": Toss SMS 어댑터 (TOSS_SMS_API_KEY, TOSS_SMS_FROM)
//   - "" or "noop": 발송하지 않고 로깅만 (개발용)
func NewFromEnv() Adapter {
	provider := strings.ToLower(os.Getenv("SMS_PROVIDER"))
	switch provider {
	case "twilio":
		return NewTwilioAdapter(
			os.Getenv("TWILIO_ACCOUNT_SID"),
			os.Getenv("TWILIO_AUTH_TOKEN"),
			os.Getenv("TWILIO_FROM"),
		)
	case "toss":
		return NewTossAdapter(
			os.Getenv("TOSS_SMS_API_KEY"),
			os.Getenv("TOSS_SMS_FROM"),
		)
	default:
		return NewNoopAdapter()
	}
}

// ============================================================================
// Noop 어댑터 (개발용)
// ============================================================================

// NoopAdapter는 발송하지 않고 로깅만 하는 어댑터입니다.
type NoopAdapter struct {
	mu       sync.Mutex
	messages []*Message
}

// NewNoopAdapter는 새 Noop 어댑터를 생성합니다.
func NewNoopAdapter() *NoopAdapter {
	return &NoopAdapter{}
}

// Send는 메시지를 메모리에 저장하고 성공 응답을 반환합니다.
func (a *NoopAdapter) Send(ctx context.Context, msg *Message) (*SendResult, error) {
	if err := ValidateMessage(msg); err != nil {
		return nil, err
	}
	a.mu.Lock()
	a.messages = append(a.messages, msg)
	a.mu.Unlock()
	return &SendResult{
		MessageID: fmt.Sprintf("noop-%d", time.Now().UnixNano()),
		Provider:  "noop",
		SentAt:    time.Now().UTC(),
		Status:    "sent",
		Currency:  "KRW",
	}, nil
}

// Provider는 프로바이더 이름을 반환합니다.
func (a *NoopAdapter) Provider() string { return "noop" }

// HealthCheck는 항상 성공합니다.
func (a *NoopAdapter) HealthCheck(_ context.Context) error { return nil }

// SentMessages는 저장된 메시지 목록을 반환합니다 (테스트 용도).
func (a *NoopAdapter) SentMessages() []*Message {
	a.mu.Lock()
	defer a.mu.Unlock()
	return append([]*Message{}, a.messages...)
}

// Clear는 저장된 메시지를 비웁니다.
func (a *NoopAdapter) Clear() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.messages = nil
}

// Count는 저장된 메시지 수를 반환합니다.
func (a *NoopAdapter) Count() int {
	a.mu.Lock()
	defer a.mu.Unlock()
	return len(a.messages)
}
