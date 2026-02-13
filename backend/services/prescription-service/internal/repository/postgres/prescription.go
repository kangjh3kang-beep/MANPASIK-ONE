// Package postgres는 prescription-service의 PostgreSQL 저장소 구현입니다.
package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manpasik/backend/services/prescription-service/internal/service"
	apperrors "github.com/manpasik/backend/shared/errors"
)

// ============================================================================
// prescription_status / drug_interaction_severity  ↔  Go enum 매핑
// ============================================================================

var statusToString = map[service.PrescriptionStatus]string{
	service.StatusDraft:     "draft",
	service.StatusActive:    "active",
	service.StatusDispensed: "dispensed",
	service.StatusCompleted: "completed",
	service.StatusCancelled: "cancelled",
	service.StatusExpired:   "expired",
}

var stringToStatus = map[string]service.PrescriptionStatus{
	"draft":     service.StatusDraft,
	"active":    service.StatusActive,
	"dispensed": service.StatusDispensed,
	"completed": service.StatusCompleted,
	"cancelled": service.StatusCancelled,
	"expired":   service.StatusExpired,
}

var stringToSeverity = map[string]service.InteractionSeverity{
	"none":            service.SeverityNone,
	"minor":           service.SeverityMinor,
	"moderate":        service.SeverityModerate,
	"major":           service.SeverityMajor,
	"contraindicated": service.SeverityContraindicated,
}

// ============================================================================
// PrescriptionRepository
// ============================================================================

// PrescriptionRepository는 PostgreSQL 기반 PrescriptionRepository 구현입니다.
type PrescriptionRepository struct {
	pool *pgxpool.Pool
}

// NewPrescriptionRepository는 PrescriptionRepository를 생성합니다.
func NewPrescriptionRepository(pool *pgxpool.Pool) *PrescriptionRepository {
	return &PrescriptionRepository{pool: pool}
}

// prescriptionColumns는 처방전 SELECT에서 사용하는 공통 컬럼입니다.
const prescriptionColumns = `id, patient_user_id, doctor_id, COALESCE(consultation_id::text,''), status::text,
	COALESCE(diagnosis,''), COALESCE(notes,''),
	COALESCE(pharmacy_id::text,''), COALESCE(fulfillment_type,''), COALESCE(shipping_address_id,''),
	COALESCE(fulfillment_token,''), COALESCE(dispensary_status,''),
	prescribed_at, expires_at, sent_to_pharmacy_at, dispensed_at, created_at`

