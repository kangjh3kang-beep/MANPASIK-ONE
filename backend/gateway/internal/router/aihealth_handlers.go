package router

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	v1 "github.com/manpasik/backend/shared/gen/go/v1"
)

// ============================================================================
// AI Inference Service Handlers
// ============================================================================

// POST /api/v1/ai/analyze
func (r *Router) handleAnalyzeMeasurement(w http.ResponseWriter, req *http.Request) {
	var body struct {
		UserID        string  `json:"user_id"`
		MeasurementID string  `json:"measurement_id"`
		Models        []int32 `json:"models"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "잘못된 요청 형식")
		return
	}

	conn, err := dialGRPC(r.aiInferenceAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "AI 추론 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewAiInferenceServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	models := make([]v1.AiModelType, 0, len(body.Models))
	for _, m := range body.Models {
		models = append(models, v1.AiModelType(m))
	}

	resp, err := client.AnalyzeMeasurement(ctx, &v1.AnalyzeMeasurementRequest{
		UserId:        body.UserID,
		MeasurementId: body.MeasurementID,
		Models:        models,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, analysisResultToMap(resp))
}

// GET /api/v1/ai/health-score/{userId}
func (r *Router) handleGetHealthScore(w http.ResponseWriter, req *http.Request) {
	userID := req.PathValue("userId")
	if userID == "" {
		writeError(w, http.StatusBadRequest, "user_id가 필요합니다")
		return
	}

	days, _ := strconv.Atoi(req.URL.Query().Get("days"))
	if days <= 0 {
		days = 30
	}

	conn, err := dialGRPC(r.aiInferenceAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "AI 추론 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewAiInferenceServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.GetHealthScore(ctx, &v1.GetHealthScoreRequest{
		UserId: userID,
		Days:   int32(days),
	})
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	result := map[string]interface{}{
		"user_id":         resp.GetUserId(),
		"overall_score":   resp.GetOverallScore(),
		"category_scores": resp.GetCategoryScores(),
		"trend":           resp.GetTrend(),
		"recommendation":  resp.GetRecommendation(),
	}
	if resp.GetCalculatedAt() != nil {
		result["calculated_at"] = resp.GetCalculatedAt().AsTime().Format(time.RFC3339)
	}

	writeJSON(w, http.StatusOK, result)
}

// POST /api/v1/ai/predict-trend
func (r *Router) handlePredictTrend(w http.ResponseWriter, req *http.Request) {
	var body struct {
		UserID         string `json:"user_id"`
		MetricName     string `json:"metric_name"`
		HistoryDays    int32  `json:"history_days"`
		PredictionDays int32  `json:"prediction_days"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "잘못된 요청 형식")
		return
	}

	conn, err := dialGRPC(r.aiInferenceAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "AI 추론 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewAiInferenceServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	resp, err := client.PredictTrend(ctx, &v1.PredictTrendRequest{
		UserId:         body.UserID,
		MetricName:     body.MetricName,
		HistoryDays:    body.HistoryDays,
		PredictionDays: body.PredictionDays,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, trendPredictionToMap(resp))
}

// GET /api/v1/ai/models
func (r *Router) handleListAiModels(w http.ResponseWriter, _ *http.Request) {
	conn, err := dialGRPC(r.aiInferenceAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "AI 추론 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewAiInferenceServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.ListModels(ctx, &v1.ListModelsRequest{})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	models := make([]map[string]interface{}, 0, len(resp.GetModels()))
	for _, m := range resp.GetModels() {
		models = append(models, modelInfoToMap(m))
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"models": models,
	})
}

// --- AI Inference helper mappers ---

