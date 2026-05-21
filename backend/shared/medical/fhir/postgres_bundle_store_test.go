package fhir_test

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/manpasik/backend/shared/medical/fhir"
)

// ============================================================================
// Mock SQLDB — PostgresBundleStore 단위 테스트용
// ============================================================================
//
// 실제 DB 없이 PostgresBundleStore 의 SQL 호출 패턴/에러 분기를 검증.

type mockResult struct{ rows int64 }

func (m mockResult) RowsAffected() (int64, error) { return m.rows, nil }

type mockRow struct {
	cols  []interface{}
	err   error
}

func (m *mockRow) Scan(dest ...interface{}) error {
	if m.err != nil {
		return m.err
	}
	if len(dest) != len(m.cols) {
		return errors.New("mock: column count mismatch")
	}
	for i, c := range m.cols {
		if err := assignAny(dest[i], c); err != nil {
			return err
		}
	}
	return nil
}

func assignAny(dest, src interface{}) error {
	switch d := dest.(type) {
	case *string:
		*d, _ = src.(string)
	case *[]byte:
		switch v := src.(type) {
		case []byte:
			*d = v
		case string:
			*d = []byte(v)
		}
	case *time.Time:
		if v, ok := src.(time.Time); ok {
			*d = v
		}
	case *int:
		if v, ok := src.(int); ok {
			*d = v
		}
	default:
		return errors.New("mock: unsupported destination type")
	}
	return nil
}

type mockRows struct {
	rows   [][]interface{}
	cursor int
	err    error
}

func (m *mockRows) Next() bool {
	m.cursor++
	return m.cursor <= len(m.rows)
}

func (m *mockRows) Scan(dest ...interface{}) error {
	if m.cursor < 1 || m.cursor > len(m.rows) {
		return errors.New("mock: out of range")
	}
	row := m.rows[m.cursor-1]
	if len(dest) != len(row) {
		return errors.New("mock: column count mismatch")
	}
	for i, c := range row {
		if err := assignAny(dest[i], c); err != nil {
			return err
		}
	}
	return nil
}

func (m *mockRows) Close() error { return nil }
func (m *mockRows) Err() error   { return m.err }

type mockDB struct {
	// 마지막 호출 기록 (검증용)
	lastExecSQL  string
	lastExecArgs []interface{}
	lastQuerySQL string
	lastQueryArgs []interface{}

	// 다음 호출 응답
	execRows     int64
	execErr      error
	rowResp      *mockRow
	rowsResp     *mockRows
}

func (m *mockDB) ExecContext(_ context.Context, sql string, args ...interface{}) (fhir.SQLResult, error) {
	m.lastExecSQL = sql
	m.lastExecArgs = args
	if m.execErr != nil {
		return nil, m.execErr
	}
	return mockResult{rows: m.execRows}, nil
}

func (m *mockDB) QueryContext(_ context.Context, sql string, args ...interface{}) (fhir.SQLRows, error) {
	m.lastQuerySQL = sql
	m.lastQueryArgs = args
	if m.rowsResp == nil {
		return &mockRows{}, nil
	}
	// 새 인스턴스 복제 (cursor 초기화)
	clone := &mockRows{rows: m.rowsResp.rows}
	return clone, nil
}

func (m *mockDB) QueryRowContext(_ context.Context, sql string, args ...interface{}) fhir.SQLRow {
	m.lastQuerySQL = sql
	m.lastQueryArgs = args
	if m.rowResp != nil {
		return m.rowResp
	}
	return &mockRow{err: errors.New("no rows in result set")}
}

// ============================================================================
// 테스트
// ============================================================================

func TestPostgresBundleStore_NilDB(t *testing.T) {
	if _, err := fhir.NewPostgresBundleStore(nil); err == nil {
		t.Error("nil DB 통과")
	}
}

func TestPostgresBundleStore_Save_UPSERT(t *testing.T) {
	db := &mockDB{execRows: 1}
	store, _ := fhir.NewPostgresBundleStore(db)
	b := fhir.NewBundle("collection")
	b.ID = "bundle-1"

	err := store.Save(context.Background(), b, "tenantA", "PAT1")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(db.lastExecSQL, "INSERT INTO bundle_store") {
		t.Errorf("INSERT 누락:\n%s", db.lastExecSQL)
	}
	if !strings.Contains(db.lastExecSQL, "ON CONFLICT") {
		t.Errorf("UPSERT (ON CONFLICT) 누락")
	}
	// 인자 순서: tenant_id, bundle_id, patient_id, payload
	if db.lastExecArgs[0] != "tenantA" {
		t.Errorf("arg[0] = %v, tenantA 기대", db.lastExecArgs[0])
	}
	if db.lastExecArgs[1] != "bundle-1" {
		t.Errorf("arg[1] = %v", db.lastExecArgs[1])
	}
	if db.lastExecArgs[2] != "PAT1" {
		t.Errorf("arg[2] = %v", db.lastExecArgs[2])
	}
	// payload 가 JSON byte 슬라이스
	if _, ok := db.lastExecArgs[3].([]byte); !ok {
		t.Errorf("arg[3] = %T (bytes 기대)", db.lastExecArgs[3])
	}
}

func TestPostgresBundleStore_Save_NilBundle(t *testing.T) {
	db := &mockDB{}
	store, _ := fhir.NewPostgresBundleStore(db)
	if err := store.Save(context.Background(), nil, "T", "P"); err == nil {
		t.Error("nil Bundle 통과")
	}
	b := fhir.NewBundle("collection")
	b.ID = ""
	if err := store.Save(context.Background(), b, "T", "P"); err == nil {
		t.Error("빈 Bundle.ID 통과")
	}
}

