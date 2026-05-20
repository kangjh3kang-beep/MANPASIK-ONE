package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manpasik/backend/services/voice-profile-service/internal/service"
)

type VoiceProfileRepository struct {
	pool *pgxpool.Pool
}

func NewVoiceProfileRepository(pool *pgxpool.Pool) *VoiceProfileRepository {
	return &VoiceProfileRepository{pool: pool}
}

func (r *VoiceProfileRepository) Create(ctx context.Context, p *service.VoiceProfile) error {
	const q = `INSERT INTO voice_profiles (profile_id, user_id, profile_name, language, status, model_url, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.pool.Exec(ctx, q, p.ProfileID, p.UserID, p.ProfileName, p.Language, p.Status, p.ModelURL, p.CreatedAt)
	return err
}

func (r *VoiceProfileRepository) GetByID(ctx context.Context, profileID string) (*service.VoiceProfile, error) {
	const q = `SELECT profile_id, user_id, profile_name, language, status, COALESCE(model_url,''), created_at
		FROM voice_profiles WHERE profile_id = $1`
	var p service.VoiceProfile
	err := r.pool.QueryRow(ctx, q, profileID).Scan(&p.ProfileID, &p.UserID, &p.ProfileName, &p.Language, &p.Status, &p.ModelURL, &p.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows { return nil, nil }
		return nil, err
	}
	return &p, nil
}

func (r *VoiceProfileRepository) ListByUser(ctx context.Context, userID string) ([]*service.VoiceProfile, error) {
	const q = `SELECT profile_id, user_id, profile_name, language, status, COALESCE(model_url,''), created_at
		FROM voice_profiles WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := r.pool.Query(ctx, q, userID)
	if err != nil { return nil, err }
	defer rows.Close()
	var list []*service.VoiceProfile
	for rows.Next() {
		var p service.VoiceProfile
		if err := rows.Scan(&p.ProfileID, &p.UserID, &p.ProfileName, &p.Language, &p.Status, &p.ModelURL, &p.CreatedAt); err != nil { return nil, err }
		list = append(list, &p)
	}
	return list, rows.Err()
}

func (r *VoiceProfileRepository) Update(ctx context.Context, p *service.VoiceProfile) error {
	const q = `UPDATE voice_profiles SET profile_name=$1, language=$2, status=$3, model_url=$4 WHERE profile_id=$5`
	_, err := r.pool.Exec(ctx, q, p.ProfileName, p.Language, p.Status, p.ModelURL, p.ProfileID)
	return err
}

func (r *VoiceProfileRepository) Delete(ctx context.Context, profileID string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM voice_profiles WHERE profile_id = $1`, profileID)
	return err
}
