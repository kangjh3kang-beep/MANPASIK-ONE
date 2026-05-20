package llm

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
)

// QuotaHTTPHandler 는 admin-service 등에서 마운트할 수 있는 HTTP 핸들러.
//
// 엔드포인트:
//   GET    /quota/{tenant_id}    config 조회 (없으면 404)
//   PUT    /quota/{tenant_id}    config UPSERT (body: QuotaConfig JSON)
//   DELETE /quota/{tenant_id}    config 제거
//   GET    /quota                전체 목록 (admin 전용)
//
// 인증/권한 검사는 호출자가 미들웨어로 별도 처리.
type QuotaHTTPHandler struct {
	store      QuotaStore
	pathPrefix string
	dynamic    *DynamicQuota // 옵션 — Set/Delete 시 cache invalidate
}

// NewQuotaHTTPHandler 생성.
func NewQuotaHTTPHandler(store QuotaStore) *QuotaHTTPHandler {
	return &QuotaHTTPHandler{store: store}
}

// SetPathPrefix 는 모든 엔드포인트 prefix 추가 (예: "/ops/tenancy").
func (h *QuotaHTTPHandler) SetPathPrefix(prefix string) {
	h.pathPrefix = strings.TrimRight(prefix, "/")
}

// SetDynamicQuota 는 cache invalidate 대상 등록 (Set/Delete 시 자동).
func (h *QuotaHTTPHandler) SetDynamicQuota(dq *DynamicQuota) {
	h.dynamic = dq
}

// RegisterRoutes 는 ServeMux 에 등록.
func (h *QuotaHTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	p := h.pathPrefix
	mux.HandleFunc(p+"/quota", h.handleList)
	mux.HandleFunc(p+"/quota/", h.handleByTenant)
}

func (h *QuotaHTTPHandler) handleList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeQuotaErr(w, http.StatusMethodNotAllowed, "GET only")
		return
	}
	configs, err := h.store.List(r.Context())
	if err != nil {
		writeQuotaErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeQuotaJSON(w, http.StatusOK, map[string]interface{}{"configs": configs})
}

func (h *QuotaHTTPHandler) handleByTenant(w http.ResponseWriter, r *http.Request) {
	tenantID := strings.TrimPrefix(r.URL.Path, h.pathPrefix+"/quota/")
	if tenantID == "" || strings.Contains(tenantID, "/") {
		writeQuotaErr(w, http.StatusBadRequest, "tenant_id 필수")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.serveGet(w, r, tenantID)
	case http.MethodPut:
		h.servePut(w, r, tenantID)
	case http.MethodDelete:
		h.serveDelete(w, r, tenantID)
	default:
		writeQuotaErr(w, http.StatusMethodNotAllowed, "GET/PUT/DELETE only")
	}
}

func (h *QuotaHTTPHandler) serveGet(w http.ResponseWriter, r *http.Request, tenantID string) {
	cfg, err := h.store.Get(r.Context(), tenantID)
	if err != nil {
		if errors.Is(err, ErrQuotaConfigNotFound) {
			writeQuotaErr(w, http.StatusNotFound, "quota not found")
			return
		}
		writeQuotaErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeQuotaJSON(w, http.StatusOK, cfg)
}

func (h *QuotaHTTPHandler) servePut(w http.ResponseWriter, r *http.Request, tenantID string) {
	var body struct {
		DailyTokenLimit   int `json:"daily_token_limit"`
		MonthlyTokenLimit int `json:"monthly_token_limit"`
		DailyRequestLimit int `json:"daily_request_limit"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeQuotaErr(w, http.StatusBadRequest, "JSON 본문 파싱 실패")
		return
	}
	if body.DailyTokenLimit < 0 || body.MonthlyTokenLimit < 0 || body.DailyRequestLimit < 0 {
		writeQuotaErr(w, http.StatusBadRequest, "한도는 음수일 수 없음")
		return
	}
	cfg := QuotaConfig{
		TenantID:          tenantID,
		DailyTokenLimit:   body.DailyTokenLimit,
		MonthlyTokenLimit: body.MonthlyTokenLimit,
		DailyRequestLimit: body.DailyRequestLimit,
		UpdatedAt:         time.Now(),
	}
	if err := h.store.Set(r.Context(), cfg); err != nil {
		writeQuotaErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	if h.dynamic != nil {
		h.dynamic.InvalidateCache(tenantID)
	}
	writeQuotaJSON(w, http.StatusOK, cfg)
}

func (h *QuotaHTTPHandler) serveDelete(w http.ResponseWriter, r *http.Request, tenantID string) {
	if err := h.store.Delete(r.Context(), tenantID); err != nil {
		writeQuotaErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	if h.dynamic != nil {
		h.dynamic.InvalidateCache(tenantID)
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeQuotaJSON(w http.ResponseWriter, code int, body interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	if body != nil {
		_ = json.NewEncoder(w).Encode(body)
	}
}

func writeQuotaErr(w http.ResponseWriter, code int, msg string) {
	writeQuotaJSON(w, code, map[string]string{"error": msg})
}

// 사용되지 않은 import 제거를 위한 더미 ref
var _ = context.Background