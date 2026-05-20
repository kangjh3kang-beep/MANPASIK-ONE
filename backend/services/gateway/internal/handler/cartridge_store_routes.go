package handler

import (
	"net/http"

	v1 "github.com/manpasik/backend/shared/gen/go/v1"
)

func (h *RestHandler) registerCartridgeStoreRoutes(mux *http.ServeMux) {
	// Developer
	mux.HandleFunc("POST /api/v1/developers/register", h.handleRegisterDeveloper)
	mux.HandleFunc("GET /api/v1/developers/profile", h.handleGetDeveloperProfile)
	mux.HandleFunc("POST /api/v1/developers/api-keys", h.handleCreateApiKey)
	mux.HandleFunc("POST /api/v1/developers/submit-cartridge", h.handleSubmitCartridge)
	mux.HandleFunc("GET /api/v1/developers/submissions", h.handleListSubmissions)
	// Cartridge Type Definition (Phase 5)
	mux.HandleFunc("POST /api/v1/developers/cartridge-types", h.handleRegisterCartridgeType)
	mux.HandleFunc("GET /api/v1/developers/cartridge-types/{definitionId}", h.handleGetCartridgeDefinition)
	// Store
	mux.HandleFunc("GET /api/v1/store/items", h.handleListStoreItems)
	mux.HandleFunc("GET /api/v1/store/search", h.handleSearchCartridges)
	mux.HandleFunc("GET /api/v1/store/items/{itemId}", h.handleGetStoreItem)
	mux.HandleFunc("POST /api/v1/store/purchase", h.handlePurchaseCartridge)
	mux.HandleFunc("GET /api/v1/store/purchases", h.handleGetPurchaseHistory)
	// Review
	mux.HandleFunc("POST /api/v1/store/review/submit", h.handleSubmitForReview)
	mux.HandleFunc("GET /api/v1/store/review/{submissionId}", h.handleGetReviewStatus)
}

func (h *RestHandler) handleRegisterDeveloper(w http.ResponseWriter, r *http.Request) {
	if h.developer == nil {
		writeError(w, http.StatusServiceUnavailable, "developer service unavailable")
		return
	}
	var body struct {
		UserID      string `json:"user_id"`
		CompanyName string `json:"company_name"`
		Tier        string `json:"tier"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.developer.RegisterDeveloper(r.Context(), &v1.RegisterDeveloperRequest{
		UserId: body.UserID, CompanyName: body.CompanyName, Tier: body.Tier,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusCreated, resp)
}

func (h *RestHandler) handleGetDeveloperProfile(w http.ResponseWriter, r *http.Request) {
	if h.developer == nil {
		writeError(w, http.StatusServiceUnavailable, "developer service unavailable")
		return
	}
	userID := r.URL.Query().Get("user_id")
	resp, err := h.developer.GetDeveloperProfile(r.Context(), &v1.GetDeveloperProfileRequest{UserId: userID})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleCreateApiKey(w http.ResponseWriter, r *http.Request) {
	if h.developer == nil {
		writeError(w, http.StatusServiceUnavailable, "developer service unavailable")
		return
	}
	var body struct {
		DeveloperID string `json:"developer_id"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.developer.CreateApiKey(r.Context(), &v1.CreateApiKeyRequest{DeveloperId: body.DeveloperID})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusCreated, resp)
}

