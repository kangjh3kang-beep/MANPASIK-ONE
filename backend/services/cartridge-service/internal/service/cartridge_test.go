package service

import (
	"context"
	"encoding/binary"
	"math"
	"testing"
	"time"

	"go.uber.org/zap"
)

// ============================================================================
// 테스트용 Fake 저장소
// ============================================================================

type fakeUsageRepo struct {
	records []*CartridgeUsageRecord
}

func newFakeUsageRepo() *fakeUsageRepo {
	return &fakeUsageRepo{
		records: make([]*CartridgeUsageRecord, 0),
	}
}

func (r *fakeUsageRepo) Create(_ context.Context, record *CartridgeUsageRecord) error {
	r.records = append(r.records, record)
	return nil
}

func (r *fakeUsageRepo) ListByUserID(_ context.Context, userID string, limit, offset int32) ([]*CartridgeUsageRecord, int32, error) {
	var filtered []*CartridgeUsageRecord
	for _, rec := range r.records {
		if rec.UserID == userID {
			filtered = append(filtered, rec)
		}
	}
	total := int32(len(filtered))
	start := int(offset)
	if start >= len(filtered) {
		return []*CartridgeUsageRecord{}, total, nil
	}
	end := start + int(limit)
	if end > len(filtered) {
		end = len(filtered)
	}
	return filtered[start:end], total, nil
}

type fakeStateRepo struct {
	byUID map[string]*CartridgeRemainingInfo
}

func newFakeStateRepo() *fakeStateRepo {
	return &fakeStateRepo{
		byUID: make(map[string]*CartridgeRemainingInfo),
	}
}

func (r *fakeStateRepo) GetByUID(_ context.Context, uid string) (*CartridgeRemainingInfo, error) {
	info, ok := r.byUID[uid]
	if !ok {
		return nil, nil
	}
	cp := *info
	return &cp, nil
}

func (r *fakeStateRepo) Upsert(_ context.Context, info *CartridgeRemainingInfo) error {
	cp := *info
	r.byUID[info.CartridgeUID] = &cp
	return nil
}

func (r *fakeStateRepo) DecrementUses(_ context.Context, uid string) (int32, error) {
	info, ok := r.byUID[uid]
	if !ok {
		return 0, nil
	}
	if info.RemainingUses > 0 {
		info.RemainingUses--
	}
	return info.RemainingUses, nil
}

func newTestService() (*CartridgeService, *fakeUsageRepo, *fakeStateRepo) {
	usageRepo := newFakeUsageRepo()
	stateRepo := newFakeStateRepo()
	svc := NewCartridgeService(zap.NewNop(), usageRepo, stateRepo)
	return svc, usageRepo, stateRepo
}

// ============================================================================
// NFC 태그 데이터 생성 헬퍼
// ============================================================================

func makeV1TagData(legacyCode byte, lotID, expiry string, remaining, maxUses uint16, alpha float64) []byte {
	data := make([]byte, 64)

	// UID (8바이트)
	copy(data[0:8], []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08})

	// 레거시 타입 코드
	data[8] = legacyCode

	// 로트 ID (8바이트, 패딩)
	lot := make([]byte, 8)
	copy(lot, []byte(lotID))
	copy(data[9:17], lot)

	// 유효 기간 (8바이트 YYYYMMDD)
	exp := make([]byte, 8)
	copy(exp, []byte(expiry))
	copy(data[17:25], exp)

	// 잔여 사용 횟수 (u16 LE)
	binary.LittleEndian.PutUint16(data[25:27], remaining)
	// 최대 사용 횟수 (u16 LE)
	binary.LittleEndian.PutUint16(data[27:29], maxUses)

	// α 계수 (f64 LE)
	binary.LittleEndian.PutUint64(data[29:37], math.Float64bits(alpha))
	// 온도 보정 계수 (f64 LE)
	binary.LittleEndian.PutUint64(data[37:45], math.Float64bits(0.02))
	// 습도 보정 계수 (f64 LE)
	binary.LittleEndian.PutUint64(data[45:53], math.Float64bits(0.01))

	return data
}

func makeV2TagData(categoryCode, typeIndex, legacyCode byte, lotID, expiry string, remaining, maxUses uint16, channels uint16, secs byte, alpha float64) []byte {
	data := make([]byte, 80)

	// UID (8바이트)
	copy(data[0:8], []byte{0xA1, 0xB2, 0xC3, 0xD4, 0xE5, 0xF6, 0x07, 0x18})

	data[8] = categoryCode
	data[9] = typeIndex
	data[10] = legacyCode
	data[11] = 0x02 // v2.0

	// 로트 ID
	lot := make([]byte, 8)
	copy(lot, []byte(lotID))
	copy(data[12:20], lot)

	// 유효 기간
	exp := make([]byte, 8)
	copy(exp, []byte(expiry))
	copy(data[20:28], exp)

	binary.LittleEndian.PutUint16(data[28:30], remaining)
	binary.LittleEndian.PutUint16(data[30:32], maxUses)

	binary.BigEndian.PutUint16(data[32:34], channels)
	data[34] = secs

	binary.LittleEndian.PutUint64(data[36:44], math.Float64bits(alpha))
	binary.LittleEndian.PutUint64(data[44:52], math.Float64bits(0.03))
	binary.LittleEndian.PutUint64(data[52:60], math.Float64bits(0.015))

	return data
}

