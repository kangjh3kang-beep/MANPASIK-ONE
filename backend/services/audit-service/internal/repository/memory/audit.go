package memory

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/manpasik/backend/services/audit-service/internal/service"
)

// AuditRepository는 감사 로그의 인메모리 저장소입니다.
type AuditRepository struct {
	mu      sync.RWMutex
	entries []*service.AuditEntry
}

// NewAuditRepository는 새 인메모리 감사 로그 저장소를 생성합니다.
func NewAuditRepository() *AuditRepository {
	return &AuditRepository{
		entries: make([]*service.AuditEntry, 0),
	}
}

func (r *AuditRepository) Store(ctx context.Context, entry *service.AuditEntry) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries = append(r.entries, entry)
	return nil
}

func (r *AuditRepository) List(ctx context.Context, filter service.AuditFilter) ([]*service.AuditEntry, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var filtered []*service.AuditEntry
	for _, e := range r.entries {
		if filter.AdminID != "" && e.AdminID != filter.AdminID {
			continue
		}
		if filter.Action != "" && e.Action != filter.Action {
			continue
		}
		if filter.ResourceType != "" && e.ResourceType != filter.ResourceType {
			continue
		}
		if !filter.StartTime.IsZero() && e.Timestamp.Before(filter.StartTime) {
			continue
		}
		if !filter.EndTime.IsZero() && e.Timestamp.After(filter.EndTime) {
			continue
		}
		filtered = append(filtered, e)
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Timestamp.After(filtered[j].Timestamp)
	})

	total := len(filtered)
	if filter.Offset >= len(filtered) {
		return nil, total, nil
	}
	filtered = filtered[filter.Offset:]
	if filter.Limit > 0 && len(filtered) > filter.Limit {
		filtered = filtered[:filter.Limit]
	}
	return filtered, total, nil
}

func (r *AuditRepository) Get(ctx context.Context, entryID string) (*service.AuditEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, e := range r.entries {
		if e.ID == entryID {
			return e, nil
		}
	}
	return nil, nil
}

// Seed는 테스트용 초기 데이터를 추가합니다.
func (r *AuditRepository) Seed() {
	now := time.Now().UTC()
	entries := []*service.AuditEntry{
		{ID: "aud-001", AdminID: "admin-1", Action: "user.create", ResourceType: "user", ResourceID: "usr-100", Description: "신규 사용자 생성", IPAddress: "10.0.0.1", Timestamp: now.Add(-1 * time.Hour)},
		{ID: "aud-002", AdminID: "admin-1", Action: "config.update", ResourceType: "system_config", ResourceID: "cfg-001", Description: "시스템 설정 변경", IPAddress: "10.0.0.1", Timestamp: now.Add(-30 * time.Minute)},
		{ID: "aud-003", AdminID: "admin-2", Action: "cartridge.approve", ResourceType: "cartridge", ResourceID: "crt-050", Description: "카트리지 승인", IPAddress: "10.0.0.2", Timestamp: now.Add(-10 * time.Minute)},
	}
	r.mu.Lock()
	r.entries = append(r.entries, entries...)
	r.mu.Unlock()
}
