// Package pg는 PG사(Toss 등) 결제 연동 구현을 제공합니다.
package pg

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const defaultTossAPIURL = "https://api.tosspayments.com"

// TossClient는 Toss Payments API 클라이언트입니다.
type TossClient struct {
	baseURL    string
	secretKey  string
	httpClient *http.Client
}

// NewTossClient는 Toss 클라이언트를 생성합니다.
// secretKey는 환경변수 등으로 주입하며 코드에 기재하지 않습니다.
// apiURL이 빈 문자열이면 defaultTossAPIURL을 사용합니다.
func NewTossClient(secretKey, apiURL string) *TossClient {
	if apiURL == "" {
		apiURL = defaultTossAPIURL
	}
	return &TossClient{
		baseURL:   strings.TrimSuffix(apiURL, "/"),
		secretKey: secretKey,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// confirmRequest는 Toss 결제 승인 요청 본문입니다.
type confirmRequest struct {
	PaymentKey string `json:"paymentKey"`
	OrderID    string `json:"orderId"`
	Amount     int32  `json:"amount"`
}

// confirmResponse는 Toss 결제 승인 응답 (필요 필드만).
type confirmResponse struct {
	PaymentKey string `json:"paymentKey"`
	OrderID    string `json:"orderId"`
	TotalAmount int32 `json:"totalAmount"`
}

// Confirm은 Toss에 결제 승인을 요청합니다.
// 성공 시 Toss에서 반환한 paymentKey를 pgTransactionID로 반환합니다 (취소 시 동일 키 사용).
func (c *TossClient) Confirm(ctx context.Context, paymentKey, orderId string, amountKRW int32) (pgTransactionID string, err error) {
	if paymentKey == "" || orderId == "" {
		return "", fmt.Errorf("paymentKey와 orderId는 필수입니다")
	}
	if amountKRW <= 0 {
		return "", fmt.Errorf("amount는 0보다 커야 합니다")
	}

	body := confirmRequest{PaymentKey: paymentKey, OrderID: orderId, Amount: amountKRW}
	raw, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/payments/confirm", bytes.NewReader(raw))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(c.secretKey+":")))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("toss confirm failed: status=%d body=%s", resp.StatusCode, string(respBody))
	}

	var parsed confirmResponse
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return "", fmt.Errorf("toss confirm response parse: %w", err)
	}
	// 취소 API에 사용할 키로 paymentKey 반환
	if parsed.PaymentKey != "" {
		return parsed.PaymentKey, nil
	}
	return paymentKey, nil
}

// cancelRequest는 Toss 결제 취소 요청 본문입니다.
type cancelRequest struct {
	CancelReason string `json:"cancelReason"`
}

// Cancel은 Toss에 결제 취소(환불)를 요청합니다.
// paymentKey는 승인 시 받은 값(PgTransactionID로 저장된 값)과 동일해야 합니다.
func (c *TossClient) Cancel(ctx context.Context, paymentKey, reason string) error {
	if paymentKey == "" {
		return fmt.Errorf("paymentKey는 필수입니다")
	}
	if reason == "" {
		reason = "고객 요청"
	}

	body := cancelRequest{CancelReason: reason}
	raw, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/payments/"+paymentKey+"/cancel", bytes.NewReader(raw))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(c.secretKey+":")))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("toss cancel failed: status=%d body=%s", resp.StatusCode, string(respBody))
	}
	return nil
}