func TestPostgresBundleStore_Get_RoundTrip(t *testing.T) {
	bundleJSON, _ := json.Marshal(&fhir.Bundle{ID: "bundle-X", Type: "collection"})
	db := &mockDB{
		rowResp: &mockRow{
			cols: []interface{}{
				"bundle-X", "tenantA", "PAT9", bundleJSON, time.Now().UTC(),
			},
		},
	}
	store, _ := fhir.NewPostgresBundleStore(db)

	got, err := store.Get(context.Background(), "bundle-X", "tenantA")
	if err != nil {
		t.Fatal(err)
	}
	if got.Bundle.ID != "bundle-X" || got.TenantID != "tenantA" || got.PatientID != "PAT9" {
		t.Errorf("got = %+v", got)
	}
	// 쿼리 인자 검증
	if db.lastQueryArgs[0] != "tenantA" || db.lastQueryArgs[1] != "bundle-X" {
		t.Errorf("query args = %v", db.lastQueryArgs)
	}
}

func TestPostgresBundleStore_Get_NotFound(t *testing.T) {
	db := &mockDB{
		rowResp: &mockRow{err: errors.New("no rows in result set")},
	}
	store, _ := fhir.NewPostgresBundleStore(db)
	_, err := store.Get(context.Background(), "missing", "T")
	if !errors.Is(err, fhir.ErrBundleNotFound) {
		t.Errorf("err = %v, ErrBundleNotFound 기대", err)
	}
}

func TestPostgresBundleStore_Delete_NotFound(t *testing.T) {
	db := &mockDB{execRows: 0} // 0 rows affected
	store, _ := fhir.NewPostgresBundleStore(db)
	err := store.Delete(context.Background(), "missing", "T")
	if !errors.Is(err, fhir.ErrBundleNotFound) {
		t.Errorf("err = %v, ErrBundleNotFound 기대", err)
	}
}

func TestPostgresBundleStore_Delete_Success(t *testing.T) {
	db := &mockDB{execRows: 1}
	store, _ := fhir.NewPostgresBundleStore(db)
	err := store.Delete(context.Background(), "bundle-1", "T")
	if err != nil {
		t.Errorf("err = %v", err)
	}
	if !strings.Contains(db.lastExecSQL, "DELETE FROM bundle_store") {
		t.Errorf("DELETE 누락: %s", db.lastExecSQL)
	}
}

func TestPostgresBundleStore_ListByPatient_AppliesLimit(t *testing.T) {
	db := &mockDB{
		rowsResp: &mockRows{
			rows: [][]interface{}{
				{"b1", "T", "PAT1", mustBundleJSON("b1"), time.Now().UTC()},
				{"b2", "T", "PAT1", mustBundleJSON("b2"), time.Now().UTC()},
			},
		},
	}
	store, _ := fhir.NewPostgresBundleStore(db)
	list, err := store.ListByPatient(context.Background(), "T", "PAT1", 5)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 2 {
		t.Errorf("len = %d, 2 기대", len(list))
	}
	if !strings.Contains(db.lastQuerySQL, "LIMIT $3") {
		t.Errorf("LIMIT 누락: %s", db.lastQuerySQL)
	}
	if db.lastQueryArgs[2] != 5 {
		t.Errorf("limit arg = %v", db.lastQueryArgs[2])
	}
}

func TestPostgresBundleStore_ListByPatient_EmptyPatient(t *testing.T) {
	db := &mockDB{}
	store, _ := fhir.NewPostgresBundleStore(db)
	if _, err := store.ListByPatient(context.Background(), "T", "", 0); err == nil {
		t.Error("빈 patientID 통과")
	}
}

func TestPostgresBundleStore_ListByTenant_NoLimit(t *testing.T) {
	db := &mockDB{
		rowsResp: &mockRows{
			rows: [][]interface{}{
				{"b1", "T", "P", mustBundleJSON("b1"), time.Now().UTC()},
			},
		},
	}
	store, _ := fhir.NewPostgresBundleStore(db)
	list, err := store.ListByTenant(context.Background(), "T", 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 {
		t.Errorf("len = %d", len(list))
	}
	if strings.Contains(db.lastQuerySQL, "LIMIT") {
		t.Errorf("limit=0 인데 LIMIT 포함됨: %s", db.lastQuerySQL)
	}
}

func TestPostgresBundleStore_Count(t *testing.T) {
	db := &mockDB{
		rowResp: &mockRow{cols: []interface{}{42}},
	}
	store, _ := fhir.NewPostgresBundleStore(db)
	n := store.Count(context.Background(), "T")
	if n != 42 {
		t.Errorf("count = %d, 42 기대", n)
	}
}

func TestPostgresBundleStore_Save_DBError(t *testing.T) {
	db := &mockDB{execErr: errors.New("connection refused")}
	store, _ := fhir.NewPostgresBundleStore(db)
	b := fhir.NewBundle("collection")
	b.ID = "b1"
	err := store.Save(context.Background(), b, "T", "P")
	if err == nil {
		t.Error("DB 에러 통과")
	}
	if !strings.Contains(err.Error(), "Save") {
		t.Errorf("에러 컨텍스트 누락: %v", err)
	}
}

// 헬퍼.

func mustBundleJSON(id string) []byte {
	b, _ := json.Marshal(&fhir.Bundle{ID: id, Type: "collection"})
	return b
}
