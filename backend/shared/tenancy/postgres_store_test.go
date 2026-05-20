package tenancy_test

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"

	"github.com/manpasik/backend/shared/tenancy"
)

// fakeDB 는 SQLDB 인터페이스의 테스트용 인메모리 구현.
//
// CRUD 동작을 SQL 분석 없이 인메모리 map 으로 흉내 — UPSERT/SELECT/DELETE/UPDATE
// 패턴을 SQL 문자열 prefix 로 분기.
type fakeDB struct {
	mu       sync.Mutex
	rows     map[string]*tenancy.Membership // key = userID::tenantID
	failExec error
	failQuery error
}

func newFakeDB() *fakeDB {
	return &fakeDB{rows: make(map[string]*tenancy.Membership)}
}

func key(userID, tenantID string) string { return userID + "::" + tenantID }

type fakeResult struct{ affected int64 }

func (r fakeResult) RowsAffected() (int64, error) { return r.affected, nil }

type fakeRows struct {
	rows []*tenancy.Membership
	idx  int
}

func (r *fakeRows) Next() bool {
	if r.idx >= len(r.rows) {
		return false
	}
	r.idx++
	return true
}

func (r *fakeRows) Scan(dest ...interface{}) error {
	row := r.rows[r.idx-1]
	if len(dest) < 5 {
		return errors.New("scan dest len")
	}
	*dest[0].(*string) = row.UserID
	*dest[1].(*string) = string(row.TenantID)
	*dest[2].(*string) = string(row.Role)
	*dest[3].(*bool) = row.Active
	*dest[4].(*int64) = row.JoinedAt
	return nil
}

func (r *fakeRows) Close() error { return nil }
func (r *fakeRows) Err() error   { return nil }

type fakeRow struct {
	row *tenancy.Membership
}

func (r *fakeRow) Scan(dest ...interface{}) error {
	if r.row == nil {
		return errors.New("no rows in result set")
	}
	*dest[0].(*string) = r.row.UserID
	*dest[1].(*string) = string(r.row.TenantID)
	*dest[2].(*string) = string(r.row.Role)
	*dest[3].(*bool) = r.row.Active
	*dest[4].(*int64) = r.row.JoinedAt
	return nil
}

func (db *fakeDB) ExecContext(_ context.Context, sql string, args ...interface{}) (tenancy.SQLResult, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	if db.failExec != nil {
		return nil, db.failExec
	}

	switch {
	case strings.HasPrefix(strings.TrimSpace(sql), "INSERT INTO tenant_memberships"):
		// args: user_id, tenant_id, role, active, joined_at
		uid := args[0].(string)
		tid := args[1].(string)
		role := args[2].(string)
		active := args[3].(bool)
		joinedAt := args[4].(int64)
		db.rows[key(uid, tid)] = &tenancy.Membership{
			UserID: uid, TenantID: tenancy.TenantID(tid),
			Role: tenancy.TenantRole(role), Active: active, JoinedAt: joinedAt,
		}
		return fakeResult{affected: 1}, nil

	case strings.HasPrefix(strings.TrimSpace(sql), "DELETE FROM tenant_memberships"):
		uid := args[0].(string)
		tid := args[1].(string)
		k := key(uid, tid)
		_, exists := db.rows[k]
		delete(db.rows, k)
		if exists {
			return fakeResult{affected: 1}, nil
		}
		return fakeResult{affected: 0}, nil

	case strings.Contains(sql, "UPDATE tenant_memberships SET role"):
		newRole := args[0].(string)
		uid := args[1].(string)
		tid := args[2].(string)
		row, ok := db.rows[key(uid, tid)]
		if !ok {
			return fakeResult{affected: 0}, nil
		}
		row.Role = tenancy.TenantRole(newRole)
		return fakeResult{affected: 1}, nil

	case strings.Contains(sql, "UPDATE tenant_memberships SET active"):
		active := args[0].(bool)
		uid := args[1].(string)
		tid := args[2].(string)
		row, ok := db.rows[key(uid, tid)]
		if !ok {
			return fakeResult{affected: 0}, nil
		}
		row.Active = active
		return fakeResult{affected: 1}, nil
	}
	return fakeResult{affected: 0}, nil
}

