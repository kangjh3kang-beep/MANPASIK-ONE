package service_test

import (
	"context"
	"testing"

	"github.com/manpasik/backend/services/emergency-service/internal/notifier"
	"github.com/manpasik/backend/services/emergency-service/internal/repository/memory"
	"github.com/manpasik/backend/services/emergency-service/internal/service"
)

func newTestService() *service.EmergencyService {
	repo := memory.NewEmergencyRepository()
	return service.NewEmergencyService(repo)
}

func TestReportEmergency_Success(t *testing.T) {
	svc := newTestService()

	id, err := svc.ReportEmergency(service.ReportEmergencyInput{
		UserID:      "user-1",
		Type:        "cardiac",
		Location:    "37.5665,126.9780",
		Description: "Chest pain reported",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id == "" {
		t.Fatal("expected non-empty emergency ID")
	}
}

func TestReportEmergency_MissingUserID(t *testing.T) {
	svc := newTestService()

	_, err := svc.ReportEmergency(service.ReportEmergencyInput{
		Type: "fall",
	})
	if err == nil {
		t.Fatal("expected error for missing user_id")
	}
}

func TestReportEmergency_MissingType(t *testing.T) {
	svc := newTestService()

	_, err := svc.ReportEmergency(service.ReportEmergencyInput{
		UserID: "user-1",
	})
	if err == nil {
		t.Fatal("expected error for missing type")
	}
}

func TestGetEmergencyContacts_Empty(t *testing.T) {
	svc := newTestService()

	contacts, err := svc.GetEmergencyContacts("user-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(contacts) != 0 {
		t.Fatalf("expected 0 contacts, got %d", len(contacts))
	}
}

func TestGetEmergencyContacts_WithContacts(t *testing.T) {
	repo := memory.NewEmergencyRepository()
	svc := service.NewEmergencyService(repo)

	_ = repo.AddContact(&service.EmergencyContact{
		UserID:       "user-1",
		Name:         "Kim",
		Phone:        "010-1234-5678",
		Relationship: "spouse",
	})
	_ = repo.AddContact(&service.EmergencyContact{
		UserID:       "user-1",
		Name:         "Lee",
		Phone:        "010-8765-4321",
		Relationship: "parent",
	})

	contacts, err := svc.GetEmergencyContacts("user-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(contacts) != 2 {
		t.Fatalf("expected 2 contacts, got %d", len(contacts))
	}
}

func TestUpdateAndGetEmergencySettings(t *testing.T) {
	svc := newTestService()

	settings := &service.EmergencySettings{
		UserID:              "user-1",
		AutoCall119:         true,
		EmergencyContactIDs: []string{"c-1", "c-2"},
		MedicalInfo:         "Allergic to penicillin",
	}

	err := svc.UpdateEmergencySettings(settings)
	if err != nil {
		t.Fatalf("unexpected error on update: %v", err)
	}

	got, err := svc.GetEmergencySettings("user-1")
	if err != nil {
		t.Fatalf("unexpected error on get: %v", err)
	}
	if !got.AutoCall119 {
		t.Fatal("expected AutoCall119 to be true")
	}
	if len(got.EmergencyContactIDs) != 2 {
		t.Fatalf("expected 2 contact IDs, got %d", len(got.EmergencyContactIDs))
	}
	if got.MedicalInfo != "Allergic to penicillin" {
		t.Fatalf("unexpected MedicalInfo: %s", got.MedicalInfo)
	}
}

func TestGetEmergencySettings_Default(t *testing.T) {
	svc := newTestService()

	got, err := svc.GetEmergencySettings("user-new")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.AutoCall119 {
		t.Fatal("expected AutoCall119 to default to false")
	}
	if got.UserID != "user-new" {
		t.Fatalf("expected UserID user-new, got %s", got.UserID)
	}
}

func TestReportEmergency_AttachesContacts(t *testing.T) {
	repo := memory.NewEmergencyRepository()
	svc := service.NewEmergencyService(repo)

	_ = repo.AddContact(&service.EmergencyContact{
		UserID:       "user-1",
		Name:         "Kim",
		Phone:        "010-1234-5678",
		Relationship: "spouse",
	})

	id, err := svc.ReportEmergency(service.ReportEmergencyInput{
		UserID:      "user-1",
		Type:        "fall",
		Location:    "37.5665,126.9780",
		Description: "Fall detected",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	emg, err := repo.GetEmergency(id)
	if err != nil {
		t.Fatalf("unexpected error getting emergency: %v", err)
	}
	if len(emg.ContactIDs) != 1 {
		t.Fatalf("expected 1 contact ID attached, got %d", len(emg.ContactIDs))
	}
}

// ============================================================================
// Phase B 신규 테스트
// ============================================================================

func TestReportEmergency_SeverityClassification(t *testing.T) {
	svc := newTestService()

	// cardiac → Critical
	id, err := svc.ReportEmergency(service.ReportEmergencyInput{
		UserID: "user-1", Type: "cardiac", Description: "Heart attack",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	repo := memory.NewEmergencyRepository()
	svcWithRepo := service.NewEmergencyService(repo)

	id2, _ := svcWithRepo.ReportEmergency(service.ReportEmergencyInput{
		UserID: "user-1", Type: "other", Description: "Minor issue",
	})
	emg, _ := repo.GetEmergency(id2)
	if emg.Severity != service.SeverityLow {
		t.Fatalf("other type expected SeverityLow, got %d", emg.Severity)
	}
	_ = id
}

func TestReportEmergency_AutoDispatch(t *testing.T) {
	repo := memory.NewEmergencyRepository()
	svc := service.NewEmergencyService(repo)

	// AutoCall119 설정
	_ = repo.SaveSettings(&service.EmergencySettings{
		UserID:      "user-1",
		AutoCall119: true,
	})

	// cardiac(Critical) + AutoCall119 → dispatched
	id, err := svc.ReportEmergency(service.ReportEmergencyInput{
		UserID: "user-1", Type: "cardiac", Description: "Emergency",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	emg, _ := repo.GetEmergency(id)
	if emg.Status != "dispatched" {
		t.Fatalf("expected status dispatched, got %s", emg.Status)
	}

	// other(Low) + AutoCall119 → reported (심각도 부족)
	id2, _ := svc.ReportEmergency(service.ReportEmergencyInput{
		UserID: "user-1", Type: "other", Description: "Minor issue",
	})
	emg2, _ := repo.GetEmergency(id2)
	if emg2.Status != "reported" {
		t.Fatalf("expected status reported for low severity, got %s", emg2.Status)
	}
}

func TestAddAndRemoveContact(t *testing.T) {
	repo := memory.NewEmergencyRepository()
	svc := service.NewEmergencyService(repo)

	// 2개 추가
	err := svc.AddContact("user-1", &service.EmergencyContact{
		Name: "Kim", Phone: "010-1111-2222", Relationship: "spouse",
	})
	if err != nil {
		t.Fatalf("add contact 1 failed: %v", err)
	}
	err = svc.AddContact("user-1", &service.EmergencyContact{
		Name: "Lee", Phone: "010-3333-4444", Relationship: "parent",
	})
	if err != nil {
		t.Fatalf("add contact 2 failed: %v", err)
	}

	contacts, _ := svc.GetEmergencyContacts("user-1")
	if len(contacts) != 2 {
		t.Fatalf("expected 2 contacts, got %d", len(contacts))
	}

	// 첫 번째 제거
	err = svc.RemoveContact("user-1", contacts[0].ID)
	if err != nil {
		t.Fatalf("remove contact failed: %v", err)
	}

	contacts, _ = svc.GetEmergencyContacts("user-1")
	if len(contacts) != 1 {
		t.Fatalf("expected 1 contact after removal, got %d", len(contacts))
	}
}

func TestResolveEmergency(t *testing.T) {
	repo := memory.NewEmergencyRepository()
	svc := service.NewEmergencyService(repo)

	id, _ := svc.ReportEmergency(service.ReportEmergencyInput{
		UserID: "user-1", Type: "fall", Description: "Fall detected",
	})

	err := svc.ResolveEmergency(id, "Patient recovered")
	if err != nil {
		t.Fatalf("resolve failed: %v", err)
	}

	emg, _ := repo.GetEmergency(id)
	if emg.Status != "resolved" {
		t.Fatalf("expected resolved, got %s", emg.Status)
	}
	if emg.Resolution != "Patient recovered" {
		t.Fatalf("expected resolution 'Patient recovered', got %s", emg.Resolution)
	}
	if emg.ResolvedAt == nil {
		t.Fatal("expected ResolvedAt to be set")
	}
}

func TestResolveEmergency_AlreadyResolved(t *testing.T) {
	repo := memory.NewEmergencyRepository()
	svc := service.NewEmergencyService(repo)

	id, _ := svc.ReportEmergency(service.ReportEmergencyInput{
		UserID: "user-1", Type: "fall", Description: "Fall",
	})

	_ = svc.ResolveEmergency(id, "Resolved")

	// 이미 해결된 응급 상황 재해결 시도
	err := svc.ResolveEmergency(id, "Resolved again")
	if err == nil {
		t.Fatal("expected error for already resolved emergency")
	}
}

func TestGetEmergencyHistory(t *testing.T) {
	repo := memory.NewEmergencyRepository()
	svc := service.NewEmergencyService(repo)

	// 3개 보고
	for i := 0; i < 3; i++ {
		_, err := svc.ReportEmergency(service.ReportEmergencyInput{
			UserID: "user-1", Type: "fall", Description: "Fall event",
		})
		if err != nil {
			t.Fatalf("report %d failed: %v", i, err)
		}
	}

	// 다른 사용자 보고 (필터 확인)
	_, _ = svc.ReportEmergency(service.ReportEmergencyInput{
		UserID: "user-2", Type: "other", Description: "Other",
	})

	history, err := svc.GetEmergencyHistory("user-1")
	if err != nil {
		t.Fatalf("get history failed: %v", err)
	}
	if len(history) != 3 {
		t.Fatalf("expected 3 history items, got %d", len(history))
	}
}

// ============================================================================
// EmergencyNotifier 통합 테스트
// ============================================================================

type mockNotifier struct {
	dispatch119Called   bool
	notifyContactsCalled bool
	dispatch119Err     error
	notifyContactsErr  error
}

func (m *mockNotifier) Dispatch119(_ context.Context, _ *notifier.EmergencyAlert) (*notifier.DispatchResult, error) {
	m.dispatch119Called = true
	if m.dispatch119Err != nil {
		return nil, m.dispatch119Err
	}
	return &notifier.DispatchResult{Success: true, DispatchID: "MOCK-001"}, nil
}

func (m *mockNotifier) NotifyContacts(_ context.Context, _ *notifier.EmergencyAlert, _ []string) error {
	m.notifyContactsCalled = true
	return m.notifyContactsErr
}

func (m *mockNotifier) ProviderName() string { return "mock" }

func TestReportEmergency_WithNotifier_HighSeverity(t *testing.T) {
	repo := memory.NewEmergencyRepository()
	svc := service.NewEmergencyService(repo)
	mock := &mockNotifier{}
	svc.SetEmergencyNotifier(mock)

	// AutoCall119 활성화
	_ = repo.SaveSettings(&service.EmergencySettings{
		UserID:      "user-notif",
		AutoCall119: true,
	})

	// cardiac = Critical → 119 자동 신고
	_, err := svc.ReportEmergency(service.ReportEmergencyInput{
		UserID:      "user-notif",
		Type:        "cardiac",
		Description: "흉통 발생",
		Location:    "서울",
	})
	if err != nil {
		t.Fatalf("ReportEmergency 실패: %v", err)
	}
	if !mock.dispatch119Called {
		t.Error("Critical 심각도에서 Dispatch119가 호출되어야 합니다")
	}
}

func TestReportEmergency_WithNotifier_LowSeverity(t *testing.T) {
	repo := memory.NewEmergencyRepository()
	svc := service.NewEmergencyService(repo)
	mock := &mockNotifier{}
	svc.SetEmergencyNotifier(mock)

	// Low severity → 119 미호출
	_, err := svc.ReportEmergency(service.ReportEmergencyInput{
		UserID:      "user-low",
		Type:        "other",
		Description: "경미한 사고",
	})
	if err != nil {
		t.Fatalf("ReportEmergency 실패: %v", err)
	}
	if mock.dispatch119Called {
		t.Error("Low 심각도에서는 Dispatch119가 호출되지 않아야 합니다")
	}
}

func TestReportEmergency_WithLogNotifier(t *testing.T) {
	repo := memory.NewEmergencyRepository()
	svc := service.NewEmergencyService(repo)
	svc.SetEmergencyNotifier(notifier.NewLogNotifier())

	_ = repo.SaveSettings(&service.EmergencySettings{
		UserID:      "user-log",
		AutoCall119: true,
	})

	id, err := svc.ReportEmergency(service.ReportEmergencyInput{
		UserID:      "user-log",
		Type:        "cardiac",
		Description: "흉통",
		Location:    "서울",
	})
	if err != nil {
		t.Fatalf("LogNotifier 사용 시 에러 발생: %v", err)
	}
	if id == "" {
		t.Fatal("응급 ID가 비어 있음")
	}
}

func TestReportEmergency_WithoutNotifier(t *testing.T) {
	svc := newTestService()

	id, err := svc.ReportEmergency(service.ReportEmergencyInput{
		UserID:      "user-no-notif",
		Type:        "cardiac",
		Description: "흉통",
	})
	if err != nil {
		t.Fatalf("Notifier 없이도 에러가 발생하지 않아야 합니다: %v", err)
	}
	if id == "" {
		t.Fatal("응급 ID가 비어 있음")
	}
}
