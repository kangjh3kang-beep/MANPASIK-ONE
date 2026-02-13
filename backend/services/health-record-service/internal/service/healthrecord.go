// Package service는 health-record-service의 비즈니스 로직을 구현합니다.
package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	apperrors "github.com/manpasik/backend/shared/errors"
	"go.uber.org/zap"
)

// HealthRecordType은 건강 기록 유형입니다.
type HealthRecordType int

const (
	RecordTypeUnknown      HealthRecordType = iota
	RecordTypeLabResult                     // 검사 결과
	RecordTypeImaging                       // 영상 진단
	RecordTypeVitalSign                     // 활력징후
	RecordTypeAllergy                       // 알레르기
	RecordTypeCondition                     // 질환 이력
	RecordTypeImmunization                  // 예방접종
	RecordTypeProcedure                     // 시술/수술
)

// FHIRResourceType은 FHIR R4 리소스 유형입니다.
type FHIRResourceType int

const (
	FHIRUnknown              FHIRResourceType = iota
	FHIRObservation                           // Observation
	FHIRCondition                             // Condition
	FHIRMedicationStatement                   // MedicationStatement
	FHIRAllergyIntolerance                    // AllergyIntolerance
	FHIRImmunization                          // Immunization
	FHIRProcedure                             // Procedure
	FHIRDiagnosticReport                      // DiagnosticReport
	FHIRPatient                               // Patient
)

// HealthRecord는 건강 기록 도메인 객체입니다.
type HealthRecord struct {
	ID             string
	UserID         string
	RecordType     HealthRecordType
	Title          string
	Description    string
	Data           string
	Source         string // "manpasik", "manual", "fhir_import"
	FHIRResourceID string
	FHIRType       FHIRResourceType
	RecordedAt     time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// HealthRecordRepository는 건강 기록 저장소 인터페이스입니다.
type HealthRecordRepository interface {
	Save(ctx context.Context, r *HealthRecord) error
	FindByID(ctx context.Context, id string) (*HealthRecord, error)
	FindByUserID(ctx context.Context, userID string, typeFilter HealthRecordType, startDate, endDate *time.Time, limit, offset int) ([]*HealthRecord, int, error)
	Update(ctx context.Context, r *HealthRecord) error
	Delete(ctx context.Context, id string) error
}

// ============================================================================
// 데이터 공유 동의 (Data Sharing Consent) 도메인
// ============================================================================

// ConsentType determines what data is shared
type ConsentType string

const (
	ConsentMeasurementShare ConsentType = "measurement_share"
	ConsentRecordShare      ConsentType = "record_share"
	ConsentFullAccess       ConsentType = "full_access"
)

// ConsentStatus tracks consent lifecycle
type ConsentStatus string

const (
	ConsentActive  ConsentStatus = "active"
	ConsentRevoked ConsentStatus = "revoked"
	ConsentExpired ConsentStatus = "expired"
)

// DataSharingConsent represents a patient's consent to share data with a provider
type DataSharingConsent struct {
	ID           string
	UserID       string
	ProviderID   string // facility_id
	ProviderName string
	ConsentType  ConsentType
	Scope        []string // ["blood_glucose", "blood_pressure", ...]
	Purpose      string   // "treatment", "research", "emergency"
	Status       ConsentStatus
	GrantedAt    time.Time
	ExpiresAt    time.Time
	RevokedAt    time.Time
	RevokeReason string
}

// DataAccessLog records each time shared data is accessed
type DataAccessLog struct {
	ID           string
	ConsentID    string
	UserID       string
	ProviderID   string
	Action       string // "view", "export", "share"
	ResourceType string // "measurement", "health_record"
	ResourceIDs  []string
	AccessedAt   time.Time
	IPAddress    string
}

// SharedDataBundle is the result of sharing data with a provider
type SharedDataBundle struct {
	ConsentID      string
	ProviderID     string
	FHIRBundleJSON string
	ResourceCount  int
	SharedAt       time.Time
}

// ConsentRepository는 데이터 공유 동의 저장소 인터페이스입니다.
type ConsentRepository interface {
	Create(ctx context.Context, consent *DataSharingConsent) error
	GetByID(ctx context.Context, consentID string) (*DataSharingConsent, error)
	ListByUser(ctx context.Context, userID string) ([]*DataSharingConsent, error)
	ListByProvider(ctx context.Context, providerID string) ([]*DataSharingConsent, error)
	Revoke(ctx context.Context, consentID string, reason string) error
	CheckAccess(ctx context.Context, userID, providerID, scope string) (bool, error)
}

// DataAccessLogRepository는 데이터 접근 로그 저장소 인터페이스입니다.
type DataAccessLogRepository interface {
	Log(ctx context.Context, entry *DataAccessLog) error
	ListByUser(ctx context.Context, userID string, limit, offset int) ([]*DataAccessLog, int, error)
}

// HealthRecordService는 건강 기록 서비스 핵심 로직입니다.
type HealthRecordService struct {
	log            *zap.Logger
	repo           HealthRecordRepository
	consentRepo    ConsentRepository
	accessLogRepo  DataAccessLogRepository
}

// NewHealthRecordService는 HealthRecordService를 생성합니다.
func NewHealthRecordService(log *zap.Logger, repo HealthRecordRepository, consentRepo ConsentRepository, accessLogRepo DataAccessLogRepository) *HealthRecordService {
	return &HealthRecordService{log: log, repo: repo, consentRepo: consentRepo, accessLogRepo: accessLogRepo}
}

// CreateRecord는 건강 기록을 생성합니다.
func (s *HealthRecordService) CreateRecord(ctx context.Context, userID string, recordType HealthRecordType, title, description string, dataJson string, source string) (*HealthRecord, error) {
	if userID == "" || title == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "입력값이 올바르지 않습니다")
	}

	if source == "" {
		source = "manual"
	}

	now := time.Now()
	fhirType := mapRecordTypeToFHIR(recordType)

	record := &HealthRecord{
		ID:         uuid.New().String(),
		UserID:     userID,
		RecordType: recordType,
		Title:      title,
		Description: description,
		Data:       dataJson,
		Source:     source,
		FHIRType:   fhirType,
		RecordedAt: now,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := s.repo.Save(ctx, record); err != nil {
		return nil, fmt.Errorf("건강 기록 저장 실패: %w", err)
	}

	s.log.Info("건강 기록 생성",
		zap.String("record_id", record.ID),
		zap.String("user_id", userID),
		zap.Int("type", int(recordType)),
	)

	return record, nil
}