func (h *RestHandler) handleSubmitCartridge(w http.ResponseWriter, r *http.Request) {
	if h.developer == nil {
		writeError(w, http.StatusServiceUnavailable, "developer service unavailable")
		return
	}
	var body struct {
		DeveloperID   string `json:"developer_id"`
		CartridgeName string `json:"cartridge_name"`
		Category      string `json:"category"`
		Description   string `json:"description"`
		PriceKrw      int32  `json:"price_krw"`
		PackageUrl    string `json:"package_url"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.developer.SubmitCartridge(r.Context(), &v1.SubmitCartridgeRequest{
		DeveloperId: body.DeveloperID, CartridgeName: body.CartridgeName,
		Category: body.Category, Description: body.Description,
		PriceKrw: body.PriceKrw, PackageUrl: body.PackageUrl,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusCreated, resp)
}

func (h *RestHandler) handleListSubmissions(w http.ResponseWriter, r *http.Request) {
	if h.developer == nil {
		writeError(w, http.StatusServiceUnavailable, "developer service unavailable")
		return
	}
	devID := r.URL.Query().Get("developer_id")
	resp, err := h.developer.ListSubmissions(r.Context(), &v1.ListSubmissionsRequest{
		DeveloperId: devID, Limit: queryInt(r, "limit", 20),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleListStoreItems(w http.ResponseWriter, r *http.Request) {
	if h.store == nil {
		writeError(w, http.StatusServiceUnavailable, "store service unavailable")
		return
	}
	resp, err := h.store.ListStoreItems(r.Context(), &v1.ListStoreItemsRequest{
		Category: r.URL.Query().Get("category"),
		Limit:    queryInt(r, "limit", 20),
		Offset:   queryInt(r, "offset", 0),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleSearchCartridges(w http.ResponseWriter, r *http.Request) {
	if h.store == nil {
		writeError(w, http.StatusServiceUnavailable, "store service unavailable")
		return
	}
	resp, err := h.store.SearchCartridges(r.Context(), &v1.SearchCartridgesRequest{
		Query: r.URL.Query().Get("q"), Limit: queryInt(r, "limit", 20),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleGetStoreItem(w http.ResponseWriter, r *http.Request) {
	if h.store == nil {
		writeError(w, http.StatusServiceUnavailable, "store service unavailable")
		return
	}
	resp, err := h.store.GetStoreItem(r.Context(), &v1.GetStoreItemRequest{ItemId: r.PathValue("itemId")})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handlePurchaseCartridge(w http.ResponseWriter, r *http.Request) {
	if h.store == nil {
		writeError(w, http.StatusServiceUnavailable, "store service unavailable")
		return
	}
	var body struct {
		UserID string `json:"user_id"`
		ItemID string `json:"item_id"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.store.PurchaseCartridge(r.Context(), &v1.PurchaseCartridgeRequest{
		UserId: body.UserID, ItemId: body.ItemID,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleGetPurchaseHistory(w http.ResponseWriter, r *http.Request) {
	if h.store == nil {
		writeError(w, http.StatusServiceUnavailable, "store service unavailable")
		return
	}
	userID := r.URL.Query().Get("user_id")
	resp, err := h.store.GetPurchaseHistory(r.Context(), &v1.GetPurchaseHistoryRequest{
		UserId: userID, Limit: queryInt(r, "limit", 20),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleSubmitForReview(w http.ResponseWriter, r *http.Request) {
	if h.cartridgeReview == nil {
		writeError(w, http.StatusServiceUnavailable, "review service unavailable")
		return
	}
	var body struct {
		SubmissionID string `json:"submission_id"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.cartridgeReview.SubmitForReview(r.Context(), &v1.SubmitForReviewRequest{
		SubmissionId: body.SubmissionID,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

// ── Cartridge Type Definition (Phase 5) ──

func (h *RestHandler) handleRegisterCartridgeType(w http.ResponseWriter, r *http.Request) {
	if h.developer == nil {
		writeError(w, http.StatusServiceUnavailable, "developer service unavailable")
		return
	}
	var body struct {
		DeveloperID          string   `json:"developer_id"`
		CategoryCode         int32    `json:"category_code"`
		TypeIndex            int32    `json:"type_index"`
		Name                 string   `json:"name"`
		Version              string   `json:"version"`
		CSIVersion           string   `json:"csi_version"`
		PinConfigJSON        string   `json:"pin_config_json"`
		AFEBlocks            []string `json:"afe_blocks"`
		MeasurementModes     []string `json:"measurement_modes"`
		CalibrationSchemaJSON string  `json:"calibration_schema_json"`
		AlphaDefault         float64  `json:"alpha_default"`
		TempCompensation     bool     `json:"temp_compensation"`
		HumidityCompensation bool     `json:"humidity_compensation"`
		RegulatoryClass      string   `json:"regulatory_class"`
		RegulatoryPath       string   `json:"regulatory_path"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.developer.RegisterCartridgeType(r.Context(), &v1.RegisterCartridgeTypeRequest{
		DeveloperId:           body.DeveloperID,
		CategoryCode:          body.CategoryCode,
		TypeIndex:             body.TypeIndex,
		Name:                  body.Name,
		Version:               body.Version,
		CsiVersion:            body.CSIVersion,
		PinConfigJson:         body.PinConfigJSON,
		AfeBlocks:             body.AFEBlocks,
		MeasurementModes:      body.MeasurementModes,
		CalibrationSchemaJson: body.CalibrationSchemaJSON,
		AlphaDefault:          body.AlphaDefault,
		TempCompensation:      body.TempCompensation,
		HumidityCompensation:  body.HumidityCompensation,
		RegulatoryClass:       body.RegulatoryClass,
		RegulatoryPath:        body.RegulatoryPath,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusCreated, resp)
}

func (h *RestHandler) handleGetCartridgeDefinition(w http.ResponseWriter, r *http.Request) {
	if h.developer == nil {
		writeError(w, http.StatusServiceUnavailable, "developer service unavailable")
		return
	}
	definitionId := r.PathValue("definitionId")
	resp, err := h.developer.GetCartridgeDefinition(r.Context(), &v1.GetCartridgeDefinitionRequest{
		DefinitionId: definitionId,
	})
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleGetReviewStatus(w http.ResponseWriter, r *http.Request) {
	if h.cartridgeReview == nil {
		writeError(w, http.StatusServiceUnavailable, "review service unavailable")
		return
	}
	resp, err := h.cartridgeReview.GetReviewStatus(r.Context(), &v1.GetReviewStatusRequest{
		SubmissionId: r.PathValue("submissionId"),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}
