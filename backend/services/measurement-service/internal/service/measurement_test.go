package service

import (
	"context"
	"testing"
	"time"

	"go.uber.org/zap"
)

// =============================================================================
// 목(Mock) 저장소
// =============================================================================

type mockSessionRepo struct {
	sessions map[string]*MeasurementSession
}

func newMockSessionRepo() *mockSessionRepo {
	return &mockSessionRepo{sessions: make(map[string]*MeasurementSession)}
}

func (r *mockSessionRepo) CreateSession(_ context.Context, session *MeasurementSession) error {
	r.sessions[session.ID] = session
	return nil
}

func (r *mockSessionRepo) GetSession(_ context.Context, sessionID string) (*MeasurementSession, error) {
	s, ok := r.sessions[sessionID]
	if !ok {
		return nil, nil
	}
	return s, nil
}

func (r *mockSessionRepo) EndSession(_ context.Context, sessionID string, total int, endedAt time.Time) error {
	s, ok := r.sessions[sessionID]
	if !ok {
		return nil
	}
	s.TotalMeasurements = total
	s.EndedAt = &endedAt
	s.Status = "completed"
	return nil
}

type mockMeasureRepo struct {
	data []*MeasurementData
}

func newMockMeasureRepo() *mockMeasureRepo {
	return &mockMeasureRepo{}
}

func (r *mockMeasureRepo) Store(_ context.Context, data *MeasurementData) error {
	r.data = append(r.data, data)
	return nil
}

func (r *mockMeasureRepo) GetHistory(_ context.Context, userID string, start, end time.Time, limit, offset int) ([]*MeasurementSummary, int, error) {
	var filtered []*MeasurementSummary
	for _, d := range r.data {
		if d.UserID != userID {
			continue
		}
		filtered = append(filtered, &MeasurementSummary{
			SessionID:     d.SessionID,
			CartridgeType: d.CartridgeType,
			PrimaryValue:  d.PrimaryValue,
			Unit:          d.Unit,
			MeasuredAt:    d.Time,
		})
	}
	total := len(filtered)
	if limit <= 0 {
		limit = 20
	}
	if offset >= total {
		return nil, total, nil
	}
	endIdx := offset + limit
	if endIdx > total {
		endIdx = total
	}
	return filtered[offset:endIdx], total, nil
}

type mockVectorRepo struct {
	vectors map[string][]float32
}

func newMockVectorRepo() *mockVectorRepo {
	return &mockVectorRepo{vectors: make(map[string][]float32)}
}

func (r *mockVectorRepo) StoreFingerprint(_ context.Context, sessionID string, vector []float32) error {
	r.vectors[sessionID] = vector
	return nil
}

func (r *mockVectorRepo) SearchSimilar(_ context.Context, _ []float32, _ int) ([]SimilarResult, error) {
	return nil, nil
}

type mockEventPublisher struct {
	events []*MeasurementCompletedEvent
}

func newMockEventPublisher() *mockEventPublisher {
	return &mockEventPublisher{}
}

func (p *mockEventPublisher) PublishMeasurementCompleted(_ context.Context, event *MeasurementCompletedEvent) error {
	p.events = append(p.events, event)
	return nil
}

// =============================================================================
// 헬퍼
// =============================================================================

func newTestMeasurementService() (*MeasurementService, *mockSessionRepo, *mockMeasureRepo, *mockVectorRepo, *mockEventPublisher) {
	logger, _ := zap.NewDevelopment()
	sessionRepo := newMockSessionRepo()
	measureRepo := newMockMeasureRepo()
	vectorRepo := newMockVectorRepo()
	eventPub := newMockEventPublisher()
	svc := NewMeasurementService(logger, sessionRepo, measureRepo, vectorRepo, eventPub)
	return svc, sessionRepo, measureRepo, vectorRepo, eventPub
}

// =============================================================================
// StartSession 테스트
// =============================================================================

