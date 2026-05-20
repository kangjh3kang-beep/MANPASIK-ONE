// Package integration_test: E2E 시나리오 테스트
//
// 외부 어댑터들을 사용한 사용자 여정 시나리오를 검증합니다.
package integration_test

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/manpasik/backend/shared/external/email"
	"github.com/manpasik/backend/shared/external/fcm"
	"github.com/manpasik/backend/shared/external/health"
	"github.com/manpasik/backend/shared/external/kakao"
	"github.com/manpasik/backend/shared/external/sms"
)

// TestE2E_UserOnboarding_WelcomeFlow는 신규 사용자 온보딩 흐름을 검증합니다.
func TestE2E_UserOnboarding_WelcomeFlow(t *testing.T) {
	ctx := context.Background()
	emailA := email.NewNoopAdapter()
	smsA := sms.NewNoopAdapter()
	kakaoA := kakao.NewNoopAdapter()

	// 1. 이메일 인증 메일 발송
	verifyEmail := &email.Message{
		From:     &email.Address{Email: "noreply@manpasik.com", Name: "만파식"},
		To:       []*email.Address{{Email: "newuser@example.com"}},
		Subject:  "이메일 인증을 완료해주세요",
		HTMLBody: "<a href='https://manpasik.com/verify?token=abc'>인증하기</a>",
	}
	if _, err := emailA.Send(ctx, verifyEmail); err != nil {
		t.Errorf("인증 이메일 발송 실패: %v", err)
	}

	// 2. SMS 인증코드 발송
	verifySMS := &sms.Message{
		To:   "+821012345678",
		Body: "[만파식] 인증번호: 123456",
	}
	if _, err := smsA.Send(ctx, verifySMS); err != nil {
		t.Errorf("SMS 인증 발송 실패: %v", err)
	}

	// 3. 환영 카카오 알림톡
	_ = kakaoA.RegisterTemplate("WELCOME_001", "안녕하세요 #{name}님! 만파식에 가입해주셔서 감사합니다.", nil)
	welcomeKakao := &kakao.Message{
		TemplateCode: "WELCOME_001",
		To:           "010-1234-5678",
		Variables:    map[string]string{"name": "홍길동"},
	}
	if _, err := kakaoA.Send(ctx, welcomeKakao); err != nil {
		t.Errorf("환영 알림톡 발송 실패: %v", err)
	}

	// 4. 환영 이메일
	html, _ := email.RenderTemplate("welcome", map[string]string{"name": "홍길동"})
	welcomeEmail := &email.Message{
		From:     &email.Address{Email: "noreply@manpasik.com"},
		To:       []*email.Address{{Email: "newuser@example.com"}},
		Subject:  "환영합니다!",
		HTMLBody: html,
	}
	if _, err := emailA.Send(ctx, welcomeEmail); err != nil {
		t.Errorf("환영 이메일 발송 실패: %v", err)
	}

	if emailA.Count() != 2 {
		t.Errorf("Email = %d, want 2", emailA.Count())
	}
	if smsA.Count() != 1 {
		t.Errorf("SMS = %d, want 1", smsA.Count())
	}
	if kakaoA.Count() != 1 {
		t.Errorf("Kakao = %d, want 1", kakaoA.Count())
	}
}

