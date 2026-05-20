package handler

import (
	"net/http"

	v1 "github.com/manpasik/backend/shared/gen/go/v1"
)

func (h *RestHandler) registerAssistantRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/assistant/command", h.handleAssistantCommand)
	mux.HandleFunc("GET /api/v1/assistant/sessions", h.handleListAssistantSessions)
	mux.HandleFunc("GET /api/v1/assistant/sessions/{sessionId}", h.handleGetAssistantSession)
	mux.HandleFunc("GET /api/v1/assistant/sessions/{sessionId}/turns", h.handleListAssistantTurns)
}

func (h *RestHandler) handleAssistantCommand(w http.ResponseWriter, r *http.Request) {
	if h.assistant == nil {
		writeError(w, http.StatusServiceUnavailable, "assistant service unavailable")
		return
	}
	var body struct {
		UserID    string `json:"user_id"`
		SessionID string `json:"session_id"`
		Text      string `json:"text"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.assistant.SendCommand(r.Context(), &v1.AssistantCommandRequest{
		UserId: body.UserID, SessionId: body.SessionID, Text: body.Text,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleListAssistantSessions(w http.ResponseWriter, r *http.Request) {
	if h.assistant == nil {
		writeError(w, http.StatusServiceUnavailable, "assistant service unavailable")
		return
	}
	userID := r.URL.Query().Get("user_id")
	resp, err := h.assistant.ListSessions(r.Context(), &v1.ListAssistantSessionsRequest{
		UserId: userID, Limit: queryInt(r, "limit", 20),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleGetAssistantSession(w http.ResponseWriter, r *http.Request) {
	if h.assistant == nil {
		writeError(w, http.StatusServiceUnavailable, "assistant service unavailable")
		return
	}
	sessionID := r.PathValue("sessionId")
	resp, err := h.assistant.GetSession(r.Context(), &v1.GetAssistantSessionRequest{SessionId: sessionID})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleListAssistantTurns(w http.ResponseWriter, r *http.Request) {
	if h.assistant == nil {
		writeError(w, http.StatusServiceUnavailable, "assistant service unavailable")
		return
	}
	sessionID := r.PathValue("sessionId")
	resp, err := h.assistant.ListTurns(r.Context(), &v1.ListAssistantTurnsRequest{
		SessionId: sessionID, Limit: queryInt(r, "limit", 50),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}
