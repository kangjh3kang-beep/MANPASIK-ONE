// Package push는 Firebase Cloud Messaging (FCM) 기반 푸시 알림 전송을 구현합니다.
package push

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// FCMConfig는 FCM 연결 설정입니다.
type FCMConfig struct {
	ServerKey string // FCM 서버 키 (Legacy HTTP API)
	ProjectID string // Firebase 프로젝트 ID
}

// FCMClient는 Firebase Cloud Messaging 클라이언트입니다.
type FCMClient struct {
	config     FCMConfig
	httpClient *http.Client
	// TODO: deviceTokenRepo — 사용자별 FCM 디바이스 토큰 조회
}

// NewFCMClient는 FCM 클라이언트를 생성합니다.
func NewFCMClient(cfg FCMConfig) *FCMClient {
	return &FCMClient{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// fcmMessage는 FCM Legacy HTTP API 요청 본문입니다.
type fcmMessage struct {
	To           string            `json:"to,omitempty"`
	Notification *fcmNotification  `json:"notification,omitempty"`
	Data         map[string]string `json:"data,omitempty"`
	Priority     string            `json:"priority,omitempty"`
}

type fcmNotification struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

// SendPush는 FCM을 통해 푸시 알림을 전송합니다.
// 현재는 userID → FCM 디바이스 토큰 매핑이 구현되면 실제 전송됩니다.
// TODO: UserDeviceTokenRepository에서 토큰 조회 후 실제 FCM API 호출
func (c *FCMClient) SendPush(ctx context.Context, userID, title, body, data string) error {
	// TODO: userID → FCM device token 조회
	// deviceToken, err := c.deviceTokenRepo.GetByUserID(ctx, userID)
	// if err != nil || deviceToken == "" {
	//     return fmt.Errorf("사용자 FCM 토큰 없음: %s", userID)
	// }

	// 현재는 로깅만 수행 (디바이스 토큰 저장소 구현 후 실제 전송)
	_ = ctx
	_ = userID
	_ = title
	_ = body
	_ = data

	// 실제 FCM 전송 로직 (디바이스 토큰 확보 후 활성화)
	// return c.sendToFCM(ctx, deviceToken, title, body, data)
	return nil
}

// sendToFCM는 FCM Legacy HTTP API로 메시지를 전송합니다.
// 디바이스 토큰 저장소 연동 후 SendPush에서 호출됩니다.
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