func analysisResultToMap(r *v1.AnalysisResult) map[string]interface{} {
	biomarkers := make([]map[string]interface{}, 0, len(r.GetBiomarkers()))
	for _, b := range r.GetBiomarkers() {
		biomarkers = append(biomarkers, map[string]interface{}{
			"biomarker_name":  b.GetBiomarkerName(),
			"value":           b.GetValue(),
			"unit":            b.GetUnit(),
			"classification":  b.GetClassification(),
			"confidence":      b.GetConfidence(),
			"risk_level":      b.GetRiskLevel().String(),
			"reference_range": b.GetReferenceRange(),
		})
	}

	anomalies := make([]map[string]interface{}, 0, len(r.GetAnomalies()))
	for _, a := range r.GetAnomalies() {
		anomalies = append(anomalies, map[string]interface{}{
			"metric_name":   a.GetMetricName(),
			"value":         a.GetValue(),
			"expected_min":  a.GetExpectedMin(),
			"expected_max":  a.GetExpectedMax(),
			"anomaly_score": a.GetAnomalyScore(),
			"description":   a.GetDescription(),
		})
	}

	m := map[string]interface{}{
		"analysis_id":          r.GetAnalysisId(),
		"user_id":              r.GetUserId(),
		"measurement_id":       r.GetMeasurementId(),
		"biomarkers":           biomarkers,
		"anomalies":            anomalies,
		"overall_health_score": r.GetOverallHealthScore(),
		"summary":              r.GetSummary(),
	}
	if r.GetAnalyzedAt() != nil {
		m["analyzed_at"] = r.GetAnalyzedAt().AsTime().Format(time.RFC3339)
	}
	return m
}

func trendPredictionToMap(t *v1.TrendPrediction) map[string]interface{} {
	historical := make([]map[string]interface{}, 0, len(t.GetHistorical()))
	for _, dp := range t.GetHistorical() {
		h := map[string]interface{}{
			"value":       dp.GetValue(),
			"lower_bound": dp.GetLowerBound(),
			"upper_bound": dp.GetUpperBound(),
		}
		if dp.GetTimestamp() != nil {
			h["timestamp"] = dp.GetTimestamp().AsTime().Format(time.RFC3339)
		}
		historical = append(historical, h)
	}

	predicted := make([]map[string]interface{}, 0, len(t.GetPredicted()))
	for _, dp := range t.GetPredicted() {
		p := map[string]interface{}{
			"value":       dp.GetValue(),
			"lower_bound": dp.GetLowerBound(),
			"upper_bound": dp.GetUpperBound(),
		}
		if dp.GetTimestamp() != nil {
			p["timestamp"] = dp.GetTimestamp().AsTime().Format(time.RFC3339)
		}
		predicted = append(predicted, p)
	}

	return map[string]interface{}{
		"user_id":     t.GetUserId(),
		"metric_name": t.GetMetricName(),
		"historical":  historical,
		"predicted":   predicted,
		"confidence":  t.GetConfidence(),
		"direction":   t.GetDirection(),
		"insight":     t.GetInsight(),
	}
}

func modelInfoToMap(m *v1.ModelInfo) map[string]interface{} {
	result := map[string]interface{}{
		"model_type":  m.GetModelType().String(),
		"name":        m.GetName(),
		"version":     m.GetVersion(),
		"description": m.GetDescription(),
		"accuracy":    m.GetAccuracy(),
		"status":      m.GetStatus(),
	}
	if m.GetLastTrained() != nil {
		result["last_trained"] = m.GetLastTrained().AsTime().Format(time.RFC3339)
	}
	return result
}

// ============================================================================
// Cartridge Service Handlers
// ============================================================================

// POST /api/v1/cartridges/read
func (r *Router) handleReadCartridge(w http.ResponseWriter, req *http.Request) {
	var body struct {
		NfcTagData []byte `json:"nfc_tag_data"` // base64-encoded
		TagVersion int32  `json:"tag_version"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "잘못된 요청 형식")
		return
	}

	conn, err := dialGRPC(r.cartridgeAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "카트리지 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewCartridgeServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.ReadCartridge(ctx, &v1.ReadCartridgeRequest{
		NfcTagData: body.NfcTagData,
		TagVersion: body.TagVersion,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, cartridgeDetailToMap(resp))
}

// POST /api/v1/cartridges/usage
func (r *Router) handleRecordCartridgeUsage(w http.ResponseWriter, req *http.Request) {
	var body struct {
		UserID       string `json:"user_id"`
		SessionID    string `json:"session_id"`
		CartridgeUID string `json:"cartridge_uid"`
		CategoryCode int32  `json:"category_code"`
		TypeIndex    int32  `json:"type_index"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "잘못된 요청 형식")
		return
	}

	conn, err := dialGRPC(r.cartridgeAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "카트리지 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewCartridgeServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.RecordUsage(ctx, &v1.RecordUsageRequest{
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

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success":           resp.GetSuccess(),
		"remaining_uses":    resp.GetRemainingUses(),
		"remaining_daily":   resp.GetRemainingDaily(),
		"remaining_monthly": resp.GetRemainingMonthly(),
	})
}

