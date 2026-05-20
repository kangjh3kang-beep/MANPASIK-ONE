package service

import (
	"context"
	"testing"
	"time"

	"go.uber.org/zap"
)

type mockAuditRepo struct {
	entries []*AuditEntry
}

func (r *mockAuditRepo) Store(ctx context.Context, entry *AuditEntry) error {
	r.entries = append(r.entries, entry)
	return nil
}
func (r *mockAuditRepo) List(ctx context.Context, filter AuditFilter) ([]*AuditEntry, int, error) {
	var filtered []*AuditEntry
	for _, e := range r.entries {
		if filter.AdminID != "" && e.AdminID != filter.AdminID {
			continue
		}
		if filter.Action != "" && e.Action != filter.Action {
			continue
		}
		filtered = append(filtered, e)
	}
	return filtered, len(filtered), nil
}
func (r *mockAuditRepo) Get(ctx context.Context, entryID string) (*AuditEntry, error) {
	for _, e := range r.entries {
		if e.ID == entryID {
			return e, nil
		}
	}
	return nil, nil
}

func TestRecordAction(t *testing.T) {
	repo := &mockAuditRepo{}
	svc := NewAuditService(zap.NewNop(), repo)

	entry, err := svc.RecordAction(context.Background(), "admin-1", "user.create", "user", "usr-100", "사용자 생성", "10.0.0.1")
	if err != nil {
		t.Fatalf("RecordAction failed: %v", err)
	}
	if entry.ID == "" {
		t.Fatal("entry ID should not be empty")
	}
	if entry.Action != "user.create" {
		t.Fatalf("expected action user.create, got %s", entry.Action)
	}
	if len(repo.entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(repo.entries))
	}
}

func TestListEntries(t *testing.T) {
	repo := &mockAuditRepo{}
	svc := NewAuditService(zap.NewNop(), repo)

	svc.RecordAction(context.Background(), "admin-1", "user.create", "user", "usr-100", "a", "")
	svc.RecordAction(context.Background(), "admin-2", "config.update", "config", "cfg-001", "b", "")
	svc.RecordAction(context.Background(), "admin-1", "cartridge.approve", "cartridge", "crt-050", "c", "")

	entries, total, err := svc.ListEntries(context.Background(), AuditFilter{AdminID: "admin-1"})
	if err != nil {
		t.Fatalf("ListEntries failed: %v", err)
	}
	if total != 2 {
		t.Fatalf("expected 2 entries for admin-1, got %d", total)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}

func TestListEntriesDefaultLimit(t *testing.T) {
	repo := &mockAuditRepo{}
	svc := NewAuditService(zap.NewNop(), repo)

	filter := AuditFilter{Limit: 0}
	_, _, err := svc.ListEntries(context.Background(), filter)
	if err != nil {
		t.Fatalf("ListEntries with default limit failed: %v", err)
	}
}

func TestGetEntry(t *testing.T) {
	repo := &mockAuditRepo{}
	svc := NewAuditService(zap.NewNop(), repo)

	created, _ := svc.RecordAction(context.Background(), "admin-1", "user.create", "user", "usr-100", "test", "")

	found, err := svc.GetEntry(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("GetEntry failed: %v", err)
	}
	if found == nil {
		t.Fatal("entry should not be nil")
	}
	if found.ID != created.ID {
		t.Fatalf("expected ID %s, got %s", created.ID, found.ID)
	}
}

func TestTimestamp(t *testing.T) {
	repo := &mockAuditRepo{}
	svc := NewAuditService(zap.NewNop(), repo)

	before := time.Now().UTC()
	entry, _ := svc.RecordAction(context.Background(), "admin-1", "test", "test", "t-1", "test", "")
	after := time.Now().UTC()

	if entry.Timestamp.Before(before) || entry.Timestamp.After(after) {
		t.Fatal("timestamp should be between before and after")
	}
}

// ============================================================================
// Phase H: PHI 접근 추적 테스트
// ============================================================================

func TestRecordPHIAccess(t *testing.T) {
	repo := &mockAuditRepo{}
	svc := NewAuditService(zap.NewNop(), repo)

	event := &PHIAccessEvent{
		UserID:         "patient-001",
		AccessorID:     "doctor-001",
		AccessorRole:   "physician",
		Action:         "view",
		ResourceType:   "health_record",
		ResourceIDs:    []string{"hr-001", "hr-002"},
		FieldsAccessed: []string{"blood_glucose", "blood_pressure"},
		Purpose:        "진료",
		ConsentID:      "consent-001",
		IPAddress:      "192.168.1.10",
	}

	err := svc.RecordPHIAccess(context.Background(), event)
	if err != nil {
		t.Fatalf("RecordPHIAccess failed: %v", err)
	}

	if event.EventID == "" {
		t.Error("EventID should be set")
	}
	if len(repo.entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(repo.entries))
	}

	entry := repo.entries[0]
	if entry.Action != "phi.view" {
		t.Errorf("action=%s, want phi.view", entry.Action)
	}
	if entry.Metadata["phi_access"] != "true" {
		t.Error("metadata should have phi_access=true")
	}
	if entry.Metadata["target_user"] != "patient-001" {
		t.Error("metadata should have target_user")
	}
	if entry.Metadata["consent_id"] != "consent-001" {
		t.Error("metadata should have consent_id")
	}
}

