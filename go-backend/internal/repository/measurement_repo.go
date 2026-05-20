package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kangjh3kang/manpasik/go-backend/internal/domain"
)

// MeasurementRepository는 DB(PostgreSQL/TimescaleDB)와 실제로 통신하는 인터페이스 계약입니다.
type MeasurementRepository interface {
	InsertMeasurement(ctx context.Context, m *domain.Measurement) error
	GetMeasurementsByUser(ctx context.Context, userID uuid.UUID, limit int) ([]*domain.Measurement, error)
}

type pgMeasurementRepository struct {
	db *pgxpool.Pool
}

// NewMeasurementRepository 생성자
func NewMeasurementRepository(db *pgxpool.Pool) MeasurementRepository {
	return &pgMeasurementRepository{db: db}
}

// InsertMeasurement: 앱에서 올라온 차동 분석 신호(16핀) 및 트윈 결과값을 로우 레벨 최적화 쿼리로 DB에 즉시 적재합니다.
func (r *pgMeasurementRepository) InsertMeasurement(ctx context.Context, m *domain.Measurement) error {
	query := `
		INSERT INTO measurements (
			id, user_id, device_mac, measured_at,
			diff_signal, fingerprint, health_score, risk_label,
			client_local_id, synced_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, CURRENT_TIMESTAMP
		)`

	_, err := r.db.Exec(ctx, query,
		m.ID, m.UserID, m.DeviceMAC, m.MeasuredAt,
		m.DiffSignal, m.Fingerprint, m.HealthScore, m.RiskLabel,
		m.ClientLocalID,
	)
	if err != nil {
		return fmt.Errorf("측정 데이터 DB 적재 실패 (InsertMeasurement): %w", err)
	}
	return nil
}

// GetMeasurementsByUser: 특정 사용자의 헬스케어 차트 및 결과 피드백을 위해 데이터를 시계열(최신순)로 조회합니다.
func (r *pgMeasurementRepository) GetMeasurementsByUser(ctx context.Context, userID uuid.UUID, limit int) ([]*domain.Measurement, error) {
	query := `
		SELECT id, user_id, device_mac, measured_at, 
		       diff_signal, fingerprint, health_score, risk_label, 
		       client_local_id, synced_at
		FROM measurements
		WHERE user_id = $1
		ORDER BY measured_at DESC
		LIMIT $2`

	rows, err := r.db.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("측정 데이터 조회 실패 (GetMeasurementsByUser): %w", err)
	}
	defer rows.Close()

	var results []*domain.Measurement
	for rows.Next() {
		var m domain.Measurement
		err := rows.Scan(
			&m.ID, &m.UserID, &m.DeviceMAC, &m.MeasuredAt,
			&m.DiffSignal, &m.Fingerprint, &m.HealthScore, &m.RiskLabel,
			&m.ClientLocalID, &m.SyncedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("로우 데이터 맵핑 실패: %w", err)
		}
		results = append(results, &m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("결과셋 순회 중 에러: %w", err)
	}

	return results, nil
}