// GetRecord는 건강 기록을 조회합니다.
func (s *HealthRecordService) GetRecord(ctx context.Context, recordID string) (*HealthRecord, error) {
	if recordID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "입력값이 올바르지 않습니다")
	}
	return s.repo.FindByID(ctx, recordID)
}

// ListRecords는 건강 기록 목록을 조회합니다.
func (s *HealthRecordService) ListRecords(ctx context.Context, userID string, typeFilter HealthRecordType, limit, offset int) ([]*HealthRecord, int, error) {
	if userID == "" {
		return nil, 0, apperrors.New(apperrors.ErrInvalidInput, "입력값이 올바르지 않습니다")
	}
	if limit <= 0 {
		limit = 20
	}
	return s.repo.FindByUserID(ctx, userID, typeFilter, nil, nil, limit, offset)
}

// UpdateRecord는 건강 기록을 업데이트합니다.
func (s *HealthRecordService) UpdateRecord(ctx context.Context, recordID, title, description string, dataJson string) (*HealthRecord, error) {
	if recordID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "입력값이 올바르지 않습니다")
	}

	record, err := s.repo.FindByID(ctx, recordID)
	if err != nil {
		return nil, err
	}

	if title != "" {
		record.Title = title
	}
	if description != "" {
		record.Description = description
	}
	if dataJson != "" {
		record.Data = dataJson
	}
	record.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, record); err != nil {
		return nil, fmt.Errorf("건강 기록 업데이트 실패: %w", err)
	}

	return record, nil
}

// DeleteRecord는 건강 기록을 삭제합니다.
func (s *HealthRecordService) DeleteRecord(ctx context.Context, recordID string) error {
	if recordID == "" {
		return apperrors.New(apperrors.ErrInvalidInput, "입력값이 올바르지 않습니다")
	}
	return s.repo.Delete(ctx, recordID)
}

