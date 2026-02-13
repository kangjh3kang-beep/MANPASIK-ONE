package service

import (
	"context"
	"math"
	"testing"
	"time"

	"go.uber.org/zap"
)

// ============================================================================
// 테스트용 fakeCalRepo (CalibrationRepository)
// ============================================================================

type fakeCalRepo struct {
	records []*CalibrationRecord
}

func newFakeCalRepo() *fakeCalRepo {
	return &fakeCalRepo{records: make([]*CalibrationRecord, 0)}
}

func (r *fakeCalRepo) Save(_ context.Context, record *CalibrationRecord) error {
	cp := *record
	if record.ChannelOffsets != nil {
		cp.ChannelOffsets = make([]float64, len(record.ChannelOffsets))
		copy(cp.ChannelOffsets, record.ChannelOffsets)
	}
	if record.ChannelGains != nil {
		cp.ChannelGains = make([]float64, len(record.ChannelGains))
		copy(cp.ChannelGains, record.ChannelGains)
	}
	r.records = append(r.records, &cp)
	return nil
}

func (r *fakeCalRepo) GetLatest(_ context.Context, deviceID string, category, typeIndex int32) (*CalibrationRecord, error) {
	var latest *CalibrationRecord
	for _, rec := range r.records {
		if rec.DeviceID == deviceID && rec.CartridgeCategory == category && rec.CartridgeTypeIndex == typeIndex {
			if latest == nil || rec.CalibratedAt.After(latest.CalibratedAt) {
				latest = rec
			}
		}
	}
	if latest == nil {
		return nil, nil
	}
	cp := *latest
	return &cp, nil
}

func (r *fakeCalRepo) GetLatestByType(_ context.Context, deviceID string, category, typeIndex int32, calType CalibrationType) (*CalibrationRecord, error) {
	var latest *CalibrationRecord
	for _, rec := range r.records {
		if rec.DeviceID == deviceID && rec.CartridgeCategory == category && rec.CartridgeTypeIndex == typeIndex && rec.CalibrationType == calType {
			if latest == nil || rec.CalibratedAt.After(latest.CalibratedAt) {
				latest = rec
			}
		}
	}
	if latest == nil {
		return nil, nil
	}
	cp := *latest
	return &cp, nil
}

func (r *fakeCalRepo) ListByDevice(_ context.Context, deviceID string, limit, offset int32) ([]*CalibrationRecord, int32, error) {
	var filtered []*CalibrationRecord
	for _, rec := range r.records {
		if rec.DeviceID == deviceID {
			filtered = append(filtered, rec)
		}
	}
	total := int32(len(filtered))
	start := int(offset)
	if start >= len(filtered) {
		return nil, total, nil
	}
	end := start + int(limit)
	if end > len(filtered) {
		end = len(filtered)
	}
	return filtered[start:end], total, nil
}

// ============================================================================
// 테스트용 fakeModelRepo (CalibrationModelRepository)
// ============================================================================

type fakeModelRepo struct {
	models []*CalibrationModel
}

func newFakeModelRepo() *fakeModelRepo {
	return &fakeModelRepo{models: DefaultCalibrationModels()}
}

func (r *fakeModelRepo) GetAll(_ context.Context) ([]*CalibrationModel, error) {
	return r.models, nil
}

func (r *fakeModelRepo) GetByCartridgeType(_ context.Context, category, typeIndex int32) (*CalibrationModel, error) {
	for _, m := range r.models {
		if m.CartridgeCategory == category && m.CartridgeTypeIndex == typeIndex {
			return m, nil
		}
	}
	return nil, nil
}

// ============================================================================
// 테스트 헬퍼
// ============================================================================

func newTestCalibrationService() (*CalibrationService, *fakeCalRepo, *fakeModelRepo) {
	calRepo := newFakeCalRepo()
	modelRepo := newFakeModelRepo()
	svc := NewCalibrationService(zap.NewNop(), calRepo, modelRepo)
	return svc, calRepo, modelRepo
}