// GET /api/v1/cartridges/types
func (r *Router) handleListCartridgeCategories(w http.ResponseWriter, _ *http.Request) {
	conn, err := dialGRPC(r.cartridgeAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "카트리지 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewCartridgeServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.ListCategories(ctx, &v1.ListCategoriesRequest{})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	categories := make([]map[string]interface{}, 0, len(resp.GetCategories()))
	for _, c := range resp.GetCategories() {
		categories = append(categories, map[string]interface{}{
			"code":        c.GetCode(),
			"name_en":     c.GetNameEn(),
			"name_ko":     c.GetNameKo(),
			"description": c.GetDescription(),
			"type_count":  c.GetTypeCount(),
			"is_active":   c.GetIsActive(),
		})
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"categories": categories,
	})
}

// GET /api/v1/cartridges/{cartridgeId}/remaining
func (r *Router) handleGetRemainingUses(w http.ResponseWriter, req *http.Request) {
	cartridgeID := req.PathValue("cartridgeId")
	if cartridgeID == "" {
		writeError(w, http.StatusBadRequest, "cartridge_id가 필요합니다")
		return
	}

	conn, err := dialGRPC(r.cartridgeAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "카트리지 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewCartridgeServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetRemainingUses(ctx, &v1.GetRemainingUsesRequest{
		CartridgeUid: cartridgeID,
	})
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"cartridge_uid":  resp.GetCartridgeUid(),
		"remaining_uses": resp.GetRemainingUses(),
		"max_uses":       resp.GetMaxUses(),
		"expiry_date":    resp.GetExpiryDate(),
		"is_expired":     resp.GetIsExpired(),
	})
}

// POST /api/v1/cartridges/validate
func (r *Router) handleValidateCartridge(w http.ResponseWriter, req *http.Request) {
	var body struct {
		CartridgeUID string `json:"cartridge_uid"`
		CategoryCode int32  `json:"category_code"`
		TypeIndex    int32  `json:"type_index"`
		UserID       string `json:"user_id"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "잘못된 요청 형식")
		return
	}

	conn, err := dialGRPC(r.cartridgeAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "카트리지 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewCartridgeServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.ValidateCartridge(ctx, &v1.ValidateCartridgeRequest{
		CartridgeUid: body.CartridgeUID,
		CategoryCode: body.CategoryCode,
		TypeIndex:    body.TypeIndex,
		UserId:       body.UserID,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	result := map[string]interface{}{
		"is_valid":       resp.GetIsValid(),
		"reason":         resp.GetReason(),
		"remaining_uses": resp.GetRemainingUses(),
		"access_level":   resp.GetAccessLevel().String(),
	}
	if resp.GetDetail() != nil {
		result["detail"] = cartridgeDetailToMap(resp.GetDetail())
	}

	writeJSON(w, http.StatusOK, result)
}

func cartridgeDetailToMap(c *v1.CartridgeDetail) map[string]interface{} {
	return map[string]interface{}{
		"cartridge_uid":        c.GetCartridgeUid(),
		"category_code":        c.GetCategoryCode(),
		"type_index":           c.GetTypeIndex(),
		"legacy_code":          c.GetLegacyCode(),
		"name_ko":              c.GetNameKo(),
		"name_en":              c.GetNameEn(),
		"lot_id":               c.GetLotId(),
		"expiry_date":          c.GetExpiryDate(),
		"remaining_uses":       c.GetRemainingUses(),
		"max_uses":             c.GetMaxUses(),
		"alpha_coefficient":    c.GetAlphaCoefficient(),
		"temp_coefficient":     c.GetTempCoefficient(),
		"humidity_coefficient": c.GetHumidityCoefficient(),
		"required_channels":    c.GetRequiredChannels(),
	}
}

// ============================================================================
// Calibration Service Handlers
// ============================================================================

// POST /api/v1/calibration/factory
func (r *Router) handleRegisterFactoryCalibration(w http.ResponseWriter, req *http.Request) {
	var body struct {
		DeviceID            string    `json:"device_id"`
		CartridgeCategory   int32     `json:"cartridge_category"`
		CartridgeTypeIndex  int32     `json:"cartridge_type_index"`
		Alpha               float64   `json:"alpha"`
		ChannelOffsets      []float64 `json:"channel_offsets"`
		ChannelGains        []float64 `json:"channel_gains"`
		TempCoefficient     float64   `json:"temp_coefficient"`
		HumidityCoefficient float64   `json:"humidity_coefficient"`
		ReferenceStandard   string    `json:"reference_standard"`
		CalibratedBy        string    `json:"calibrated_by"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "잘못된 요청 형식")
		return
	}

	conn, err := dialGRPC(r.calibrationAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "보정 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewCalibrationServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.RegisterFactoryCalibration(ctx, &v1.RegisterFactoryCalibrationRequest{
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

	writeJSON(w, http.StatusCreated, calibrationRecordToMap(resp))
}

