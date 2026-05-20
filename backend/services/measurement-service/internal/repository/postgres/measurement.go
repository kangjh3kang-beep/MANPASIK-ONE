package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manpasik/backend/services/measurement-service/internal/service"
	"github.com/manpasik/backend/shared/assay"
	"github.com/manpasik/backend/shared/tenancy"
)

// MeasurementRepository는 PostgreSQL/TimescaleDB 기반 MeasurementRepository 구현입니다.
type MeasurementRepository struct {
	pool *pgxpool.Pool
}

// NewMeasurementRepository는 MeasurementRepository를 생성합니다.
func NewMeasurementRepository(pool *pgxpool.Pool) *MeasurementRepository {
	return &MeasurementRepository{pool: pool}
}

// Store는 측정 데이터를 TimescaleDB에 저장합니다.
//
// Phase AG-1: ctx 의 tenant_id 를 자동으로 함께 저장. tenant 미설정 시
// NULL 로 저장 (legacy/personal 모드 호환). 격리는 SELECT 시점 SQLFilter 로
// 강제.
func (r *MeasurementRepository) Store(ctx context.Context, data *service.MeasurementData) error {
	const q = `INSERT INTO measurement_data
		(time, session_id, device_id, user_id, cartridge_type,
		 raw_channels, s_det, s_ref, alpha, s_corrected,
		 primary_value, unit, confidence, fingerprint_dim,
		 evidence_status, diagnostic_ready, evidence_gaps,
		 temp_c, humidity_pct, battery_pct, tenant_id)
		VALUES ($1, $2, $3, $4, $5,
		        $6, $7, $8, $9, $10,
		        $11, $12, $13, $14,
		        $15, $16, $17,
		        $18, $19, $20, $21)`
	var tenantID interface{} // NULL 대응
	if tid, ok := tenancy.TenantFromContext(ctx); ok && !tid.IsZero() {
		tenantID = string(tid)
	}
	evidenceStatus := data.EvidenceStatus
	if evidenceStatus == "" {
		evidenceStatus = assay.EvidenceStatus("unknown")
	}
	_, err := r.pool.Exec(ctx, q,
		data.Time, data.SessionID, data.DeviceID, data.UserID, data.CartridgeType,
		data.RawChannels, data.SDet, data.SRef, data.Alpha, data.SCorrected,
		data.PrimaryValue, data.Unit, data.Confidence, len(data.FingerprintVector),
		string(evidenceStatus), data.DiagnosticReady, data.EvidenceGaps,
		data.TempC, data.HumidityPct, data.BatteryPct, tenantID,
	)
	return err
}

// GetHistory는 사용자의 측정 기록을 시간 역순으로 조회합니다.
// TimescaleDB 환경에서는 time 컬럼의 하이퍼테이블 인덱스를 활용합니다.
//
// Phase AG-1: ctx 의 tenant_id 가 있으면 해당 조직 데이터만 반환. tenant
// 미설정 시 tenant_id IS NULL 인 (legacy/personal) 데이터만 반환 — 보안 우선.
// 교차 조직 접근은 빈 결과로 자연스럽게 차단.
func (r *MeasurementRepository) GetHistory(
	ctx context.Context,
	userID string,
	start, end time.Time,
	limit, offset int,
) ([]*service.MeasurementSummary, int, error) {
	tenantClause, tenantArgs := buildTenantClause(ctx, "md")

	// 전체 건수 조회
	countQ := `SELECT COUNT(*) FROM measurement_data md
		JOIN measurement_sessions ms ON md.session_id = ms.id
		WHERE md.user_id = $1 AND ms.status = 'completed'` + tenantClause
	args := append([]interface{}{userID}, tenantArgs...)
	argIdx := len(args) + 1

	if !start.IsZero() {
		countQ += ` AND md.time >= $` + itoa(argIdx)
		args = append(args, start)
		argIdx++
	}
	if !end.IsZero() {
		countQ += ` AND md.time <= $` + itoa(argIdx)
		args = append(args, end)
		argIdx++
	}

	var total int
	if err := r.pool.QueryRow(ctx, countQ, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	if total == 0 || offset >= total {
		return nil, total, nil
	}

	// 데이터 조회 (time DESC)
	dataQ := `SELECT md.session_id, md.cartridge_type, md.primary_value, md.unit,
			md.evidence_status, md.diagnostic_ready, md.evidence_gaps, md.time
		FROM measurement_data md
		JOIN measurement_sessions ms ON md.session_id = ms.id
		WHERE md.user_id = $1 AND ms.status = 'completed'` + tenantClause
	dataArgs := append([]interface{}{userID}, tenantArgs...)
	dataIdx := len(dataArgs) + 1

	if !start.IsZero() {
		dataQ += ` AND md.time >= $` + itoa(dataIdx)
		dataArgs = append(dataArgs, start)
		dataIdx++
	}
	if !end.IsZero() {
		dataQ += ` AND md.time <= $` + itoa(dataIdx)
		dataArgs = append(dataArgs, end)
		dataIdx++
	}

	dataQ += ` ORDER BY md.time DESC LIMIT $` + itoa(dataIdx) + ` OFFSET $` + itoa(dataIdx+1)
	dataArgs = append(dataArgs, limit, offset)

	rows, err := r.pool.Query(ctx, dataQ, dataArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var results []*service.MeasurementSummary
	for rows.Next() {
		var m service.MeasurementSummary
		var evidenceStatus string
		if err := rows.Scan(
			&m.SessionID,
			&m.CartridgeType,
			&m.PrimaryValue,
			&m.Unit,
			&evidenceStatus,
			&m.DiagnosticReady,
			&m.EvidenceGaps,
			&m.MeasuredAt,
		); err != nil {
			return nil, 0, err
		}
		if evidenceStatus == "" {
			evidenceStatus = "unknown"
		}
		m.EvidenceStatus = assay.EvidenceStatus(evidenceStatus)
		results = append(results, &m)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

// itoa는 간단한 정수→문자열 변환입니다 (쿼리 파라미터 인덱스용).
func itoa(n int) string {
	if n < 10 {
		return string(rune('0' + n))
	}
	return string(rune('0'+n/10)) + string(rune('0'+n%10))
}

// buildTenantClause 는 ctx 의 tenant_id 에 따라 SELECT 의 추가 WHERE 절을 반환.
//
// 동작:
//   - ctx 에 tenant 있음    → " AND {tableAlias}.tenant_id = $N"
//   - ctx 에 tenant 없음    → " AND {tableAlias}.tenant_id IS NULL"
//
// 두 번째 동작이 보안 우선 — tenant 미명시 시 legacy/personal 데이터만 반환,
// 다른 조직 데이터 자동 차단.
//
// 반환된 args 는 첫 번째 WHERE 인자 ($1) 다음 위치부터 추가됨 (호출자가 placeholder
// 인덱스 직접 관리).
func buildTenantClause(ctx context.Context, tableAlias string) (string, []interface{}) {
	col := "tenant_id"
	if tableAlias != "" {
		col = tableAlias + ".tenant_id"
	}
	tid, ok := tenancy.TenantFromContext(ctx)
	if !ok || tid.IsZero() {
		// legacy/personal 만 조회 — 다른 조직 데이터 자동 차단
		return " AND " + col + " IS NULL", nil
	}
	// $2 placeholder (호출자가 args[0]=userID 를 가정)
	return " AND " + col + " = $2", []interface{}{string(tid)}
}
