// Package health은 HealthKit/Google Fit 데이터의 백엔드 표준화 브릿지입니다.
//
// Flutter `health` 패키지에서 수신한 데이터를 만파식 표준 측정 모델로 변환하고,
// 양방향 동기화 + 충돌 해결을 제공합니다.
package health

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
)

// ============================================================================
// 도메인 모델
// ============================================================================

// Source는 데이터 소스 플랫폼입니다.
const (
	SourceAppleHealthKit       = "apple_healthkit"
	SourceGoogleHealthConnect  = "google_health_connect"
	SourceManPaSik             = "manpasik"
	SourceManualEntry          = "manual"
)

// DataType은 건강 데이터 유형입니다.
type DataType string

const (
	TypeSteps              DataType = "steps"
	TypeHeartRate          DataType = "heart_rate"
	TypeBloodPressureSys   DataType = "blood_pressure_systolic"
	TypeBloodPressureDia   DataType = "blood_pressure_diastolic"
	TypeBloodGlucose       DataType = "blood_glucose"
	TypeWeight             DataType = "weight"
	TypeSleep              DataType = "sleep"
	TypeOxygenSaturation   DataType = "oxygen_saturation"
	TypeBodyTemperature    DataType = "body_temperature"
	TypeRespiratoryRate    DataType = "respiratory_rate"
)

// HealthSample은 건강 데이터 샘플입니다.
type HealthSample struct {
	UserID    string
	Type      DataType
	Value     float64
	Unit      string    // "count", "bpm", "mmHg", "mg/dL", "kg" 등
	Source    string    // "apple_healthkit", "google_health_connect", "manpasik"
	SourceID  string    // 외부 시스템의 고유 ID (중복 동기화 방지)
	Timestamp time.Time // 측정 시각
	EndTime   *time.Time // 구간 데이터 (수면, 걸음수 시간대)
	DeviceID  string    // 측정 디바이스
	Metadata  map[string]string
}

// SyncBatch는 일괄 동기화 요청입니다.
type SyncBatch struct {
	UserID       string
	Source       string
	Samples      []*HealthSample
	StartTime    time.Time
	EndTime      time.Time
	RequestID    string
}

// SyncResult는 동기화 결과입니다.
type SyncResult struct {
	UserID         string
	Source         string
	ImportedCount  int
	DuplicateCount int  // 이미 동기화된 샘플
	ConflictCount  int  // 충돌 처리 필요
	FailedCount    int
	NewSampleIDs   []string
	Conflicts      []*Conflict
	SyncedAt       time.Time
}

// Conflict는 동기화 충돌입니다.
type Conflict struct {
	ExistingSampleID string
	NewSample        *HealthSample
	Resolution       string // "kept_existing", "replaced", "merged"
	Reason           string
}

// ============================================================================
// 데이터 정규화
// ============================================================================

// SampleNormalizer는 외부 플랫폼 데이터를 만파식 표준으로 정규화합니다.
type SampleNormalizer struct{}

// NewSampleNormalizer는 새 정규화자를 생성합니다.
func NewSampleNormalizer() *SampleNormalizer {
	return &SampleNormalizer{}
}

// Normalize는 샘플을 표준 단위로 정규화합니다.
func (n *SampleNormalizer) Normalize(s *HealthSample) (*HealthSample, error) {
	if s == nil {
		return nil, errors.New("sample is nil")
	}

	normalized := *s

	// 단위 표준화 (HealthKit/Google Fit 차이 처리)
	switch s.Type {
	case TypeBloodGlucose:
		// HealthKit: mg/dL 또는 mmol/L → mg/dL로 통일
		if strings.EqualFold(s.Unit, "mmol/L") {
			normalized.Value = s.Value * 18.0182
			normalized.Unit = "mg/dL"
		}
	case TypeWeight:
		// HealthKit: kg, lb → kg로 통일
		if strings.EqualFold(s.Unit, "lb") || strings.EqualFold(s.Unit, "lbs") {
			normalized.Value = s.Value * 0.453592
			normalized.Unit = "kg"
		}
	case TypeBodyTemperature:
		// °F → °C
		if strings.EqualFold(s.Unit, "°F") || strings.EqualFold(s.Unit, "F") {
			normalized.Value = (s.Value - 32) * 5 / 9
			normalized.Unit = "°C"
		}
	case TypeHeartRate:
		if normalized.Unit == "" {
			normalized.Unit = "bpm"
		}
	}

	// 타임스탬프 UTC 변환
	if !normalized.Timestamp.IsZero() {
		normalized.Timestamp = normalized.Timestamp.UTC()
	}

	return &normalized, nil
}

