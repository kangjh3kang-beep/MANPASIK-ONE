package handler

import (
	"net/http"

	v1 "github.com/manpasik/backend/shared/gen/go/v1"
)

// registerMeasurementRoutes는 측정/디바이스/카트리지/캘리브레이션 관련 REST 엔드포인트를 등록합니다.
func (h *RestHandler) registerMeasurementRoutes(mux *http.ServeMux) {
	// Measurement
	mux.HandleFunc("POST /api/v1/measurements/sessions", h.handleStartSession)
	mux.HandleFunc("POST /api/v1/measurements/sessions/{sessionId}/end", h.handleEndSession)
	mux.HandleFunc("GET /api/v1/measurements/history", h.handleGetMeasurementHistory)

	// Device
	mux.HandleFunc("POST /api/v1/devices", h.handleRegisterDevice)
	mux.HandleFunc("GET /api/v1/devices", h.handleListDevices)

	// Cartridge
	mux.HandleFunc("POST /api/v1/cartridges/read", h.handleReadCartridge)
	mux.HandleFunc("POST /api/v1/cartridges/usage", h.handleRecordCartridgeUsage)
	mux.HandleFunc("GET /api/v1/cartridges/types", h.handleListCartridgeTypes)
	mux.HandleFunc("GET /api/v1/cartridges/{cartridgeId}/remaining", h.handleGetRemainingUses)
	mux.HandleFunc("POST /api/v1/cartridges/validate", h.handleValidateCartridge)

	// Calibration
	mux.HandleFunc("POST /api/v1/calibration/factory", h.handleRegisterFactoryCalibration)
	mux.HandleFunc("POST /api/v1/calibration/field", h.handlePerformFieldCalibration)
	mux.HandleFunc("GET /api/v1/calibration/{deviceId}/status", h.handleCheckCalibrationStatus)
	mux.HandleFunc("GET /api/v1/calibration/models", h.handleListCalibrationModels)

	// Health Record
	mux.HandleFunc("POST /api/v1/health-records", h.handleCreateHealthRecord)
	mux.HandleFunc("GET /api/v1/health-records", h.handleListHealthRecords)
	mux.HandleFunc("GET /api/v1/health-records/{recordId}", h.handleGetHealthRecord)
	mux.HandleFunc("POST /api/v1/health-records/export/fhir", h.handleExportToFHIR)
	mux.HandleFunc("POST /api/v1/health-records/import", h.handleImportHealthData)
}

// ── Measurement ──

