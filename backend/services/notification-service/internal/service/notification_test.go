package service_test

import (
	"context"
	"testing"

	"github.com/manpasik/backend/services/notification-service/internal/repository/memory"
	"github.com/manpasik/backend/services/notification-service/internal/service"
	"go.uber.org/zap"
)

func setupTestService() *service.NotificationService {
	logger := zap.NewNop()
	notiRepo := memory.NewNotificationRepository()
	prefRepo := memory.NewPreferencesRepository()
	return service.NewNotificationService(logger, notiRepo, prefRepo)
}

func TestSendNotification_Success(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	noti, err := svc.SendNotification(ctx, "user-1", service.TypeHealthAlert, service.ChannelPush, service.PriorityHigh, "건강 경고", "혈당이 높습니다", "")
	if err != nil {
		t.Fatalf("알림 발송 실패: %v", err)
	}
	if noti.ID == "" {
		t.Fatal("알림 ID가 비어 있음")
	}
	if noti.UserID != "user-1" {
		t.Fatalf("UserID 불일치: got %s, want user-1", noti.UserID)
	}
	if noti.Type != service.TypeHealthAlert {
		t.Fatalf("Type 불일치: got %d, want %d", noti.Type, service.TypeHealthAlert)
	}
	if noti.Channel != service.ChannelPush {
		t.Fatalf("Channel 불일치: got %d, want %d", noti.Channel, service.ChannelPush)
	}
	if noti.IsRead {
		t.Fatal("새 알림은 읽지 않음 상태여야 함")
	}
}

func TestSendNotification_EmptyUserID(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	_, err := svc.SendNotification(ctx, "", service.TypeSystem, service.ChannelInApp, service.PriorityNormal, "제목", "내용", "")
	if err == nil {
		t.Fatal("빈 user_id에 에러가 반환되어야 함")
	}
}

func TestSendNotification_EmptyTitle(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	_, err := svc.SendNotification(ctx, "user-1", service.TypeSystem, service.ChannelInApp, service.PriorityNormal, "", "내용", "")
	if err == nil {
		t.Fatal("빈 title에 에러가 반환되어야 함")
	}
}

func TestSendNotification_AutoChannel(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	noti, err := svc.SendNotification(ctx, "user-2", service.TypeCommunity, service.ChannelUnknown, service.PriorityNormal, "코칭", "내용", "")
	if err != nil {
		t.Fatalf("자동 채널 선택 실패: %v", err)
	}
	if noti.Channel != service.ChannelInApp {
		t.Fatalf("기본 채널이 InApp이어야 함: got %d", noti.Channel)
	}
}

func TestListNotifications(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		_, err := svc.SendNotification(ctx, "user-list", service.TypeSystem, service.ChannelInApp, service.PriorityNormal, "제목", "내용", "")
		if err != nil {
			t.Fatalf("알림 생성 실패: %v", err)
		}
	}

	notis, total, _, err := svc.ListNotifications(ctx, "user-list", service.TypeUnknown, false, 10, 0)
	if err != nil {
		t.Fatalf("목록 조회 실패: %v", err)
	}
	if total != 3 {
		t.Fatalf("총 개수 불일치: got %d, want 3", total)
	}
	if len(notis) != 3 {
		t.Fatalf("반환 개수 불일치: got %d, want 3", len(notis))
	}
}

func TestListNotifications_TypeFilter(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	svc.SendNotification(ctx, "user-filter", service.TypeSystem, service.ChannelInApp, service.PriorityNormal, "시스템", "내용", "")
	svc.SendNotification(ctx, "user-filter", service.TypeCommunity, service.ChannelInApp, service.PriorityNormal, "코칭", "내용", "")
	svc.SendNotification(ctx, "user-filter", service.TypeSystem, service.ChannelInApp, service.PriorityNormal, "시스템2", "내용", "")

	notis, total, _, err := svc.ListNotifications(ctx, "user-filter", service.TypeSystem, false, 10, 0)
	if err != nil {
		t.Fatalf("필터 조회 실패: %v", err)
	}
	if total != 2 {
		t.Fatalf("시스템 알림 수 불일치: got %d, want 2", total)
	}
	if len(notis) != 2 {
		t.Fatalf("반환 개수 불일치: got %d, want 2", len(notis))
	}
}

