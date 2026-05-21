package fhir

import (
	"context"
	"errors"
	"sort"
	"sync"
	"time"
)

// ============================================================================
// BundleStore — HL7 v2 외부 검사 결과 (FHIR Bundle) 영속화 (Phase AV)
// ============================================================================
//
// 외부 LIS/EHR 가 송신한 HL7 메시지를 FHIR Bundle 로 변환한 후 만파식 내부에
// 보관하여 운영자가 다시 조회할 수 있게 합니다. 의료 표준 정보 모델 측면에서
// Bundle 은 트랜잭션 단위로 그대로 유지하는 것이 추적성/감사 가능성 면에서
// 가장 안전합니다.
//
// tenant 격리: 멀티테넌트 환경에서 다른 조직의 Bundle 을 볼 수 없도록
// tenantID 별 namespace 로 격리됩니다. 빈 tenant ("") 는 단일 테넌트 모드.

// ErrBundleNotFound 는 Get/Delete 시 Bundle 미존재.
var ErrBundleNotFound = errors.New("Bundle 항목 없음")

// StoredBundle 은 저장된 Bundle + 메타.
type StoredBundle struct {
	Bundle    *Bundle   `json:"bundle"`
	TenantID  string    `json:"tenant_id,omitempty"`
	PatientID string    `json:"patient_id,omitempty"` // 빠른 환자별 조회용 (PID-3 추출값)
	StoredAt  time.Time `json:"stored_at"`
}

// BundleStore 는 Bundle 저장/조회 인터페이스 — PostgreSQL/Memory/외부 LIS 백엔드
// 모두 같은 인터페이스로 교체 가능합니다.
type BundleStore interface {
	// Save 는 Bundle 저장. Bundle.ID 가 이미 존재하면 멱등 (update 또는 noop).
	// patientID 는 PID-3 값. tenant 미설정 시 ""를 전달.
	Save(ctx context.Context, b *Bundle, tenantID, patientID string) error

	// Get 은 단일 Bundle 조회. tenant 가 일치하지 않으면 NotFound.
	Get(ctx context.Context, bundleID, tenantID string) (*StoredBundle, error)

	// ListByPatient 는 tenant 내 환자별 Bundle 목록 (최신순). limit 0 = 모두.
	ListByPatient(ctx context.Context, tenantID, patientID string, limit int) ([]*StoredBundle, error)

	// ListByTenant 는 tenant 의 전체 Bundle (최신순). limit 0 = 모두.
	ListByTenant(ctx context.Context, tenantID string, limit int) ([]*StoredBundle, error)

	// Delete 는 단일 Bundle 제거. tenant 검사 포함.
	Delete(ctx context.Context, bundleID, tenantID string) error

	// Count 는 tenant 내 Bundle 개수 (운영 메트릭).
	Count(ctx context.Context, tenantID string) int
}

// MemoryBundleStore 는 인메모리 BundleStore (개발/테스트/소규모 운영).
//
// 본격 운영은 PostgreSQL 백엔드로 교체 권장 (별도 Phase 에서 구현).
type MemoryBundleStore struct {
	mu       sync.RWMutex
	bundles  map[string]*StoredBundle // key: tenantID + "::" + bundleID
	maxSize  int                       // tenant 당 최대 보관 수 (0 = 무제한)
}

// NewMemoryBundleStore 생성. maxPerTenant 가 양수이면 tenant 별 FIFO 자동 제거.
func NewMemoryBundleStore(maxPerTenant int) *MemoryBundleStore {
	return &MemoryBundleStore{
		bundles: make(map[string]*StoredBundle),
		maxSize: maxPerTenant,
	}
}

func keyOf(tenantID, bundleID string) string {
	return tenantID + "::" + bundleID
}

