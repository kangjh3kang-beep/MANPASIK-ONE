package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manpasik/backend/services/assistant-service/internal/service"
)

type AssistantRepository struct {
	pool *pgxpool.Pool
}

func NewAssistantRepository(pool *pgxpool.Pool) *AssistantRepository {
	return &AssistantRepository{pool: pool}
}

func (r *AssistantRepository) CreateSession(ctx context.Context, s *service.AssistantSession) error {
	const q = `INSERT INTO assistant_sessions (id, user_id, title, status, turn_count, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.pool.Exec(ctx, q, s.ID, s.UserID, s.Title, s.Status, s.TurnCount, s.CreatedAt, s.UpdatedAt)
	return err
}

func (r *AssistantRepository) GetSession(ctx context.Context, id string) (*service.AssistantSession, error) {
	const q = `SELECT id, user_id, title, status, turn_count, created_at, updated_at
		FROM assistant_sessions WHERE id = $1`
	var s service.AssistantSession
	err := r.pool.QueryRow(ctx, q, id).Scan(&s.ID, &s.UserID, &s.Title, &s.Status, &s.TurnCount, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *AssistantRepository) ListSessions(ctx context.Context, userID string, limit, offset int32) ([]*service.AssistantSession, int32, error) {
	var total int32
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM assistant_sessions WHERE user_id = $1`, userID).Scan(&total)

	const q = `SELECT id, user_id, title, status, turn_count, created_at, updated_at
		FROM assistant_sessions WHERE user_id = $1 ORDER BY updated_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.pool.Query(ctx, q, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var list []*service.AssistantSession
	for rows.Next() {
		var s service.AssistantSession
		if err := rows.Scan(&s.ID, &s.UserID, &s.Title, &s.Status, &s.TurnCount, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, 0, err
		}
		list = append(list, &s)
	}
	return list, total, rows.Err()
}

func (r *AssistantRepository) DeleteSession(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM assistant_turns WHERE session_id = $1`, id)
	if err != nil {
		return err
	}
	_, err = r.pool.Exec(ctx, `DELETE FROM assistant_sessions WHERE id = $1`, id)
	return err
}

func (r *AssistantRepository) AddTurn(ctx context.Context, t *service.AssistantTurn) error {
	const q = `INSERT INTO assistant_turns (id, session_id, role, content, intent, action_type, action_result, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.pool.Exec(ctx, q, t.ID, t.SessionID, t.Role, t.Content, t.Intent, t.ActionType, t.ActionResult, t.CreatedAt)
	return err
}

func (r *AssistantRepository) ListTurns(ctx context.Context, sessionID string, limit, offset int32) ([]*service.AssistantTurn, int32, error) {
	var total int32
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM assistant_turns WHERE session_id = $1`, sessionID).Scan(&total)

	const q = `SELECT id, session_id, role, content, COALESCE(intent,''), COALESCE(action_type,''), COALESCE(action_result,''), created_at
		FROM assistant_turns WHERE session_id = $1 ORDER BY created_at ASC LIMIT $2 OFFSET $3`
	rows, err := r.pool.Query(ctx, q, sessionID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var list []*service.AssistantTurn
	for rows.Next() {
		var t service.AssistantTurn
		if err := rows.Scan(&t.ID, &t.SessionID, &t.Role, &t.Content, &t.Intent, &t.ActionType, &t.ActionResult, &t.CreatedAt); err != nil {
			return nil, 0, err
		}
		list = append(list, &t)
	}
	return list, total, rows.Err()
}

func (r *AssistantRepository) UpdateSession(ctx context.Context, s *service.AssistantSession) error {
	const q = `UPDATE assistant_sessions SET title=$1, status=$2, turn_count=$3, updated_at=$4 WHERE id = $5`
	_, err := r.pool.Exec(ctx, q, s.Title, s.Status, s.TurnCount, s.UpdatedAt, s.ID)
	return err
}
