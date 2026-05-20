// Package service는 audit-service의 비즈니스 로직을 구현합니다.
//
// 감사 로그 수집, 저장, 조회를 담당합니다.
package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// AuditEntry는 감사 로그 엔트리입니다.
type AuditEntry struct {
	ID           string
	AdminID      string
	Action       string // "user.create", "config.update", "cartridge.approve" 등
	ResourceType string // "user", "system_config", "cartridge" 등
	ResourceID   string
	Description  string
	IPAddress    string
	UserAgent    string
	Metadata     map[string]string
	Timestamp    time.Time
}

// AuditFilter는 감사 로그 필터링 조건입니다.
type AuditFilter struct {
	AdminID      string
	Action       string
	ResourceType string
	StartTime    time.Time
	EndTime      time.Time
	Limit        int
	Offset       int
}

// AuditRepository는 감사 로그 저장소 인터페이스입니다.
type AuditRepository interface {
	Store(ctx context.Context, entry *AuditEntry) error
	List(ctx context.Context, filter AuditFilter) ([]*AuditEntry, int, error)
	Get(ctx context.Context, entryID string) (*AuditEntry, error)
}

// AuditService는 감사 로그 서비스입니다.
type AuditService struct {
	logger *zap.Logger
	repo   AuditRepository
}

// NewAuditService는 새 AuditService를 생성합니다.
func NewAuditService(logger *zap.Logger, repo AuditRepository) *AuditService {
	return &AuditService{logger: logger, repo: repo}
}

type RecordActionInput struct {
	AdminID      string
	Action       string
	ResourceType string
	ResourceID   string
	Description  string
	IPAddress    string
	UserAgent    string
	Metadata     map[string]string
	Timestamp    time.Time
}

// RecordAction은 감사 로그를 기록합니다.
func (s *AuditService) RecordAction(ctx context.Context, adminID, action, resourceType, resourceID, description, ipAddr string) (*AuditEntry, error) {
	return s.RecordActionWithMetadata(ctx, RecordActionInput{
		AdminID:      adminID,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Description:  description,
		IPAddress:    ipAddr,
	})
}

func (s *AuditService) RecordActionWithMetadata(ctx context.Context, input RecordActionInput) (*AuditEntry, error) {
	timestamp := input.Timestamp
	if timestamp.IsZero() {
		timestamp = time.Now().UTC()
	}
	entry := &AuditEntry{
		ID:           uuid.New().String(),
		AdminID:      input.AdminID,
		Action:       input.Action,
		ResourceType: input.ResourceType,
		ResourceID:   input.ResourceID,
		Description:  input.Description,
		IPAddress:    input.IPAddress,
		UserAgent:    input.UserAgent,
		Metadata:     input.Metadata,
		Timestamp:    timestamp,
	}

	if err := s.repo.Store(ctx, entry); err != nil {
		s.logger.Error("감사 로그 저장 실패", zap.Error(err))
		return nil, err
	}

	s.logger.Info("감사 로그 기록",
		zap.String("entry_id", entry.ID),
		zap.String("action", input.Action),
		zap.String("admin_id", input.AdminID),
	)
	return entry, nil
}

// ListEntries는 감사 로그를 조회합니다.
func (s *AuditService) ListEntries(ctx context.Context, filter AuditFilter) ([]*AuditEntry, int, error) {
	if filter.Limit <= 0 {
		filter.Limit = 50
	}
	if filter.Limit > 500 {
		filter.Limit = 500
	}
	return s.repo.List(ctx, filter)
}

// GetEntry는 단일 감사 로그를 조회합니다.
func (s *AuditService) GetEntry(ctx context.Context, entryID string) (*AuditEntry, error) {
	return s.repo.Get(ctx, entryID)
}

// ============================================================================
// PHI 접근 추적 (HIPAA §164.312(b), GDPR Art.15)
// ============================================================================

// PHIAccessEvent는 개인건강정보(PHI) 접근 이벤트입니다.
type PHIAccessEvent struct {
	EventID        string
	UserID         string   // 접근 대상 사용자
	AccessorID     string   // 접근자 (관리자/의사)
	AccessorRole   string   // 접근자 역할
	Action         string   // "view", "export", "share", "modify", "delete"
	ResourceType   string   // "health_record", "measurement", "prescription"
	ResourceIDs    []string // 접근한 리소스 ID 목록
	FieldsAccessed []string // 접근한 PHI 필드 목록
	Purpose        string   // 접근 목적
	ConsentID      string   // 관련 동의 ID (있는 경우)
	IPAddress      string
	Timestamp      time.Time
}

