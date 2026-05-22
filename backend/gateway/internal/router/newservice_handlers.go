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
// IoT Gateway Service Handlers
// ============================================================================

// GET /api/v1/iot/devices?device_id=...
func (r *Router) handleListIoTDevices(w http.ResponseWriter, req *http.Request) {
	deviceID := req.URL.Query().Get("device_id")

	cb, err := r.dialWithBreaker("iot-gateway", r.iotGatewayAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "IoT 게이트웨이 서비스 연결 실패")
		return
	}
	defer cb.Close()

	client := v1.NewIoTGatewayServiceClient(cb.conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if deviceID != "" {
		resp, err := client.GetDevice(ctx, &v1.GetIoTDeviceRequest{
			DeviceId: deviceID,
		})
		if err != nil {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"device_id": resp.DeviceId,
			"protocol":  resp.Protocol,
			"status":    resp.Status,
			"metadata":  resp.Metadata,
		})
		return
	}

	// No device_id: list device data (fallback)
	limit, _ := strconv.Atoi(req.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 20
	}
	offset, _ := strconv.Atoi(req.URL.Query().Get("offset"))

	resp, err := client.ListDeviceData(ctx, &v1.ListDeviceDataRequest{
		DeviceId: "",
		Limit:    int32(limit),
		Offset:   int32(offset),
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	data := make([]map[string]interface{}, 0, len(resp.Data))
	for _, d := range resp.Data {
		entry := map[string]interface{}{
			"id":        d.Id,
			"device_id": d.DeviceId,
			"data_type": d.DataType,
			"value":     d.Value,
			"unit":      d.Unit,
		}
		if d.ReceivedAt != nil {
			entry["received_at"] = d.ReceivedAt.AsTime().Format(time.RFC3339)
		}
		data = append(data, entry)
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data": data,
	})
}

// GET /api/v1/iot/telemetry?device_id=...&limit=...&offset=...
func (r *Router) handleGetIoTTelemetry(w http.ResponseWriter, req *http.Request) {
	deviceID := req.URL.Query().Get("device_id")
	if deviceID == "" {
		writeError(w, http.StatusBadRequest, "device_id 쿼리 파라미터가 필요합니다")
		return
	}

	limit, _ := strconv.Atoi(req.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 20
	}
	offset, _ := strconv.Atoi(req.URL.Query().Get("offset"))

	cb, err := r.dialWithBreaker("iot-gateway", r.iotGatewayAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "IoT 게이트웨이 서비스 연결 실패")
		return
	}
	defer cb.Close()

	client := v1.NewIoTGatewayServiceClient(cb.conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetDeviceTelemetry(ctx, &v1.GetDeviceTelemetryRequest{
		DeviceId: deviceID,
		Limit:    int32(limit),
		Offset:   int32(offset),
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	data := make([]map[string]interface{}, 0, len(resp.Data))
	for _, d := range resp.Data {
		entry := map[string]interface{}{
			"id":        d.Id,
			"device_id": d.DeviceId,
			"data_type": d.DataType,
			"value":     d.Value,
			"unit":      d.Unit,
		}
		if d.ReceivedAt != nil {
			entry["received_at"] = d.ReceivedAt.AsTime().Format(time.RFC3339)
		}
		data = append(data, entry)
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"telemetry": data,
	})
}

// ============================================================================
// NLP Service Handlers
// ============================================================================

// POST /api/v1/nlp/parse
func (r *Router) handleNlpParse(w http.ResponseWriter, req *http.Request) {
	var body struct {
		UserID string `json:"user_id"`
		Text   string `json:"text"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "잘못된 요청 형식")
		return
	}

	cb, err := r.dialWithBreaker("nlp", r.nlpAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "NLP 서비스 연결 실패")
		return
	}
	defer cb.Close()

	client := v1.NewNLPServiceClient(cb.conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.ParseHealthQuery(ctx, &v1.ParseHealthQueryRequest{
		UserId: body.UserID,
		Text:   body.Text,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"id":         resp.Id,
		"user_id":    resp.UserId,
		"raw_text":   resp.RawText,
		"intent":     resp.Intent,
		"entities":   resp.Entities,
		"confidence": resp.Confidence,
		"urgency":    resp.Urgency,
	})
}

// POST /api/v1/nlp/symptoms
func (r *Router) handleNlpSymptoms(w http.ResponseWriter, req *http.Request) {
	var body struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "잘못된 요청 형식")
		return
	}

	cb, err := r.dialWithBreaker("nlp", r.nlpAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "NLP 서비스 연결 실패")
		return
	}
	defer cb.Close()

	client := v1.NewNLPServiceClient(cb.conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.ExtractSymptoms(ctx, &v1.ExtractSymptomsRequest{
		Text: body.Text,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	symptoms := make([]map[string]interface{}, 0, len(resp.Symptoms))
	for _, s := range resp.Symptoms {
		symptoms = append(symptoms, map[string]interface{}{
			"name":       s.Name,
			"severity":   s.Severity,
			"body_part":  s.BodyPart,
			"confidence": s.Confidence,
		})
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"id":       resp.Id,
		"text":     resp.Text,
		"symptoms": symptoms,
		"urgency":  resp.Urgency,
	})
}

// ============================================================================
// Marketplace Service Handlers
// ============================================================================

// GET /api/v1/marketplace/products?partner_id=...&category=...&limit=...&offset=...
func (r *Router) handleListMarketplaceProducts(w http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	limit, _ := strconv.Atoi(q.Get("limit"))
	if limit <= 0 {
		limit = 20
	}
	offset, _ := strconv.Atoi(q.Get("offset"))

	cb, err := r.dialWithBreaker("marketplace", r.marketplaceAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "마켓플레이스 서비스 연결 실패")
		return
	}
	defer cb.Close()

	client := v1.NewMarketplaceServiceClient(cb.conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.ListProducts(ctx, &v1.ListMarketplaceProductsRequest{
		PartnerId: q.Get("partner_id"),
		Category:  q.Get("category"),
		Limit:     int32(limit),
		Offset:    int32(offset),
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	products := make([]map[string]interface{}, 0, len(resp.Products))
	for _, p := range resp.Products {
		entry := map[string]interface{}{
			"id":          p.Id,
			"partner_id":  p.PartnerId,
			"name":        p.Name,
			"description": p.Description,
			"price":       p.Price,
			"category":    p.Category,
			"image_url":   p.ImageUrl,
			"is_active":   p.IsActive,
		}
		if p.CreatedAt != nil {
			entry["created_at"] = p.CreatedAt.AsTime().Format(time.RFC3339)
		}
		products = append(products, entry)
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"products": products,
		"total":    resp.Total,
	})
}

// GET /api/v1/marketplace/partners?partner_id=...
func (r *Router) handleListMarketplacePartners(w http.ResponseWriter, req *http.Request) {
	partnerID := req.URL.Query().Get("partner_id")
	if partnerID == "" {
		writeError(w, http.StatusBadRequest, "partner_id 쿼리 파라미터가 필요합니다")
		return
	}

	cb, err := r.dialWithBreaker("marketplace", r.marketplaceAddr)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "마켓플레이스 서비스 연결 실패")
		return
	}
	defer cb.Close()

	client := v1.NewMarketplaceServiceClient(cb.conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetPartnerStats(ctx, &v1.GetPartnerStatsRequest{
		PartnerId: partnerID,
	})
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"partner_id":     resp.PartnerId,
		"total_products": resp.TotalProducts,
		"total_orders":   resp.TotalOrders,
		"revenue":        resp.Revenue,
	})
}