func TestMarkAsRead(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	noti, _ := svc.SendNotification(ctx, "user-read", service.TypeSystem, service.ChannelInApp, service.PriorityNormal, "제목", "내용", "")

	err := svc.MarkAsRead(ctx, noti.ID)
	if err != nil {
		t.Fatalf("읽음 처리 실패: %v", err)
	}

	notis, _, _, _ := svc.ListNotifications(ctx, "user-read", service.TypeUnknown, true, 10, 0)
	if len(notis) != 0 {
		t.Fatalf("읽음 처리 후 미읽음 목록에 나오면 안 됨: got %d", len(notis))
	}
}

func TestMarkAllAsRead(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		svc.SendNotification(ctx, "user-allread", service.TypeSystem, service.ChannelInApp, service.PriorityNormal, "제목", "내용", "")
	}

	count, err := svc.MarkAllAsRead(ctx, "user-allread")
	if err != nil {
		t.Fatalf("전체 읽음 실패: %v", err)
	}
	if count != 5 {
		t.Fatalf("읽음 처리 수 불일치: got %d, want 5", count)
	}

	unread, _, _ := svc.GetUnreadCount(ctx, "user-allread")
	if unread != 0 {
		t.Fatalf("전체 읽음 후 미읽음 수 불일치: got %d, want 0", unread)
	}
}

func TestGetUnreadCount(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	svc.SendNotification(ctx, "user-unread", service.TypeHealthAlert, service.ChannelInApp, service.PriorityHigh, "경고1", "내용", "")
	svc.SendNotification(ctx, "user-unread", service.TypeHealthAlert, service.ChannelInApp, service.PriorityHigh, "경고2", "내용", "")
	svc.SendNotification(ctx, "user-unread", service.TypeCommunity, service.ChannelInApp, service.PriorityNormal, "코칭1", "내용", "")

	total, byType, err := svc.GetUnreadCount(ctx, "user-unread")
	if err != nil {
		t.Fatalf("미읽음 수 조회 실패: %v", err)
	}
	if total != 3 {
		t.Fatalf("총 미읽음 수 불일치: got %d, want 3", total)
	}
	if byType["health_alert"] != 2 {
		t.Fatalf("health_alert 미읽음 수 불일치: got %d, want 2", byType["health_alert"])
	}
	if byType["community"] != 1 {
		t.Fatalf("community 미읽음 수 불일치: got %d, want 1", byType["community"])
	}
}

func TestUpdateAndGetPreferences(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	pref := &service.NotificationPreferences{
		UserID:             "user-pref",
		PushEnabled:        true,
		EmailEnabled:       false,
		SMSEnabled:         false,
		InAppEnabled:       true,
		HealthAlertEnabled: true,
		CoachingEnabled:    true,
		PromotionEnabled:   false,
		QuietHoursStart:    "22:00",
		QuietHoursEnd:      "07:00",
		Language:           "ko",
	}

	result, err := svc.UpdatePreferences(ctx, pref)
	if err != nil {
		t.Fatalf("설정 업데이트 실패: %v", err)
	}
	if !result.PushEnabled {
		t.Fatal("PushEnabled가 true여야 함")
	}

	got, err := svc.GetPreferences(ctx, "user-pref")
	if err != nil {
		t.Fatalf("설정 조회 실패: %v", err)
	}
	if got.QuietHoursStart != "22:00" {
		t.Fatalf("QuietHoursStart 불일치: got %s, want 22:00", got.QuietHoursStart)
	}
}