// ============================================================================
// TestReadCartridgeV1 — v1.0 NFC 태그 파싱
// ============================================================================

func TestReadCartridgeV1(t *testing.T) {
	svc, _, _ := newTestService()
	ctx := context.Background()

	tagData := makeV1TagData(0x01, "LOT12345", "20271231", 50, 100, 0.95)

	detail, err := svc.ReadCartridge(ctx, tagData, 1)
	if err != nil {
		t.Fatalf("ReadCartridge v1 실패: %v", err)
	}

	if detail.CartridgeUID != "0102030405060708" {
		t.Errorf("UID: got %s, want 0102030405060708", detail.CartridgeUID)
	}
	if detail.CategoryCode != 1 {
		t.Errorf("CategoryCode: got %d, want 1", detail.CategoryCode)
	}
	if detail.TypeIndex != 1 {
		t.Errorf("TypeIndex: got %d, want 1", detail.TypeIndex)
	}
	if detail.LegacyCode != 0x01 {
		t.Errorf("LegacyCode: got %d, want 1", detail.LegacyCode)
	}
	if detail.NameKO != "혈당" {
		t.Errorf("NameKO: got %s, want 혈당", detail.NameKO)
	}
	if detail.NameEN != "Glucose" {
		t.Errorf("NameEN: got %s, want Glucose", detail.NameEN)
	}
	if detail.LotID != "LOT12345" {
		t.Errorf("LotID: got %s, want LOT12345", detail.LotID)
	}
	if detail.ExpiryDate != "20271231" {
		t.Errorf("ExpiryDate: got %s, want 20271231", detail.ExpiryDate)
	}
	if detail.RemainingUses != 50 {
		t.Errorf("RemainingUses: got %d, want 50", detail.RemainingUses)
	}
	if detail.MaxUses != 100 {
		t.Errorf("MaxUses: got %d, want 100", detail.MaxUses)
	}
	if math.Abs(detail.AlphaCoefficient-0.95) > 0.001 {
		t.Errorf("AlphaCoefficient: got %f, want 0.95", detail.AlphaCoefficient)
	}
	if detail.RequiredChannels != 88 {
		t.Errorf("RequiredChannels: got %d, want 88", detail.RequiredChannels)
	}
	if detail.MeasurementSecs != 15 {
		t.Errorf("MeasurementSecs: got %d, want 15", detail.MeasurementSecs)
	}
	if detail.Unit != "mg/dL" {
		t.Errorf("Unit: got %s, want mg/dL", detail.Unit)
	}
	if !detail.IsValid {
		t.Error("IsValid should be true")
	}
}

func TestReadCartridgeV1_Environmental(t *testing.T) {
	svc, _, _ := newTestService()
	ctx := context.Background()

	// WaterQuality = 0x20
	tagData := makeV1TagData(0x20, "ENVLOT01", "20280630", 30, 50, 0.92)

	detail, err := svc.ReadCartridge(ctx, tagData, 1)
	if err != nil {
		t.Fatalf("ReadCartridge v1 환경 실패: %v", err)
	}

	if detail.CategoryCode != 2 {
		t.Errorf("CategoryCode: got %d, want 2", detail.CategoryCode)
	}
	if detail.TypeIndex != 1 {
		t.Errorf("TypeIndex: got %d, want 1", detail.TypeIndex)
	}
	if detail.NameKO != "수질 검사" {
		t.Errorf("NameKO: got %s, want 수질 검사", detail.NameKO)
	}
}

func TestReadCartridgeV1_TooShort(t *testing.T) {
	svc, _, _ := newTestService()
	ctx := context.Background()

	_, err := svc.ReadCartridge(ctx, make([]byte, 20), 1)
	if err == nil {
		t.Error("짧은 v1 태그 데이터가 허용되었습니다")
	}
}

func TestReadCartridgeV1_Empty(t *testing.T) {
	svc, _, _ := newTestService()
	ctx := context.Background()

	_, err := svc.ReadCartridge(ctx, nil, 1)
	if err == nil {
		t.Error("빈 태그 데이터가 허용되었습니다")
	}
}

