// Package integration_test contains end-to-end integration tests for cross-service flows.
//
// 이 패키지는 외부 어댑터(SMS/Email/FCM/Kakao)와 도메인 서비스의 통합을 검증합니다.
package integration_test

import (
	"context"
	"strings"
	"testing"

	"github.com/manpasik/backend/shared/external/email"
	"github.com/manpasik/backend/shared/external/fcm"
	"github.com/manpasik/backend/shared/external/kakao"
	"github.com/manpasik/backend/shared/external/sms"
)

// TestNotificationFlow_HealthAlert_AllChannels는 건강 이상 발생 시 다중 채널 발송을 검증합니다.
func TestNotificationFlow_HealthAlert_AllChannels(t *testing.T) {
	ctx := context.Background()
	smsAdapter := sms.NewNoopAdapter()
	emailAdapter := email.NewNoopAdapter()
	fcmAdapter := fcm.NewNoopAdapter()
	kakaoAdapter := kakao.NewNoopAdapter()

	userID := "user-001"
	phone := "+821012345678"
	emailAddr := "user@example.com"
	deviceToken := "fcm-device-001"

	// 시나리오: 혈당 이상 감지 → 4채널 발송
	smsMsg := &sms.Message{To: phone, Body: "혈당이 정상 범위를 초과했습니다."}
	emailMsg := &email.Message{
		From:     &email.Address{Email: "alert@manpasik.com", Name: "만파식 알림"},
		To:       []*email.Address{{Email: emailAddr}},
		Subject:  "건강 알림: 혈당 이상",
		HTMLBody: "<p>혈당 수치가 높습니다.</p>",
	}
	fcmMsg := &fcm.Message{
		Token:        deviceToken,
		Notification: &fcm.Notification{Title: "건강 알림", Body: "혈당 이상"},
		Data:         map[string]string{"user_id": userID, "type": "glucose_alert"},
	}
	kakaoMsg := &kakao.Message{
		TemplateCode: "ALERT_GLUCOSE",
		To:           "010-1234-5678",
		Variables:    map[string]string{"name": "홍길동"},
	}

	// 모든 채널로 발송
	if _, err := smsAdapter.Send(ctx, smsMsg); err != nil {
		t.Errorf("SMS 발송 실패: %v", err)
	}
	if _, err := emailAdapter.Send(ctx, emailMsg); err != nil {
		t.Errorf("Email 발송 실패: %v", err)
	}
	if _, err := fcmAdapter.Send(ctx, fcmMsg); err != nil {
		t.Errorf("FCM 발송 실패: %v", err)
	}
	if _, err := kakaoAdapter.Send(ctx, kakaoMsg); err != nil {
		t.Errorf("Kakao 발송 실패: %v", err)
	}

	if smsAdapter.Count() != 1 {
		t.Errorf("SMS Count = %d, want 1", smsAdapter.Count())
	}
	if emailAdapter.Count() != 1 {
		t.Errorf("Email Count = %d, want 1", emailAdapter.Count())
	}
	if fcmAdapter.Count() != 1 {
		t.Errorf("FCM Count = %d, want 1", fcmAdapter.Count())
	}
	if kakaoAdapter.Count() != 1 {
		t.Errorf("Kakao Count = %d, want 1", kakaoAdapter.Count())
	}
}

// TestNotificationFlow_FamilyBroadcast는 가족 그룹 다중 발송을 검증합니다.
func TestNotificationFlow_FamilyBroadcast(t *testing.T) {
	ctx := context.Background()
	fcmAdapter := fcm.NewNoopAdapter()

	familyTokens := []string{"fcm-mom", "fcm-dad", "fcm-son", "fcm-daughter"}
	msg := &fcm.Message{
		Notification: &fcm.Notification{
			Title: "긴급 알림",
			Body:  "어머니의 혈당이 위험 수준에 도달했습니다.",
		},
		Data: map[string]string{"alert_id": "alert-001", "severity": "critical"},
	}

	result, err := fcmAdapter.SendMulticast(ctx, msg, familyTokens)
	if err != nil {
		t.Fatalf("Multicast 실패: %v", err)
	}
	if result.SuccessCount != 4 {
		t.Errorf("SuccessCount = %d, want 4", result.SuccessCount)
	}
}