// POST /api/v1/calibration/field
func (r *Router) handlePerformFieldCalibration(w http.ResponseWriter, req *http.Request) {
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
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "잘못된 요청 형식")
		return
	}

	conn, err := dialGRPC(r.calibrationAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "보정 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewCalibrationServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.PerformFieldCalibration(ctx, &v1.PerformFieldCalibrationRequest{
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

	writeJSON(w, http.StatusCreated, calibrationRecordToMap(resp))
}

// GET /api/v1/calibration/{deviceId}/status?cartridge_category=...&cartridge_type_index=...
func (r *Router) handleCheckCalibrationStatus(w http.ResponseWriter, req *http.Request) {
	deviceID := req.PathValue("deviceId")
	if deviceID == "" {
		writeError(w, http.StatusBadRequest, "device_id가 필요합니다")
		return
	}

	q := req.URL.Query()
	cartridgeCategory, _ := strconv.Atoi(q.Get("cartridge_category"))
	cartridgeTypeIndex, _ := strconv.Atoi(q.Get("cartridge_type_index"))

	conn, err := dialGRPC(r.calibrationAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "보정 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewCalibrationServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.CheckCalibrationStatus(ctx, &v1.CheckCalibrationStatusRequest{
		DeviceId:           deviceID,
		CartridgeCategory:  int32(cartridgeCategory),
		CartridgeTypeIndex: int32(cartridgeTypeIndex),
	})
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	result := map[string]interface{}{
		"status":            resp.GetStatus().String(),
		"device_id":         resp.GetDeviceId(),
		"message":           resp.GetMessage(),
		"days_until_expiry": resp.GetDaysUntilExpiry(),
	}
	if resp.GetLastCalibratedAt() != nil {
		result["last_calibrated_at"] = resp.GetLastCalibratedAt().AsTime().Format(time.RFC3339)
	}
	if resp.GetExpiresAt() != nil {
		result["expires_at"] = resp.GetExpiresAt().AsTime().Format(time.RFC3339)
	}
	if resp.GetLatestRecord() != nil {
		result["latest_record"] = calibrationRecordToMap(resp.GetLatestRecord())
	}

	writeJSON(w, http.StatusOK, result)
}

// GET /api/v1/calibration/models
func (r *Router) handleListCalibrationModels(w http.ResponseWriter, _ *http.Request) {
	conn, err := dialGRPC(r.calibrationAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "보정 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewCalibrationServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.ListCalibrationModels(ctx, &v1.ListCalibrationModelsRequest{})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	models := make([]map[string]interface{}, 0, len(resp.GetModels()))
	for _, m := range resp.GetModels() {
		mm := map[string]interface{}{
			"model_id":              m.GetModelId(),
			"cartridge_category":    m.GetCartridgeCategory(),
			"cartridge_type_index":  m.GetCartridgeTypeIndex(),
			"name":                  m.GetName(),
			"version":              m.GetVersion(),
			"default_alpha":         m.GetDefaultAlpha(),
			"validity_days":         m.GetValidityDays(),
			"description":           m.GetDescription(),
		}
		if m.GetCreatedAt() != nil {
			mm["created_at"] = m.GetCreatedAt().AsTime().Format(time.RFC3339)
		}
		models = append(models, mm)
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"models": models,
	})
}

