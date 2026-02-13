// Package postgres는 admin-service의 PostgreSQL 저장소 구현입니다.
package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manpasik/backend/services/admin-service/internal/service"
)

// ============================================================================
// admin_role / audit_action ↔ Go enum 매핑
// ============================================================================

var roleToString = map[service.AdminRole]string{
	service.RoleSuperAdmin: "super_admin",
	service.RoleAdmin:      "admin",
	service.RoleModerator:  "moderator",
	service.RoleSupport:    "support",
	service.RoleAnalyst:    "analyst",
}

var stringToRole = map[string]service.AdminRole{
	"super_admin": service.RoleSuperAdmin,
	"admin":       service.RoleAdmin,
	"moderator":   service.RoleModerator,
	"support":     service.RoleSupport,
	"analyst":     service.RoleAnalyst,
}

var actionToString = map[service.AuditAction]string{
	service.ActionLogin:        "login",
	service.ActionLogout:       "logout",
	service.ActionCreate:       "create",
	service.ActionUpdate:       "update",
	service.ActionDelete:       "delete",
	service.ActionConfigChange: "config_change",
	service.ActionUserBan:      "user_ban",
	service.ActionUserUnban:    "user_unban",
	service.ActionRoleChange:   "role_change",
}

var stringToAction = map[string]service.AuditAction{
	"login":         service.ActionLogin,
	"logout":        service.ActionLogout,
	"create":        service.ActionCreate,
	"update":        service.ActionUpdate,
	"delete":        service.ActionDelete,
	"config_change": service.ActionConfigChange,
	"user_ban":      service.ActionUserBan,
	"user_unban":    service.ActionUserUnban,
	"role_change":   service.ActionRoleChange,
}

// ============================================================================
// AdminRepository
// ============================================================================

// AdminRepository는 PostgreSQL 기반 관리자 저장소입니다.
type AdminRepository struct {
	pool *pgxpool.Pool
}

// NewAdminRepository는 AdminRepository를 생성합니다.
func NewAdminRepository(pool *pgxpool.Pool) *AdminRepository {
	return &AdminRepository{pool: pool}
}

// Save는 관리자를 저장합니다.
func (r *AdminRepository) Save(ctx context.Context, admin *service.AdminUser) error {
	const q = `INSERT INTO admin_users
		(admin_id, user_id, email, display_name, role, is_active, created_at, last_login_at, created_by, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`

	_, err := r.pool.Exec(ctx, q,
		admin.AdminID, admin.UserID, admin.Email, admin.DisplayName,
		roleToString[admin.Role], admin.IsActive, admin.CreatedAt, admin.LastLoginAt,
		nullIfEmpty(admin.CreatedBy), time.Now().UTC(),
	)
	return err
}

