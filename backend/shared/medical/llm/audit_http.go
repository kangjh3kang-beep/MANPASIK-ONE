package llm

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// AuditHTTPHandler 는 LLM Audit 통계 REST API (Phase AM-2).
//
// 엔드포인트:
//   GET /audit/stats?tenant=<id>&hours=<n>  → tenant 별 통합 통계
//
// 인증/권한 검사는 호출자 책임.
type AuditHTTPHandler struct {
	auditLog   *PostgresAuditLog
	pathPrefix string
}

// NewAuditHTTPHandler 생성. auditLog 는 PostgresAuditLog 만 지원
// (querier 가 필요한 통계 쿼리 사용).
func NewAuditHTTPHandler(auditLog *PostgresAuditLog) *AuditHTTPHandler {
	return &AuditHTTPHandler{auditLog: auditLog}
}

// SetPathPrefix 는 라우트 prefix 설정.
func (h *AuditHTTPHandler) SetPathPrefix(prefix string) {
	h.pathPrefix = strings.TrimRight(prefix, "/")
}

// RegisterRoutes 는 ServeMux 에 등록.
func (h *AuditHTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	p := h.pathPrefix
	mux.HandleFunc(p+"/audit/stats", h.handleStats)
	mux.HandleFunc(p+"/audit/failures", h.handleFailures)
}

func (h *AuditHTTPHandler) handleStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeAuditErr(w, http.StatusMethodNotAllowed, "GET only")
		return
	}
	if h.auditLog == nil {
		writeAuditErr(w, http.StatusServiceUnavailable, "audit log 미설정")
		return
	}
	tenantID := r.URL.Query().Get("tenant")
	if tenantID == "" {
		writeAuditErr(w, http.StatusBadRequest, "tenant 파라미터 필수")
		return
	}
	hours := 24
	if hStr := r.URL.Query().Get("hours"); hStr != "" {
		if v, err := strconv.Atoi(hStr); err == nil && v > 0 {
			hours = v
		}
	}
	stats, err := h.auditLog.StatsByTenant(r.Context(), tenantID, hours)
	if err != nil {
		writeAuditErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeAuditJSON(w, http.StatusOK, stats)
}

// handleFailures 는 GET /audit/failures?tenant=X&limit=N (Phase AN-3).
func (h *AuditHTTPHandler) handleFailures(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeAuditErr(w, http.StatusMethodNotAllowed, "GET only")
		return
	}
	if h.auditLog == nil {
		writeAuditErr(w, http.StatusServiceUnavailable, "audit log 미설정")
		return
	}
	tenantID := r.URL.Query().Get("tenant")
	if tenantID == "" {
		writeAuditErr(w, http.StatusBadRequest, "tenant 파라미터 필수")
		return
	}
	limit := 10
	if lStr := r.URL.Query().Get("limit"); lStr != "" {
		if v, err := strconv.Atoi(lStr); err == nil && v > 0 {
			limit = v
		}
	}
	failures, err := h.auditLog.RecentFailures(r.Context(), tenantID, limit)
	if err != nil {
		writeAuditErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeAuditJSON(w, http.StatusOK, map[string]interface{}{
		"tenant":   tenantID,
		"failures": failures,
		"count":    len(failures),
	})
}

func writeAuditJSON(w http.ResponseWriter, code int, body interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	if body != nil {
		_ = json.NewEncoder(w).Encode(body)
	}
}

func writeAuditErr(w http.ResponseWriter, code int, msg string) {
	writeAuditJSON(w, code, map[string]string{"error": msg})
}

// 사용되지 않은 import 방지용 (context 는 매개변수로 사용됨)
var _ = context.Background