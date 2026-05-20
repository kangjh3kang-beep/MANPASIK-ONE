// Package kakao는 카카오 알림톡 어댑터입니다.
//
// 알림톡은 사전 승인된 템플릿을 통해 카카오톡으로 메시지를 발송하는 한국 한정 채널입니다.
// 본 구현은 NHN Cloud, KakaoBiz 등 카카오 알림톡 중계사 API를 추상화합니다.
package kakao

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

// ============================================================================
// 도메인 모델
// ============================================================================

// Message는 알림톡 메시지입니다.
type Message struct {
	TemplateCode string            // 사전 등록된 템플릿 코드
	To           string            // 수신자 전화번호 (E.164 또는 010-XXXX-XXXX)
	Variables    map[string]string // 템플릿 변수 (#{변수명})
	FailbackSMS  *FailbackSMS      // 실패 시 SMS 자동 발송 옵션
	RequestID    string            // 추적용
}

// FailbackSMS는 알림톡 실패 시 폴백 SMS 옵션입니다.
type FailbackSMS struct {
	Body  string
	IsLMS bool
	From  string
}

// Button은 알림톡 버튼입니다.
type Button struct {
	Type   string // "DS" (배송), "WL" (URL), "AL" (앱), "BK" (봇 키워드)
	Name   string
	URLPC  string
	URLMobile string
	SchemeAndroid string
	SchemeIOS     string
}

// SendResult는 발송 결과입니다.
type SendResult struct {
	MessageID string
	Provider  string
	SentAt    time.Time
	Status    string // "sent", "failed", "fallback_sms"
	ErrorCode string
	ErrorMsg  string
}

// ============================================================================
// 어댑터 인터페이스
// ============================================================================

// Adapter는 카카오 알림톡 어댑터 인터페이스입니다.
type Adapter interface {
	Send(ctx context.Context, msg *Message) (*SendResult, error)
	Provider() string
	HealthCheck(ctx context.Context) error
	RegisterTemplate(code string, body string, buttons []Button) error
	GetTemplate(code string) (*Template, error)
}

// Template는 사전 등록된 알림톡 템플릿입니다.
type Template struct {
	Code       string
	Body       string
	Buttons    []Button
	Variables  []string  // body에서 추출된 변수 목록
	ApprovedAt time.Time
	Status     string    // "pending", "approved", "rejected"
}

// ============================================================================
// 검증
// ============================================================================

// phoneRegex는 한국/E.164 전화번호 검증 정규식입니다.
var phoneRegex = regexp.MustCompile(`^(010-\d{4}-\d{4}|01\d{8,9}|\+82\d{9,10})$`)

// templateVarRegex는 템플릿 변수 추출 정규식 (#{변수명})입니다.
var templateVarRegex = regexp.MustCompile(`#\{([^}]+)\}`)

// ValidateMessage는 메시지를 검증합니다.
func ValidateMessage(msg *Message) error {
	if msg == nil {
		return errors.New("message is nil")
	}
	if msg.TemplateCode == "" {
		return errors.New("template_code required")
	}
	if msg.To == "" {
		return errors.New("recipient phone required")
	}
	if !phoneRegex.MatchString(msg.To) {
		return fmt.Errorf("invalid phone format: %q", msg.To)
	}
	return nil
}

// ExtractTemplateVariables는 템플릿 본문에서 #{변수} 형태의 변수를 추출합니다.
func ExtractTemplateVariables(body string) []string {
	matches := templateVarRegex.FindAllStringSubmatch(body, -1)
	seen := make(map[string]bool)
	vars := []string{}
	for _, m := range matches {
		if len(m) > 1 && !seen[m[1]] {
			seen[m[1]] = true
			vars = append(vars, m[1])
		}
	}
	return vars
}

// RenderTemplate는 템플릿 본문에 변수를 치환합니다.
func RenderTemplate(body string, vars map[string]string) string {
	return templateVarRegex.ReplaceAllStringFunc(body, func(match string) string {
		// "#{name}" → "name"
		varName := match[2 : len(match)-1]
		if v, ok := vars[varName]; ok {
			return v
		}
		return match // 미치환 변수는 그대로 유지
	})
}