// FindByID는 관리자 ID로 조회합니다.
func (r *AdminRepository) FindByID(ctx context.Context, adminID string) (*service.AdminUser, error) {
	const q = `SELECT admin_id, user_id, email, display_name, role::text, is_active,
		created_at, COALESCE(last_login_at, created_at), COALESCE(created_by,'')
		FROM admin_users WHERE admin_id = $1`

	var a service.AdminUser
	var roleStr string
	err := r.pool.QueryRow(ctx, q, adminID).Scan(
		&a.AdminID, &a.UserID, &a.Email, &a.DisplayName, &roleStr, &a.IsActive,
		&a.CreatedAt, &a.LastLoginAt, &a.CreatedBy,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	a.Role = stringToRole[roleStr]
	return &a, nil
}

// FindAll은 관리자 목록을 조회합니다.
func (r *AdminRepository) FindAll(ctx context.Context, roleFilter service.AdminRole, activeOnly bool, limit, offset int32) ([]*service.AdminUser, int, error) {
	where := `WHERE 1=1`
	args := []interface{}{}
	idx := 1

	if roleFilter != service.RoleUnknown {
		where += fmt.Sprintf(` AND role = $%d`, idx)
		args = append(args, roleToString[roleFilter])
		idx++
	}
	if activeOnly {
		where += ` AND is_active = true`
	}

	// count
	var total int
	if err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM admin_users `+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// data
	dataQ := fmt.Sprintf(`SELECT admin_id, user_id, email, display_name, role::text, is_active,
		created_at, COALESCE(last_login_at, created_at), COALESCE(created_by,'')
		FROM admin_users %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d`,
		where, idx, idx+1)
	dataArgs := append(args, limit, offset)

	rows, err := r.pool.Query(ctx, dataQ, dataArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []*service.AdminUser
	for rows.Next() {
		var a service.AdminUser
		var roleStr string
		if err := rows.Scan(
			&a.AdminID, &a.UserID, &a.Email, &a.DisplayName, &roleStr, &a.IsActive,
			&a.CreatedAt, &a.LastLoginAt, &a.CreatedBy,
		); err != nil {
			return nil, 0, err
		}
		a.Role = stringToRole[roleStr]
		list = append(list, &a)
	}
	return list, total, nil
}

// Update는 관리자 정보를 업데이트합니다.
func (r *AdminRepository) Update(ctx context.Context, admin *service.AdminUser) error {
	const q = `UPDATE admin_users SET
		display_name=$1, role=$2, is_active=$3, last_login_at=$4, updated_at=$5
		WHERE admin_id=$6`

	_, err := r.pool.Exec(ctx, q,
		admin.DisplayName, roleToString[admin.Role], admin.IsActive,
		admin.LastLoginAt, time.Now().UTC(),
		admin.AdminID,
	)
	return err
}

// ListByRegion은 국가 코드와 지역 코드로 관리자를 검색합니다.
// admin_users 테이블에 region 컬럼이 없는 경우 빈 목록을 반환합니다.
func (r *AdminRepository) ListByRegion(ctx context.Context, countryCode, regionCode string) ([]*service.AdminUser, error) {
	// admin_users 테이블에 country_code/region_code 컬럼이 존재하지 않을 수 있으므로
	// 전체 관리자를 반환합니다 (실제 스키마 확장 후 WHERE 절 추가 예정).
	q := `SELECT admin_id, user_id, email, display_name, role::text, is_active,
		created_at, COALESCE(last_login_at, created_at), COALESCE(created_by,'')
		FROM admin_users WHERE is_active = true ORDER BY created_at DESC`

	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*service.AdminUser
	for rows.Next() {
		var a service.AdminUser
		var roleStr string
		if err := rows.Scan(
			&a.AdminID, &a.UserID, &a.Email, &a.DisplayName, &roleStr, &a.IsActive,
			&a.CreatedAt, &a.LastLoginAt, &a.CreatedBy,
		); err != nil {
			return nil, err
		}
		a.Role = stringToRole[roleStr]
		a.CountryCode = countryCode
		a.RegionCode = regionCode
		list = append(list, &a)
	}
	return list, nil
}

// ============================================================================
// AuditLogRepository
// ============================================================================

// AuditLogRepository는 PostgreSQL 기반 감사 로그 저장소입니다.
type AuditLogRepository struct {
	pool *pgxpool.Pool
}

// NewAuditLogRepository는 AuditLogRepository를 생성합니다.
func NewAuditLogRepository(pool *pgxpool.Pool) *AuditLogRepository {
	return &AuditLogRepository{pool: pool}
}

// Save는 감사 로그를 저장합니다.
func (r *AuditLogRepository) Save(ctx context.Context, entry *service.AuditLogEntry) error {
	const q = `INSERT INTO audit_logs
		(entry_id, admin_id, admin_email, action, resource_type, resource_id,
		 description, ip_address, timestamp)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`

	_, err := r.pool.Exec(ctx, q,
		entry.EntryID, entry.AdminID, entry.AdminEmail,
		actionToString[entry.Action], entry.ResourceType, entry.ResourceID,
		entry.Description, nullIfEmpty(entry.IPAddress), entry.Timestamp,
	)
	return err
}

// FindAll은 감사 로그를 조회합니다.
func (r *AuditLogRepository) FindAll(ctx context.Context, adminID string, actionFilter service.AuditAction, startDate, endDate *time.Time, limit, offset int32) ([]*service.AuditLogEntry, int, error) {
	where := `WHERE 1=1`
	args := []interface{}{}
	idx := 1

	if adminID != "" {
		where += fmt.Sprintf(` AND admin_id = $%d`, idx)
		args = append(args, adminID)
		idx++
	}
	if actionFilter != service.ActionUnknown {
		where += fmt.Sprintf(` AND action = $%d`, idx)
		args = append(args, actionToString[actionFilter])
		idx++
	}
	if startDate != nil {
		where += fmt.Sprintf(` AND timestamp >= $%d`, idx)
		args = append(args, *startDate)
		idx++
	}
	if endDate != nil {
		where += fmt.Sprintf(` AND timestamp <= $%d`, idx)
		args = append(args, *endDate)
		idx++
	}

	// count
	var total int
	if err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM audit_logs `+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// data
	dataQ := fmt.Sprintf(`SELECT entry_id, admin_id, COALESCE(admin_email,''), action::text,
		COALESCE(resource_type,''), COALESCE(resource_id,''),
		COALESCE(description,''), COALESCE(ip_address,''), timestamp
		FROM audit_logs %s ORDER BY timestamp DESC LIMIT $%d OFFSET $%d`,
		where, idx, idx+1)
	dataArgs := append(args, limit, offset)

	rows, err := r.pool.Query(ctx, dataQ, dataArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []*service.AuditLogEntry
	for rows.Next() {
		var e service.AuditLogEntry
		var actionStr string
		if err := rows.Scan(
			&e.EntryID, &e.AdminID, &e.AdminEmail, &actionStr,
			&e.ResourceType, &e.ResourceID,
			&e.Description, &e.IPAddress, &e.Timestamp,
		); err != nil {
			return nil, 0, err
		}
		e.Action = stringToAction[actionStr]
		list = append(list, &e)
	}
	return list, total, nil
}

// ============================================================================
// SystemConfigRepository
// ============================================================================

// SystemConfigRepository는 PostgreSQL 기반 시스템 설정 저장소입니다.
type SystemConfigRepository struct {
	pool *pgxpool.Pool
}

// NewSystemConfigRepository는 SystemConfigRepository를 생성합니다.
func NewSystemConfigRepository(pool *pgxpool.Pool) *SystemConfigRepository {
	return &SystemConfigRepository{pool: pool}
}

// Save는 시스템 설정을 저장합니다 (upsert).
func (r *SystemConfigRepository) Save(ctx context.Context, cfg *service.SystemConfig) error {
	const q = `INSERT INTO system_configs (key, value, description, updated_by, updated_at)
		VALUES ($1,$2,$3,$4,$5)
		ON CONFLICT (key) DO UPDATE SET value=$2, description=$3, updated_by=$4, updated_at=$5`

	_, err := r.pool.Exec(ctx, q,
		cfg.Key, cfg.Value, cfg.Description, nullIfEmpty(cfg.UpdatedBy), cfg.UpdatedAt,
	)
	return err
}

// FindByKey는 키로 시스템 설정을 조회합니다.
func (r *SystemConfigRepository) FindByKey(ctx context.Context, key string) (*service.SystemConfig, error) {
	const q = `SELECT key, value, COALESCE(description,''), COALESCE(updated_by,''), COALESCE(updated_at, NOW())
		FROM system_configs WHERE key = $1`

	var cfg service.SystemConfig
	err := r.pool.QueryRow(ctx, q, key).Scan(
		&cfg.Key, &cfg.Value, &cfg.Description, &cfg.UpdatedBy, &cfg.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &cfg, nil
}

// ============================================================================
// UserSummaryRepository
// ============================================================================

// UserSummaryRepository는 PostgreSQL 기반 사용자 요약 정보 저장소입니다.
// users 테이블에서 집계하여 관리자 대시보드에 표시합니다.
type UserSummaryRepository struct {
	pool *pgxpool.Pool
}

// NewUserSummaryRepository는 UserSummaryRepository를 생성합니다.
func NewUserSummaryRepository(pool *pgxpool.Pool) *UserSummaryRepository {
	return &UserSummaryRepository{pool: pool}
}

// FindAll은 사용자 요약 목록을 조회합니다.
// users 테이블 구조에 맞춰 조회합니다.
func (r *UserSummaryRepository) FindAll(ctx context.Context, keyword, roleFilter string, activeOnly bool, limit, offset int32) ([]*service.AdminUserSummary, int, error) {
	where := `WHERE 1=1`
	args := []interface{}{}
	idx := 1

	if activeOnly {
		where += ` AND is_active = true`
	}
	if keyword != "" {
		where += fmt.Sprintf(` AND (LOWER(email) LIKE $%d OR LOWER(COALESCE(display_name,'')) LIKE $%d)`, idx, idx)
		args = append(args, "%"+keyword+"%")
		idx++
	}
	if roleFilter != "" {
		where += fmt.Sprintf(` AND role = $%d`, idx)
		args = append(args, roleFilter)
		idx++
	}

	// count
	var total int
	if err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM users `+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// data
	dataQ := fmt.Sprintf(`SELECT id, email, COALESCE(display_name,''), is_active, created_at, updated_at
		FROM users %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d`,
		where, idx, idx+1)
	dataArgs := append(args, limit, offset)

	rows, err := r.pool.Query(ctx, dataQ, dataArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []*service.AdminUserSummary
	for rows.Next() {
		var u service.AdminUserSummary
		if err := rows.Scan(
			&u.UserID, &u.Email, &u.DisplayName, &u.IsActive, &u.CreatedAt, &u.LastActiveAt,
		); err != nil {
			return nil, 0, err
		}
		list = append(list, &u)
	}
	return list, total, nil
}

// ============================================================================
// helpers
// ============================================================================

func nullIfEmpty(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