// ============================================================================
// TestReadCartridgeV2 — v2.0 NFC 태그 파싱
// ============================================================================

func TestReadCartridgeV2(t *testing.T) {
	svc, _, _ := newTestService()
	ctx := context.Background()

	// Industrial category (v2.0 전용, 레거시 없음)
	tagData := makeV2TagData(0x06, 0x01, 0x00, "LOT2XXXX", "20281231", 30, 50, 88, 30, 0.93)

	detail, err := svc.ReadCartridge(ctx, tagData, 2)
	if err != nil {
		t.Fatalf("ReadCartridge v2 실패: %v", err)
	}

	if detail.CartridgeUID != "a1b2c3d4e5f60718" {
		t.Errorf("UID: got %s, want a1b2c3d4e5f60718", detail.CartridgeUID)
	}
	if detail.CategoryCode != 6 {
		t.Errorf("CategoryCode: got %d, want 6", detail.CategoryCode)
	}
	if detail.TypeIndex != 1 {
		t.Errorf("TypeIndex: got %d, want 1", detail.TypeIndex)
	}
	if detail.LegacyCode != 0 {
		t.Errorf("LegacyCode: got %d, want 0", detail.LegacyCode)
	}
	if detail.RemainingUses != 30 {
		t.Errorf("RemainingUses: got %d, want 30", detail.RemainingUses)
	}
	if detail.MaxUses != 50 {
		t.Errorf("MaxUses: got %d, want 50", detail.MaxUses)
	}
	if detail.RequiredChannels != 88 {
		t.Errorf("RequiredChannels: got %d, want 88", detail.RequiredChannels)
	}
	if detail.MeasurementSecs != 30 {
		t.Errorf("MeasurementSecs: got %d, want 30", detail.MeasurementSecs)
	}
	if math.Abs(detail.AlphaCoefficient-0.93) > 0.001 {
		t.Errorf("AlphaCoefficient: got %f, want 0.93", detail.AlphaCoefficient)
	}
	if !detail.IsValid {
		t.Error("IsValid should be true")
	}
}

func TestReadCartridgeV2_WithLegacy(t *testing.T) {
	svc, _, _ := newTestService()
	ctx := context.Background()

	// HealthBiomarker Glucose with v2.0 format and legacy code
	tagData := makeV2TagData(0x01, 0x01, 0x01, "LOT2G001", "20271231", 40, 100, 88, 15, 0.95)

	detail, err := svc.ReadCartridge(ctx, tagData, 2)
	if err != nil {
		t.Fatalf("ReadCartridge v2 레거시 실패: %v", err)
	}

	if detail.CategoryCode != 1 {
		t.Errorf("CategoryCode: got %d, want 1", detail.CategoryCode)
	}
	if detail.LegacyCode != 1 {
		t.Errorf("LegacyCode: got %d, want 1", detail.LegacyCode)
	}
	if detail.NameKO != "혈당" {
		t.Errorf("NameKO: got %s, want 혈당", detail.NameKO)
	}
}

func TestReadCartridgeV2_NonTarget1792(t *testing.T) {
	svc, _, _ := newTestService()
	ctx := context.Background()

	// AdvancedAnalysis NonTarget1792 (cat=5, idx=3)
	tagData := makeV2TagData(0x05, 0x03, 0x52, "LOTN1792", "20291231", 10, 20, 1792, 180, 0.99)

	detail, err := svc.ReadCartridge(ctx, tagData, 2)
	if err != nil {
		t.Fatalf("ReadCartridge v2 NonTarget1792 실패: %v", err)
	}

	if detail.NameKO != "비표적 1792차원(궁극)" {
		t.Errorf("NameKO: got %s, want 비표적 1792차원(궁극)", detail.NameKO)
	}
	if detail.RequiredChannels != 1792 {
		t.Errorf("RequiredChannels: got %d, want 1792", detail.RequiredChannels)
	}
	if detail.MeasurementSecs != 180 {
		t.Errorf("MeasurementSecs: got %d, want 180", detail.MeasurementSecs)
	}
}

func TestReadCartridgeV2_TooShort(t *testing.T) {
	svc, _, _ := newTestService()
	ctx := context.Background()

	_, err := svc.ReadCartridge(ctx, make([]byte, 60), 2)
	if err == nil {
		t.Error("짧은 v2 태그 데이터가 허용되었습니다")
	}
}

// ============================================================================
// TestRecordUsage — 사용 기록 + 잔여 횟수 감소
// ============================================================================