// ============================================================================
// 팩토리
// ============================================================================

// NewFromEnv는 환경변수에 따라 어댑터를 생성합니다.
//
// KAKAO_PROVIDER:
//   - "nhncloud": NHN Cloud Kakao Alimtalk
//   - "kakaobiz": Kakao Biz API
//   - "" or "noop": Noop
func NewFromEnv() Adapter {
	provider := strings.ToLower(os.Getenv("KAKAO_PROVIDER"))
	switch provider {
	case "nhncloud":
		return NewNHNCloudAdapter(
			os.Getenv("NHNCLOUD_APP_KEY"),
			os.Getenv("NHNCLOUD_SECRET_KEY"),
			os.Getenv("KAKAO_SENDER_KEY"),
		)
	case "kakaobiz":
		return NewKakaoBizAdapter(
			os.Getenv("KAKAO_BIZ_API_KEY"),
			os.Getenv("KAKAO_SENDER_KEY"),
		)
	default:
		return NewNoopAdapter()
	}
}

// ============================================================================
// Noop 어댑터
// ============================================================================

// NoopAdapter는 발송하지 않고 메모리에 저장만 하는 어댑터입니다.
type NoopAdapter struct {
	mu        sync.Mutex
	messages  []*Message
	templates map[string]*Template
}

// NewNoopAdapter는 새 Noop 어댑터를 생성합니다.
func NewNoopAdapter() *NoopAdapter {
	return &NoopAdapter{
		templates: make(map[string]*Template),
	}
}

// Send는 메시지를 메모리에 저장합니다.
func (a *NoopAdapter) Send(_ context.Context, msg *Message) (*SendResult, error) {
	if err := ValidateMessage(msg); err != nil {
		return nil, err
	}
	a.mu.Lock()
	a.messages = append(a.messages, msg)
	a.mu.Unlock()
	return &SendResult{
		MessageID: fmt.Sprintf("noop-kakao-%d", time.Now().UnixNano()),
		Provider:  "noop",
		Status:    "sent",
		SentAt:    time.Now().UTC(),
	}, nil
}

// Provider는 프로바이더 이름을 반환합니다.
func (a *NoopAdapter) Provider() string { return "noop" }

// HealthCheck는 항상 성공합니다.
func (a *NoopAdapter) HealthCheck(_ context.Context) error { return nil }

// RegisterTemplate는 템플릿을 등록합니다.
func (a *NoopAdapter) RegisterTemplate(code, body string, buttons []Button) error {
	if code == "" || body == "" {
		return errors.New("code and body required")
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	a.templates[code] = &Template{
		Code:       code,
		Body:       body,
		Buttons:    buttons,
		Variables:  ExtractTemplateVariables(body),
		ApprovedAt: time.Now().UTC(),
		Status:     "approved",
	}
	return nil
}

// GetTemplate는 템플릿을 조회합니다.
func (a *NoopAdapter) GetTemplate(code string) (*Template, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	tmpl, ok := a.templates[code]
	if !ok {
		return nil, fmt.Errorf("template %q not found", code)
	}
	return tmpl, nil
}

// SentMessages는 저장된 메시지 목록을 반환합니다.
func (a *NoopAdapter) SentMessages() []*Message {
	a.mu.Lock()
	defer a.mu.Unlock()
	return append([]*Message{}, a.messages...)
}

// Count는 저장된 메시지 수를 반환합니다.
func (a *NoopAdapter) Count() int {
	a.mu.Lock()
	defer a.mu.Unlock()
	return len(a.messages)
}

// Clear는 메시지를 비웁니다.
func (a *NoopAdapter) Clear() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.messages = nil
}

// ============================================================================
// NHN Cloud 어댑터
// ============================================================================

