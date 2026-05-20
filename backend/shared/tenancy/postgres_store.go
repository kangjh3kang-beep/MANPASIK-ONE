package tenancy

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// SQLDB 는 pgx 또는 database/sql 호환 인터페이스.
//
// 외부 의존을 피하기 위해 최소 메서드만 정의 — 실제로는 pgxpool.Pool 또는
// *sql.DB 가 이 인터페이스를 만족.
type SQLDB interface {
	ExecContext(ctx context.Context, sql string, args ...interface{}) (SQLResult, error)
	QueryContext(ctx context.Context, sql string, args ...interface{}) (SQLRows, error)
	QueryRowContext(ctx context.Context, sql string, args ...interface{}) SQLRow
}

// SQLResult 는 ExecContext 결과 (RowsAffected 만 사용).
type SQLResult interface {
	RowsAffected() (int64, error)
}

// SQLRows 는 다중 행 반환.
type SQLRows interface {
	Next() bool
	Scan(dest ...interface{}) error
	Close() error
	Err() error
}

// SQLRow 는 단일 행 반환.
type SQLRow interface {
	Scan(dest ...interface{}) error
}

// PostgresMembershipStore 는 SQLDB 인터페이스 기반 영속 저장소.
//
// 41-tenancy.sql 의 tenant_memberships 테이블 사용.
type PostgresMembershipStore struct {
	db SQLDB
}

// NewPostgresMembershipStore 생성. db=nil 이면 nil + 에러 반환.
func NewPostgresMembershipStore(db SQLDB) (*PostgresMembershipStore, error) {
	if db == nil {
		return nil, errors.New("db 필수")
	}
	return &PostgresMembershipStore{db: db}, nil
}

// Add 는 새 멤버십 등록 (UPSERT 동작).
//
// PostgreSQL ON CONFLICT 사용 — 동일 (user_id, tenant_id) 이면 역할/active 갱신.
func (s *PostgresMembershipStore) Add(m Membership) error {
	if err := m.Validate(); err != nil {
		return err
	}
	if m.JoinedAt == 0 {
		m.JoinedAt = time.Now().Unix()
	}
	// 신규 추가는 기본 활성 (Active false 명시 시에만 비활성)
	active := m.Active
	if !active && m.JoinedAt == time.Now().Unix() {
		active = true
	}

	_, err := s.db.ExecContext(context.Background(), `
		INSERT INTO tenant_memberships (user_id, tenant_id, role, active, joined_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id, tenant_id) DO UPDATE
		SET role = EXCLUDED.role, active = EXCLUDED.active
	`, m.UserID, string(m.TenantID), string(m.Role), active, m.JoinedAt)
	if err != nil {
		return fmt.Errorf("Add: %w", err)
	}
	return nil
}

// Remove 는 멤버십 삭제.
func (s *PostgresMembershipStore) Remove(userID string, tenantID TenantID) error {
	_, err := s.db.ExecContext(context.Background(),
		`DELETE FROM tenant_memberships WHERE user_id = $1 AND tenant_id = $2`,
		userID, string(tenantID))
	return err
}

// Get 은 멤버십 조회.
func (s *PostgresMembershipStore) Get(userID string, tenantID TenantID) (*Membership, error) {
	row := s.db.QueryRowContext(context.Background(), `
		SELECT user_id, tenant_id, role, active, joined_at
		FROM tenant_memberships
		WHERE user_id = $1 AND tenant_id = $2
	`, userID, string(tenantID))
	m, err := scanMembership(row)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// ListUserTenants 는 사용자 소속 모든 조직.
func (s *PostgresMembershipStore) ListUserTenants(userID string) []*Membership {
	rows, err := s.db.QueryContext(context.Background(), `
		SELECT user_id, tenant_id, role, active, joined_at
		FROM tenant_memberships WHERE user_id = $1
	`, userID)
	if err != nil {
		return nil
	}
	defer rows.Close()
	return scanMembershipRows(rows)
}

// ListTenantMembers 는 조직의 모든 멤버.
func (s *PostgresMembershipStore) ListTenantMembers(tenantID TenantID) []*Membership {
	rows, err := s.db.QueryContext(context.Background(), `
		SELECT user_id, tenant_id, role, active, joined_at
		FROM tenant_memberships WHERE tenant_id = $1
	`, string(tenantID))
	if err != nil {
		return nil
	}
	defer rows.Close()
	return scanMembershipRows(rows)
}

// UpdateRole 멤버 역할 변경.
func (s *PostgresMembershipStore) UpdateRole(userID string, tenantID TenantID, newRole TenantRole) error {
	if !newRole.IsKnown() {
		return errors.New("알 수 없는 역할: " + string(newRole))
	}
	res, err := s.db.ExecContext(context.Background(), `
		UPDATE tenant_memberships SET role = $1
		WHERE user_id = $2 AND tenant_id = $3
	`, string(newRole), userID, string(tenantID))
	if err != nil {
		return fmt.Errorf("UpdateRole: %w", err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrNoMembership
	}
	return nil
}

// SetActive 활성/비활성.
func (s *PostgresMembershipStore) SetActive(userID string, tenantID TenantID, active bool) error {
	res, err := s.db.ExecContext(context.Background(), `
		UPDATE tenant_memberships SET active = $1
		WHERE user_id = $2 AND tenant_id = $3
	`, active, userID, string(tenantID))
	if err != nil {
		return fmt.Errorf("SetActive: %w", err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrNoMembership
	}
	return nil
}

// scanMembership 은 단일 행에서 Membership 으로 매핑.
func scanMembership(row SQLRow) (*Membership, error) {
	var m Membership
	var tenantStr, roleStr string
	if err := row.Scan(&m.UserID, &tenantStr, &roleStr, &m.Active, &m.JoinedAt); err != nil {
		// pgx.ErrNoRows 또는 sql.ErrNoRows 에 대한 일관 처리:
		// 호출자가 ErrNoMembership 으로 변환할 수 있도록 그대로 전달.
		// 단, 직관성을 위해 "no rows" 메시지 패턴이면 ErrNoMembership 변환.
		errStr := err.Error()
		if errStr == "no rows in result set" || errStr == "sql: no rows in result set" {
			return nil, ErrNoMembership
		}
		return nil, err
	}
	m.TenantID = TenantID(tenantStr)
	m.Role = TenantRole(roleStr)
	return &m, nil
}

// scanMembershipRows 는 다중 행 매핑 (skip on scan error).
func scanMembershipRows(rows SQLRows) []*Membership {
	out := []*Membership{}
	for rows.Next() {
		var m Membership
		var tenantStr, roleStr string
		if err := rows.Scan(&m.UserID, &tenantStr, &roleStr, &m.Active, &m.JoinedAt); err != nil {
			continue
		}
		m.TenantID = TenantID(tenantStr)
		m.Role = TenantRole(roleStr)
		mc := m
		out = append(out, &mc)
	}
	return out
}