// ============================================================================
// 테스트
// ============================================================================

func TestRegisterFactoryCalibration(t *testing.T) {
	svc, calRepo, _ := newTestCalibrationService()
	ctx := context.Background()

	record, err := svc.RegisterFactoryCalibration(
		ctx,
		"device-001",
		1, 1, // HealthBiomarker Glucose
		0.95,
		[]float64{0.01, -0.02, 0.005},
		[]float64{1.01, 0.99, 1.005},
		0.003,
		0.001,
		"NIST-SRM-917c",
		"factory-technician-1",
	)
	if err != nil {
		t.Fatalf("RegisterFactoryCalibration 실패: %v", err)
	}

	// 기본 필드 확인
	if record.ID == "" {
		t.Error("ID가 비어 있습니다")
	}
	if record.DeviceID != "device-001" {
		t.Errorf("DeviceID: got %s, want device-001", record.DeviceID)
	}
	if record.CalibrationType != CalibrationTypeFactory {
		t.Errorf("CalibrationType: got %d, want %d", record.CalibrationType, CalibrationTypeFactory)
	}
	if record.Alpha != 0.95 {
		t.Errorf("Alpha: got %f, want 0.95", record.Alpha)
	}
	if record.Status != CalibrationStatusValid {
		t.Errorf("Status: got %d, want %d", record.Status, CalibrationStatusValid)
	}

	// 채널 오프셋/게인 저장 확인
	if len(record.ChannelOffsets) != 3 {
		t.Errorf("ChannelOffsets 길이: got %d, want 3", len(record.ChannelOffsets))
	}
	if len(record.ChannelGains) != 3 {
		t.Errorf("ChannelGains 길이: got %d, want 3", len(record.ChannelGains))
	}

	// 유효기간 확인 (HealthBiomarker: 90일)
	expectedExpiry := record.CalibratedAt.AddDate(0, 0, 90)
	if !record.ExpiresAt.Equal(expectedExpiry) {
		t.Errorf("ExpiresAt: got %v, want %v", record.ExpiresAt, expectedExpiry)
	}

	// 정확도 점수 확인 (alpha=0.95, 모델 기본=0.95 → 정확도 높음)
	if record.AccuracyScore < 0.9 {
		t.Errorf("AccuracyScore가 너무 낮습니다: %f", record.AccuracyScore)
	}

	// ReferenceStandard 확인
	if record.ReferenceStandard != "NIST-SRM-917c" {
		t.Errorf("ReferenceStandard: got %s, want NIST-SRM-917c", record.ReferenceStandard)
	}

	// 저장소에 기록되었는지 확인
	if len(calRepo.records) != 1 {
		t.Errorf("저장소 기록 수: got %d, want 1", len(calRepo.records))
	}
}

func TestPerformFieldCalibration(t *testing.T) {
	svc, _, _ := newTestCalibrationService()
	ctx := context.Background()

	// 기준값 [100, 200, 300], 측정값 [95, 190, 285]
	// Alpha = mean(95/100, 190/200, 285/300) = mean(0.95, 0.95, 0.95) = 0.95
	referenceValues := []float64{100.0, 200.0, 300.0}
	measuredValues := []float64{95.0, 190.0, 285.0}

	record, err := svc.PerformFieldCalibration(
		ctx,
		"device-002",
		"user-001",
		1, 1,
		referenceValues,
		measuredValues,
		25.0,
		50.0,
	)
	if err != nil {
		t.Fatalf("PerformFieldCalibration 실패: %v", err)
	}

	if record.CalibrationType != CalibrationTypeField {
		t.Errorf("CalibrationType: got %d, want %d", record.CalibrationType, CalibrationTypeField)
	}

	// Alpha 확인: mean(95/100, 190/200, 285/300) = 0.95
	expectedAlpha := 0.95
	if math.Abs(record.Alpha-expectedAlpha) > 0.001 {
		t.Errorf("Alpha: got %f, want ~%f", record.Alpha, expectedAlpha)
	}

	// 유효기간 확인 (HealthBiomarker model validity=90, field=90/3=30일)
	expectedExpiry := record.CalibratedAt.AddDate(0, 0, 30)
	if !record.ExpiresAt.Equal(expectedExpiry) {
		t.Errorf("ExpiresAt: got %v, want %v", record.ExpiresAt, expectedExpiry)
	}

	// 보정 후 잔차가 0이므로 정확도 높음
	if record.AccuracyScore < 0.9 {
		t.Errorf("AccuracyScore가 너무 낮습니다: %f (잔차가 0이므로 1.0에 가까워야 함)", record.AccuracyScore)
	}

	// CalibratedBy = user_id
	if record.CalibratedBy != "user-001" {
		t.Errorf("CalibratedBy: got %s, want user-001", record.CalibratedBy)
	}
}