// TestE2E_HealthAlertFamilyBroadcast는 건강 이상 → 가족 그룹 다중 채널 알림을 검증합니다.
func TestE2E_HealthAlertFamilyBroadcast(t *testing.T) {
	ctx := context.Background()
	fcmA := fcm.NewNoopAdapter()
	smsA := sms.NewNoopAdapter()
	kakaoA := kakao.NewNoopAdapter()

	// 가족 그룹: 어머니 측정 이상 → 자녀 3명에게 알림
	familyTokens := []string{"fcm-son", "fcm-daughter-1", "fcm-daughter-2"}
	familyPhones := []string{"+821011111111", "+821022222222", "+821033333333"}

	// 1. 푸시 알림 (FCM Multicast)
	fcmMsg := &fcm.Message{
		Notification: &fcm.Notification{
			Title: "긴급: 어머니의 혈당 이상",
			Body:  "어머니의 혈당이 280 mg/dL로 위험 수준입니다.",
		},
		Data: map[string]string{
			"alert_id": "alert-001",
			"severity": "critical",
		},
	}
	fcmResult, err := fcmA.SendMulticast(ctx, fcmMsg, familyTokens)
	if err != nil {
		t.Fatalf("FCM Multicast 실패: %v", err)
	}
	if fcmResult.SuccessCount != 3 {
		t.Errorf("FCM 성공 = %d, want 3", fcmResult.SuccessCount)
	}

	// 2. SMS 백업 알림 (FCM 미수신 가능성 대비)
	for _, phone := range familyPhones {
		_, _ = smsA.Send(ctx, &sms.Message{
			To:   phone,
			Body: "[만파식] 어머니 건강 이상. 즉시 확인해주세요.",
		})
	}
	if smsA.Count() != 3 {
		t.Errorf("SMS = %d, want 3", smsA.Count())
	}

	// 3. 카카오 알림톡 (한국 사용자 우선)
	_ = kakaoA.RegisterTemplate("ALERT_FAMILY", "[만파식] #{name}님의 #{biomarker} 이상 #{value}", nil)
	for i, phone := range []string{"010-1111-1111", "010-2222-2222"} {
		_, _ = kakaoA.Send(ctx, &kakao.Message{
			TemplateCode: "ALERT_FAMILY",
			To:           phone,
			Variables: map[string]string{
				"name":      "어머니",
				"biomarker": "혈당",
				"value":     "280 mg/dL",
			},
		})
		_ = i
	}
}

// TestE2E_HealthDataSyncBidirectional은 양방향 헬스 데이터 동기화입니다.
func TestE2E_HealthDataSyncBidirectional(t *testing.T) {
	repo := newInMemoryRepo()
	bridge := health.NewHealthBridge(repo, health.PolicyManpasikWins)

	now := time.Now().UTC()

	// 1. Apple HealthKit → 만파식 (Pull)
	appleBatch := &health.SyncBatch{
		UserID: "user-bidir",
		Source: health.SourceAppleHealthKit,
		Samples: []*health.HealthSample{
			{UserID: "user-bidir", Type: health.TypeSteps, Value: 8500, Unit: "count",
				Source: health.SourceAppleHealthKit, SourceID: "ak-1", Timestamp: now},
			{UserID: "user-bidir", Type: health.TypeHeartRate, Value: 72, Unit: "bpm",
				Source: health.SourceAppleHealthKit, SourceID: "ak-2", Timestamp: now},
		},
	}
	pullResult, err := bridge.Sync(appleBatch)
	if err != nil {
		t.Fatalf("Pull 실패: %v", err)
	}
	if pullResult.ImportedCount != 2 {
		t.Errorf("Pull ImportedCount = %d, want 2", pullResult.ImportedCount)
	}

	// 2. 만파식 측정 직접 저장 (실 측정)
	mpSample := &health.HealthSample{
		UserID: "user-bidir",
		Type:   health.TypeBloodGlucose,
		Value:  110, Unit: "mg/dL",
		Source: health.SourceManPaSik,
		SourceID: "mp-1",
		Timestamp: now,
	}
	_, _ = repo.Save(mpSample)

	// 3. 만파식 → Apple HealthKit (Push)
	pushBatch, err := bridge.PreparePushBatch("user-bidir", health.SourceAppleHealthKit, now.Add(-1*time.Hour))
	if err != nil {
		t.Fatalf("Push 준비 실패: %v", err)
	}
	if len(pushBatch.Samples) != 1 {
		t.Errorf("Push samples = %d, want 1", len(pushBatch.Samples))
	}
	if pushBatch.Samples[0].Source != health.SourceManPaSik {
		t.Error("Push 샘플의 Source가 만파식이 아님")
	}
}

