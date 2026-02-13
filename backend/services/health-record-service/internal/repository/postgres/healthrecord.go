// Package postgres는 health-record-service의 PostgreSQL 저장소 구현입니다.
package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manpasik/backend/services/health-record-service/internal/service"
	apperrors "github.com/manpasik/backend/shared/errors"
)

// ============================================================================
// health_record_type / fhir_resource_type ↔ Go enum 매핑
// ============================================================================

var recordTypeToString = map[service.HealthRecordType]string{
	service.RecordTypeLabResult:    "lab_result",
	service.RecordTypeImaging:      "imaging",
	service.RecordTypeVitalSign:    "vital_sign",
	service.RecordTypeAllergy:      "allergy",
	service.RecordTypeCondition:    "condition",
	service.RecordTypeImmunization: "immunization",
	service.RecordTypeProcedure:    "procedure",
}

var stringToRecordType = map[string]service.HealthRecordType{
	"measurement":  service.RecordTypeUnknown,
	"medication":   service.RecordTypeUnknown,
	"symptom":      service.RecordTypeUnknown,
	"vital_sign":   service.RecordTypeVitalSign,
	"lab_result":   service.RecordTypeLabResult,
	"allergy":      service.RecordTypeAllergy,
	"condition":    service.RecordTypeCondition,
	"immunization": service.RecordTypeImmunization,
	"procedure":    service.RecordTypeProcedure,
	"note":         service.RecordTypeUnknown,
	"imaging":      service.RecordTypeImaging,
}

var fhirTypeToString = map[service.FHIRResourceType]string{
	service.FHIRObservation:          "observation",
	service.FHIRCondition:            "condition",
	service.FHIRMedicationStatement:  "medication_statement",
	service.FHIRAllergyIntolerance:   "allergy_intolerance",
	service.FHIRImmunization:         "immunization",
	service.FHIRProcedure:            "procedure",
	service.FHIRDiagnosticReport:     "diagnostic_report",
	service.FHIRPatient:              "patient",
}

var stringToFHIRType = map[string]service.FHIRResourceType{
	"observation":           service.FHIRObservation,
	"condition":             service.FHIRCondition,
	"medication_statement":  service.FHIRMedicationStatement,
	"allergy_intolerance":   service.FHIRAllergyIntolerance,
	"immunization":          service.FHIRImmunization,
	"procedure":             service.FHIRProcedure,
	"diagnostic_report":     service.FHIRDiagnosticReport,
	"patient":               service.FHIRPatient,
}

// recordTypeDBValue returns a DB-safe string for the record type enum.
// Falls back to "note" for unknown types.
func recordTypeDBValue(t service.HealthRecordType) string {
	if s, ok := recordTypeToString[t]; ok {
		return s
	}
	return "note"
}

// ============================================================================
// HealthRecordRepository
// ============================================================================

// HealthRecordRepository는 PostgreSQL 기반 건강 기록 저장소입니다.
type HealthRecordRepository struct {
	pool *pgxpool.Pool
}

// NewHealthRecordRepository는 HealthRecordRepository를 생성합니다.
func NewHealthRecordRepository(pool *pgxpool.Pool) *HealthRecordRepository {
	return &HealthRecordRepository{pool: pool}
}

// Save는 건강 기록을 저장합니다.
func (r *HealthRecordRepository) Save(ctx context.Context, rec *service.HealthRecord) error {
	const q = `INSERT INTO health_records
		(id, user_id, record_type, title, description, data, source,
		 fhir_resource_id, fhir_type, recorded_at, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`

	_, err := r.pool.Exec(ctx, q,
		rec.ID, rec.UserID, recordTypeDBValue(rec.RecordType), rec.Title, rec.Description,
		rec.Data, rec.Source,
		nullIfEmpty(rec.FHIRResourceID), nullFHIRType(rec.FHIRType),
		rec.RecordedAt, rec.CreatedAt, rec.UpdatedAt,
	)
	return err
}

// FindByID는 ID로 건강 기록을 조회합니다.
func (r *HealthRecordRepository) FindByID(ctx context.Context, id string) (*service.HealthRecord, error) {
	const q = `SELECT id, user_id, record_type::text, title, COALESCE(description,''),
		COALESCE(data::text,'{}'), COALESCE(source,'manual'),
		COALESCE(fhir_resource_id,''), COALESCE(fhir_type::text,''),
		recorded_at, created_at, updated_at
		FROM health_records WHERE id = $1 AND deleted_at IS NULL`

	rec, err := r.scanOne(ctx, q, id)
	if err != nil {
		return nil, err
	}
	if rec == nil {
		return nil, apperrors.New(apperrors.ErrNotFound, "리소스를 찾을 수 없습니다")
	}
	return rec, nil
}

