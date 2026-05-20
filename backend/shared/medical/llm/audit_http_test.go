package llm

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// fakeStatsRow 는 StatsByTenant 의 5개 필드를 채움.
type fakeStatsRow struct {
	total, success, failed, totalTok, windowTok int64
	scanErr                                     error
}

func (r *fakeStatsRow) Scan(dest ...interface{}) error {
	if r.scanErr != nil {
		return r.scanErr
	}
	if len(dest) < 5 {
		return errors.New("dest len")
	}
	*dest[0].(*int64) = r.total
	*dest[1].(*int64) = r.success
	*dest[2].(*int64) = r.failed
	*dest[3].(*int64) = r.totalTok
	*dest[4].(*int64) = r.windowTok
	return nil
}

type statsQuerier struct {
	row *fakeStatsRow
}

func (q *statsQuerier) QueryRowContext(_ context.Context, _ string, _ ...interface{}) AuditRow {
	return q.row
}

func TestAuditHTTPHandler_StatsBasic(t *testing.T) {
	exec := &fakeExec{}
	q := &statsQuerier{row: &fakeStatsRow{
		total: 100, success: 95, failed: 5, totalTok: 50000, windowTok: 12000,
	}}
	pgLog, _ := NewPostgresAuditLog(exec, q)
	h := NewAuditHTTPHandler(pgLog)
	h.SetPathPrefix("/ops/tenancy")
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	r := httptest.NewRequest("GET", "/ops/tenancy/audit/stats?tenant=hospA&hours=24", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("Code = %d, body = %s", w.Code, w.Body.String())
	}
	var stats AuditTenantStats
	_ = json.Unmarshal(w.Body.Bytes(), &stats)
	if stats.TotalCalls != 100 || stats.SuccessfulCalls != 95 {
		t.Errorf("stats = %+v", stats)
	}
	if stats.WindowHours != 24 {
		t.Errorf("WindowHours = %d", stats.WindowHours)
	}
}

func TestAuditHTTPHandler_MissingTenant(t *testing.T) {
	pgLog, _ := NewPostgresAuditLog(&fakeExec{}, &statsQuerier{row: &fakeStatsRow{}})
	h := NewAuditHTTPHandler(pgLog)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	r := httptest.NewRequest("GET", "/audit/stats", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Code = %d", w.Code)
	}
}

func TestAuditHTTPHandler_DefaultHours(t *testing.T) {
	q := &statsQuerier{row: &fakeStatsRow{}}
	pgLog, _ := NewPostgresAuditLog(&fakeExec{}, q)
	h := NewAuditHTTPHandler(pgLog)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	r := httptest.NewRequest("GET", "/audit/stats?tenant=x", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Errorf("Code = %d", w.Code)
	}
	var stats AuditTenantStats
	_ = json.Unmarshal(w.Body.Bytes(), &stats)
	if stats.WindowHours != 24 {
		t.Errorf("기본 hours = %d, want 24", stats.WindowHours)
	}
}

