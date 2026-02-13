// Package postgres는 measurement-service의 PostgreSQL 저장소 구현입니다.
package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manpasik/backend/services/measurement-service/internal/service"
)

// SessionRepository는 PostgreSQL 기반 SessionRepository 구현입니다.
type SessionRepository struct {
	pool *pgxpool.Pool
}

// NewSessionRepository는 SessionRepository를 생성합니다.
func NewSessionRepository(pool *pgxpool.Pool) *SessionRepository {
	return &SessionRepository{pool: pool}
}

// CreateSession은 새 측정 세션을 생성합니다.
func (r *SessionRepository) CreateSession(ctx context.Context, session *service.MeasurementSession) error {
	const q = `INSERT INTO measurement_sessions (id, device_id, cartridge_id, user_id, status, total_measurements, started_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.pool.Exec(ctx, q,
		session.ID, session.DeviceID, session.CartridgeID, session.UserID,
		session.Status, session.TotalMeasurements, session.StartedAt,
	)
	return err
}

// GetSession은 세션 ID로 측정 세션을 조회합니다.
func (r *SessionRepository) GetSession(ctx context.Context, sessionID string) (*service.MeasurementSession, error) {
	const q = `SELECT id, device_id, cartridge_id, user_id, status, total_measurements, started_at, ended_at
		FROM measurement_sessions WHERE id = $1`
	var s service.MeasurementSession
	err := r.pool.QueryRow(ctx, q, sessionID).Scan(
		&s.ID, &s.DeviceID, &s.CartridgeID, &s.UserID,
		&s.Status, &s.TotalMeasurements, &s.StartedAt, &s.EndedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

// EndSession은 측정 세션을 종료합니다.
func (r *SessionRepository) EndSession(ctx context.Context, sessionID string, totalMeasurements int, endedAt time.Time) error {
	const q = `UPDATE measurement_sessions SET status = 'completed', total_measurements = $1, ended_at = $2 WHERE id = $3`
	_, err := r.pool.Exec(ctx, q, totalMeasurements, endedAt, sessionID)
	return err
}
