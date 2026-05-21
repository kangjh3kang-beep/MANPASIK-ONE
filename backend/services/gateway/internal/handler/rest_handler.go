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

	// 감사/디지털 트윈 전용 서비스 (H6 실패격리: 미연결 시 기존 서비스 fallback)
	audit         v1.AdminServiceClient       // audit-service (AdminService 위임)
	auditRecorder auditEventRecorder          // audit-service HTTP write intake
	digitalTwin   v1.MeasurementServiceClient // digital-twin-service (MeasurementService 위임)

	// Sprint 16-24 신규 서비스
	assistant       v1.AssistantServiceClient
	vision          v1.VisionServiceClient
	concept         v1.ConceptServiceClient
	organization    v1.OrganizationServiceClient
	developer       v1.DeveloperServiceClient
	store           v1.StoreServiceClient
	cartridgeReview v1.CartridgeReviewServiceClient
	revenue         v1.RevenueServiceClient
	cartridgeAnalyt v1.CartridgeAnalyticsServiceClient
	locationStats   v1.LocationStatsServiceClient
	dataProvision   v1.DataProvisionServiceClient
	voiceProfile    v1.VoiceProfileServiceClient

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

// SetAssistantClient는 AI 비서 gRPC 클라이언트를 설정합니다.
func (h *RestHandler) SetAssistantClient(c v1.AssistantServiceClient) { h.assistant = c }

// SetVisionClient는 비전 gRPC 클라이언트를 설정합니다.
func (h *RestHandler) SetVisionClient(c v1.VisionServiceClient) { h.vision = c }

// SetConceptClients는 컨셉/조직 gRPC 클라이언트를 설정합니다.
func (h *RestHandler) SetConceptClients(c v1.ConceptServiceClient, o v1.OrganizationServiceClient) {
	h.concept = c
	h.organization = o
}

// SetCartridgeStoreClients는 카트리지 스토어 gRPC 클라이언트를 설정합니다.
func (h *RestHandler) SetCartridgeStoreClients(
	dev v1.DeveloperServiceClient,
	store v1.StoreServiceClient,
	review v1.CartridgeReviewServiceClient,
	rev v1.RevenueServiceClient,
	analyt v1.CartridgeAnalyticsServiceClient,
) {
	h.developer = dev
	h.store = store
	h.cartridgeReview = review
	h.revenue = rev
	h.cartridgeAnalyt = analyt
}

// SetDataPlatformClients는 데이터 플랫폼 gRPC 클라이언트를 설정합니다.
func (h *RestHandler) SetDataPlatformClients(ls v1.LocationStatsServiceClient, dp v1.DataProvisionServiceClient) {
	h.locationStats = ls
	h.dataProvision = dp
}

// SetVoiceProfileClient는 음성 프로필 gRPC 클라이언트를 설정합니다.
func (h *RestHandler) SetVoiceProfileClient(c v1.VoiceProfileServiceClient) { h.voiceProfile = c }

// SetAuditClient는 감사 서비스 gRPC 클라이언트를 설정합니다. (H6: 미연결 시 admin fallback)
func (h *RestHandler) SetAuditClient(c v1.AdminServiceClient) { h.audit = c }

func (h *RestHandler) SetAuditEventRecorder(recorder auditEventRecorder) {
	h.auditRecorder = recorder
}

// SetDigitalTwinClient는 디지털 트윈 gRPC 클라이언트를 설정합니다. (H6: 미연결 시 measurement fallback)
func (h *RestHandler) SetDigitalTwinClient(c v1.MeasurementServiceClient) { h.digitalTwin = c }

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
	h.registerNotificationRoutes(mux)
	h.registerTranslationRoutes(mux)
	h.registerSubscriptionRoutes(mux)
	h.registerCoachingRoutes(mux)
	h.registerAdminRoutes(mux)

	// Sprint 16-24 신규 라우트
	h.registerAssistantRoutes(mux)
	h.registerVisionRoutes(mux)
	h.registerConceptRoutes(mux)
	h.registerCartridgeStoreRoutes(mux)
	h.registerDataPlatformRoutes(mux)
	h.registerVoiceProfileRoutes(mux)

	// 외부 HL7 v2 통합 (Phase AU-2)
	h.registerExternalHL7Routes(mux)

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
		UseProtoNames:   true,
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

// readJSON은 요청 본문을 JSON으로 파싱합니다. 바디 크기를 10MB로 제한합니다.
func readJSON(r *http.Request, v interface{}) error {
	reader := io.LimitReader(r.Body, 10<<20)
	body, err := io.ReadAll(reader)
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
