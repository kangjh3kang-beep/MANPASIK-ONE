package postgres

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manpasik/backend/services/analytics-service/internal/service"
)

type AnalyticsRepository struct {
	pool *pgxpool.Pool
}

func NewAnalyticsRepository(pool *pgxpool.Pool) *AnalyticsRepository {
	return &AnalyticsRepository{pool: pool}
}

func (r *AnalyticsRepository) TrackEvent(ctx context.Context, event *service.AnalyticsEvent) error {
	propsJSON, _ := json.Marshal(event.Properties)
	const q = `INSERT INTO analytics_events (id, user_id, event_type, properties, created_at)
		VALUES ($1, $2, $3, $4, $5)`
	_, err := r.pool.Exec(ctx, q, event.ID, event.UserID, event.EventType, propsJSON, event.CreatedAt)
	return err
}

func (r *AnalyticsRepository) GetUserAnalytics(ctx context.Context, userID string) (*service.UserAnalytics, error) {
	const q = `SELECT COUNT(*), MAX(created_at) FROM analytics_events WHERE user_id = $1`
	var total int
	var lastAt *time.Time
	err := r.pool.QueryRow(ctx, q, userID).Scan(&total, &lastAt)
	if err != nil {
		return nil, err
	}
	ua := &service.UserAnalytics{UserID: userID, TotalEvents: total}
	if lastAt != nil {
		ua.LastEventAt = *lastAt
	}
	// top event types
	const q2 = `SELECT event_type, COUNT(*) as cnt FROM analytics_events WHERE user_id = $1
		GROUP BY event_type ORDER BY cnt DESC LIMIT 5`
	rows, err := r.pool.Query(ctx, q2, userID)
	if err != nil {
		return ua, nil
	}
	defer rows.Close()
	for rows.Next() {
		var et string
		var c int
		if err := rows.Scan(&et, &c); err == nil {
			ua.TopEventTypes = append(ua.TopEventTypes, et)
		}
	}
	return ua, nil
}

func (r *AnalyticsRepository) GetDailyStats(ctx context.Context, date string) (*service.DailyStats, error) {
	const q = `SELECT stat_date, total_events, unique_users, events_by_type FROM analytics_daily_stats WHERE stat_date = $1`
	var ds service.DailyStats
	var ebtJSON []byte
	err := r.pool.QueryRow(ctx, q, date).Scan(&ds.Date, &ds.TotalEvents, &ds.UniqueUsers, &ebtJSON)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if len(ebtJSON) > 0 {
		json.Unmarshal(ebtJSON, &ds.EventsByType)
	}
	return &ds, nil
}

func (r *AnalyticsRepository) ListRecentEvents(ctx context.Context, limit int) ([]*service.AnalyticsEvent, error) {
	const q = `SELECT id, user_id, event_type, properties, created_at FROM analytics_events
		ORDER BY created_at DESC LIMIT $1`
	rows, err := r.pool.Query(ctx, q, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*service.AnalyticsEvent
	for rows.Next() {
		var e service.AnalyticsEvent
		var propsJSON []byte
		if err := rows.Scan(&e.ID, &e.UserID, &e.EventType, &propsJSON, &e.CreatedAt); err != nil {
			return nil, err
		}
		if len(propsJSON) > 0 {
			json.Unmarshal(propsJSON, &e.Properties)
		}
		list = append(list, &e)
	}
	return list, rows.Err()
}

func (r *AnalyticsRepository) ListEventsByUser(ctx context.Context, userID string) ([]*service.AnalyticsEvent, error) {
	const q = `SELECT id, user_id, event_type, properties, created_at FROM analytics_events
		WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := r.pool.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*service.AnalyticsEvent
	for rows.Next() {
		var e service.AnalyticsEvent
		var propsJSON []byte
		if err := rows.Scan(&e.ID, &e.UserID, &e.EventType, &propsJSON, &e.CreatedAt); err != nil {
			return nil, err
		}
		if len(propsJSON) > 0 {
			json.Unmarshal(propsJSON, &e.Properties)
		}
		list = append(list, &e)
	}
	return list, rows.Err()
}
