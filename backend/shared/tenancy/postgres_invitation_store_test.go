package tenancy_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/manpasik/backend/shared/tenancy"
)

// fakeInvDB 는 SQLDB 인터페이스의 인메모리 가짜 구현 (invitation 테이블 한정).
type fakeInvDB struct {
	mu       sync.Mutex
	rows     map[string]*tenancy.Invitation
	failExec error
}

func newFakeInvDB() *fakeInvDB {
	return &fakeInvDB{rows: make(map[string]*tenancy.Invitation)}
}

type fakeInvResult struct{ affected int64 }

func (r fakeInvResult) RowsAffected() (int64, error) { return r.affected, nil }

type fakeInvRows struct {
	rows []*tenancy.Invitation
	idx  int
}

func (r *fakeInvRows) Next() bool {
	if r.idx >= len(r.rows) {
		return false
	}
	r.idx++
	return true
}

func (r *fakeInvRows) Scan(dest ...interface{}) error {
	row := r.rows[r.idx-1]
	if len(dest) < 10 {
		return errors.New("scan dest len")
	}
	*dest[0].(*string) = row.Token
	*dest[1].(*string) = string(row.TenantID)
	*dest[2].(*string) = row.InviterID

	hint, _ := dest[3].(*sql.NullString)
	if row.InviteeHint != "" {
		*hint = sql.NullString{String: row.InviteeHint, Valid: true}
	} else {
		*hint = sql.NullString{}
	}

	*dest[4].(*string) = string(row.Role)
	*dest[5].(*string) = string(row.Status)
	*dest[6].(*time.Time) = row.IssuedAt
	*dest[7].(*time.Time) = row.ExpiresAt

	accBy, _ := dest[8].(*sql.NullString)
	if row.AcceptedBy != "" {
		*accBy = sql.NullString{String: row.AcceptedBy, Valid: true}
	} else {
		*accBy = sql.NullString{}
	}

	accAt, _ := dest[9].(*sql.NullTime)
	if row.AcceptedAt != nil {
		*accAt = sql.NullTime{Time: *row.AcceptedAt, Valid: true}
	} else {
		*accAt = sql.NullTime{}
	}
	return nil
}

func (r *fakeInvRows) Close() error { return nil }
func (r *fakeInvRows) Err() error   { return nil }

type fakeInvRow struct{ row *tenancy.Invitation }

func (r *fakeInvRow) Scan(dest ...interface{}) error {
	if r.row == nil {
		return errors.New("no rows in result set")
	}
	rows := &fakeInvRows{rows: []*tenancy.Invitation{r.row}}
	rows.idx = 1 // Scan 은 현재 행 매핑
	return rows.Scan(dest...)
}

func (db *fakeInvDB) ExecContext(_ context.Context, sqlStr string, args ...interface{}) (tenancy.SQLResult, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	if db.failExec != nil {
		return nil, db.failExec
	}
	t := strings.TrimSpace(sqlStr)
	switch {
	case strings.HasPrefix(t, "INSERT INTO tenant_invitations"):
		// 10 args: token, tenant, inviter, hint, role, status, issued, expires, accBy, accAt
		token := args[0].(string)
		if _, exists := db.rows[token]; exists {
			return nil, errors.New("duplicate token")
		}
		inv := &tenancy.Invitation{
			Token:     token,
			TenantID:  tenancy.TenantID(args[1].(string)),
			InviterID: args[2].(string),
			Role:      tenancy.TenantRole(args[4].(string)),
			Status:    tenancy.InvitationStatus(args[5].(string)),
			IssuedAt:  args[6].(time.Time),
			ExpiresAt: args[7].(time.Time),
		}
		if hint, ok := args[3].(string); ok {
			inv.InviteeHint = hint
		}
		if accBy, ok := args[8].(string); ok {
			inv.AcceptedBy = accBy
		}
		if accAt, ok := args[9].(time.Time); ok {
			inv.AcceptedAt = &accAt
		}
		db.rows[token] = inv
		return fakeInvResult{affected: 1}, nil

	case strings.HasPrefix(t, "UPDATE tenant_invitations"):
		// 5 args: status, accBy, accAt, expires, token
		status := args[0].(string)
		token := args[4].(string)
		row, ok := db.rows[token]
		if !ok {
			return fakeInvResult{affected: 0}, nil
		}
		row.Status = tenancy.InvitationStatus(status)
		if accBy, ok := args[1].(string); ok {
			row.AcceptedBy = accBy
		}
		if accAt, ok := args[2].(time.Time); ok {
			row.AcceptedAt = &accAt
		}
		row.ExpiresAt = args[3].(time.Time)
		return fakeInvResult{affected: 1}, nil
	}
	return fakeInvResult{affected: 0}, nil
}