func TestPerformFieldCalibration_MismatchedLengths(t *testing.T) {
	svc, _, _ := newTestCalibrationService()
	ctx := context.Background()

	_, err := svc.PerformFieldCalibration(
		ctx,
		"device-003",
		"user-001",
		1, 1,
		[]float64{100.0, 200.0},
		[]float64{95.0},
		25.0,
		50.0,
	)
	if err == nil {
		t.Error("기준값과 측정값 개수 불일치 시 에러가 반환되어야 합니다")
	}
}

func TestGetCalibration(t *testing.T) {
	svc, _, _ := newTestCalibrationService()
	ctx := context.Background()

	// 팩토리 보정 등록
	_, err := svc.RegisterFactoryCalibration(
		ctx,
		"device-010",
		1, 1,
		0.96,
		[]float64{0.01},
		[]float64{1.0},
		0.003,
		0.001,
		"NIST-SRM-917c",
		"tech-1",
	)
	if err != nil {
		t.Fatalf("팩토리 보정 등록 실패: %v", err)
	}

	// 현장 보정도 등록
	_, err = svc.PerformFieldCalibration(
		ctx,
		"device-010",
		"user-001",
		1, 1,
		[]float64{100.0},
		[]float64{94.0},
		25.0,
		50.0,
	)
	if err != nil {
		t.Fatalf("현장 보정 등록 실패: %v", err)
	}

	// GetCalibration → 팩토리 우선
	record, err := svc.GetCalibration(ctx, "device-010", 1, 1)
	if err != nil {
		t.Fatalf("GetCalibration 실패: %v", err)
	}
	if record.CalibrationType != CalibrationTypeFactory {
		t.Errorf("팩토리 보정이 우선되어야 합니다: got type %d, want %d", record.CalibrationType, CalibrationTypeFactory)
	}
	if record.Alpha != 0.96 {
		t.Errorf("Alpha: got %f, want 0.96", record.Alpha)
	}
}

func TestGetCalibration_NotFound(t *testing.T) {
	svc, _, _ := newTestCalibrationService()
	ctx := context.Background()

	_, err := svc.GetCalibration(ctx, "non-existent-device", 1, 1)
	if err == nil {
		t.Error("존재하지 않는 보정 데이터 조회가 성공했습니다")
	}
}

func TestListCalibrationHistory(t *testing.T) {
	svc, _, _ := newTestCalibrationService()
	ctx := context.Background()

	// 3개의 보정 기록 생성
	for i := 0; i < 3; i++ {
		_, err := svc.RegisterFactoryCalibration(
			ctx,
			"device-020",
			1, 1,
			0.95+float64(i)*0.01,
			nil, nil,
			0.003, 0.001,
			"NIST-SRM-917c",
			"tech-1",
		)
		if err != nil {
			t.Fatalf("보정 등록 실패 (i=%d): %v", i, err)
		}
	}

	records, total, err := svc.ListCalibrationHistory(ctx, "device-020", 10, 0)
	if err != nil {
		t.Fatalf("ListCalibrationHistory 실패: %v", err)
	}
	if total != 3 {
		t.Errorf("total: got %d, want 3", total)
	}
	if len(records) != 3 {
		t.Errorf("records 수: got %d, want 3", len(records))
	}
}

