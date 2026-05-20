// gateway: REST → gRPC 브릿지 API 게이트웨이
//
// 포트: HTTP :8080
// 역할: Flutter REST Client ↔ 백엔드 gRPC 서비스 브릿지
//
// 기능:
// - REST API 엔드포인트 → gRPC 호출 변환
// - JWT 인증 미들웨어 (Bearer Token)
// - CORS 설정
// - 요청/응답 JSON 직렬화
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manpasik/backend/services/gateway/internal/handler"
	gw "github.com/manpasik/backend/services/gateway/internal/middleware"
	"github.com/manpasik/backend/shared/config"
	v1 "github.com/manpasik/backend/shared/gen/go/v1"
	"github.com/manpasik/backend/shared/tenancy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const serviceName = "gateway"

// ServiceClients는 모든 gRPC 서비스 클라이언트를 보관합니다.
type ServiceClients struct {
	Auth         v1.AuthServiceClient
	User         v1.UserServiceClient
	Measurement  v1.MeasurementServiceClient
	Device       v1.DeviceServiceClient
	Subscription v1.SubscriptionServiceClient
	Shop         v1.ShopServiceClient
	Payment      v1.PaymentServiceClient
	AiInference  v1.AiInferenceServiceClient
	Cartridge    v1.CartridgeServiceClient
	Calibration  v1.CalibrationServiceClient
	Coaching     v1.CoachingServiceClient
	Reservation  v1.ReservationServiceClient
	Admin        v1.AdminServiceClient
	Family       v1.FamilyServiceClient
	HealthRecord v1.HealthRecordServiceClient
	Prescription v1.PrescriptionServiceClient
	Community    v1.CommunityServiceClient
	Video        v1.VideoServiceClient
	Notification v1.NotificationServiceClient
	Translation  v1.TranslationServiceClient
	Telemedicine v1.TelemedicineServiceClient

	// 신규 전용 서비스 (H6: nil이면 기존 서비스로 fallback)
	Audit       v1.AdminServiceClient       // audit-service → AdminService 위임
	DigitalTwin v1.MeasurementServiceClient // digital-twin-service → MeasurementService 위임

	// Phase B 서비스 (H6: nil이면 503 ServiceUnavailable)
	Assistant     v1.AssistantServiceClient
	Vision        v1.VisionServiceClient
	Concept       v1.ConceptServiceClient
	Organization  v1.OrganizationServiceClient
	DataPlatform  v1.LocationStatsServiceClient
	DataProvision v1.DataProvisionServiceClient
	VoiceProfile  v1.VoiceProfileServiceClient
}