// FindByUserID는 사용자 건강 기록 목록을 조회합니다.
func (r *HealthRecordRepository) FindByUserID(ctx context.Context, userID string, typeFilter service.HealthRecordType, startDate, endDate *time.Time, limit, offset int) ([]*service.HealthRecord, int, error) {
	where := `WHERE user_id = $1 AND deleted_at IS NULL`
	args := []interface{}{userID}
	idx := 2

	if typeFilter != service.RecordTypeUnknown {
		where += fmt.Sprintf(` AND record_type = $%d`, idx)
		args = append(args, recordTypeDBValue(typeFilter))
		idx++
	}
	if startDate != nil {
		where += fmt.Sprintf(` AND recorded_at >= $%d`, idx)
		args = append(args, *startDate)
		idx++
	}
	if endDate != nil {
		where += fmt.Sprintf(` AND recorded_at <= $%d`, idx)
		args = append(args, *endDate)
		idx++
	}

	// count
	var total int
	if err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM health_records `+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// data
	dataQ := fmt.Sprintf(`SELECT id, user_id, record_type::text, title, COALESCE(description,''),
		COALESCE(data::text,'{}'), COALESCE(source,'manual'),
		COALESCE(fhir_resource_id,''), COALESCE(fhir_type::text,''),
		recorded_at, created_at, updated_at
		FROM health_records %s ORDER BY recorded_at DESC LIMIT $%d OFFSET $%d`,
		where, idx, idx+1)
	dataArgs := append(args, limit, offset)

	rows, err := r.pool.Query(ctx, dataQ, dataArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []*service.HealthRecord
	for rows.Next() {
		rec, err := scanHealthRecordRow(rows)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, rec)
	}
	return list, total, nil
}

// Update는 건강 기록을 업데이트합니다.
func (r *HealthRecordRepository) Update(ctx context.Context, rec *service.HealthRecord) error {
	const q = `UPDATE health_records SET
		title=$1, description=$2, data=$3, source=$4,
		fhir_resource_id=$5, fhir_type=$6, updated_at=$7
		WHERE id=$8 AND deleted_at IS NULL`

	tag, err := r.pool.Exec(ctx, q,
		rec.Title, rec.Description, rec.Data, rec.Source,
		nullIfEmpty(rec.FHIRResourceID), nullFHIRType(rec.FHIRType), rec.UpdatedAt,
		rec.ID,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return apperrors.New(apperrors.ErrNotFound, "리소스를 찾을 수 없습니다")
	}
	return nil
}

// Delete는 건강 기록을 소프트 삭제합니다.
func (r *HealthRecordRepository) Delete(ctx context.Context, id string) error {
	const q = `UPDATE health_records SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`
	tag, err := r.pool.Exec(ctx, q, time.Now().UTC(), id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return apperrors.New(apperrors.ErrNotFound, "리소스를 찾을 수 없습니다")
	}
	return nil
}

// ============================================================================
// ConsentRepository
// ============================================================================

// ConsentRepository는 PostgreSQL 기반 데이터 공유 동의 저장소입니다.
type ConsentRepository struct {
	pool *pgxpool.Pool
}

// NewConsentRepository는 ConsentRepository를 생성합니다.
func NewConsentRepository(pool *pgxpool.Pool) *ConsentRepository {
	return &ConsentRepository{pool: pool}
}

// Create는 데이터 공유 동의를 저장합니다.
func (r *ConsentRepository) Create(ctx context.Context, consent *service.DataSharingConsent) error {
	const q = `INSERT INTO data_sharing_consents
		(id, user_id, facility_id, purpose, data_scope, is_active, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`

	purpose := "treatment"
	if consent.Purpose != "" {
		purpose = consent.Purpose
	}
	dataScope := "last_30_days"
	if len(consent.Scope) > 0 {
		dataScope = strings.Join(consent.Scope, ",")
	}

	_, err := r.pool.Exec(ctx, q,
		consent.ID, consent.UserID, consent.ProviderID,
		purpose, dataScope,
		consent.Status == service.ConsentActive,
		consent.GrantedAt, consent.GrantedAt,
	)
	return err
}

// GetByID는 ID로 동의를 조회합니다.
func (r *ConsentRepository) GetByID(ctx context.Context, consentID string) (*service.DataSharingConsent, error) {
	const q = `SELECT id, user_id, facility_id, COALESCE(doctor_id,''),
		purpose::text, COALESCE(data_scope,''), is_active,
		revoked_at, created_at, updated_at
		FROM data_sharing_consents WHERE id = $1`

	c, err := r.scanConsent(ctx, q, consentID)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, apperrors.New(apperrors.ErrNotFound, "동의를 찾을 수 없습니다")
	}
	return c, nil
}

// ListByUser는 사용자의 동의 목록을 조회합니다.
func (r *ConsentRepository) ListByUser(ctx context.Context, userID string) ([]*service.DataSharingConsent, error) {
	const q = `SELECT id, user_id, facility_id, COALESCE(doctor_id,''),
		purpose::text, COALESCE(data_scope,''), is_active,
		revoked_at, created_at, updated_at
		FROM data_sharing_consents WHERE user_id = $1 ORDER BY created_at DESC`

	return r.scanConsentList(ctx, q, userID)
}

// ListByProvider는 제공자의 동의 목록을 조회합니다.
func (r *ConsentRepository) ListByProvider(ctx context.Context, providerID string) ([]*service.DataSharingConsent, error) {
	const q = `SELECT id, user_id, facility_id, COALESCE(doctor_id,''),
		purpose::text, COALESCE(data_scope,''), is_active,
		revoked_at, created_at, updated_at
		FROM data_sharing_consents WHERE facility_id = $1 ORDER BY created_at DESC`

	return r.scanConsentList(ctx, q, providerID)
}

// Revoke는 동의를 철회합니다.
func (r *ConsentRepository) Revoke(ctx context.Context, consentID string, reason string) error {
	now := time.Now().UTC()
	const q = `UPDATE data_sharing_consents SET is_active = false, revoked_at = $1, updated_at = $2 WHERE id = $3`
	tag, err := r.pool.Exec(ctx, q, now, now, consentID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return apperrors.New(apperrors.ErrNotFound, "동의를 찾을 수 없습니다")
	}
	return nil
}

// CheckAccess는 특정 사용자-제공자-scope 조합에 대한 접근 권한을 확인합니다.
func (r *ConsentRepository) CheckAccess(ctx context.Context, userID, providerID, scope string) (bool, error) {
	const q = `SELECT COUNT(*) FROM data_sharing_consents
		WHERE user_id = $1 AND facility_id = $2 AND is_active = true
		AND (data_scope LIKE '%' || $3 || '%' OR data_scope = 'all')`

	var cnt int
	if err := r.pool.QueryRow(ctx, q, userID, providerID, scope).Scan(&cnt); err != nil {
		return false, err
	}
	return cnt > 0, nil
}

// ============================================================================
// DataAccessLogRepository
// ============================================================================

// DataAccessLogRepository는 PostgreSQL 기반 데이터 접근 로그 저장소입니다.
type DataAccessLogRepository struct {
	pool *pgxpool.Pool
}

// NewDataAccessLogRepository는 DataAccessLogRepository를 생성합니다.
func NewDataAccessLogRepository(pool *pgxpool.Pool) *DataAccessLogRepository {
	return &DataAccessLogRepository{pool: pool}
}

// Log는 데이터 접근 로그를 기록합니다.
func (r *DataAccessLogRepository) Log(ctx context.Context, entry *service.DataAccessLog) error {
	const q = `INSERT INTO shared_data_access_logs
		(id, user_id, facility_id, doctor_id, consent_id, access_type, resource_count, accessed_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`

	_, err := r.pool.Exec(ctx, q,
		entry.ID, entry.UserID, entry.ProviderID, nullIfEmpty(entry.Action),
		nullIfEmpty(entry.ConsentID), entry.ResourceType, len(entry.ResourceIDs),
		entry.AccessedAt,
	)
	return err
}

// ListByUser는 사용자의 접근 로그를 조회합니다.
func (r *DataAccessLogRepository) ListByUser(ctx context.Context, userID string, limit, offset int) ([]*service.DataAccessLog, int, error) {
	// count
	var total int
	if err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM shared_data_access_logs WHERE user_id = $1`, userID).Scan(&total); err != nil {
		return nil, 0, err
	}

	const q = `SELECT id, COALESCE(consent_id::text,''), user_id, facility_id,
		COALESCE(doctor_id,''), COALESCE(access_type,''), COALESCE(resource_count,0), accessed_at
		FROM shared_data_access_logs WHERE user_id = $1
		ORDER BY accessed_at DESC LIMIT $2 OFFSET $3`

	rows, err := r.pool.Query(ctx, q, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []*service.DataAccessLog
	for rows.Next() {
		var entry service.DataAccessLog
		var resourceCount int
		if err := rows.Scan(&entry.ID, &entry.ConsentID, &entry.UserID, &entry.ProviderID,
			&entry.Action, &entry.ResourceType, &resourceCount, &entry.AccessedAt); err != nil {
			return nil, 0, err
		}
		list = append(list, &entry)
	}
	return list, total, nil
}

// ============================================================================
// internal helpers
// ============================================================================

func (r *HealthRecordRepository) scanOne(ctx context.Context, query string, args ...interface{}) (*service.HealthRecord, error) {
	var rec service.HealthRecord
	var recordTypeStr, fhirTypeStr string

	err := r.pool.QueryRow(ctx, query, args...).Scan(
		&rec.ID, &rec.UserID, &recordTypeStr, &rec.Title, &rec.Description,
		&rec.Data, &rec.Source,
		&rec.FHIRResourceID, &fhirTypeStr,
		&rec.RecordedAt, &rec.CreatedAt, &rec.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	rec.RecordType = stringToRecordType[recordTypeStr]
	rec.FHIRType = stringToFHIRType[fhirTypeStr]
	return &rec, nil
}

func scanHealthRecordRow(rows pgx.Rows) (*service.HealthRecord, error) {
	var rec service.HealthRecord
	var recordTypeStr, fhirTypeStr string

	err := rows.Scan(
		&rec.ID, &rec.UserID, &recordTypeStr, &rec.Title, &rec.Description,
		&rec.Data, &rec.Source,
		&rec.FHIRResourceID, &fhirTypeStr,
		&rec.RecordedAt, &rec.CreatedAt, &rec.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	rec.RecordType = stringToRecordType[recordTypeStr]
	rec.FHIRType = stringToFHIRType[fhirTypeStr]
	return &rec, nil
}

func (r *ConsentRepository) scanConsent(ctx context.Context, query string, args ...interface{}) (*service.DataSharingConsent, error) {
	var c service.DataSharingConsent
	var purposeStr, dataScope string
	var isActive bool
	var revokedAt *time.Time
	var createdAt, updatedAt time.Time
	var doctorID string

	err := r.pool.QueryRow(ctx, query, args...).Scan(
		&c.ID, &c.UserID, &c.ProviderID, &doctorID,
		&purposeStr, &dataScope, &isActive,
		&revokedAt, &createdAt, &updatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	c.ProviderName = doctorID
	c.ConsentType = service.ConsentType(purposeStr)
	c.Purpose = purposeStr
	if dataScope != "" {
		c.Scope = strings.Split(dataScope, ",")
	}
	if isActive {
		c.Status = service.ConsentActive
	} else {
		c.Status = service.ConsentRevoked
	}
	c.GrantedAt = createdAt
	c.ExpiresAt = createdAt.AddDate(1, 0, 0)
	if revokedAt != nil {
		c.RevokedAt = *revokedAt
	}
	return &c, nil
}

func (r *ConsentRepository) scanConsentList(ctx context.Context, query string, args ...interface{}) ([]*service.DataSharingConsent, error) {
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*service.DataSharingConsent
	for rows.Next() {
		var c service.DataSharingConsent
		var purposeStr, dataScope string
		var isActive bool
		var revokedAt *time.Time
		var createdAt, updatedAt time.Time
		var doctorID string

		if err := rows.Scan(
			&c.ID, &c.UserID, &c.ProviderID, &doctorID,
			&purposeStr, &dataScope, &isActive,
			&revokedAt, &createdAt, &updatedAt,
		); err != nil {
			return nil, err
		}

		c.ProviderName = doctorID
		c.ConsentType = service.ConsentType(purposeStr)
		c.Purpose = purposeStr
		if dataScope != "" {
			c.Scope = strings.Split(dataScope, ",")
		}
		if isActive {
			c.Status = service.ConsentActive
		} else {
			c.Status = service.ConsentRevoked
		}
		c.GrantedAt = createdAt
		c.ExpiresAt = createdAt.AddDate(1, 0, 0)
		if revokedAt != nil {
			c.RevokedAt = *revokedAt
		}
		list = append(list, &c)
	}
	return list, nil
}

func nullIfEmpty(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

func nullFHIRType(t service.FHIRResourceType) interface{} {
	if s, ok := fhirTypeToString[t]; ok {
		return s
	}
	return nil
}