func TestListCalibrationHistory_Pagination(t *testing.T) {
	svc, _, _ := newTestCalibrationService()
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		_, _ = svc.RegisterFactoryCalibration(
			ctx,
			"device-021",
			1, 1,
			0.95,
			nil, nil,
			0.003, 0.001,
			"", "tech",
		)
	}

	records, total, err := svc.ListCalibrationHistory(ctx, "device-021", 2, 1)
	if err != nil {
		t.Fatalf("ListCalibrationHistory 페이지네이션 실패: %v", err)
	}
	if total != 5 {
		t.Errorf("total: got %d, want 5", total)
	}
	if len(records) != 2 {
		t.Errorf("records 수: got %d, want 2", len(records))
	}
}

func TestCheckCalibrationStatus_Valid(t *testing.T) {
	svc, _, _ := newTestCalibrationService()
	ctx := context.Background()

	_, err := svc.RegisterFactoryCalibration(
		ctx,
		"device-030",
		1, 1,
		0.95,
		nil, nil,
		0.003, 0.001,
		"", "tech",
	)
	if err != nil {
		t.Fatalf("보정 등록 실패: %v", err)
	}

	st, msg, latest, err := svc.CheckCalibrationStatus(ctx, "device-030", 1, 1)
	if err != nil {
		t.Fatalf("CheckCalibrationStatus 실패: %v", err)
	}
	if st != CalibrationStatusValid {
		t.Errorf("Status: got %d, want %d (Valid)", st, CalibrationStatusValid)
	}
	if msg == "" {
		t.Error("메시지가 비어 있습니다")
	}
	if latest == nil {
		t.Error("latest가 nil입니다")
	}
}

func TestCheckCalibrationStatus_Expired(t *testing.T) {
	svc, calRepo, _ := newTestCalibrationService()
	ctx := context.Background()

	// 이미 만료된 보정 기록을 직접 삽입
	expiredRecord := &CalibrationRecord{
		ID:                 "expired-cal-001",
		DeviceID:           "device-031",
		CartridgeCategory:  1,
		CartridgeTypeIndex: 1,
		CalibrationType:    CalibrationTypeFactory,
		Alpha:              0.95,
		CalibratedAt:       time.Now().UTC().AddDate(0, 0, -100),
		ExpiresAt:          time.Now().UTC().AddDate(0, 0, -10), // 10일 전 만료
		Status:             CalibrationStatusValid,
	}
	_ = calRepo.Save(ctx, expiredRecord)

	st, msg, _, err := svc.CheckCalibrationStatus(ctx, "device-031", 1, 1)
	if err != nil {
		t.Fatalf("CheckCalibrationStatus 실패: %v", err)
	}
	if st != CalibrationStatusExpired {
		t.Errorf("Status: got %d, want %d (Expired)", st, CalibrationStatusExpired)
	}
	if msg == "" {
		t.Error("만료 메시지가 비어 있습니다")
	}
}

func TestCheckCalibrationStatus_Expiring(t *testing.T) {
	svc, calRepo, _ := newTestCalibrationService()
	ctx := context.Background()

	// 5일 후 만료되는 보정 기록 삽입 (7일 이내 = EXPIRING)
	expiringRecord := &CalibrationRecord{
		ID:                 "expiring-cal-001",
		DeviceID:           "device-032",
		CartridgeCategory:  1,
		CartridgeTypeIndex: 1,
		CalibrationType:    CalibrationTypeFactory,
		Alpha:              0.95,
		CalibratedAt:       time.Now().UTC().AddDate(0, 0, -85),
		ExpiresAt:          time.Now().UTC().AddDate(0, 0, 5), // 5일 후 만료
		Status:             CalibrationStatusValid,
	}
	_ = calRepo.Save(ctx, expiringRecord)

	st, msg, _, err := svc.CheckCalibrationStatus(ctx, "device-032", 1, 1)
	if err != nil {
		t.Fatalf("CheckCalibrationStatus 실패: %v", err)
	}
	if st != CalibrationStatusExpiring {
		t.Errorf("Status: got %d, want %d (Expiring)", st, CalibrationStatusExpiring)
	}
	if msg == "" {
		t.Error("만료 임박 메시지가 비어 있습니다")
	}
}

