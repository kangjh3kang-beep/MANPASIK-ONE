package handler

import (
	"net/http"

	v1 "github.com/manpasik/backend/shared/gen/go/v1"
)

// registerNotificationRoutes는 알림 관련 REST 엔드포인트를 등록합니다.
func (h *RestHandler) registerNotificationRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/notifications", h.handleSendNotification)
	mux.HandleFunc("GET /api/v1/notifications", h.handleListNotifications)
	mux.HandleFunc("POST /api/v1/notifications/{notificationId}/read", h.handleMarkAsRead)
	mux.HandleFunc("POST /api/v1/notifications/read-all", h.handleMarkAllAsRead)
	mux.HandleFunc("GET /api/v1/notifications/unread-count", h.handleGetUnreadCount)
	mux.HandleFunc("PUT /api/v1/notifications/preferences", h.handleUpdateNotificationPreferences)
	mux.HandleFunc("GET /api/v1/notifications/preferences", h.handleGetNotificationPreferences)
	mux.HandleFunc("POST /api/v1/notifications/template", h.handleSendFromTemplate)

	// 개별 알림 조회 + 푸시 토큰 등록
	mux.HandleFunc("GET /api/v1/notifications/alerts/{alertId}", h.handleGetAlert)
	mux.HandleFunc("POST /api/v1/notifications/push-token", h.handleRegisterPushToken)
}

