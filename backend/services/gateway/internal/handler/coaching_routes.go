package handler

import (
	"net/http"

	v1 "github.com/manpasik/backend/shared/gen/go/v1"
)

// registerCoachingRoutes는 AI 코칭 관련 REST 엔드포인트를 등록합니다.
func (h *RestHandler) registerCoachingRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/coaching/goals", h.handleSetHealthGoal)
	mux.HandleFunc("GET /api/v1/coaching/goals", h.handleGetHealthGoals)
	mux.HandleFunc("POST /api/v1/coaching/generate", h.handleGenerateCoaching)
	mux.HandleFunc("GET /api/v1/coaching/messages", h.handleListCoachingMessages)
	mux.HandleFunc("POST /api/v1/coaching/daily-report", h.handleGenerateDailyReport)
	mux.HandleFunc("GET /api/v1/coaching/weekly-report", h.handleGetWeeklyReport)
	mux.HandleFunc("GET /api/v1/coaching/recommendations", h.handleGetRecommendations)
}

func (h *RestHandler) handleSetHealthGoal(w http.ResponseWriter, r *http.Request) {
	if h.coaching == nil {
		writeError(w, http.StatusServiceUnavailable, "coaching service unavailable")
		return
	}
	var body struct {
		UserID      string  `json:"user_id"`
		Category    int32   `json:"category"`
		MetricName  string  `json:"metric_name"`
		TargetValue float64 `json:"target_value"`
		Unit        string  `json:"unit"`
		Description string  `json:"description"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.coaching.SetHealthGoal(r.Context(), &v1.SetHealthGoalRequest{
		UserId:      body.UserID,
		Category:    v1.GoalCategory(body.Category),
		MetricName:  body.MetricName,
		TargetValue: body.TargetValue,
		Unit:        body.Unit,
		Description: body.Description,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusCreated, resp)
}

func (h *RestHandler) handleGetHealthGoals(w http.ResponseWriter, r *http.Request) {
	if h.coaching == nil {
		writeError(w, http.StatusServiceUnavailable, "coaching service unavailable")
		return
	}
	resp, err := h.coaching.GetHealthGoals(r.Context(), &v1.GetHealthGoalsRequest{
		UserId:       r.URL.Query().Get("user_id"),
		StatusFilter: v1.GoalStatus(queryInt(r, "status", 0)),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleGenerateCoaching(w http.ResponseWriter, r *http.Request) {
	if h.coaching == nil {
		writeError(w, http.StatusServiceUnavailable, "coaching service unavailable")
		return
	}
	var body struct {
		UserID        string `json:"user_id"`
		MeasurementID string `json:"measurement_id"`
		CoachingType  int32  `json:"coaching_type"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.coaching.GenerateCoaching(r.Context(), &v1.GenerateCoachingRequest{
		UserId:        body.UserID,
		MeasurementId: body.MeasurementID,
		CoachingType:  v1.CoachingType(body.CoachingType),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleListCoachingMessages(w http.ResponseWriter, r *http.Request) {
	if h.coaching == nil {
		writeError(w, http.StatusServiceUnavailable, "coaching service unavailable")
		return
	}
	resp, err := h.coaching.ListCoachingMessages(r.Context(), &v1.ListCoachingMessagesRequest{
		UserId:     r.URL.Query().Get("user_id"),
		TypeFilter: v1.CoachingType(queryInt(r, "type", 0)),
		Limit:      queryInt(r, "limit", 20),
		Offset:     queryInt(r, "offset", 0),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleGenerateDailyReport(w http.ResponseWriter, r *http.Request) {
	if h.coaching == nil {
		writeError(w, http.StatusServiceUnavailable, "coaching service unavailable")
		return
	}
	var body struct {
		UserID string `json:"user_id"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.coaching.GenerateDailyReport(r.Context(), &v1.GenerateDailyReportRequest{
		UserId: body.UserID,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleGetWeeklyReport(w http.ResponseWriter, r *http.Request) {
	if h.coaching == nil {
		writeError(w, http.StatusServiceUnavailable, "coaching service unavailable")
		return
	}
	resp, err := h.coaching.GetWeeklyReport(r.Context(), &v1.GetWeeklyReportRequest{
		UserId: r.URL.Query().Get("user_id"),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleGetRecommendations(w http.ResponseWriter, r *http.Request) {
	if h.coaching == nil {
		writeError(w, http.StatusServiceUnavailable, "coaching service unavailable")
		return
	}
	resp, err := h.coaching.GetRecommendations(r.Context(), &v1.GetRecommendationsRequest{
		UserId:     r.URL.Query().Get("user_id"),
		TypeFilter: v1.RecommendationType(queryInt(r, "type", 0)),
		Limit:      queryInt(r, "limit", 10),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}