func TestRecordUsage(t *testing.T) {
	svc, usageRepo, stateRepo := newTestService()
	ctx := context.Background()

	// 카트리지 상태 초기화
	stateRepo.byUID["cart-001"] = &CartridgeRemainingInfo{
		CartridgeUID:  "cart-001",
		RemainingUses: 10,
		MaxUses:       50,
		ExpiryDate:    "20271231",
	}

	remaining, err := svc.RecordUsage(ctx, "user-1", "session-1", "cart-001", 1, 1)
	if err != nil {
		t.Fatalf("RecordUsage 실패: %v", err)
	}

	if remaining != 9 {
		t.Errorf("잔여 횟수: got %d, want 9", remaining)
	}

	// 사용 기록 확인
	if len(usageRepo.records) != 1 {
		t.Fatalf("사용 기록 수: got %d, want 1", len(usageRepo.records))
	}
	record := usageRepo.records[0]
	if record.UserID != "user-1" {
		t.Errorf("UserID: got %s, want user-1", record.UserID)
	}
	if record.CartridgeUID != "cart-001" {
		t.Errorf("CartridgeUID: got %s, want cart-001", record.CartridgeUID)
	}
	if record.TypeNameKO != "혈당" {
		t.Errorf("TypeNameKO: got %s, want 혈당", record.TypeNameKO)
	}
	if record.RecordID == "" {
		t.Error("RecordID가 비어 있습니다")
	}

	// 추가 사용
	remaining2, err := svc.RecordUsage(ctx, "user-1", "session-2", "cart-001", 1, 1)
	if err != nil {
		t.Fatalf("RecordUsage 2차 실패: %v", err)
	}
	if remaining2 != 8 {
		t.Errorf("2차 잔여 횟수: got %d, want 8", remaining2)
	}
}

func TestRecordUsage_InvalidInput(t *testing.T) {
	svc, _, _ := newTestService()
	ctx := context.Background()

	_, err := svc.RecordUsage(ctx, "", "", "", 0, 0)
	if err == nil {
		t.Error("빈 입력으로 RecordUsage가 성공했습니다")
	}
}

// ============================================================================
// TestGetUsageHistory — 사용 이력 조회
// ============================================================================

func TestGetUsageHistory(t *testing.T) {
	svc, _, stateRepo := newTestService()
	ctx := context.Background()

	// 카트리지 상태 초기화
	stateRepo.byUID["cart-hist-001"] = &CartridgeRemainingInfo{
		CartridgeUID:  "cart-hist-001",
		RemainingUses: 20,
		MaxUses:       50,
	}

	// 3건의 사용 기록
	for i := 0; i < 3; i++ {
		_, err := svc.RecordUsage(ctx, "user-hist", "session-"+string(rune('A'+i)), "cart-hist-001", 1, 1)
		if err != nil {
			t.Fatalf("RecordUsage[%d] 실패: %v", i, err)
		}
	}

	// 전체 조회
	records, total, err := svc.GetUsageHistory(ctx, "user-hist", 10, 0)
	if err != nil {
		t.Fatalf("GetUsageHistory 실패: %v", err)
	}
	if total != 3 {
		t.Errorf("총 건수: got %d, want 3", total)
	}
	if len(records) != 3 {
		t.Errorf("반환 건수: got %d, want 3", len(records))
	}

	// 페이지네이션
	records2, total2, err := svc.GetUsageHistory(ctx, "user-hist", 2, 1)
	if err != nil {
		t.Fatalf("GetUsageHistory 페이지네이션 실패: %v", err)
	}
	if total2 != 3 {
		t.Errorf("총 건수: got %d, want 3", total2)
	}
	if len(records2) != 2 {
		t.Errorf("반환 건수: got %d, want 2", len(records2))
	}

	// 다른 사용자는 빈 결과
	records3, total3, err := svc.GetUsageHistory(ctx, "other-user", 10, 0)
	if err != nil {
		t.Fatalf("GetUsageHistory 다른 사용자 실패: %v", err)
	}
	if total3 != 0 {
		t.Errorf("다른 사용자 총 건수: got %d, want 0", total3)
	}
	if len(records3) != 0 {
		t.Errorf("다른 사용자 반환 건수: got %d, want 0", len(records3))
	}
}

func TestGetUsageHistory_InvalidInput(t *testing.T) {
	svc, _, _ := newTestService()
	ctx := context.Background()

	_, _, err := svc.GetUsageHistory(ctx, "", 10, 0)
	if err == nil {
		t.Error("빈 user_id로 GetUsageHistory가 성공했습니다")
	}
}

// ============================================================================
// TestGetCartridgeType — 타입 정보 조회
// ============================================================================

