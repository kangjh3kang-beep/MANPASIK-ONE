package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/manpasik/backend/services/audit-service/internal/service"
	"go.uber.org/zap"
)

type HTTPHandler struct {
	svc *service.AuditService
	log *zap.Logger
}

type recordAuditEventRequest struct {
	AdminID      string            `json:"admin_id"`
	Action       string            `json:"action"`
	ResourceType string            `json:"resource_type"`
	ResourceID   string            `json:"resource_id"`
	Description  string            `json:"description"`
	IPAddress    string            `json:"ip_address"`
	UserAgent    string            `json:"user_agent"`
	Metadata     map[string]string `json:"metadata"`
	OccurredAt   string            `json:"occurred_at"`
}

func NewHTTPHandler(svc *service.AuditService, log *zap.Logger) *HTTPHandler {
	return &HTTPHandler{svc: svc, log: log}
}

func (h *HTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /audit/events", h.handleRecordAuditEvent)
}

func (h *HTTPHandler) handleRecordAuditEvent(w http.ResponseWriter, r *http.Request) {
	var body recordAuditEventRequest
	if err := readHTTPJSON(r, &body); err != nil {
		writeHTTPJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if body.Action == "" || body.ResourceType == "" {
		writeHTTPJSON(w, http.StatusBadRequest, map[string]string{"error": "action and resource_type are required"})
		return
	}

	var occurredAt time.Time
	if body.OccurredAt != "" {
		parsed, err := time.Parse(time.RFC3339Nano, body.OccurredAt)
		if err != nil {
			writeHTTPJSON(w, http.StatusBadRequest, map[string]string{"error": "occurred_at must be RFC3339"})
			return
		}
		occurredAt = parsed.UTC()
	}

	entry, err := h.svc.RecordActionWithMetadata(r.Context(), service.RecordActionInput{
		AdminID:      defaultString(body.AdminID, "system:gateway"),
		Action:       body.Action,
		ResourceType: body.ResourceType,
		ResourceID:   body.ResourceID,
		Description:  body.Description,
		IPAddress:    body.IPAddress,
		UserAgent:    body.UserAgent,
		Metadata:     body.Metadata,
		Timestamp:    occurredAt,
	})
	if err != nil {
		h.log.Error("audit event store failed", zap.Error(err))
		writeHTTPJSON(w, http.StatusInternalServerError, map[string]string{"error": "audit event store failed"})
		return
	}

	writeHTTPJSON(w, http.StatusCreated, map[string]interface{}{
		"accepted": true,
		"status":   "persisted",
		"entry_id": entry.ID,
	})
}

func readHTTPJSON(r *http.Request, v interface{}) error {
	reader := io.LimitReader(r.Body, 1<<20)
	body, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.Unmarshal(body, v)
}

func writeHTTPJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func defaultString(value, fallback string) string {
	if value != "" {
		return value
	}
	return fallback
}
