// Package fcm은 Firebase Cloud Messaging HTTP v1 API 어댑터입니다.
//
// FCM legacy API는 deprecated이며, 본 구현은 HTTP v1 API (OAuth2 인증)를 사용합니다.
// 환경변수: GOOGLE_APPLICATION_CREDENTIALS (서비스 계정 JSON 경로)
package fcm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// ============================================================================
// 도메인 모델
// ============================================================================

// Notification은 알림 본문입니다.
type Notification struct {
	Title    string
	Body     string
	ImageURL string // 옵션 (큰 이미지)
}

// Android는 Android 전용 옵션입니다.
type Android struct {
	Priority      string            // "high", "normal"
	TTL           time.Duration     // 메시지 보관 기간
	CollapseKey   string            // 동일 키 메시지 그룹화
	ChannelID     string            // 알림 채널 (Android 8+)
	Sound         string            // 사운드 파일명
	ClickAction   string            // 탭 시 액션
	Data          map[string]string
}

// IOS는 iOS APNs 전용 옵션입니다.
type IOS struct {
	BadgeCount int
	Sound      string
	Category   string
	ThreadID   string
	ContentAvailable bool // 백그라운드 갱신
	MutableContent   bool // 알림 내용 변경 가능
}

// Message는 FCM 메시지입니다.
type Message struct {
	Token        string             // 단일 디바이스 토큰
	Topic        string             // 토픽 (예: "all-users")
	Condition    string             // 조건 ("'topic1' in topics && 'topic2' in topics")
	Notification *Notification
	Data         map[string]string  // 데이터 페이로드 (백그라운드 처리용)
	Android      *Android
	IOS          *IOS
	DryRun       bool               // 실제 발송 없이 검증만
}

// SendResult는 발송 결과입니다.
type SendResult struct {
	MessageID string
	Provider  string
	SentAt    time.Time
	Status    string // "sent", "failed", "invalid_token"
	ErrorCode string
	ErrorMsg  string
}

// MulticastResult는 다중 발송 결과입니다.
type MulticastResult struct {
	SuccessCount int
	FailureCount int
	Results      []*SendResult
	InvalidTokens []string // 무효 토큰 (재발송 제외 + DB에서 정리)
}

// ============================================================================
// 어댑터 인터페이스
// ============================================================================

// Adapter는 FCM 어댑터 인터페이스입니다.
type Adapter interface {
	Send(ctx context.Context, msg *Message) (*SendResult, error)
	SendMulticast(ctx context.Context, msg *Message, tokens []string) (*MulticastResult, error)
	Provider() string
	HealthCheck(ctx context.Context) error
}

// ============================================================================
// 검증
// ============================================================================

// ValidateMessage는 FCM 메시지를 검증합니다.
func ValidateMessage(msg *Message) error {
	if msg == nil {
		return errors.New("message is nil")
	}
	// 타겟 검증: token, topic, condition 중 정확히 하나
	count := 0
	if msg.Token != "" {
		count++
	}
	if msg.Topic != "" {
		count++
	}
	if msg.Condition != "" {
		count++
	}
	if count == 0 {
		return errors.New("one of token/topic/condition required")
	}
	if count > 1 {
		return errors.New("only one of token/topic/condition allowed")
	}
	if msg.Notification == nil && len(msg.Data) == 0 {
		return errors.New("notification or data required")
	}
	if msg.Notification != nil && msg.Notification.Title == "" && msg.Notification.Body == "" {
		return errors.New("notification title or body required")
	}
	return nil
}

// ============================================================================
// 팩토리
// ============================================================================

