package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/manpasik/backend/services/emergency-service/internal/service"
)

// EmergencyHandler handles emergency-related HTTP requests.
type EmergencyHandler struct {
	svc *service.EmergencyService
}

// NewEmergencyHandler creates a new EmergencyHandler.
func NewEmergencyHandler(svc *service.EmergencyService) *EmergencyHandler {
	return &EmergencyHandler{svc: svc}
}

// RegisterRoutes registers all emergency routes on the given mux.
func (h *EmergencyHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/emergency/report", h.handleReportEmergency)
	mux.HandleFunc("POST /api/v1/emergency/resolve", h.handleResolveEmergency)
	mux.HandleFunc("GET /api/v1/emergency/history", h.handleGetHistory)
	mux.HandleFunc("GET /api/v1/emergency/contacts", h.handleGetContacts)
	mux.HandleFunc("POST /api/v1/emergency/contacts", h.handleAddContact)
	mux.HandleFunc("DELETE /api/v1/emergency/contacts", h.handleRemoveContact)
	mux.HandleFunc("GET /api/v1/emergency/settings", h.handleGetSettings)
	mux.HandleFunc("PUT /api/v1/emergency/settings", h.handleUpdateSettings)
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
	var req reportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
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
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, reportResponse{EmergencyID: id})
}

// ----------- Resolve Emergency -----------

type resolveRequest struct {
	EmergencyID string `json:"emergency_id"`
	Resolution  string `json:"resolution"`
}

func (h *EmergencyHandler) handleResolveEmergency(w http.ResponseWriter, r *http.Request) {
	var req resolveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.svc.ResolveEmergency(req.EmergencyID, req.Resolution); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

// ----------- Get Emergency History -----------

func (h *EmergencyHandler) handleGetHistory(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		writeError(w, http.StatusBadRequest, "user_id is required")
		return
	}

	history, err := h.svc.GetEmergencyHistory(userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"emergencies": history, "count": len(history)})
}

// ----------- Get Contacts -----------

func (h *EmergencyHandler) handleGetContacts(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	contacts, err := h.svc.GetEmergencyContacts(userID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, contacts)
}

// ----------- Add Contact -----------

type addContactRequest struct {
	UserID       string `json:"user_id"`
	Name         string `json:"name"`
	Phone        string `json:"phone"`
	Relationship string `json:"relationship"`
}

func (h *EmergencyHandler) handleAddContact(w http.ResponseWriter, r *http.Request) {
	var req addContactRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.svc.AddContact(req.UserID, &service.EmergencyContact{
		UserID:       req.UserID,
		Name:         req.Name,
		Phone:        req.Phone,
		Relationship: req.Relationship,
	}); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]bool{"ok": true})
}

// ----------- Remove Contact -----------

func (h *EmergencyHandler) handleRemoveContact(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	contactID := r.URL.Query().Get("contact_id")

	if err := h.svc.RemoveContact(userID, contactID); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

// ----------- Get Settings -----------

func (h *EmergencyHandler) handleGetSettings(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	settings, err := h.svc.GetEmergencySettings(userID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, settings)
}

// ----------- Update Settings -----------

func (h *EmergencyHandler) handleUpdateSettings(w http.ResponseWriter, r *http.Request) {
	var s service.EmergencySettings
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.svc.UpdateEmergencySettings(&s); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
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
