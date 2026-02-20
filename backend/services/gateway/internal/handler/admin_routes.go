package handler

import (
	"net/http"
	"strconv"

	v1 "github.com/manpasik/backend/shared/gen/go/v1"
)

// registerAdminRoutes는 관리자 관련 REST 엔드포인트를 등록합니다.
func (h *RestHandler) registerAdminRoutes(mux *http.ServeMux) {
	// Admin CRUD
	mux.HandleFunc("POST /api/v1/admin/admins", h.handleCreateAdmin)
	mux.HandleFunc("GET /api/v1/admin/admins", h.handleListAdmins)
	mux.HandleFunc("GET /api/v1/admin/admins/{adminId}", h.handleGetAdmin)
	mux.HandleFunc("PUT /api/v1/admin/admins/{adminId}/role", h.handleUpdateAdminRole)
	mux.HandleFunc("POST /api/v1/admin/admins/{adminId}/deactivate", h.handleDeactivateAdmin)
	mux.HandleFunc("GET /api/v1/admin/admins/by-region", h.handleListAdminsByRegion)

	// User management
	mux.HandleFunc("GET /api/v1/admin/users", h.handleAdminListUsers)
	mux.HandleFunc("PUT /api/v1/admin/users/{userId}/role", h.handleAdminChangeRole)
	mux.HandleFunc("POST /api/v1/admin/users/bulk", h.handleAdminBulkAction)

	// System stats & audit
	mux.HandleFunc("GET /api/v1/admin/stats", h.handleGetSystemStats)
	mux.HandleFunc("GET /api/v1/admin/audit-log", h.handleGetAuditLog)
	mux.HandleFunc("GET /api/v1/admin/audit-log/details", h.handleGetAuditLogDetails)

	// System config
	mux.HandleFunc("PUT /api/v1/admin/config", h.handleSetSystemConfig)
	mux.HandleFunc("GET /api/v1/admin/config", h.handleGetSystemConfig)
	mux.HandleFunc("GET /api/v1/admin/configs", h.handleListSystemConfigs)
	mux.HandleFunc("GET /api/v1/admin/configs/{key}", h.handleGetConfigWithMeta)
	mux.HandleFunc("POST /api/v1/admin/configs/validate", h.handleValidateConfigValue)
	mux.HandleFunc("POST /api/v1/admin/configs/bulk", h.handleBulkSetConfigs)

	// Revenue & inventory
	mux.HandleFunc("GET /api/v1/admin/revenue", h.handleGetRevenueStats)
	mux.HandleFunc("GET /api/v1/admin/inventory", h.handleGetInventoryStats)

	// Metrics / hierarchy / compliance
	mux.HandleFunc("GET /api/v1/admin/metrics", h.handleGetAdminMetrics)
	mux.HandleFunc("GET /api/v1/admin/hierarchy", h.handleGetAdminHierarchy)
	mux.HandleFunc("GET /api/v1/admin/compliance", h.handleGetComplianceReport)
}

// ── Admin CRUD ──

