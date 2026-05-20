// audit-service: 감사 로그 전담 마이크로서비스
//
// 포트: gRPC :50072
// 의존: PostgreSQL(선택) — 미설정 시 인메모리 저장소 사용
//
// 기능:
// - 감사 로그 수집 (Kafka consumer로 이벤트 수신)
// - 감사 로그 저장 (PostgreSQL)
// - 감사 로그 검색/필터링 (시간 범위, 사용자, 액션 타입)
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
	"github.com/manpasik/backend/services/audit-service/internal/handler"
	"github.com/manpasik/backend/services/audit-service/internal/repository/memory"
	"github.com/manpasik/backend/services/audit-service/internal/repository/postgres"
	"github.com/manpasik/backend/services/audit-service/internal/service"
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

const serviceName = "audit-service"

func main() {
	cfg := config.LoadFromEnv(serviceName)
	logger, err := zap.NewProduction()
	if err != nil {
		logger = zap.NewNop()
	}
	defer logger.Sync()

	metrics := observability.NewMetrics()
	healthCheck := observability.NewHealthCheck(serviceName, cfg.Version)
	log.Printf("[%s] Starting v%s...", serviceName, cfg.Version)

	var repo service.AuditRepository

	if _, dbHostSet := os.LookupEnv("DB_HOST"); dbHostSet && cfg.DB.Host != "" && cfg.DB.DBName != "" {
		connCtx, connCancel := context.WithTimeout(context.Background(), 5*time.Second)
		pool, poolErr := pgxpool.New(connCtx, cfg.DB.DSN())
		connCancel()
		if poolErr != nil {
			log.Printf("[%s] DB connection failed, using memory: %v", serviceName, poolErr)
			memRepo := memory.NewAuditRepository()
			memRepo.Seed()
			repo = memRepo
		} else {
			pingCtx, pingCancel := context.WithTimeout(context.Background(), 3*time.Second)
			if pingErr := pool.Ping(pingCtx); pingErr != nil {
				pingCancel()
				pool.Close()
				log.Printf("[%s] DB ping failed, using memory: %v", serviceName, pingErr)
				memRepo := memory.NewAuditRepository()
				memRepo.Seed()
				repo = memRepo
			} else {
				pingCancel()
				defer pool.Close()
				log.Printf("[%s] Connected to PostgreSQL", serviceName)
				repo = postgres.NewAuditRepository(pool)
			}
		}
	} else {
		memRepo := memory.NewAuditRepository()
		memRepo.Seed()
		repo = memRepo
	}

	svc := service.NewAuditService(logger, repo)
	h := handler.NewAuditHandler(svc, logger)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.RequestIDInterceptor(),
			observability.UnaryServerInterceptor(metrics),
		),
	)
	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus(serviceName, healthpb.HealthCheckResponse_SERVING)
	v1.RegisterAdminServiceServer(grpcServer, h)
	reflection.Register(grpcServer)

	grpcPort := cfg.GRPCPort
	if grpcPort == "" || grpcPort == ":50051" {
		grpcPort = ":50072"
	}
	lis, err := net.Listen("tcp", grpcPort)
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
		go func() { time.Sleep(cfg.ShutdownTimeout); os.Exit(1) }()
		grpcServer.GracefulStop()
		cancel()
	}()
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/metrics", metrics.PrometheusHandler())
		mux.HandleFunc("/health", healthCheck.Handler())
		handler.NewHTTPHandler(svc, logger).RegisterRoutes(mux)
		logger.Info("Metrics server starting", zap.String("addr", ":9100"))
		http.ListenAndServe(":9100", mux)
	}()
	log.Printf("[%s] gRPC server listening on %s", serviceName, grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("[%s] Failed to serve: %v", serviceName, err)
	}
	<-ctx.Done()
}