func TestStartSession_성공(t *testing.T) {
	svc, sessionRepo, _, _, _ := newTestMeasurementService()
	ctx := context.Background()

	session, err := svc.StartSession(ctx, "device-1", "cartridge-glucose", "user-1")
	if err != nil {
		t.Fatalf("세션 시작 실패: %v", err)
	}
	if session.ID == "" {
		t.Error("세션 ID가 비어있습니다")
	}
	if session.DeviceID != "device-1" {
		t.Errorf("디바이스 ID 불일치: got %s", session.DeviceID)
	}
	if session.CartridgeID != "cartridge-glucose" {
		t.Errorf("카트리지 ID 불일치: got %s", session.CartridgeID)
	}
	if session.Status != "active" {
		t.Errorf("초기 상태가 active여야 합니다: got %s", session.Status)
	}
	if len(sessionRepo.sessions) != 1 {
		t.Errorf("세션이 1개 저장되어야 합니다: got %d", len(sessionRepo.sessions))
	}
}

func TestStartSession_빈_입력(t *testing.T) {
	svc, _, _, _, _ := newTestMeasurementService()
	ctx := context.Background()

	_, err := svc.StartSession(ctx, "", "cartridge-1", "user-1")
	if err == nil {
		t.Fatal("빈 device_id에 대해 에러가 발생해야 합니다")
	}

	_, err = svc.StartSession(ctx, "device-1", "", "user-1")
	if err == nil {
		t.Fatal("빈 cartridge_id에 대해 에러가 발생해야 합니다")
	}

	_, err = svc.StartSession(ctx, "device-1", "cartridge-1", "")
	if err == nil {
		t.Fatal("빈 user_id에 대해 에러가 발생해야 합니다")
	}
}

// =============================================================================
// ProcessMeasurement 테스트
// =============================================================================

func TestProcessMeasurement_성공(t *testing.T) {
	svc, _, measureRepo, vectorRepo, _ := newTestMeasurementService()
	ctx := context.Background()

	session, _ := svc.StartSession(ctx, "device-1", "cartridge-glucose", "user-1")

	result, err := svc.ProcessMeasurement(ctx, &MeasurementData{
		SessionID:         session.ID,
		DeviceID:          "device-1",
		UserID:            "user-1",
		CartridgeType:     "glucose",
		SDet:              1.5,
		SRef:              0.1,
		Alpha:             0.95,
		SCorrected:        1.405,
		PrimaryValue:      105.3,
		Unit:              "mg/dL",
		Confidence:        0.95,
		FingerprintVector: []float32{0.1, 0.2, 0.3},
		Time:              time.Now().UTC(),
	})
	if err != nil {
		t.Fatalf("측정 데이터 처리 실패: %v", err)
	}
	if result.PrimaryValue != 105.3 {
		t.Errorf("PrimaryValue 불일치: got %f, want 105.3", result.PrimaryValue)
	}
	if result.Unit != "mg/dL" {
		t.Errorf("Unit 불일치: got %s, want mg/dL", result.Unit)
	}
	if len(measureRepo.data) != 1 {
		t.Errorf("측정 데이터가 1개 저장되어야 합니다: got %d", len(measureRepo.data))
	}
	if len(vectorRepo.vectors) != 1 {
		t.Errorf("벡터가 1개 저장되어야 합니다: got %d", len(vectorRepo.vectors))
	}
}

func TestProcessMeasurement_빈_세션ID(t *testing.T) {
	svc, _, _, _, _ := newTestMeasurementService()
	ctx := context.Background()

	_, err := svc.ProcessMeasurement(ctx, &MeasurementData{SessionID: ""})
	if err == nil {
		t.Fatal("빈 session_id에 대해 에러가 발생해야 합니다")
	}
}

// =============================================================================
// EndSession 테스트
// =============================================================================

