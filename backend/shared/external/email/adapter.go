// Package email은 Email 발송을 위한 통합 어댑터 인터페이스를 제공합니다.
//
// 지원 프로바이더: SendGrid, AWS SES, Noop
// 환경변수 EMAIL_PROVIDER로 스위칭 (기본: noop)
package email

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

// Address는 이메일 주소입니다.
type Address struct {
	Email string
	Name  string
}

// String은 RFC 2822 형식 문자열을 반환합니다.
func (a *Address) String() string {
	if a.Name != "" {
		return fmt.Sprintf("%s <%s>", a.Name, a.Email)
	}
	return a.Email
}

// Message는 이메일 발송 메시지입니다.
type Message struct {
	From       *Address
	To         []*Address
	Cc         []*Address
	Bcc        []*Address
	ReplyTo    *Address
	Subject    string
	TextBody   string
	HTMLBody   string
	Attachments []*Attachment
	Headers    map[string]string
	TemplateID string            // SendGrid 템플릿 ID (옵션)
	TemplateData map[string]string // 템플릿 변수
}

// Attachment는 이메일 첨부 파일입니다.
type Attachment struct {
	Filename    string
	ContentType string
	Data        []byte
}

// SendResult는 발송 결과입니다.
type SendResult struct {
	MessageID  string
	Provider   string
	SentAt     time.Time
	Status     string // "queued", "sent", "delivered", "failed", "bounced"
	ErrorCode  string
	ErrorMsg   string
}

// ============================================================================
// 어댑터 인터페이스
// ============================================================================

// Adapter는 이메일 프로바이더 어댑터 인터페이스입니다.
type Adapter interface {
	Send(ctx context.Context, msg *Message) (*SendResult, error)
	Provider() string
	HealthCheck(ctx context.Context) error
}

// ============================================================================
// 검증
// ============================================================================

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// ValidateAddress는 이메일 주소를 검증합니다.
func ValidateAddress(addr *Address) error {
	if addr == nil || addr.Email == "" {
		return errors.New("address is empty")
	}
	if !emailRegex.MatchString(addr.Email) {
		return fmt.Errorf("invalid email format: %q", addr.Email)
	}
	return nil
}

// ValidateMessage는 메시지를 검증합니다.
func ValidateMessage(msg *Message) error {
	if msg == nil {
		return errors.New("message is nil")
	}
	if msg.From == nil {
		return errors.New("from address required")
	}
	if err := ValidateAddress(msg.From); err != nil {
		return fmt.Errorf("invalid from: %w", err)
	}
	if len(msg.To) == 0 {
		return errors.New("at least one to address required")
	}
	for i, addr := range msg.To {
		if err := ValidateAddress(addr); err != nil {
			return fmt.Errorf("to[%d]: %w", i, err)
		}
	}
	if strings.TrimSpace(msg.Subject) == "" {
		return errors.New("subject required")
	}
	if msg.TextBody == "" && msg.HTMLBody == "" && msg.TemplateID == "" {
		return errors.New("at least one of TextBody, HTMLBody, or TemplateID required")
	}
	return nil
}

// ============================================================================
// 팩토리
// ============================================================================

// NewFromEnv는 환경변수에 따라 어댑터를 생성합니다.
//
// EMAIL_PROVIDER:
//   - "sendgrid": SendGrid (SENDGRID_API_KEY)
//   - "ses": AWS SES (AWS_REGION, SES_FROM_EMAIL)
//   - "" or "noop": Noop (개발용)
func NewFromEnv() Adapter {
	provider := strings.ToLower(os.Getenv("EMAIL_PROVIDER"))
	switch provider {
	case "sendgrid":
		return NewSendGridAdapter(
			os.Getenv("SENDGRID_API_KEY"),
			os.Getenv("SENDGRID_FROM_EMAIL"),
			os.Getenv("SENDGRID_FROM_NAME"),
		)
	case "ses":
		return NewSESAdapter(
			os.Getenv("AWS_REGION"),
			os.Getenv("SES_FROM_EMAIL"),
		)
	default:
		return NewNoopAdapter()
	}
}

// ============================================================================
// HTML 템플릿
// ============================================================================

// BuiltInTemplates는 내장 HTML 템플릿입니다.
var BuiltInTemplates = map[string]string{
	"welcome": `<!DOCTYPE html>
<html><body style="font-family: sans-serif; padding: 20px;">
  <h1>만파식에 오신 것을 환영합니다</h1>
  <p>안녕하세요 {{name}}님,</p>
  <p>만파식 가입을 진심으로 환영합니다. 건강한 일상을 함께 만들어가요.</p>
</body></html>`,

	"alert": `<!DOCTYPE html>
<html><body style="font-family: sans-serif; padding: 20px; background: #fff3cd;">
  <h2 style="color: #856404;">건강 알림</h2>
  <p>{{message}}</p>
  <p><strong>측정 시각:</strong> {{measured_at}}</p>
</body></html>`,

	"receipt": `<!DOCTYPE html>
<html><body style="font-family: sans-serif; padding: 20px;">
  <h1>주문 영수증</h1>
  <p>주문번호: {{order_id}}</p>
  <p>총 결제금액: {{total_amount}}</p>
</body></html>`,

	"password_reset": `<!DOCTYPE html>
<html><body style="font-family: sans-serif; padding: 20px;">
  <h1>비밀번호 재설정</h1>
  <p>아래 링크를 클릭하여 비밀번호를 재설정하세요 (24시간 내 만료):</p>
  <a href="{{reset_link}}">비밀번호 재설정</a>
  <p>요청하지 않으셨다면 이 메일을 무시하세요.</p>
</body></html>`,
}

// RenderTemplate는 템플릿에 변수를 치환합니다.
func RenderTemplate(name string, vars map[string]string) (string, error) {
	tmpl, ok := BuiltInTemplates[name]
	if !ok {
		return "", fmt.Errorf("template %q not found", name)
	}
	rendered := tmpl
	for k, v := range vars {
		rendered = strings.ReplaceAll(rendered, "{{"+k+"}}", v)
	}
	return rendered, nil
}

// ============================================================================
// Noop 어댑터
// ============================================================================

// NoopAdapter는 발송하지 않고 메모리에 저장만 하는 어댑터입니다.
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
	}, nil
}

// Provider는 프로바이더 이름을 반환합니다.
func (a *NoopAdapter) Provider() string { return "noop" }

// HealthCheck는 항상 성공합니다.
func (a *NoopAdapter) HealthCheck(_ context.Context) error { return nil }

// SentMessages는 저장된 메시지 목록을 반환합니다.
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