func calibrationRecordToMap(c *v1.CalibrationRecord) map[string]interface{} {
	m := map[string]interface{}{
		"calibration_id":       c.GetCalibrationId(),
		"device_id":            c.GetDeviceId(),
		"cartridge_category":   c.GetCartridgeCategory(),
		"cartridge_type_index": c.GetCartridgeTypeIndex(),
		"calibration_type":     c.GetCalibrationType().String(),
		"alpha":                c.GetAlpha(),
		"channel_offsets":      c.GetChannelOffsets(),
		"channel_gains":        c.GetChannelGains(),
		"temp_coefficient":     c.GetTempCoefficient(),
		"humidity_coefficient": c.GetHumidityCoefficient(),
		"accuracy_score":       c.GetAccuracyScore(),
		"reference_standard":   c.GetReferenceStandard(),
		"calibrated_by":        c.GetCalibratedBy(),
		"status":               c.GetStatus().String(),
	}
	if c.GetCalibratedAt() != nil {
		m["calibrated_at"] = c.GetCalibratedAt().AsTime().Format(time.RFC3339)
	}
	if c.GetExpiresAt() != nil {
		m["expires_at"] = c.GetExpiresAt().AsTime().Format(time.RFC3339)
	}
	return m
}

// ============================================================================
// Coaching Service Handlers
// ============================================================================

// POST /api/v1/coaching/goals
func (r *Router) handleSetHealthGoal(w http.ResponseWriter, req *http.Request) {
	var body struct {
		UserID      string  `json:"user_id"`
		Category    int32   `json:"category"`
		MetricName  string  `json:"metric_name"`
		TargetValue float64 `json:"target_value"`
		Unit        string  `json:"unit"`
		Description string  `json:"description"`
		TargetDate  string  `json:"target_date"` // RFC3339
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "잘못된 요청 형식")
		return
	}

	conn, err := dialGRPC(r.coachingAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "코칭 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewCoachingServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	grpcReq := &v1.SetHealthGoalRequest{
		UserId:      body.UserID,
		Category:    v1.GoalCategory(body.Category),
		MetricName:  body.MetricName,
		TargetValue: body.TargetValue,
		Unit:        body.Unit,
		Description: body.Description,
	}

	resp, err := client.SetHealthGoal(ctx, grpcReq)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, healthGoalToMap(resp))
}

// GET /api/v1/coaching/goals/{userId}?status_filter=...
func (r *Router) handleGetHealthGoals(w http.ResponseWriter, req *http.Request) {
	userID := req.PathValue("userId")
	if userID == "" {
		writeError(w, http.StatusBadRequest, "user_id가 필요합니다")
		return
	}

	statusFilter, _ := strconv.Atoi(req.URL.Query().Get("status_filter"))

	conn, err := dialGRPC(r.coachingAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "코칭 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewCoachingServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetHealthGoals(ctx, &v1.GetHealthGoalsRequest{
		UserId:       userID,
		StatusFilter: v1.GoalStatus(statusFilter),
	})
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	goals := make([]map[string]interface{}, 0, len(resp.GetGoals()))
	for _, g := range resp.GetGoals() {
		goals = append(goals, healthGoalToMap(g))
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"goals": goals,
	})
}