func main() {
	cfg := config.LoadFromEnv(serviceName)
	httpPort := cfg.HTTPPort
	if httpPort == "" || httpPort == ":8080" {
		httpPort = ":8080"
	}

	log.Printf("[%s] Starting REST Gateway v%s on %s...", serviceName, cfg.Version, httpPort)

	// gRPC 연결 옵션
	dialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	// 각 서비스에 gRPC 연결
	conns, clients := connectServices(dialOpts)
	defer func() {
		for _, conn := range conns {
			conn.Close()
		}
	}()

	// REST 핸들러 초기화
	restHandler := handler.NewRestHandler(clients.Auth, clients.User, clients.Measurement,
		clients.Device, clients.Subscription, clients.Shop, clients.Payment,
		clients.AiInference, clients.Cartridge, clients.Calibration, clients.Coaching,
		clients.Reservation, clients.Admin, clients.Family, clients.HealthRecord,
		clients.Prescription, clients.Community, clients.Video, clients.Notification,
		clients.Translation, clients.Telemedicine, cfg.JWT.Secret)

	// 신규 서비스 클라이언트 등록 (H6: nil이면 기존 서비스로 fallback)
	if clients.Audit != nil {
		restHandler.SetAuditClient(clients.Audit)
	}
	if auditIntakeURL := os.Getenv("AUDIT_INTAKE_URL"); auditIntakeURL != "" {
		restHandler.SetAuditEventRecorder(handler.NewHTTPAuditEventRecorder(auditIntakeURL, nil))
	}
	if clients.DigitalTwin != nil {
		restHandler.SetDigitalTwinClient(clients.DigitalTwin)
	}

	// Phase B 서비스 클라이언트 등록
	if clients.Assistant != nil {
		restHandler.SetAssistantClient(clients.Assistant)
	}
	if clients.Vision != nil {
		restHandler.SetVisionClient(clients.Vision)
	}
	if clients.Concept != nil {
		restHandler.SetConceptClients(clients.Concept, clients.Organization)
	}
	if clients.DataPlatform != nil {
		restHandler.SetDataPlatformClients(clients.DataPlatform, clients.DataProvision)
	}
	if clients.VoiceProfile != nil {
		restHandler.SetVoiceProfileClient(clients.VoiceProfile)
	}

	// 라우터 설정
	mux := restHandler.SetupRoutes()

	// 멀티테넌트 멤버십 + 초대 REST API (Phase Z + AA).
	// DB_HOST 설정 시 PostgreSQL store 사용 (영속화), 미설정 시 인메모리 fallback.
	tenancyMemStore, tenancyInvStore, tenancyPool := buildTenancyStores(cfg, serviceName)
	if tenancyPool != nil {
		defer tenancyPool.Close()
	}
	tenancyEngine := tenancy.NewPolicyEngine(tenancyMemStore)
	tenancyInvSvc, _ := tenancy.NewInvitationService(tenancyInvStore, tenancyMemStore, tenancy.InvitationServiceConfig{})
	tenancyInvSvc.SetPolicyEngine(tenancyEngine)

	// Webhook dispatcher (Phase AM-1) — WEBHOOK_URL 환경변수 시 활성
	webhookDispatcher := buildGatewayWebhook(serviceName)
	if webhookDispatcher != nil {
		tenancyInvSvc.SetWebhookDispatcher(webhookDispatcher)
		webhookDispatcher.Start(context.Background())
		defer webhookDispatcher.Stop()
	}

	tenancyHTTP := tenancy.NewHTTPHandler(tenancyInvSvc, tenancyMemStore, tenancyEngine)
	tenancyHTTP.SetPathPrefix("/api/v1")
	tenancyHTTP.RegisterRoutes(mux)
	log.Printf("[%s] Tenancy REST API mounted on /api/v1/tenancy/*", serviceName)

	// 미들웨어 체인: 보안헤더 → CORS → Rate Limit → 바디크기제한 → 로깅 → 테넌시 전파
	// (테넌시는 가장 안쪽 = 핸들러 직전에 두어, 라우트 분기 후 ctx 에 안전히 보관)
	finalHandler := gw.SecurityHeaders(gw.CORS(gw.RateLimit(gw.MaxBodySize(10 << 20)(gw.Logging(gw.TenantPropagation(mux))))))

	server := &http.Server{
		Addr:         httpPort,
		Handler:      finalHandler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigCh
		log.Printf("[%s] Received signal %v, shutting down...", serviceName, sig)

		ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("[%s] Shutdown error: %v", serviceName, err)
		}
	}()

	log.Printf("[%s] REST Gateway listening on %s", serviceName, httpPort)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("[%s] Failed to serve: %v", serviceName, err)
	}
	log.Printf("[%s] Shutdown complete", serviceName)
}

// buildGatewayWebhook 는 WEBHOOK_URL 환경변수 시 dispatcher 생성 (Phase AM-1).
//
// gateway 와 admin-service 모두 동일 invitation 흐름을 처리하므로 같은 webhook
// 발송 로직을 공유. 둘 다 활성 시 동일 이벤트가 두 번 발송될 수 있으므로 운영에서는
// 하나만 활성화 권장.
func buildGatewayWebhook(serviceName string) *tenancy.WebhookDispatcher {
	url := os.Getenv("WEBHOOK_URL")
	if url == "" {
		return nil
	}
	mode := os.Getenv("WEBHOOK_MODE")
	if mode == "" {
		mode = "generic"
	}
	cfg := tenancy.WebhookConfig{
		URL:        url,
		Mode:       mode,
		MaxRetries: 3,
		OnError: func(_ tenancy.Event, err error) {
			log.Printf("[%s] webhook 에러: %v", serviceName, err)
		},
	}
	d, err := tenancy.NewWebhookDispatcher(cfg, nil)
	if err != nil {
		log.Printf("[%s] WebhookDispatcher 생성 실패: %v", serviceName, err)
		return nil
	}
	log.Printf("[%s] Webhook 활성 (mode=%s)", serviceName, mode)
	return d
}

