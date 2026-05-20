package service_test

import (
	"context"
	"testing"
	"time"

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

// =============================================================================
// Phase F 테스트 보강: 엣지 케이스
// =============================================================================

func TestMarkAsRead_EmptyID(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()
	err := svc.MarkAsRead(ctx, "")
	if err == nil {
		t.Fatal("빈 ID에 에러가 반환되어야 함")
	}
}

func TestMarkAllAsRead_EmptyUserID(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()
	_, err := svc.MarkAllAsRead(ctx, "")
	if err == nil {
		t.Fatal("빈 user_id에 에러가 반환되어야 함")
	}
}

func TestGetUnreadCount_EmptyUserID(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()
	_, _, err := svc.GetUnreadCount(ctx, "")
	if err == nil {
		t.Fatal("빈 user_id에 에러가 반환되어야 함")
	}
}

func TestListNotifications_EmptyUserID(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()
	_, _, _, err := svc.ListNotifications(ctx, "", service.TypeUnknown, false, 10, 0)
	if err == nil {
		t.Fatal("빈 user_id에 에러가 반환되어야 함")
	}
}

func TestUpdatePreferences_NilPref(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()
	_, err := svc.UpdatePreferences(ctx, nil)
	if err == nil {
		t.Fatal("nil 설정에 에러가 반환되어야 함")
	}
}

func TestUpdatePreferences_EmptyUserID(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()
	pref := &service.NotificationPreferences{UserID: ""}
	_, err := svc.UpdatePreferences(ctx, pref)
	if err == nil {
		t.Fatal("빈 user_id에 에러가 반환되어야 함")
	}
}

func TestSendNotification_HighPriorityAutoChannel(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	// 사용자 설정: push 활성화
	pref := &service.NotificationPreferences{
		UserID:       "user-hpac",
		PushEnabled:  true,
		InAppEnabled: true,
		Language:     "ko",
	}
	svc.UpdatePreferences(ctx, pref)

	// ChannelUnknown + PriorityHigh → push가 선택되어야 함
	noti, err := svc.SendNotification(ctx, "user-hpac", service.TypeHealthAlert, service.ChannelUnknown, service.PriorityHigh, "긴급", "내용", "")
	if err != nil {
		t.Fatalf("발송 실패: %v", err)
	}
	if noti.Channel != service.ChannelPush {
		t.Fatalf("High priority + push 활성화 시 push 채널이어야 함: got %d", noti.Channel)
	}
}

func TestSendNotification_DefaultPriority(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	noti, err := svc.SendNotification(ctx, "user-dp", service.TypeSystem, service.ChannelInApp, service.PriorityUnknown, "제목", "내용", "")
	if err != nil {
		t.Fatalf("발송 실패: %v", err)
	}
	if noti.Priority != service.PriorityNormal {
		t.Fatalf("기본 우선순위가 Normal이어야 함: got %d", noti.Priority)
	}
}

func TestGetPreferences_EmptyUserID(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()
	_, err := svc.GetPreferences(ctx, "")
	if err == nil {
		t.Fatal("빈 user_id에 에러가 반환되어야 함")
	}
}

func TestListNotifications_DefaultLimit(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		svc.SendNotification(ctx, "user-dl", service.TypeSystem, service.ChannelInApp, service.PriorityNormal, "제목", "내용", "")
	}

	// limit=0이면 기본값(20) 적용
	notis, total, _, err := svc.ListNotifications(ctx, "user-dl", service.TypeUnknown, false, 0, 0)
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

func TestSendFromTemplate_EmptyUserID(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()
	err := svc.SendFromTemplate(ctx, "", "prescription_created", "김의사")
	if err == nil {
		t.Fatal("빈 user_id에 에러가 반환되어야 함")
	}
}

// =============================================================================
// Phase F 테스트 보강: 에스컬레이션 서비스
// =============================================================================

func setupTestEscalationService() (*service.EscalationService, *service.NotificationService) {
	notiSvc := setupTestService()
	logger := zap.NewNop()
	escSvc := service.NewEscalationService(logger, notiSvc)
	return escSvc, notiSvc
}

func TestTriggerEscalation_Success(t *testing.T) {
	escSvc, notiSvc := setupTestEscalationService()
	// 빠른 테스트를 위해 딜레이를 최소화
	escSvc.SetConfig(service.EscalationConfig{
		Stage2Delay: 24 * time.Hour,
		Stage3Delay: 24 * time.Hour,
		Stage4Delay: 24 * time.Hour,
	})

	ctx := context.Background()
	event, err := escSvc.TriggerEscalation(ctx, "user-esc1", "health_critical", "m-123", "350 mg/dL")
	if err != nil {
		t.Fatalf("에스컬레이션 시작 실패: %v", err)
	}
	if event.ID == "" {
		t.Fatal("이벤트 ID가 비어있음")
	}
	if event.Stage != service.StageInAppAlert {
		t.Fatalf("초기 스테이지 불일치: got %d, want %d", event.Stage, service.StageInAppAlert)
	}
	if event.UserID != "user-esc1" {
		t.Fatalf("UserID 불일치: got %s", event.UserID)
	}

	// Stage1에서 인앱 알림이 발송되었는지 확인
	notis, total, _, _ := notiSvc.ListNotifications(ctx, "user-esc1", service.TypeUnknown, false, 10, 0)
	if total != 1 {
		t.Fatalf("Stage1 알림 수 불일치: got %d, want 1", total)
	}
	if notis[0].Priority != service.PriorityUrgent {
		t.Fatalf("긴급 알림 우선순위 불일치: got %d", notis[0].Priority)
	}
}

func TestTriggerEscalation_EmptyUserID(t *testing.T) {
	escSvc, _ := setupTestEscalationService()
	ctx := context.Background()
	_, err := escSvc.TriggerEscalation(ctx, "", "health_critical", "m-1", "350")
	if err == nil {
		t.Fatal("빈 userID에 에러가 반환되어야 함")
	}
}

func TestTriggerEscalation_EmptyAlertType(t *testing.T) {
	escSvc, _ := setupTestEscalationService()
	ctx := context.Background()
	_, err := escSvc.TriggerEscalation(ctx, "user-1", "", "m-1", "350")
	if err == nil {
		t.Fatal("빈 alertType에 에러가 반환되어야 함")
	}
}

func TestAcknowledgeEscalation_Success(t *testing.T) {
	escSvc, _ := setupTestEscalationService()
	escSvc.SetConfig(service.EscalationConfig{
		Stage2Delay: 24 * time.Hour,
		Stage3Delay: 24 * time.Hour,
		Stage4Delay: 24 * time.Hour,
	})

	ctx := context.Background()
	event, _ := escSvc.TriggerEscalation(ctx, "user-ack", "fall_detected", "", "")

	err := escSvc.AcknowledgeEscalation(ctx, event.ID, "user")
	if err != nil {
		t.Fatalf("에스컬레이션 해제 실패: %v", err)
	}

	// 해제 후 active 목록에서 제거되어야 함
	active := escSvc.GetActiveEscalations("user-ack")
	if len(active) != 0 {
		t.Fatalf("해제 후 활성 에스컬레이션이 남아있음: got %d", len(active))
	}
}

func TestAcknowledgeEscalation_NotFound(t *testing.T) {
	escSvc, _ := setupTestEscalationService()
	ctx := context.Background()
	err := escSvc.AcknowledgeEscalation(ctx, "nonexistent-id", "user")
	if err == nil {
		t.Fatal("존재하지 않는 이벤트에 에러가 반환되어야 함")
	}
}

func TestGetActiveEscalations(t *testing.T) {
	escSvc, _ := setupTestEscalationService()
	escSvc.SetConfig(service.EscalationConfig{
		Stage2Delay: 24 * time.Hour,
		Stage3Delay: 24 * time.Hour,
		Stage4Delay: 24 * time.Hour,
	})

	ctx := context.Background()
	escSvc.TriggerEscalation(ctx, "user-multi", "health_critical", "m-1", "400")
	escSvc.TriggerEscalation(ctx, "user-multi", "fall_detected", "", "")
	escSvc.TriggerEscalation(ctx, "user-other", "no_response", "", "")

	active := escSvc.GetActiveEscalations("user-multi")
	if len(active) != 2 {
		t.Fatalf("활성 에스컬레이션 수 불일치: got %d, want 2", len(active))
	}

	activeOther := escSvc.GetActiveEscalations("user-other")
	if len(activeOther) != 1 {
		t.Fatalf("다른 사용자 활성 에스컬레이션 수 불일치: got %d, want 1", len(activeOther))
	}
}

func TestSetGuardians(t *testing.T) {
	escSvc, notiSvc := setupTestEscalationService()
	escSvc.SetConfig(service.EscalationConfig{
		Stage2Delay: 10 * time.Millisecond,
		Stage3Delay: 24 * time.Hour,
		Stage4Delay: 24 * time.Hour,
	})

	// 보호자 설정
	escSvc.SetGuardians("user-guard", []string{"guardian-1", "guardian-2"})

	ctx := context.Background()
	escSvc.TriggerEscalation(ctx, "user-guard", "health_critical", "m-1", "400")

	// Stage2 비동기 실행 대기
	time.Sleep(100 * time.Millisecond)

	// 보호자들에게 알림이 전송되었는지 확인
	g1Notis, g1Total, _, _ := notiSvc.ListNotifications(ctx, "guardian-1", service.TypeUnknown, false, 10, 0)
	if g1Total < 1 {
		t.Fatalf("보호자1에게 알림이 전송되어야 함: got %d", g1Total)
	}
	_ = g1Notis

	g2Notis, g2Total, _, _ := notiSvc.ListNotifications(ctx, "guardian-2", service.TypeUnknown, false, 10, 0)
	if g2Total < 1 {
		t.Fatalf("보호자2에게 알림이 전송되어야 함: got %d", g2Total)
	}
	_ = g2Notis
}

func TestEscalationConfig_Default(t *testing.T) {
	cfg := service.DefaultEscalationConfig
	if cfg.Stage2Delay != 3*time.Minute {
		t.Fatalf("Stage2Delay 기본값 불일치: got %v, want 3m", cfg.Stage2Delay)
	}
	if cfg.Stage3Delay != 3*time.Minute {
		t.Fatalf("Stage3Delay 기본값 불일치: got %v, want 3m", cfg.Stage3Delay)
	}
	if cfg.Stage4Delay != 4*time.Minute {
		t.Fatalf("Stage4Delay 기본값 불일치: got %v, want 4m", cfg.Stage4Delay)
	}
}

func TestEscalationStageConstants(t *testing.T) {
	if service.StageNone != 0 {
		t.Fatalf("StageNone = %d, want 0", service.StageNone)
	}
	if service.StageInAppAlert != 1 {
		t.Fatalf("StageInAppAlert = %d, want 1", service.StageInAppAlert)
	}
	if service.StageGuardianPush != 2 {
		t.Fatalf("StageGuardianPush = %d, want 2", service.StageGuardianPush)
	}
	if service.StageAIVoiceCall != 3 {
		t.Fatalf("StageAIVoiceCall = %d, want 3", service.StageAIVoiceCall)
	}
	if service.StageEmergencyCall != 4 {
		t.Fatalf("StageEmergencyCall = %d, want 4", service.StageEmergencyCall)
	}
	if service.StageResolved != 5 {
		t.Fatalf("StageResolved = %d, want 5", service.StageResolved)
	}
	if service.StageCancelled != 6 {
		t.Fatalf("StageCancelled = %d, want 6", service.StageCancelled)
	}
}