func TestEndSession_성공(t *testing.T) {
	svc, _, _, _, eventPub := newTestMeasurementService()
	ctx := context.Background()

	session, _ := svc.StartSession(ctx, "device-1", "cartridge-glucose", "user-1")

	result, err := svc.EndSession(ctx, session.ID)
	if err != nil {
		t.Fatalf("세션 종료 실패: %v", err)
	}
	if result.SessionID != session.ID {
		t.Errorf("세션 ID 불일치: got %s", result.SessionID)
	}
	if result.EndedAt.IsZero() {
		t.Error("종료 시각이 설정되어야 합니다")
	}
	if len(eventPub.events) != 1 {
		t.Errorf("완료 이벤트가 1개 발행되어야 합니다: got %d", len(eventPub.events))
	}
}

func TestEndSession_빈_세션ID(t *testing.T) {
	svc, _, _, _, _ := newTestMeasurementService()
	ctx := context.Background()

	_, err := svc.EndSession(ctx, "")
	if err == nil {
		t.Fatal("빈 session_id에 대해 에러가 발생해야 합니다")
	}
}

func TestEndSession_존재하지_않는_세션(t *testing.T) {
	svc, _, _, _, _ := newTestMeasurementService()
	ctx := context.Background()

	_, err := svc.EndSession(ctx, "nonexistent-session")
	if err == nil {
		t.Fatal("존재하지 않는 세션에 대해 에러가 발생해야 합니다")
	}
}

// =============================================================================
// GetHistory 테스트
// =============================================================================

func TestGetHistory_성공(t *testing.T) {
	svc, _, measureRepo, _, _ := newTestMeasurementService()
	ctx := context.Background()

	// 시드 데이터
	now := time.Now().UTC()
	measureRepo.data = append(measureRepo.data,
		&MeasurementData{
			SessionID:     "sess-1",
			UserID:        "user-1",
			CartridgeType: "glucose",
			PrimaryValue:  100.0,
			Unit:          "mg/dL",
			Time:          now.Add(-1 * time.Hour),
		},
		&MeasurementData{
			SessionID:     "sess-2",
			UserID:        "user-1",
			CartridgeType: "lipid",
			PrimaryValue:  200.0,
			Unit:          "mg/dL",
			Time:          now.Add(-30 * time.Minute),
		},
		&MeasurementData{
			SessionID:     "sess-3",
			UserID:        "user-2", // 다른 사용자
			CartridgeType: "glucose",
			PrimaryValue:  110.0,
			Unit:          "mg/dL",
			Time:          now,
		},
	)

	summaries, total, err := svc.GetHistory(ctx, "user-1", time.Time{}, time.Time{}, 10, 0)
	if err != nil {
		t.Fatalf("측정 기록 조회 실패: %v", err)
	}
	if total != 2 {
		t.Errorf("user-1의 측정 기록은 2개여야 합니다: got %d", total)
	}
	if len(summaries) != 2 {
		t.Errorf("반환된 요약은 2개여야 합니다: got %d", len(summaries))
	}
}

func TestGetHistory_빈_유저ID(t *testing.T) {
	svc, _, _, _, _ := newTestMeasurementService()
	ctx := context.Background()

	_, _, err := svc.GetHistory(ctx, "", time.Time{}, time.Time{}, 10, 0)
	if err == nil {
		t.Fatal("빈 user_id에 대해 에러가 발생해야 합니다")
	}
}

// =============================================================================
// ExportSingleMeasurement 테스트
// =============================================================================

