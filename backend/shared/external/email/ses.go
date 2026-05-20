package email

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// ============================================================================
// AWS SES 어댑터 (간략화 버전)
// ============================================================================

// SESAdapter는 AWS SES 어댑터입니다.
//
// 실제 호출은 https://email.{region}.amazonaws.com 의 SES API.
// 인증은 AWS Signature V4를 사용해야 하지만, 본 구현은 AWS SDK 사용을 가정합니다.
// 본 어댑터는 인터페이스만 제공하며 실 호출은 외부 SDK 통합 시 구현됩니다.
type SESAdapter struct {
	region     string
	from       string
	httpClient HTTPDoer
}

// NewSESAdapter는 새 SES 어댑터를 생성합니다.
func NewSESAdapter(region, fromEmail string) *SESAdapter {
	return &SESAdapter{
		region:     region,
		from:       fromEmail,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// SetHTTPClient는 테스트용 HTTP 클라이언트를 설정합니다.
func (a *SESAdapter) SetHTTPClient(c HTTPDoer) {
	a.httpClient = c
}

// Provider는 프로바이더 이름을 반환합니다.
func (a *SESAdapter) Provider() string { return "ses" }

// Send는 SES를 통해 이메일을 전송합니다.
//
// 본 구현은 AWS SDK를 사용하지 않는 기본 발송 골격입니다.
// 운영 시에는 github.com/aws/aws-sdk-go-v2/service/ses 사용을 권장합니다.
func (a *SESAdapter) Send(ctx context.Context, msg *Message) (*SendResult, error) {
	if err := ValidateMessage(msg); err != nil {
		return nil, err
	}
	if a.region == "" {
		return nil, errors.New("ses region not configured")
	}
	if a.from == "" && (msg.From == nil || msg.From.Email == "") {
		return nil, errors.New("ses from address required")
	}

	// 실 SES SDK 통합 시 호출 위치
	// _, err := sesClient.SendEmail(ctx, &ses.SendEmailInput{...})

	// 본 구현은 골격만 제공
	return &SendResult{
		MessageID: fmt.Sprintf("ses-%d", time.Now().UnixNano()),
		Provider:  "ses",
		Status:    "queued",
		SentAt:    time.Now().UTC(),
	}, nil
}

// HealthCheck는 SES 설정을 확인합니다.
func (a *SESAdapter) HealthCheck(ctx context.Context) error {
	if a.region == "" {
		return errors.New("ses region not configured")
	}
	return nil
}
