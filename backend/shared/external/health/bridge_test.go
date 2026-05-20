package health_test

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/manpasik/backend/shared/external/health"
)

// ============================================================================
// In-memory mock repository
// ============================================================================

type mockRepo struct {
	mu      sync.Mutex
	samples map[string]*health.HealthSample
}

func newMockRepo() *mockRepo {
	return &mockRepo{samples: make(map[string]*health.HealthSample)}
}

func (m *mockRepo) FindBySourceID(userID, source, sourceID string) (*health.HealthSample, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	key := userID + "|" + source + "|" + sourceID
	s, ok := m.samples[key]
	if !ok {
		return nil, nil
	}
	return s, nil
}

func (m *mockRepo) FindByTimeRange(userID string, dataType health.DataType, start, end time.Time) ([]*health.HealthSample, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var out []*health.HealthSample
	for _, s := range m.samples {
		if s.UserID == userID && s.Type == dataType &&
			!s.Timestamp.Before(start) && !s.Timestamp.After(end) {
			out = append(out, s)
		}
	}
	return out, nil
}

func (m *mockRepo) Save(s *health.HealthSample) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if s.SourceID == "" {
		return "", errors.New("source_id required for mock")
	}
	key := s.UserID + "|" + s.Source + "|" + s.SourceID
	m.samples[key] = s
	return s.SourceID, nil
}

func (m *mockRepo) Update(id string, s *health.HealthSample) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for k, existing := range m.samples {
		if existing.SourceID == id {
			m.samples[k] = s
			return nil
		}
	}
	return errors.New("not found")
}

// ============================================================================
// 테스트
// ============================================================================

// TestNormalize_BloodGlucose는 mmol/L → mg/dL 변환을 검증합니다.
func TestNormalize_BloodGlucose(t *testing.T) {
	n := health.NewSampleNormalizer()
	s := &health.HealthSample{
		UserID:    "u1",
		Type:      health.TypeBloodGlucose,
		Value:     5.5, // mmol/L
		Unit:      "mmol/L",
		Source:    health.SourceAppleHealthKit,
		Timestamp: time.Now().UTC(),
	}

	normalized, err := n.Normalize(s)
	if err != nil {
		t.Fatalf("Normalize 실패: %v", err)
	}
	if normalized.Unit != "mg/dL" {
		t.Errorf("Unit = %q, want mg/dL", normalized.Unit)
	}
	expected := 5.5 * 18.0182
	if normalized.Value < expected-0.1 || normalized.Value > expected+0.1 {
		t.Errorf("Value = %f, want ~%f", normalized.Value, expected)
	}
}

// TestNormalize_Weight는 lb → kg 변환을 검증합니다.
func TestNormalize_Weight(t *testing.T) {
	n := health.NewSampleNormalizer()
	s := &health.HealthSample{
		UserID: "u1", Type: health.TypeWeight,
		Value: 154.0, Unit: "lb",
		Source: health.SourceAppleHealthKit,
		Timestamp: time.Now().UTC(),
	}
	normalized, _ := n.Normalize(s)
	if normalized.Unit != "kg" {
		t.Errorf("Unit = %q, want kg", normalized.Unit)
	}
	if normalized.Value < 69 || normalized.Value > 71 {
		t.Errorf("Value = %f, want ~70", normalized.Value)
	}
}

// TestNormalize_Temperature는 °F → °C 변환을 검증합니다.
func TestNormalize_Temperature(t *testing.T) {
	n := health.NewSampleNormalizer()
	s := &health.HealthSample{
		UserID: "u1", Type: health.TypeBodyTemperature,
		Value: 98.6, Unit: "°F",
		Source: health.SourceAppleHealthKit,
		Timestamp: time.Now().UTC(),
	}
	normalized, _ := n.Normalize(s)
	if normalized.Unit != "°C" {
		t.Errorf("Unit = %q, want °C", normalized.Unit)
	}
	if normalized.Value < 36.9 || normalized.Value > 37.1 {
		t.Errorf("Value = %f, want ~37.0", normalized.Value)
	}
}

// TestValidateSample_OutOfRange는 범위 검증을 확인합니다.
func TestValidateSample_OutOfRange(t *testing.T) {
	cases := []*health.HealthSample{
		{UserID: "u", Type: health.TypeHeartRate, Value: 250, Unit: "bpm", Source: "x", Timestamp: time.Now()},
		{UserID: "u", Type: health.TypeBloodGlucose, Value: 1000, Unit: "mg/dL", Source: "x", Timestamp: time.Now()},
		{UserID: "u", Type: health.TypeOxygenSaturation, Value: 30, Unit: "%", Source: "x", Timestamp: time.Now()},
	}
	for i, c := range cases {
		if err := health.ValidateSample(c); err == nil {
			t.Errorf("case %d: 범위 초과 통과됨", i)
		}
	}
}

