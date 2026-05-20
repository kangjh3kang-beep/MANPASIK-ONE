package integration_test

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/manpasik/backend/shared/external/health"
)

// inMemoryRepo는 health 동기화 테스트용 인메모리 저장소입니다.
type inMemoryRepo struct {
	mu      sync.Mutex
	samples map[string]*health.HealthSample
}

func newInMemoryRepo() *inMemoryRepo {
	return &inMemoryRepo{samples: make(map[string]*health.HealthSample)}
}

func (m *inMemoryRepo) FindBySourceID(userID, source, sourceID string) (*health.HealthSample, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	key := userID + "|" + source + "|" + sourceID
	s, ok := m.samples[key]
	if !ok {
		return nil, nil
	}
	return s, nil
}

func (m *inMemoryRepo) FindByTimeRange(userID string, dataType health.DataType, start, end time.Time) ([]*health.HealthSample, error) {
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

func (m *inMemoryRepo) Save(s *health.HealthSample) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if s.SourceID == "" {
		return "", errors.New("source_id required")
	}
	key := s.UserID + "|" + s.Source + "|" + s.SourceID
	m.samples[key] = s
	return s.SourceID, nil
}

func (m *inMemoryRepo) Update(id string, s *health.HealthSample) error {
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

// TestHealthSyncFlow_AppleHealthKit_DailySync는 일일 HealthKit 동기화 시나리오입니다.
func TestHealthSyncFlow_AppleHealthKit_DailySync(t *testing.T) {
	repo := newInMemoryRepo()
	bridge := health.NewHealthBridge(repo, health.PolicyManpasikWins)

	now := time.Now().UTC()
	batch := &health.SyncBatch{
		UserID: "user-001",
		Source: health.SourceAppleHealthKit,
		Samples: []*health.HealthSample{
			{UserID: "user-001", Type: health.TypeSteps, Value: 8500, Unit: "count",
				Source: health.SourceAppleHealthKit, SourceID: "ak-step-1", Timestamp: now.Add(-3 * time.Hour)},
			{UserID: "user-001", Type: health.TypeHeartRate, Value: 72, Unit: "bpm",
				Source: health.SourceAppleHealthKit, SourceID: "ak-hr-1", Timestamp: now.Add(-2 * time.Hour)},
			{UserID: "user-001", Type: health.TypeBloodGlucose, Value: 5.5, Unit: "mmol/L",
				Source: health.SourceAppleHealthKit, SourceID: "ak-glu-1", Timestamp: now.Add(-1 * time.Hour)},
		},
	}

	result, err := bridge.Sync(batch)
	if err != nil {
		t.Fatalf("Sync 실패: %v", err)
	}
	if result.ImportedCount != 3 {
		t.Errorf("ImportedCount = %d, want 3", result.ImportedCount)
	}

	// 단위 정규화 검증: 5.5 mmol/L → ~99.1 mg/dL
	saved, _ := repo.FindBySourceID("user-001", health.SourceAppleHealthKit, "ak-glu-1")
	if saved == nil {
		t.Fatal("저장된 샘플 없음")
	}
	if saved.Unit != "mg/dL" {
		t.Errorf("Unit = %q, want mg/dL", saved.Unit)
	}
}

// TestHealthSyncFlow_GoogleFit_BatchUpload은 Google Fit 배치 업로드를 검증합니다.
func TestHealthSyncFlow_GoogleFit_BatchUpload(t *testing.T) {
	repo := newInMemoryRepo()
	bridge := health.NewHealthBridge(repo, health.PolicyLatestWins)

	now := time.Now().UTC()
	samples := make([]*health.HealthSample, 24)
	for i := 0; i < 24; i++ {
		samples[i] = &health.HealthSample{
			UserID: "user-002",
			Type:   health.TypeHeartRate,
			Value:  float64(70 + i),
			Unit:   "bpm",
			Source: health.SourceGoogleHealthConnect,
			SourceID: "gf-hr-" + string(rune('a'+i)),
			Timestamp: now.Add(time.Duration(-i) * time.Hour),
		}
	}

	batch := &health.SyncBatch{
		UserID:  "user-002",
		Source:  health.SourceGoogleHealthConnect,
		Samples: samples,
	}

	result, err := bridge.Sync(batch)
	if err != nil {
		t.Fatalf("Sync 실패: %v", err)
	}
	if result.ImportedCount != 24 {
		t.Errorf("ImportedCount = %d, want 24", result.ImportedCount)
	}
}

// TestHealthSyncFlow_DuplicateSync는 동일 source로 중복 동기화 시 스킵을 검증합니다.
func TestHealthSyncFlow_DuplicateSync(t *testing.T) {
	repo := newInMemoryRepo()
	bridge := health.NewHealthBridge(repo, health.PolicyKeepExisting)

	now := time.Now().UTC()

	batch := &health.SyncBatch{
		UserID: "u3",
		Source: health.SourceAppleHealthKit,
		Samples: []*health.HealthSample{
			{UserID: "u3", Type: health.TypeBloodGlucose,
				Value: 130, Unit: "mg/dL",
				Source: health.SourceAppleHealthKit,
				SourceID: "duplicate-id",
				Timestamp: now},
		},
	}

	// 첫 동기화 → 신규 저장
	r1, _ := bridge.Sync(batch)
	if r1.ImportedCount != 1 {
		t.Errorf("first ImportedCount = %d, want 1", r1.ImportedCount)
	}

	// 두 번째 동기화 → 중복 스킵
	r2, _ := bridge.Sync(batch)
	if r2.DuplicateCount != 1 {
		t.Errorf("DuplicateCount = %d, want 1", r2.DuplicateCount)
	}
	if r2.ImportedCount != 0 {
		t.Errorf("재동기화 ImportedCount = %d, want 0", r2.ImportedCount)
	}
}

// TestHealthSyncFlow_OutOfRangeFiltered는 범위 외 데이터 필터링을 검증합니다.
func TestHealthSyncFlow_OutOfRangeFiltered(t *testing.T) {
	repo := newInMemoryRepo()
	bridge := health.NewHealthBridge(repo, health.PolicyKeepExisting)

	now := time.Now().UTC()
	batch := &health.SyncBatch{
		UserID: "u4",
		Source: health.SourceAppleHealthKit,
		Samples: []*health.HealthSample{
			{UserID: "u4", Type: health.TypeHeartRate, Value: 250, Unit: "bpm",
				Source: health.SourceAppleHealthKit, SourceID: "bad-1", Timestamp: now}, // 범위 초과
			{UserID: "u4", Type: health.TypeHeartRate, Value: 75, Unit: "bpm",
				Source: health.SourceAppleHealthKit, SourceID: "good-1", Timestamp: now},
		},
	}

	result, _ := bridge.Sync(batch)
	if result.ImportedCount != 1 {
		t.Errorf("ImportedCount = %d, want 1", result.ImportedCount)
	}
	if result.FailedCount != 1 {
		t.Errorf("FailedCount = %d, want 1", result.FailedCount)
	}
}

// TestHealthSyncFlow_MultipleSources는 여러 소스의 동기화를 검증합니다.
func TestHealthSyncFlow_MultipleSources(t *testing.T) {
	repo := newInMemoryRepo()
	bridge := health.NewHealthBridge(repo, health.PolicyLatestWins)

	now := time.Now().UTC()
	apple := &health.SyncBatch{
		UserID: "u5",
		Source: health.SourceAppleHealthKit,
		Samples: []*health.HealthSample{
			{UserID: "u5", Type: health.TypeSteps, Value: 8000, Unit: "count",
				Source: health.SourceAppleHealthKit, SourceID: "apple-step", Timestamp: now},
		},
	}
	google := &health.SyncBatch{
		UserID: "u5",
		Source: health.SourceGoogleHealthConnect,
		Samples: []*health.HealthSample{
			{UserID: "u5", Type: health.TypeSteps, Value: 8200, Unit: "count",
				Source: health.SourceGoogleHealthConnect, SourceID: "google-step", Timestamp: now},
		},
	}

	r1, _ := bridge.Sync(apple)
	r2, _ := bridge.Sync(google)

	// 다른 source_id이므로 둘 다 저장됨
	if r1.ImportedCount != 1 {
		t.Errorf("Apple ImportedCount = %d", r1.ImportedCount)
	}
	if r2.ImportedCount != 1 {
		t.Errorf("Google ImportedCount = %d", r2.ImportedCount)
	}
}

// TestHealthSyncFlow_PreparePushBatch는 만파식 → 외부 푸시 배치 준비를 검증합니다.
func TestHealthSyncFlow_PreparePushBatch(t *testing.T) {
	repo := newInMemoryRepo()
	bridge := health.NewHealthBridge(repo, health.PolicyManpasikWins)

	now := time.Now().UTC()
	for i := 0; i < 5; i++ {
		_, _ = repo.Save(&health.HealthSample{
			UserID: "u6", Type: health.TypeBloodGlucose,
			Value: float64(100 + i*5), Unit: "mg/dL",
			Source: health.SourceManPaSik,
			SourceID: "mp-" + string(rune('a'+i)),
			Timestamp: now.Add(time.Duration(-i) * time.Hour),
		})
	}

	batch, err := bridge.PreparePushBatch("u6", health.SourceAppleHealthKit, now.Add(-24*time.Hour))
	if err != nil {
		t.Fatalf("PreparePushBatch 실패: %v", err)
	}
	if len(batch.Samples) != 5 {
		t.Errorf("Samples = %d, want 5", len(batch.Samples))
	}
	for _, s := range batch.Samples {
		if s.Source != health.SourceManPaSik {
			t.Errorf("Source = %q, want manpasik", s.Source)
		}
	}
}
