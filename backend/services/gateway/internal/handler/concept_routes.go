package handler

import (
	"net/http"

	v1 "github.com/manpasik/backend/shared/gen/go/v1"
)

func (h *RestHandler) registerConceptRoutes(mux *http.ServeMux) {
	// Concept
	mux.HandleFunc("GET /api/v1/concepts", h.handleListConcepts)
	mux.HandleFunc("GET /api/v1/concepts/{conceptId}", h.handleGetConcept)
	mux.HandleFunc("POST /api/v1/concepts", h.handleCreateConcept)
	mux.HandleFunc("POST /api/v1/concepts/{conceptId}/devices", h.handleAssignDevice)
	mux.HandleFunc("GET /api/v1/concepts/{conceptId}/stats", h.handleGetConceptStats)
	mux.HandleFunc("GET /api/v1/concepts/{conceptId}/dashboard", h.handleGetConceptDashboard)
	// Organization
	mux.HandleFunc("GET /api/v1/organizations", h.handleListOrganizations)
	mux.HandleFunc("GET /api/v1/organizations/{orgId}", h.handleGetOrganization)
	mux.HandleFunc("POST /api/v1/organizations", h.handleCreateOrganization)
	mux.HandleFunc("POST /api/v1/organizations/{orgId}/members", h.handleAddMember)
	mux.HandleFunc("DELETE /api/v1/organizations/{orgId}/members/{userId}", h.handleRemoveMember)
}

func (h *RestHandler) handleListConcepts(w http.ResponseWriter, r *http.Request) {
	if h.concept == nil {
		writeError(w, http.StatusServiceUnavailable, "concept service unavailable")
		return
	}
	resp, err := h.concept.ListConcepts(r.Context(), &v1.ListConceptsRequest{})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleGetConcept(w http.ResponseWriter, r *http.Request) {
	if h.concept == nil {
		writeError(w, http.StatusServiceUnavailable, "concept service unavailable")
		return
	}
	resp, err := h.concept.GetConcept(r.Context(), &v1.GetConceptRequest{ConceptId: r.PathValue("conceptId")})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleCreateConcept(w http.ResponseWriter, r *http.Request) {
	if h.concept == nil {
		writeError(w, http.StatusServiceUnavailable, "concept service unavailable")
		return
	}
	var body struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Category    string `json:"category"`
		OwnerID     string `json:"owner_id"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.concept.CreateConcept(r.Context(), &v1.CreateConceptRequest{
		Name: body.Name, Description: body.Description, Category: body.Category, OwnerId: body.OwnerID,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusCreated, resp)
}

func (h *RestHandler) handleAssignDevice(w http.ResponseWriter, r *http.Request) {
	if h.concept == nil {
		writeError(w, http.StatusServiceUnavailable, "concept service unavailable")
		return
	}
	var body struct {
		DeviceID string `json:"device_id"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.concept.AssignDeviceToConcept(r.Context(), &v1.AssignDeviceToConceptRequest{
		ConceptId: r.PathValue("conceptId"), DeviceId: body.DeviceID,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleGetConceptStats(w http.ResponseWriter, r *http.Request) {
	if h.concept == nil {
		writeError(w, http.StatusServiceUnavailable, "concept service unavailable")
		return
	}
	resp, err := h.concept.GetConceptStats(r.Context(), &v1.GetConceptStatsRequest{ConceptId: r.PathValue("conceptId")})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleGetConceptDashboard(w http.ResponseWriter, r *http.Request) {
	if h.concept == nil {
		writeError(w, http.StatusServiceUnavailable, "concept service unavailable")
		return
	}
	resp, err := h.concept.GetConceptDashboard(r.Context(), &v1.GetConceptDashboardRequest{ConceptId: r.PathValue("conceptId")})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleListOrganizations(w http.ResponseWriter, r *http.Request) {
	if h.organization == nil {
		writeError(w, http.StatusServiceUnavailable, "organization service unavailable")
		return
	}
	resp, err := h.organization.ListOrganizations(r.Context(), &v1.ListOrganizationsRequest{})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleGetOrganization(w http.ResponseWriter, r *http.Request) {
	if h.organization == nil {
		writeError(w, http.StatusServiceUnavailable, "organization service unavailable")
		return
	}
	resp, err := h.organization.GetOrganization(r.Context(), &v1.GetOrganizationRequest{OrganizationId: r.PathValue("orgId")})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleCreateOrganization(w http.ResponseWriter, r *http.Request) {
	if h.organization == nil {
		writeError(w, http.StatusServiceUnavailable, "organization service unavailable")
		return
	}
	var body struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.organization.CreateOrganization(r.Context(), &v1.CreateOrganizationRequest{
		Name: body.Name, Description: body.Description,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusCreated, resp)
}

func (h *RestHandler) handleAddMember(w http.ResponseWriter, r *http.Request) {
	if h.organization == nil {
		writeError(w, http.StatusServiceUnavailable, "organization service unavailable")
		return
	}
	var body struct {
		UserID string `json:"user_id"`
		Role   string `json:"role"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.organization.AddMember(r.Context(), &v1.AddOrgMemberRequest{
		OrganizationId: r.PathValue("orgId"), UserId: body.UserID, Role: body.Role,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleRemoveMember(w http.ResponseWriter, r *http.Request) {
	if h.organization == nil {
		writeError(w, http.StatusServiceUnavailable, "organization service unavailable")
		return
	}
	resp, err := h.organization.RemoveMember(r.Context(), &v1.RemoveOrgMemberRequest{
		OrganizationId: r.PathValue("orgId"), UserId: r.PathValue("userId"),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}