// NewFromEnv는 환경변수에 따라 어댑터를 생성합니다.
//
// FCM_ENABLED=true 이고 GOOGLE_APPLICATION_CREDENTIALS, FCM_PROJECT_ID가 설정되어 있으면 v1 어댑터.
// 그렇지 않으면 Noop.
func NewFromEnv() Adapter {
	enabled := strings.ToLower(os.Getenv("FCM_ENABLED")) == "true"
	projectID := os.Getenv("FCM_PROJECT_ID")
	credPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")

	if enabled && projectID != "" && credPath != "" {
		return NewV1Adapter(projectID, credPath)
	}
	return NewNoopAdapter()
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

// Send는 메시지를 메모리에 저장합니다.
func (a *NoopAdapter) Send(_ context.Context, msg *Message) (*SendResult, error) {
	if err := ValidateMessage(msg); err != nil {
		return nil, err
	}
	a.mu.Lock()
	a.messages = append(a.messages, msg)
	a.mu.Unlock()
	return &SendResult{
		MessageID: fmt.Sprintf("noop-fcm-%d", time.Now().UnixNano()),
		Provider:  "noop",
		Status:    "sent",
		SentAt:    time.Now().UTC(),
	}, nil
}

// SendMulticast는 다중 발송을 시뮬레이션합니다.
func (a *NoopAdapter) SendMulticast(ctx context.Context, msg *Message, tokens []string) (*MulticastResult, error) {
	result := &MulticastResult{Results: make([]*SendResult, 0, len(tokens))}
	for _, token := range tokens {
		clone := *msg
		clone.Token = token
		clone.Topic = ""
		clone.Condition = ""
		r, err := a.Send(ctx, &clone)
		if err != nil {
			result.FailureCount++
			result.Results = append(result.Results, &SendResult{
				Provider:  "noop",
				Status:    "failed",
				ErrorMsg:  err.Error(),
			})
			continue
		}
		result.SuccessCount++
		result.Results = append(result.Results, r)
	}
	return result, nil
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

// ============================================================================
// FCM v1 어댑터
// ============================================================================

// HTTPDoer는 테스트 가능한 HTTP 클라이언트 인터페이스입니다.
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// TokenProvider는 OAuth2 액세스 토큰 제공자입니다.
type TokenProvider interface {
	AccessToken(ctx context.Context) (string, error)
}

// V1Adapter는 FCM HTTP v1 API 어댑터입니다.
//
// 엔드포인트: https://fcm.googleapis.com/v1/projects/{PROJECT_ID}/messages:send
// 인증: Bearer 토큰 (서비스 계정 OAuth2)
type V1Adapter struct {
	projectID  string
	credPath   string
	endpoint   string
	httpClient HTTPDoer
	tokenProv  TokenProvider
}

// NewV1Adapter는 새 FCM v1 어댑터를 생성합니다.
func NewV1Adapter(projectID, credentialsPath string) *V1Adapter {
	return &V1Adapter{
		projectID:  projectID,
		credPath:   credentialsPath,
		endpoint:   fmt.Sprintf("https://fcm.googleapis.com/v1/projects/%s/messages:send", projectID),
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// SetHTTPClient는 테스트용 HTTP 클라이언트를 설정합니다.
func (a *V1Adapter) SetHTTPClient(c HTTPDoer) {
	a.httpClient = c
}

// SetEndpoint는 테스트용 엔드포인트를 설정합니다.
func (a *V1Adapter) SetEndpoint(endpoint string) {
	a.endpoint = endpoint
}

// SetTokenProvider는 액세스 토큰 제공자를 설정합니다.
func (a *V1Adapter) SetTokenProvider(p TokenProvider) {
	a.tokenProv = p
}

// Provider는 프로바이더 이름을 반환합니다.
func (a *V1Adapter) Provider() string { return "fcm_v1" }

// fcmRequest는 FCM v1 API 요청 본문입니다.
type fcmRequest struct {
	Message struct {
		Token        string                 `json:"token,omitempty"`
		Topic        string                 `json:"topic,omitempty"`
		Condition    string                 `json:"condition,omitempty"`
		Notification *fcmNotification       `json:"notification,omitempty"`
		Data         map[string]string      `json:"data,omitempty"`
		Android      map[string]interface{} `json:"android,omitempty"`
		APNS         map[string]interface{} `json:"apns,omitempty"`
	} `json:"message"`
	ValidateOnly bool `json:"validate_only,omitempty"`
}

type fcmNotification struct {
	Title string `json:"title,omitempty"`
	Body  string `json:"body,omitempty"`
	Image string `json:"image,omitempty"`
}

type fcmResponse struct {
	Name  string `json:"name,omitempty"`
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error,omitempty"`
}

// Send는 FCM v1 API를 호출하여 메시지를 전송합니다.
func (a *V1Adapter) Send(ctx context.Context, msg *Message) (*SendResult, error) {
	if err := ValidateMessage(msg); err != nil {
		return nil, err
	}
	if a.tokenProv == nil {
		return nil, errors.New("token provider not configured")
	}

	accessToken, err := a.tokenProv.AccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	body := buildFCMRequest(msg)
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.endpoint, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fcm request failed: %w", err)
	}
	defer resp.Body.Close()

	var fr fcmResponse
	if err := json.NewDecoder(resp.Body).Decode(&fr); err != nil {
		return nil, fmt.Errorf("failed to decode: %w", err)
	}

	if resp.StatusCode >= 400 {
		status := "failed"
		// UNREGISTERED, INVALID_ARGUMENT 등 무효 토큰 식별
		if fr.Error.Status == "UNREGISTERED" || fr.Error.Status == "INVALID_ARGUMENT" {
			status = "invalid_token"
		}
		return &SendResult{
			Provider:  "fcm_v1",
			Status:    status,
			ErrorCode: fr.Error.Status,
			ErrorMsg:  fr.Error.Message,
			SentAt:    time.Now().UTC(),
		}, fmt.Errorf("fcm error: %s", fr.Error.Message)
	}

	return &SendResult{
		MessageID: fr.Name,
		Provider:  "fcm_v1",
		Status:    "sent",
		SentAt:    time.Now().UTC(),
	}, nil
}

// SendMulticast는 다수의 토큰에 동일 메시지를 발송합니다 (순차).
//
// FCM v1은 batch API를 별도 제공하지만 본 구현은 단순 순차 발송을 합니다.
func (a *V1Adapter) SendMulticast(ctx context.Context, msg *Message, tokens []string) (*MulticastResult, error) {
	result := &MulticastResult{
		Results: make([]*SendResult, 0, len(tokens)),
	}

	for _, token := range tokens {
		clone := *msg
		clone.Token = token
		clone.Topic = ""
		clone.Condition = ""

		r, err := a.Send(ctx, &clone)
		if err != nil {
			result.FailureCount++
			if r != nil && r.Status == "invalid_token" {
				result.InvalidTokens = append(result.InvalidTokens, token)
			}
			result.Results = append(result.Results, r)
			continue
		}
		result.SuccessCount++
		result.Results = append(result.Results, r)
	}
	return result, nil
}

// HealthCheck는 자격증명을 확인합니다.
func (a *V1Adapter) HealthCheck(ctx context.Context) error {
	if a.projectID == "" {
		return errors.New("fcm project_id not configured")
	}
	if a.tokenProv == nil {
		return errors.New("token provider not configured")
	}
	return nil
}

// buildFCMRequest는 Message를 FCM v1 요청 본문으로 변환합니다.
func buildFCMRequest(msg *Message) *fcmRequest {
	req := &fcmRequest{
		ValidateOnly: msg.DryRun,
	}
	req.Message.Token = msg.Token
	req.Message.Topic = msg.Topic
	req.Message.Condition = msg.Condition
	req.Message.Data = msg.Data

	if msg.Notification != nil {
		req.Message.Notification = &fcmNotification{
			Title: msg.Notification.Title,
			Body:  msg.Notification.Body,
			Image: msg.Notification.ImageURL,
		}
	}

	if msg.Android != nil {
		req.Message.Android = buildAndroidConfig(msg.Android)
	}

	if msg.IOS != nil {
		req.Message.APNS = buildAPNSConfig(msg.IOS)
	}

	return req
}

func buildAndroidConfig(a *Android) map[string]interface{} {
	cfg := map[string]interface{}{}
	if a.Priority != "" {
		cfg["priority"] = strings.ToUpper(a.Priority)
	}
	if a.TTL > 0 {
		cfg["ttl"] = fmt.Sprintf("%ds", int(a.TTL.Seconds()))
	}
	if a.CollapseKey != "" {
		cfg["collapse_key"] = a.CollapseKey
	}
	if a.Data != nil {
		cfg["data"] = a.Data
	}

	if a.ChannelID != "" || a.Sound != "" || a.ClickAction != "" {
		notif := map[string]interface{}{}
		if a.ChannelID != "" {
			notif["channel_id"] = a.ChannelID
		}
		if a.Sound != "" {
			notif["sound"] = a.Sound
		}
		if a.ClickAction != "" {
			notif["click_action"] = a.ClickAction
		}
		cfg["notification"] = notif
	}

	return cfg
}

func buildAPNSConfig(i *IOS) map[string]interface{} {
	aps := map[string]interface{}{}
	if i.BadgeCount > 0 {
		aps["badge"] = i.BadgeCount
	}
	if i.Sound != "" {
		aps["sound"] = i.Sound
	}
	if i.Category != "" {
		aps["category"] = i.Category
	}
	if i.ThreadID != "" {
		aps["thread-id"] = i.ThreadID
	}
	if i.ContentAvailable {
		aps["content-available"] = 1
	}
	if i.MutableContent {
		aps["mutable-content"] = 1
	}

	return map[string]interface{}{
		"payload": map[string]interface{}{
			"aps": aps,
		},
	}
}

// ============================================================================
// 정적 토큰 제공자 (테스트/단순 환경용)
// ============================================================================

// StaticTokenProvider는 정적 액세스 토큰을 반환합니다 (테스트용).
type StaticTokenProvider struct {
	Token string
}

// AccessToken은 토큰을 반환합니다.
func (s *StaticTokenProvider) AccessToken(_ context.Context) (string, error) {
	if s.Token == "" {
		return "", errors.New("static token is empty")
	}
	return s.Token, nil
}