func TestAuditHTTPHandler_CustomHours(t *testing.T) {
	q := &statsQuerier{row: &fakeStatsRow{}}
	pgLog, _ := NewPostgresAuditLog(&fakeExec{}, q)
	h := NewAuditHTTPHandler(pgLog)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	r := httptest.NewRequest("GET", "/audit/stats?tenant=x&hours=168", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	var stats AuditTenantStats
	_ = json.Unmarshal(w.Body.Bytes(), &stats)
	if stats.WindowHours != 168 {
		t.Errorf("hours = %d, want 168", stats.WindowHours)
	}
}

func TestAuditHTTPHandler_InvalidHours(t *testing.T) {
	q := &statsQuerier{row: &fakeStatsRow{}}
	pgLog, _ := NewPostgresAuditLog(&fakeExec{}, q)
	h := NewAuditHTTPHandler(pgLog)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	// 음수/잘못된 값 → 기본값 24
	r := httptest.NewRequest("GET", "/audit/stats?tenant=x&hours=-5", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	var stats AuditTenantStats
	_ = json.Unmarshal(w.Body.Bytes(), &stats)
	if stats.WindowHours != 24 {
		t.Errorf("음수 hours 무시 안됨: %d", stats.WindowHours)
	}
}

func TestAuditHTTPHandler_NilAuditLog(t *testing.T) {
	h := NewAuditHTTPHandler(nil)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	r := httptest.NewRequest("GET", "/audit/stats?tenant=x", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("Code = %d", w.Code)
	}
}

func TestAuditHTTPHandler_MethodNotAllowed(t *testing.T) {
	q := &statsQuerier{row: &fakeStatsRow{}}
	pgLog, _ := NewPostgresAuditLog(&fakeExec{}, q)
	h := NewAuditHTTPHandler(pgLog)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	r := httptest.NewRequest("POST", "/audit/stats?tenant=x", strings.NewReader(""))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Code = %d", w.Code)
	}
}

func TestPostgresAuditLog_StatsByTenant_QuerierError(t *testing.T) {
	q := &statsQuerier{row: &fakeStatsRow{scanErr: errors.New("db down")}}
	pgLog, _ := NewPostgresAuditLog(&fakeExec{}, q)
	if _, err := pgLog.StatsByTenant(context.Background(), "x", 24); err == nil {
		t.Error("DB 에러 통과")
	}
}

func TestPostgresAuditLog_StatsByTenant_NoQuerier(t *testing.T) {
	pgLog, _ := NewPostgresAuditLog(&fakeExec{}, nil)
	if _, err := pgLog.StatsByTenant(context.Background(), "x", 24); err == nil {
		t.Error("querier nil 통과")
	}
}

// fakeFailureQuerier 는 QueryContext 로 실패 행 반환 (AN-3).
type fakeFailureRows struct {
	rows []AuditFailureEntry
	idx  int
}

func (r *fakeFailureRows) Next() bool {
	if r.idx >= len(r.rows) {
		return false
	}
	r.idx++
	return true
}

func (r *fakeFailureRows) Scan(dest ...interface{}) error {
	row := r.rows[r.idx-1]
	*dest[0].(*string) = row.TenantID
	if ns, ok := dest[1].(*sql.NullString); ok {
		*ns = sql.NullString{String: row.UserID, Valid: row.UserID != ""}
	}
	if ns, ok := dest[2].(*sql.NullString); ok {
		*ns = sql.NullString{String: row.Provider, Valid: row.Provider != ""}
	}
	if ns, ok := dest[3].(*sql.NullString); ok {
		*ns = sql.NullString{String: row.Model, Valid: row.Model != ""}
	}
	if ns, ok := dest[4].(*sql.NullString); ok {
		*ns = sql.NullString{String: row.ErrorMessage, Valid: row.ErrorMessage != ""}
	}
	*dest[5].(*time.Time) = row.RecordedAt
	return nil
}

func (r *fakeFailureRows) Close() error { return nil }
func (r *fakeFailureRows) Err() error   { return nil }

type fakeFailureQuerier struct {
	row  *fakeStatsRow // for QueryRowContext (StatsByTenant)
	rows *fakeFailureRows
}

func (q *fakeFailureQuerier) QueryRowContext(_ context.Context, _ string, _ ...interface{}) AuditRow {
	return q.row
}
func (q *fakeFailureQuerier) QueryContext(_ context.Context, _ string, _ ...interface{}) (QuotaRows, error) {
	return q.rows, nil
}

func TestPostgresAuditLog_RecentFailures_Basic(t *testing.T) {
	q := &fakeFailureQuerier{
		row: &fakeStatsRow{},
		rows: &fakeFailureRows{rows: []AuditFailureEntry{
			{TenantID: "hospA", UserID: "u1", Provider: "openai", Model: "gpt-4",
				ErrorMessage: "rate limited", RecordedAt: time.Now()},
			{TenantID: "hospA", UserID: "u2", Provider: "anthropic", Model: "claude-sonnet",
				ErrorMessage: "context too long", RecordedAt: time.Now().Add(-time.Hour)},
		}},
	}
	pgLog, _ := NewPostgresAuditLog(&fakeExec{}, q)
	failures, err := pgLog.RecentFailures(context.Background(), "hospA", 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(failures) != 2 {
		t.Errorf("failures = %d", len(failures))
	}
	if failures[0].ErrorMessage != "rate limited" {
		t.Errorf("err msg = %q", failures[0].ErrorMessage)
	}
}

func TestPostgresAuditLog_RecentFailures_DefaultLimit(t *testing.T) {
	// Default 10, > 100 clamped — only test default
	q := &fakeFailureQuerier{
		rows: &fakeFailureRows{rows: nil},
	}
	pgLog, _ := NewPostgresAuditLog(&fakeExec{}, q)
	failures, err := pgLog.RecentFailures(context.Background(), "x", 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(failures) != 0 {
		t.Errorf("len = %d", len(failures))
	}
}

func TestPostgresAuditLog_RecentFailures_NoQuerier(t *testing.T) {
	pgLog, _ := NewPostgresAuditLog(&fakeExec{}, nil)
	if _, err := pgLog.RecentFailures(context.Background(), "x", 10); err == nil {
		t.Error("querier nil 통과")
	}
}

func TestAuditHTTPHandler_FailuresEndpoint(t *testing.T) {
	q := &fakeFailureQuerier{
		rows: &fakeFailureRows{rows: []AuditFailureEntry{
			{TenantID: "hospA", ErrorMessage: "test error", RecordedAt: time.Now()},
		}},
	}
	pgLog, _ := NewPostgresAuditLog(&fakeExec{}, q)
	h := NewAuditHTTPHandler(pgLog)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	r := httptest.NewRequest("GET", "/audit/failures?tenant=hospA&limit=5", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("Code = %d, body = %s", w.Code, w.Body.String())
	}
	var resp struct {
		Tenant   string              `json:"tenant"`
		Failures []AuditFailureEntry `json:"failures"`
		Count    int                 `json:"count"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Count != 1 {
		t.Errorf("count = %d", resp.Count)
	}
	if resp.Tenant != "hospA" {
		t.Errorf("tenant = %q", resp.Tenant)
	}
}

func TestAuditHTTPHandler_FailuresMissingTenant(t *testing.T) {
	q := &fakeFailureQuerier{rows: &fakeFailureRows{}}
	pgLog, _ := NewPostgresAuditLog(&fakeExec{}, q)
	h := NewAuditHTTPHandler(pgLog)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	r := httptest.NewRequest("GET", "/audit/failures", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Code = %d", w.Code)
	}
}
