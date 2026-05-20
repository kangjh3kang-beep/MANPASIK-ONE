// Package push는 Firebase Cloud Messaging (FCM) 기반 푸시 알림 전송을 구현합니다.
package push

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// ============================================================================
// 디바이스 토큰 저장소
// ============================================================================

// DeviceTokenStore는 사용자별 FCM 디바이스 토큰 저장소 인터페이스입니다.
type DeviceTokenStore interface {
	GetTokens(ctx context.Context, userID string) ([]string, error)
	SaveToken(ctx context.Context, userID, token string) error
	RemoveToken(ctx context.Context, userID, token string) error
}

// MemoryTokenStore는 인메모리 디바이스 토큰 저장소입니다.
type MemoryTokenStore struct {
	mu     sync.RWMutex
	tokens map[string][]string // userID → []token
}

// NewMemoryTokenStore는 인메모리 토큰 저장소를 생성합니다.
func NewMemoryTokenStore() *MemoryTokenStore {
	return &MemoryTokenStore{tokens: make(map[string][]string)}
}

// GetTokens는 사용자의 FCM 토큰 목록을 반환합니다.
func (s *MemoryTokenStore) GetTokens(_ context.Context, userID string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.tokens[userID], nil
}

// SaveToken은 사용자의 FCM 토큰을 저장합니다 (중복 방지).
func (s *MemoryTokenStore) SaveToken(_ context.Context, userID, token string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, t := range s.tokens[userID] {
		if t == token {
			return nil
		}
	}
	s.tokens[userID] = append(s.tokens[userID], token)
	return nil
}

// RemoveToken은 사용자의 FCM 토큰을 제거합니다.
func (s *MemoryTokenStore) RemoveToken(_ context.Context, userID, token string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	tokens := s.tokens[userID]
	for i, t := range tokens {
		if t == token {
			s.tokens[userID] = append(tokens[:i], tokens[i+1:]...)
			return nil
		}
	}
	return nil
}

// ============================================================================
// FCM 클라이언트
// ============================================================================

// FCMConfig는 FCM 연결 설정입니다.
type FCMConfig struct {
	ServerKey string // FCM 서버 키 (Legacy HTTP API)
	ProjectID string // Firebase 프로젝트 ID
}

// FCMClient는 Firebase Cloud Messaging 클라이언트입니다.
type FCMClient struct {
	config     FCMConfig
	httpClient *http.Client
	tokenStore DeviceTokenStore
}

// NewFCMClient는 FCM 클라이언트를 생성합니다.
func NewFCMClient(cfg FCMConfig) *FCMClient {
	return &FCMClient{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		tokenStore: NewMemoryTokenStore(),
	}
}

// SetTokenStore는 디바이스 토큰 저장소를 설정합니다.
func (c *FCMClient) SetTokenStore(store DeviceTokenStore) {
	c.tokenStore = store
}

// RegisterToken은 사용자의 FCM 디바이스 토큰을 등록합니다.
func (c *FCMClient) RegisterToken(ctx context.Context, userID, token string) error {
	return c.tokenStore.SaveToken(ctx, userID, token)
}

// UnregisterToken은 사용자의 FCM 디바이스 토큰을 제거합니다.
func (c *FCMClient) UnregisterToken(ctx context.Context, userID, token string) error {
	return c.tokenStore.RemoveToken(ctx, userID, token)
}

// fcmMessage는 FCM Legacy HTTP API 요청 본문입니다.
type fcmMessage struct {
	To           string            `json:"to,omitempty"`
	Topic        string            `json:"topic,omitempty"`
	Notification *fcmNotification  `json:"notification,omitempty"`
	Data         map[string]string `json:"data,omitempty"`
	Priority     string            `json:"priority,omitempty"`
}

type fcmNotification struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