// NormalizeBatch는 배치 샘플을 정규화합니다.
func (n *SampleNormalizer) NormalizeBatch(samples []*HealthSample) ([]*HealthSample, []error) {
	results := make([]*HealthSample, 0, len(samples))
	errs := make([]error, 0)
	for i, s := range samples {
		norm, err := n.Normalize(s)
		if err != nil {
			errs = append(errs, fmt.Errorf("sample[%d]: %w", i, err))
			continue
		}
		results = append(results, norm)
	}
	return results, errs
}

// ============================================================================
// 충돌 해결 정책
// ============================================================================

// ConflictResolver는 동기화 충돌을 해결합니다.
type ConflictResolver struct {
	policy ConflictPolicy
}

// ConflictPolicy는 충돌 해결 정책입니다.
type ConflictPolicy string

const (
	PolicyKeepExisting ConflictPolicy = "keep_existing" // 기존 데이터 유지
	PolicyReplaceNew   ConflictPolicy = "replace_new"   // 새 데이터로 교체
	PolicyMerge        ConflictPolicy = "merge"          // 평균/병합
	PolicyManpasikWins ConflictPolicy = "manpasik_wins" // 만파식 측정이 우선
	PolicyLatestWins   ConflictPolicy = "latest_wins"   // 최신 timestamp 우선
)

// NewConflictResolver는 새 충돌 해결자를 생성합니다.
func NewConflictResolver(policy ConflictPolicy) *ConflictResolver {
	if policy == "" {
		policy = PolicyManpasikWins
	}
	return &ConflictResolver{policy: policy}
}

// Resolve는 두 샘플의 충돌을 해결하여 적용할 샘플을 반환합니다.
//
// 결과 reasoning은 Conflict.Resolution에 기록됩니다.
func (r *ConflictResolver) Resolve(existing, incoming *HealthSample) (*HealthSample, *Conflict) {
	if existing == nil {
		return incoming, nil
	}
	if incoming == nil {
		return existing, nil
	}

	conflict := &Conflict{
		ExistingSampleID: existing.SourceID,
		NewSample:        incoming,
	}

	switch r.policy {
	case PolicyKeepExisting:
		conflict.Resolution = "kept_existing"
		conflict.Reason = "정책: 기존 데이터 유지"
		return existing, conflict

	case PolicyReplaceNew:
		conflict.Resolution = "replaced"
		conflict.Reason = "정책: 새 데이터로 교체"
		return incoming, conflict

	case PolicyManpasikWins:
		if existing.Source == SourceManPaSik && incoming.Source != SourceManPaSik {
			conflict.Resolution = "kept_existing"
			conflict.Reason = "만파식 측정이 우선"
			return existing, conflict
		}
		if incoming.Source == SourceManPaSik && existing.Source != SourceManPaSik {
			conflict.Resolution = "replaced"
			conflict.Reason = "만파식 측정이 우선"
			return incoming, conflict
		}
		// 둘 다 같은 출처면 최신 우선
		fallthrough

	case PolicyLatestWins:
		if incoming.Timestamp.After(existing.Timestamp) {
			conflict.Resolution = "replaced"
			conflict.Reason = "더 최신 timestamp"
			return incoming, conflict
		}
		conflict.Resolution = "kept_existing"
		conflict.Reason = "기존 데이터가 더 최신"
		return existing, conflict

	case PolicyMerge:
		merged := *existing
		merged.Value = (existing.Value + incoming.Value) / 2
		merged.Source = "merged"
		conflict.Resolution = "merged"
		conflict.Reason = "두 값의 평균"
		return &merged, conflict
	}

	return existing, conflict
}

// ============================================================================
// 검증
// ============================================================================

// ValidateSample은 샘플을 검증합니다.
func ValidateSample(s *HealthSample) error {
	if s == nil {
		return errors.New("sample is nil")
	}
	if s.UserID == "" {
		return errors.New("user_id required")
	}
	if s.Type == "" {
		return errors.New("type required")
	}
	if s.Source == "" {
		return errors.New("source required")
	}
	if s.Timestamp.IsZero() {
		return errors.New("timestamp required")
	}
	if s.Value < 0 {
		return errors.New("value cannot be negative")
	}
	return validateRangeForType(s.Type, s.Value)
}

func validateRangeForType(t DataType, v float64) error {
	ranges := map[DataType][2]float64{
		TypeSteps:            {0, 100000},
		TypeHeartRate:        {30, 220},
		TypeBloodPressureSys: {60, 250},
		TypeBloodPressureDia: {30, 150},
		TypeBloodGlucose:     {20, 600},
		TypeWeight:           {10, 500},
		TypeOxygenSaturation: {50, 100},
		TypeBodyTemperature:  {30, 45},
		TypeRespiratoryRate:  {5, 50},
	}
	if r, ok := ranges[t]; ok {
		if v < r[0] || v > r[1] {
			return fmt.Errorf("value %f out of range [%f, %f] for %s", v, r[0], r[1], t)
		}
	}
	return nil
}

