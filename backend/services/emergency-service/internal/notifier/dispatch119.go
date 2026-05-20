package notifier

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Dispatch119Config는 119 긴급 서비스 API 설정입니다.
type Dispatch119Config struct {
	BaseURL string // 119 API 엔드포인트
	APIKey  string // 인증 키
}

// Dispatch119Notifier는 119 긴급 서비스 API를 호출하는 알림 제공자입니다.
type Dispatch119Notifier struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

// NewDispatch119Notifier는 119 API 기반 알림 제공자를 생성합니다.
func NewDispatch119Notifier(cfg Dispatch119Config) *Dispatch119Notifier {
	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://api.119.go.kr"
	}
	return &Dispatch119Notifier{
		baseURL: cfg.BaseURL,
		apiKey:  cfg.APIKey,
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

type dispatch119Request struct {
	EmergencyType string  `json:"emergency_type"`
	Severity      int     `json:"severity"`
	Location      string  `json:"location"`
	Description   string  `json:"description"`
	MedicalInfo   string  `json:"medical_info,omitempty"`
	CallerID      string  `json:"caller_id"`
	Timestamp     string  `json:"timestamp"`
}

type dispatch119Response struct {
	Success      bool   `json:"success"`
	DispatchID   string `json:"dispatch_id"`
	EstimatedETA int    `json:"estimated_eta"`
	Message      string `json:"message"`
}

// Dispatch119는 119 긴급 서비스에 자동 신고합니다.
func (d *Dispatch119Notifier) Dispatch119(ctx context.Context, alert *EmergencyAlert) (*DispatchResult, error) {
	if d.apiKey == "" {
		return nil, fmt.Errorf("119 API 키가 설정되지 않았습니다")
	}

	reqBody := dispatch119Request{
		EmergencyType: alert.Type,
		Severity:      alert.Severity,
		Location:      alert.Location,
		Description:   alert.Description,
		MedicalInfo:   alert.MedicalInfo,
		CallerID:      alert.UserID,
		Timestamp:     time.Now().Format(time.RFC3339),
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("119 요청 직렬화 실패: %w", err)
	}

	url := d.baseURL + "/v1/dispatch"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("HTTP 요청 생성 실패: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+d.apiKey)

	resp, err := d.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("119 API 호출 실패: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("119 API 응답 오류: status=%d", resp.StatusCode)
	}

	var result dispatch119Response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("119 응답 파싱 실패: %w", err)
	}

	return &DispatchResult{
		Success:      result.Success,
		DispatchID:   result.DispatchID,
		EstimatedETA: result.EstimatedETA,
		Message:      result.Message,
	}, nil
}

// NotifyContacts는 비상 연락처에 SMS를 발송합니다.
func (d *Dispatch119Notifier) NotifyContacts(ctx context.Context, alert *EmergencyAlert, contactPhones []string) error {
	if d.apiKey == "" {
		return fmt.Errorf("119 API 키가 설정되지 않았습니다")
	}

	reqBody := map[string]interface{}{
		"emergency_id": alert.EmergencyID,
		"phones":       contactPhones,
		"message":      fmt.Sprintf("[만파식 응급] %s님의 응급 상황이 발생했습니다. 위치: %s", alert.UserID, alert.Location),
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("알림 요청 직렬화 실패: %w", err)
	}

	url := d.baseURL + "/v1/notify"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonBody))
	if err != nil {
		return fmt.Errorf("HTTP 요청 생성 실패: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+d.apiKey)

	resp, err := d.client.Do(req)
	if err != nil {
		return fmt.Errorf("알림 API 호출 실패: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("알림 API 응답 오류: status=%d", resp.StatusCode)
	}

	return nil
}

// ProviderName은 "dispatch119"를 반환합니다.
func (d *Dispatch119Notifier) ProviderName() string {
	return "dispatch119"
}
