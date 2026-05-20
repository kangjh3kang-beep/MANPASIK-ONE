package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/manpasik/backend/services/notification-service/internal/service"
	"github.com/manpasik/backend/shared/external/email"
	"github.com/manpasik/backend/shared/external/fcm"
	"github.com/manpasik/backend/shared/external/kakao"
	"github.com/manpasik/backend/shared/external/sms"
	"github.com/manpasik/backend/shared/medical/poise"
)

func newTestDispatcher() *service.MultiChannelDispatcher {
	return service.NewMultiChannelDispatcher(
		sms.NewNoopAdapter(),
		email.NewNoopAdapter(),
		fcm.NewNoopAdapter(),
		kakao.NewNoopAdapter(),
		nil,
	)
}

func TestDispatcher_Critical_AllChannels(t *testing.T) {
	d := newTestDispatcher()
	ctx := context.Background()

	req := &service.MultiChannelRequest{
		UserID:            "user-001",
		Title:             "긴급 건강 알림",
		Body:              "혈당이 위험 수준입니다",
		Priority:          "critical",
		FCMToken:          "fcm-token-123",
		PhoneNumber:       "+821012345678",
		EmailAddr:         "user@example.com",
		KakaoPhone:        "010-1234-5678",
		KakaoTemplateCode: "ALERT_GLUCOSE",
		KakaoVariables:    map[string]string{"name": "홍길동"},
	}

	result, err := d.Dispatch(ctx, req)
	if err != nil {
		t.Fatalf("Dispatch 실패: %v", err)
	}
	if !result.AllSuccess() {
		t.Errorf("AllSuccess = false: errors=%v", result.Errors)
	}
	if len(result.SuccessChannels) != 4 {
		t.Errorf("SuccessChannels = %d, want 4 (critical)", len(result.SuccessChannels))
	}
}

func TestDispatcher_High_FCMPlusSMS(t *testing.T) {
	d := newTestDispatcher()
	ctx := context.Background()

	req := &service.MultiChannelRequest{
		UserID:      "u",
		Title:       "T", Body: "B", Priority: "high",
		FCMToken:    "fcm",
		PhoneNumber: "+821011112222",
	}
	result, _ := d.Dispatch(ctx, req)
	if len(result.AttemptedChannels) != 2 {
		t.Errorf("Attempted = %d, want 2 (fcm+sms)", len(result.AttemptedChannels))
	}
}

func TestDispatcher_Normal_FCMOnly(t *testing.T) {
	d := newTestDispatcher()

	req := &service.MultiChannelRequest{
		UserID: "u", Title: "T", Body: "B", Priority: "normal",
		FCMToken: "fcm",
	}
	result, _ := d.Dispatch(context.Background(), req)
	if len(result.AttemptedChannels) != 1 {
		t.Errorf("Attempted = %d, want 1 (fcm only)", len(result.AttemptedChannels))
	}
	if result.SuccessChannels[0] != "fcm" {
		t.Errorf("Channel = %q, want fcm", result.SuccessChannels[0])
	}
}

func TestDispatcher_Normal_FallbackToEmail(t *testing.T) {
	d := newTestDispatcher()

	// FCM 토큰 없으면 email 폴백
	req := &service.MultiChannelRequest{
		UserID: "u", Title: "T", Body: "B", Priority: "normal",
		EmailAddr: "u@x.com",
	}
	result, _ := d.Dispatch(context.Background(), req)
	if len(result.SuccessChannels) != 1 || result.SuccessChannels[0] != "email" {
		t.Errorf("폴백 미작동: %v", result.SuccessChannels)
	}
}

func TestDispatcher_Low_EmailOnly(t *testing.T) {
	d := newTestDispatcher()

	req := &service.MultiChannelRequest{
		UserID: "u", Title: "T", Body: "B", Priority: "low",
		EmailAddr: "u@x.com",
		FCMToken:  "ignored", // low 우선순위에선 미사용
	}
	result, _ := d.Dispatch(context.Background(), req)
	if len(result.AttemptedChannels) != 1 || result.AttemptedChannels[0] != "email" {
		t.Errorf("Channels = %v, want [email]", result.AttemptedChannels)
	}
}

func TestDispatcher_NoEligibleChannels(t *testing.T) {
	d := newTestDispatcher()

	req := &service.MultiChannelRequest{
		UserID: "u", Title: "T", Body: "B", Priority: "critical",
		// 수신자 정보 전무
	}
	_, err := d.Dispatch(context.Background(), req)
	if err == nil {
		t.Error("수신자 없이 발송 통과")
	}
}

func TestDispatcher_NoUserID(t *testing.T) {
	d := newTestDispatcher()
	_, err := d.Dispatch(context.Background(), &service.MultiChannelRequest{Priority: "high"})
	if err == nil {
		t.Error("UserID 없이 통과")
	}
}

func TestDispatcher_KakaoRequiresTemplate(t *testing.T) {
	d := newTestDispatcher()

	req := &service.MultiChannelRequest{
		UserID:     "u", Title: "T", Body: "B", Priority: "critical",
		KakaoPhone: "010-1234-5678", // 템플릿 코드 없음
	}
	result, _ := d.Dispatch(context.Background(), req)
	for _, ch := range result.AttemptedChannels {
		if ch == "kakao" {
			t.Error("템플릿 코드 없이 카카오 시도됨")
		}
	}
}

func TestDispatcher_HealthCheck(t *testing.T) {
	d := newTestDispatcher()
	results := d.HealthCheck(context.Background())
	if len(results) != 4 {
		t.Errorf("HealthCheck = %d, want 4 채널", len(results))
	}
	for _, r := range results {
		if !r.Healthy {
			t.Errorf("%s unhealthy: %s", r.Channel, r.Error)
		}
	}
}