// Save는 처방전과 연관 약물을 저장합니다.
func (r *PrescriptionRepository) Save(ctx context.Context, p *service.Prescription) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	statusStr := statusToString[p.Status]

	const q = `INSERT INTO prescriptions
		(id, patient_user_id, doctor_id, consultation_id, status, diagnosis, notes,
		 pharmacy_id, fulfillment_type, shipping_address_id, fulfillment_token, dispensary_status,
		 prescribed_at, expires_at, sent_to_pharmacy_at, dispensed_at, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$17)`

	_, err = tx.Exec(ctx, q,
		p.ID, p.PatientUserID, p.DoctorID, nullIfEmpty(p.ConsultationID),
		statusStr, p.Diagnosis, p.Notes,
		nullIfEmpty(p.PharmacyID), nullIfEmpty(string(p.FulfillmentType)), nullIfEmpty(p.ShippingAddress),
		nullIfEmpty(p.FulfillmentToken), nullIfEmpty(string(p.DispensaryStatus)),
		p.PrescribedAt, p.ExpiresAt, nullTimeIfZero(p.SentToPharmacyAt), nullTimeIfZero(p.DispensedAt),
		p.CreatedAt,
	)
	if err != nil {
		return err
	}

	for _, m := range p.Medications {
		if err := insertMedication(ctx, tx, p.ID, m); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

// FindByID는 ID로 처방전을 조회합니다.
func (r *PrescriptionRepository) FindByID(ctx context.Context, id string) (*service.Prescription, error) {
	q := `SELECT ` + prescriptionColumns + ` FROM prescriptions WHERE id = $1`
	p, err := r.scanOne(ctx, q, id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, apperrors.New(apperrors.ErrNotFound, "처방전을 찾을 수 없습니다")
	}
	meds, err := r.loadMedications(ctx, p.ID)
	if err != nil {
		return nil, err
	}
	p.Medications = meds
	return p, nil
}

// FindByUserID는 환자의 처방전 목록을 조회합니다 (paginated).
func (r *PrescriptionRepository) FindByUserID(ctx context.Context, userID string, statusFilter service.PrescriptionStatus, limit, offset int) ([]*service.Prescription, int, error) {
	where := `WHERE patient_user_id = $1`
	args := []interface{}{userID}
	idx := 2

	if statusFilter != service.StatusUnknown {
		where += fmt.Sprintf(` AND status = $%d`, idx)
		args = append(args, statusToString[statusFilter])
		idx++
	}

	// count
	var total int
	if err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM prescriptions `+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// data
	dataQ := fmt.Sprintf(`SELECT %s FROM prescriptions %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d`,
		prescriptionColumns, where, idx, idx+1)
	dataArgs := append(args, limit, offset)

	rows, err := r.pool.Query(ctx, dataQ, dataArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []*service.Prescription
	for rows.Next() {
		p, err := scanPrescriptionRow(rows)
		if err != nil {
			return nil, 0, err
		}
		meds, err := r.loadMedications(ctx, p.ID)
		if err != nil {
			return nil, 0, err
		}
		p.Medications = meds
		list = append(list, p)
	}

	return list, total, nil
}

// Update는 처방전 정보를 업데이트합니다.
func (r *PrescriptionRepository) Update(ctx context.Context, p *service.Prescription) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	statusStr := statusToString[p.Status]

	const q = `UPDATE prescriptions SET
		status=$1, diagnosis=$2, notes=$3,
		pharmacy_id=$4, fulfillment_type=$5, shipping_address_id=$6,
		fulfillment_token=$7, dispensary_status=$8,
		sent_to_pharmacy_at=$9, dispensed_at=$10, updated_at=$11
		WHERE id=$12`

	_, err = tx.Exec(ctx, q,
		statusStr, p.Diagnosis, p.Notes,
		nullIfEmpty(p.PharmacyID), nullIfEmpty(string(p.FulfillmentType)), nullIfEmpty(p.ShippingAddress),
		nullIfEmpty(p.FulfillmentToken), nullIfEmpty(string(p.DispensaryStatus)),
		nullTimeIfZero(p.SentToPharmacyAt), nullTimeIfZero(p.DispensedAt), time.Now().UTC(),
		p.ID,
	)
	if err != nil {
		return err
	}

	// 약물 목록 교체: 삭제 후 재삽입
	if _, err := tx.Exec(ctx, `DELETE FROM prescription_medications WHERE prescription_id = $1`, p.ID); err != nil {
		return err
	}
	for _, m := range p.Medications {
		if err := insertMedication(ctx, tx, p.ID, m); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

// FindByPharmacyID는 약국 ID로 처방전 목록을 조회합니다.
func (r *PrescriptionRepository) FindByPharmacyID(ctx context.Context, pharmacyID string) ([]*service.Prescription, error) {
	q := `SELECT ` + prescriptionColumns + ` FROM prescriptions WHERE pharmacy_id = $1 ORDER BY created_at DESC`

	rows, err := r.pool.Query(ctx, q, pharmacyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*service.Prescription
	for rows.Next() {
		p, err := scanPrescriptionRow(rows)
		if err != nil {
			return nil, err
		}
		meds, err := r.loadMedications(ctx, p.ID)
		if err != nil {
			return nil, err
		}
		p.Medications = meds
		list = append(list, p)
	}
	return list, nil
}

// FindByFulfillmentToken은 조제 토큰으로 처방전을 조회합니다.
func (r *PrescriptionRepository) FindByFulfillmentToken(ctx context.Context, token string) (*service.Prescription, error) {
	q := `SELECT ` + prescriptionColumns + ` FROM prescriptions WHERE fulfillment_token = $1`
	p, err := r.scanOne(ctx, q, token)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, apperrors.New(apperrors.ErrNotFound, "해당 토큰의 처방전을 찾을 수 없습니다")
	}
	meds, err := r.loadMedications(ctx, p.ID)
	if err != nil {
		return nil, err
	}
	p.Medications = meds
	return p, nil
}

// ============================================================================
// TokenRepository
// ============================================================================

// TokenRepository는 PostgreSQL 기반 조제 토큰 저장소입니다.
// fulfillment_token은 prescriptions 테이블의 컬럼으로 관리됩니다.
type TokenRepository struct {
	pool *pgxpool.Pool
}

// NewTokenRepository는 TokenRepository를 생성합니다.
func NewTokenRepository(pool *pgxpool.Pool) *TokenRepository {
	return &TokenRepository{pool: pool}
}

// Create는 조제 토큰을 처방전에 연결합니다.
func (r *TokenRepository) Create(ctx context.Context, token *service.FulfillmentToken) error {
	const q = `UPDATE prescriptions SET fulfillment_token=$1, sent_to_pharmacy_at=$2, updated_at=$3 WHERE id=$4`
	now := time.Now().UTC()
	_, err := r.pool.Exec(ctx, q, token.Token, token.CreatedAt, now, token.PrescriptionID)
	return err
}

// GetByToken은 토큰으로 FulfillmentToken 메타데이터를 조회합니다.
func (r *TokenRepository) GetByToken(ctx context.Context, token string) (*service.FulfillmentToken, error) {
	const q = `SELECT id, COALESCE(pharmacy_id::text,''), COALESCE(fulfillment_token,''),
		COALESCE(sent_to_pharmacy_at, NOW()), expires_at,
		COALESCE(dispensary_status,'') = 'dispensed', COALESCE(dispensed_at, '0001-01-01'::timestamptz)
		FROM prescriptions WHERE fulfillment_token = $1`

	var ft service.FulfillmentToken
	var isUsed bool
	err := r.pool.QueryRow(ctx, q, token).Scan(
		&ft.PrescriptionID, &ft.PharmacyID, &ft.Token,
		&ft.CreatedAt, &ft.ExpiresAt, &isUsed, &ft.UsedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, apperrors.New(apperrors.ErrNotFound, "토큰을 찾을 수 없습니다")
		}
		return nil, err
	}
	ft.IsUsed = isUsed
	return &ft, nil
}

// MarkUsed는 토큰을 사용 완료 처리합니다.
func (r *TokenRepository) MarkUsed(ctx context.Context, token string) error {
	now := time.Now().UTC()
	const q = `UPDATE prescriptions SET dispensary_status='dispensed', dispensed_at=$1, updated_at=$1
		WHERE fulfillment_token=$2`
	tag, err := r.pool.Exec(ctx, q, now, token)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return apperrors.New(apperrors.ErrNotFound, "토큰을 찾을 수 없습니다")
	}
	return nil
}

// ============================================================================
// DrugInteractionRepository
// ============================================================================

// DrugInteractionRepository는 PostgreSQL 기반 약물 상호작용 저장소입니다.
type DrugInteractionRepository struct {
	pool *pgxpool.Pool
}

// NewDrugInteractionRepository는 DrugInteractionRepository를 생성합니다.
func NewDrugInteractionRepository(pool *pgxpool.Pool) *DrugInteractionRepository {
	return &DrugInteractionRepository{pool: pool}
}

// CheckInteractions는 약물 코드 목록에서 상호작용을 검색합니다.
func (r *DrugInteractionRepository) CheckInteractions(ctx context.Context, drugCodes []string) ([]*service.DrugInteraction, error) {
	if len(drugCodes) < 2 {
		return nil, nil
	}

	const q = `SELECT drug_a_code, drug_b_code, severity::text, COALESCE(description,''), COALESCE(recommendation,'')
		FROM drug_interactions
		WHERE drug_a_code = ANY($1) AND drug_b_code = ANY($1)`

	rows, err := r.pool.Query(ctx, q, drugCodes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*service.DrugInteraction
	for rows.Next() {
		var di service.DrugInteraction
		var sevStr string
		if err := rows.Scan(&di.DrugA, &di.DrugB, &sevStr, &di.Description, &di.Recommendation); err != nil {
			return nil, err
		}
		di.Severity = stringToSeverity[sevStr]
		result = append(result, &di)
	}
	return result, nil
}

// ============================================================================
// internal helpers
// ============================================================================

func (r *PrescriptionRepository) scanOne(ctx context.Context, query string, args ...interface{}) (*service.Prescription, error) {
	row := r.pool.QueryRow(ctx, query, args...)
	p, err := scanPrescriptionFromRow(row)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return p, nil
}

type scannable interface {
	Scan(dest ...interface{}) error
}

func scanPrescriptionFromRow(row scannable) (*service.Prescription, error) {
	var p service.Prescription
	var statusStr, fulfillmentType, dispensaryStatus string
	var sentAt, dispensedAt *time.Time

	err := row.Scan(
		&p.ID, &p.PatientUserID, &p.DoctorID, &p.ConsultationID, &statusStr,
		&p.Diagnosis, &p.Notes,
		&p.PharmacyID, &fulfillmentType, &p.ShippingAddress,
		&p.FulfillmentToken, &dispensaryStatus,
		&p.PrescribedAt, &p.ExpiresAt, &sentAt, &dispensedAt, &p.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	p.Status = stringToStatus[statusStr]
	p.FulfillmentType = service.FulfillmentType(fulfillmentType)
	p.DispensaryStatus = service.DispensaryStatus(dispensaryStatus)
	if sentAt != nil {
		p.SentToPharmacyAt = *sentAt
	}
	if dispensedAt != nil {
		p.DispensedAt = *dispensedAt
	}
	return &p, nil
}

func scanPrescriptionRow(rows pgx.Rows) (*service.Prescription, error) {
	return scanPrescriptionFromRow(rows)
}

func (r *PrescriptionRepository) loadMedications(ctx context.Context, prescriptionID string) ([]*service.Medication, error) {
	const q = `SELECT id, drug_name, COALESCE(drug_code,''), COALESCE(dosage,''),
		COALESCE(frequency,''), COALESCE(duration_days,0), COALESCE(route,''),
		COALESCE(instructions,''), COALESCE(quantity,0), COALESCE(refills_remaining,0),
		COALESCE(is_generic_allowed,true)
		FROM prescription_medications WHERE prescription_id = $1`
	rows, err := r.pool.Query(ctx, q, prescriptionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var meds []*service.Medication
	for rows.Next() {
		var m service.Medication
		if err := rows.Scan(&m.ID, &m.DrugName, &m.DrugCode, &m.Dosage,
			&m.Frequency, &m.DurationDays, &m.Route,
			&m.Instructions, &m.Quantity, &m.RefillsRemaining, &m.IsGenericAllowed); err != nil {
			return nil, err
		}
		meds = append(meds, &m)
	}
	return meds, nil
}

func insertMedication(ctx context.Context, tx pgx.Tx, prescriptionID string, m *service.Medication) error {
	const q = `INSERT INTO prescription_medications
		(id, prescription_id, drug_name, drug_code, dosage, frequency, duration_days, route,
		 instructions, quantity, refills_remaining, is_generic_allowed)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`
	_, err := tx.Exec(ctx, q,
		m.ID, prescriptionID, m.DrugName, m.DrugCode, m.Dosage, m.Frequency,
		m.DurationDays, m.Route, m.Instructions, m.Quantity, m.RefillsRemaining, m.IsGenericAllowed,
	)
	return err
}

func nullIfEmpty(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

func nullTimeIfZero(t time.Time) interface{} {
	if t.IsZero() {
		return nil
	}
	return t
}