// buildTenancyStores 는 DB_HOST 설정 여부에 따라 PostgreSQL/메모리 store 선택.
//
// 반환된 pool 은 호출자가 defer Close 책임. 미연결 시 nil 반환.
func buildTenancyStores(cfg *config.ServiceConfig, serviceName string) (
	tenancy.MembershipStore, tenancy.InvitationStore, *pgxpool.Pool) {
	dbHostSet := false
	if _, ok := os.LookupEnv("DB_HOST"); ok {
		dbHostSet = true
	}
	if !dbHostSet || cfg.DB.Host == "" || cfg.DB.DBName == "" {
		log.Printf("[%s] Tenancy 인메모리 store 사용 (DB_HOST 미설정)", serviceName)
		return tenancy.NewMemoryMembershipStore(),
			tenancy.NewMemoryInvitationStore(),
			nil
	}

	connCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	pool, err := pgxpool.New(connCtx, cfg.DB.DSN())
	cancel()
	if err != nil {
		log.Printf("[%s] Tenancy DB 풀 생성 실패, 인메모리 fallback: %v", serviceName, err)
		return tenancy.NewMemoryMembershipStore(),
			tenancy.NewMemoryInvitationStore(),
			nil
	}
	pingCtx, pingCancel := context.WithTimeout(context.Background(), 3*time.Second)
	if err := pool.Ping(pingCtx); err != nil {
		pingCancel()
		pool.Close()
		log.Printf("[%s] Tenancy DB Ping 실패, 인메모리 fallback: %v", serviceName, err)
		return tenancy.NewMemoryMembershipStore(),
			tenancy.NewMemoryInvitationStore(),
			nil
	}
	pingCancel()

	adapter := tenancy.NewPgxAdapter(pool)
	memStore, err := tenancy.NewPostgresMembershipStore(adapter)
	if err != nil {
		pool.Close()
		log.Printf("[%s] PostgresMembershipStore 생성 실패, 인메모리 fallback: %v", serviceName, err)
		return tenancy.NewMemoryMembershipStore(),
			tenancy.NewMemoryInvitationStore(),
			nil
	}
	invStore, err := tenancy.NewPostgresInvitationStore(adapter)
	if err != nil {
		pool.Close()
		log.Printf("[%s] PostgresInvitationStore 생성 실패: %v", serviceName, err)
		return memStore, tenancy.NewMemoryInvitationStore(), nil
	}
	log.Printf("[%s] Tenancy PostgreSQL store 활성화 (db=%s)", serviceName, cfg.DB.DBName)
	return memStore, invStore, pool
}

