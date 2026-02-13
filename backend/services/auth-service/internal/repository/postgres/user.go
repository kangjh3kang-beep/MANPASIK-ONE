// Package postgres는 auth-service의 PostgreSQL 저장소 구현입니다.
package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manpasik/backend/services/auth-service/internal/service"
)

// UserRepository는 PostgreSQL 기반 UserRepository 구현입니다.
type UserRepository struct {
	pool *pgxpool.Pool
}

// NewUserRepository는 UserRepository를 생성합니다.
func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

// GetByID는 ID로 사용자를 조회합니다.
func (r *UserRepository) GetByID(ctx context.Context, id string) (*service.User, error) {
	const q = `SELECT id, email, password_hash, COALESCE(display_name, ''), role, is_active, created_at, updated_at
		FROM users WHERE id = $1`
	var u service.User
	var displayName string
	err := r.pool.QueryRow(ctx, q, id).Scan(
		&u.ID, &u.Email, &u.HashedPassword, &displayName, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	u.DisplayName = displayName
	return &u, nil
}

// GetByEmail은 이메일로 사용자를 조회합니다.
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*service.User, error) {
	const q = `SELECT id, email, password_hash, COALESCE(display_name, ''), role, is_active, created_at, updated_at
		FROM users WHERE email = $1`
	var u service.User
	var displayName string
	err := r.pool.QueryRow(ctx, q, email).Scan(
		&u.ID, &u.Email, &u.HashedPassword, &displayName, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	u.DisplayName = displayName
	return &u, nil
}

// Create는 사용자를 생성합니다.
func (r *UserRepository) Create(ctx context.Context, user *service.User) error {
	const q = `INSERT INTO users (id, email, password_hash, display_name, role, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.pool.Exec(ctx, q,
		user.ID, user.Email, user.HashedPassword, user.DisplayName, user.Role, user.IsActive, user.CreatedAt, user.UpdatedAt,
	)
	return err
}

// UpdatePassword는 비밀번호를 업데이트합니다.
func (r *UserRepository) UpdatePassword(ctx context.Context, id string, hashedPassword string) error {
	const q = `UPDATE users SET password_hash = $1, updated_at = $2 WHERE id = $3`
	_, err := r.pool.Exec(ctx, q, hashedPassword, time.Now().UTC(), id)
	return err
}