func (db *fakeDB) QueryContext(_ context.Context, sql string, args ...interface{}) (tenancy.SQLRows, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	if db.failQuery != nil {
		return nil, db.failQuery
	}
	var matched []*tenancy.Membership
	switch {
	case strings.Contains(sql, "WHERE user_id = $1"):
		uid := args[0].(string)
		for _, m := range db.rows {
			if m.UserID == uid {
				cp := *m
				matched = append(matched, &cp)
			}
		}
	case strings.Contains(sql, "WHERE tenant_id = $1"):
		tid := args[0].(string)
		for _, m := range db.rows {
			if string(m.TenantID) == tid {
				cp := *m
				matched = append(matched, &cp)
			}
		}
	}
	return &fakeRows{rows: matched}, nil
}

func (db *fakeDB) QueryRowContext(_ context.Context, _ string, args ...interface{}) tenancy.SQLRow {
	db.mu.Lock()
	defer db.mu.Unlock()
	uid := args[0].(string)
	tid := args[1].(string)
	row := db.rows[key(uid, tid)]
	if row == nil {
		return &fakeRow{}
	}
	cp := *row
	return &fakeRow{row: &cp}
}

func TestPostgresMembershipStore_AddGet(t *testing.T) {
	db := newFakeDB()
	s, err := tenancy.NewPostgresMembershipStore(db)
	if err != nil {
		t.Fatal(err)
	}

	if err := s.Add(tenancy.Membership{
		UserID: "u1", TenantID: "t1", Role: tenancy.TenantRoleAdmin,
	}); err != nil {
		t.Fatal(err)
	}

	got, err := s.Get("u1", "t1")
	if err != nil {
		t.Fatal(err)
	}
	if got.Role != tenancy.TenantRoleAdmin {
		t.Errorf("Role = %q", got.Role)
	}
}

func TestPostgresMembershipStore_GetNotFound(t *testing.T) {
	db := newFakeDB()
	s, _ := tenancy.NewPostgresMembershipStore(db)
	if _, err := s.Get("ghost", "t1"); err != tenancy.ErrNoMembership {
		t.Errorf("err = %v, want ErrNoMembership", err)
	}
}

func TestPostgresMembershipStore_AddValidation(t *testing.T) {
	db := newFakeDB()
	s, _ := tenancy.NewPostgresMembershipStore(db)
	// 빈 UserID
	if err := s.Add(tenancy.Membership{TenantID: "t1", Role: tenancy.TenantRoleMember}); err == nil {
		t.Error("validation 통과")
	}
}

func TestPostgresMembershipStore_Remove(t *testing.T) {
	db := newFakeDB()
	s, _ := tenancy.NewPostgresMembershipStore(db)
	_ = s.Add(tenancy.Membership{UserID: "u1", TenantID: "t1", Role: tenancy.TenantRoleMember})
	if err := s.Remove("u1", "t1"); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Get("u1", "t1"); err != tenancy.ErrNoMembership {
		t.Errorf("Remove 후 Get err = %v", err)
	}
}

func TestPostgresMembershipStore_ListUserTenants(t *testing.T) {
	db := newFakeDB()
	s, _ := tenancy.NewPostgresMembershipStore(db)
	_ = s.Add(tenancy.Membership{UserID: "u1", TenantID: "t1", Role: tenancy.TenantRoleAdmin})
	_ = s.Add(tenancy.Membership{UserID: "u1", TenantID: "t2", Role: tenancy.TenantRoleViewer})
	_ = s.Add(tenancy.Membership{UserID: "u2", TenantID: "t1", Role: tenancy.TenantRoleMember})

	list := s.ListUserTenants("u1")
	if len(list) != 2 {
		t.Errorf("len = %d, want 2", len(list))
	}
}