// SendPush는 FCM을 통해 사용자에게 푸시 알림을 전송합니다.
// 토큰이 없으면 전송을 건너뛰고 nil을 반환합니다 (H6 실패격리).
func (c *FCMClient) SendPush(ctx context.Context, userID, title, body, data string) error {
	tokens, err := c.tokenStore.GetTokens(ctx, userID)
	if err != nil {
		return fmt.Errorf("FCM 토큰 조회 실패: %w", err)
	}

	if len(tokens) == 0 {
		// 디바이스 토큰 미등록 — 조용히 무시 (H6 실패격리)
		return nil
	}

	// 모든 디바이스에 전송
	var lastErr error
	for _, token := range tokens {
		if err := c.sendToFCM(ctx, token, title, body, data); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

// SendToTopic은 FCM 토픽으로 푸시 알림을 전송합니다.
func (c *FCMClient) SendToTopic(ctx context.Context, topic, title, body, data string) error {
	msg := fcmMessage{
		Topic: topic,
		Notification: &fcmNotification{
			Title: title,
			Body:  body,
		},
		Priority: "high",
	}
	if data != "" {
		msg.Data = map[string]string{"payload": data}
	}
	return c.sendMessage(ctx, msg)
}

// sendToFCM는 FCM Legacy HTTP API로 개별 디바이스에 메시지를 전송합니다.
func (c *FCMClient) sendToFCM(ctx context.Context, deviceToken, title, body, data string) error {
	msg := fcmMessage{
		To: deviceToken,
		Notification: &fcmNotification{
			Title: title,
			Body:  body,
		},
		Priority: "high",
	}
	if data != "" {
		msg.Data = map[string]string{"payload": data}
	}
	return c.sendMessage(ctx, msg)
}

// sendMessage는 FCM HTTP API로 메시지를 전송합니다.
func (c *FCMClient) sendMessage(ctx context.Context, msg fcmMessage) error {
	jsonBody, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("FCM 메시지 직렬화 실패: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://fcm.googleapis.com/fcm/send", bytes.NewReader(jsonBody))
	if err != nil {
		return fmt.Errorf("FCM 요청 생성 실패: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "key="+c.config.ServerKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("FCM 전송 실패: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("FCM 응답 오류: status=%d", resp.StatusCode)
	}

	return nil
}

// ============================================================================
// 대체 전송기 (FCM 미설정 시)
// ============================================================================

// LogPushSender는 전송 내용을 로그에 기록하는 푸시 전송기입니다.
// FCM 미설정 시 사용하며, 디버깅과 테스트에 유용합니다.
type LogPushSender struct{}

// NewLogPushSender는 로그 기반 푸시 전송기를 생성합니다.
func NewLogPushSender() *LogPushSender {
	return &LogPushSender{}
}

// SendPush는 알림 내용을 로그에 기록합니다.
func (l *LogPushSender) SendPush(_ context.Context, userID, title, body, _ string) error {
	log.Printf("[PUSH-LOG] user=%s title=%q body=%q (FCM 미설정, 로그만 기록)", userID, title, body)
	return nil
}

// NoopPushSender는 아무 동작도 하지 않는 푸시 전송기입니다 (개발용).
type NoopPushSender struct{}

// NewNoopPushSender는 No-op 푸시 전송기를 생성합니다.
func NewNoopPushSender() *NoopPushSender {
	return &NoopPushSender{}
}

// SendPush는 아무 동작도 하지 않습니다.
func (n *NoopPushSender) SendPush(_ context.Context, _, _, _, _ string) error {
	return nil
}

// NoopEmailSender는 아무 동작도 하지 않는 이메일 전송기입니다 (개발용).
type NoopEmailSender struct{}

// NewNoopEmailSender는 No-op 이메일 전송기를 생성합니다.
func NewNoopEmailSender() *NoopEmailSender {
	return &NoopEmailSender{}
}

// SendEmail은 아무 동작도 하지 않습니다.
func (n *NoopEmailSender) SendEmail(_ context.Context, _, _, _ string) error {
	return nil
}