// TestNotificationFlow_TemplateRendering은 카카오 알림톡 템플릿 발송을 검증합니다.
func TestNotificationFlow_TemplateRendering(t *testing.T) {
	ctx := context.Background()
	a := kakao.NewNoopAdapter()

	body := "안녕하세요 #{name}님, 주문 #{order_id}이 #{status}되었습니다."
	if err := a.RegisterTemplate("ORDER_STATUS", body, nil); err != nil {
		t.Fatalf("RegisterTemplate 실패: %v", err)
	}

	tmpl, _ := a.GetTemplate("ORDER_STATUS")
	rendered := kakao.RenderTemplate(tmpl.Body, map[string]string{
		"name":     "홍길동",
		"order_id": "ORD-12345",
		"status":   "배송 시작",
	})

	if !strings.Contains(rendered, "홍길동") || !strings.Contains(rendered, "배송 시작") {
		t.Errorf("렌더링 실패: %q", rendered)
	}

	msg := &kakao.Message{
		TemplateCode: "ORDER_STATUS",
		To:           "010-1234-5678",
		Variables: map[string]string{
			"name":     "홍길동",
			"order_id": "ORD-12345",
			"status":   "배송 시작",
		},
	}
	if _, err := a.Send(ctx, msg); err != nil {
		t.Errorf("Send 실패: %v", err)
	}
}

// TestNotificationFlow_FailoverToNoop은 환경변수 미설정 시 Noop 폴백을 검증합니다.
func TestNotificationFlow_FailoverToNoop(t *testing.T) {
	t.Setenv("SMS_PROVIDER", "")
	t.Setenv("EMAIL_PROVIDER", "")
	t.Setenv("KAKAO_PROVIDER", "")

	if sms.NewFromEnv().Provider() != "noop" {
		t.Error("SMS noop 폴백 실패")
	}
	if email.NewFromEnv().Provider() != "noop" {
		t.Error("Email noop 폴백 실패")
	}
	if kakao.NewFromEnv().Provider() != "noop" {
		t.Error("Kakao noop 폴백 실패")
	}
}

// TestNotificationFlow_EmailTemplateChain은 이메일 템플릿 → 발송 흐름을 검증합니다.
func TestNotificationFlow_EmailTemplateChain(t *testing.T) {
	ctx := context.Background()
	a := email.NewNoopAdapter()

	html, err := email.RenderTemplate("welcome", map[string]string{"name": "홍길동"})
	if err != nil {
		t.Fatalf("Render 실패: %v", err)
	}

	msg := &email.Message{
		From:     &email.Address{Email: "noreply@manpasik.com", Name: "만파식"},
		To:       []*email.Address{{Email: "user@example.com"}},
		Subject:  "환영합니다",
		HTMLBody: html,
	}
	result, err := a.Send(ctx, msg)
	if err != nil {
		t.Fatalf("Send 실패: %v", err)
	}
	if !strings.Contains(result.Provider, "noop") {
		t.Errorf("Provider = %q", result.Provider)
	}
}

// TestNotificationFlow_BulkEmail는 다수 이메일 발송을 검증합니다.
func TestNotificationFlow_BulkEmail(t *testing.T) {
	ctx := context.Background()
	a := email.NewNoopAdapter()

	emails := []string{"u1@x.com", "u2@x.com", "u3@x.com", "u4@x.com", "u5@x.com"}

	for _, e := range emails {
		msg := &email.Message{
			From:     &email.Address{Email: "f@x.com"},
			To:       []*email.Address{{Email: e}},
			Subject:  "테스트",
			TextBody: "Bulk send test",
		}
		if _, err := a.Send(ctx, msg); err != nil {
			t.Errorf("%q 발송 실패: %v", e, err)
		}
	}

	if a.Count() != 5 {
		t.Errorf("Count = %d, want 5", a.Count())
	}
}

// TestNotificationFlow_ChannelPriority는 채널 우선순위 (FCM → SMS → Email) 폴백을 검증합니다.
func TestNotificationFlow_ChannelPriority(t *testing.T) {
	ctx := context.Background()
	fcmA := fcm.NewNoopAdapter()
	smsA := sms.NewNoopAdapter()
	emailA := email.NewNoopAdapter()

	// 시나리오: 사용자가 FCM 토큰이 있으면 FCM 발송, 없으면 SMS, 그 다음 Email
	type recipient struct {
		userID string
		fcmToken string
		phone string
		email string
	}

	recipients := []recipient{
		{userID: "u1", fcmToken: "f1"},                          // FCM
		{userID: "u2", phone: "+821011112222"},                   // SMS
		{userID: "u3", email: "u3@x.com"},                        // Email
	}

	for _, r := range recipients {
		switch {
		case r.fcmToken != "":
			_, _ = fcmA.Send(ctx, &fcm.Message{
				Token:        r.fcmToken,
				Notification: &fcm.Notification{Title: "T", Body: "B"},
			})
		case r.phone != "":
			_, _ = smsA.Send(ctx, &sms.Message{To: r.phone, Body: "T"})
		case r.email != "":
			_, _ = emailA.Send(ctx, &email.Message{
				From:     &email.Address{Email: "f@x.com"},
				To:       []*email.Address{{Email: r.email}},
				Subject:  "T",
				TextBody: "T",
			})
		}
	}

	if fcmA.Count() != 1 || smsA.Count() != 1 || emailA.Count() != 1 {
		t.Errorf("Counts = fcm:%d sms:%d email:%d", fcmA.Count(), smsA.Count(), emailA.Count())
	}
}