func TestExportSingleMeasurement_성공(t *testing.T) {
	svc, _, measureRepo, _, _ := newTestMeasurementService()
	ctx := context.Background()

	// 세션 생성
	session, err := svc.StartSession(ctx, "device-export", "cartridge-glucose", "user-export")
	if err != nil {
		t.Fatalf("세션 시작 실패: %v", err)
	}

	// 측정 데이터 시드 (같은 세션에 2개)
	now := time.Now().UTC()
	measureRepo.data = append(measureRepo.data,
		&MeasurementData{
			SessionID:     session.ID,
			UserID:        "user-export",
			CartridgeType: "blood_glucose",
			PrimaryValue:  105.3,
			Unit:          "mg/dL",
			Time:          now.Add(-5 * time.Minute),
		},
		&MeasurementData{
			SessionID:     session.ID,
			UserID:        "user-export",
			CartridgeType: "blood_glucose",
			PrimaryValue:  108.1,
			Unit:          "mg/dL",
			Time:          now,
		},
		&MeasurementData{
			SessionID:     "other-session",
			UserID:        "user-export",
			CartridgeType: "cholesterol_total",
			PrimaryValue:  200.0,
			Unit:          "mg/dL",
			Time:          now,
		},
	)

	bundleJSON, err := svc.ExportSingleMeasurement(ctx, session.ID)
	if err != nil {
		t.Fatalf("ExportSingleMeasurement 실패: %v", err)
	}

	if bundleJSON == "" {
		t.Fatal("FHIR Bundle JSON이 비어 있습니다")
	}

	// Bundle에 해당 세션의 측정만 포함되는지 확인
	if !contains(bundleJSON, session.ID) {
		t.Error("Bundle에 세션 ID가 포함되어야 합니다")
	}
	if contains(bundleJSON, "other-session") {
		t.Error("Bundle에 다른 세션 데이터가 포함되면 안 됩니다")
	}

	// LOINC 코드 확인 (blood_glucose → 15074-8)
	if !contains(bundleJSON, "15074-8") {
		t.Error("Bundle에 blood_glucose LOINC 코드(15074-8)가 포함되어야 합니다")
	}
}

func TestExportSingleMeasurement_빈_세션ID(t *testing.T) {
	svc, _, _, _, _ := newTestMeasurementService()
	ctx := context.Background()

	_, err := svc.ExportSingleMeasurement(ctx, "")
	if err == nil {
		t.Fatal("빈 session_id에 대해 에러가 발생해야 합니다")
	}
}

func TestExportSingleMeasurement_존재하지_않는_세션(t *testing.T) {
	svc, _, _, _, _ := newTestMeasurementService()
	ctx := context.Background()

	_, err := svc.ExportSingleMeasurement(ctx, "nonexistent-session")
	if err == nil {
		t.Fatal("존재하지 않는 세션에 대해 에러가 발생해야 합니다")
	}
}

func TestDefaultLOINCMap(t *testing.T) {
	expected := map[string]string{
		"blood_glucose":     "15074-8",
		"blood_pressure":    "85354-9",
		"cholesterol_total": "2093-3",
		"hemoglobin_a1c":    "4548-4",
		"heart_rate":        "8867-4",
		"body_temperature":  "8310-5",
		"oxygen_saturation": "2708-6",
	}
	for biomarker, wantCode := range expected {
		got := loincCodeFor(biomarker)
		if got != wantCode {
			t.Errorf("loincCodeFor(%s) = %s, want %s", biomarker, got, wantCode)
		}
	}

	// 알 수 없는 바이오마커 → 기본 코드
	got := loincCodeFor("unknown_biomarker")
	if got != "29463-7" {
		t.Errorf("loincCodeFor(unknown) = %s, want 29463-7", got)
	}
}

// contains는 문자열에 부분 문자열이 포함되어 있는지 확인합니다.
func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestGetHistory_기본_Limit(t *testing.T) {
	svc, _, _, _, _ := newTestMeasurementService()
	ctx := context.Background()

	// limit <= 0이면 기본값 20 사용
	summaries, total, err := svc.GetHistory(ctx, "user-empty", time.Time{}, time.Time{}, 0, 0)
	if err != nil {
		t.Fatalf("기록 조회 실패: %v", err)
	}
	if total != 0 {
		t.Errorf("빈 사용자의 기록은 0이어야 합니다: got %d", total)
	}
	if len(summaries) != 0 {
		t.Errorf("빈 결과여야 합니다: got %d", len(summaries))
	}
}