// TestValidateSample_Required는 필수 필드 검증을 확인합니다.
func TestValidateSample_Required(t *testing.T) {
	cases := []*health.HealthSample{
		{},
		{UserID: "u"},                                                    // type 누락
		{UserID: "u", Type: health.TypeSteps},                            // source 누락
		{UserID: "u", Type: health.TypeSteps, Source: "x"},              // timestamp 누락
	}
	for i, c := range cases {
		if err := health.ValidateSample(c); err == nil {
			t.Errorf("case %d: 필수 필드 누락 통과됨", i)
		}
	}
}

// TestConflictResolver_KeepExisting은 기존 유지 정책을 검증합니다.
func TestConflictResolver_KeepExisting(t *testing.T) {
	r := health.NewConflictResolver(health.PolicyKeepExisting)
	existing := &health.HealthSample{Value: 100, SourceID: "old"}
	incoming := &health.HealthSample{Value: 120, SourceID: "new"}

	resolved, conflict := r.Resolve(existing, incoming)
	if resolved != existing {
		t.Error("기존이 유지되지 않음")
	}
	if conflict.Resolution != "kept_existing" {
		t.Errorf("Resolution = %q", conflict.Resolution)
	}
}

// TestConflictResolver_ManpasikWins는 만파식 우선 정책을 검증합니다.
func TestConflictResolver_ManpasikWins(t *testing.T) {
	r := health.NewConflictResolver(health.PolicyManpasikWins)

	// 기존: 만파식 → 유지
	existing := &health.HealthSample{Source: health.SourceManPaSik, Value: 100, SourceID: "m"}
	incoming := &health.HealthSample{Source: health.SourceAppleHealthKit, Value: 120}

	resolved, _ := r.Resolve(existing, incoming)
	if resolved.Source != health.SourceManPaSik {
		t.Error("만파식 데이터가 유지되지 않음")
	}

	// 신규: 만파식 → 교체
	existing2 := &health.HealthSample{Source: health.SourceAppleHealthKit, Value: 100}
	incoming2 := &health.HealthSample{Source: health.SourceManPaSik, Value: 120}

	resolved2, _ := r.Resolve(existing2, incoming2)
	if resolved2.Source != health.SourceManPaSik {
		t.Error("만파식 데이터가 교체되지 않음")
	}
}

// TestConflictResolver_LatestWins는 최신 우선 정책을 검증합니다.
func TestConflictResolver_LatestWins(t *testing.T) {
	r := health.NewConflictResolver(health.PolicyLatestWins)
	old := time.Now().Add(-1 * time.Hour)
	new := time.Now()

	existing := &health.HealthSample{Timestamp: old, Value: 100}
	incoming := &health.HealthSample{Timestamp: new, Value: 120}

	resolved, _ := r.Resolve(existing, incoming)
	if resolved.Value != 120 {
		t.Errorf("Value = %f, want 120 (newer)", resolved.Value)
	}
}

// TestConflictResolver_Merge는 병합 정책을 검증합니다.
func TestConflictResolver_Merge(t *testing.T) {
	r := health.NewConflictResolver(health.PolicyMerge)
	existing := &health.HealthSample{Value: 100, SourceID: "e"}
	incoming := &health.HealthSample{Value: 120, SourceID: "i"}

	resolved, _ := r.Resolve(existing, incoming)
	if resolved.Value != 110 {
		t.Errorf("Merged value = %f, want 110", resolved.Value)
	}
	if resolved.Source != "merged" {
		t.Errorf("Source = %q, want merged", resolved.Source)
	}
}

// TestHealthBridge_Sync_HappyPath는 정상 동기화를 검증합니다.
func TestHealthBridge_Sync_HappyPath(t *testing.T) {
	repo := newMockRepo()
	bridge := health.NewHealthBridge(repo, health.PolicyManpasikWins)

	now := time.Now().UTC()
	batch := &health.SyncBatch{
		UserID: "u1",
		Source: health.SourceAppleHealthKit,
		Samples: []*health.HealthSample{
			{UserID: "u1", Type: health.TypeSteps, Value: 8000, Unit: "count", Source: health.SourceAppleHealthKit, SourceID: "ak-1", Timestamp: now},
			{UserID: "u1", Type: health.TypeHeartRate, Value: 72, Unit: "bpm", Source: health.SourceAppleHealthKit, SourceID: "ak-2", Timestamp: now},
		},
	}

	result, err := bridge.Sync(batch)
	if err != nil {
		t.Fatalf("Sync 실패: %v", err)
	}
	if result.ImportedCount != 2 {
		t.Errorf("ImportedCount = %d, want 2", result.ImportedCount)
	}
	if result.FailedCount != 0 {
		t.Errorf("FailedCount = %d, want 0", result.FailedCount)
	}
}

// TestHealthBridge_Sync_DuplicateSkip은 중복 스킵을 검증합니다.
func TestHealthBridge_Sync_DuplicateSkip(t *testing.T) {
	repo := newMockRepo()
	bridge := health.NewHealthBridge(repo, health.PolicyKeepExisting)

	now := time.Now().UTC()
	batch1 := &health.SyncBatch{
		UserID: "u1",
		Source: health.SourceAppleHealthKit,
		Samples: []*health.HealthSample{
			{UserID: "u1", Type: health.TypeSteps, Value: 8000, Unit: "count", Source: health.SourceAppleHealthKit, SourceID: "ak-dup", Timestamp: now},
		},
	}
	_, _ = bridge.Sync(batch1)

	// 동일 SourceID 재동기화
	result, _ := bridge.Sync(batch1)
	if result.ImportedCount != 0 {
		t.Errorf("재동기화 ImportedCount = %d, want 0", result.ImportedCount)
	}
	if result.DuplicateCount != 1 {
		t.Errorf("DuplicateCount = %d, want 1", result.DuplicateCount)
	}
}