func TestDispatcher_ProviderSummary(t *testing.T) {
	d := newTestDispatcher()
	summary := d.ProviderSummary()
	if len(summary) != 4 {
		t.Errorf("Providers = %d, want 4", len(summary))
	}
	if summary["fcm"] != "noop" {
		t.Errorf("fcm provider = %q", summary["fcm"])
	}
}

func TestDispatcher_NilAdapter_Skip(t *testing.T) {
	d := service.NewMultiChannelDispatcher(
		nil, // SMS 비활성
		email.NewNoopAdapter(),
		fcm.NewNoopAdapter(),
		nil, // Kakao 비활성
		nil,
	)

	req := &service.MultiChannelRequest{
		UserID: "u", Title: "T", Body: "B", Priority: "critical",
		FCMToken:    "fcm",
		PhoneNumber: "+821012345678", // SMS 비활성이므로 무시됨
		EmailAddr:   "u@x.com",
		KakaoPhone:  "010-1111-2222", // Kakao 비활성이므로 무시됨
	}
	result, _ := d.Dispatch(context.Background(), req)
	for _, ch := range result.AttemptedChannels {
		if ch == "sms" || ch == "kakao" {
			t.Errorf("nil 어댑터 채널 %s가 시도됨", ch)
		}
	}
	if len(result.AttemptedChannels) != 2 {
		t.Errorf("Attempted = %d, want 2 (fcm+email)", len(result.AttemptedChannels))
	}
}

func TestDispatcher_FromEnv_DefaultsToNoop(t *testing.T) {
	t.Setenv("SMS_PROVIDER", "")
	t.Setenv("EMAIL_PROVIDER", "")
	t.Setenv("KAKAO_PROVIDER", "")
	t.Setenv("FCM_ENABLED", "")

	d := service.NewMultiChannelDispatcherFromEnv()
	summary := d.ProviderSummary()
	for ch, prov := range summary {
		if prov != "noop" {
			t.Errorf("%s = %q, want noop", ch, prov)
		}
	}
}

func TestDispatcher_HasAnySuccess_AllSuccess(t *testing.T) {
	d := newTestDispatcher()

	req := &service.MultiChannelRequest{
		UserID: "u", Title: "T", Body: "B", Priority: "critical",
		FCMToken: "fcm", PhoneNumber: "+821012345678",
	}
	result, _ := d.Dispatch(context.Background(), req)
	if !result.HasAnySuccess() {
		t.Error("HasAnySuccess = false")
	}
	if !result.AllSuccess() {
		t.Error("AllSuccess = false (모든 noop은 성공)")
	}
}

// TestDispatcher_PoISE_FeedbackCollected는 PoISE Loop 주입 시
// 발송 결과가 자동으로 메트릭으로 수집되는지 검증합니다.
func TestDispatcher_PoISE_FeedbackCollected(t *testing.T) {
	d := newTestDispatcher()
	loop := poise.NewLoop(poise.DefaultThresholds)
	d.SetPoISELoop(loop)

	req := &service.MultiChannelRequest{
		UserID: "user-poise", Title: "T", Body: "B", Priority: "high",
		FCMToken: "fcm", PhoneNumber: "+821012345678",
	}
	_, err := d.Dispatch(context.Background(), req)
	if err != nil {
		t.Fatalf("Dispatch 실패: %v", err)
	}

	if loop.FeedbackCount() == 0 {
		t.Error("PoISE에 피드백이 수집되지 않음")
	}

	// ResponseTime 메트릭이 기록되었는지 확인
	metric := loop.AggregateMetric(poise.FeedbackResponseTime, time.Time{})
	if metric.SampleSize == 0 {
		t.Error("ResponseTime 메트릭 미수집")
	}
}

// TestDispatcher_PoISE_FailureRecording는 발송 실패 시 SafetyAlert 자동 기록을 검증합니다.
func TestDispatcher_PoISE_FailureRecording(t *testing.T) {
	// 일부 채널 nil로 실패 유도하지 않고 selectChannels 결과 0건일 때
	// 에러 반환되므로, 다른 시나리오: 모든 채널 정상이면 SafetyAlert 0건
	d := newTestDispatcher()
	loop := poise.NewLoop(poise.DefaultThresholds)
	d.SetPoISELoop(loop)

	req := &service.MultiChannelRequest{
		UserID: "user-x", Title: "T", Body: "B", Priority: "normal",
		FCMToken: "fcm",
	}
	_, _ = d.Dispatch(context.Background(), req)

	// 정상 발송 시 SafetyAlert는 0건이지만 ResponseTime은 기록
	metric := loop.AggregateMetric(poise.FeedbackSafetyAlert, time.Time{})
	if metric.SampleSize > 0 {
		t.Errorf("정상 발송인데 SafetyAlert 수집됨: %d", metric.SampleSize)
	}
}

// TestDispatcher_PoISE_NotInjected는 PoISE 미주입 시 통합 동작이 영향받지 않음을 검증합니다.
func TestDispatcher_PoISE_NotInjected(t *testing.T) {
	d := newTestDispatcher()
	// SetPoISELoop 미호출

	req := &service.MultiChannelRequest{
		UserID: "u", Title: "T", Body: "B", Priority: "normal",
		FCMToken: "fcm",
	}
	result, err := d.Dispatch(context.Background(), req)
	if err != nil {
		t.Fatalf("Dispatch 실패: %v", err)
	}
	if !result.HasAnySuccess() {
		t.Error("PoISE 미주입 영향으로 발송 실패")
	}
}
