package email

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
// SendGrid 어댑터
// ============================================================================

// HTTPDoer는 테스트 가능한 HTTP 클라이언트 인터페이스입니다.
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// SendGridAdapter는 SendGrid 어댑터입니다.
type SendGridAdapter struct {
	apiKey     string
	from       *Address
	httpClient HTTPDoer
	endpoint   string
}

// NewSendGridAdapter는 새 SendGrid 어댑터를 생성합니다.
func NewSendGridAdapter(apiKey, fromEmail, fromName string) *SendGridAdapter {
	return &SendGridAdapter{
		apiKey:     apiKey,
		from:       &Address{Email: fromEmail, Name: fromName},
		httpClient: &http.Client{Timeout: 10 * time.Second},
		endpoint:   "https://api.sendgrid.com/v3/mail/send",
	}
}

// SetHTTPClient는 테스트용 HTTP 클라이언트를 설정합니다.
func (a *SendGridAdapter) SetHTTPClient(c HTTPDoer) {
	a.httpClient = c
}

// SetEndpoint는 테스트용 엔드포인트를 설정합니다.
func (a *SendGridAdapter) SetEndpoint(endpoint string) {
	a.endpoint = endpoint
}

// Provider는 프로바이더 이름을 반환합니다.
func (a *SendGridAdapter) Provider() string { return "sendgrid" }

// sendGridRequest는 SendGrid API 요청 본문 (v3 mail/send)입니다.
type sendGridRequest struct {
	Personalizations []sgPersonalization `json:"personalizations"`
	From             sgEmail             `json:"from"`
	Subject          string              `json:"subject,omitempty"`
	Content          []sgContent         `json:"content,omitempty"`
	TemplateID       string              `json:"template_id,omitempty"`
}

type sgPersonalization struct {
	To                  []sgEmail         `json:"to"`
	Cc                  []sgEmail         `json:"cc,omitempty"`
	Bcc                 []sgEmail         `json:"bcc,omitempty"`
	Subject             string            `json:"subject,omitempty"`
	DynamicTemplateData map[string]string `json:"dynamic_template_data,omitempty"`
}

type sgEmail struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

type sgContent struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// Send는 SendGrid API를 호출하여 이메일을 전송합니다.
func (a *SendGridAdapter) Send(ctx context.Context, msg *Message) (*SendResult, error) {
	if err := ValidateMessage(msg); err != nil {
		return nil, err
	}
	if a.apiKey == "" {
		return nil, errors.New("sendgrid api_key not configured")
	}

	from := msg.From
	if from == nil {
		from = a.from
	}

	body := sendGridRequest{
		From:    sgEmail{Email: from.Email, Name: from.Name},
		Subject: msg.Subject,
	}

	personalization := sgPersonalization{
		To: convertAddresses(msg.To),
	}
	if len(msg.Cc) > 0 {
		personalization.Cc = convertAddresses(msg.Cc)
	}
	if len(msg.Bcc) > 0 {
		personalization.Bcc = convertAddresses(msg.Bcc)
	}
	if msg.TemplateID != "" {
		body.TemplateID = msg.TemplateID
		personalization.DynamicTemplateData = msg.TemplateData
	}
	body.Personalizations = []sgPersonalization{personalization}

	if msg.TextBody != "" {
		body.Content = append(body.Content, sgContent{Type: "text/plain", Value: msg.TextBody})
	}
	if msg.HTMLBody != "" {
		body.Content = append(body.Content, sgContent{Type: "text/html", Value: msg.HTMLBody})
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
		return nil, fmt.Errorf("sendgrid request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return &SendResult{
			Provider:  "sendgrid",
			Status:    "failed",
			ErrorCode: fmt.Sprintf("HTTP_%d", resp.StatusCode),
			SentAt:    time.Now().UTC(),
		}, fmt.Errorf("sendgrid returned %d", resp.StatusCode)
	}

	// SendGrid는 X-Message-Id 헤더로 메시지 ID 반환
	messageID := resp.Header.Get("X-Message-Id")
	if messageID == "" {
		messageID = fmt.Sprintf("sg-%d", time.Now().UnixNano())
	}

	return &SendResult{
		MessageID: messageID,
		Provider:  "sendgrid",
		Status:    "queued",
		SentAt:    time.Now().UTC(),
	}, nil
}

// HealthCheck는 자격증명을 확인합니다.
func (a *SendGridAdapter) HealthCheck(ctx context.Context) error {
	if a.apiKey == "" {
		return errors.New("sendgrid api_key not configured")
	}
	return nil
}

func convertAddresses(addrs []*Address) []sgEmail {
	out := make([]sgEmail, len(addrs))
	for i, a := range addrs {
		out[i] = sgEmail{Email: a.Email, Name: a.Name}
	}
	return out
}
