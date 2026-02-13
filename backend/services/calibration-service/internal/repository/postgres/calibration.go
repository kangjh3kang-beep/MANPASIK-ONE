// Package postgres는 calibration-service의 PostgreSQL 저장소 구현입니다.
package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manpasik/backend/services/calibration-service/internal/service"
)

// ============================================================================
// CalibrationRepository — PostgreSQL 기반
// ============================================================================

// CalibrationRepository는 PostgreSQL 기반 보정 기록 저장소입니다.
type CalibrationRepository struct {
	pool *pgxpool.Pool
}

// NewCalibrationRepository는 PostgreSQL CalibrationRepository를 생성합니다.
func NewCalibrationRepository(pool *pgxpool.Pool) *CalibrationRepository {
	return &CalibrationRepository{pool: pool}
}

// Save는 보정 기록을 저장합니다.
func (r *CalibrationRepository) Save(ctx context.Context, record *service.CalibrationRecord) error {
	const q = `INSERT INTO calibration_records
		(id, device_id, cartridge_category, cartridge_type_index, calibration_type,
		 alpha, channel_offsets, channel_gains, temp_coefficient, humidity_coefficient,
		 accuracy_score, reference_standard, calibrated_by, calibrated_at, expires_at, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`
	_, err := r.pool.Exec(ctx, q,
		record.ID,
		record.DeviceID,
		record.CartridgeCategory,
		record.CartridgeTypeIndex,
		calibrationTypeToString(record.CalibrationType),
		record.Alpha,
		record.ChannelOffsets,
		record.ChannelGains,
		record.TempCoefficient,
		record.HumidityCoefficient,
		record.AccuracyScore,
		record.ReferenceStandard,
		record.CalibratedBy,
		record.CalibratedAt,
		record.ExpiresAt,
		calibrationStatusToString(record.Status),
	)
	return err
}

// GetLatest는 디바이스 + 카트리지 타입별 최신 보정 기록을 반환합니다.
func (r *CalibrationRepository) GetLatest(ctx context.Context, deviceID string, category, typeIndex int32) (*service.CalibrationRecord, error) {
	const q = `SELECT id, device_id, cartridge_category, cartridge_type_index, calibration_type,
		alpha, channel_offsets, channel_gains, temp_coefficient, humidity_coefficient,
		accuracy_score, COALESCE(reference_standard, ''), COALESCE(calibrated_by, ''),
		calibrated_at, expires_at, status
		FROM calibration_records
		WHERE device_id = $1 AND cartridge_category = $2 AND cartridge_type_index = $3
		ORDER BY calibrated_at DESC LIMIT 1`

	return r.scanRecord(ctx, q, deviceID, category, typeIndex)
}

// GetLatestByType은 디바이스 + 카트리지 타입 + 보정 유형별 최신 기록을 반환합니다.
func (r *CalibrationRepository) GetLatestByType(ctx context.Context, deviceID string, category, typeIndex int32, calType service.CalibrationType) (*service.CalibrationRecord, error) {
	const q = `SELECT id, device_id, cartridge_category, cartridge_type_index, calibration_type,
		alpha, channel_offsets, channel_gains, temp_coefficient, humidity_coefficient,
		accuracy_score, COALESCE(reference_standard, ''), COALESCE(calibrated_by, ''),
		calibrated_at, expires_at, status
		FROM calibration_records
		WHERE device_id = $1 AND cartridge_category = $2 AND cartridge_type_index = $3 AND calibration_type = $4
		ORDER BY calibrated_at DESC LIMIT 1`

	return r.scanRecord(ctx, q, deviceID, category, typeIndex, calibrationTypeToString(calType))
}