// HTTPDoer는 테스트 가능한 HTTP 클라이언트 인터페이스입니다.
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// NHNCloudAdapter는 NHN Cloud Alimtalk 어댑터입니다.
type NHNCloudAdapter struct {
	appKey     string
	secretKey  string
	senderKey  string
	httpClient HTTPDoer
	endpoint   string
	templates  map[string]*Template
	tmplMu     sync.RWMutex
}

// NewNHNCloudAdapter는 새 NHN Cloud 어댑터를 생성합니다.
func NewNHNCloudAdapter(appKey, secretKey, senderKey string) *NHNCloudAdapter {
	return &NHNCloudAdapter{
		appKey:     appKey,
		secretKey:  secretKey,
		senderKey:  senderKey,
		httpClient: &http.Client{Timeout: 10 * time.Second},
		endpoint:   "https://api-alimtalk.cloud.toast.com/alimtalk/v2.3",
		templates:  make(map[string]*Template),
	}
}

// SetHTTPClient는 테스트용 HTTP 클라이언트를 설정합니다.
func (a *NHNCloudAdapter) SetHTTPClient(c HTTPDoer) {
	a.httpClient = c
}

// SetEndpoint는 테스트용 엔드포인트를 설정합니다.
func (a *NHNCloudAdapter) SetEndpoint(endpoint string) {
	a.endpoint = endpoint
}

// Provider는 프로바이더 이름을 반환합니다.
func (a *NHNCloudAdapter) Provider() string { return "nhncloud" }

type nhnRequest struct {
	SenderKey    string            `json:"senderKey"`
	TemplateCode string            `json:"templateCode"`
	Recipients   []nhnRecipient    `json:"recipientList"`
}

type nhnRecipient struct {
	RecipientNo  string            `json:"recipientNo"`
	TemplateData map[string]string `json:"templateParameter,omitempty"`
}

type nhnResponse struct {
	Header struct {
		IsSuccessful bool   `json:"isSuccessful"`
		ResultCode   int    `json:"resultCode"`
		ResultMsg    string `json:"resultMessage"`
	} `json:"header"`
	Message struct {
		RequestID string `json:"requestId"`
	} `json:"message"`
}

// Send는 NHN Cloud Alimtalk API를 호출합니다.
func (a *NHNCloudAdapter) Send(ctx context.Context, msg *Message) (*SendResult, error) {
	if err := ValidateMessage(msg); err != nil {
		return nil, err
	}
	if a.appKey == "" {
		return nil, errors.New("nhncloud appKey not configured")
	}

	body := nhnRequest{
		SenderKey:    a.senderKey,
		TemplateCode: msg.TemplateCode,
		Recipients: []nhnRecipient{
			{
				RecipientNo:  normalizePhone(msg.To),
				TemplateData: msg.Variables,
			},
		},
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal: %w", err)
	}

	url := fmt.Sprintf("%s/appkeys/%s/messages", a.endpoint, a.appKey)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Secret-Key", a.secretKey)

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("nhncloud request failed: %w", err)
	}
	defer resp.Body.Close()

	var nr nhnResponse
	if err := json.NewDecoder(resp.Body).Decode(&nr); err != nil {
		return nil, fmt.Errorf("failed to decode: %w", err)
	}

	if !nr.Header.IsSuccessful {
		return &SendResult{
			Provider:  "nhncloud",
			Status:    "failed",
			ErrorCode: fmt.Sprintf("%d", nr.Header.ResultCode),
			ErrorMsg:  nr.Header.ResultMsg,
			SentAt:    time.Now().UTC(),
		}, fmt.Errorf("nhncloud error: %s", nr.Header.ResultMsg)
	}

	return &SendResult{
		MessageID: nr.Message.RequestID,
		Provider:  "nhncloud",
		Status:    "sent",
		SentAt:    time.Now().UTC(),
	}, nil
}

// HealthCheck는 자격증명을 확인합니다.
func (a *NHNCloudAdapter) HealthCheck(ctx context.Context) error {
	if a.appKey == "" || a.secretKey == "" {
		return errors.New("nhncloud credentials not configured")
	}
	if a.senderKey == "" {
		return errors.New("kakao sender_key not configured")
	}
	return nil
}

