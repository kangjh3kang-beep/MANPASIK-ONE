package llm

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setupQuotaHTTP(t *testing.T) (*http.ServeMux, *MemoryQuotaStore) {
	t.Helper()
	store := NewMemoryQuotaStore()
	handler := NewQuotaHTTPHandler(store)
	handler.SetPathPrefix("/ops/tenancy")
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	return mux, store
}

func TestQuotaHTTPHandler_GetNotFound(t *testing.T) {
	mux, _ := setupQuotaHTTP(t)
	r := httptest.NewRequest("GET", "/ops/tenancy/quota/hospA", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusNotFound {
		t.Errorf("Code = %d", w.Code)
	}
}

func TestQuotaHTTPHandler_PutAndGet(t *testing.T) {
	mux, store := setupQuotaHTTP(t)

	// PUT
	r := httptest.NewRequest("PUT", "/ops/tenancy/quota/hospA",
		strings.NewReader(`{"daily_token_limit":10000,"monthly_token_limit":300000}`))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("PUT Code = %d, body = %s", w.Code, w.Body.String())
	}

	// store 확인
	cfg, err := store.Get(context.Background(), "hospA")
	if err != nil {
		t.Fatal(err)
	}
	if cfg.DailyTokenLimit != 10000 {
		t.Errorf("DailyTokenLimit = %d", cfg.DailyTokenLimit)
	}

	// GET
	r2 := httptest.NewRequest("GET", "/ops/tenancy/quota/hospA", nil)
	w2 := httptest.NewRecorder()
	mux.ServeHTTP(w2, r2)
	if w2.Code != http.StatusOK {
		t.Errorf("GET Code = %d", w2.Code)
	}
	var got QuotaConfig
	_ = json.Unmarshal(w2.Body.Bytes(), &got)
	if got.MonthlyTokenLimit != 300000 {
		t.Errorf("MonthlyTokenLimit = %d", got.MonthlyTokenLimit)
	}
}

func TestQuotaHTTPHandler_Delete(t *testing.T) {
	mux, store := setupQuotaHTTP(t)
	_ = store.Set(context.Background(), QuotaConfig{TenantID: "hospA", DailyTokenLimit: 100})

	r := httptest.NewRequest("DELETE", "/ops/tenancy/quota/hospA", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusNoContent {
		t.Errorf("Code = %d", w.Code)
	}
	if _, err := store.Get(context.Background(), "hospA"); err != ErrQuotaConfigNotFound {
		t.Error("delete 후에도 존재")
	}
}

func TestQuotaHTTPHandler_List(t *testing.T) {
	mux, store := setupQuotaHTTP(t)
	_ = store.Set(context.Background(), QuotaConfig{TenantID: "A", DailyTokenLimit: 100})
	_ = store.Set(context.Background(), QuotaConfig{TenantID: "B", DailyTokenLimit: 200})

	r := httptest.NewRequest("GET", "/ops/tenancy/quota", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("Code = %d", w.Code)
	}
	var resp struct {
		Configs []QuotaConfig `json:"configs"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if len(resp.Configs) != 2 {
		t.Errorf("configs = %d", len(resp.Configs))
	}
}

func TestQuotaHTTPHandler_PutBadJSON(t *testing.T) {
	mux, _ := setupQuotaHTTP(t)
	r := httptest.NewRequest("PUT", "/ops/tenancy/quota/hospA", strings.NewReader("bad json"))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Code = %d", w.Code)
	}
}

func TestQuotaHTTPHandler_PutNegativeRejected(t *testing.T) {
	mux, _ := setupQuotaHTTP(t)
	r := httptest.NewRequest("PUT", "/ops/tenancy/quota/hospA",
		strings.NewReader(`{"daily_token_limit":-100}`))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Code = %d", w.Code)
	}
}

func TestQuotaHTTPHandler_MethodNotAllowed(t *testing.T) {
	mux, _ := setupQuotaHTTP(t)
	r := httptest.NewRequest("PATCH", "/ops/tenancy/quota/hospA", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Code = %d", w.Code)
	}
}

func TestQuotaHTTPHandler_InvalidatesDynamicCache(t *testing.T) {
	store := NewMemoryQuotaStore()
	dynamic := NewDynamicQuota(store, nil, 60*0)
	handler := NewQuotaHTTPHandler(store)
	handler.SetDynamicQuota(dynamic)
	handler.SetPathPrefix("")
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	// PUT → 캐시 invalidate (panic 안 나면 OK)
	r := httptest.NewRequest("PUT", "/quota/hospA",
		strings.NewReader(`{"daily_token_limit":1000}`))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("Code = %d", w.Code)
	}
}

// PostgresQuotaStore 통합 테스트 (인메모리 fake)
type fakeQuotaQuerier struct {
	row    *fakeRow
	rows   *fakeQuotaRows
}

func (q *fakeQuotaQuerier) QueryRowContext(_ context.Context, _ string, _ ...interface{}) AuditRow {
	return q.row
}
func (q *fakeQuotaQuerier) QueryContext(_ context.Context, _ string, _ ...interface{}) (QuotaRows, error) {
	return q.rows, nil
}

type fakeQuotaRows struct {
	configs []*QuotaConfig
	idx     int
}

func (r *fakeQuotaRows) Next() bool {
	if r.idx >= len(r.configs) {
		return false
	}
	r.idx++
	return true
}

func (r *fakeQuotaRows) Scan(dest ...interface{}) error {
	c := r.configs[r.idx-1]
	*dest[0].(*string) = c.TenantID
	*dest[1].(*int) = c.DailyTokenLimit
	*dest[2].(*int) = c.MonthlyTokenLimit
	*dest[3].(*int) = c.DailyRequestLimit
	return nil
}

func (r *fakeQuotaRows) Close() error { return nil }
func (r *fakeQuotaRows) Err() error   { return nil }

func TestPostgresQuotaStore_Validation(t *testing.T) {
	if _, err := NewPostgresQuotaStore(nil, nil); err == nil {
		t.Error("nil exec/querier 통과")
	}
}

func TestPostgresQuotaStore_SetExec(t *testing.T) {
	exec := &fakeExec{}
	q := &fakeQuotaQuerier{row: &fakeRow{}}
	store, _ := NewPostgresQuotaStore(exec, q)
	cfg := QuotaConfig{TenantID: "hospA", DailyTokenLimit: 5000}
	if err := store.Set(context.Background(), cfg); err != nil {
		t.Fatal(err)
	}
	if len(exec.calls) != 1 {
		t.Errorf("exec calls = %d", len(exec.calls))
	}
	if !strings.Contains(exec.calls[0].sql, "ON CONFLICT") {
		t.Errorf("UPSERT 미사용: %s", exec.calls[0].sql)
	}
}

func TestPostgresQuotaStore_SetEmptyTenantID(t *testing.T) {
	exec := &fakeExec{}
	q := &fakeQuotaQuerier{}
	store, _ := NewPostgresQuotaStore(exec, q)
	if err := store.Set(context.Background(), QuotaConfig{}); err == nil {
		t.Error("빈 TenantID 통과")
	}
}
