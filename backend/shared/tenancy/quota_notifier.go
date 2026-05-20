package tenancy

import (
	"fmt"
)

// WebhookQuotaNotifier 는 WebhookDispatcher 를 사용해 한도 초과를 알림 (Phase AN-2).
//
// llm.DynamicQuota.SetViolationNotifier 에 주입할 수 있도록 method 시그니처를
// QuotaViolationNotifier 인터페이스에 맞춤 (모듈 간 직접 의존 없이).
type WebhookQuotaNotifier struct {
	dispatcher *WebhookDispatcher
}

// NewWebhookQuotaNotifier 생성.
func NewWebhookQuotaNotifier(d *WebhookDispatcher) *WebhookQuotaNotifier {
	return &WebhookQuotaNotifier{dispatcher: d}
}

// NotifyQuotaExceeded 는 한도 초과 시 비동기 webhook 발송.
//
// dispatcher=nil 이면 no-op. 호출자가 LLM 호출 hot path 에서 호출하므로
// DispatchAsync 사용으로 응답 시간에 영향 없음.
func (n *WebhookQuotaNotifier) NotifyQuotaExceeded(tenantID string, limit int, used int64, limitType string) {
	if n == nil || n.dispatcher == nil {
		return
	}
	n.dispatcher.DispatchAsync(Event{
		Type:     EventQuotaExceeded,
		TenantID: tenantID,
		Payload: map[string]string{
			"limit_type": limitType,
			"limit":      fmt.Sprintf("%d", limit),
			"used":       fmt.Sprintf("%d", used),
		},
	})
}