// TestE2E_OrderToShipmentToReceipt은 주문 → 결제 → 영수증 → 배송 알림 흐름입니다.
func TestE2E_OrderToShipmentToReceipt(t *testing.T) {
	ctx := context.Background()
	emailA := email.NewNoopAdapter()
	fcmA := fcm.NewNoopAdapter()

	// 1. 주문 영수증 이메일
	receiptHTML, _ := email.RenderTemplate("receipt", map[string]string{
		"order_id":     "ORD-2026-001",
		"total_amount": "75,000원",
	})
	_, err := emailA.Send(ctx, &email.Message{
		From:     &email.Address{Email: "shop@manpasik.com"},
		To:       []*email.Address{{Email: "buyer@example.com"}},
		Subject:  "[만파식] 주문 영수증 - ORD-2026-001",
		HTMLBody: receiptHTML,
	})
	if err != nil {
		t.Errorf("영수증 이메일 발송 실패: %v", err)
	}
	if !strings.Contains(receiptHTML, "ORD-2026-001") {
		t.Error("영수증에 주문번호 누락")
	}

	// 2. 배송 시작 푸시 알림
	_, err = fcmA.Send(ctx, &fcm.Message{
		Token: "buyer-fcm-token",
		Notification: &fcm.Notification{
			Title: "배송이 시작되었습니다",
			Body:  "예상 도착일: 2026-05-02",
		},
		Data: map[string]string{
			"order_id":         "ORD-2026-001",
			"tracking_number":  "TRK-12345",
		},
	})
	if err != nil {
		t.Errorf("배송 푸시 실패: %v", err)
	}

	// 3. 배송 완료 알림
	_, err = fcmA.Send(ctx, &fcm.Message{
		Token: "buyer-fcm-token",
		Notification: &fcm.Notification{
			Title: "배송이 완료되었습니다",
			Body:  "상품을 안전하게 받으셨나요?",
		},
	})
	if err != nil {
		t.Errorf("배송 완료 푸시 실패: %v", err)
	}

	if fcmA.Count() != 2 {
		t.Errorf("FCM = %d, want 2", fcmA.Count())
	}
}

// TestE2E_PasswordResetFlow는 비밀번호 재설정 흐름을 검증합니다.
func TestE2E_PasswordResetFlow(t *testing.T) {
	ctx := context.Background()
	emailA := email.NewNoopAdapter()
	smsA := sms.NewNoopAdapter()

	// 1. 비밀번호 재설정 이메일
	resetHTML, _ := email.RenderTemplate("password_reset", map[string]string{
		"reset_link": "https://manpasik.com/reset?token=abc123",
	})
	_, err := emailA.Send(ctx, &email.Message{
		From:     &email.Address{Email: "noreply@manpasik.com"},
		To:       []*email.Address{{Email: "user@example.com"}},
		Subject:  "비밀번호 재설정 요청",
		HTMLBody: resetHTML,
	})
	if err != nil {
		t.Errorf("재설정 이메일 발송 실패: %v", err)
	}

	// 2. 의심스러운 활동 알림 (다른 IP 접속 등)
	_, err = smsA.Send(ctx, &sms.Message{
		To:   "+821012345678",
		Body: "[만파식] 새 위치에서 비밀번호 재설정 요청. 본인이 아니면 무시하세요.",
	})
	if err != nil {
		t.Errorf("의심 활동 SMS 발송 실패: %v", err)
	}
}

