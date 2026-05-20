// Package handler는 analytics-service의 REST 핸들러를 제공합니다.
package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/manpasik/backend/services/analytics-service/internal/service"
)

// RegisterRoutes는 analytics REST 라우트를 등록합니다.
func (h *AnalyticsHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/analytics/events", h.handleTrackEvent)
	mux.HandleFunc("GET /api/v1/analytics/events", h.handleListRecentEvents)
	mux.HandleFunc("GET /api/v1/analytics/users/{userID}", h.handleGetUserAnalytics)
	mux.HandleFunc("GET /api/v1/analytics/daily/{date}", h.handleGetDailyStats)
	mux.HandleFunc("GET /api/v1/analytics/retention/{userID}", h.handleGetUserRetention)
	mux.HandleFunc("POST /api/v1/analytics/funnel", h.handleGetEventFunnel)
}

// ----------- Track Event -----------

type trackEventRequest struct {
	UserID     string            `json:"user_id"`
	EventType  string            `json:"event_type"`
	Properties map[string]string `json:"properties"`
}

func (h *AnalyticsHandler) handleTrackEvent(w http.ResponseWriter, r *http.Request) {
	var req trackEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	eventID, err := h.svc.TrackEvent(r.Context(), req.UserID, req.EventType, req.Properties)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"event_id": eventID})
}

// ----------- List Recent Events -----------

func (h *AnalyticsHandler) handleListRecentEvents(w http.ResponseWriter, r *http.Request) {
	limit := 50
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			limit = n
		}
	}

	events, err := h.svc.ListRecentEvents(r.Context(), limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"events": events, "count": len(events)})
}

// ----------- Get User Analytics -----------

func (h *AnalyticsHandler) handleGetUserAnalytics(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userID")
	if userID == "" {
		writeError(w, http.StatusBadRequest, "user_id is required")
		return
	}

	analytics, err := h.svc.GetUserAnalytics(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, analytics)
}

// ----------- Get Daily Stats -----------

func (h *AnalyticsHandler) handleGetDailyStats(w http.ResponseWriter, r *http.Request) {
	date := r.PathValue("date")
	if date == "" {
		writeError(w, http.StatusBadRequest, "date is required")
		return
	}

	stats, err := h.svc.GetDailyStats(r.Context(), date)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if stats == nil {
		writeJSON(w, http.StatusOK, &service.DailyStats{Date: date})
		return
	}

	writeJSON(w, http.StatusOK, stats)
}

// ----------- Get User Retention -----------

func (h *AnalyticsHandler) handleGetUserRetention(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userID")
	if userID == "" {
		writeError(w, http.StatusBadRequest, "user_id is required")
		return
	}

	days := 7
	if v := r.URL.Query().Get("days"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			days = n
		}
	}

	retention, err := h.svc.GetUserRetention(r.Context(), userID, days)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, retention)
}

// ----------- Get Event Funnel -----------

type funnelRequest struct {
	UserID     string   `json:"user_id"`
	EventTypes []string `json:"event_types"`
}

func (h *AnalyticsHandler) handleGetEventFunnel(w http.ResponseWriter, r *http.Request) {
	var req funnelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	steps, err := h.svc.GetEventFunnel(r.Context(), req.UserID, req.EventTypes)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"steps": steps})
}

// ----------- Helpers -----------

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("writeJSON error: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
