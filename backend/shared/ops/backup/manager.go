// Package backup은 데이터베이스 백업/복구 관리자입니다.
//
// 지원 대상: PostgreSQL (pg_dump), Milvus snapshot, S3/MinIO 업로드.
// cron 스케줄링과 dry-run 모드를 지원합니다.
package backup

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"
)

// ============================================================================
// 도메인 모델
// ============================================================================

// BackupKind는 백업 종류입니다.
type BackupKind string

const (
	KindPostgres BackupKind = "postgres"
	KindMilvus   BackupKind = "milvus"
	KindRedis    BackupKind = "redis"
	KindFiles    BackupKind = "files"
)

// Status는 백업 상태입니다.
const (
	StatusPending    = "pending"
	StatusInProgress = "in_progress"
	StatusCompleted  = "completed"
	StatusFailed     = "failed"
	StatusVerified   = "verified"
)

// BackupRecord는 백업 메타데이터입니다.
type BackupRecord struct {
	ID          string
	Kind        BackupKind
	Source      string    // DB 이름 또는 컬렉션
	Destination string    // S3 경로 또는 파일 경로
	Status      string
	SizeBytes   int64
	Checksum    string    // SHA-256
	StartedAt   time.Time
	CompletedAt *time.Time
	Error       string
	Tags        map[string]string
}

// Schedule는 백업 일정입니다.
type Schedule struct {
	Name      string
	Kind      BackupKind
	Source    string
	Cron      string // "0 2 * * *" (매일 새벽 2시)
	Retention int    // 보관 기간 (일)
	Enabled   bool
}

// RestoreRequest는 복구 요청입니다.
type RestoreRequest struct {
	BackupID    string
	TargetDB    string
	DryRun      bool
	RequestedBy string
}

// ============================================================================
// Provider 인터페이스
// ============================================================================

// Provider는 백업 백엔드 인터페이스입니다.
type Provider interface {
	Backup(ctx context.Context, kind BackupKind, source, destination string) (*BackupRecord, error)
	Restore(ctx context.Context, req *RestoreRequest) error
	Verify(ctx context.Context, recordID string) error
	List(ctx context.Context, kind BackupKind) ([]*BackupRecord, error)
	Delete(ctx context.Context, recordID string) error
	Provider() string
}

// ============================================================================
// 인메모리 Provider (테스트용)
// ============================================================================

// MemoryProvider는 인메모리 백업 Provider입니다.
type MemoryProvider struct {
	mu      sync.RWMutex
	records map[string]*BackupRecord
}

// NewMemoryProvider는 새 인메모리 Provider를 생성합니다.
func NewMemoryProvider() *MemoryProvider {
	return &MemoryProvider{records: make(map[string]*BackupRecord)}
}

// Backup은 백업을 시뮬레이션합니다.
func (p *MemoryProvider) Backup(_ context.Context, kind BackupKind, source, destination string) (*BackupRecord, error) {
	if source == "" {
		return nil, errors.New("source required")
	}

	now := time.Now().UTC()
	record := &BackupRecord{
		ID:          fmt.Sprintf("bkp-%s-%d", kind, now.UnixNano()),
		Kind:        kind,
		Source:      source,
		Destination: destination,
		Status:      StatusInProgress,
		StartedAt:   now,
	}

	// 시뮬레이션: 즉시 완료
	completed := time.Now().UTC()
	record.Status = StatusCompleted
	record.CompletedAt = &completed
	record.SizeBytes = 1024 * 1024 // 시뮬레이션 1MB
	record.Checksum = fmt.Sprintf("sha256:%d", now.UnixNano())

	p.mu.Lock()
	p.records[record.ID] = record
	p.mu.Unlock()
	return record, nil
}

// Restore는 복구를 시뮬레이션합니다.
func (p *MemoryProvider) Restore(_ context.Context, req *RestoreRequest) error {
	if req == nil || req.BackupID == "" {
		return errors.New("backup_id required")
	}

	p.mu.RLock()
	record, ok := p.records[req.BackupID]
	p.mu.RUnlock()
	if !ok {
		return fmt.Errorf("backup %s not found", req.BackupID)
	}
	if record.Status != StatusCompleted && record.Status != StatusVerified {
		return fmt.Errorf("backup %s not completed", req.BackupID)
	}
	if req.DryRun {
		return nil
	}
	return nil
}