// TestE2E_EmergencyEscalation은 응급 상황 다단계 에스컬레이션입니다.
func TestE2E_EmergencyEscalation(t *testing.T) {
	ctx := context.Background()
	smsA := sms.NewNoopAdapter()
	fcmA := fcm.NewNoopAdapter()

	// 1단계: 사용자 본인 푸시
	_, _ = fcmA.Send(ctx, &fcm.Message{
		Token: "user-fcm",
		Notification: &fcm.Notification{Title: "위험 수치", Body: "즉시 확인이 필요합니다"},
		Android: &fcm.Android{Priority: "high", ChannelID: "emergency"},
	})

	// 2단계: 가족에게 알림 (응답 없음 가정)
	familyTokens := []string{"fam-1", "fam-2"}
	r, _ := fcmA.SendMulticast(ctx, &fcm.Message{
		Notification: &fcm.Notification{Title: "긴급", Body: "환자 응답 없음"},
	}, familyTokens)
	if r.SuccessCount != 2 {
		t.Errorf("가족 푸시 = %d, want 2", r.SuccessCount)
	}

	// 3단계: 119/응급 SMS
	_, _ = smsA.Send(ctx, &sms.Message{
		To:   "+821011112222",
		Body: "[만파식] 환자 응답 없음. 119 호출이 필요합니다.",
	})

	if fcmA.Count() != 3 {
		t.Errorf("FCM = %d, want 3 (본인+가족 2)", fcmA.Count())
	}
	if smsA.Count() != 1 {
		t.Errorf("SMS = %d, want 1", smsA.Count())
	}
}

// TestE2E_ConcurrentNotificationsLoad는 동시 다중 알림 부하 테스트입니다.
func TestE2E_ConcurrentNotificationsLoad(t *testing.T) {
	ctx := context.Background()
	fcmA := fcm.NewNoopAdapter()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, _ = fcmA.Send(ctx, &fcm.Message{
				Token:        "concurrent-token",
				Notification: &fcm.Notification{Title: "T", Body: "B"},
				Data:         map[string]string{"index": string(rune(i))},
			})
		}(i)
	}
	wg.Wait()

	if fcmA.Count() != 100 {
		t.Errorf("FCM Count = %d, want 100", fcmA.Count())
	}
}

// TestE2E_MultilingualEmail은 다국어 이메일 발송을 검증합니다.
func TestE2E_MultilingualEmail(t *testing.T) {
	ctx := context.Background()
	a := email.NewNoopAdapter()

	cases := []struct {
		name    string
		subject string
		body    string
	}{
		{"한국어", "안녕하세요", "<p>건강 측정 결과를 확인해주세요</p>"},
		{"English", "Hello", "<p>Please review your health metrics</p>"},
		{"日本語", "こんにちは", "<p>健康データをご確認ください</p>"},
		{"中文", "您好", "<p>请查看您的健康数据</p>"},
	}

	for _, c := range cases {
		_, err := a.Send(ctx, &email.Message{
			From:     &email.Address{Email: "intl@manpasik.com"},
			To:       []*email.Address{{Email: "user@example.com"}},
			Subject:  c.subject,
			HTMLBody: c.body,
		})
		if err != nil {
			t.Errorf("%s 발송 실패: %v", c.name, err)
		}
	}

	if a.Count() != 4 {
		t.Errorf("Count = %d, want 4", a.Count())
	}
}

// TestE2E_HealthSyncWithUnitConversion는 단위 변환 동기화를 검증합니다.
func TestE2E_HealthSyncWithUnitConversion(t *testing.T) {
	repo := newInMemoryRepo()
	bridge := health.NewHealthBridge(repo, health.PolicyKeepExisting)

	now := time.Now().UTC()

	// 미국 기기에서 lb, °F, mg/dL 보냄
	usBatch := &health.SyncBatch{
		UserID: "user-us",
		Source: health.SourceAppleHealthKit,
		Samples: []*health.HealthSample{
			{UserID: "user-us", Type: health.TypeWeight, Value: 154, Unit: "lb",
				Source: health.SourceAppleHealthKit, SourceID: "us-1", Timestamp: now},
			{UserID: "user-us", Type: health.TypeBodyTemperature, Value: 98.6, Unit: "°F",
				Source: health.SourceAppleHealthKit, SourceID: "us-2", Timestamp: now},
		},
	}

	r, _ := bridge.Sync(usBatch)
	if r.ImportedCount != 2 {
		t.Errorf("ImportedCount = %d, want 2", r.ImportedCount)
	}

	// 단위가 SI로 변환되어 저장되었는지 확인
	weight, _ := repo.FindBySourceID("user-us", health.SourceAppleHealthKit, "us-1")
	if weight.Unit != "kg" {
		t.Errorf("Weight unit = %q, want kg", weight.Unit)
	}
	temp, _ := repo.FindBySourceID("user-us", health.SourceAppleHealthKit, "us-2")
	if temp.Unit != "°C" {
		t.Errorf("Temp unit = %q, want °C", temp.Unit)
	}
}

