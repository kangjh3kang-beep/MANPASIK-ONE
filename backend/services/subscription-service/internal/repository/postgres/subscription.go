// Package postgres는 subscription-service의 PostgreSQL 저장소 구현입니다.
package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manpasik/backend/services/subscription-service/internal/service"
)

// ============================================================================
// SubscriptionRepository
// ============================================================================

// SubscriptionRepository는 PostgreSQL 기반 SubscriptionRepository 구현입니다.
type SubscriptionRepository struct {
	pool *pgxpool.Pool
}

// NewSubscriptionRepository는 SubscriptionRepository를 생성합니다.
func NewSubscriptionRepository(pool *pgxpool.Pool) *SubscriptionRepository {
	return &SubscriptionRepository{pool: pool}
}

// GetByUserID는 사용자 ID로 구독을 조회합니다 (활성 구독 우선).
func (r *SubscriptionRepository) GetByUserID(ctx context.Context, userID string) (*service.Subscription, error) {
	const q = `SELECT id, user_id, tier, status, started_at, expires_at, cancelled_at,
		max_devices, max_family_members, ai_coaching_enabled, telemedicine_enabled,
		monthly_price_krw, auto_renew, created_at, updated_at
		FROM subscriptions WHERE user_id = $1 ORDER BY status ASC, created_at DESC LIMIT 1`

	var s service.Subscription
	err := r.pool.QueryRow(ctx, q, userID).Scan(
		&s.ID, &s.UserID, &s.Tier, &s.Status, &s.StartedAt, &s.ExpiresAt, &s.CancelledAt,
		&s.MaxDevices, &s.MaxFamilyMembers, &s.AICoachingEnabled, &s.TelemedicineEnabled,
		&s.MonthlyPriceKRW, &s.AutoRenew, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

// GetByID는 구독 ID로 조회합니다.
func (r *SubscriptionRepository) GetByID(ctx context.Context, id string) (*service.Subscription, error) {
	const q = `SELECT id, user_id, tier, status, started_at, expires_at, cancelled_at,
		max_devices, max_family_members, ai_coaching_enabled, telemedicine_enabled,
		monthly_price_krw, auto_renew, created_at, updated_at
		FROM subscriptions WHERE id = $1`

	var s service.Subscription
	err := r.pool.QueryRow(ctx, q, id).Scan(
		&s.ID, &s.UserID, &s.Tier, &s.Status, &s.StartedAt, &s.ExpiresAt, &s.CancelledAt,
		&s.MaxDevices, &s.MaxFamilyMembers, &s.AICoachingEnabled, &s.TelemedicineEnabled,
		&s.MonthlyPriceKRW, &s.AutoRenew, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

// Create는 새 구독을 생성합니다.
func (r *SubscriptionRepository) Create(ctx context.Context, sub *service.Subscription) error {
	const q = `INSERT INTO subscriptions
		(id, user_id, tier, status, started_at, expires_at, cancelled_at,
		 max_devices, max_family_members, ai_coaching_enabled, telemedicine_enabled,
		 monthly_price_krw, auto_renew, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)`
	_, err := r.pool.Exec(ctx, q,
		sub.ID, sub.UserID, int32(sub.Tier), int32(sub.Status),
		sub.StartedAt, sub.ExpiresAt, sub.CancelledAt,
		sub.MaxDevices, sub.MaxFamilyMembers, sub.AICoachingEnabled, sub.TelemedicineEnabled,
		sub.MonthlyPriceKRW, sub.AutoRenew, sub.CreatedAt, sub.UpdatedAt,
	)
	return err
}

// Update는 구독을 업데이트합니다.
func (r *SubscriptionRepository) Update(ctx context.Context, sub *service.Subscription) error {
	const q = `UPDATE subscriptions SET
		tier=$1, status=$2, expires_at=$3, cancelled_at=$4,
		max_devices=$5, max_family_members=$6, ai_coaching_enabled=$7, telemedicine_enabled=$8,
		monthly_price_krw=$9, auto_renew=$10, updated_at=$11
		WHERE id=$12`
	_, err := r.pool.Exec(ctx, q,
		int32(sub.Tier), int32(sub.Status), sub.ExpiresAt, sub.CancelledAt,
		sub.MaxDevices, sub.MaxFamilyMembers, sub.AICoachingEnabled, sub.TelemedicineEnabled,
		sub.MonthlyPriceKRW, sub.AutoRenew, sub.UpdatedAt,
		sub.ID,
	)
	return err
}

// ============================================================================
// SubscriptionHistoryRepository
// ============================================================================

// SubscriptionHistoryRepository는 PostgreSQL 기반 SubscriptionHistoryRepository 구현입니다.
type SubscriptionHistoryRepository struct {
	pool *pgxpool.Pool
}

// NewSubscriptionHistoryRepository는 SubscriptionHistoryRepository를 생성합니다.
func NewSubscriptionHistoryRepository(pool *pgxpool.Pool) *SubscriptionHistoryRepository {
	return &SubscriptionHistoryRepository{pool: pool}
}

// Record는 구독 변경 이력을 기록합니다.
func (r *SubscriptionHistoryRepository) Record(ctx context.Context, entry *service.SubscriptionHistoryEntry) error {
	const q = `INSERT INTO subscription_history (id, user_id, old_tier, new_tier, action, reason, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7)`
	_, err := r.pool.Exec(ctx, q,
		entry.ID, entry.UserID, int32(entry.OldTier), int32(entry.NewTier),
		entry.Action, entry.Reason, entry.CreatedAt,
	)
	return err
}

// ListByUserID는 사용자의 구독 변경 이력을 조회합니다.
func (r *SubscriptionHistoryRepository) ListByUserID(ctx context.Context, userID string, limit, offset int32) ([]*service.SubscriptionHistoryEntry, error) {
	const q = `SELECT id, user_id, old_tier, new_tier, action, COALESCE(reason, ''), created_at
		FROM subscription_history WHERE user_id = $1
		ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.pool.Query(ctx, q, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []*service.SubscriptionHistoryEntry
	for rows.Next() {
		var e service.SubscriptionHistoryEntry
		if err := rows.Scan(
			&e.ID, &e.UserID, &e.OldTier, &e.NewTier,
			&e.Action, &e.Reason, &e.CreatedAt,
		); err != nil {
			return nil, err
		}
		entries = append(entries, &e)
	}
	return entries, rows.Err()
}