// Save 구현.
func (s *MemoryBundleStore) Save(_ context.Context, b *Bundle, tenantID, patientID string) error {
	if b == nil || b.ID == "" {
		return errors.New("Bundle.ID 가 필요")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	key := keyOf(tenantID, b.ID)
	// 이미 존재하면 timestamp 만 갱신 (멱등 update 의미)
	if existing, ok := s.bundles[key]; ok {
		existing.Bundle = b
		existing.PatientID = patientID
		// StoredAt 은 최초 저장 시점 유지
		return nil
	}
	s.bundles[key] = &StoredBundle{
		Bundle:    b,
		TenantID:  tenantID,
		PatientID: patientID,
		StoredAt:  time.Now().UTC(),
	}
	// maxSize 초과 시 가장 오래된 항목 제거 (FIFO)
	if s.maxSize > 0 {
		s.evictIfNeeded(tenantID)
	}
	return nil
}

// evictIfNeeded 는 tenant 의 항목이 maxSize 를 초과하면 가장 오래된 항목 제거.
//
// 내부 호출 — 호출 측이 mu 락 보유 가정.
func (s *MemoryBundleStore) evictIfNeeded(tenantID string) {
	tenantPrefix := tenantID + "::"
	tenantBundles := make([]*StoredBundle, 0)
	for k, v := range s.bundles {
		if len(k) >= len(tenantPrefix) && k[:len(tenantPrefix)] == tenantPrefix {
			tenantBundles = append(tenantBundles, v)
		}
	}
	if len(tenantBundles) <= s.maxSize {
		return
	}
	// StoredAt 오름차순 정렬 → 오래된 것이 앞
	sort.Slice(tenantBundles, func(i, j int) bool {
		return tenantBundles[i].StoredAt.Before(tenantBundles[j].StoredAt)
	})
	excess := len(tenantBundles) - s.maxSize
	for i := 0; i < excess; i++ {
		delete(s.bundles, keyOf(tenantBundles[i].TenantID, tenantBundles[i].Bundle.ID))
	}
}

// Get 구현.
func (s *MemoryBundleStore) Get(_ context.Context, bundleID, tenantID string) (*StoredBundle, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	entry, ok := s.bundles[keyOf(tenantID, bundleID)]
	if !ok {
		return nil, ErrBundleNotFound
	}
	return entry, nil
}

// ListByPatient 구현.
func (s *MemoryBundleStore) ListByPatient(_ context.Context, tenantID, patientID string, limit int) ([]*StoredBundle, error) {
	if patientID == "" {
		return nil, errors.New("patientID 필수")
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := s.filterAndSort(tenantID, func(b *StoredBundle) bool {
		return b.PatientID == patientID
	})
	return applyLimit(out, limit), nil
}

// ListByTenant 구현.
func (s *MemoryBundleStore) ListByTenant(_ context.Context, tenantID string, limit int) ([]*StoredBundle, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := s.filterAndSort(tenantID, func(*StoredBundle) bool { return true })
	return applyLimit(out, limit), nil
}

func (s *MemoryBundleStore) filterAndSort(tenantID string, match func(*StoredBundle) bool) []*StoredBundle {
	tenantPrefix := tenantID + "::"
	out := make([]*StoredBundle, 0)
	for k, v := range s.bundles {
		if len(k) >= len(tenantPrefix) && k[:len(tenantPrefix)] == tenantPrefix && match(v) {
			out = append(out, v)
		}
	}
	// 최신순
	sort.Slice(out, func(i, j int) bool {
		return out[i].StoredAt.After(out[j].StoredAt)
	})
	return out
}

func applyLimit(arr []*StoredBundle, limit int) []*StoredBundle {
	if limit <= 0 || limit >= len(arr) {
		return arr
	}
	return arr[:limit]
}

// Delete 구현.
func (s *MemoryBundleStore) Delete(_ context.Context, bundleID, tenantID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := keyOf(tenantID, bundleID)
	if _, ok := s.bundles[key]; !ok {
		return ErrBundleNotFound
	}
	delete(s.bundles, key)
	return nil
}

// Count 구현.
func (s *MemoryBundleStore) Count(_ context.Context, tenantID string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	tenantPrefix := tenantID + "::"
	n := 0
	for k := range s.bundles {
		if len(k) >= len(tenantPrefix) && k[:len(tenantPrefix)] == tenantPrefix {
			n++
		}
	}
	return n
}
