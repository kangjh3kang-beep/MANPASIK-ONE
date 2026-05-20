package tenancy

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// PostgresInvitationStore 는 SQLDB 인터페이스 기반 영속 InvitationStore.
//
// 41-tenancy.sql 의 tenant_invitations 테이블 사용.
//
// timestamp 칼럼은 TIMESTAMPTZ 이며, Go 의 time.Time 으로 매핑.
// pgx 와 database/sql 모두 time.Time scan 을 지원.
type PostgresInvitationStore struct {
	db SQLDB
}

// NewPostgresInvitationStore 생성. db=nil 이면 nil + 에러.
func NewPostgresInvitationStore(db SQLDB) (*PostgresInvitationStore, error) {
	if db == nil {
		return nil, errors.New("db 필수")
	}
	return &PostgresInvitationStore{db: db}, nil
}

// Add 는 신규 초대 등록 (token PK 충돌 시 에러).
func (s *PostgresInvitationStore) Add(inv Invitation) error {
	if inv.Token == "" {
		return errors.New("Token 필수")
	}
	_, err := s.db.ExecContext(context.Background(), `
		INSERT INTO tenant_invitations
		    (token, tenant_id, inviter_id, invitee_hint, role,
		     status, issued_at, expires_at, accepted_by, accepted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`,
		inv.Token, string(inv.TenantID), inv.InviterID, inv.InviteeHint, string(inv.Role),
		string(inv.Status), inv.IssuedAt, inv.ExpiresAt,
		nullableString(inv.AcceptedBy), nullableTime(inv.AcceptedAt),
	)
	if err != nil {
		return fmt.Errorf("Add: %w", err)
	}
	return nil
}

// Get 은 토큰으로 조회.
func (s *PostgresInvitationStore) Get(token string) (*Invitation, error) {
	row := s.db.QueryRowContext(context.Background(), `
		SELECT token, tenant_id, inviter_id, invitee_hint, role,
		       status, issued_at, expires_at, accepted_by, accepted_at
		FROM tenant_invitations
		WHERE token = $1
	`, token)
	return scanInvitation(row)
}

// Update 는 status, AcceptedBy, AcceptedAt 갱신.
//
// 다른 필드(tenant_id, role 등)는 불변으로 간주하여 UPDATE 대상에서 제외.
func (s *PostgresInvitationStore) Update(inv Invitation) error {
	res, err := s.db.ExecContext(context.Background(), `
		UPDATE tenant_invitations
		SET status = $1, accepted_by = $2, accepted_at = $3, expires_at = $4
		WHERE token = $5
	`,
		string(inv.Status), nullableString(inv.AcceptedBy),
		nullableTime(inv.AcceptedAt), inv.ExpiresAt, inv.Token,
	)
	if err != nil {
		return fmt.Errorf("Update: %w", err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrInvitationNotFound
	}
	return nil
}

// ListByTenant 는 조직별 초대 목록 (만료/수락 모두 포함).
//
// 활성 (pending+미만료) 만 원하면 호출자가 InvitationService.ListPendingByTenant 사용.
func (s *PostgresInvitationStore) ListByTenant(tenantID TenantID) []*Invitation {
	rows, err := s.db.QueryContext(context.Background(), `
		SELECT token, tenant_id, inviter_id, invitee_hint, role,
		       status, issued_at, expires_at, accepted_by, accepted_at
		FROM tenant_invitations
		WHERE tenant_id = $1
		ORDER BY issued_at DESC
	`, string(tenantID))
	if err != nil {
		return nil
	}
	defer rows.Close()

	var out []*Invitation
	for rows.Next() {
		inv, err := scanInvitationRow(rows)
		if err != nil {
			continue
		}
		out = append(out, inv)
	}
	return out
}

// scanInvitation 은 단일 행 매핑 (QueryRow).
func scanInvitation(row SQLRow) (*Invitation, error) {
	var (
		inv          Invitation
		tenantStr    string
		roleStr      string
		statusStr    string
		inviteeHint  sql.NullString
		acceptedBy   sql.NullString
		acceptedAt   sql.NullTime
	)
	err := row.Scan(
		&inv.Token, &tenantStr, &inv.InviterID, &inviteeHint, &roleStr,
		&statusStr, &inv.IssuedAt, &inv.ExpiresAt, &acceptedBy, &acceptedAt,
	)
	if err != nil {
		msg := err.Error()
		if msg == "no rows in result set" || msg == "sql: no rows in result set" {
			return nil, ErrInvitationNotFound
		}
		return nil, err
	}
	inv.TenantID = TenantID(tenantStr)
	inv.Role = TenantRole(roleStr)
	inv.Status = InvitationStatus(statusStr)
	if inviteeHint.Valid {
		inv.InviteeHint = inviteeHint.String
	}
	if acceptedBy.Valid {
		inv.AcceptedBy = acceptedBy.String
	}
	if acceptedAt.Valid {
		t := acceptedAt.Time
		inv.AcceptedAt = &t
	}
	return &inv, nil
}

// scanInvitationRow 는 다중 행 매핑 (Query).
func scanInvitationRow(rows SQLRows) (*Invitation, error) {
	var (
		inv          Invitation
		tenantStr    string
		roleStr      string
		statusStr    string
		inviteeHint  sql.NullString
		acceptedBy   sql.NullString
		acceptedAt   sql.NullTime
	)
	if err := rows.Scan(
		&inv.Token, &tenantStr, &inv.InviterID, &inviteeHint, &roleStr,
		&statusStr, &inv.IssuedAt, &inv.ExpiresAt, &acceptedBy, &acceptedAt,
	); err != nil {
		return nil, err
	}
	inv.TenantID = TenantID(tenantStr)
	inv.Role = TenantRole(roleStr)
	inv.Status = InvitationStatus(statusStr)
	if inviteeHint.Valid {
		inv.InviteeHint = inviteeHint.String
	}
	if acceptedBy.Valid {
		inv.AcceptedBy = acceptedBy.String
	}
	if acceptedAt.Valid {
		t := acceptedAt.Time
		inv.AcceptedAt = &t
	}
	return &inv, nil
}

func nullableString(v string) interface{} {
	if v == "" {
		return nil
	}
	return v
}

func nullableTime(t *time.Time) interface{} {
	if t == nil {
		return nil
	}
	return *t
}