// ExportToFHIR는 건강 기록을 FHIR R4 Bundle JSON으로 내보냅니다.
func (s *HealthRecordService) ExportToFHIR(ctx context.Context, userID string, recordTypes []HealthRecordType, startDate, endDate *time.Time) (string, int, []FHIRResourceType, error) {
	if userID == "" {
		return "", 0, nil, apperrors.New(apperrors.ErrInvalidInput, "입력값이 올바르지 않습니다")
	}

	// 모든 기록 조회 (typeFilter 0 = 전체)
	records, _, err := s.repo.FindByUserID(ctx, userID, RecordTypeUnknown, startDate, endDate, 1000, 0)
	if err != nil {
		return "", 0, nil, err
	}

	// 타입 필터링
	var filtered []*HealthRecord
	typeSet := make(map[HealthRecordType]bool)
	for _, t := range recordTypes {
		typeSet[t] = true
	}
	for _, r := range records {
		if len(typeSet) == 0 || typeSet[r.RecordType] {
			filtered = append(filtered, r)
		}
	}

	// FHIR Bundle 생성 (간소화된 형태)
	bundle := map[string]interface{}{
		"resourceType": "Bundle",
		"type":         "collection",
		"total":        len(filtered),
		"entry":        make([]map[string]interface{}, 0, len(filtered)),
	}

	resourceTypes := make(map[FHIRResourceType]bool)
	entries := make([]map[string]interface{}, 0, len(filtered))

	for _, r := range filtered {
		fhirType := mapRecordTypeToFHIR(r.RecordType)
		resourceTypes[fhirType] = true

		entry := map[string]interface{}{
			"resource": map[string]interface{}{
				"resourceType": fhirResourceTypeName(fhirType),
				"id":           r.ID,
				"subject": map[string]interface{}{
					"reference": "Patient/" + r.UserID,
				},
				"effectiveDateTime": r.RecordedAt.Format(time.RFC3339),
				"code": map[string]interface{}{
					"text": r.Title,
				},
			},
		}
		entries = append(entries, entry)
	}
	bundle["entry"] = entries

	jsonBytes, err := json.Marshal(bundle)
	if err != nil {
		return "", 0, nil, fmt.Errorf("FHIR JSON 직렬화 실패: %w", err)
	}

	var types []FHIRResourceType
	for t := range resourceTypes {
		types = append(types, t)
	}

	return string(jsonBytes), len(filtered), types, nil
}

// ImportFromFHIR는 FHIR R4 Bundle JSON에서 건강 기록을 가져옵니다.
func (s *HealthRecordService) ImportFromFHIR(ctx context.Context, userID, bundleJSON string) ([]*HealthRecord, int, int, []string, error) {
	if userID == "" || bundleJSON == "" {
		return nil, 0, 0, nil, apperrors.New(apperrors.ErrInvalidInput, "입력값이 올바르지 않습니다")
	}

	var bundle map[string]interface{}
	if err := json.Unmarshal([]byte(bundleJSON), &bundle); err != nil {
		return nil, 0, 0, []string{"JSON 파싱 실패: " + err.Error()}, fmt.Errorf("FHIR JSON 파싱 실패: %w", err)
	}

	entries, ok := bundle["entry"].([]interface{})
	if !ok {
		return nil, 0, 0, []string{"entry 필드가 없거나 잘못됨"}, apperrors.New(apperrors.ErrInvalidInput, "입력값이 올바르지 않습니다")
	}

	var imported []*HealthRecord
	var errors []string
	skipped := 0

	for i, e := range entries {
		entry, ok := e.(map[string]interface{})
		if !ok {
			errors = append(errors, fmt.Sprintf("entry[%d]: 잘못된 형식", i))
			skipped++
			continue
		}

		resource, ok := entry["resource"].(map[string]interface{})
		if !ok {
			errors = append(errors, fmt.Sprintf("entry[%d]: resource 필드 없음", i))
			skipped++
			continue
		}

		resourceType, _ := resource["resourceType"].(string)
		fhirType := parseFHIRResourceType(resourceType)
		recordType := mapFHIRToRecordType(fhirType)

		title := "가져온 기록"
		if code, ok := resource["code"].(map[string]interface{}); ok {
			if text, ok := code["text"].(string); ok {
				title = text
			}
		}

		fhirID := ""
		if id, ok := resource["id"].(string); ok {
			fhirID = id
		}

		record := &HealthRecord{
			ID:             uuid.New().String(),
			UserID:         userID,
			RecordType:     recordType,
			Title:          title,
			Description:    "FHIR R4 가져오기",
			Data:           `{"fhir_resource_type":"` + resourceType + `"}`,
			Source:         "fhir_import",
			FHIRResourceID: fhirID,
			FHIRType:       fhirType,
			RecordedAt:     time.Now(),
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		if err := s.repo.Save(ctx, record); err != nil {
			errors = append(errors, fmt.Sprintf("entry[%d]: 저장 실패 - %v", i, err))
			skipped++
			continue
		}

		imported = append(imported, record)
	}

	return imported, len(imported), skipped, errors, nil
}

// GetHealthSummary는 건강 요약을 생성합니다.
func (s *HealthRecordService) GetHealthSummary(ctx context.Context, userID string, days int) (int, map[string]int, []*HealthRecord, string, error) {
	if userID == "" {
		return 0, nil, nil, "", apperrors.New(apperrors.ErrInvalidInput, "입력값이 올바르지 않습니다")
	}
	if days <= 0 {
		days = 30
	}

	startDate := time.Now().AddDate(0, 0, -days)
	records, total, err := s.repo.FindByUserID(ctx, userID, RecordTypeUnknown, &startDate, nil, 1000, 0)
	if err != nil {
		return 0, nil, nil, "", err
	}

	// 타입별 레코드 수 집계
	byType := make(map[string]int)
	for _, r := range records {
		typeName := recordTypeName(r.RecordType)
		byType[typeName]++
	}

	// 최근 5개 기록
	recent := records
	if len(recent) > 5 {
		recent = recent[:5]
	}

	// 요약 텍스트 생성
	summary := fmt.Sprintf("최근 %d일간 총 %d건의 건강 기록이 있습니다.", days, total)
	if total > 0 {
		summary += " 주요 기록 유형: "
		for k, v := range byType {
			summary += fmt.Sprintf("%s(%d건) ", k, v)
		}
	}

	return total, byType, recent, summary, nil
}

// ============================================================================
// 데이터 공유 동의 서비스 메서드
// ============================================================================

// CreateDataSharingConsent는 데이터 공유 동의를 생성합니다.
func (s *HealthRecordService) CreateDataSharingConsent(ctx context.Context, consent *DataSharingConsent) (*DataSharingConsent, error) {
	if consent.UserID == "" || consent.ProviderID == "" || consent.ConsentType == "" || len(consent.Scope) == 0 {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "UserID, ProviderID, ConsentType, Scope는 필수입니다")
	}

	consent.ID = uuid.New().String()
	consent.Status = ConsentActive
	consent.GrantedAt = time.Now()
	if consent.ExpiresAt.IsZero() {
		consent.ExpiresAt = consent.GrantedAt.AddDate(1, 0, 0) // 기본 1년
	}

	if err := s.consentRepo.Create(ctx, consent); err != nil {
		return nil, fmt.Errorf("동의 저장 실패: %w", err)
	}

	s.log.Info("데이터 공유 동의 생성",
		zap.String("consent_id", consent.ID),
		zap.String("user_id", consent.UserID),
		zap.String("provider_id", consent.ProviderID),
	)

	return consent, nil
}