func (h *RestHandler) handleStartSession(w http.ResponseWriter, r *http.Request) {
	if h.measurement == nil {
		writeError(w, http.StatusServiceUnavailable, "measurement service unavailable")
		return
	}
	var body struct {
		DeviceID           string `json:"device_id"`
		UserID             string `json:"user_id"`
		CartridgeID        string `json:"cartridge_id"`
		CartridgeCategory  int32  `json:"cartridge_category"`
		CartridgeTypeIndex int32  `json:"cartridge_type_index"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.measurement.StartSession(r.Context(), &v1.StartSessionRequest{
		DeviceId:           body.DeviceID,
		UserId:             body.UserID,
		CartridgeId:        body.CartridgeID,
		CartridgeCategory:  body.CartridgeCategory,
		CartridgeTypeIndex: body.CartridgeTypeIndex,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusCreated, resp)
}

func (h *RestHandler) handleEndSession(w http.ResponseWriter, r *http.Request) {
	if h.measurement == nil {
		writeError(w, http.StatusServiceUnavailable, "measurement service unavailable")
		return
	}
	sessionId := r.PathValue("sessionId")
	resp, err := h.measurement.EndSession(r.Context(), &v1.EndSessionRequest{SessionId: sessionId})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleGetMeasurementHistory(w http.ResponseWriter, r *http.Request) {
	if h.measurement == nil {
		writeError(w, http.StatusServiceUnavailable, "measurement service unavailable")
		return
	}
	resp, err := h.measurement.GetMeasurementHistory(r.Context(), &v1.GetHistoryRequest{
		UserId: r.URL.Query().Get("user_id"),
		Limit:  queryInt(r, "limit", 20),
		Offset: queryInt(r, "offset", 0),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

// ── Device ──

func (h *RestHandler) handleRegisterDevice(w http.ResponseWriter, r *http.Request) {
	if h.device == nil {
		writeError(w, http.StatusServiceUnavailable, "device service unavailable")
		return
	}
	var body struct {
		DeviceID        string `json:"device_id"`
		UserID          string `json:"user_id"`
		SerialNumber    string `json:"serial_number"`
		FirmwareVersion string `json:"firmware_version"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.device.RegisterDevice(r.Context(), &v1.RegisterDeviceRequest{
		DeviceId:        body.DeviceID,
		UserId:          body.UserID,
		SerialNumber:    body.SerialNumber,
		FirmwareVersion: body.FirmwareVersion,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusCreated, resp)
}

func (h *RestHandler) handleListDevices(w http.ResponseWriter, r *http.Request) {
	if h.device == nil {
		writeError(w, http.StatusServiceUnavailable, "device service unavailable")
		return
	}
	resp, err := h.device.ListDevices(r.Context(), &v1.ListDevicesRequest{
		UserId: r.URL.Query().Get("user_id"),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

// ── Cartridge ──

func (h *RestHandler) handleReadCartridge(w http.ResponseWriter, r *http.Request) {
	if h.cartridge == nil {
		writeError(w, http.StatusServiceUnavailable, "cartridge service unavailable")
		return
	}
	var body struct {
		NfcTagData []byte `json:"nfc_tag_data"`
		TagVersion int32  `json:"tag_version"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.cartridge.ReadCartridge(r.Context(), &v1.ReadCartridgeRequest{
		NfcTagData: body.NfcTagData,
		TagVersion: body.TagVersion,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleRecordCartridgeUsage(w http.ResponseWriter, r *http.Request) {
	if h.cartridge == nil {
		writeError(w, http.StatusServiceUnavailable, "cartridge service unavailable")
		return
	}
	var body struct {
		UserID       string `json:"user_id"`
		SessionID    string `json:"session_id"`
		CartridgeUID string `json:"cartridge_uid"`
		CategoryCode int32  `json:"category_code"`
		TypeIndex    int32  `json:"type_index"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.cartridge.RecordUsage(r.Context(), &v1.RecordUsageRequest{
		UserId:       body.UserID,
		SessionId:    body.SessionID,
		CartridgeUid: body.CartridgeUID,
		CategoryCode: body.CategoryCode,
		TypeIndex:    body.TypeIndex,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleListCartridgeTypes(w http.ResponseWriter, r *http.Request) {
	if h.cartridge == nil {
		writeError(w, http.StatusServiceUnavailable, "cartridge service unavailable")
		return
	}
	resp, err := h.cartridge.ListCategories(r.Context(), &v1.ListCategoriesRequest{})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleGetRemainingUses(w http.ResponseWriter, r *http.Request) {
	if h.cartridge == nil {
		writeError(w, http.StatusServiceUnavailable, "cartridge service unavailable")
		return
	}
	cartridgeId := r.PathValue("cartridgeId")
	resp, err := h.cartridge.GetRemainingUses(r.Context(), &v1.GetRemainingUsesRequest{
		CartridgeUid: cartridgeId,
	})
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleValidateCartridge(w http.ResponseWriter, r *http.Request) {
	if h.cartridge == nil {
		writeError(w, http.StatusServiceUnavailable, "cartridge service unavailable")
		return
	}
	var body struct {
		CartridgeUID string `json:"cartridge_uid"`
		CategoryCode int32  `json:"category_code"`
		TypeIndex    int32  `json:"type_index"`
		UserID       string `json:"user_id"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.cartridge.ValidateCartridge(r.Context(), &v1.ValidateCartridgeRequest{
		CartridgeUid: body.CartridgeUID,
		CategoryCode: body.CategoryCode,
		TypeIndex:    body.TypeIndex,
		UserId:       body.UserID,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

// ── Calibration ──

func (h *RestHandler) handleRegisterFactoryCalibration(w http.ResponseWriter, r *http.Request) {
	if h.calibration == nil {
		writeError(w, http.StatusServiceUnavailable, "calibration service unavailable")
		return
	}
	var body struct {
		DeviceID              string    `json:"device_id"`
		CartridgeCategory     int32     `json:"cartridge_category"`
		CartridgeTypeIndex    int32     `json:"cartridge_type_index"`
		Alpha                 float64   `json:"alpha"`
		ChannelOffsets        []float64 `json:"channel_offsets"`
		ChannelGains          []float64 `json:"channel_gains"`
		TempCoefficient       float64   `json:"temp_coefficient"`
		HumidityCoefficient   float64   `json:"humidity_coefficient"`
		ReferenceStandard     string    `json:"reference_standard"`
		CalibratedBy          string    `json:"calibrated_by"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.calibration.RegisterFactoryCalibration(r.Context(), &v1.RegisterFactoryCalibrationRequest{
		DeviceId:            body.DeviceID,
		CartridgeCategory:   body.CartridgeCategory,
		CartridgeTypeIndex:  body.CartridgeTypeIndex,
		Alpha:               body.Alpha,
		ChannelOffsets:      body.ChannelOffsets,
		ChannelGains:        body.ChannelGains,
		TempCoefficient:     body.TempCoefficient,
		HumidityCoefficient: body.HumidityCoefficient,
		ReferenceStandard:   body.ReferenceStandard,
		CalibratedBy:        body.CalibratedBy,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusCreated, resp)
}

func (h *RestHandler) handlePerformFieldCalibration(w http.ResponseWriter, r *http.Request) {
	if h.calibration == nil {
		writeError(w, http.StatusServiceUnavailable, "calibration service unavailable")
		return
	}
	var body struct {
		DeviceID           string    `json:"device_id"`
		UserID             string    `json:"user_id"`
		CartridgeCategory  int32     `json:"cartridge_category"`
		CartridgeTypeIndex int32     `json:"cartridge_type_index"`
		ReferenceValues    []float64 `json:"reference_values"`
		MeasuredValues     []float64 `json:"measured_values"`
		TemperatureC       float64   `json:"temperature_c"`
		HumidityPct        float64   `json:"humidity_pct"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.calibration.PerformFieldCalibration(r.Context(), &v1.PerformFieldCalibrationRequest{
		DeviceId:           body.DeviceID,
		UserId:             body.UserID,
		CartridgeCategory:  body.CartridgeCategory,
		CartridgeTypeIndex: body.CartridgeTypeIndex,
		ReferenceValues:    body.ReferenceValues,
		MeasuredValues:     body.MeasuredValues,
		TemperatureC:       body.TemperatureC,
		HumidityPct:        body.HumidityPct,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusCreated, resp)
}

func (h *RestHandler) handleCheckCalibrationStatus(w http.ResponseWriter, r *http.Request) {
	if h.calibration == nil {
		writeError(w, http.StatusServiceUnavailable, "calibration service unavailable")
		return
	}
	deviceId := r.PathValue("deviceId")
	resp, err := h.calibration.CheckCalibrationStatus(r.Context(), &v1.CheckCalibrationStatusRequest{
		DeviceId:           deviceId,
		CartridgeCategory:  queryInt(r, "cartridge_category", 0),
		CartridgeTypeIndex: queryInt(r, "cartridge_type_index", 0),
	})
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleListCalibrationModels(w http.ResponseWriter, r *http.Request) {
	if h.calibration == nil {
		writeError(w, http.StatusServiceUnavailable, "calibration service unavailable")
		return
	}
	resp, err := h.calibration.ListCalibrationModels(r.Context(), &v1.ListCalibrationModelsRequest{})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

// ── Health Record ──

func (h *RestHandler) handleCreateHealthRecord(w http.ResponseWriter, r *http.Request) {
	if h.healthRecord == nil {
		writeError(w, http.StatusServiceUnavailable, "health-record service unavailable")
		return
	}
	var body struct {
		UserID        string            `json:"user_id"`
		RecordType    int32             `json:"record_type"`
		Title         string            `json:"title"`
		Description   string            `json:"description"`
		Provider      string            `json:"provider"`
		Metadata      map[string]string `json:"metadata"`
		MeasurementID string            `json:"measurement_id"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.healthRecord.CreateRecord(r.Context(), &v1.CreateHealthRecordRequest{
		UserId:        body.UserID,
		RecordType:    v1.HealthRecordType(body.RecordType),
		Title:         body.Title,
		Description:   body.Description,
		Provider:      body.Provider,
		Metadata:      body.Metadata,
		MeasurementId: body.MeasurementID,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusCreated, resp)
}

func (h *RestHandler) handleListHealthRecords(w http.ResponseWriter, r *http.Request) {
	if h.healthRecord == nil {
		writeError(w, http.StatusServiceUnavailable, "health-record service unavailable")
		return
	}
	resp, err := h.healthRecord.ListRecords(r.Context(), &v1.ListHealthRecordsRequest{
		UserId:     r.URL.Query().Get("user_id"),
		TypeFilter: v1.HealthRecordType(queryInt(r, "type_filter", 0)),
		Limit:      queryInt(r, "limit", 20),
		Offset:     queryInt(r, "offset", 0),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleGetHealthRecord(w http.ResponseWriter, r *http.Request) {
	if h.healthRecord == nil {
		writeError(w, http.StatusServiceUnavailable, "health-record service unavailable")
		return
	}
	recordId := r.PathValue("recordId")
	resp, err := h.healthRecord.GetRecord(r.Context(), &v1.GetHealthRecordRequest{RecordId: recordId})
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleExportToFHIR(w http.ResponseWriter, r *http.Request) {
	if h.healthRecord == nil {
		writeError(w, http.StatusServiceUnavailable, "health-record service unavailable")
		return
	}
	var body struct {
		UserID     string   `json:"user_id"`
		RecordIDs  []string `json:"record_ids"`
		TargetType int32    `json:"target_type"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.healthRecord.ExportToFHIR(r.Context(), &v1.ExportToFHIRRequest{
		UserId:     body.UserID,
		RecordIds:  body.RecordIDs,
		TargetType: v1.FHIRResourceType(body.TargetType),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (h *RestHandler) handleImportHealthData(w http.ResponseWriter, r *http.Request) {
	if h.healthRecord == nil {
		writeError(w, http.StatusServiceUnavailable, "health-record service unavailable")
		return
	}
	var body struct {
		UserID  string                   `json:"user_id"`
		Source  string                   `json:"source"`
		Records []map[string]interface{} `json:"records"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.healthRecord.ImportFromFHIR(r.Context(), &v1.ImportFromFHIRRequest{
		UserId:         body.UserID,
		FhirBundleJson: body.Source,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}