// RecordPHIAccess는 PHI 접근 이벤트를 기록합니다.
func (s *AuditService) RecordPHIAccess(ctx context.Context, event *PHIAccessEvent) error {
	if event == nil {
		return nil
	}
	if event.UserID == "" || event.AccessorID == "" {
		s.logger.Warn("PHI 접근 이벤트: 필수 필드 누락",
			zap.String("user_id", event.UserID),
			zap.String("accessor_id", event.AccessorID),
		)
	}

	event.EventID = uuid.New().String()
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now().UTC()
	}

	// 감사 로그로도 기록
	metadata := map[string]string{
		"phi_access":    "true",
		"target_user":   event.UserID,
		"accessor_role": event.AccessorRole,
		"purpose":       event.Purpose,
	}
	if event.ConsentID != "" {
		metadata["consent_id"] = event.ConsentID
	}

	entry := &AuditEntry{
		ID:           event.EventID,
		AdminID:      event.AccessorID,
		Action:       "phi." + event.Action,
		ResourceType: event.ResourceType,
		ResourceID:   event.UserID,
		Description:  "PHI 접근: " + event.Action + " (" + event.Purpose + ")",
		IPAddress:    event.IPAddress,
		Metadata:     metadata,
		Timestamp:    event.Timestamp,
	}

	if err := s.repo.Store(ctx, entry); err != nil {
		s.logger.Error("PHI 접근 로그 저장 실패", zap.Error(err))
		return err
	}

	s.logger.Info("PHI 접근 기록",
		zap.String("event_id", event.EventID),
		zap.String("accessor", event.AccessorID),
		zap.String("target_user", event.UserID),
		zap.String("action", event.Action),
	)
	return nil
}

// ListPHIAccess는 특정 사용자의 PHI 접근 이력을 조회합니다.
// GDPR Art.15 (접근권) 지원: 누가, 언제, 어떤 데이터에 접근했는지 조회 가능.
func (s *AuditService) ListPHIAccess(ctx context.Context, userID string, limit, offset int) ([]*AuditEntry, int, error) {
	filter := AuditFilter{
		ResourceType: "phi_access",
		Limit:        limit,
		Offset:       offset,
	}

	// PHI 접근은 ResourceID에 대상 사용자 ID가 저장됨
	entries, total, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// 대상 사용자 필터링
	var result []*AuditEntry
	for _, e := range entries {
		if e.Metadata != nil && e.Metadata["target_user"] == userID {
			result = append(result, e)
		}
		if e.ResourceID == userID && e.Metadata != nil && e.Metadata["phi_access"] == "true" {
			result = append(result, e)
		}
	}

	return result, total, nil
}

// ============================================================================
// 변경 이력 추적 (Change Tracking)
// ============================================================================

// ChangeRecord는 데이터 변경 이력입니다.
type ChangeRecord struct {
	ChangeID     string
	EntityType   string // "user", "health_record", "config", "admin"
	EntityID     string
	ChangedBy    string
	ChangeType   string // "create", "update", "delete"
	FieldChanges []FieldChange
	Timestamp    time.Time
}

// FieldChange는 개별 필드 변경 내역입니다.
type FieldChange struct {
	FieldName string
	OldValue  string
	NewValue  string
}

// RecordChange는 데이터 변경 이력을 기록합니다.
func (s *AuditService) RecordChange(ctx context.Context, record *ChangeRecord) error {
	if record == nil {
		return nil
	}

	record.ChangeID = uuid.New().String()
	if record.Timestamp.IsZero() {
		record.Timestamp = time.Now().UTC()
	}

	// 변경 내역 요약 생성
	desc := record.ChangeType + " " + record.EntityType + ":" + record.EntityID
	if len(record.FieldChanges) > 0 {
		desc += " (" + string(rune('0'+len(record.FieldChanges))) + " fields changed)"
	}

	metadata := map[string]string{
		"change_tracking": "true",
		"change_type":     record.ChangeType,
		"entity_type":     record.EntityType,
	}

	// 각 필드 변경을 메타데이터에 기록
	for i, fc := range record.FieldChanges {
		if i >= 10 {
			break // 최대 10개 필드까지 메타데이터에 기록
		}
		key := "field_" + fc.FieldName
		metadata[key+"_old"] = fc.OldValue
		metadata[key+"_new"] = fc.NewValue
	}

	entry := &AuditEntry{
		ID:           record.ChangeID,
		AdminID:      record.ChangedBy,
		Action:       "change." + record.ChangeType,
		ResourceType: record.EntityType,
		ResourceID:   record.EntityID,
		Description:  desc,
		Metadata:     metadata,
		Timestamp:    record.Timestamp,
	}

	if err := s.repo.Store(ctx, entry); err != nil {
		s.logger.Error("변경 이력 저장 실패", zap.Error(err))
		return err
	}

	s.logger.Info("변경 이력 기록",
		zap.String("change_id", record.ChangeID),
		zap.String("entity", record.EntityType+":"+record.EntityID),
		zap.String("change_type", record.ChangeType),
		zap.Int("fields_changed", len(record.FieldChanges)),
	)
	return nil
}
