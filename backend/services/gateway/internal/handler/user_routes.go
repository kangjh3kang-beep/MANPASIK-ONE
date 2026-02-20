package handler

import (
	"net/http"

	v1 "github.com/manpasik/backend/shared/gen/go/v1"
)

// registerUserRoutes는 사용자/AI 추론 관련 REST 엔드포인트를 등록합니다.
func (h *RestHandler) registerUserRoutes(mux *http.ServeMux) {
	// User
	mux.HandleFunc("GET /api/v1/users/{userId}/profile", h.handleGetProfile)
	mux.HandleFunc("PUT /api/v1/users/{userId}/profile", h.handleUpdateProfile)
	mux.HandleFunc("PUT /api/v1/users/{userId}/emergency-settings", h.handleSaveEmergencySettings)

	// Support
	mux.HandleFunc("POST /api/v1/support/inquiries", h.handleCreateSupportInquiry)

	// AI Inference
	mux.HandleFunc("POST /api/v1/ai/analyze", h.handleAnalyzeMeasurement)
	mux.HandleFunc("GET /api/v1/ai/health-score/{userId}", h.handleGetHealthScore)
	mux.HandleFunc("POST /api/v1/ai/predict-trend", h.handlePredictTrend)
	mux.HandleFunc("GET /api/v1/ai/models", h.handleListAiModels)
}

// ── User ──

func (h *RestHandler) handleGetProfile(w http.ResponseWriter, r *http.Request) {
	if h.user == nil {
		writeError(w, http.StatusServiceUnavailable, "user service unavailable")
		return
	}
	userId := r.PathValue("userId")
	resp, err := h.user.GetProfile(r.Context(), &v1.GetProfileRequest{UserId: userId})
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

// ── Support ──

func (h *RestHandler) handleCreateSupportInquiry(w http.ResponseWriter, r *http.Request) {
	var body struct {
		UserID   string `json:"user_id"`
		Category string `json:"category"`
		Subject  string `json:"subject"`
		Content  string `json:"content"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	// 알림 서비스로 문의 전달
	if h.notification != nil {
		h.notification.SendNotification(r.Context(), &v1.SendNotificationRequest{
			UserId:   body.UserID,
			Type:     v1.NotificationType_NOTIFICATION_TYPE_SYSTEM,
			Title:    "문의 접수: " + body.Subject,
			Body:     body.Content,
			Priority: v1.NotificationPriority_NOTIFICATION_PRIORITY_NORMAL,
			Data:     map[string]string{"category": body.Category},
		})
	}
	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"success":    true,
		"inquiry_id": "inq_" + body.UserID,
		"message":    "문의가 접수되었습니다",
	})
}

func (h *RestHandler) handleUpdateProfile(w http.ResponseWriter, r *http.Request) {
	if h.user == nil {
		writeError(w, http.StatusServiceUnavailable, "user service unavailable")
		return
	}
	userId := r.PathValue("userId")
	var body struct {
		DisplayName string `json:"display_name"`
		AvatarURL   string `json:"avatar_url"`
		Language    string `json:"language"`
		Timezone    string `json:"timezone"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.user.UpdateProfile(r.Context(), &v1.UpdateProfileRequest{
		UserId:      userId,
		DisplayName: body.DisplayName,
		AvatarUrl:   body.AvatarURL,
		Language:    body.Language,
		Timezone:    body.Timezone,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleSaveEmergencySettings(w http.ResponseWriter, r *http.Request) {
	if h.user == nil {
		writeError(w, http.StatusServiceUnavailable, "user service unavailable")
		return
	}
	userId := r.PathValue("userId")
	var body struct {
		Contacts      []map[string]string `json:"contacts"`
		AutoCall119   bool                `json:"auto_call_119"`
		SafetyMode    string              `json:"safety_mode"`
		HighThreshold float64             `json:"high_threshold"`
		LowThreshold  float64             `json:"low_threshold"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.user.UpdateProfile(r.Context(), &v1.UpdateProfileRequest{
		UserId: userId,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

// ── AI Inference ──

func (h *RestHandler) handleAnalyzeMeasurement(w http.ResponseWriter, r *http.Request) {
	if h.aiInference == nil {
		writeError(w, http.StatusServiceUnavailable, "ai-inference service unavailable")
		return
	}
	var body struct {
		UserID        string  `json:"user_id"`
		MeasurementID string  `json:"measurement_id"`
		Models        []int32 `json:"models"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	var models []v1.AiModelType
	for _, m := range body.Models {
		models = append(models, v1.AiModelType(m))
	}
	resp, err := h.aiInference.AnalyzeMeasurement(r.Context(), &v1.AnalyzeMeasurementRequest{
		UserId:        body.UserID,
		MeasurementId: body.MeasurementID,
		Models:        models,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleGetHealthScore(w http.ResponseWriter, r *http.Request) {
	if h.aiInference == nil {
		writeError(w, http.StatusServiceUnavailable, "ai-inference service unavailable")
		return
	}
	userId := r.PathValue("userId")
	resp, err := h.aiInference.GetHealthScore(r.Context(), &v1.GetHealthScoreRequest{
		UserId: userId,
		Days:   queryInt(r, "days", 30),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handlePredictTrend(w http.ResponseWriter, r *http.Request) {
	if h.aiInference == nil {
		writeError(w, http.StatusServiceUnavailable, "ai-inference service unavailable")
		return
	}
	var body struct {
		UserID         string `json:"user_id"`
		MetricName     string `json:"metric_name"`
		HistoryDays    int32  `json:"history_days"`
		PredictionDays int32  `json:"prediction_days"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.aiInference.PredictTrend(r.Context(), &v1.PredictTrendRequest{
		UserId:         body.UserID,
		MetricName:     body.MetricName,
		HistoryDays:    body.HistoryDays,
		PredictionDays: body.PredictionDays,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleListAiModels(w http.ResponseWriter, r *http.Request) {
	if h.aiInference == nil {
		writeError(w, http.StatusServiceUnavailable, "ai-inference service unavailable")
		return
	}
	resp, err := h.aiInference.ListModels(r.Context(), &v1.ListModelsRequest{})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}