func TestPostgresMembershipStore_ListTenantMembers(t *testing.T) {
	db := newFakeDB()
	s, _ := tenancy.NewPostgresMembershipStore(db)
	_ = s.Add(tenancy.Membership{UserID: "u1", TenantID: "t1", Role: tenancy.TenantRoleAdmin})
	_ = s.Add(tenancy.Membership{UserID: "u2", TenantID: "t1", Role: tenancy.TenantRoleMember})

	list := s.ListTenantMembers("t1")
	if len(list) != 2 {
		t.Errorf("len = %d, want 2", len(list))
	}
}

func TestPostgresMembershipStore_UpdateRole(t *testing.T) {
	db := newFakeDB()
	s, _ := tenancy.NewPostgresMembershipStore(db)
	_ = s.Add(tenancy.Membership{UserID: "u1", TenantID: "t1", Role: tenancy.TenantRoleViewer})

	if err := s.UpdateRole("u1", "t1", tenancy.TenantRoleAdmin); err != nil {
		t.Fatal(err)
	}
	got, _ := s.Get("u1", "t1")
	if got.Role != tenancy.TenantRoleAdmin {
		t.Errorf("Role = %q", got.Role)
	}
}

func TestPostgresMembershipStore_UpdateRole_Unknown(t *testing.T) {
	db := newFakeDB()
	s, _ := tenancy.NewPostgresMembershipStore(db)
	_ = s.Add(tenancy.Membership{UserID: "u1", TenantID: "t1", Role: tenancy.TenantRoleViewer})
	if err := s.UpdateRole("u1", "t1", "ghost"); err == nil {
		t.Error("unknown role 통과")
	}
}

func TestPostgresMembershipStore_UpdateRole_NotFound(t *testing.T) {
	db := newFakeDB()
	s, _ := tenancy.NewPostgresMembershipStore(db)
	if err := s.UpdateRole("ghost", "t1", tenancy.TenantRoleAdmin); err != tenancy.ErrNoMembership {
		t.Errorf("err = %v", err)
	}
}

func TestPostgresMembershipStore_SetActive(t *testing.T) {
	db := newFakeDB()
	s, _ := tenancy.NewPostgresMembershipStore(db)
	_ = s.Add(tenancy.Membership{UserID: "u1", TenantID: "t1", Role: tenancy.TenantRoleMember})

	if err := s.SetActive("u1", "t1", false); err != nil {
		t.Fatal(err)
	}
	got, _ := s.Get("u1", "t1")
	if got.Active {
		t.Error("Active=false 적용 안됨")
	}
}

func TestPostgresMembershipStore_NilDB(t *testing.T) {
	if _, err := tenancy.NewPostgresMembershipStore(nil); err == nil {
		t.Error("nil DB 통과")
	}
}

func TestPostgresMembershipStore_ExecError(t *testing.T) {
	db := newFakeDB()
	db.failExec = errors.New("conn lost")
	s, _ := tenancy.NewPostgresMembershipStore(db)
	if err := s.Add(tenancy.Membership{
		UserID: "u1", TenantID: "t1", Role: tenancy.TenantRoleMember,
	}); err == nil {
		t.Error("DB 에러 통과")
	}
}

// PolicyEngine 과 함께 동작 검증 (영속 store 로 정책 평가)
func TestPolicyEngine_WithPostgresStore(t *testing.T) {
	db := newFakeDB()
	s, _ := tenancy.NewPostgresMembershipStore(db)
	_ = s.Add(tenancy.Membership{UserID: "doc1", TenantID: "hospA", Role: tenancy.TenantRoleAdmin})

	eng := tenancy.NewPolicyEngine(s)
	res := &tenancy.Resource{TenantID: "hospA", OwnerID: "patient1"}
	d := eng.Evaluate("doc1", "hospA", res, tenancy.ActionWrite)
	if !d.Allowed {
		t.Errorf("권한 거부: %v", d)
	}
}
