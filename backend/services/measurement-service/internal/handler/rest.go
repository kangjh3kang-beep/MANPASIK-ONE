// Package handler는 measurement-service의 REST 핸들러를 제공합니다.
package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/manpasik/backend/services/measurement-service/internal/service"
)

// RESTHandler는 측정 REST API 핸들러입니다.
type RESTHandler struct {
	svc *service.MeasurementService
}

// NewRESTHandler는 새 RESTHandler를 생성합니다.
func NewRESTHandler(svc *service.MeasurementService) *RESTHandler {
	return &RESTHandler{svc: svc}
}

// RegisterRoutes는 측정 REST 라우트를 등록합니다.
func (h *RESTHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/measurements/history", h.handleGetHistory)
	mux.HandleFunc("POST /api/v1/measurements/search/similar", h.handleSearchSimilar)
}

// ----------- Get Measurement History -----------

func (h *RESTHandler) handleGetHistory(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		restWriteError(w, http.StatusBadRequest, "user_id is required")
		return
	}

	limit := 20
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			limit = n
		}
	}

	offset := 0
	if v := r.URL.Query().Get("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			offset = n
		}
	}

	// 기본 시간 범위: 최근 30일
	end := time.Now().UTC()
	start := end.AddDate(0, 0, -30)

	history, total, err := h.svc.GetHistory(r.Context(), userID, start, end, limit, offset)
	if err != nil {
		restWriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	restWriteJSON(w, http.StatusOK, map[string]interface{}{
		"measurements": history,
		"total_count":  total,
	})
}

// ----------- Search Similar Fingerprints -----------

type searchSimilarRequest struct {
	Vector []float32 `json:"vector"`
	TopK   int       `json:"top_k"`
}

func (h *RESTHandler) handleSearchSimilar(w http.ResponseWriter, r *http.Request) {
	var req searchSimilarRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		restWriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if len(req.Vector) == 0 {
		restWriteError(w, http.StatusBadRequest, "vector is required")
		return
	}

	results, err := h.svc.SearchSimilarFingerprints(r.Context(), req.Vector, req.TopK)
	if err != nil {
		restWriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	restWriteJSON(w, http.StatusOK, map[string]interface{}{
		"results": results,
		"count":   len(results),
	})
}

// ----------- Helpers -----------

func restWriteJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("restWriteJSON error: %v", err)
	}
}

func restWriteError(w http.ResponseWriter, status int, msg string) {
	restWriteJSON(w, status, map[string]string{"error": msg})
}