// ValidateBatch는 배치를 검증하고 잘못된 샘플 인덱스를 반환합니다.
func ValidateBatch(batch *SyncBatch) ([]int, error) {
	if batch == nil {
		return nil, errors.New("batch is nil")
	}
	if batch.UserID == "" {
		return nil, errors.New("user_id required")
	}

	var invalidIndexes []int
	for i, s := range batch.Samples {
		if err := ValidateSample(s); err != nil {
			invalidIndexes = append(invalidIndexes, i)
		}
	}
	return invalidIndexes, nil
}

// ============================================================================
// 동기화 브릿지
// ============================================================================

// SyncRepository는 동기화 데이터 저장소 인터페이스입니다.
type SyncRepository interface {
	FindBySourceID(userID, source, sourceID string) (*HealthSample, error)
	FindByTimeRange(userID string, dataType DataType, start, end time.Time) ([]*HealthSample, error)
	Save(sample *HealthSample) (string, error) // 저장 후 sample ID 반환
	Update(id string, sample *HealthSample) error
}

// HealthBridge는 외부 헬스 데이터 동기화 브릿지입니다.
type HealthBridge struct {
	normalizer *SampleNormalizer
	resolver   *ConflictResolver
	repo       SyncRepository
}

// NewHealthBridge는 새 동기화 브릿지를 생성합니다.
func NewHealthBridge(repo SyncRepository, policy ConflictPolicy) *HealthBridge {
	return &HealthBridge{
		normalizer: NewSampleNormalizer(),
		resolver:   NewConflictResolver(policy),
		repo:       repo,
	}
}

// Sync는 외부 플랫폼에서 받은 배치를 동기화합니다.
func (b *HealthBridge) Sync(batch *SyncBatch) (*SyncResult, error) {
	if _, err := ValidateBatch(batch); err != nil {
		return nil, err
	}

	result := &SyncResult{
		UserID:   batch.UserID,
		Source:   batch.Source,
		SyncedAt: time.Now().UTC(),
	}

	for _, raw := range batch.Samples {
		// 정규화
		normalized, err := b.normalizer.Normalize(raw)
		if err != nil {
			result.FailedCount++
			continue
		}

		if err := ValidateSample(normalized); err != nil {
			result.FailedCount++
			continue
		}

		// 중복 체크 (source_id 기반)
		if normalized.SourceID != "" {
			existing, _ := b.repo.FindBySourceID(batch.UserID, normalized.Source, normalized.SourceID)
			if existing != nil {
				// 충돌 처리
				resolved, conflict := b.resolver.Resolve(existing, normalized)
				if conflict != nil && conflict.Resolution == "replaced" {
					if err := b.repo.Update(existing.SourceID, resolved); err == nil {
						result.ConflictCount++
						result.Conflicts = append(result.Conflicts, conflict)
					} else {
						result.FailedCount++
					}
				} else {
					result.DuplicateCount++
				}
				continue
			}
		}

		// 신규 저장
		id, err := b.repo.Save(normalized)
		if err != nil {
			result.FailedCount++
			continue
		}
		result.ImportedCount++
		result.NewSampleIDs = append(result.NewSampleIDs, id)
	}

	sort.Strings(result.NewSampleIDs)
	return result, nil
}

// ============================================================================
// 양방향 동기화 (만파식 → 외부 푸시)
// ============================================================================

// PushBatch는 만파식 측정을 외부 플랫폼으로 푸시할 배치입니다.
type PushBatch struct {
	UserID  string
	Target  string // SourceAppleHealthKit, SourceGoogleHealthConnect
	Samples []*HealthSample
}

// PreparePushBatch는 만파식 측정을 외부 플랫폼 형식으로 준비합니다.
func (b *HealthBridge) PreparePushBatch(userID string, target string, since time.Time) (*PushBatch, error) {
	if userID == "" || target == "" {
		return nil, errors.New("user_id and target required")
	}

	batch := &PushBatch{
		UserID: userID,
		Target: target,
	}

	// 만파식 측정 데이터 수집 (저장소가 지원하는 경우)
	for _, dataType := range []DataType{TypeBloodGlucose, TypeBloodPressureSys, TypeBloodPressureDia, TypeHeartRate, TypeOxygenSaturation} {
		samples, err := b.repo.FindByTimeRange(userID, dataType, since, time.Now().UTC())
		if err != nil {
			continue
		}
		for _, s := range samples {
			if s.Source == SourceManPaSik {
				batch.Samples = append(batch.Samples, s)
			}
		}
	}

	return batch, nil
}
