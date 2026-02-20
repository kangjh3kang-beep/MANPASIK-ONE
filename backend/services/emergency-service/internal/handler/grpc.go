package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/manpasik/backend/services/emergency-service/internal/repository/memory"
	"github.com/manpasik/backend/services/emergency-service/internal/service"
)

// EmergencyHandler handles emergency-related HTTP/gRPC requests.
// In this skeleton, only HTTP JSON handlers are provided (no proto dependency).
type EmergencyHandler struct {
	svc *service.EmergencyService
}

// NewEmergencyHandler creates a new EmergencyHandler.
func NewEmergencyHandler(svc *service.EmergencyService) *EmergencyHandler {
	return &EmergencyHandler{svc: svc}
}

// RegisterRoutes registers all emergency routes on the given mux.
func (h *EmergencyHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/emergency/report", h.handleReportEmergency)
	mux.HandleFunc("/api/v1/emergency/contacts", h.handleGetContacts)
	mux.HandleFunc("/api/v1/emergency/settings", h.handleSettings)
}

// ----------- Report Emergency -----------

type reportRequest struct {
	UserID      string `json:"user_id"`
	Type        string `json:"type"`
	Location    string `json:"location"`
	Description string `json:"description"`
}

type reportResponse struct {
	EmergencyID string `json:"emergency_id"`
}

func (h *EmergencyHandler) handleReportEmergency(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req reportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	id, err := h.svc.ReportEmergency(service.ReportEmergencyInput{
		UserID:      req.UserID,
		Type:        req.Type,
		Location:    req.Location,
		Description: req.Description,
	})
	if err != nil {
		log.Printf("ReportEmergency error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusCreated, reportResponse{EmergencyID: id})
}

// ----------- Get Contacts -----------

func (h *EmergencyHandler) handleGetContacts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("user_id")
	contacts, err := h.svc.GetEmergencyContacts(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusOK, contacts)
}

// ----------- Settings -----------

func (h *EmergencyHandler) handleSettings(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		userID := r.URL.Query().Get("user_id")
		settings, err := h.svc.GetEmergencySettings(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		writeJSON(w, http.StatusOK, settings)

	case http.MethodPut:
		var s memory.EmergencySettings
		if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		if err := h.svc.UpdateEmergencySettings(&s); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		writeJSON(w, http.StatusOK, map[string]bool{"ok": true})

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// ----------- Helpers -----------

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("writeJSON error: %v", err)
	}
}