func TestGetCartridgeType(t *testing.T) {
	svc, _, _ := newTestService()
	ctx := context.Background()

	// Glucose (cat=1, idx=1)
	typeInfo, err := svc.GetCartridgeType(ctx, 1, 1)
	if err != nil {
		t.Fatalf("GetCartridgeType Glucose 실패: %v", err)
	}
	if typeInfo.NameEN != "Glucose" {
		t.Errorf("NameEN: got %s, want Glucose", typeInfo.NameEN)
	}
	if typeInfo.NameKO != "혈당" {
		t.Errorf("NameKO: got %s, want 혈당", typeInfo.NameKO)
	}
	if typeInfo.RequiredChannels != 88 {
		t.Errorf("RequiredChannels: got %d, want 88", typeInfo.RequiredChannels)
	}
	if typeInfo.MeasurementSecs != 15 {
		t.Errorf("MeasurementSecs: got %d, want 15", typeInfo.MeasurementSecs)
	}
	if typeInfo.Unit != "mg/dL" {
		t.Errorf("Unit: got %s, want mg/dL", typeInfo.Unit)
	}
	if typeInfo.Manufacturer != "ManPaSik" {
		t.Errorf("Manufacturer: got %s, want ManPaSik", typeInfo.Manufacturer)
	}

	// ENose (cat=4, idx=1)
	enose, err := svc.GetCartridgeType(ctx, 4, 1)
	if err != nil {
		t.Fatalf("GetCartridgeType ENose 실패: %v", err)
	}
	if enose.NameEN != "ENose" {
		t.Errorf("NameEN: got %s, want ENose", enose.NameEN)
	}
	if enose.RequiredChannels != 8 {
		t.Errorf("RequiredChannels: got %d, want 8", enose.RequiredChannels)
	}
	if enose.MeasurementSecs != 30 {
		t.Errorf("MeasurementSecs: got %d, want 30", enose.MeasurementSecs)
	}

	// NonTarget1792 (cat=5, idx=3)
	nt1792, err := svc.GetCartridgeType(ctx, 5, 3)
	if err != nil {
		t.Fatalf("GetCartridgeType NonTarget1792 실패: %v", err)
	}
	if nt1792.NameEN != "NonTarget1792" {
		t.Errorf("NameEN: got %s, want NonTarget1792", nt1792.NameEN)
	}
	if nt1792.RequiredChannels != 1792 {
		t.Errorf("RequiredChannels: got %d, want 1792", nt1792.RequiredChannels)
	}
	if nt1792.MeasurementSecs != 180 {
		t.Errorf("MeasurementSecs: got %d, want 180", nt1792.MeasurementSecs)
	}

	// CustomResearch (cat=255, idx=1)
	custom, err := svc.GetCartridgeType(ctx, 255, 1)
	if err != nil {
		t.Fatalf("GetCartridgeType CustomResearch 실패: %v", err)
	}
	if custom.NameEN != "CustomResearch" {
		t.Errorf("NameEN: got %s, want CustomResearch", custom.NameEN)
	}
}

func TestGetCartridgeType_NotFound(t *testing.T) {
	svc, _, _ := newTestService()
	ctx := context.Background()

	_, err := svc.GetCartridgeType(ctx, 99, 99)
	if err == nil {
		t.Error("존재하지 않는 타입 조회가 성공했습니다")
	}
}

// ============================================================================
// TestListCategories — 카테고리 목록 (15개)
// ============================================================================

func TestListCategories(t *testing.T) {
	svc, _, _ := newTestService()

	categories := svc.ListCategories()

	// 15개 카테고리: Unknown(0) + HealthBiomarker~Marine(1~12) + Beta(254) + CustomResearch(255)
	if len(categories) != 15 {
		t.Fatalf("카테고리 수: got %d, want 15", len(categories))
	}

	// 첫 번째는 Unknown
	if categories[0].Code != 0 {
		t.Errorf("첫 번째 카테고리 코드: got %d, want 0", categories[0].Code)
	}

	// 두 번째는 HealthBiomarker
	if categories[1].Code != 1 {
		t.Errorf("두 번째 카테고리 코드: got %d, want 1", categories[1].Code)
	}
	if categories[1].NameEN != "HealthBiomarker" {
		t.Errorf("두 번째 카테고리: got %s, want HealthBiomarker", categories[1].NameEN)
	}
	if categories[1].TypeCount != 14 {
		t.Errorf("HealthBiomarker 타입 수: got %d, want 14", categories[1].TypeCount)
	}
	if !categories[1].IsActive {
		t.Error("HealthBiomarker는 활성화 상태여야 합니다")
	}

	// AdvancedAnalysis (idx=5, code=5)
	if categories[5].Code != 5 {
		t.Errorf("6번째 카테고리 코드: got %d, want 5", categories[5].Code)
	}
	if categories[5].TypeCount != 4 {
		t.Errorf("AdvancedAnalysis 타입 수: got %d, want 4", categories[5].TypeCount)
	}

	// Industrial (idx=6, code=6)은 비활성
	if categories[6].Code != 6 {
		t.Errorf("7번째 카테고리 코드: got %d, want 6", categories[6].Code)
	}
	if categories[6].IsActive {
		t.Error("Industrial은 비활성 상태여야 합니다 (Phase 3-4)")
	}

	// Beta (code=254)
	betaIdx := len(categories) - 2
	if categories[betaIdx].Code != 254 {
		t.Errorf("Beta 카테고리 코드: got %d, want 254", categories[betaIdx].Code)
	}
	if categories[betaIdx].NameEN != "Beta" {
		t.Errorf("Beta 카테고리 이름: got %s, want Beta", categories[betaIdx].NameEN)
	}

	// CustomResearch (마지막, code=255)
	lastIdx := len(categories) - 1
	if categories[lastIdx].Code != 255 {
		t.Errorf("마지막 카테고리 코드: got %d, want 255", categories[lastIdx].Code)
	}
	if categories[lastIdx].TypeCount != 1 {
		t.Errorf("CustomResearch 타입 수: got %d, want 1", categories[lastIdx].TypeCount)
	}
}