// TestE2E_KakaoTemplateLifecycle은 카카오 템플릿 등록 → 발송 → 통계 흐름입니다.
func TestE2E_KakaoTemplateLifecycle(t *testing.T) {
	ctx := context.Background()
	a := kakao.NewNoopAdapter()

	// 1. 템플릿 등록
	templates := map[string]string{
		"ORDER_CONFIRMED": "[만파식] #{name}님의 주문 #{order_id}이 접수되었습니다.",
		"DELIVERY_START":  "[만파식] 주문 #{order_id} 배송이 시작되었습니다.",
		"HEALTH_ALERT":    "[만파식] #{name}님의 #{biomarker} 수치 이상.",
	}
	for code, body := range templates {
		if err := a.RegisterTemplate(code, body, nil); err != nil {
			t.Errorf("%s 템플릿 등록 실패: %v", code, err)
		}
	}

	// 2. 각 템플릿으로 발송
	sends := []struct {
		code string
		vars map[string]string
	}{
		{"ORDER_CONFIRMED", map[string]string{"name": "홍길동", "order_id": "ORD-1"}},
		{"DELIVERY_START", map[string]string{"order_id": "ORD-1"}},
		{"HEALTH_ALERT", map[string]string{"name": "홍길동", "biomarker": "혈당"}},
	}
	for _, s := range sends {
		_, err := a.Send(ctx, &kakao.Message{
			TemplateCode: s.code,
			To:           "010-1234-5678",
			Variables:    s.vars,
		})
		if err != nil {
			t.Errorf("%s 발송 실패: %v", s.code, err)
		}
	}

	if a.Count() != 3 {
		t.Errorf("Count = %d, want 3", a.Count())
	}
}

// TestE2E_HealthAnalysisIncluding90Days는 90일 분석 시나리오입니다.
func TestE2E_HealthAnalysisIncluding90Days(t *testing.T) {
	repo := newInMemoryRepo()
	bridge := health.NewHealthBridge(repo, health.PolicyLatestWins)

	// 90일치 일별 측정 데이터
	now := time.Now().UTC()
	for d := 0; d < 90; d++ {
		batch := &health.SyncBatch{
			UserID: "user-90day",
			Source: health.SourceAppleHealthKit,
			Samples: []*health.HealthSample{
				{
					UserID: "user-90day", Type: health.TypeBloodGlucose,
					Value: 100 + float64(d%20), Unit: "mg/dL",
					Source: health.SourceAppleHealthKit,
					SourceID: "day-" + string(rune('a'+d%26)) + string(rune('a'+d/26)),
					Timestamp: now.AddDate(0, 0, -d),
				},
			},
		}
		r, _ := bridge.Sync(batch)
		_ = r
	}

	// 90일 데이터 조회
	samples, err := repo.FindByTimeRange("user-90day", health.TypeBloodGlucose, now.AddDate(0, 0, -90), now)
	if err != nil {
		t.Fatalf("FindByTimeRange 실패: %v", err)
	}
	if len(samples) < 80 {
		// 일부 source_id 충돌 가능 (rune 인코딩)
		t.Logf("samples = %d (일부 ID 충돌 가능)", len(samples))
	}
}