// ListByDevice는 디바이스의 보정 이력을 반환합니다 (최신순, 페이지네이션).
func (r *CalibrationRepository) ListByDevice(ctx context.Context, deviceID string, limit, offset int32) ([]*service.CalibrationRecord, int32, error) {
	// 총 개수 조회
	const countQ = `SELECT COUNT(*) FROM calibration_records WHERE device_id = $1`
	var total int32
	if err := r.pool.QueryRow(ctx, countQ, deviceID).Scan(&total); err != nil {
		return nil, 0, err
	}

	const q = `SELECT id, device_id, cartridge_category, cartridge_type_index, calibration_type,
		alpha, channel_offsets, channel_gains, temp_coefficient, humidity_coefficient,
		accuracy_score, COALESCE(reference_standard, ''), COALESCE(calibrated_by, ''),
		calibrated_at, expires_at, status
		FROM calibration_records
		WHERE device_id = $1
		ORDER BY calibrated_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.pool.Query(ctx, q, deviceID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var records []*service.CalibrationRecord
	for rows.Next() {
		rec, err := r.scanRecordFromRow(rows)
		if err != nil {
			return nil, 0, err
		}
		records = append(records, rec)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

// scanRecord는 단일 보정 기록을 스캔합니다.
func (r *CalibrationRepository) scanRecord(ctx context.Context, query string, args ...any) (*service.CalibrationRecord, error) {
	var rec service.CalibrationRecord
	var calType, status string
	err := r.pool.QueryRow(ctx, query, args...).Scan(
		&rec.ID, &rec.DeviceID, &rec.CartridgeCategory, &rec.CartridgeTypeIndex, &calType,
		&rec.Alpha, &rec.ChannelOffsets, &rec.ChannelGains,
		&rec.TempCoefficient, &rec.HumidityCoefficient,
		&rec.AccuracyScore, &rec.ReferenceStandard, &rec.CalibratedBy,
		&rec.CalibratedAt, &rec.ExpiresAt, &status,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	rec.CalibrationType = calibrationTypeFromString(calType)
	rec.Status = calibrationStatusFromString(status)
	return &rec, nil
}

// scanRecordFromRow는 rows.Next() 루프 내에서 보정 기록을 스캔합니다.
func (r *CalibrationRepository) scanRecordFromRow(rows pgx.Rows) (*service.CalibrationRecord, error) {
	var rec service.CalibrationRecord
	var calType, status string
	if err := rows.Scan(
		&rec.ID, &rec.DeviceID, &rec.CartridgeCategory, &rec.CartridgeTypeIndex, &calType,
		&rec.Alpha, &rec.ChannelOffsets, &rec.ChannelGains,
		&rec.TempCoefficient, &rec.HumidityCoefficient,
		&rec.AccuracyScore, &rec.ReferenceStandard, &rec.CalibratedBy,
		&rec.CalibratedAt, &rec.ExpiresAt, &status,
	); err != nil {
		return nil, err
	}
	rec.CalibrationType = calibrationTypeFromString(calType)
	rec.Status = calibrationStatusFromString(status)
	return &rec, nil
}

// ============================================================================
// CalibrationModelRepository — PostgreSQL 기반
// ============================================================================

// CalibrationModelRepository는 PostgreSQL 기반 보정 모델 저장소입니다.
type CalibrationModelRepository struct {
	pool *pgxpool.Pool
}

// NewCalibrationModelRepository는 PostgreSQL CalibrationModelRepository를 생성합니다.
func NewCalibrationModelRepository(pool *pgxpool.Pool) *CalibrationModelRepository {
	return &CalibrationModelRepository{pool: pool}
}

// GetAll은 전체 보정 모델 목록을 반환합니다.
func (r *CalibrationModelRepository) GetAll(ctx context.Context) ([]*service.CalibrationModel, error) {
	const q = `SELECT id, cartridge_category, cartridge_type_index, name, version,
		default_alpha, validity_days, COALESCE(description, ''), created_at
		FROM calibration_models WHERE is_active = true
		ORDER BY cartridge_category, cartridge_type_index`

	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var models []*service.CalibrationModel
	for rows.Next() {
		var m service.CalibrationModel
		if err := rows.Scan(
			&m.ID, &m.CartridgeCategory, &m.CartridgeTypeIndex,
			&m.Name, &m.Version, &m.DefaultAlpha, &m.ValidityDays,
			&m.Description, &m.CreatedAt,
		); err != nil {
			return nil, err
		}
		models = append(models, &m)
	}
	return models, rows.Err()
}

// GetByCartridgeType은 카트리지 타입별 보정 모델을 반환합니다.
func (r *CalibrationModelRepository) GetByCartridgeType(ctx context.Context, category, typeIndex int32) (*service.CalibrationModel, error) {
	const q = `SELECT id, cartridge_category, cartridge_type_index, name, version,
		default_alpha, validity_days, COALESCE(description, ''), created_at
		FROM calibration_models
		WHERE cartridge_category = $1 AND cartridge_type_index = $2 AND is_active = true`

	var m service.CalibrationModel
	err := r.pool.QueryRow(ctx, q, category, typeIndex).Scan(
		&m.ID, &m.CartridgeCategory, &m.CartridgeTypeIndex,
		&m.Name, &m.Version, &m.DefaultAlpha, &m.ValidityDays,
		&m.Description, &m.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

// ============================================================================
// ENUM 변환 헬퍼
// ============================================================================

func calibrationTypeToString(ct service.CalibrationType) string {
	switch ct {
	case service.CalibrationTypeFactory:
		return "FACTORY"
	case service.CalibrationTypeField:
		return "FIELD"
	case service.CalibrationTypeAuto:
		return "AUTO"
	default:
		return "FACTORY"
	}
}

func calibrationTypeFromString(s string) service.CalibrationType {
	switch s {
	case "FACTORY":
		return service.CalibrationTypeFactory
	case "FIELD":
		return service.CalibrationTypeField
	case "AUTO":
		return service.CalibrationTypeAuto
	default:
		return service.CalibrationTypeUnknown
	}
}

func calibrationStatusToString(cs service.CalibrationStatus) string {
	switch cs {
	case service.CalibrationStatusValid:
		return "VALID"
	case service.CalibrationStatusExpiring:
		return "EXPIRING"
	case service.CalibrationStatusExpired:
		return "EXPIRED"
	case service.CalibrationStatusNeeded:
		return "NEEDED"
	default:
		return "VALID"
	}
}

func calibrationStatusFromString(s string) service.CalibrationStatus {
	switch s {
	case "VALID":
		return service.CalibrationStatusValid
	case "EXPIRING":
		return service.CalibrationStatusExpiring
	case "EXPIRED":
		return service.CalibrationStatusExpired
	case "NEEDED":
		return service.CalibrationStatusNeeded
	default:
		return service.CalibrationStatusUnknown
	}
}