func TestRecordPHIAccess_Nil(t *testing.T) {
	repo := &mockAuditRepo{}
	svc := NewAuditService(zap.NewNop(), repo)

	err := svc.RecordPHIAccess(context.Background(), nil)
	if err != nil {
		t.Error("nil event should not return error")
	}
	if len(repo.entries) != 0 {
		t.Error("nil event should not create entries")
	}
}

func TestRecordPHIAccess_ExportAction(t *testing.T) {
	repo := &mockAuditRepo{}
	svc := NewAuditService(zap.NewNop(), repo)

	event := &PHIAccessEvent{
		UserID:       "patient-002",
		AccessorID:   "admin-001",
		AccessorRole: "admin",
		Action:       "export",
		ResourceType: "measurement",
		Purpose:      "연구 데이터 내보내기",
	}

	err := svc.RecordPHIAccess(context.Background(), event)
	if err != nil {
		t.Fatal(err)
	}
	if repo.entries[0].Action != "phi.export" {
		t.Errorf("action=%s, want phi.export", repo.entries[0].Action)
	}
}

// ============================================================================
// Phase H: 변경 이력 추적 테스트
// ============================================================================

func TestRecordChange(t *testing.T) {
	repo := &mockAuditRepo{}
	svc := NewAuditService(zap.NewNop(), repo)

	record := &ChangeRecord{
		EntityType: "user",
		EntityID:   "usr-001",
		ChangedBy:  "admin-001",
		ChangeType: "update",
		FieldChanges: []FieldChange{
			{FieldName: "email", OldValue: "old@test.com", NewValue: "new@test.com"},
			{FieldName: "display_name", OldValue: "홍길동", NewValue: "김길동"},
		},
	}

	err := svc.RecordChange(context.Background(), record)
	if err != nil {
		t.Fatalf("RecordChange failed: %v", err)
	}

	if record.ChangeID == "" {
		t.Error("ChangeID should be set")
	}
	if len(repo.entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(repo.entries))
	}

	entry := repo.entries[0]
	if entry.Action != "change.update" {
		t.Errorf("action=%s, want change.update", entry.Action)
	}
	if entry.Metadata["change_tracking"] != "true" {
		t.Error("metadata should have change_tracking=true")
	}
	if entry.Metadata["field_email_old"] != "old@test.com" {
		t.Error("metadata should have field_email_old")
	}
	if entry.Metadata["field_email_new"] != "new@test.com" {
		t.Error("metadata should have field_email_new")
	}
}

func TestRecordChange_Nil(t *testing.T) {
	repo := &mockAuditRepo{}
	svc := NewAuditService(zap.NewNop(), repo)

	err := svc.RecordChange(context.Background(), nil)
	if err != nil {
		t.Error("nil record should not return error")
	}
}

func TestRecordChange_Delete(t *testing.T) {
	repo := &mockAuditRepo{}
	svc := NewAuditService(zap.NewNop(), repo)

	record := &ChangeRecord{
		EntityType: "health_record",
		EntityID:   "hr-001",
		ChangedBy:  "admin-001",
		ChangeType: "delete",
	}

	err := svc.RecordChange(context.Background(), record)
	if err != nil {
		t.Fatal(err)
	}
	if repo.entries[0].Action != "change.delete" {
		t.Errorf("action=%s, want change.delete", repo.entries[0].Action)
	}
}