func (db *fakeInvDB) QueryContext(_ context.Context, sqlStr string, args ...interface{}) (tenancy.SQLRows, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	if !strings.Contains(sqlStr, "WHERE tenant_id = $1") {
		return &fakeInvRows{}, nil
	}
	tid := args[0].(string)
	var matched []*tenancy.Invitation
	for _, inv := range db.rows {
		if string(inv.TenantID) == tid {
			cp := *inv
			matched = append(matched, &cp)
		}
	}
	return &fakeInvRows{rows: matched}, nil
}

func (db *fakeInvDB) QueryRowContext(_ context.Context, _ string, args ...interface{}) tenancy.SQLRow {
	db.mu.Lock()
	defer db.mu.Unlock()
	token := args[0].(string)
	row := db.rows[token]
	if row == nil {
		return &fakeInvRow{}
	}
	cp := *row
	return &fakeInvRow{row: &cp}
}

func TestPostgresInvitationStore_AddGet(t *testing.T) {
	db := newFakeInvDB()
	s, err := tenancy.NewPostgresInvitationStore(db)
	if err != nil {
		t.Fatal(err)
	}

	now := time.Now()
	inv := tenancy.Invitation{
		Token:       "tok-1",
		TenantID:    "t1",
		InviterID:   "u1",
		InviteeHint: "user@kakao",
		Role:        tenancy.TenantRoleAdmin,
		Status:      tenancy.InvitationPending,
		IssuedAt:    now,
		ExpiresAt:   now.Add(7 * 24 * time.Hour),
	}
	if err := s.Add(inv); err != nil {
		t.Fatal(err)
	}

	got, err := s.Get("tok-1")
	if err != nil {
		t.Fatal(err)
	}
	if got.TenantID != "t1" || got.Role != tenancy.TenantRoleAdmin {
		t.Errorf("got = %+v", got)
	}
	if got.InviteeHint != "user@kakao" {
		t.Errorf("InviteeHint = %q", got.InviteeHint)
	}
}

func TestPostgresInvitationStore_GetNotFound(t *testing.T) {
	db := newFakeInvDB()
	s, _ := tenancy.NewPostgresInvitationStore(db)
	if _, err := s.Get("missing"); err != tenancy.ErrInvitationNotFound {
		t.Errorf("err = %v", err)
	}
}

func TestPostgresInvitationStore_AddDuplicate(t *testing.T) {
	db := newFakeInvDB()
	s, _ := tenancy.NewPostgresInvitationStore(db)
	now := time.Now()
	inv := tenancy.Invitation{
		Token: "dup", TenantID: "t", InviterID: "u",
		Role: tenancy.TenantRoleMember, Status: tenancy.InvitationPending,
		IssuedAt: now, ExpiresAt: now.Add(time.Hour),
	}
	if err := s.Add(inv); err != nil {
		t.Fatal(err)
	}
	if err := s.Add(inv); err == nil {
		t.Error("중복 token 통과")
	}
}

func TestPostgresInvitationStore_AddNoToken(t *testing.T) {
	db := newFakeInvDB()
	s, _ := tenancy.NewPostgresInvitationStore(db)
	if err := s.Add(tenancy.Invitation{}); err == nil {
		t.Error("빈 token 통과")
	}
}

