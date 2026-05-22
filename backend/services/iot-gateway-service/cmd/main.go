// iot-gateway-service: IoT 디바이스 게이트웨이 마이크로서비스
//
// 포트: gRPC :50072
// 의존: PostgreSQL(선택) — 미설정 시 인메모리 저장소 사용
//
// 기능:
// - IoT 디바이스 등록 / 조회
// - 디바이스 명령 전송
// - 디바이스 텔레메트리 수신/조회
// - 디바이스 하트비트
// - gRPC IoTGatewayService
package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manpasik/backend/services/iot-gateway-service/internal/handler"
	"github.com/manpasik/backend/services/iot-gateway-service/internal/repository/memory"
	"github.com/manpasik/backend/services/iot-gateway-service/internal/repository/postgres"
	"github.com/manpasik/backend/services/iot-gateway-service/internal/service"
	"github.com/manpasik/backend/shared/config"
	v1 "github.com/manpasik/backend/shared/gen/go/v1"
	"github.com/manpasik/backend/shared/middleware"
	"github.com/manpasik/backend/shared/observability"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

const serviceName = "iot-gateway-service"

func main() {
	cfg := config.LoadFromEnv(serviceName)
	if cfg.GRPCPort == "" || cfg.GRPCPort == ":50051" {
		cfg.GRPCPort = ":50072"
	}

	logger, err := zap.NewProduction()
	if err != nil {
		logger = zap.NewNop()
	}
	defer logger.Sync()

	metrics := observability.NewMetrics()
	healthCheck := observability.NewHealthCheck(serviceName, cfg.Version)

	log.Printf("[%s] Starting v%s...", serviceName, cfg.Version)
	log.Printf("[%s] gRPC port: %s", serviceName, cfg.GRPCPort)

	// IoTRepository: PostgreSQL 또는 인메모리
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
		log.Printf("[%s] 인메모리 저장소 사용", serviceName)
	}

	svc := service.NewIoTGatewayService(repo)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.RequestIDInterceptor(),
			observability.UnaryServerInterceptor(metrics),
		),
	)

	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus(serviceName, healthpb.HealthCheckResponse_SERVING)

	iotHandler := handler.NewIoTGatewayHandler(svc, logger)
	v1.RegisterIoTGatewayServiceServer(grpcServer, iotHandler)

	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", cfg.GRPCPort)
	if err != nil {
		log.Fatalf("[%s] Failed to listen: %v", serviceName, err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigCh
		log.Printf("[%s] Received signal %v, shutting down...", serviceName, sig)
		healthServer.SetServingStatus(serviceName, healthpb.HealthCheckResponse_NOT_SERVING)
		go func() {
			time.Sleep(cfg.ShutdownTimeout)
			os.Exit(1)
		}()
		grpcServer.GracefulStop()
		cancel()
	}()

	// Observability HTTP server
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/metrics", metrics.PrometheusHandler())
		mux.HandleFunc("/health", healthCheck.Handler())
		metricsAddr := ":9100"
		logger.Info("Metrics server starting", zap.String("addr", metricsAddr))
		if err := http.ListenAndServe(metricsAddr, mux); err != nil {
			logger.Error("Metrics server failed", zap.Error(err))
		}
	}()

	log.Printf("[%s] gRPC server listening on %s", serviceName, cfg.GRPCPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("[%s] Failed to serve: %v", serviceName, err)
	}
	<-ctx.Done()
	log.Printf("[%s] Shutdown complete", serviceName)
}
