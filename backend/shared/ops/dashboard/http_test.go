package dashboard_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/manpasik/backend/shared/ops/dashboard"
)

func newAggForHTTP() *dashboard.Aggregator {
	a := dashboard.NewAggregator()
	a.Register(dashboard.NewSimpleProvider("auth", func(ctx context.Context) dashboard.ModuleSnapshot {
		return dashboard.ModuleSnapshot{Status: dashboard.StatusHealthy, Metrics: map[string]interface{}{"qps": 12.3}}
	}))
	a.Register(dashboard.NewSimpleProvider("db", func(ctx context.Context) dashboard.ModuleSnapshot {
		return dashboard.ModuleSnapshot{Status: dashboard.StatusDegraded, Message: "lag 200ms",
			Metrics: map[string]interface{}{"lag_ms": 200}}
	}))
	return a
}

func doReq(t *testing.T, h *dashboard.HTTPHandler, path string) (*httptest.ResponseRecorder, []byte) {
	t.Helper()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	r := httptest.NewRequest("GET", path, nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	return w, w.Body.Bytes()
}

func TestHTTPHandler_Healthz_Healthy(t *testing.T) {
	a := dashboard.NewAggregator()
	a.Register(dashboard.NewSimpleProvider("ok", func(ctx context.Context) dashboard.ModuleSnapshot {
		return dashboard.ModuleSnapshot{Status: dashboard.StatusHealthy}
	}))
	h := dashboard.NewHTTPHandler(a)
	w, body := doReq(t, h, "/healthz")
	if w.Code != 200 {
		t.Errorf("Code = %d", w.Code)
	}
	var resp map[string]interface{}
	_ = json.Unmarshal(body, &resp)
	if resp["status"] != "healthy" {
		t.Errorf("status = %v", resp["status"])
	}
}

func TestHTTPHandler_Healthz_Degraded200(t *testing.T) {
	// degraded 는 라이브니스에서 200
	h := dashboard.NewHTTPHandler(newAggForHTTP())
	w, _ := doReq(t, h, "/healthz")
	if w.Code != 200 {
		t.Errorf("degraded healthz code = %d, want 200", w.Code)
	}
}

func TestHTTPHandler_Healthz_Failed503(t *testing.T) {
	a := dashboard.NewAggregator()
	a.Register(dashboard.NewSimpleProvider("dead", func(ctx context.Context) dashboard.ModuleSnapshot {
		return dashboard.ModuleSnapshot{Status: dashboard.StatusFailed}
	}))
	h := dashboard.NewHTTPHandler(a)
	w, _ := doReq(t, h, "/healthz")
	if w.Code != 503 {
		t.Errorf("failed healthz code = %d, want 503", w.Code)
	}
}

func TestHTTPHandler_Readyz_DegradedTo503(t *testing.T) {
	h := dashboard.NewHTTPHandler(newAggForHTTP())
	w, _ := doReq(t, h, "/readyz")
	if w.Code != 503 {
		t.Errorf("degraded readyz code = %d, want 503", w.Code)
	}
}

func TestHTTPHandler_Dashboard_JSON(t *testing.T) {
	h := dashboard.NewHTTPHandler(newAggForHTTP())
	w, body := doReq(t, h, "/dashboard")
	if w.Code != 200 {
		t.Errorf("Code = %d", w.Code)
	}
	var resp map[string]interface{}
	if err := json.Unmarshal(body, &resp); err != nil {
		t.Fatalf("JSON parse: %v", err)
	}
	if resp["overall_status"] != "degraded" {
		t.Errorf("overall = %v", resp["overall_status"])
	}
	mods, ok := resp["modules"].([]interface{})
	if !ok || len(mods) != 2 {
		t.Errorf("modules = %v", resp["modules"])
	}
}

func TestHTTPHandler_Metrics_Format(t *testing.T) {
	h := dashboard.NewHTTPHandler(newAggForHTTP())
	w, body := doReq(t, h, "/metrics")
	if w.Code != 200 {
		t.Errorf("Code = %d", w.Code)
	}
	ct := w.Header().Get("Content-Type")
	if !strings.HasPrefix(ct, "text/plain") {
		t.Errorf("Content-Type = %q", ct)
	}
	s := string(body)
	if !strings.Contains(s, "manpasik_module_status{module=\"auth\",status=\"healthy\"} 1") {
		t.Errorf("auth healthy metric 누락:\n%s", s)
	}
	if !strings.Contains(s, "manpasik_module_status{module=\"db\",status=\"degraded\"} 1") {
		t.Errorf("db degraded metric 누락:\n%s", s)
	}
	if !strings.Contains(s, "manpasik_module_count{status=\"degraded\"} 1") {
		t.Errorf("count metric 누락")
	}
	if !strings.Contains(s, "manpasik_module_metric{module=\"auth\",key=\"qps\"}") {
		t.Errorf("custom metric 누락:\n%s", s)
	}
}

func TestHTTPHandler_Metrics_LabelEscape(t *testing.T) {
	a := dashboard.NewAggregator()
	a.Register(dashboard.NewSimpleProvider(`name"with"quotes`, func(ctx context.Context) dashboard.ModuleSnapshot {
		return dashboard.ModuleSnapshot{Status: dashboard.StatusHealthy}
	}))
	h := dashboard.NewHTTPHandler(a)
	_, body := doReq(t, h, "/metrics")
	if !strings.Contains(string(body), `module="name\"with\"quotes"`) {
		t.Errorf("레이블 이스케이프 실패:\n%s", string(body))
	}
}