func TestGetPreferences_Default(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	pref, err := svc.GetPreferences(ctx, "user-noprefs")
	if err != nil {
		t.Fatalf("기본 설정 조회 실패: %v", err)
	}
	if !pref.PushEnabled {
		t.Fatal("기본 PushEnabled가 true여야 함")
	}
	if pref.Language != "ko" {
		t.Fatalf("기본 Language 불일치: got %s, want ko", pref.Language)
	}
}

func TestSendNotification_WithData(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	dataStr := `{"deep_link":"/measurement/123","measurement_id":"m-456"}`

	noti, err := svc.SendNotification(ctx, "user-data", service.TypeMeasurement, service.ChannelPush, service.PriorityNormal, "측정 완료", "결과를 확인하세요", dataStr)
	if err != nil {
		t.Fatalf("알림 발송 실패: %v", err)
	}
	if noti.Data != dataStr {
		t.Fatalf("Data 불일치: got %s", noti.Data)
	}
}

func TestNotificationTypeToString(t *testing.T) {
	tests := []struct {
		t    service.NotificationType
		want string
	}{
		{service.TypeMeasurement, "measurement"},
		{service.TypeHealthAlert, "health_alert"},
		{service.TypeCommunity, "community"},
		{service.TypeAppointment, "appointment"},
		{service.TypePrescription, "prescription"},
		{service.TypeSystem, "system"},
		{service.TypePromotion, "promotion"},
		{service.TypeUnknown, "unknown"},
	}

	for _, tt := range tests {
		got := service.NotificationTypeToString(tt.t)
		if got != tt.want {
			t.Errorf("NotificationTypeToString(%d) = %s, want %s", tt.t, got, tt.want)
		}
	}
}

func TestSendFromTemplate_Success(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	// prescription_created 템플릿 사용
	err := svc.SendFromTemplate(ctx, "user-tmpl", "prescription_created", "김의사")
	if err != nil {
		t.Fatalf("SendFromTemplate 실패: %v", err)
	}

	// 알림이 생성되었는지 확인
	notis, total, _, err := svc.ListNotifications(ctx, "user-tmpl", service.TypeUnknown, false, 10, 0)
	if err != nil {
		t.Fatalf("알림 목록 조회 실패: %v", err)
	}
	if total != 1 {
		t.Fatalf("알림 수 불일치: got %d, want 1", total)
	}
	if notis[0].Title != "새 처방전 발행" {
		t.Fatalf("알림 제목 불일치: got %s, want '새 처방전 발행'", notis[0].Title)
	}
	if notis[0].Type != service.TypePrescription {
		t.Fatalf("알림 타입 불일치: got %d, want %d", notis[0].Type, service.TypePrescription)
	}
	if notis[0].Priority != service.PriorityHigh {
		t.Fatalf("알림 우선순위 불일치: got %d, want %d", notis[0].Priority, service.PriorityHigh)
	}
	if notis[0].Channel != service.ChannelPush {
		t.Fatalf("알림 채널 불일치: got %d, want %d", notis[0].Channel, service.ChannelPush)
	}
}

func TestSendFromTemplate_HealthAlert(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	err := svc.SendFromTemplate(ctx, "user-alert", "health_alert_critical", "혈당", "350 mg/dL")
	if err != nil {
		t.Fatalf("SendFromTemplate 실패: %v", err)
	}

	notis, _, _, _ := svc.ListNotifications(ctx, "user-alert", service.TypeUnknown, false, 10, 0)
	if len(notis) != 1 {
		t.Fatalf("알림 수 불일치: got %d, want 1", len(notis))
	}
	if notis[0].Priority != service.PriorityUrgent {
		t.Fatalf("urgent 알림 우선순위 불일치: got %d, want %d", notis[0].Priority, service.PriorityUrgent)
	}
	if notis[0].Type != service.TypeHealthAlert {
		t.Fatalf("알림 타입 불일치: got %d, want %d", notis[0].Type, service.TypeHealthAlert)
	}
}