// connectServices는 모든 gRPC 서비스에 연결합니다.
func connectServices(opts []grpc.DialOption) ([]*grpc.ClientConn, *ServiceClients) {
	type svcInfo struct {
		name string
		port string
	}

	services := []svcInfo{
		{"auth", getEnv("AUTH_SERVICE_ADDR", "localhost:50051")},
		{"user", getEnv("USER_SERVICE_ADDR", "localhost:50052")},
		{"device", getEnv("DEVICE_SERVICE_ADDR", "localhost:50053")},
		{"measurement", getEnv("MEASUREMENT_SERVICE_ADDR", "localhost:50054")},
		{"subscription", getEnv("SUBSCRIPTION_SERVICE_ADDR", "localhost:50055")},
		{"shop", getEnv("SHOP_SERVICE_ADDR", "localhost:50056")},
		{"payment", getEnv("PAYMENT_SERVICE_ADDR", "localhost:50057")},
		{"ai-inference", getEnv("AI_INFERENCE_SERVICE_ADDR", "localhost:50058")},
		{"cartridge", getEnv("CARTRIDGE_SERVICE_ADDR", "localhost:50059")},
		{"calibration", getEnv("CALIBRATION_SERVICE_ADDR", "localhost:50060")},
		{"coaching", getEnv("COACHING_SERVICE_ADDR", "localhost:50061")},
		{"notification", getEnv("NOTIFICATION_SERVICE_ADDR", "localhost:50062")},
		{"family", getEnv("FAMILY_SERVICE_ADDR", "localhost:50063")},
		{"health-record", getEnv("HEALTH_RECORD_SERVICE_ADDR", "localhost:50064")},
		{"telemedicine", getEnv("TELEMEDICINE_SERVICE_ADDR", "localhost:50065")},
		{"reservation", getEnv("RESERVATION_SERVICE_ADDR", "localhost:50066")},
		{"community", getEnv("COMMUNITY_SERVICE_ADDR", "localhost:50067")},
		{"admin", getEnv("ADMIN_SERVICE_ADDR", "localhost:50068")},
		{"prescription", getEnv("PRESCRIPTION_SERVICE_ADDR", "localhost:50069")},
		{"translation", getEnv("TRANSLATION_SERVICE_ADDR", "localhost:50070")},
		{"video", getEnv("VIDEO_SERVICE_ADDR", "localhost:50071")},
		{"audit", getEnv("AUDIT_SERVICE_ADDR", "localhost:50072")},
		{"digital-twin", getEnv("DIGITAL_TWIN_SERVICE_ADDR", "localhost:50073")},
		{"assistant", getEnv("ASSISTANT_SERVICE_ADDR", "localhost:50074")},
		{"vision", getEnv("VISION_SERVICE_ADDR", "localhost:50075")},
		{"concept", getEnv("CONCEPT_SERVICE_ADDR", "localhost:50076")},
		{"data-platform", getEnv("DATA_PLATFORM_SERVICE_ADDR", "localhost:50077")},
		{"voice-profile", getEnv("VOICE_PROFILE_SERVICE_ADDR", "localhost:50078")},
	}

	connMap := make(map[string]*grpc.ClientConn)
	var conns []*grpc.ClientConn

	for _, svc := range services {
		conn, err := grpc.NewClient(svc.port, opts...)
		if err != nil {
			log.Printf("[%s] Warning: failed to connect to %s at %s: %v", serviceName, svc.name, svc.port, err)
			continue
		}
		connMap[svc.name] = conn
		conns = append(conns, conn)
		log.Printf("[%s] Connected to %s at %s", serviceName, svc.name, svc.port)
	}

	clients := &ServiceClients{}
	if c, ok := connMap["auth"]; ok {
		clients.Auth = v1.NewAuthServiceClient(c)
	}
	if c, ok := connMap["user"]; ok {
		clients.User = v1.NewUserServiceClient(c)
	}
	if c, ok := connMap["measurement"]; ok {
		clients.Measurement = v1.NewMeasurementServiceClient(c)
	}
	if c, ok := connMap["device"]; ok {
		clients.Device = v1.NewDeviceServiceClient(c)
	}
	if c, ok := connMap["subscription"]; ok {
		clients.Subscription = v1.NewSubscriptionServiceClient(c)
	}
	if c, ok := connMap["shop"]; ok {
		clients.Shop = v1.NewShopServiceClient(c)
	}
	if c, ok := connMap["payment"]; ok {
		clients.Payment = v1.NewPaymentServiceClient(c)
	}
	if c, ok := connMap["ai-inference"]; ok {
		clients.AiInference = v1.NewAiInferenceServiceClient(c)
	}
	if c, ok := connMap["cartridge"]; ok {
		clients.Cartridge = v1.NewCartridgeServiceClient(c)
	}
	if c, ok := connMap["calibration"]; ok {
		clients.Calibration = v1.NewCalibrationServiceClient(c)
	}
	if c, ok := connMap["coaching"]; ok {
		clients.Coaching = v1.NewCoachingServiceClient(c)
	}
	if c, ok := connMap["reservation"]; ok {
		clients.Reservation = v1.NewReservationServiceClient(c)
	}
	if c, ok := connMap["admin"]; ok {
		clients.Admin = v1.NewAdminServiceClient(c)
	}
	if c, ok := connMap["family"]; ok {
		clients.Family = v1.NewFamilyServiceClient(c)
	}
	if c, ok := connMap["health-record"]; ok {
		clients.HealthRecord = v1.NewHealthRecordServiceClient(c)
	}
	if c, ok := connMap["prescription"]; ok {
		clients.Prescription = v1.NewPrescriptionServiceClient(c)
	}
	if c, ok := connMap["community"]; ok {
		clients.Community = v1.NewCommunityServiceClient(c)
	}
	if c, ok := connMap["video"]; ok {
		clients.Video = v1.NewVideoServiceClient(c)
	}
	if c, ok := connMap["notification"]; ok {
		clients.Notification = v1.NewNotificationServiceClient(c)
	}
	if c, ok := connMap["translation"]; ok {
		clients.Translation = v1.NewTranslationServiceClient(c)
	}
	if c, ok := connMap["telemedicine"]; ok {
		clients.Telemedicine = v1.NewTelemedicineServiceClient(c)
	}

	// 신규 서비스: audit-service + digital-twin-service (H6 실패격리)
	if c, ok := connMap["audit"]; ok {
		clients.Audit = v1.NewAdminServiceClient(c)
	}
	if c, ok := connMap["digital-twin"]; ok {
		clients.DigitalTwin = v1.NewMeasurementServiceClient(c)
	}

	// Phase B 서비스 연결 (H6: nil이면 503 반환)
	if c, ok := connMap["assistant"]; ok {
		clients.Assistant = v1.NewAssistantServiceClient(c)
	}
	if c, ok := connMap["vision"]; ok {
		clients.Vision = v1.NewVisionServiceClient(c)
	}
	if c, ok := connMap["concept"]; ok {
		clients.Concept = v1.NewConceptServiceClient(c)
		clients.Organization = v1.NewOrganizationServiceClient(c)
	}
	if c, ok := connMap["data-platform"]; ok {
		clients.DataPlatform = v1.NewLocationStatsServiceClient(c)
		clients.DataProvision = v1.NewDataProvisionServiceClient(c)
	}
	if c, ok := connMap["voice-profile"]; ok {
		clients.VoiceProfile = v1.NewVoiceProfileServiceClient(c)
	}

	connected := 0
	for _, c := range conns {
		_ = c
		connected++
	}
	log.Printf("[%s] Connected to %d/%d services", serviceName, connected, len(services))

	return conns, clients
}

func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

// init은 fmt 패키지 사용을 보장합니다.
var _ = fmt.Sprintf
