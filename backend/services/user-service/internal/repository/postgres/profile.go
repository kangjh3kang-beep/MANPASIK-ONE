// Package postgres는 user-service의 PostgreSQL 저장소 구현입니다.
package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manpasik/backend/services/user-service/internal/service"
)

// ProfileRepository는 PostgreSQL 기반 ProfileRepository 구현입니다.
type ProfileRepository struct {
	pool *pgxpool.Pool
}

// NewProfileRepository는 ProfileRepository를 생성합니다.
func NewProfileRepository(pool *pgxpool.Pool) *ProfileRepository {
	return &ProfileRepository{pool: pool}
}

// GetByID는 사용자 ID로 프로필을 조회합니다.
func (r *ProfileRepository) GetByID(ctx context.Context, userID string) (*service.UserProfile, error) {
	const q = `SELECT user_id, email, COALESCE(display_name, ''), COALESCE(avatar_url, ''),
		COALESCE(language, 'ko'), COALESCE(timezone, 'Asia/Seoul'),
		COALESCE(subscription_tier, 'free'), created_at, updated_at
		FROM user_profiles WHERE user_id = $1`
	var p service.UserProfile
	err := r.pool.QueryRow(ctx, q, userID).Scan(
		&p.UserID, &p.Email, &p.DisplayName, &p.AvatarURL,
		&p.Language, &p.Timezone, &p.SubscriptionTier,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

// Update는 사용자 프로필을 업데이트합니다 (upsert).
func (r *ProfileRepository) Update(ctx context.Context, profile *service.UserProfile) error {
	const q = `INSERT INTO user_profiles (user_id, email, display_name, avatar_url, language, timezone, subscription_tier, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (user_id) DO UPDATE SET
			email = EXCLUDED.email,
			display_name = EXCLUDED.display_name,
			avatar_url = EXCLUDED.avatar_url,
			language = EXCLUDED.language,
			timezone = EXCLUDED.timezone,
			subscription_tier = EXCLUDED.subscription_tier,
			updated_at = EXCLUDED.updated_at`
	_, err := r.pool.Exec(ctx, q,
		profile.UserID, profile.Email, profile.DisplayName, profile.AvatarURL,
		profile.Language, profile.Timezone, profile.SubscriptionTier,
		profile.CreatedAt, profile.UpdatedAt,
	)
	return err
}
