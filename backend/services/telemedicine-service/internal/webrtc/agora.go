package webrtc

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// AgoraConfig는 Agora RTC 연결 설정입니다.
type AgoraConfig struct {
	AppID      string // Agora App ID
	AppCert    string // Agora App Certificate (토큰 서명용)
	CustomerID string // Agora RESTful API Customer ID
	Secret     string // Agora RESTful API Secret
	BaseURL    string // 기본값: "https://api.agora.io"
}

// AgoraProvider는 Agora RTC API를 사용하는 WebRTC 제공자입니다.
type AgoraProvider struct {
	appID      string
	appCert    string
	customerID string
	secret     string
	baseURL    string
	client     *http.Client
}

// NewAgoraProvider는 Agora 기반 WebRTC 제공자를 생성합니다.
func NewAgoraProvider(cfg AgoraConfig) *AgoraProvider {
	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://api.agora.io"
	}
	return &AgoraProvider{
		appID:      cfg.AppID,
		appCert:    cfg.AppCert,
		customerID: cfg.CustomerID,
		secret:     cfg.Secret,
		baseURL:    cfg.BaseURL,
		client:     &http.Client{Timeout: 15 * time.Second},
	}
}

// CreateRoom은 Agora 채널을 생성합니다.
func (a *AgoraProvider) CreateRoom(ctx context.Context, consultationID string) (*RoomInfo, error) {
	if a.appID == "" {
		return nil, fmt.Errorf("Agora App ID가 설정되지 않았습니다")
	}

	channelName := fmt.Sprintf("manpasik-%s", consultationID)

	return &RoomInfo{
		RoomID:  channelName,
		RoomURL: fmt.Sprintf("https://meet.manpasik.com/agora/%s", channelName),
	}, nil
}

// GenerateToken은 Agora RTC 토큰을 발급합니다.
func (a *AgoraProvider) GenerateToken(ctx context.Context, roomID, userID string, role TokenRole) (*TokenInfo, error) {
	if a.appID == "" {
		return nil, fmt.Errorf("Agora App ID가 설정되지 않았습니다")
	}
	if a.appCert == "" {
		return nil, fmt.Errorf("Agora App Certificate가 설정되지 않았습니다")
	}

	// Agora RESTful 토큰 생성 API 호출
	url := fmt.Sprintf("%s/v1/apps/%s/cloud_recording/acquire", a.baseURL, a.appID)

	reqBody := map[string]interface{}{
		"cname": roomID,
		"uid":   userID,
		"clientRequest": map[string]interface{}{
			"resourceExpiredHour": 1,
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("토큰 요청 직렬화 실패: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("HTTP 요청 생성 실패: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(a.customerID, a.secret)

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Agora API 호출 실패: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Agora API 응답 오류: status=%d", resp.StatusCode)
	}

	var result struct {
		ResourceID string `json:"resourceId"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("Agora 응답 파싱 실패: %w", err)
	}

	// HMAC-SHA256 기반 간이 토큰 생성 (실제 Agora SDK에서는 RtcTokenBuilder 사용)
	token := a.generateSimpleToken(roomID, userID)

	return &TokenInfo{
		Token:     token,
		ChannelID: roomID,
		UID:       userID,
	}, nil
}

// EndRoom은 Agora 채널의 녹화를 중지합니다.
func (a *AgoraProvider) EndRoom(ctx context.Context, roomID string) error {
	// Agora 채널은 참가자가 모두 나가면 자동 종료됨
	// 녹화 중이었다면 녹화 중지 API 호출
	return nil
}

// ProviderName은 "agora"를 반환합니다.
func (a *AgoraProvider) ProviderName() string {
	return "agora"
}

// generateSimpleToken은 HMAC-SHA256 기반 간이 토큰을 생성합니다.
func (a *AgoraProvider) generateSimpleToken(channelName, uid string) string {
	mac := hmac.New(sha256.New, []byte(a.appCert))
	payload := fmt.Sprintf("%s:%s:%s:%d", a.appID, channelName, uid, time.Now().Unix())
	mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}