// RevokeDataSharingConsent는 데이터 공유 동의를 철회합니다.
func (s *HealthRecordService) RevokeDataSharingConsent(ctx context.Context, consentID, reason string) error {
	if consentID == "" {
		return apperrors.New(apperrors.ErrInvalidInput, "consentID는 필수입니다")
	}

	consent, err := s.consentRepo.GetByID(ctx, consentID)
	if err != nil {
		return err
	}

	if consent.Status == ConsentRevoked {
		return apperrors.New(apperrors.ErrConflict, "이미 철회된 동의입니다")
	}

	if err := s.consentRepo.Revoke(ctx, consentID, reason); err != nil {
		return fmt.Errorf("동의 철회 실패: %w", err)
	}

	s.log.Info("데이터 공유 동의 철회",
		zap.String("consent_id", consentID),
		zap.String("reason", reason),
	)

	return nil
}

// ListDataSharingConsents는 사용자의 데이터 공유 동의 목록을 조회합니다.
func (s *HealthRecordService) ListDataSharingConsents(ctx context.Context, userID string) ([]*DataSharingConsent, error) {
	if userID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "userID는 필수입니다")
	}
	return s.consentRepo.ListByUser(ctx, userID)
}

// ShareWithProvider는 동의에 따라 데이터를 제공자에게 공유합니다.
func (s *HealthRecordService) ShareWithProvider(ctx context.Context, consentID string) (*SharedDataBundle, error) {
	if consentID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "consentID는 필수입니다")
	}

	consent, err := s.consentRepo.GetByID(ctx, consentID)
	if err != nil {
		return nil, err
	}

	if consent.Status != ConsentActive {
		return nil, apperrors.New(apperrors.ErrForbidden, "활성 상태가 아닌 동의입니다")
	}
	if time.Now().After(consent.ExpiresAt) {
		return nil, apperrors.New(apperrors.ErrForbidden, "만료된 동의입니다")
	}

	// 동의 scope에 따른 데이터 수집 — 모든 건강 기록 조회 후 FHIR 변환
	records, _, err := s.repo.FindByUserID(ctx, consent.UserID, RecordTypeUnknown, nil, nil, 1000, 0)
	if err != nil {
		return nil, fmt.Errorf("건강 기록 조회 실패: %w", err)
	}

	// scope에 맞는 FHIR Observation 생성
	var observations []map[string]interface{}
	var resourceIDs []string
	patientRef := "Patient/" + consent.UserID

	for _, r := range records {
		for _, scope := range consent.Scope {
			obs := MeasurementToFHIRObservation(scope, 0, "", patientRef, r.RecordedAt)
			observations = append(observations, obs)
			resourceIDs = append(resourceIDs, r.ID)
		}
	}

	bundleJSON, err := BuildFHIRBundle(observations)
	if err != nil {
		return nil, err
	}

	// 접근 로그 기록
	accessLog := &DataAccessLog{
		ID:           uuid.New().String(),
		ConsentID:    consentID,
		UserID:       consent.UserID,
		ProviderID:   consent.ProviderID,
		Action:       "share",
		ResourceType: "health_record",
		ResourceIDs:  resourceIDs,
		AccessedAt:   time.Now(),
	}
	if err := s.accessLogRepo.Log(ctx, accessLog); err != nil {
		s.log.Warn("접근 로그 기록 실패", zap.Error(err))
	}

	return &SharedDataBundle{
		ConsentID:      consentID,
		ProviderID:     consent.ProviderID,
		FHIRBundleJSON: bundleJSON,
		ResourceCount:  len(observations),
		SharedAt:       time.Now(),
	}, nil
}

