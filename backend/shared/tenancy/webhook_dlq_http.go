package tenancy

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// WebhookDLQHandler 는 DLQ REST API (Phase AO-2).
//
// 엔드포인트:
//
//	GET    {prefix}/webhook/dlq             — 목록 조회 (count + entries)
//	POST   {prefix}/webhook/dlq/{id}/replay — 단일 항목 재발송 (성공 시 제거)
//	DELETE {prefix}/webhook/dlq/{id}        — 단일 항목 폐기 (재발송 없음)
//	DELETE {prefix}/webhook/dlq             — 전체 비움 (감사 로그 권장)
//
// 인증/권한 검사는 호출자 책임.
type WebhookDLQHandler struct {
	dispatcher *WebhookDispatcher
	pathPrefix string
	replayCtx  func(parent context.Context) (context.Context, context.CancelFunc)
}

// NewWebhookDLQHandler 생성. dispatcher nil 이면 503 반환.
func NewWebhookDLQHandler(d *WebhookDispatcher) *WebhookDLQHandler {
	return &WebhookDLQHandler{
		dispatcher: d,
		replayCtx: func(parent context.Context) (context.Context, context.CancelFunc) {
			return context.WithTimeout(parent, 10*time.Second)
		},
	}
}

// SetPathPrefix 는 라우트 prefix 설정.
func (h *WebhookDLQHandler) SetPathPrefix(prefix string) {
	h.pathPrefix = strings.TrimRight(prefix, "/")
}

// RegisterRoutes 는 ServeMux 에 등록.
func (h *WebhookDLQHandler) RegisterRoutes(mux *http.ServeMux) {
	p := h.pathPrefix
	mux.HandleFunc(p+"/webhook/dlq", h.handleCollection)
	mux.HandleFunc(p+"/webhook/dlq/", h.handleItem)
}

func (h *WebhookDLQHandler) handleCollection(w http.ResponseWriter, r *http.Request) {
	if h.dispatcher == nil {
		writeWebhookErr(w, http.StatusServiceUnavailable, "dispatcher 미설정")
		return
	}
	switch r.Method {
	case http.MethodGet:
		entries := h.dispatcher.DLQEntries()
		writeWebhookJSON(w, http.StatusOK, map[string]interface{}{
			"count":   len(entries),
			"entries": entries,
		})
	case http.MethodDelete:
		removed := h.dispatcher.DLQClear()
		writeWebhookJSON(w, http.StatusOK, map[string]interface{}{"cleared": removed})
	default:
		writeWebhookErr(w, http.StatusMethodNotAllowed, "GET 또는 DELETE")
	}
}

func (h *WebhookDLQHandler) handleItem(w http.ResponseWriter, r *http.Request) {
	if h.dispatcher == nil {
		writeWebhookErr(w, http.StatusServiceUnavailable, "dispatcher 미설정")
		return
	}
	rest := strings.TrimPrefix(r.URL.Path, h.pathPrefix+"/webhook/dlq/")
	parts := strings.Split(rest, "/")
	if len(parts) == 0 || parts[0] == "" {
		writeWebhookErr(w, http.StatusBadRequest, "ID 필수")
		return
	}
	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		writeWebhookErr(w, http.StatusBadRequest, "ID 정수 형식 필수")
		return
	}

	// POST {id}/replay
	if len(parts) >= 2 && parts[1] == "replay" {
		if r.Method != http.MethodPost {
			writeWebhookErr(w, http.StatusMethodNotAllowed, "POST only")
			return
		}
		ctx, cancel := h.replayCtx(r.Context())
		defer cancel()
		if rerr := h.dispatcher.DLQReplay(ctx, id); rerr != nil {
			if errors.Is(rerr, ErrDLQNotFound) {
				writeWebhookErr(w, http.StatusNotFound, rerr.Error())
				return
			}
			writeWebhookErr(w, http.StatusBadGateway, rerr.Error())
			return
		}
		writeWebhookJSON(w, http.StatusOK, map[string]interface{}{
			"replayed": id,
			"status":   "sent",
		})
		return
	}

	// DELETE {id}
	if r.Method == http.MethodDelete {
		if !h.dispatcher.DLQDrop(id) {
			writeWebhookErr(w, http.StatusNotFound, "DLQ 항목 없음")
			return
		}
		writeWebhookJSON(w, http.StatusOK, map[string]interface{}{"dropped": id})
		return
	}
	writeWebhookErr(w, http.StatusMethodNotAllowed, "DELETE 또는 POST .../replay")
}

func writeWebhookJSON(w http.ResponseWriter, code int, body interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	if body != nil {
		_ = json.NewEncoder(w).Encode(body)
	}
}

func writeWebhookErr(w http.ResponseWriter, code int, msg string) {
	writeWebhookJSON(w, code, map[string]string{"error": msg})
}