func TestSendFromTemplate_InvalidTemplate(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	err := svc.SendFromTemplate(ctx, "user-1", "nonexistent_template")
	if err == nil {
		t.Fatal("존재하지 않는 템플릿에 에러가 반환되어야 함")
	}
}

func TestPredefinedTemplatesExist(t *testing.T) {
	expectedKeys := []string{
		"prescription_created", "prescription_sent", "prescription_ready", "prescription_dispensed",
		"delivery_started", "delivery_arrived",
		"appointment_reminder", "appointment_cancelled",
		"health_alert_critical", "health_alert_warning",
		"measurement_complete",
		"family_data_shared",
	}

	for _, key := range expectedKeys {
		tmpl, exists := service.PredefinedTemplates[key]
		if !exists {
			t.Errorf("템플릿 '%s'가 PredefinedTemplates에 없습니다", key)
			continue
		}
		if tmpl.Key != key {
			t.Errorf("템플릿 '%s'의 Key 불일치: got %s", key, tmpl.Key)
		}
		if tmpl.Title == "" {
			t.Errorf("템플릿 '%s'의 Title이 비어 있습니다", key)
		}
		if tmpl.BodyFmt == "" {
			t.Errorf("템플릿 '%s'의 BodyFmt가 비어 있습니다", key)
		}
		if tmpl.Type == "" {
			t.Errorf("템플릿 '%s'의 Type이 비어 있습니다", key)
		}
		if tmpl.Priority == "" {
			t.Errorf("템플릿 '%s'의 Priority가 비어 있습니다", key)
		}
		if tmpl.Channel == "" {
			t.Errorf("템플릿 '%s'의 Channel이 비어 있습니다", key)
		}
	}

	// 총 템플릿 수 확인
	if len(service.PredefinedTemplates) != len(expectedKeys) {
		t.Errorf("PredefinedTemplates 수: got %d, want %d", len(service.PredefinedTemplates), len(expectedKeys))
	}
}

func TestEndToEnd_NotificationFlow(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	pref := &service.NotificationPreferences{
		UserID:             "user-e2e",
		PushEnabled:        true,
		EmailEnabled:       true,
		SMSEnabled:         false,
		InAppEnabled:       true,
		HealthAlertEnabled: true,
		CoachingEnabled:    true,
		PromotionEnabled:   false,
		Language:           "en",
	}
	svc.UpdatePreferences(ctx, pref)

	svc.SendNotification(ctx, "user-e2e", service.TypeHealthAlert, service.ChannelPush, service.PriorityUrgent, "Critical Alert", "Check now", "")
	svc.SendNotification(ctx, "user-e2e", service.TypeCommunity, service.ChannelInApp, service.PriorityNormal, "Daily Tip", "Drink water", "")
	noti3, _ := svc.SendNotification(ctx, "user-e2e", service.TypePromotion, service.ChannelEmail, service.PriorityLow, "Order Shipped", "Your order is on the way", "")

	unread, _, _ := svc.GetUnreadCount(ctx, "user-e2e")
	if unread != 3 {
		t.Fatalf("미읽음 수 불일치: got %d, want 3", unread)
	}

	svc.MarkAsRead(ctx, noti3.ID)

	unread2, _, _ := svc.GetUnreadCount(ctx, "user-e2e")
	if unread2 != 2 {
		t.Fatalf("읽음 후 미읽음 수 불일치: got %d, want 2", unread2)
	}

	notis, total, _, _ := svc.ListNotifications(ctx, "user-e2e", service.TypeUnknown, false, 10, 0)
	if total != 3 {
		t.Fatalf("총 알림 수 불일치: got %d, want 3", total)
	}
	if len(notis) != 3 {
		t.Fatalf("반환 수 불일치: got %d, want 3", len(notis))
	}

	count, _ := svc.MarkAllAsRead(ctx, "user-e2e")
	if count != 2 {
		t.Fatalf("전체 읽음 처리 수 불일치: got %d, want 2 (이미 1개 읽음)", count)
	}
}
