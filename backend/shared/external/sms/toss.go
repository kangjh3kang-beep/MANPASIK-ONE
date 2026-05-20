package sms

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// ============================================================================
// Toss SMS 어댑터 (한국)
// ============================================================================

// TossAdapter는 Toss SMS 어댑터입니다.
//
// 실제 호출은 https://api.tosspayments.com/v1/sms (가상 엔드포인트)
type TossAdapter struct {
	apiKey     string
	from       string
	httpClient HTTPDoer
	endpoint   string
}

// NewTossAdapter는 새 Toss 어댑터를 생성합니다.
func NewTossAdapter(apiKey, from string) *TossAdapter {
	return &TossAdapter{
		apiKey:     apiKey,
		from:       from,
		httpClient: &http.Client{Timeout: 10 * time.Second},
		endpoint:   "https://api.tosspayments.com/v1/sms",
	}
}

// SetHTTPClient는 테스트용 HTTP 클라이언트를 설정합니다.
func (a *TossAdapter) SetHTTPClient(c HTTPDoer) {
	a.httpClient = c
}

// SetEndpoint는 테스트용 엔드포인트를 설정합니다.
func (a *TossAdapter) SetEndpoint(endpoint string) {
	a.endpoint = endpoint
}

// Provider는 프로바이더 이름을 반환합니다.
func (a *TossAdapter) Provider() string { return "toss" }

// tossRequest는 Toss SMS 요청 본문입니다.
type tossRequest struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject,omitempty"`
	Body    string `json:"body"`
	IsLMS   bool   `json:"is_lms"`
}

// tossResponse는 Toss SMS 응답입니다.
type tossResponse struct {
	MessageID string  `json:"message_id"`
	Status    string  `json:"status"`
	Cost      float64 `json:"cost"`
	ErrorCode string  `json:"error_code,omitempty"`
	ErrorMsg  string  `json:"error_msg,omitempty"`
}

// Send는 Toss API를 호출하여 SMS를 전송합니다.
func (a *TossAdapter) Send(ctx context.Context, msg *Message) (*SendResult, error) {
	if err := ValidateMessage(msg); err != nil {
		return nil, err
	}
	if a.apiKey == "" {
		return nil, errors.New("toss api_key not configured")
	}

	from := msg.From
	if from == "" {
		from = a.from
	}

	body := tossRequest{
		From:    from,
		To:      msg.To,
		Subject: msg.Subject,
		Body:    msg.Body,
		IsLMS:   msg.IsLMS,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.endpoint, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+a.apiKey)

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("toss request failed: %w", err)
	}
	defer resp.Body.Close()

	var tr tossResponse
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return &SendResult{
			Provider:  "toss",
			Status:    "failed",
			ErrorCode: tr.ErrorCode,
			ErrorMsg:  tr.ErrorMsg,
			SentAt:    time.Now().UTC(),
		}, fmt.Errorf("toss returned %d: %s", resp.StatusCode, tr.ErrorMsg)
	}

	return &SendResult{
		MessageID: tr.MessageID,
		Provider:  "toss",
		Status:    tr.Status,
		Cost:      tr.Cost,
		Currency:  "KRW",
		SentAt:    time.Now().UTC(),
	}, nil
}

// HealthCheck는 Toss API 자격증명을 확인합니다.
func (a *TossAdapter) HealthCheck(ctx context.Context) error {
	if a.apiKey == "" {
		return errors.New("toss api_key not configured")
	}
	return nil
}
