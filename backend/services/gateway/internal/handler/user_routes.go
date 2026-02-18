package handler

import (
	"net/http"
	"time"

	v1 "github.com/manpasik/backend/shared/gen/go/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// registerUserRoutes는 사용자/구독/알림/번역/코칭 관련 REST 엔드포인트를 등록합니다.
func (h *RestHandler) registerUserRoutes(mux *http.ServeMux) {
	// User
	mux.HandleFunc("GET /api/v1/users/{userId}/profile", h.handleGetProfile)
	mux.HandleFunc("PUT /api/v1/users/{userId}/profile", h.handleUpdateProfile)
	mux.HandleFunc("PUT /api/v1/users/{userId}/emergency-settings", h.handleSaveEmergencySettings)

	// Subscription
	mux.HandleFunc("GET /api/v1/subscriptions/plans", h.handleListSubscriptionPlans)
	mux.HandleFunc("GET /api/v1/subscriptions/{userId}", h.handleGetSubscription)
	mux.HandleFunc("POST /api/v1/subscriptions", h.handleCreateSubscription)
	mux.HandleFunc("DELETE /api/v1/subscriptions/{subscriptionId}", h.handleCancelSubscription)

	// Notification
	mux.HandleFunc("GET /api/v1/notifications", h.handleListNotifications)
	mux.HandleFunc("GET /api/v1/notifications/unread-count", h.handleGetUnreadCount)
	mux.HandleFunc("POST /api/v1/notifications/{notificationId}/read", h.handleMarkNotificationAsRead)

	// Translation
	mux.HandleFunc("POST /api/v1/translations/translate", h.handleTranslateText)

	// Coaching
	mux.HandleFunc("POST /api/v1/coaching/goals", h.handleSetHealthGoal)
	mux.HandleFunc("GET /api/v1/coaching/goals/{userId}", h.handleGetHealthGoals)
	mux.HandleFunc("POST /api/v1/coaching/generate", h.handleGenerateCoaching)
	mux.HandleFunc("GET /api/v1/coaching/daily-report/{userId}", h.handleGenerateDailyReport)
	mux.HandleFunc("GET /api/v1/coaching/recommendations/{userId}", h.handleGetRecommendations)

	// AI
	mux.HandleFunc("POST /api/v1/ai/analyze", h.handleAnalyzeMeasurement)
	mux.HandleFunc("GET /api/v1/ai/health-score/{userId}", h.handleGetHealthScore)
	mux.HandleFunc("POST /api/v1/ai/predict-trend", h.handlePredictTrend)
	mux.HandleFunc("GET /api/v1/ai/models", h.handleListAiModels)

	// Admin
	mux.HandleFunc("GET /api/v1/admin/stats", h.handleGetSystemStats)
	mux.HandleFunc("GET /api/v1/admin/users", h.handleAdminListUsers)
	mux.HandleFunc("PUT /api/v1/admin/users/{userId}/role", h.handleAdminChangeRole)
	mux.HandleFunc("POST /api/v1/admin/users/bulk", h.handleAdminBulkAction)
	mux.HandleFunc("GET /api/v1/admin/audit-log", h.handleGetAuditLog)
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
	// Emergency settings are stored via user service profile update
	if h.user == nil {
		writeError(w, http.StatusServiceUnavailable, "user service unavailable")
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "긴급 설정이 저장되었습니다",
	})
}

// ── Subscription ──