// ============================================================================
// TestListTypesByCategory — 카테고리별 타입 목록
// ============================================================================

func TestListTypesByCategory(t *testing.T) {
	svc, _, _ := newTestService()
	ctx := context.Background()

	// HealthBiomarker: 14종
	types, err := svc.ListTypesByCategory(ctx, 1)
	if err != nil {
		t.Fatalf("ListTypesByCategory HealthBiomarker 실패: %v", err)
	}
	if len(types) != 14 {
		t.Errorf("HealthBiomarker 타입 수: got %d, want 14", len(types))
	}

	// ElectronicSensor: 3종
	sensorTypes, err := svc.ListTypesByCategory(ctx, 4)
	if err != nil {
		t.Fatalf("ListTypesByCategory ElectronicSensor 실패: %v", err)
	}
	if len(sensorTypes) != 3 {
		t.Errorf("ElectronicSensor 타입 수: got %d, want 3", len(sensorTypes))
	}

	// AdvancedAnalysis: 4종
	advTypes, err := svc.ListTypesByCategory(ctx, 5)
	if err != nil {
		t.Fatalf("ListTypesByCategory AdvancedAnalysis 실패: %v", err)
	}
	if len(advTypes) != 4 {
		t.Errorf("AdvancedAnalysis 타입 수: got %d, want 4", len(advTypes))
	}

	// 빈 카테고리
	_, err = svc.ListTypesByCategory(ctx, 99)
	if err == nil {
		t.Error("빈 카테고리 조회가 성공했습니다")
	}
}

// ============================================================================
// TestValidateCartridge — 유효성 검증
// ============================================================================

func TestValidateCartridge(t *testing.T) {
	svc, _, stateRepo := newTestService()
	ctx := context.Background()

	// 유효한 카트리지
	stateRepo.byUID["cart-valid"] = &CartridgeRemainingInfo{
		CartridgeUID:  "cart-valid",
		RemainingUses: 10,
		MaxUses:       50,
		ExpiryDate:    "20291231",
	}

	result, err := svc.ValidateCartridge(ctx, "cart-valid", 1, 1, "user-1")
	if err != nil {
		t.Fatalf("ValidateCartridge 유효 카트리지 실패: %v", err)
	}
	if !result.IsValid {
		t.Errorf("유효한 카트리지가 무효로 판정: reason=%s", result.Reason)
	}
	if result.Reason != "ok" {
		t.Errorf("Reason: got %s, want ok", result.Reason)
	}
	if result.RemainingUses != 10 {
		t.Errorf("RemainingUses: got %d, want 10", result.RemainingUses)
	}

	// 만료된 카트리지
	stateRepo.byUID["cart-expired"] = &CartridgeRemainingInfo{
		CartridgeUID:  "cart-expired",
		RemainingUses: 10,
		MaxUses:       50,
		ExpiryDate:    "20200101",
	}

	result2, err := svc.ValidateCartridge(ctx, "cart-expired", 1, 1, "user-1")
	if err != nil {
		t.Fatalf("ValidateCartridge 만료 카트리지 실패: %v", err)
	}
	if result2.IsValid {
		t.Error("만료된 카트리지가 유효로 판정되었습니다")
	}
	if result2.Reason != "expired" {
		t.Errorf("Reason: got %s, want expired", result2.Reason)
	}

	// 사용 횟수 소진
	stateRepo.byUID["cart-no-uses"] = &CartridgeRemainingInfo{
		CartridgeUID:  "cart-no-uses",
		RemainingUses: 0,
		MaxUses:       50,
		ExpiryDate:    "20291231",
	}

	result3, err := svc.ValidateCartridge(ctx, "cart-no-uses", 1, 1, "user-1")
	if err != nil {
		t.Fatalf("ValidateCartridge 사용 소진 실패: %v", err)
	}
	if result3.IsValid {
		t.Error("사용 소진 카트리지가 유효로 판정되었습니다")
	}
	if result3.Reason != "no_uses" {
		t.Errorf("Reason: got %s, want no_uses", result3.Reason)
	}

	// 알 수 없는 타입
	result4, err := svc.ValidateCartridge(ctx, "cart-valid", 99, 99, "user-1")
	if err != nil {
		t.Fatalf("ValidateCartridge 알 수 없는 타입 실패: %v", err)
	}
	if result4.IsValid {
		t.Error("알 수 없는 타입이 유효로 판정되었습니다")
	}
	if result4.Reason != "unknown_type" {
		t.Errorf("Reason: got %s, want unknown_type", result4.Reason)
	}
}

