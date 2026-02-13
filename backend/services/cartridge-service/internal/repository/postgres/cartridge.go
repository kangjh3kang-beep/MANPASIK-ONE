// Package postgres는 cartridge-service의 PostgreSQL 저장소 구현입니다.
package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manpasik/backend/services/cartridge-service/internal/service"
)

// ============================================================================
// CartridgeUsageRepository — PostgreSQL 기반
// ============================================================================

// CartridgeUsageRepository는 PostgreSQL 기반 카트리지 사용 기록 저장소입니다.
type CartridgeUsageRepository struct {
	pool *pgxpool.Pool
}

// NewCartridgeUsageRepository는 PostgreSQL CartridgeUsageRepository를 생성합니다.
func NewCartridgeUsageRepository(pool *pgxpool.Pool) *CartridgeUsageRepository {
	return &CartridgeUsageRepository{pool: pool}
}

// Create는 카트리지 사용 기록을 저장합니다.
func (r *CartridgeUsageRepository) Create(ctx context.Context, record *service.CartridgeUsageRecord) error {
	const q = `INSERT INTO cartridge_usage_log (id, user_id, session_id, cartridge_uid, category_code, type_index, tier_at_usage, access_level, used_at)
		VALUES ($1, $2, $3, $4, $5, $6, 0, 'included', $7)`
	_, err := r.pool.Exec(ctx, q,
		record.RecordID,
		record.UserID,
		record.SessionID,
		record.CartridgeUID,
		record.CategoryCode,
		record.TypeIndex,
		record.UsedAt,
	)
	return err
}

// ListByUserID는 사용자의 카트리지 사용 이력을 조회합니다 (페이지네이션).
func (r *CartridgeUsageRepository) ListByUserID(ctx context.Context, userID string, limit, offset int32) ([]*service.CartridgeUsageRecord, int32, error) {
	// 총 개수 조회
	const countQ = `SELECT COUNT(*) FROM cartridge_usage_log WHERE user_id = $1`
	var totalCount int32
	if err := r.pool.QueryRow(ctx, countQ, userID).Scan(&totalCount); err != nil {
		return nil, 0, err
	}

	const q = `SELECT ul.id, ul.user_id, ul.session_id, ul.cartridge_uid,
			ul.category_code, ul.type_index, COALESCE(ct.name_ko, ''), ul.used_at
		FROM cartridge_usage_log ul
		LEFT JOIN cartridge_types ct ON ct.category_code = ul.category_code AND ct.type_index = ul.type_index
		WHERE ul.user_id = $1
		ORDER BY ul.used_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.pool.Query(ctx, q, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var records []*service.CartridgeUsageRecord
	for rows.Next() {
		var rec service.CartridgeUsageRecord
		if err := rows.Scan(
			&rec.RecordID, &rec.UserID, &rec.SessionID, &rec.CartridgeUID,
			&rec.CategoryCode, &rec.TypeIndex, &rec.TypeNameKO, &rec.UsedAt,
		); err != nil {
			return nil, 0, err
		}
		records = append(records, &rec)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return records, totalCount, nil
}

// ============================================================================
// CartridgeStateRepository — PostgreSQL 기반
// ============================================================================

// CartridgeStateRepository는 PostgreSQL 기반 카트리지 잔여 사용 상태 저장소입니다.
// cartridge_states 테이블: cartridge_uid (PK), remaining_uses, max_uses, expiry_date
type CartridgeStateRepository struct {
	pool *pgxpool.Pool
}

// NewCartridgeStateRepository는 PostgreSQL CartridgeStateRepository를 생성합니다.
func NewCartridgeStateRepository(pool *pgxpool.Pool) *CartridgeStateRepository {
	return &CartridgeStateRepository{pool: pool}
}

// GetByUID는 카트리지 UID로 잔여 사용 정보를 조회합니다.
func (r *CartridgeStateRepository) GetByUID(ctx context.Context, uid string) (*service.CartridgeRemainingInfo, error) {
	const q = `SELECT cartridge_uid, remaining_uses, max_uses, COALESCE(expiry_date, '')
		FROM cartridge_states WHERE cartridge_uid = $1`

	var info service.CartridgeRemainingInfo
	err := r.pool.QueryRow(ctx, q, uid).Scan(
		&info.CartridgeUID, &info.RemainingUses, &info.MaxUses, &info.ExpiryDate,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &info, nil
}

// Upsert는 카트리지 상태를 생성하거나 업데이트합니다.
func (r *CartridgeStateRepository) Upsert(ctx context.Context, info *service.CartridgeRemainingInfo) error {
	const q = `INSERT INTO cartridge_states (cartridge_uid, remaining_uses, max_uses, expiry_date)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (cartridge_uid) DO UPDATE SET
			remaining_uses = EXCLUDED.remaining_uses,
			max_uses = EXCLUDED.max_uses,
			expiry_date = EXCLUDED.expiry_date`
	_, err := r.pool.Exec(ctx, q,
		info.CartridgeUID, info.RemainingUses, info.MaxUses, info.ExpiryDate,
	)
	return err
}

// DecrementUses는 카트리지 잔여 사용 횟수를 1 감소시키고 갱신된 값을 반환합니다.
func (r *CartridgeStateRepository) DecrementUses(ctx context.Context, uid string) (int32, error) {
	const q = `UPDATE cartridge_states
		SET remaining_uses = GREATEST(remaining_uses - 1, 0)
		WHERE cartridge_uid = $1
		RETURNING remaining_uses`

	var remaining int32
	err := r.pool.QueryRow(ctx, q, uid).Scan(&remaining)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return remaining, nil
}
