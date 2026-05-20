// Package notifier는 응급 상황 알림 제공자 인터페이스를 정의합니다.
package notifier

import "context"

// EmergencyAlert는 응급 알림 정보입니다.
type EmergencyAlert struct {
	EmergencyID string
	UserID      string
	Type        string
	Severity    int
	Location    string
	Description string
	MedicalInfo string
	ContactIDs  []string
}

// DispatchResult는 응급 서비스 디스패치 결과입니다.
type DispatchResult struct {
	Success      bool
	DispatchID   string // 119 접수 번호
	EstimatedETA int    // 예상 도착 시간 (분)
	Message      string
}

// EmergencyNotifier는 응급 상황 외부 알림 인터페이스입니다.
type EmergencyNotifier interface {
	// Dispatch119는 119 긴급 서비스에 신고합니다.
	Dispatch119(ctx context.Context, alert *EmergencyAlert) (*DispatchResult, error)

	// NotifyContacts는 비상 연락처에 알림을 발송합니다.
	NotifyContacts(ctx context.Context, alert *EmergencyAlert, contactPhones []string) error

	// ProviderName은 제공자 이름을 반환합니다.
	ProviderName() string
}

// LogNotifier는 실제 API 호출 없이 로그만 기록하는 폴백 구현입니다.
type LogNotifier struct{}

// NewLogNotifier는 LogNotifier를 생성합니다.
func NewLogNotifier() *LogNotifier {
	return &LogNotifier{}
}

// Dispatch119는 로그만 기록하고 성공을 반환합니다.
func (n *LogNotifier) Dispatch119(_ context.Context, alert *EmergencyAlert) (*DispatchResult, error) {
	return &DispatchResult{
		Success:    true,
		DispatchID: "LOG-" + alert.EmergencyID,
		Message:    "119 API 미설정 — 로그 기록만 수행",
	}, nil
}

// NotifyContacts는 로그만 기록합니다.
func (n *LogNotifier) NotifyContacts(_ context.Context, _ *EmergencyAlert, _ []string) error {
	return nil
}

// ProviderName은 "log"를 반환합니다.
func (n *LogNotifier) ProviderName() string {
	return "log"
}
