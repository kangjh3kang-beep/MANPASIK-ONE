package orchestrator

import (
	"context"
	"fmt"
	"log"
	"time"
)

// DataSharingFlowOrchestrator coordinates consent → health record → FHIR export
type DataSharingFlowOrchestrator struct {
	consentManager ConsentManager
	recordProvider RecordProvider
	fhirExporter   FHIRExporter
	auditLogger    AuditLogger
}

// ConsentManager manages data sharing consents
type ConsentManager interface {
	CheckConsent(ctx context.Context, userID, providerID, scope string) (bool, string, error)
	GetConsent(ctx context.Context, consentID string) (*ConsentInfo, error)
}

// ConsentInfo represents consent details
type ConsentInfo struct {
	ConsentID  string
	UserID     string
	ProviderID string
	Scope      []string
	Status     string
	ExpiresAt  time.Time
}

// RecordProvider provides health records
type RecordProvider interface {
	GetRecordsByScope(ctx context.Context, userID string, scope []string, limit int) ([]HealthRecordItem, error)
}

// HealthRecordItem represents a health record for sharing
type HealthRecordItem struct {
	RecordID   string
	Type       string
	Title      string
	DataJSON   string
	RecordedAt time.Time
}

// FHIRExporter exports data as FHIR bundle
type FHIRExporter interface {
	ExportBundle(ctx context.Context, records []HealthRecordItem, patientRef string) (string, int, error)
}

// AuditLogger logs data access events
type AuditLogger interface {
	LogAccess(ctx context.Context, userID, providerID, action, resourceType string, resourceIDs []string) error
}

// NewDataSharingFlowOrchestrator creates a new data sharing orchestrator
func NewDataSharingFlowOrchestrator(cm ConsentManager, rp RecordProvider, fe FHIRExporter, al AuditLogger) *DataSharingFlowOrchestrator {
	return &DataSharingFlowOrchestrator{
		consentManager: cm,
		recordProvider: rp,
		fhirExporter:   fe,
		auditLogger:    al,
	}
}

// ShareData executes the full data sharing flow for a consent
func (o *DataSharingFlowOrchestrator) ShareData(ctx context.Context, consentID string) (*SharingResult, error) {
	log.Printf("[DataSharing] Processing consent: %s", consentID)

	// Step 1: Verify consent
	consent, err := o.consentManager.GetConsent(ctx, consentID)
	if err != nil {
		return nil, fmt.Errorf("동의 정보 조회 실패: %w", err)
	}
	if consent.Status != "active" {
		return nil, fmt.Errorf("동의가 활성 상태가 아닙니다: %s", consent.Status)
	}
	if time.Now().After(consent.ExpiresAt) {
		return nil, fmt.Errorf("동의 기간이 만료되었습니다")
	}
	log.Printf("[DataSharing] Step 1: Consent verified (scope=%v)", consent.Scope)

	// Step 2: Get health records matching scope
	records, err := o.recordProvider.GetRecordsByScope(ctx, consent.UserID, consent.Scope, 100)
	if err != nil {
		return nil, fmt.Errorf("건강 기록 조회 실패: %w", err)
	}
	if len(records) == 0 {
		return nil, fmt.Errorf("공유 가능한 건강 기록이 없습니다")
	}
	log.Printf("[DataSharing] Step 2: Found %d records", len(records))

	// Step 3: Export as FHIR bundle
	patientRef := fmt.Sprintf("Patient/%s", consent.UserID)
	bundleJSON, resourceCount, err := o.fhirExporter.ExportBundle(ctx, records, patientRef)
	if err != nil {
		return nil, fmt.Errorf("FHIR 번들 생성 실패: %w", err)
	}
	log.Printf("[DataSharing] Step 3: FHIR bundle created (%d resources)", resourceCount)

	// Step 4: Log the access
	recordIDs := make([]string, len(records))
	for i, r := range records {
		recordIDs[i] = r.RecordID
	}
	if err := o.auditLogger.LogAccess(ctx, consent.UserID, consent.ProviderID, "share", "health_record", recordIDs); err != nil {
		log.Printf("[DataSharing] Audit log failed (non-fatal): %v", err)
	}
	log.Printf("[DataSharing] Step 4: Access logged")

	return &SharingResult{
		ConsentID:      consentID,
		ProviderID:     consent.ProviderID,
		FHIRBundleJSON: bundleJSON,
		ResourceCount:  resourceCount,
		RecordCount:    len(records),
		SharedAt:       time.Now(),
	}, nil
}

// SharingResult is the result of a data sharing operation
type SharingResult struct {
	ConsentID      string
	ProviderID     string
	FHIRBundleJSON string
	ResourceCount  int
	RecordCount    int
	SharedAt       time.Time
}
