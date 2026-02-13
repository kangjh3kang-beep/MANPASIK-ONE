package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manpasik/backend/services/user-service/internal/service"
)

// SubscriptionRepository는 PostgreSQL 기반 SubscriptionRepository 구현입니다.
type SubscriptionRepository struct {
	pool *pgxpool.Pool
}

// NewSubscriptionRepository는 SubscriptionRepository를 생성합니다.
func NewSubscriptionRepository(pool *pgxpool.Pool) *SubscriptionRepository {
	return &SubscriptionRepository{pool: pool}
}

// GetByUserID는 사용자 ID로 구독 정보를 조회합니다.
func (r *SubscriptionRepository) GetByUserID(ctx context.Context, userID string) (*service.Subscription, error) {
	const q = `SELECT id, user_id, tier, started_at, COALESCE(expires_at, '0001-01-01'::timestamptz),
		max_devices, max_family_members, ai_coaching_enabled, telemedicine_enabled
		FROM subscriptions WHERE user_id = $1
		ORDER BY created_at DESC LIMIT 1`
	var s service.Subscription
	var tier string
	err := r.pool.QueryRow(ctx, q, userID).Scan(
		&s.ID, &s.UserID, &tier, &s.StartedAt, &s.ExpiresAt,
		&s.MaxDevices, &s.MaxFamilyMembers, &s.AICoachingEnabled, &s.TelemedicineEnabled,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	s.Tier = service.SubscriptionTier(tier)
	return &s, nil
}

// Create는 새 구독을 생성합니다.
func (r *SubscriptionRepository) Create(ctx context.Context, sub *service.Subscription) error {
	const q = `INSERT INTO subscriptions
		(id, user_id, tier, started_at, expires_at, max_devices, max_family_members, ai_coaching_enabled, telemedicine_enabled)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.pool.Exec(ctx, q,
		sub.ID, sub.UserID, string(sub.Tier), sub.StartedAt, sub.ExpiresAt,
		sub.MaxDevices, sub.MaxFamilyMembers, sub.AICoachingEnabled, sub.TelemedicineEnabled,
	)
	return err
}

// Update는 구독 정보를 업데이트합니다.
func (r *SubscriptionRepository) Update(ctx context.Context, sub *service.Subscription) error {
	const q = `UPDATE subscriptions SET
		tier = $1, expires_at = $2, max_devices = $3, max_family_members = $4,
		ai_coaching_enabled = $5, telemedicine_enabled = $6, updated_at = NOW()
		WHERE id = $7`
	_, err := r.pool.Exec(ctx, q,
		string(sub.Tier), sub.ExpiresAt, sub.MaxDevices, sub.MaxFamilyMembers,
		sub.AICoachingEnabled, sub.TelemedicineEnabled, sub.ID,
	)
	return err
}
