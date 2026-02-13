package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manpasik/backend/services/measurement-service/internal/service"
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
func (r *MeasurementRepository) Store(ctx context.Context, data *service.MeasurementData) error {
	const q = `INSERT INTO measurement_data
		(time, session_id, device_id, user_id, cartridge_type,
		 raw_channels, s_det, s_ref, alpha, s_corrected,
		 primary_value, unit, confidence, fingerprint_dim,
		 temp_c, humidity_pct, battery_pct)
		VALUES ($1, $2, $3, $4, $5,
		        $6, $7, $8, $9, $10,
		        $11, $12, $13, $14,
		        $15, $16, $17)`
	_, err := r.pool.Exec(ctx, q,
		data.Time, data.SessionID, data.DeviceID, data.UserID, data.CartridgeType,
		data.RawChannels, data.SDet, data.SRef, data.Alpha, data.SCorrected,
		data.PrimaryValue, data.Unit, data.Confidence, len(data.FingerprintVector),
		data.TempC, data.HumidityPct, data.BatteryPct,
	)
	return err
}

// GetHistory는 사용자의 측정 기록을 시간 역순으로 조회합니다.
// TimescaleDB 환경에서는 time 컬럼의 하이퍼테이블 인덱스를 활용합니다.
func (r *MeasurementRepository) GetHistory(
	ctx context.Context,
	userID string,
	start, end time.Time,
	limit, offset int,
) ([]*service.MeasurementSummary, int, error) {
	// 전체 건수 조회
	countQ := `SELECT COUNT(*) FROM measurement_data md
		JOIN measurement_sessions ms ON md.session_id = ms.id
		WHERE md.user_id = $1 AND ms.status = 'completed'`
	args := []interface{}{userID}
	argIdx := 2

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
	dataQ := `SELECT md.session_id, md.cartridge_type, md.primary_value, md.unit, md.time
		FROM measurement_data md
		JOIN measurement_sessions ms ON md.session_id = ms.id
		WHERE md.user_id = $1 AND ms.status = 'completed'`
	dataArgs := []interface{}{userID}
	dataIdx := 2

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
		if err := rows.Scan(&m.SessionID, &m.CartridgeType, &m.PrimaryValue, &m.Unit, &m.MeasuredAt); err != nil {
			return nil, 0, err
		}
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
