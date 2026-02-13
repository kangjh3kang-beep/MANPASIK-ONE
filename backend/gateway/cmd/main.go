package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/manpasik/backend/gateway/internal/router"
	"github.com/manpasik/backend/shared/observability"
	"github.com/manpasik/backend/shared/storage"
)

const serviceName = "gateway"

var (
	Version   = "1.0.0"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

func main() {
	log.Printf("[%s] Starting REST-to-gRPC gateway...", serviceName)

	port := getEnvOrDefault("GATEWAY_PORT", "8080")

	cfg := router.Config{
		AuthAddr:         getEnvOrDefault("AUTH_SERVICE_ADDR", "localhost:50051"),
		MeasurementAddr:  getEnvOrDefault("MEASUREMENT_SERVICE_ADDR", "localhost:50054"),
		UserAddr:         getEnvOrDefault("USER_SERVICE_ADDR", "localhost:50052"),
		DeviceAddr:       getEnvOrDefault("DEVICE_SERVICE_ADDR", "localhost:50053"),
		ReservationAddr:  getEnvOrDefault("RESERVATION_SERVICE_ADDR", "localhost:50055"),
		PrescriptionAddr: getEnvOrDefault("PRESCRIPTION_SERVICE_ADDR", "localhost:50062"),
		SubscriptionAddr: getEnvOrDefault("SUBSCRIPTION_SERVICE_ADDR", "localhost:50055"),
		ShopAddr:         getEnvOrDefault("SHOP_SERVICE_ADDR", "localhost:50056"),
		PaymentAddr:      getEnvOrDefault("PAYMENT_SERVICE_ADDR", "localhost:50057"),
		HealthRecordAddr: getEnvOrDefault("HEALTH_RECORD_SERVICE_ADDR", "localhost:50064"),
		NotificationAddr: getEnvOrDefault("NOTIFICATION_SERVICE_ADDR", "localhost:50068"),
		CommunityAddr:    getEnvOrDefault("COMMUNITY_SERVICE_ADDR", "localhost:50065"),
		AdminAddr:        getEnvOrDefault("ADMIN_SERVICE_ADDR", "localhost:50067"),
		AiInferenceAddr:  getEnvOrDefault("AI_INFERENCE_SERVICE_ADDR", "localhost:50058"),
		CartridgeAddr:    getEnvOrDefault("CARTRIDGE_SERVICE_ADDR", "localhost:50059"),
		CalibrationAddr:  getEnvOrDefault("CALIBRATION_SERVICE_ADDR", "localhost:50060"),
		CoachingAddr:     getEnvOrDefault("COACHING_SERVICE_ADDR", "localhost:50061"),
	}

	metrics := observability.NewMetrics()
	healthCheck := observability.NewHealthCheck(serviceName, Version)

	r := router.NewRouter(cfg)

	// S3/MinIO 파일 스토리지 (선택)
	if s3Endpoint := os.Getenv("S3_ENDPOINT"); s3Endpoint != "" {
		useSSL, _ := strconv.ParseBool(getEnvOrDefault("S3_USE_SSL", "false"))
		s3Client, err := storage.NewS3Client(
			s3Endpoint,
			getEnvOrDefault("S3_ACCESS_KEY", "minioadmin"),
			getEnvOrDefault("S3_SECRET_KEY", "minioadmin"),
			getEnvOrDefault("S3_BUCKET", "manpasik"),
			getEnvOrDefault("S3_REGION", "us-east-1"),
			useSSL,
		)
		if err != nil {
			log.Printf("[%s] S3 연결 실패, 파일 업로드 비활성화: %v", serviceName, err)
		} else {
			r.SetS3Client(s3Client)
			log.Printf("[%s] S3/MinIO 연결됨: %s (버킷: %s)", serviceName, s3Endpoint, getEnvOrDefault("S3_BUCKET", "manpasik"))
		}
	}

	// Wrap router with observability routes
	mux := http.NewServeMux()
	mux.HandleFunc("/metrics", metrics.PrometheusHandler())
	mux.HandleFunc("/health/obs", healthCheck.Handler())
	mux.Handle("/", r)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      loggingMiddleware(mux),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		log.Printf("[%s] Shutting down...", serviceName)
		server.Close()
	}()

	log.Printf("[%s] HTTP server listening on :%s", serviceName, port)
	log.Printf("[%s] Endpoints:", serviceName)
	log.Printf("[%s]   /health, /health/live, /health/ready", serviceName)
	log.Printf("[%s]   /api/v1/auth/*           → Auth (%s)", serviceName, cfg.AuthAddr)
	log.Printf("[%s]   /api/v1/users/*          → User (%s)", serviceName, cfg.UserAddr)
	log.Printf("[%s]   /api/v1/devices/*        → Device (%s)", serviceName, cfg.DeviceAddr)
	log.Printf("[%s]   /api/v1/measurements/*   → Measurement (%s)", serviceName, cfg.MeasurementAddr)
	log.Printf("[%s]   /api/v1/facilities/*, /api/v1/reservations/* → Reservation (%s)", serviceName, cfg.ReservationAddr)
	log.Printf("[%s]   /api/v1/prescriptions/*  → Prescription (%s)", serviceName, cfg.PrescriptionAddr)
	log.Printf("[%s]   /api/v1/subscriptions/*  → Subscription (%s)", serviceName, cfg.SubscriptionAddr)
	log.Printf("[%s]   /api/v1/products/*, /api/v1/cart/*, /api/v1/orders/* → Shop (%s)", serviceName, cfg.ShopAddr)
	log.Printf("[%s]   /api/v1/payments/*       → Payment (%s)", serviceName, cfg.PaymentAddr)
	log.Printf("[%s]   /api/v1/health-records/* → HealthRecord (%s)", serviceName, cfg.HealthRecordAddr)
	log.Printf("[%s]   /api/v1/notifications/*  → Notification (%s)", serviceName, cfg.NotificationAddr)
	log.Printf("[%s]   /api/v1/posts/*          → Community (%s)", serviceName, cfg.CommunityAddr)
	log.Printf("[%s]   /api/v1/admin/*          → Admin (%s)", serviceName, cfg.AdminAddr)
	log.Printf("[%s]   /api/v1/ai/*            → AiInference (%s)", serviceName, cfg.AiInferenceAddr)
	log.Printf("[%s]   /api/v1/cartridges/*    → Cartridge (%s)", serviceName, cfg.CartridgeAddr)
	log.Printf("[%s]   /api/v1/calibration/*   → Calibration (%s)", serviceName, cfg.CalibrationAddr)
	log.Printf("[%s]   /api/v1/coaching/*       → Coaching (%s)", serviceName, cfg.CoachingAddr)

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server error: %v", err)
	}

	log.Printf("[%s] Shutdown complete", serviceName)
}

func getEnvOrDefault(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

// loggingMiddleware logs every HTTP request with method, path and duration.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("[%s] %s %s %s", serviceName, r.Method, r.URL.Path, time.Since(start))
	})
}
