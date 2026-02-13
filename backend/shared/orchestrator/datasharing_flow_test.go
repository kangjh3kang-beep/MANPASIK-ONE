package orchestrator

import (
	"context"
	"testing"
	"time"
)

// --- Mock implementations for DataSharing Flow ---

type mockConsentManager struct {
	consent *ConsentInfo
	err     error
}

func (m *mockConsentManager) CheckConsent(ctx context.Context, userID, providerID, scope string) (bool, string, error) {
	if m.err != nil {
		return false, "", m.err
	}
	return true, m.consent.ConsentID, nil
}

func (m *mockConsentManager) GetConsent(ctx context.Context, consentID string) (*ConsentInfo, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.consent, nil
}

type mockRecordProvider struct {
	records []HealthRecordItem
	err     error
}

func (m *mockRecordProvider) GetRecordsByScope(ctx context.Context, userID string, scope []string, limit int) ([]HealthRecordItem, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.records, nil
}

type mockFHIRExporter struct {
	bundleJSON    string
	resourceCount int
	err           error
}

func (m *mockFHIRExporter) ExportBundle(ctx context.Context, records []HealthRecordItem, patientRef string) (string, int, error) {
	if m.err != nil {
		return "", 0, m.err
	}
	return m.bundleJSON, m.resourceCount, nil
}

type mockAuditLogger struct {
	err error
}

func (m *mockAuditLogger) LogAccess(ctx context.Context, userID, providerID, action, resourceType string, resourceIDs []string) error {
	return m.err
}

// --- Tests ---

func TestShareData_Success(t *testing.T) {
	cm := &mockConsentManager{consent: &ConsentInfo{
		ConsentID:  "consent-001",
		UserID:     "user-100",
		ProviderID: "provider-A",
		Scope:      []string{"blood_test", "vitals"},
		Status:     "active",
		ExpiresAt:  time.Now().Add(24 * time.Hour),
	}}
	rp := &mockRecordProvider{records: []HealthRecordItem{
		{RecordID: "rec-1", Type: "blood_test", Title: "Blood Test", DataJSON: `{"wbc":7.5}`, RecordedAt: time.Now()},
		{RecordID: "rec-2", Type: "vitals", Title: "Vitals", DataJSON: `{"hr":72}`, RecordedAt: time.Now()},
	}}
	fe := &mockFHIRExporter{bundleJSON: `{"resourceType":"Bundle","entry":[]}`, resourceCount: 3}
	al := &mockAuditLogger{}

	orchestrator := NewDataSharingFlowOrchestrator(cm, rp, fe, al)
	result, err := orchestrator.ShareData(context.Background(), "consent-001")

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if result.ConsentID != "consent-001" {
		t.Errorf("expected ConsentID=consent-001, got %s", result.ConsentID)
	}
	if result.ProviderID != "provider-A" {
		t.Errorf("expected ProviderID=provider-A, got %s", result.ProviderID)
	}
	if result.ResourceCount != 3 {
		t.Errorf("expected ResourceCount=3, got %d", result.ResourceCount)
	}
	if result.RecordCount != 2 {
		t.Errorf("expected RecordCount=2, got %d", result.RecordCount)
	}
	if result.FHIRBundleJSON == "" {
		t.Error("expected non-empty FHIR bundle JSON")
	}
}

func TestShareData_InactiveConsent(t *testing.T) {
	cm := &mockConsentManager{consent: &ConsentInfo{
		ConsentID:  "consent-002",
		UserID:     "user-200",
		ProviderID: "provider-B",
		Scope:      []string{"blood_test"},
		Status:     "revoked",
		ExpiresAt:  time.Now().Add(24 * time.Hour),
	}}
	rp := &mockRecordProvider{}
	fe := &mockFHIRExporter{}
	al := &mockAuditLogger{}

	orchestrator := NewDataSharingFlowOrchestrator(cm, rp, fe, al)
	result, err := orchestrator.ShareData(context.Background(), "consent-002")

	if err == nil {
		t.Fatal("expected error for inactive consent, got nil")
	}
	if result != nil {
		t.Errorf("expected nil result, got %+v", result)
	}
	expected := "동의가 활성 상태가 아닙니다: revoked"
	if err.Error() != expected {
		t.Errorf("expected error %q, got %q", expected, err.Error())
	}
}

func TestShareData_ExpiredConsent(t *testing.T) {
	cm := &mockConsentManager{consent: &ConsentInfo{
		ConsentID:  "consent-003",
		UserID:     "user-300",
		ProviderID: "provider-C",
		Scope:      []string{"vitals"},
		Status:     "active",
		ExpiresAt:  time.Now().Add(-1 * time.Hour), // expired 1 hour ago
	}}
	rp := &mockRecordProvider{}
	fe := &mockFHIRExporter{}
	al := &mockAuditLogger{}

	orchestrator := NewDataSharingFlowOrchestrator(cm, rp, fe, al)
	result, err := orchestrator.ShareData(context.Background(), "consent-003")

	if err == nil {
		t.Fatal("expected error for expired consent, got nil")
	}
	if result != nil {
		t.Errorf("expected nil result, got %+v", result)
	}
	if err.Error() != "동의 기간이 만료되었습니다" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}

func TestShareData_NoRecords(t *testing.T) {
	cm := &mockConsentManager{consent: &ConsentInfo{
		ConsentID:  "consent-004",
		UserID:     "user-400",
		ProviderID: "provider-D",
		Scope:      []string{"genetic"},
		Status:     "active",
		ExpiresAt:  time.Now().Add(24 * time.Hour),
	}}
	rp := &mockRecordProvider{records: []HealthRecordItem{}} // empty records
	fe := &mockFHIRExporter{}
	al := &mockAuditLogger{}

	orchestrator := NewDataSharingFlowOrchestrator(cm, rp, fe, al)
	result, err := orchestrator.ShareData(context.Background(), "consent-004")

	if err == nil {
		t.Fatal("expected error for no records, got nil")
	}
	if result != nil {
		t.Errorf("expected nil result, got %+v", result)
	}
	if err.Error() != "공유 가능한 건강 기록이 없습니다" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}