// GetDataAccessLog는 사용자의 데이터 접근 로그를 조회합니다.
func (s *HealthRecordService) GetDataAccessLog(ctx context.Context, userID string, limit, offset int) ([]*DataAccessLog, int, error) {
	if userID == "" {
		return nil, 0, apperrors.New(apperrors.ErrInvalidInput, "userID는 필수입니다")
	}
	if limit <= 0 {
		limit = 20
	}
	return s.accessLogRepo.ListByUser(ctx, userID, limit, offset)
}

// ============================================================================
// 헬퍼 함수
// ============================================================================

func mapRecordTypeToFHIR(t HealthRecordType) FHIRResourceType {
	switch t {
	case RecordTypeVitalSign:
		return FHIRObservation
	case RecordTypeCondition:
		return FHIRCondition
	case RecordTypeLabResult, RecordTypeImaging:
		return FHIRDiagnosticReport
	case RecordTypeAllergy:
		return FHIRAllergyIntolerance
	case RecordTypeImmunization:
		return FHIRImmunization
	case RecordTypeProcedure:
		return FHIRProcedure
	default:
		return FHIRObservation
	}
}

func mapFHIRToRecordType(f FHIRResourceType) HealthRecordType {
	switch f {
	case FHIRObservation:
		return RecordTypeVitalSign
	case FHIRCondition:
		return RecordTypeCondition
	case FHIRMedicationStatement:
		return RecordTypeUnknown
	case FHIRAllergyIntolerance:
		return RecordTypeAllergy
	case FHIRImmunization:
		return RecordTypeImmunization
	case FHIRProcedure:
		return RecordTypeProcedure
	case FHIRDiagnosticReport:
		return RecordTypeLabResult
	default:
		return RecordTypeUnknown
	}
}

func fhirResourceTypeName(t FHIRResourceType) string {
	switch t {
	case FHIRObservation:
		return "Observation"
	case FHIRCondition:
		return "Condition"
	case FHIRMedicationStatement:
		return "MedicationStatement"
	case FHIRAllergyIntolerance:
		return "AllergyIntolerance"
	case FHIRImmunization:
		return "Immunization"
	case FHIRProcedure:
		return "Procedure"
	case FHIRDiagnosticReport:
		return "DiagnosticReport"
	case FHIRPatient:
		return "Patient"
	default:
		return "Unknown"
	}
}

func parseFHIRResourceType(name string) FHIRResourceType {
	switch name {
	case "Observation":
		return FHIRObservation
	case "Condition":
		return FHIRCondition
	case "MedicationStatement":
		return FHIRMedicationStatement
	case "AllergyIntolerance":
		return FHIRAllergyIntolerance
	case "Immunization":
		return FHIRImmunization
	case "Procedure":
		return FHIRProcedure
	case "DiagnosticReport":
		return FHIRDiagnosticReport
	case "Patient":
		return FHIRPatient
	default:
		return FHIRUnknown
	}
}

func recordTypeName(t HealthRecordType) string {
	switch t {
	case RecordTypeLabResult:
		return "lab_result"
	case RecordTypeImaging:
		return "imaging"
	case RecordTypeVitalSign:
		return "vital_sign"
	case RecordTypeAllergy:
		return "allergy"
	case RecordTypeCondition:
		return "condition"
	case RecordTypeImmunization:
		return "immunization"
	case RecordTypeProcedure:
		return "procedure"
	default:
		return "unknown"
	}
}
