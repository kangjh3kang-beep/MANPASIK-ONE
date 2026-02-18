// Package handler는 Gateway REST 핸들러를 제공합니다.
package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	v1 "github.com/manpasik/backend/shared/gen/go/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// RestHandler는 모든 REST API 핸들러를 관리합니다.
type RestHandler struct {
	auth         v1.AuthServiceClient
	user         v1.UserServiceClient
	measurement  v1.MeasurementServiceClient
	device       v1.DeviceServiceClient
	subscription v1.SubscriptionServiceClient
	shop         v1.ShopServiceClient
	payment      v1.PaymentServiceClient
	aiInference  v1.AiInferenceServiceClient
	cartridge    v1.CartridgeServiceClient
	calibration  v1.CalibrationServiceClient
	coaching     v1.CoachingServiceClient
	reservation  v1.ReservationServiceClient
	admin        v1.AdminServiceClient
	family       v1.FamilyServiceClient
	healthRecord v1.HealthRecordServiceClient
	prescription v1.PrescriptionServiceClient
	community    v1.CommunityServiceClient
	video        v1.VideoServiceClient
	notification v1.NotificationServiceClient
	translation  v1.TranslationServiceClient
	telemedicine v1.TelemedicineServiceClient

	jwtSecret string
}

// NewRestHandler는 RestHandler를 생성합니다.
func NewRestHandler(
	auth v1.AuthServiceClient,
	user v1.UserServiceClient,
	measurement v1.MeasurementServiceClient,
	device v1.DeviceServiceClient,
	subscription v1.SubscriptionServiceClient,
	shop v1.ShopServiceClient,
	payment v1.PaymentServiceClient,
	aiInference v1.AiInferenceServiceClient,
	cartridge v1.CartridgeServiceClient,
	calibration v1.CalibrationServiceClient,
	coaching v1.CoachingServiceClient,
	reservation v1.ReservationServiceClient,
	admin v1.AdminServiceClient,
	family v1.FamilyServiceClient,
	healthRecord v1.HealthRecordServiceClient,
	prescription v1.PrescriptionServiceClient,
	community v1.CommunityServiceClient,
	video v1.VideoServiceClient,
	notification v1.NotificationServiceClient,
	translation v1.TranslationServiceClient,
	telemedicine v1.TelemedicineServiceClient,
	jwtSecret string,
) *RestHandler {
	return &RestHandler{
		auth: auth, user: user, measurement: measurement, device: device,
		subscription: subscription, shop: shop, payment: payment,
		aiInference: aiInference, cartridge: cartridge, calibration: calibration,
		coaching: coaching, reservation: reservation, admin: admin,
		family: family, healthRecord: healthRecord, prescription: prescription,
		community: community, video: video, notification: notification,
		translation: translation, telemedicine: telemedicine,
		jwtSecret: jwtSecret,
	}
}

// SetupRoutes는 모든 REST 라우트를 등록합니다.
func (h *RestHandler) SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	// Auth 라우트 (인증 불필요)
	h.registerAuthRoutes(mux)

	// 인증 필요 라우트
	h.registerUserRoutes(mux)
	h.registerMeasurementRoutes(mux)
	h.registerMarketRoutes(mux)
	h.registerCommunityRoutes(mux)

	return mux
}

// ============================================================================
// JSON 유틸리티
// ============================================================================

// writeJSON은 JSON 응답을 작성합니다.
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// writeProtoJSON은 Protobuf 메시지를 JSON으로 직렬화하여 응답합니다.
func writeProtoJSON(w http.ResponseWriter, status int, msg proto.Message) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	marshaler := protojson.MarshalOptions{
		UseProtoNames:   false,
		EmitUnpopulated: true,
	}
	b, err := marshaler.Marshal(msg)
	if err != nil {
		w.Write([]byte(`{"error":"serialization failed"}`))
		return
	}
	w.Write(b)
}

// writeError는 에러 JSON 응답을 작성합니다.
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

// readJSON은 요청 본문을 JSON으로 파싱합니다.
func readJSON(r *http.Request, v interface{}) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.Unmarshal(body, v)
}

// pathParam은 URL 경로에서 파라미터를 추출합니다.
// 예: "/api/v1/users/abc/profile" 에서 pathSegment(r, 4) → "abc"
func pathSegment(r *http.Request, index int) string {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if index < len(parts) {
		return parts[index]
	}
	return ""
}

// queryInt는 쿼리 파라미터를 int로 변환합니다.
func queryInt(r *http.Request, key string, defaultVal int) int32 {
	v := r.URL.Query().Get(key)
	if v == "" {
		return int32(defaultVal)
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return int32(defaultVal)
	}
	return int32(i)
}

// queryFloat는 쿼리 파라미터를 float64로 변환합니다.
func queryFloat(r *http.Request, key string, defaultVal float64) float64 {
	v := r.URL.Query().Get(key)
	if v == "" {
		return defaultVal
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return defaultVal
	}
	return f
}

// queryBool은 쿼리 파라미터를 bool로 변환합니다.
func queryBool(r *http.Request, key string) bool {
	v := r.URL.Query().Get(key)
	return v == "true" || v == "1"
}