// TestHealthBridge_Sync_ConflictReplaced는 충돌 시 교체를 검증합니다.
func TestHealthBridge_Sync_ConflictReplaced(t *testing.T) {
	repo := newMockRepo()
	bridge := health.NewHealthBridge(repo, health.PolicyReplaceNew)

	now := time.Now().UTC()
	// 기존 저장
	_ = mustSave(repo, &health.HealthSample{
		UserID: "u1", Type: health.TypeSteps, Value: 5000, Unit: "count",
		Source: health.SourceAppleHealthKit, SourceID: "dup-1", Timestamp: now,
	})

	// 동일 SourceID에 다른 값 재동기화
	batch := &health.SyncBatch{
		UserID: "u1", Source: health.SourceAppleHealthKit,
		Samples: []*health.HealthSample{
			{UserID: "u1", Type: health.TypeSteps, Value: 6000, Unit: "count",
				Source: health.SourceAppleHealthKit, SourceID: "dup-1", Timestamp: now},
		},
	}
	result, _ := bridge.Sync(batch)
	if result.ConflictCount != 1 {
		t.Errorf("ConflictCount = %d, want 1", result.ConflictCount)
	}
}

// TestHealthBridge_Sync_NormalizesUnits는 동기화 시 단위 정규화를 검증합니다.
func TestHealthBridge_Sync_NormalizesUnits(t *testing.T) {
	repo := newMockRepo()
	bridge := health.NewHealthBridge(repo, health.PolicyKeepExisting)

	now := time.Now().UTC()
	batch := &health.SyncBatch{
		UserID: "u1", Source: health.SourceAppleHealthKit,
		Samples: []*health.HealthSample{
			{UserID: "u1", Type: health.TypeBloodGlucose, Value: 5.5, Unit: "mmol/L",
				Source: health.SourceAppleHealthKit, SourceID: "g-1", Timestamp: now},
		},
	}
	_, _ = bridge.Sync(batch)

	saved, _ := repo.FindBySourceID("u1", health.SourceAppleHealthKit, "g-1")
	if saved == nil {
		t.Fatal("저장된 샘플 없음")
	}
	if saved.Unit != "mg/dL" {
		t.Errorf("Unit = %q, want mg/dL (정규화됨)", saved.Unit)
	}
}

// TestHealthBridge_Sync_InvalidSamplesFiltered는 잘못된 샘플 필터링을 검증합니다.
func TestHealthBridge_Sync_InvalidSamplesFiltered(t *testing.T) {
	repo := newMockRepo()
	bridge := health.NewHealthBridge(repo, health.PolicyKeepExisting)

	now := time.Now().UTC()
	batch := &health.SyncBatch{
		UserID: "u1", Source: "x",
		Samples: []*health.HealthSample{
			{UserID: "u1", Type: health.TypeHeartRate, Value: 72, Unit: "bpm", Source: "x", SourceID: "good", Timestamp: now},
			{UserID: "u1", Type: health.TypeHeartRate, Value: 999, Unit: "bpm", Source: "x", SourceID: "bad", Timestamp: now}, // 범위 초과
		},
	}
	result, _ := bridge.Sync(batch)
	if result.FailedCount != 1 {
		t.Errorf("FailedCount = %d, want 1", result.FailedCount)
	}
	if result.ImportedCount != 1 {
		t.Errorf("ImportedCount = %d, want 1", result.ImportedCount)
	}
}

// TestPreparePushBatch는 만파식 → 외부 푸시 배치 준비를 검증합니다.
func TestPreparePushBatch(t *testing.T) {
	repo := newMockRepo()
	bridge := health.NewHealthBridge(repo, health.PolicyManpasikWins)

	now := time.Now().UTC()
	_ = mustSave(repo, &health.HealthSample{
		UserID: "u1", Type: health.TypeBloodGlucose, Value: 110, Unit: "mg/dL",
		Source: health.SourceManPaSik, SourceID: "mp-1", Timestamp: now,
	})

	pushBatch, err := bridge.PreparePushBatch("u1", health.SourceAppleHealthKit, now.Add(-1*time.Hour))
	if err != nil {
		t.Fatalf("PreparePushBatch 실패: %v", err)
	}
	if pushBatch.Target != health.SourceAppleHealthKit {
		t.Errorf("Target = %q", pushBatch.Target)
	}
	if len(pushBatch.Samples) != 1 {
		t.Errorf("Samples = %d, want 1", len(pushBatch.Samples))
	}
}

// 헬퍼: repo에 저장
func mustSave(repo *mockRepo, s *health.HealthSample) string {
	id, _ := repo.Save(s)
	return id
}
