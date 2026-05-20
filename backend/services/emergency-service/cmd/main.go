// emergency-service: 응급 상황 관리 마이크로서비스
//
// 포트: HTTP :8080 (Health + REST API)
// 의존: 인메모리 저장소
//
// 기능:
// - 응급 상황 신고 (ReportEmergency)
// - 비상 연락처 조회 (GetEmergencyContacts)
// - 응급 설정 관리 (Get/UpdateEmergencySettings)
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manpasik/backend/services/emergency-service/internal/handler"
	"github.com/manpasik/backend/services/emergency-service/internal/notifier"
	"github.com/manpasik/backend/services/emergency-service/internal/repository/memory"
	"github.com/manpasik/backend/services/emergency-service/internal/repository/postgres"
	"github.com/manpasik/backend/services/emergency-service/internal/service"
	"github.com/manpasik/backend/shared/config"
)

const serviceName = "emergency-service"

func main() {
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = ":8080"
	}

	cfg := config.LoadFromEnv(serviceName)
	log.Printf("[%s] Starting v%s...", serviceName, cfg.Version)

	var repo service.EmergencyRepository

	if _, dbHostSet := os.LookupEnv("DB_HOST"); dbHostSet && cfg.DB.Host != "" && cfg.DB.DBName != "" {
		connCtx, connCancel := context.WithTimeout(context.Background(), 5*time.Second)
		pool, poolErr := pgxpool.New(connCtx, cfg.DB.DSN())
		connCancel()
		if poolErr != nil {
			log.Printf("[%s] DB connection failed, using memory: %v", serviceName, poolErr)
			repo = memory.NewEmergencyRepository()
		} else {
			pingCtx, pingCancel := context.WithTimeout(context.Background(), 3*time.Second)
			if pingErr := pool.Ping(pingCtx); pingErr != nil {
				pingCancel()
				pool.Close()
				log.Printf("[%s] DB ping failed, using memory: %v", serviceName, pingErr)
				repo = memory.NewEmergencyRepository()
			} else {
				pingCancel()
				defer pool.Close()
				log.Printf("[%s] Connected to PostgreSQL", serviceName)
				repo = postgres.NewEmergencyRepository(pool)
			}
		}
	} else {
		repo = memory.NewEmergencyRepository()
	}

	svc := service.NewEmergencyService(repo)

	// 119 응급 알림 제공자 초기화
	if apiKey := os.Getenv("EMERGENCY_119_API_KEY"); apiKey != "" {
		baseURL := os.Getenv("EMERGENCY_119_BASE_URL")
		n := notifier.NewDispatch119Notifier(notifier.Dispatch119Config{
			BaseURL: baseURL,
			APIKey:  apiKey,
		})
		svc.SetEmergencyNotifier(n)
		log.Printf("[%s] 119 응급 알림 제공자 활성화", serviceName)
	} else {
		log.Printf("[%s] 119 응급 알림 미설정, 로그만 기록", serviceName)
	}

	// HTTP 핸들러 초기화
	h := handler.NewEmergencyHandler(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"serving","service":"emergency-service"}`))
	})
	h.RegisterRoutes(mux)

	go func() {
		log.Printf("[%s] HTTP server on %s", serviceName, httpPort)
		if err := http.ListenAndServe(httpPort, mux); err != nil {
			log.Fatalf("[%s] HTTP server error: %v", serviceName, err)
		}
	}()

	// 시그널 대기
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigCh
	log.Printf("[%s] Received signal %v, shutting down...", serviceName, sig)
	log.Printf("[%s] Shutdown complete", serviceName)
}
