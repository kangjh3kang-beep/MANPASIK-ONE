// iot-gateway-service: IoT 디바이스 게이트웨이 마이크로서비스
//
// 포트: HTTP :8080
// 의존: 없음 — 인메모리 저장소 사용
//
// 기능:
// - IoT 디바이스 등록 / 조회
// - 디바이스 명령 전송
// - 디바이스 데이터 수신
// - 헬스 체크 엔드포인트
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
	"github.com/manpasik/backend/services/iot-gateway-service/internal/repository/memory"
	"github.com/manpasik/backend/services/iot-gateway-service/internal/repository/postgres"
	"github.com/manpasik/backend/shared/config"
	"github.com/manpasik/backend/services/iot-gateway-service/internal/service"
)

const serviceName = "iot-gateway-service"

func main() {
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = ":8080"
	}

	cfg := config.LoadFromEnv(serviceName)
	log.Printf("[%s] Starting v%s...", serviceName, cfg.Version)

	var repo service.IoTRepository

	if _, dbHostSet := os.LookupEnv("DB_HOST"); dbHostSet && cfg.DB.Host != "" && cfg.DB.DBName != "" {
		connCtx, connCancel := context.WithTimeout(context.Background(), 5*time.Second)
		pool, poolErr := pgxpool.New(connCtx, cfg.DB.DSN())
		connCancel()
		if poolErr != nil {
			log.Printf("[%s] DB connection failed, using memory: %v", serviceName, poolErr)
			repo = memory.NewIoTRepository()
		} else {
			pingCtx, pingCancel := context.WithTimeout(context.Background(), 3*time.Second)
			if pingErr := pool.Ping(pingCtx); pingErr != nil {
				pingCancel()
				pool.Close()
				log.Printf("[%s] DB ping failed, using memory: %v", serviceName, pingErr)
				repo = memory.NewIoTRepository()
			} else {
				pingCancel()
				defer pool.Close()
				log.Printf("[%s] Connected to PostgreSQL", serviceName)
				repo = postgres.NewIoTRepository(pool)
			}
		}
	} else {
		repo = memory.NewIoTRepository()
	}

	svc := service.NewIoTGatewayService(repo)
	_ = svc

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"serving","service":"iot-gateway-service"}`))
	})

	go func() {
		log.Printf("[%s] HTTP server on %s", serviceName, httpPort)
		if err := http.ListenAndServe(httpPort, mux); err != nil {
			log.Fatalf("[%s] HTTP server error: %v", serviceName, err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigCh
	log.Printf("[%s] Received signal %v, shutting down...", serviceName, sig)
}