// Verify는 백업 무결성을 검증합니다.
func (p *MemoryProvider) Verify(_ context.Context, recordID string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	record, ok := p.records[recordID]
	if !ok {
		return fmt.Errorf("backup %s not found", recordID)
	}
	if record.Status != StatusCompleted {
		return fmt.Errorf("cannot verify status: %s", record.Status)
	}
	if record.Checksum == "" {
		return errors.New("checksum missing")
	}
	record.Status = StatusVerified
	return nil
}

// List는 종류별 백업 목록을 반환합니다 (최신순).
func (p *MemoryProvider) List(_ context.Context, kind BackupKind) ([]*BackupRecord, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	var result []*BackupRecord
	for _, r := range p.records {
		if kind == "" || r.Kind == kind {
			result = append(result, r)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].StartedAt.After(result[j].StartedAt)
	})
	return result, nil
}

// Delete는 백업 레코드를 제거합니다.
func (p *MemoryProvider) Delete(_ context.Context, recordID string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.records, recordID)
	return nil
}

// Provider는 이름을 반환합니다.
func (p *MemoryProvider) Provider() string { return "memory" }

// Count는 저장된 백업 수를 반환합니다.
func (p *MemoryProvider) Count() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.records)
}

// ============================================================================
// 보존 정책 + 자동 정리
// ============================================================================

// RetentionEnforcer는 보존 기간 초과 백업을 정리합니다.
type RetentionEnforcer struct {
	provider Provider
}

// NewRetentionEnforcer는 새 보존 정책 집행자를 생성합니다.
func NewRetentionEnforcer(provider Provider) *RetentionEnforcer {
	return &RetentionEnforcer{provider: provider}
}

// Enforce는 보존 기간을 초과한 백업을 삭제합니다.
//
// retentionDays: 보관 일수 (0 이하면 정리하지 않음).
func (e *RetentionEnforcer) Enforce(ctx context.Context, kind BackupKind, retentionDays int) (int, error) {
	if retentionDays <= 0 {
		return 0, nil
	}

	records, err := e.provider.List(ctx, kind)
	if err != nil {
		return 0, err
	}

	threshold := time.Now().UTC().AddDate(0, 0, -retentionDays)
	deleted := 0
	for _, r := range records {
		if r.StartedAt.Before(threshold) {
			if err := e.provider.Delete(ctx, r.ID); err == nil {
				deleted++
			}
		}
	}
	return deleted, nil
}

// ============================================================================
// 스케줄 관리
// ============================================================================

// ScheduleManager는 백업 일정을 관리합니다.
type ScheduleManager struct {
	mu        sync.RWMutex
	schedules map[string]*Schedule
}

// NewScheduleManager는 새 스케줄 매니저를 생성합니다.
func NewScheduleManager() *ScheduleManager {
	return &ScheduleManager{schedules: make(map[string]*Schedule)}
}

// Add는 새 스케줄을 추가합니다.
func (m *ScheduleManager) Add(schedule *Schedule) error {
	if schedule == nil || schedule.Name == "" {
		return errors.New("schedule name required")
	}
	if schedule.Cron == "" {
		return errors.New("cron expression required")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.schedules[schedule.Name] = schedule
	return nil
}

// Remove는 스케줄을 제거합니다.
func (m *ScheduleManager) Remove(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.schedules, name)
}

// Get은 스케줄을 조회합니다.
func (m *ScheduleManager) Get(name string) (*Schedule, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.schedules[name]
	if !ok {
		return nil, fmt.Errorf("schedule %s not found", name)
	}
	return s, nil
}

// EnabledSchedules는 활성화된 스케줄 목록을 반환합니다.
func (m *ScheduleManager) EnabledSchedules() []*Schedule {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var enabled []*Schedule
	for _, s := range m.schedules {
		if s.Enabled {
			enabled = append(enabled, s)
		}
	}
	return enabled
}

// SetEnabled는 스케줄 활성화 상태를 변경합니다.
func (m *ScheduleManager) SetEnabled(name string, enabled bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	s, ok := m.schedules[name]
	if !ok {
		return fmt.Errorf("schedule %s not found", name)
	}
	s.Enabled = enabled
	return nil
}