// POST /api/v1/coaching/generate
func (r *Router) handleGenerateCoaching(w http.ResponseWriter, req *http.Request) {
	var body struct {
		UserID        string `json:"user_id"`
		MeasurementID string `json:"measurement_id"`
		CoachingType  int32  `json:"coaching_type"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "잘못된 요청 형식")
		return
	}

	conn, err := dialGRPC(r.coachingAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "코칭 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewCoachingServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	resp, err := client.GenerateCoaching(ctx, &v1.GenerateCoachingRequest{
		UserId:        body.UserID,
		MeasurementId: body.MeasurementID,
		CoachingType:  v1.CoachingType(body.CoachingType),
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, coachingMessageToMap(resp))
}

// GET /api/v1/coaching/daily-report/{userId}
func (r *Router) handleGenerateDailyReport(w http.ResponseWriter, req *http.Request) {
	userID := req.PathValue("userId")
	if userID == "" {
		writeError(w, http.StatusBadRequest, "user_id가 필요합니다")
		return
	}

	conn, err := dialGRPC(r.coachingAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "코칭 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewCoachingServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	resp, err := client.GenerateDailyReport(ctx, &v1.GenerateDailyReportRequest{
		UserId: userID,
	})
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	highlights := make([]map[string]interface{}, 0, len(resp.GetHighlights()))
	for _, h := range resp.GetHighlights() {
		highlights = append(highlights, coachingMessageToMap(h))
	}

	result := map[string]interface{}{
		"report_id":          resp.GetReportId(),
		"user_id":            resp.GetUserId(),
		"overall_score":      resp.GetOverallScore(),
		"measurements_count": resp.GetMeasurementsCount(),
		"highlights":         highlights,
		"summary":            resp.GetSummary(),
		"recommendations":    resp.GetRecommendations(),
	}
	if resp.GetReportDate() != nil {
		result["report_date"] = resp.GetReportDate().AsTime().Format(time.RFC3339)
	}

	writeJSON(w, http.StatusOK, result)
}

// GET /api/v1/coaching/recommendations/{userId}?type_filter=...&limit=...
func (r *Router) handleGetRecommendations(w http.ResponseWriter, req *http.Request) {
	userID := req.PathValue("userId")
	if userID == "" {
		writeError(w, http.StatusBadRequest, "user_id가 필요합니다")
		return
	}

	q := req.URL.Query()
	typeFilter, _ := strconv.Atoi(q.Get("type_filter"))
	limit, _ := strconv.Atoi(q.Get("limit"))
	if limit <= 0 {
		limit = 10
	}

	conn, err := dialGRPC(r.coachingAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "코칭 서비스 연결 실패")
		return
	}
	defer conn.Close()

	client := v1.NewCoachingServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.GetRecommendations(ctx, &v1.GetRecommendationsRequest{
		UserId:     userID,
		TypeFilter: v1.RecommendationType(typeFilter),
		Limit:      int32(limit),
	})
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	recommendations := make([]map[string]interface{}, 0, len(resp.GetRecommendations()))
	for _, rec := range resp.GetRecommendations() {
		rm := map[string]interface{}{
			"recommendation_id": rec.GetRecommendationId(),
			"type":              rec.GetType().String(),
			"title":             rec.GetTitle(),
			"description":       rec.GetDescription(),
			"reason":            rec.GetReason(),
			"priority":          rec.GetPriority().String(),
			"action_steps":      rec.GetActionSteps(),
			"related_metric":    rec.GetRelatedMetric(),
		}
		if rec.GetCreatedAt() != nil {
			rm["created_at"] = rec.GetCreatedAt().AsTime().Format(time.RFC3339)
		}
		recommendations = append(recommendations, rm)
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"recommendations": recommendations,
	})
}

// --- Coaching helper mappers ---

func healthGoalToMap(g *v1.HealthGoal) map[string]interface{} {
	m := map[string]interface{}{
		"goal_id":       g.GetGoalId(),
		"user_id":       g.GetUserId(),
		"category":      g.GetCategory().String(),
		"metric_name":   g.GetMetricName(),
		"target_value":  g.GetTargetValue(),
		"current_value": g.GetCurrentValue(),
		"unit":          g.GetUnit(),
		"progress_pct":  g.GetProgressPct(),
		"status":        g.GetStatus().String(),
		"description":   g.GetDescription(),
	}
	if g.GetCreatedAt() != nil {
		m["created_at"] = g.GetCreatedAt().AsTime().Format(time.RFC3339)
	}
	if g.GetTargetDate() != nil {
		m["target_date"] = g.GetTargetDate().AsTime().Format(time.RFC3339)
	}
	if g.GetAchievedAt() != nil {
		m["achieved_at"] = g.GetAchievedAt().AsTime().Format(time.RFC3339)
	}
	return m
}

func coachingMessageToMap(c *v1.CoachingMessage) map[string]interface{} {
	m := map[string]interface{}{
		"message_id":     c.GetMessageId(),
		"user_id":        c.GetUserId(),
		"coaching_type":  c.GetCoachingType().String(),
		"title":          c.GetTitle(),
		"body":           c.GetBody(),
		"risk_level":     c.GetRiskLevel().String(),
		"action_items":   c.GetActionItems(),
		"related_metric": c.GetRelatedMetric(),
		"related_value":  c.GetRelatedValue(),
	}
	if c.GetCreatedAt() != nil {
		m["created_at"] = c.GetCreatedAt().AsTime().Format(time.RFC3339)
	}
	return m
}
