package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manpasik/backend/services/digital-twin-service/internal/service"
)

type TwinRepository struct {
	pool *pgxpool.Pool
}

func NewTwinRepository(pool *pgxpool.Pool) *TwinRepository {
	return &TwinRepository{pool: pool}
}

func (r *TwinRepository) Store(ctx context.Context, state *service.TwinState) error {
	const q = `INSERT INTO twin_states (id, session_id, user_id, device_id, ewma_value, cusum_pos, cusum_neg,
		health_state, drift_score, remaining_measurements, measurement_count, last_synced_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		ON CONFLICT (id) DO UPDATE SET
			ewma_value = EXCLUDED.ewma_value, cusum_pos = EXCLUDED.cusum_pos, cusum_neg = EXCLUDED.cusum_neg,
			health_state = EXCLUDED.health_state, drift_score = EXCLUDED.drift_score,
			remaining_measurements = EXCLUDED.remaining_measurements, measurement_count = EXCLUDED.measurement_count,
			last_synced_at = EXCLUDED.last_synced_at`
	_, err := r.pool.Exec(ctx, q,
		state.ID, state.SessionID, state.UserID, state.DeviceID,
		state.EWMAValue, state.CUSUMPos, state.CUSUMNeg,
		state.HealthState, state.DriftScore,
		state.RemainingMeasurements, state.MeasurementCount,
		state.LastSyncedAt, state.CreatedAt,
	)
	return err
}

func (r *TwinRepository) Get(ctx context.Context, sessionID string) (*service.TwinState, error) {
	const q = `SELECT id, session_id, user_id, device_id, ewma_value, cusum_pos, cusum_neg,
		health_state, drift_score, remaining_measurements, measurement_count, last_synced_at, created_at
		FROM twin_states WHERE session_id = $1 ORDER BY last_synced_at DESC LIMIT 1`
	var s service.TwinState
	err := r.pool.QueryRow(ctx, q, sessionID).Scan(
		&s.ID, &s.SessionID, &s.UserID, &s.DeviceID,
		&s.EWMAValue, &s.CUSUMPos, &s.CUSUMNeg,
		&s.HealthState, &s.DriftScore,
		&s.RemainingMeasurements, &s.MeasurementCount,
		&s.LastSyncedAt, &s.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *TwinRepository) ListByDevice(ctx context.Context, deviceID string, limit int) ([]*service.TwinState, error) {
	const q = `SELECT id, session_id, user_id, device_id, ewma_value, cusum_pos, cusum_neg,
		health_state, drift_score, remaining_measurements, measurement_count, last_synced_at, created_at
		FROM twin_states WHERE device_id = $1 ORDER BY last_synced_at DESC LIMIT $2`
	rows, err := r.pool.Query(ctx, q, deviceID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*service.TwinState
	for rows.Next() {
		var s service.TwinState
		if err := rows.Scan(
			&s.ID, &s.SessionID, &s.UserID, &s.DeviceID,
			&s.EWMAValue, &s.CUSUMPos, &s.CUSUMNeg,
			&s.HealthState, &s.DriftScore,
			&s.RemainingMeasurements, &s.MeasurementCount,
			&s.LastSyncedAt, &s.CreatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, &s)
	}
	return list, rows.Err()
}