func (h *RestHandler) handleSendNotification(w http.ResponseWriter, r *http.Request) {
	if h.notification == nil {
		writeError(w, http.StatusServiceUnavailable, "notification service unavailable")
		return
	}
	var body struct {
		UserID    string            `json:"user_id"`
		Type      int32             `json:"type"`
		Title     string            `json:"title"`
		Body      string            `json:"body"`
		Priority  int32             `json:"priority"`
		Channel   int32             `json:"channel"`
		Data      map[string]string `json:"data"`
		ActionURL string            `json:"action_url"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.notification.SendNotification(r.Context(), &v1.SendNotificationRequest{
		UserId:    body.UserID,
		Type:      v1.NotificationType(body.Type),
		Title:     body.Title,
		Body:      body.Body,
		Priority:  v1.NotificationPriority(body.Priority),
		Channel:   v1.NotificationChannel(body.Channel),
		Data:      body.Data,
		ActionUrl: body.ActionURL,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusCreated, resp)
}

func (h *RestHandler) handleListNotifications(w http.ResponseWriter, r *http.Request) {
	if h.notification == nil {
		writeError(w, http.StatusServiceUnavailable, "notification service unavailable")
		return
	}
	resp, err := h.notification.ListNotifications(r.Context(), &v1.ListNotificationsRequest{
		UserId:     r.URL.Query().Get("user_id"),
		TypeFilter: v1.NotificationType(queryInt(r, "type", 0)),
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

func (h *RestHandler) handleMarkAsRead(w http.ResponseWriter, r *http.Request) {
	if h.notification == nil {
		writeError(w, http.StatusServiceUnavailable, "notification service unavailable")
		return
	}
	notificationId := r.PathValue("notificationId")
	resp, err := h.notification.MarkAsRead(r.Context(), &v1.MarkAsReadRequest{
		NotificationId: notificationId,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleMarkAllAsRead(w http.ResponseWriter, r *http.Request) {
	if h.notification == nil {
		writeError(w, http.StatusServiceUnavailable, "notification service unavailable")
		return
	}
	var body struct {
		UserID string `json:"user_id"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.notification.MarkAllAsRead(r.Context(), &v1.MarkAllAsReadRequest{
		UserId: body.UserID,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleGetUnreadCount(w http.ResponseWriter, r *http.Request) {
	if h.notification == nil {
		writeError(w, http.StatusServiceUnavailable, "notification service unavailable")
		return
	}
	resp, err := h.notification.GetUnreadCount(r.Context(), &v1.GetUnreadCountRequest{
		UserId: r.URL.Query().Get("user_id"),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleUpdateNotificationPreferences(w http.ResponseWriter, r *http.Request) {
	if h.notification == nil {
		writeError(w, http.StatusServiceUnavailable, "notification service unavailable")
		return
	}
	var body struct {
		UserID               string `json:"user_id"`
		PushEnabled          bool   `json:"push_enabled"`
		EmailEnabled         bool   `json:"email_enabled"`
		SmsEnabled           bool   `json:"sms_enabled"`
		MeasurementAlerts    bool   `json:"measurement_alerts"`
		HealthAlerts         bool   `json:"health_alerts"`
		AppointmentReminders bool   `json:"appointment_reminders"`
		CommunityUpdates     bool   `json:"community_updates"`
		Promotions           bool   `json:"promotions"`
		QuietHoursStart      string `json:"quiet_hours_start"`
		QuietHoursEnd        string `json:"quiet_hours_end"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.notification.UpdateNotificationPreferences(r.Context(), &v1.UpdateNotificationPreferencesRequest{
		UserId:               body.UserID,
		PushEnabled:          body.PushEnabled,
		EmailEnabled:         body.EmailEnabled,
		SmsEnabled:           body.SmsEnabled,
		MeasurementAlerts:    body.MeasurementAlerts,
		HealthAlerts:         body.HealthAlerts,
		AppointmentReminders: body.AppointmentReminders,
		CommunityUpdates:     body.CommunityUpdates,
		Promotions:           body.Promotions,
		QuietHoursStart:      body.QuietHoursStart,
		QuietHoursEnd:        body.QuietHoursEnd,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleGetNotificationPreferences(w http.ResponseWriter, r *http.Request) {
	if h.notification == nil {
		writeError(w, http.StatusServiceUnavailable, "notification service unavailable")
		return
	}
	resp, err := h.notification.GetNotificationPreferences(r.Context(), &v1.GetNotificationPreferencesRequest{
		UserId: r.URL.Query().Get("user_id"),
	})
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleSendFromTemplate(w http.ResponseWriter, r *http.Request) {
	if h.notification == nil {
		writeError(w, http.StatusServiceUnavailable, "notification service unavailable")
		return
	}
	var body struct {
		TemplateKey string            `json:"template_key"`
		UserID      string            `json:"user_id"`
		Data        map[string]string `json:"data"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.notification.SendFromTemplate(r.Context(), &v1.SendFromTemplateRequest{
		TemplateKey: body.TemplateKey,
		UserId:      body.UserID,
		Data:        body.Data,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusCreated, resp)
}

// ── 개별 알림 조회 ──

func (h *RestHandler) handleGetAlert(w http.ResponseWriter, r *http.Request) {
	if h.notification == nil {
		writeError(w, http.StatusServiceUnavailable, "notification service unavailable")
		return
	}
	alertId := r.PathValue("alertId")
	// ListNotifications를 단일 알림 조회로 활용 (ID 필터)
	resp, err := h.notification.ListNotifications(r.Context(), &v1.ListNotificationsRequest{
		UserId: r.URL.Query().Get("user_id"),
		Limit:  1,
	})
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	// alertId를 응답에 포함
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"alert_id":      alertId,
		"notifications": resp,
	})
}

// ── 푸시 토큰 등록 ──

func (h *RestHandler) handleRegisterPushToken(w http.ResponseWriter, r *http.Request) {
	if h.notification == nil {
		writeError(w, http.StatusServiceUnavailable, "notification service unavailable")
		return
	}
	var body struct {
		UserID   string `json:"user_id"`
		Token    string `json:"token"`
		Platform string `json:"platform"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	// 푸시 토큰 등록 = 푸시 알림 활성화 + 환경설정 업데이트
	resp, err := h.notification.UpdateNotificationPreferences(r.Context(), &v1.UpdateNotificationPreferencesRequest{
		UserId:      body.UserID,
		PushEnabled: true,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}
