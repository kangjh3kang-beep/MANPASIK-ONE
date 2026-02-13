// Package memory는 인메모리 카트리지 저장소입니다 (개발/테스트용).
package memory

import (
	"context"
	"sync"

	"github.com/manpasik/backend/services/cartridge-service/internal/service"
)

// ============================================================================
// CartridgeUsageRepository — 카트리지 사용 기록 저장소
// ============================================================================

// CartridgeUsageRepository는 인메모리 사용 기록 저장소입니다.
type CartridgeUsageRepository struct {
	mu      sync.RWMutex
	records []*service.CartridgeUsageRecord
}

// NewCartridgeUsageRepository는 인메모리 사용 기록 저장소를 생성합니다.
func NewCartridgeUsageRepository() *CartridgeUsageRepository {
	return &CartridgeUsageRepository{
		records: make([]*service.CartridgeUsageRecord, 0),
	}
}

// Create는 사용 기록을 추가합니다.
func (r *CartridgeUsageRepository) Create(_ context.Context, record *service.CartridgeUsageRecord) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *record
	r.records = append(r.records, &cp)
	return nil
}

// ListByUserID는 사용자의 사용 기록을 조회합니다 (페이지네이션).
func (r *CartridgeUsageRepository) ListByUserID(_ context.Context, userID string, limit, offset int32) ([]*service.CartridgeUsageRecord, int32, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// 사용자 ID로 필터링
	var filtered []*service.CartridgeUsageRecord
	for _, rec := range r.records {
		if rec.UserID == userID {
			cp := *rec
			filtered = append(filtered, &cp)
		}
	}

	totalCount := int32(len(filtered))

	// 페이지네이션
	start := int(offset)
	if start >= len(filtered) {
		return []*service.CartridgeUsageRecord{}, totalCount, nil
	}
	end := start + int(limit)
	if end > len(filtered) {
		end = len(filtered)
	}

	return filtered[start:end], totalCount, nil
}

// ============================================================================
// CartridgeStateRepository — 카트리지 잔여 사용 상태 저장소
// ============================================================================

// CartridgeStateRepository는 인메모리 카트리지 상태 저장소입니다.
type CartridgeStateRepository struct {
	mu    sync.RWMutex
	byUID map[string]*service.CartridgeRemainingInfo
}

// NewCartridgeStateRepository는 인메모리 상태 저장소를 생성합니다.
func NewCartridgeStateRepository() *CartridgeStateRepository {
	return &CartridgeStateRepository{
		byUID: make(map[string]*service.CartridgeRemainingInfo),
	}
}

// GetByUID는 카트리지 UID로 상태를 조회합니다.
func (r *CartridgeStateRepository) GetByUID(_ context.Context, uid string) (*service.CartridgeRemainingInfo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	info, ok := r.byUID[uid]
	if !ok {
		return nil, nil
	}
	cp := *info
	return &cp, nil
}

// Upsert는 카트리지 상태를 생성하거나 업데이트합니다.
func (r *CartridgeStateRepository) Upsert(_ context.Context, info *service.CartridgeRemainingInfo) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *info
	r.byUID[info.CartridgeUID] = &cp
	return nil
}

// DecrementUses는 카트리지 잔여 사용 횟수를 1 감소시킵니다.
func (r *CartridgeStateRepository) DecrementUses(_ context.Context, uid string) (int32, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	info, ok := r.byUID[uid]
	if !ok {
		return 0, nil
	}
	if info.RemainingUses > 0 {
		info.RemainingUses--
	}
	return info.RemainingUses, nil
}
