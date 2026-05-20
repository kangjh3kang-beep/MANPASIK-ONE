package handler

import (
	"net/http"

	v1 "github.com/manpasik/backend/shared/gen/go/v1"
)

func (h *RestHandler) registerDataPlatformRoutes(mux *http.ServeMux) {
	// LocationStats
	mux.HandleFunc("GET /api/v1/data/region-stats", h.handleGetRegionStats)
	mux.HandleFunc("GET /api/v1/data/hotspots", h.handleListHotspots)
	mux.HandleFunc("GET /api/v1/data/region-trend", h.handleGetTrendByRegion)
	// DataProvision
	mux.HandleFunc("GET /api/v1/data/datasets", h.handleListDatasets)
	mux.HandleFunc("GET /api/v1/data/datasets/{datasetId}", h.handleGetAnonymizedDataset)
	mux.HandleFunc("POST /api/v1/data/access-request", h.handleRequestDataAccess)
	mux.HandleFunc("GET /api/v1/data/consent", h.handleGetDataSharingConsent)
	mux.HandleFunc("PUT /api/v1/data/consent", h.handleUpdateDataSharingConsent)
	mux.HandleFunc("POST /api/v1/data/aggregate", h.handleTriggerAggregation)
	mux.HandleFunc("GET /api/v1/data/insights", h.handleGetAggregatedInsights)
}

func (h *RestHandler) handleGetRegionStats(w http.ResponseWriter, r *http.Request) {
	if h.locationStats == nil {
		writeError(w, http.StatusServiceUnavailable, "location stats service unavailable")
		return
	}
	resp, err := h.locationStats.GetRegionStats(r.Context(), &v1.GetRegionStatsRequest{
		RegionCode: r.URL.Query().Get("region_code"),
		Period:     r.URL.Query().Get("period"),
		Biomarker:  r.URL.Query().Get("biomarker"),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleListHotspots(w http.ResponseWriter, r *http.Request) {
	if h.locationStats == nil {
		writeError(w, http.StatusServiceUnavailable, "location stats service unavailable")
		return
	}
	resp, err := h.locationStats.ListHotspots(r.Context(), &v1.ListHotspotsRequest{
		Biomarker: r.URL.Query().Get("biomarker"),
		RiskLevel: r.URL.Query().Get("risk_level"),
		Limit:     queryInt(r, "limit", 10),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleGetTrendByRegion(w http.ResponseWriter, r *http.Request) {
	if h.locationStats == nil {
		writeError(w, http.StatusServiceUnavailable, "location stats service unavailable")
		return
	}
	resp, err := h.locationStats.GetTrendByRegion(r.Context(), &v1.GetTrendByRegionRequest{
		RegionCode: r.URL.Query().Get("region_code"),
		Biomarker:  r.URL.Query().Get("biomarker"),
		StartDate:  r.URL.Query().Get("start_date"),
		EndDate:    r.URL.Query().Get("end_date"),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleListDatasets(w http.ResponseWriter, r *http.Request) {
	if h.dataProvision == nil {
		writeError(w, http.StatusServiceUnavailable, "data provision service unavailable")
		return
	}
	resp, err := h.dataProvision.ListDatasets(r.Context(), &v1.ListDatasetsRequest{
		DatasetType: r.URL.Query().Get("type"),
		Limit:       queryInt(r, "limit", 20),
		Offset:      queryInt(r, "offset", 0),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleGetAnonymizedDataset(w http.ResponseWriter, r *http.Request) {
	if h.dataProvision == nil {
		writeError(w, http.StatusServiceUnavailable, "data provision service unavailable")
		return
	}
	resp, err := h.dataProvision.GetAnonymizedDataset(r.Context(), &v1.GetAnonymizedDatasetRequest{
		DatasetId: r.PathValue("datasetId"),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleRequestDataAccess(w http.ResponseWriter, r *http.Request) {
	if h.dataProvision == nil {
		writeError(w, http.StatusServiceUnavailable, "data provision service unavailable")
		return
	}
	var body struct {
		RequesterID string `json:"requester_id"`
		DatasetID   string `json:"dataset_id"`
		Purpose     string `json:"purpose"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.dataProvision.RequestDataAccess(r.Context(), &v1.RequestDataAccessRequest{
		RequesterId: body.RequesterID, DatasetId: body.DatasetID, Purpose: body.Purpose,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleGetDataSharingConsent(w http.ResponseWriter, r *http.Request) {
	if h.dataProvision == nil {
		writeError(w, http.StatusServiceUnavailable, "data provision service unavailable")
		return
	}
	resp, err := h.dataProvision.GetDataSharingConsent(r.Context(), &v1.GetConsentRequest{
		UserId: r.URL.Query().Get("user_id"),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleUpdateDataSharingConsent(w http.ResponseWriter, r *http.Request) {
	if h.dataProvision == nil {
		writeError(w, http.StatusServiceUnavailable, "data provision service unavailable")
		return
	}
	var body struct {
		UserID            string `json:"user_id"`
		AnonymousResearch bool   `json:"anonymous_research"`
		RegionStats       bool   `json:"region_stats"`
		EnterpriseData    bool   `json:"enterprise_data"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.dataProvision.UpdateDataSharingConsent(r.Context(), &v1.UpdateConsentRequest{
		UserId: body.UserID, AnonymousResearch: body.AnonymousResearch,
		RegionStats: body.RegionStats, EnterpriseData: body.EnterpriseData,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleTriggerAggregation(w http.ResponseWriter, r *http.Request) {
	if h.dataProvision == nil {
		writeError(w, http.StatusServiceUnavailable, "data provision service unavailable")
		return
	}
	var body struct {
		Biomarker string `json:"biomarker"`
		Date      string `json:"date"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.dataProvision.TriggerAnonymousAggregation(r.Context(), &v1.TriggerAggregationRequest{
		Biomarker: body.Biomarker, Date: body.Date,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleGetAggregatedInsights(w http.ResponseWriter, r *http.Request) {
	if h.dataProvision == nil {
		writeError(w, http.StatusServiceUnavailable, "data provision service unavailable")
		return
	}
	resp, err := h.dataProvision.GetAggregatedInsights(r.Context(), &v1.GetInsightsRequest{
		Biomarker: r.URL.Query().Get("biomarker"),
		Limit:     queryInt(r, "limit", 10),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}