// RegisterTemplate는 템플릿을 메모리에 등록합니다 (실제 등록은 NHN Cloud 콘솔).
func (a *NHNCloudAdapter) RegisterTemplate(code, body string, buttons []Button) error {
	a.tmplMu.Lock()
	defer a.tmplMu.Unlock()
	a.templates[code] = &Template{
		Code:       code,
		Body:       body,
		Buttons:    buttons,
		Variables:  ExtractTemplateVariables(body),
		ApprovedAt: time.Now().UTC(),
		Status:     "approved",
	}
	return nil
}

// GetTemplate는 등록된 템플릿을 조회합니다.
func (a *NHNCloudAdapter) GetTemplate(code string) (*Template, error) {
	a.tmplMu.RLock()
	defer a.tmplMu.RUnlock()
	tmpl, ok := a.templates[code]
	if !ok {
		return nil, fmt.Errorf("template %q not registered", code)
	}
	return tmpl, nil
}

// ============================================================================
// KakaoBiz 어댑터 (간략 골격)
// ============================================================================

// KakaoBizAdapter는 KakaoBiz API 어댑터입니다.
type KakaoBizAdapter struct {
	apiKey     string
	senderKey  string
	httpClient HTTPDoer
	templates  map[string]*Template
}

// NewKakaoBizAdapter는 새 KakaoBiz 어댑터를 생성합니다.
func NewKakaoBizAdapter(apiKey, senderKey string) *KakaoBizAdapter {
	return &KakaoBizAdapter{
		apiKey:     apiKey,
		senderKey:  senderKey,
		httpClient: &http.Client{Timeout: 10 * time.Second},
		templates:  make(map[string]*Template),
	}
}

// Provider는 프로바이더 이름을 반환합니다.
func (a *KakaoBizAdapter) Provider() string { return "kakaobiz" }

// Send는 KakaoBiz API를 호출합니다 (골격).
func (a *KakaoBizAdapter) Send(ctx context.Context, msg *Message) (*SendResult, error) {
	if err := ValidateMessage(msg); err != nil {
		return nil, err
	}
	if a.apiKey == "" {
		return nil, errors.New("kakaobiz api_key not configured")
	}
	return &SendResult{
		MessageID: fmt.Sprintf("kakaobiz-%d", time.Now().UnixNano()),
		Provider:  "kakaobiz",
		Status:    "sent",
		SentAt:    time.Now().UTC(),
	}, nil
}

// HealthCheck는 자격증명을 확인합니다.
func (a *KakaoBizAdapter) HealthCheck(_ context.Context) error {
	if a.apiKey == "" {
		return errors.New("kakaobiz api_key not configured")
	}
	return nil
}

// RegisterTemplate는 템플릿을 등록합니다.
func (a *KakaoBizAdapter) RegisterTemplate(code, body string, buttons []Button) error {
	a.templates[code] = &Template{
		Code:      code,
		Body:      body,
		Buttons:   buttons,
		Variables: ExtractTemplateVariables(body),
		Status:    "approved",
		ApprovedAt: time.Now().UTC(),
	}
	return nil
}

// GetTemplate는 템플릿을 조회합니다.
func (a *KakaoBizAdapter) GetTemplate(code string) (*Template, error) {
	tmpl, ok := a.templates[code]
	if !ok {
		return nil, fmt.Errorf("template %q not found", code)
	}
	return tmpl, nil
}

// ============================================================================
// 헬퍼
// ============================================================================

// normalizePhone은 전화번호를 NHN Cloud 형식 (하이픈 제거)으로 정규화합니다.
func normalizePhone(phone string) string {
	// "010-1234-5678" → "01012345678"
	cleaned := strings.ReplaceAll(phone, "-", "")
	cleaned = strings.ReplaceAll(cleaned, " ", "")
	if strings.HasPrefix(cleaned, "+82") {
		// "+821012345678" → "01012345678"
		cleaned = "0" + cleaned[3:]
	}
	return cleaned
}
