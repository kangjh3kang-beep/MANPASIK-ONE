// Package postgres는 translation-service의 PostgreSQL 저장소 구현입니다.
//
// DB 스키마: infrastructure/database/init/20-translation.sql
// 테이블: translation_records, translation_usage, medical_terms, supported_languages
package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manpasik/backend/services/translation-service/internal/service"
)

// ============================================================================
// TranslationRepository
// ============================================================================

// TranslationRepository는 PostgreSQL 기반 번역 이력 저장소입니다.
type TranslationRepository struct {
	pool *pgxpool.Pool
}

// NewTranslationRepository는 TranslationRepository를 생성합니다.
func NewTranslationRepository(pool *pgxpool.Pool) *TranslationRepository {
	return &TranslationRepository{pool: pool}
}

// Save는 번역 이력을 저장합니다.
func (r *TranslationRepository) Save(ctx context.Context, record *service.TranslationRecord) error {
	const q = `INSERT INTO translation_records (id, user_id, source_text, translated_text, source_language, target_language, confidence, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	var userID *string
	if record.UserID != "" {
		userID = &record.UserID
	}
	_, err := r.pool.Exec(ctx, q,
		record.ID, userID, record.SourceText, record.TranslatedText,
		record.SourceLanguage, record.TargetLanguage, record.Confidence, record.CreatedAt,
	)
	return err
}

// FindByUserID는 사용자의 번역 이력을 조회합니다.
func (r *TranslationRepository) FindByUserID(ctx context.Context, userID string, limit, offset int) ([]*service.TranslationRecord, int, error) {
	var total int
	if err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM translation_records WHERE user_id = $1", userID).Scan(&total); err != nil {
		return nil, 0, err
	}

	const q = `SELECT id, COALESCE(user_id::text,''), source_text, translated_text, source_language, target_language, confidence, created_at
		FROM translation_records WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.pool.Query(ctx, q, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var records []*service.TranslationRecord
	for rows.Next() {
		var rec service.TranslationRecord
		if err := rows.Scan(
			&rec.ID, &rec.UserID, &rec.SourceText, &rec.TranslatedText,
			&rec.SourceLanguage, &rec.TargetLanguage, &rec.Confidence, &rec.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		records = append(records, &rec)
	}
	return records, total, rows.Err()
}

// ============================================================================
// UsageRepository
// ============================================================================

// UsageRepository는 PostgreSQL 기반 번역 사용량 저장소입니다.
type UsageRepository struct {
	pool *pgxpool.Pool
}

// NewUsageRepository는 UsageRepository를 생성합니다.
func NewUsageRepository(pool *pgxpool.Pool) *UsageRepository {
	return &UsageRepository{pool: pool}
}

// IncrementUsage는 사용량을 증가시킵니다. UPSERT 사용.
func (r *UsageRepository) IncrementUsage(ctx context.Context, userID string, characters int) error {
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)

	const q = `INSERT INTO translation_usage (user_id, total_characters, monthly_characters, total_requests, monthly_requests, month_start, updated_at)
		VALUES ($1, $2, $2, 1, 1, $3, $4)
		ON CONFLICT (user_id) DO UPDATE SET
			total_characters = translation_usage.total_characters + $2,
			monthly_characters = CASE
				WHEN translation_usage.month_start = $3 THEN translation_usage.monthly_characters + $2
				ELSE $2
			END,
			total_requests = translation_usage.total_requests + 1,
			monthly_requests = CASE
				WHEN translation_usage.month_start = $3 THEN translation_usage.monthly_requests + 1
				ELSE 1
			END,
			month_start = $3,
			updated_at = $4`
	_, err := r.pool.Exec(ctx, q, userID, characters, monthStart, now)
	return err
}

// GetUsage는 사용량을 조회합니다.
func (r *UsageRepository) GetUsage(ctx context.Context, userID string) (*service.UsageStats, error) {
	const q = `SELECT total_characters, monthly_characters, monthly_limit, total_requests, monthly_requests
		FROM translation_usage WHERE user_id = $1`
	var stats service.UsageStats
	err := r.pool.QueryRow(ctx, q, userID).Scan(
		&stats.TotalCharacters, &stats.MonthlyCharacters, &stats.MonthlyLimit,
		&stats.TotalRequests, &stats.MonthlyRequests,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return &service.UsageStats{MonthlyLimit: 100000}, nil
		}
		return nil, err
	}
	return &stats, nil
}