func TestValidateCartridge_NewCartridge(t *testing.T) {
	svc, _, _ := newTestService()
	ctx := context.Background()

	// 상태 정보가 없는 카트리지 (아직 ReadCartridge 전)
	result, err := svc.ValidateCartridge(ctx, "new-cart", 1, 1, "user-1")
	if err != nil {
		t.Fatalf("ValidateCartridge 새 카트리지 실패: %v", err)
	}
	if !result.IsValid {
		t.Error("새 카트리지(상태 없음)가 무효로 판정되었습니다")
	}
}

func TestValidateCartridge_EmptyUID(t *testing.T) {
	svc, _, _ := newTestService()
	ctx := context.Background()

	_, err := svc.ValidateCartridge(ctx, "", 1, 1, "user-1")
	if err == nil {
		t.Error("빈 UID로 ValidateCartridge가 성공했습니다")
	}
}

// ============================================================================
// TestGetRemainingUses — 잔여 사용 횟수 조회
// ============================================================================

func TestGetRemainingUses(t *testing.T) {
	svc, _, stateRepo := newTestService()
	ctx := context.Background()

	stateRepo.byUID["cart-rem-001"] = &CartridgeRemainingInfo{
		CartridgeUID:  "cart-rem-001",
		RemainingUses: 25,
		MaxUses:       100,
		ExpiryDate:    "20291231",
	}

	info, err := svc.GetRemainingUses(ctx, "cart-rem-001")
	if err != nil {
		t.Fatalf("GetRemainingUses 실패: %v", err)
	}
	if info.RemainingUses != 25 {
		t.Errorf("RemainingUses: got %d, want 25", info.RemainingUses)
	}
	if info.MaxUses != 100 {
		t.Errorf("MaxUses: got %d, want 100", info.MaxUses)
	}
	if info.IsExpired {
		t.Error("만료되지 않은 카트리지가 만료로 표시되었습니다")
	}
}

func TestGetRemainingUses_Expired(t *testing.T) {
	svc, _, stateRepo := newTestService()
	ctx := context.Background()

	stateRepo.byUID["cart-exp-001"] = &CartridgeRemainingInfo{
		CartridgeUID:  "cart-exp-001",
		RemainingUses: 10,
		MaxUses:       50,
		ExpiryDate:    "20200101",
	}

	info, err := svc.GetRemainingUses(ctx, "cart-exp-001")
	if err != nil {
		t.Fatalf("GetRemainingUses 만료 실패: %v", err)
	}
	if !info.IsExpired {
		t.Error("만료된 카트리지가 유효로 표시되었습니다")
	}
}

func TestGetRemainingUses_NotFound(t *testing.T) {
	svc, _, _ := newTestService()
	ctx := context.Background()

	_, err := svc.GetRemainingUses(ctx, "non-existent")
	if err == nil {
		t.Error("존재하지 않는 카트리지 잔여 조회가 성공했습니다")
	}
}

func TestGetRemainingUses_EmptyUID(t *testing.T) {
	svc, _, _ := newTestService()
	ctx := context.Background()

	_, err := svc.GetRemainingUses(ctx, "")
	if err == nil {
		t.Error("빈 UID로 GetRemainingUses가 성공했습니다")
	}
}

// ============================================================================
// 레지스트리 전체 확인
// ============================================================================

func TestRegistryTotalCount(t *testing.T) {
	svc, _, _ := newTestService()

	totalTypes := len(svc.typeRegistry)
	// 14(Health) + 4(Env) + 4(Food) + 3(Sensor) + 4(Advanced: 448+896+1792+MultiBio) + 1(Custom) = 30
	if totalTypes != 30 {
		t.Errorf("레지스트리 총 타입 수: got %d, want 30", totalTypes)
	}

	totalCategories := len(svc.categoryRegistry)
	if totalCategories != 15 {
		t.Errorf("레지스트리 총 카테고리 수: got %d, want 15", totalCategories)
	}
}