func (h *RestHandler) handleCreateAdmin(w http.ResponseWriter, r *http.Request) {
	if h.admin == nil {
		writeError(w, http.StatusServiceUnavailable, "admin service unavailable")
		return
	}
	var body struct {
		Email       string `json:"email"`
		Password    string `json:"password"`
		DisplayName string `json:"display_name"`
		Role        int32  `json:"role"`
		Region      string `json:"region"`
		Branch      string `json:"branch"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.admin.CreateAdmin(r.Context(), &v1.CreateAdminRequest{
		Email:       body.Email,
		Password:    body.Password,
		DisplayName: body.DisplayName,
		Role:        v1.AdminRole(body.Role),
		Region:      body.Region,
		Branch:      body.Branch,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusCreated, resp)
}

func (h *RestHandler) handleGetAdmin(w http.ResponseWriter, r *http.Request) {
	if h.admin == nil {
		writeError(w, http.StatusServiceUnavailable, "admin service unavailable")
		return
	}
	adminId := r.PathValue("adminId")
	resp, err := h.admin.GetAdmin(r.Context(), &v1.GetAdminRequest{
		AdminId: adminId,
	})
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleListAdmins(w http.ResponseWriter, r *http.Request) {
	if h.admin == nil {
		writeError(w, http.StatusServiceUnavailable, "admin service unavailable")
		return
	}
	resp, err := h.admin.ListAdmins(r.Context(), &v1.ListAdminsRequest{
		Limit:  queryInt(r, "limit", 20),
		Offset: queryInt(r, "offset", 0),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleUpdateAdminRole(w http.ResponseWriter, r *http.Request) {
	if h.admin == nil {
		writeError(w, http.StatusServiceUnavailable, "admin service unavailable")
		return
	}
	adminId := r.PathValue("adminId")
	var body struct {
		NewRole int32 `json:"new_role"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.admin.UpdateAdminRole(r.Context(), &v1.UpdateAdminRoleRequest{
		AdminId: adminId,
		NewRole: v1.AdminRole(body.NewRole),
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleDeactivateAdmin(w http.ResponseWriter, r *http.Request) {
	if h.admin == nil {
		writeError(w, http.StatusServiceUnavailable, "admin service unavailable")
		return
	}
	adminId := r.PathValue("adminId")
	resp, err := h.admin.DeactivateAdmin(r.Context(), &v1.DeactivateAdminRequest{
		AdminId: adminId,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleListAdminsByRegion(w http.ResponseWriter, r *http.Request) {
	if h.admin == nil {
		writeError(w, http.StatusServiceUnavailable, "admin service unavailable")
		return
	}
	resp, err := h.admin.ListAdminsByRegion(r.Context(), &v1.ListAdminsByRegionRequest{
		CountryCode: r.URL.Query().Get("country_code"),
		RegionCode:  r.URL.Query().Get("region_code"),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

// ── User management ──

func (h *RestHandler) handleAdminListUsers(w http.ResponseWriter, r *http.Request) {
	if h.admin == nil {
		writeError(w, http.StatusServiceUnavailable, "admin service unavailable")
		return
	}
	resp, err := h.admin.ListUsers(r.Context(), &v1.AdminListUsersRequest{
		Query:  r.URL.Query().Get("query"),
		Limit:  queryInt(r, "limit", 20),
		Offset: queryInt(r, "offset", 0),
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
	var body struct {
		Action  string   `json:"action"`
		UserIDs []string `json:"user_ids"`
		Value   string   `json:"value"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	processed := 0
	for _, uid := range body.UserIDs {
		switch body.Action {
		case "change_role":
			_, err := h.admin.UpdateAdminRole(r.Context(), &v1.UpdateAdminRoleRequest{
				AdminId: uid,
				NewRole: v1.AdminRole(v1.AdminRole_value[body.Value]),
			})
			if err == nil {
				processed++
			}
		default:
			processed++
		}
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success":   true,
		"processed": processed,
		"total":     len(body.UserIDs),
		"action":    body.Action,
	})
}

// ── System stats & audit ──

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

func (h *RestHandler) handleGetAuditLog(w http.ResponseWriter, r *http.Request) {
	if h.admin == nil {
		writeError(w, http.StatusServiceUnavailable, "admin service unavailable")
		return
	}
	resp, err := h.admin.GetAuditLog(r.Context(), &v1.GetAuditLogRequest{
		Limit:  queryInt(r, "limit", 50),
		Offset: queryInt(r, "offset", 0),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleGetAuditLogDetails(w http.ResponseWriter, r *http.Request) {
	if h.admin == nil {
		writeError(w, http.StatusServiceUnavailable, "admin service unavailable")
		return
	}
	resp, err := h.admin.GetAuditLogDetails(r.Context(), &v1.GetAuditLogDetailsRequest{
		Limit:   queryInt(r, "limit", 50),
		Offset:  queryInt(r, "offset", 0),
		AdminId: r.URL.Query().Get("admin_id"),
		Action:  r.URL.Query().Get("action"),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

// ── System config ──

func (h *RestHandler) handleSetSystemConfig(w http.ResponseWriter, r *http.Request) {
	if h.admin == nil {
		writeError(w, http.StatusServiceUnavailable, "admin service unavailable")
		return
	}
	var body struct {
		Key         string `json:"key"`
		Value       string `json:"value"`
		Description string `json:"description"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.admin.SetSystemConfig(r.Context(), &v1.SetSystemConfigRequest{
		Key:         body.Key,
		Value:       body.Value,
		Description: body.Description,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleGetSystemConfig(w http.ResponseWriter, r *http.Request) {
	if h.admin == nil {
		writeError(w, http.StatusServiceUnavailable, "admin service unavailable")
		return
	}
	resp, err := h.admin.GetSystemConfig(r.Context(), &v1.GetSystemConfigRequest{
		Key: r.URL.Query().Get("key"),
	})
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleListSystemConfigs(w http.ResponseWriter, r *http.Request) {
	if h.admin == nil {
		writeError(w, http.StatusServiceUnavailable, "admin service unavailable")
		return
	}
	resp, err := h.admin.ListSystemConfigs(r.Context(), &v1.ListSystemConfigsRequest{
		LanguageCode:   r.URL.Query().Get("language"),
		Category:       r.URL.Query().Get("category"),
		IncludeSecrets: queryBool(r, "include_secrets"),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleGetConfigWithMeta(w http.ResponseWriter, r *http.Request) {
	if h.admin == nil {
		writeError(w, http.StatusServiceUnavailable, "admin service unavailable")
		return
	}
	key := r.PathValue("key")
	resp, err := h.admin.GetConfigWithMeta(r.Context(), &v1.GetConfigWithMetaRequest{
		Key:          key,
		LanguageCode: r.URL.Query().Get("language"),
	})
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleValidateConfigValue(w http.ResponseWriter, r *http.Request) {
	if h.admin == nil {
		writeError(w, http.StatusServiceUnavailable, "admin service unavailable")
		return
	}
	var body struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.admin.ValidateConfigValue(r.Context(), &v1.ValidateConfigValueRequest{
		Key:   body.Key,
		Value: body.Value,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleBulkSetConfigs(w http.ResponseWriter, r *http.Request) {
	if h.admin == nil {
		writeError(w, http.StatusServiceUnavailable, "admin service unavailable")
		return
	}
	var body struct {
		Configs []struct {
			Key         string `json:"key"`
			Value       string `json:"value"`
			Description string `json:"description"`
		} `json:"configs"`
		Reason string `json:"reason"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	configs := make([]*v1.SetSystemConfigRequest, len(body.Configs))
	for i, c := range body.Configs {
		configs[i] = &v1.SetSystemConfigRequest{
			Key:         c.Key,
			Value:       c.Value,
			Description: c.Description,
		}
	}
	resp, err := h.admin.BulkSetConfigs(r.Context(), &v1.BulkSetConfigsRequest{
		Configs: configs,
		Reason:  body.Reason,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

// ── Revenue & Inventory ──

func (h *RestHandler) handleGetRevenueStats(w http.ResponseWriter, r *http.Request) {
	if h.admin == nil {
		writeError(w, http.StatusServiceUnavailable, "admin service unavailable")
		return
	}
	resp, err := h.admin.GetRevenueStats(r.Context(), &v1.GetRevenueStatsRequest{
		Period:    r.URL.Query().Get("period"),
		StartDate: r.URL.Query().Get("start_date"),
		EndDate:   r.URL.Query().Get("end_date"),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleGetInventoryStats(w http.ResponseWriter, r *http.Request) {
	if h.admin == nil {
		writeError(w, http.StatusServiceUnavailable, "admin service unavailable")
		return
	}
	categoryFilter := int32(0)
	if cf := r.URL.Query().Get("category"); cf != "" {
		if v, err := strconv.ParseInt(cf, 10, 32); err == nil {
			categoryFilter = int32(v)
		}
	}
	resp, err := h.admin.GetInventoryStats(r.Context(), &v1.GetInventoryStatsRequest{
		CategoryFilter: categoryFilter,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

// ── Metrics / Hierarchy / Compliance ──

func (h *RestHandler) handleGetAdminMetrics(w http.ResponseWriter, r *http.Request) {
	if h.admin == nil {
		writeError(w, http.StatusServiceUnavailable, "admin service unavailable")
		return
	}
	// GetSystemStats를 메트릭스 엔드포인트로 재사용
	resp, err := h.admin.GetSystemStats(r.Context(), &v1.GetSystemStatsRequest{})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleGetAdminHierarchy(w http.ResponseWriter, r *http.Request) {
	if h.admin == nil {
		writeError(w, http.StatusServiceUnavailable, "admin service unavailable")
		return
	}
	// ListAdminsByRegion을 조직도 조회로 활용
	resp, err := h.admin.ListAdminsByRegion(r.Context(), &v1.ListAdminsByRegionRequest{
		CountryCode: r.URL.Query().Get("country_code"),
		RegionCode:  r.URL.Query().Get("region_code"),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleGetComplianceReport(w http.ResponseWriter, r *http.Request) {
	if h.admin == nil {
		writeError(w, http.StatusServiceUnavailable, "admin service unavailable")
		return
	}
	// GetAuditLogDetails를 준수 보고서로 활용
	resp, err := h.admin.GetAuditLogDetails(r.Context(), &v1.GetAuditLogDetailsRequest{
		Limit:   queryInt(r, "limit", 100),
		Offset:  queryInt(r, "offset", 0),
		AdminId: r.URL.Query().Get("admin_id"),
		Action:  r.URL.Query().Get("action"),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}
