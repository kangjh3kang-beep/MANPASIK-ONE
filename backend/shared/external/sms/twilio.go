package sms

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// ============================================================================
// Twilio 어댑터
// ============================================================================

// TwilioAdapter는 Twilio SMS 어댑터입니다.
//
// 실제 호출은 https://api.twilio.com/2010-04-01/Accounts/{AccountSID}/Messages.json
type TwilioAdapter struct {
	accountSID string
	authToken  string
	from       string
	httpClient HTTPDoer
	endpoint   string
}

// HTTPDoer는 테스트 가능한 HTTP 클라이언트 인터페이스입니다.
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// NewTwilioAdapter는 새 Twilio 어댑터를 생성합니다.
func NewTwilioAdapter(accountSID, authToken, from string) *TwilioAdapter {
	return &TwilioAdapter{
		accountSID: accountSID,
		authToken:  authToken,
		from:       from,
		httpClient: &http.Client{Timeout: 10 * time.Second},
		endpoint:   "https://api.twilio.com/2010-04-01",
	}
}

// SetHTTPClient는 테스트용 HTTP 클라이언트를 설정합니다.
func (a *TwilioAdapter) SetHTTPClient(c HTTPDoer) {
	a.httpClient = c
}

// SetEndpoint는 테스트용 엔드포인트를 설정합니다.
func (a *TwilioAdapter) SetEndpoint(endpoint string) {
	a.endpoint = endpoint
}

// Provider는 프로바이더 이름을 반환합니다.
func (a *TwilioAdapter) Provider() string { return "twilio" }

// Send는 Twilio API를 호출하여 SMS를 전송합니다.
func (a *TwilioAdapter) Send(ctx context.Context, msg *Message) (*SendResult, error) {
	if err := ValidateMessage(msg); err != nil {
		return nil, err
	}
	if a.accountSID == "" || a.authToken == "" {
		return nil, errors.New("twilio credentials not configured")
	}

	from := msg.From
	if from == "" {
		from = a.from
	}

	form := url.Values{}
	form.Set("From", from)
	form.Set("To", msg.To)
	form.Set("Body", msg.Body)

	endpoint := fmt.Sprintf("%s/Accounts/%s/Messages.json", a.endpoint, a.accountSID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(a.accountSID, a.authToken)

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("twilio request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return &SendResult{
			Provider:  "twilio",
			Status:    "failed",
			ErrorCode: fmt.Sprintf("HTTP_%d", resp.StatusCode),
			ErrorMsg:  resp.Status,
			SentAt:    time.Now().UTC(),
		}, fmt.Errorf("twilio returned %d", resp.StatusCode)
	}

	return &SendResult{
		MessageID: fmt.Sprintf("twilio-%d", time.Now().UnixNano()),
		Provider:  "twilio",
		Status:    "sent",
		SentAt:    time.Now().UTC(),
		Currency:  "USD",
		Cost:      0.0075, // 일반 미국 SMS 단가 (대략)
	}, nil
}

// HealthCheck는 Twilio 자격증명을 확인합니다.
func (a *TwilioAdapter) HealthCheck(ctx context.Context) error {
	if a.accountSID == "" || a.authToken == "" {
		return errors.New("twilio credentials not configured")
	}
	return nil
}