func TestPostgresInvitationStore_Update(t *testing.T) {
	db := newFakeInvDB()
	s, _ := tenancy.NewPostgresInvitationStore(db)
	now := time.Now()
	inv := tenancy.Invitation{
		Token: "tok-2", TenantID: "t", InviterID: "u",
		Role: tenancy.TenantRoleMember, Status: tenancy.InvitationPending,
		IssuedAt: now, ExpiresAt: now.Add(time.Hour),
	}
	_ = s.Add(inv)

	accAt := now.Add(10 * time.Minute)
	inv.Status = tenancy.InvitationAccepted
	inv.AcceptedBy = "newUser"
	inv.AcceptedAt = &accAt
	if err := s.Update(inv); err != nil {
		t.Fatal(err)
	}

	got, _ := s.Get("tok-2")
	if got.Status != tenancy.InvitationAccepted {
		t.Errorf("Status = %q", got.Status)
	}
	if got.AcceptedBy != "newUser" {
		t.Errorf("AcceptedBy = %q", got.AcceptedBy)
	}
	if got.AcceptedAt == nil {
		t.Error("AcceptedAt nil")
	}
}

func TestPostgresInvitationStore_UpdateNotFound(t *testing.T) {
	db := newFakeInvDB()
	s, _ := tenancy.NewPostgresInvitationStore(db)
	if err := s.Update(tenancy.Invitation{Token: "ghost", ExpiresAt: time.Now()}); err != tenancy.ErrInvitationNotFound {
		t.Errorf("err = %v", err)
	}
}

func TestPostgresInvitationStore_ListByTenant(t *testing.T) {
	db := newFakeInvDB()
	s, _ := tenancy.NewPostgresInvitationStore(db)
	now := time.Now()
	for i, tid := range []string{"A", "A", "B"} {
		token := fmt.Sprintf("tok-%d", i)
		_ = s.Add(tenancy.Invitation{
			Token: token,
			TenantID: tenancy.TenantID(tid), InviterID: "u",
			Role: tenancy.TenantRoleMember, Status: tenancy.InvitationPending,
			IssuedAt: now, ExpiresAt: now.Add(time.Hour),
		})
	}
	list := s.ListByTenant("A")
	if len(list) != 2 {
		t.Errorf("ListByTenant len = %d", len(list))
	}
}

func TestPostgresInvitationStore_NilDB(t *testing.T) {
	if _, err := tenancy.NewPostgresInvitationStore(nil); err == nil {
		t.Error("nil DB 통과")
	}
}

func TestPostgresInvitationStore_ExecError(t *testing.T) {
	db := newFakeInvDB()
	db.failExec = errors.New("conn lost")
	s, _ := tenancy.NewPostgresInvitationStore(db)
	now := time.Now()
	if err := s.Add(tenancy.Invitation{
		Token: "x", TenantID: "t", InviterID: "u",
		Role: tenancy.TenantRoleMember, Status: tenancy.InvitationPending,
		IssuedAt: now, ExpiresAt: now,
	}); err == nil {
		t.Error("DB 에러 통과")
	}
}

// InvitationService 가 PostgresInvitationStore 와 함께 동작 검증
func TestInvitationService_WithPostgresStore(t *testing.T) {
	db := newFakeInvDB()
	invStore, _ := tenancy.NewPostgresInvitationStore(db)
	memStore := tenancy.NewMemoryMembershipStore()

	svc, _ := tenancy.NewInvitationService(invStore, memStore, tenancy.InvitationServiceConfig{
		TokenGenerator: fixedTokens("z-1"),
	})

	inv, err := svc.Invite(tenancy.InviteRequest{
		InviterID: "doc1", TenantID: "hospA", Role: tenancy.TenantRoleMedicalStaff,
	})
	if err != nil {
		t.Fatal(err)
	}

	m, err := svc.Accept(inv.Token, "newDoc")
	if err != nil {
		t.Fatal(err)
	}
	if m.Role != tenancy.TenantRoleMedicalStaff {
		t.Errorf("Role = %q", m.Role)
	}

	// Accept 후 status 가 영속화되었는지 확인
	got, err := invStore.Get(inv.Token)
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != tenancy.InvitationAccepted {
		t.Errorf("Status = %q", got.Status)
	}
}
