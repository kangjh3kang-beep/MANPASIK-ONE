package dashboard

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// HTTPHandler 는 Aggregator 를 HTTP 로 노출.
//
// 기본 엔드포인트 (prefix 없음):
//   - GET /healthz   → 200 (healthy/degraded) / 503 (failed/unknown)
//   - GET /readyz    → 200 (healthy) / 503 (degraded/failed)
//   - GET /dashboard → 전체 모듈 스냅샷 JSON
//   - GET /metrics   → Prometheus exposition format (text/plain)
//
// SetPathPrefix("/admin") 호출 시 모든 경로에 prefix 가 붙음
// (예: /admin/dashboard).
type HTTPHandler struct {
	agg *Aggregator
	// snapshotTimeout 는 HTTP 호출당 최대 집계 시간.
	snapshotTimeout time.Duration
	pathPrefix      string
}

// NewHTTPHandler 생성. 기본 timeout 5초.
func NewHTTPHandler(agg *Aggregator) *HTTPHandler {
	return &HTTPHandler{agg: agg, snapshotTimeout: 5 * time.Second}
}

// SetSnapshotTimeout 설정.
func (h *HTTPHandler) SetSnapshotTimeout(d time.Duration) {
	if d > 0 {
		h.snapshotTimeout = d
	}
}

// SetPathPrefix 는 모든 라우트에 붙을 prefix 설정 (예: "/ops").
//
// 빈 문자열이거나 "/" 면 prefix 없음.
func (h *HTTPHandler) SetPathPrefix(prefix string) {
	prefix = strings.TrimRight(prefix, "/")
	h.pathPrefix = prefix
}

// RegisterRoutes 는 net/http ServeMux 에 4개 엔드포인트 등록.
func (h *HTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	p := h.pathPrefix
	mux.HandleFunc(p+"/healthz", h.healthz)
	mux.HandleFunc(p+"/readyz", h.readyz)
	mux.HandleFunc(p+"/dashboard", h.dashboard)
	mux.HandleFunc(p+"/metrics", h.metrics)
}

func (h *HTTPHandler) snapshot(r *http.Request) DashboardSnapshot {
	ctx, cancel := context.WithTimeout(r.Context(), h.snapshotTimeout)
	defer cancel()
	return h.agg.Snapshot(ctx)
}

// healthz 는 라이브니스 — failed 외에는 모두 200.
func (h *HTTPHandler) healthz(w http.ResponseWriter, r *http.Request) {
	snap := h.snapshot(r)
	code := http.StatusOK
	if snap.OverallStatus == StatusFailed || snap.OverallStatus == StatusUnknown {
		code = http.StatusServiceUnavailable
	}
	writeJSON(w, code, map[string]interface{}{
		"status":    snap.OverallStatus,
		"timestamp": snap.GeneratedAt.Format(time.RFC3339),
	})
}

// readyz 는 레디니스 — degraded 부터 503.
func (h *HTTPHandler) readyz(w http.ResponseWriter, r *http.Request) {
	snap := h.snapshot(r)
	code := http.StatusOK
	if snap.OverallStatus == StatusDegraded || snap.OverallStatus == StatusFailed {
		code = http.StatusServiceUnavailable
	}
	writeJSON(w, code, map[string]interface{}{
		"status":    snap.OverallStatus,
		"timestamp": snap.GeneratedAt.Format(time.RFC3339),
	})
}

// dashboard 는 모듈별 상세.
func (h *HTTPHandler) dashboard(w http.ResponseWriter, r *http.Request) {
	snap := h.snapshot(r)
	body := map[string]interface{}{
		"overall_status": snap.OverallStatus,
		"counts": map[string]int{
			"healthy":  snap.HealthyCount,
			"degraded": snap.DegradedCount,
			"failed":   snap.FailedCount,
			"unknown":  snap.UnknownCount,
		},
		"modules":      moduleListJSON(snap.Modules),
		"generated_at": snap.GeneratedAt.Format(time.RFC3339),
	}
	writeJSON(w, http.StatusOK, body)
}

func moduleListJSON(ms []ModuleSnapshot) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(ms))
	for _, m := range ms {
		entry := map[string]interface{}{
			"name":       m.Name,
			"status":     m.Status,
			"message":    m.Message,
			"updated_at": m.UpdatedAt.Format(time.RFC3339),
		}
		if len(m.Metrics) > 0 {
			entry["metrics"] = m.Metrics
		}
		out = append(out, entry)
	}
	return out
}

// metrics 는 Prometheus text exposition.
//
// 노출 메트릭:
//   manpasik_module_status{module="X",status="healthy|degraded|failed|unknown"} 0|1
//   manpasik_module_count{status="healthy|..."}  N
//   manpasik_module_metric{module="X",key="latency_p99"} <value>  (수치 메트릭에 한해)
func (h *HTTPHandler) metrics(w http.ResponseWriter, r *http.Request) {
	snap := h.snapshot(r)
	w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")

	var sb strings.Builder
	sb.WriteString("# HELP manpasik_module_status 모듈 헬스 상태 (1=현재 상태)\n")
	sb.WriteString("# TYPE manpasik_module_status gauge\n")
	for _, m := range snap.Modules {
		for _, s := range []Status{StatusHealthy, StatusDegraded, StatusFailed, StatusUnknown} {
			val := 0
			if m.Status == s {
				val = 1
			}
			fmt.Fprintf(&sb, "manpasik_module_status{module=\"%s\",status=\"%s\"} %d\n",
				escapeLabel(m.Name), s, val)
		}
	}

	sb.WriteString("# HELP manpasik_module_count 상태별 모듈 개수\n")
	sb.WriteString("# TYPE manpasik_module_count gauge\n")
	fmt.Fprintf(&sb, "manpasik_module_count{status=\"healthy\"} %d\n", snap.HealthyCount)
	fmt.Fprintf(&sb, "manpasik_module_count{status=\"degraded\"} %d\n", snap.DegradedCount)
	fmt.Fprintf(&sb, "manpasik_module_count{status=\"failed\"} %d\n", snap.FailedCount)
	fmt.Fprintf(&sb, "manpasik_module_count{status=\"unknown\"} %d\n", snap.UnknownCount)

	sb.WriteString("# HELP manpasik_module_metric 모듈 수치 메트릭\n")
	sb.WriteString("# TYPE manpasik_module_metric gauge\n")
	for _, m := range snap.Modules {
		for k, v := range m.Metrics {
			if num, ok := numericValue(v); ok {
				fmt.Fprintf(&sb, "manpasik_module_metric{module=\"%s\",key=\"%s\"} %g\n",
					escapeLabel(m.Name), escapeLabel(k), num)
			}
		}
	}

	_, _ = w.Write([]byte(sb.String()))
}

func numericValue(v interface{}) (float64, bool) {
	switch x := v.(type) {
	case float64:
		return x, true
	case float32:
		return float64(x), true
	case int:
		return float64(x), true
	case int32:
		return float64(x), true
	case int64:
		return float64(x), true
	case uint:
		return float64(x), true
	case uint32:
		return float64(x), true
	case uint64:
		return float64(x), true
	}
	return 0, false
}

// escapeLabel 은 Prometheus 레이블 값 escape (백슬래시/줄바꿈/따옴표).
func escapeLabel(s string) string {
	r := strings.NewReplacer(`\`, `\\`, `"`, `\"`, "\n", `\n`)
	return r.Replace(s)
}

func writeJSON(w http.ResponseWriter, code int, body interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(body)
}