func TestLegacyToFullCode(t *testing.T) {
	tests := []struct {
		legacy   int32
		wantCat  int32
		wantType int32
	}{
		{0x01, 1, 1},   // Glucose
		{0x0E, 1, 14},  // Insulin
		{0x20, 2, 1},   // WaterQuality
		{0x23, 2, 4},   // Radiation
		{0x30, 3, 1},   // PesticideResidue
		{0x33, 3, 4},   // DateDrug
		{0x40, 4, 1},   // ENose
		{0x42, 4, 3},   // EhdGas
		{0x50, 5, 1},   // NonTarget448
		{0x52, 5, 3},   // NonTarget1792
		{0xFF, 255, 1}, // CustomResearch
	}

	for _, tt := range tests {
		cat, idx := legacyToFullCode(tt.legacy)
		if cat != tt.wantCat || idx != tt.wantType {
			t.Errorf("legacyToFullCode(0x%02X): got (%d,%d), want (%d,%d)",
				tt.legacy, cat, idx, tt.wantCat, tt.wantType)
		}
	}
}

// ============================================================================
// InitCartridgeState 확인
// ============================================================================

func TestInitCartridgeState(t *testing.T) {
	svc, _, stateRepo := newTestService()
	ctx := context.Background()

	detail := &CartridgeDetail{
		CartridgeUID:  "cart-init-001",
		RemainingUses: 50,
		MaxUses:       100,
		ExpiryDate:    "20291231",
	}

	err := svc.InitCartridgeState(ctx, detail)
	if err != nil {
		t.Fatalf("InitCartridgeState 실패: %v", err)
	}

	state, ok := stateRepo.byUID["cart-init-001"]
	if !ok {
		t.Fatal("카트리지 상태가 저장되지 않았습니다")
	}
	if state.RemainingUses != 50 {
		t.Errorf("RemainingUses: got %d, want 50", state.RemainingUses)
	}
	if state.MaxUses != 100 {
		t.Errorf("MaxUses: got %d, want 100", state.MaxUses)
	}
}

// ============================================================================
// 전체 플로우 테스트: ReadCartridge → InitState → RecordUsage → GetRemaining
// ============================================================================

func TestFullFlow(t *testing.T) {
	svc, _, _ := newTestService()
	ctx := context.Background()

	// 1. NFC 태그 읽기 (v1.0 Glucose)
	tagData := makeV1TagData(0x01, "LOT00001", "20291231", 5, 10, 0.95)
	detail, err := svc.ReadCartridge(ctx, tagData, 1)
	if err != nil {
		t.Fatalf("ReadCartridge 실패: %v", err)
	}

	// 2. 상태 초기화
	err = svc.InitCartridgeState(ctx, detail)
	if err != nil {
		t.Fatalf("InitCartridgeState 실패: %v", err)
	}

	// 3. 유효성 검증
	result, err := svc.ValidateCartridge(ctx, detail.CartridgeUID, detail.CategoryCode, detail.TypeIndex, "user-flow")
	if err != nil {
		t.Fatalf("ValidateCartridge 실패: %v", err)
	}
	if !result.IsValid {
		t.Fatalf("유효한 카트리지가 무효로 판정: %s", result.Reason)
	}

	// 4. 사용 기록 3회
	for i := 0; i < 3; i++ {
		_, err = svc.RecordUsage(ctx, "user-flow", "session-flow", detail.CartridgeUID, detail.CategoryCode, detail.TypeIndex)
		if err != nil {
			t.Fatalf("RecordUsage[%d] 실패: %v", i, err)
		}
	}

	// 5. 잔여 횟수 확인 (5 - 3 = 2)
	info, err := svc.GetRemainingUses(ctx, detail.CartridgeUID)
	if err != nil {
		t.Fatalf("GetRemainingUses 실패: %v", err)
	}
	if info.RemainingUses != 2 {
		t.Errorf("잔여 횟수: got %d, want 2", info.RemainingUses)
	}

	// 6. 사용 이력 확인
	records, total, err := svc.GetUsageHistory(ctx, "user-flow", 10, 0)
	if err != nil {
		t.Fatalf("GetUsageHistory 실패: %v", err)
	}
	if total != 3 {
		t.Errorf("이력 총 건수: got %d, want 3", total)
	}
	if len(records) != 3 {
		t.Errorf("이력 반환 건수: got %d, want 3", len(records))
	}

	// 7. 타임스탬프 확인
	for _, rec := range records {
		if rec.UsedAt.IsZero() {
			t.Error("UsedAt이 설정되지 않았습니다")
		}
		if time.Since(rec.UsedAt) > 5*time.Second {
			t.Error("UsedAt이 너무 오래되었습니다")
		}
	}
}