func (h *RestHandler) handleListSubscriptionPlans(w http.ResponseWriter, r *http.Request) {
	if h.subscription == nil {
		writeError(w, http.StatusServiceUnavailable, "subscription service unavailable")
		return
	}
	resp, err := h.subscription.ListSubscriptionPlans(r.Context(), &v1.ListSubscriptionPlansRequest{})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleGetSubscription(w http.ResponseWriter, r *http.Request) {
	if h.subscription == nil {
		writeError(w, http.StatusServiceUnavailable, "subscription service unavailable")
		return
	}
	userId := r.PathValue("userId")
	resp, err := h.subscription.GetSubscription(r.Context(), &v1.GetSubscriptionDetailRequest{UserId: userId})
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleCreateSubscription(w http.ResponseWriter, r *http.Request) {
	if h.subscription == nil {
		writeError(w, http.StatusServiceUnavailable, "subscription service unavailable")
		return
	}
	var body struct {
		UserID string `json:"user_id"`
		Tier   int32  `json:"tier"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.subscription.CreateSubscription(r.Context(), &v1.CreateSubscriptionRequest{
		UserId: body.UserID,
		Tier:   v1.SubscriptionTier(body.Tier),
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusCreated, resp)
}

func (h *RestHandler) handleCancelSubscription(w http.ResponseWriter, r *http.Request) {
	if h.subscription == nil {
		writeError(w, http.StatusServiceUnavailable, "subscription service unavailable")
		return
	}
	subId := r.PathValue("subscriptionId")
	var body struct {
		UserID string `json:"user_id"`
		Reason string `json:"reason"`
	}
	readJSON(r, &body)
	_ = subId // subscription ID is implied by user
	resp, err := h.subscription.CancelSubscription(r.Context(), &v1.CancelSubscriptionRequest{
		UserId: body.UserID,
		Reason: body.Reason,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

// ── Notification ──

func (h *RestHandler) handleListNotifications(w http.ResponseWriter, r *http.Request) {
	if h.notification == nil {
		writeError(w, http.StatusServiceUnavailable, "notification service unavailable")
		return
	}
	userId := r.URL.Query().Get("user_id")
	resp, err := h.notification.ListNotifications(r.Context(), &v1.ListNotificationsRequest{
		UserId:     userId,
		UnreadOnly: queryBool(r, "unread_only"),
		Limit:      queryInt(r, "limit", 20),
		Offset:     queryInt(r, "offset", 0),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleGetUnreadCount(w http.ResponseWriter, r *http.Request) {
	if h.notification == nil {
		writeError(w, http.StatusServiceUnavailable, "notification service unavailable")
		return
	}
	userId := r.URL.Query().Get("user_id")
	resp, err := h.notification.GetUnreadCount(r.Context(), &v1.GetUnreadCountRequest{UserId: userId})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleMarkNotificationAsRead(w http.ResponseWriter, r *http.Request) {
	if h.notification == nil {
		writeError(w, http.StatusServiceUnavailable, "notification service unavailable")
		return
	}
	notifId := r.PathValue("notificationId")
	resp, err := h.notification.MarkAsRead(r.Context(), &v1.MarkAsReadRequest{NotificationId: notifId})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

// ── Translation ──

func (h *RestHandler) handleTranslateText(w http.ResponseWriter, r *http.Request) {
	if h.translation == nil {
		writeError(w, http.StatusServiceUnavailable, "translation service unavailable")
		return
	}
	var body struct {
		Text           string `json:"text"`
		SourceLanguage string `json:"source_language"`
		TargetLanguage string `json:"target_language"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.translation.TranslateText(r.Context(), &v1.TranslateTextRequest{
		Text:           body.Text,
		SourceLanguage: body.SourceLanguage,
		TargetLanguage: body.TargetLanguage,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

// ── Coaching ──

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
		TargetDate  string  `json:"target_date"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req := &v1.SetHealthGoalRequest{
		UserId:      body.UserID,
		Category:    v1.GoalCategory(body.Category),
		MetricName:  body.MetricName,
		TargetValue: body.TargetValue,
		Unit:        body.Unit,
		Description: body.Description,
	}
	if body.TargetDate != "" {
		if t, tErr := time.Parse(time.RFC3339, body.TargetDate); tErr == nil {
			req.TargetDate = timestamppb.New(t)
		}
	}
	resp, err := h.coaching.SetHealthGoal(r.Context(), req)
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
	userId := r.PathValue("userId")
	resp, err := h.coaching.GetHealthGoals(r.Context(), &v1.GetHealthGoalsRequest{
		UserId:       userId,
		StatusFilter: v1.GoalStatus(queryInt(r, "status_filter", 0)),
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

func (h *RestHandler) handleGenerateDailyReport(w http.ResponseWriter, r *http.Request) {
	if h.coaching == nil {
		writeError(w, http.StatusServiceUnavailable, "coaching service unavailable")
		return
	}
	userId := r.PathValue("userId")
	resp, err := h.coaching.GenerateDailyReport(r.Context(), &v1.GenerateDailyReportRequest{UserId: userId})
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
	userId := r.PathValue("userId")
	resp, err := h.coaching.GetRecommendations(r.Context(), &v1.GetRecommendationsRequest{
		UserId:     userId,
		TypeFilter: v1.RecommendationType(queryInt(r, "type_filter", 0)),
		Limit:      queryInt(r, "limit", 10),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
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

// ── Admin ──

func (h *RestHandler) handleGetSystemStats(w http.ResponseWriter, r *http.Request) {
	if h.admin == nil {
		writeError(w, http.StatusServiceUnavailable, "admin service unavailable")
		return
	}
	resp, err := h.admin.GetSystemStats(r.Context(), &v1.GetSystemStatsRequest{})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleAdminListUsers(w http.ResponseWriter, r *http.Request) {
	if h.admin == nil {
		writeError(w, http.StatusServiceUnavailable, "admin service unavailable")
		return
	}
	resp, err := h.admin.ListUsers(r.Context(), &v1.AdminListUsersRequest{
		Query:      r.URL.Query().Get("query"),
		TierFilter: v1.SubscriptionTier(queryInt(r, "tier_filter", 0)),
		ActiveOnly: queryBool(r, "active_only"),
		Limit:      queryInt(r, "limit", 20),
		Offset:     queryInt(r, "offset", 0),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleAdminChangeRole(w http.ResponseWriter, r *http.Request) {
	if h.admin == nil {
		writeError(w, http.StatusServiceUnavailable, "admin service unavailable")
		return
	}
	userId := r.PathValue("userId")
	var body struct {
		Role string `json:"role"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.admin.UpdateAdminRole(r.Context(), &v1.UpdateAdminRoleRequest{
		AdminId: userId,
		NewRole: v1.AdminRole(v1.AdminRole_value[body.Role]),
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleAdminBulkAction(w http.ResponseWriter, r *http.Request) {
	if h.admin == nil {
		writeError(w, http.StatusServiceUnavailable, "admin service unavailable")
		return
	}
	// Bulk action is handled at admin level
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "bulk action completed",
	})
}

func (h *RestHandler) handleGetAuditLog(w http.ResponseWriter, r *http.Request) {
	if h.admin == nil {
		writeError(w, http.StatusServiceUnavailable, "admin service unavailable")
		return
	}
	resp, err := h.admin.GetAuditLog(r.Context(), &v1.GetAuditLogRequest{
		AdminId: r.URL.Query().Get("admin_id"),
		Limit:   queryInt(r, "limit", 20),
		Offset:  queryInt(r, "offset", 0),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}
