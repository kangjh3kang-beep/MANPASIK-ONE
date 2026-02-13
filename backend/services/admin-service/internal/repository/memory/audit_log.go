// Package memory는 인메모리 감사 로그 저장소입니다 (개발/테스트용).
// AS-9: 관리자 액션 감사 로그 기능
package memory

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

// AuditLog는 관리자 액션 감사 로그 엔티티입니다.
type AuditLog struct {
	ID        string
	AdminID   string
	Action    string // "config_update", "config_create", "config_delete", "user_ban", "system_restart" 등
	Resource  string // 대상 리소스 (예: "config:security.jwt_ttl_hours")
	OldValue  string
	NewValue  string
	IPAddress string
	UserAgent string
	CreatedAt time.Time
}

// AuditLogStore는 확장된 감사 로그 저장소 인터페이스입니다.
// NOTE: 기존 AuditLogRepository 타입과의 이름 충돌을 피해 AuditLogStore로 명명합니다.
// 기존 AuditLogRepository는 service.AuditLogEntry 기반이며,
// AuditLogStore는 OldValue/NewValue 추적이 가능한 확장 감사 로그입니다.
type AuditLogStore interface {
	Create(ctx context.Context, log *AuditLog) error
	ListByAdmin(ctx context.Context, adminID string, limit, offset int) ([]*AuditLog, error)
	ListByAction(ctx context.Context, action string, limit, offset int) ([]*AuditLog, error)
	ListAll(ctx context.Context, limit, offset int) ([]*AuditLog, error)
	Count(ctx context.Context) (int, error)
}

// InMemoryAuditLogStore는 인메모리 감사 로그 저장소 구현체입니다.
// sync.RWMutex를 사용하여 동시성을 보장합니다.
type InMemoryAuditLogStore struct {
	mu   sync.RWMutex
	logs []*AuditLog
}

// NewAuditLogStore는 인메모리 AuditLogStore를 생성합니다.
func NewAuditLogStore() *InMemoryAuditLogStore {
	return &InMemoryAuditLogStore{
		logs: make([]*AuditLog, 0),
	}
}

// Create는 새 감사 로그 항목을 저장합니다.
// ID가 비어 있으면 자동 생성합니다.
func (s *InMemoryAuditLogStore) Create(_ context.Context, log *AuditLog) error {
	if log == nil {
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// 깊은 복사
	cp := *log
	if cp.ID == "" {
		cp.ID = uuid.New().String()
	}
	if cp.CreatedAt.IsZero() {
		cp.CreatedAt = time.Now().UTC()
	}

	s.logs = append(s.logs, &cp)
	return nil
}

// ListByAdmin은 특정 관리자의 감사 로그를 최신순으로 조회합니다.
func (s *InMemoryAuditLogStore) ListByAdmin(_ context.Context, adminID string, limit, offset int) ([]*AuditLog, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var filtered []*AuditLog
	// 역순(최신순) 조회
	for i := len(s.logs) - 1; i >= 0; i-- {
		if s.logs[i].AdminID == adminID {
			cp := *s.logs[i]
			filtered = append(filtered, &cp)
		}
	}

	return paginate(filtered, limit, offset), nil
}

// ListByAction은 특정 액션 유형의 감사 로그를 최신순으로 조회합니다.
func (s *InMemoryAuditLogStore) ListByAction(_ context.Context, action string, limit, offset int) ([]*AuditLog, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var filtered []*AuditLog
	for i := len(s.logs) - 1; i >= 0; i-- {
		if s.logs[i].Action == action {
			cp := *s.logs[i]
			filtered = append(filtered, &cp)
		}
	}

	return paginate(filtered, limit, offset), nil
}

// ListAll은 모든 감사 로그를 최신순으로 조회합니다.
func (s *InMemoryAuditLogStore) ListAll(_ context.Context, limit, offset int) ([]*AuditLog, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 역순(최신순) 복사
	result := make([]*AuditLog, 0, len(s.logs))
	for i := len(s.logs) - 1; i >= 0; i-- {
		cp := *s.logs[i]
		result = append(result, &cp)
	}

	return paginate(result, limit, offset), nil
}

// Count는 전체 감사 로그 수를 반환합니다.
func (s *InMemoryAuditLogStore) Count(_ context.Context) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.logs), nil
}

// paginate는 슬라이스에 limit/offset 페이지네이션을 적용합니다.
func paginate(items []*AuditLog, limit, offset int) []*AuditLog {
	total := len(items)
	if offset >= total {
		return nil
	}

	end := offset + limit
	if end > total {
		end = total
	}

	if limit <= 0 {
		// limit이 0 이하이면 offset 이후 전체 반환
		return items[offset:]
	}

	return items[offset:end]
}
