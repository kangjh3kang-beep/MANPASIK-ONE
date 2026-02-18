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

	"github.com/manpasik/backend/services/gateway/internal/handler"
	gw "github.com/manpasik/backend/services/gateway/internal/middleware"
	"github.com/manpasik/backend/shared/config"
	v1 "github.com/manpasik/backend/shared/gen/go/v1"
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

	// 라우터 설정
	mux := restHandler.SetupRoutes()

	// CORS + 로깅 미들웨어 적용
	finalHandler := gw.CORS(gw.Logging(mux))

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