func TestCheckCalibrationStatus_NoCalibration(t *testing.T) {
	svc, _, _ := newTestCalibrationService()
	ctx := context.Background()

	st, msg, latest, err := svc.CheckCalibrationStatus(ctx, "device-no-cal", 1, 1)
	if err != nil {
		t.Fatalf("CheckCalibrationStatus 실패: %v", err)
	}
	if st != CalibrationStatusNeeded {
		t.Errorf("Status: got %d, want %d (Needed)", st, CalibrationStatusNeeded)
	}
	if msg == "" {
		t.Error("보정 필요 메시지가 비어 있습니다")
	}
	if latest != nil {
		t.Error("보정 기록이 없어야 합니다")
	}
}

func TestListCalibrationModels(t *testing.T) {
	svc, _, _ := newTestCalibrationService()
	ctx := context.Background()

	models, err := svc.ListCalibrationModels(ctx)
	if err != nil {
		t.Fatalf("ListCalibrationModels 실패: %v", err)
	}

	if len(models) == 0 {
		t.Fatal("보정 모델이 하나도 없습니다")
	}

	// 시드 데이터 검증
	expectedCount := 8 // 3 (Health) + 2 (Electronic) + 2 (Advanced) + 1 (Custom)
	if len(models) != expectedCount {
		t.Errorf("모델 수: got %d, want %d", len(models), expectedCount)
	}

	// HealthBiomarker 모델 확인 (카테고리 1)
	healthModels := 0
	for _, m := range models {
		if m.CartridgeCategory == 1 {
			healthModels++
			if m.DefaultAlpha != 0.95 {
				t.Errorf("HealthBiomarker DefaultAlpha: got %f, want 0.95", m.DefaultAlpha)
			}
			if m.ValidityDays != 90 {
				t.Errorf("HealthBiomarker ValidityDays: got %d, want 90", m.ValidityDays)
			}
		}
	}
	if healthModels != 3 {
		t.Errorf("HealthBiomarker 모델 수: got %d, want 3", healthModels)
	}

	// ElectronicSensor 모델 확인 (카테고리 4)
	for _, m := range models {
		if m.CartridgeCategory == 4 {
			if m.DefaultAlpha != 0.92 {
				t.Errorf("ElectronicSensor DefaultAlpha: got %f, want 0.92", m.DefaultAlpha)
			}
			if m.ValidityDays != 60 {
				t.Errorf("ElectronicSensor ValidityDays: got %d, want 60", m.ValidityDays)
			}
		}
	}

	// AdvancedAnalysis 모델 확인 (카테고리 5)
	for _, m := range models {
		if m.CartridgeCategory == 5 {
			if m.DefaultAlpha != 0.97 {
				t.Errorf("AdvancedAnalysis DefaultAlpha: got %f, want 0.97", m.DefaultAlpha)
			}
			if m.ValidityDays != 120 {
				t.Errorf("AdvancedAnalysis ValidityDays: got %d, want 120", m.ValidityDays)
			}
		}
	}

	// CustomResearch 모델 확인 (카테고리 254)
	for _, m := range models {
		if m.CartridgeCategory == 254 {
			if m.DefaultAlpha != 0.95 {
				t.Errorf("CustomResearch DefaultAlpha: got %f, want 0.95", m.DefaultAlpha)
			}
			if m.ValidityDays != 30 {
				t.Errorf("CustomResearch ValidityDays: got %d, want 30", m.ValidityDays)
			}
		}
	}
}
