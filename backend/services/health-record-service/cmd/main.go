// health-record-service: 건강 기록 관리 마이크로서비스
//
// 포트: gRPC :50064
// 의존: PostgreSQL(선택) — 미설정 시 인메모리 저장소 사용
//
// 기능:
// - 건강 기록 CRUD (측정, 복약, 증상, 활력징후, 검사결과, 알레르기, 질환, 예방접종, 시술, 자유 기록)
// - FHIR R4 형식 내보내기/가져오기
// - 건강 요약 조회
package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manpasik/backend/services/health-record-service/internal/handler"
	"github.com/manpasik/backend/services/health-record-service/internal/repository/memory"
	"github.com/manpasik/backend/services/health-record-service/internal/repository/postgres"
	"github.com/manpasik/backend/services/health-record-service/internal/service"
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

const serviceName = "health-record-service"

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
	log.Printf("[%s] gRPC port: %s", serviceName, cfg.GRPCPort)

	var recordRepo service.HealthRecordRepository
	var consentRepo service.ConsentRepository
	var accessLogRepo service.DataAccessLogRepository

	if _, dbHostSet := os.LookupEnv("DB_HOST"); dbHostSet && cfg.DB.Host != "" && cfg.DB.DBName != "" {
		connCtx, connCancel := context.WithTimeout(context.Background(), 5*time.Second)
		pool, err := pgxpool.New(connCtx, cfg.DB.DSN())
		connCancel()
		if err != nil {
			log.Printf("[%s] DB 풀 생성 실패, 인메모리 사용: %v", serviceName, err)
			recordRepo = memory.NewHealthRecordRepository()
			consentRepo = memory.NewConsentRepository()
			accessLogRepo = memory.NewDataAccessLogRepository()
		} else {
			pingCtx, pingCancel := context.WithTimeout(context.Background(), 3*time.Second)
			if pingErr := pool.Ping(pingCtx); pingErr != nil {
				pingCancel()
				pool.Close()
				log.Printf("[%s] DB Ping 실패, 인메모리 사용: %v", serviceName, pingErr)
				recordRepo = memory.NewHealthRecordRepository()
				consentRepo = memory.NewConsentRepository()
				accessLogRepo = memory.NewDataAccessLogRepository()
			} else {
				pingCancel()
				defer pool.Close()
				recordRepo = postgres.NewHealthRecordRepository(pool)
				consentRepo = postgres.NewConsentRepository(pool)
				accessLogRepo = postgres.NewDataAccessLogRepository(pool)
				log.Printf("[%s] DB 연결됨: %s", serviceName, cfg.DB.DBName)
			}
		}
	} else {
		recordRepo = memory.NewHealthRecordRepository()
		consentRepo = memory.NewConsentRepository()
		accessLogRepo = memory.NewDataAccessLogRepository()
		log.Printf("[%s] 인메모리 저장소 사용", serviceName)
	}

	hrSvc := service.NewHealthRecordService(logger, recordRepo, consentRepo, accessLogRepo)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.RequestIDInterceptor(),
			observability.UnaryServerInterceptor(metrics),
		),
	)

	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus(serviceName, healthpb.HealthCheckResponse_SERVING)

	hrHandler := handler.NewHealthRecordHandler(hrSvc, logger)
	v1.RegisterHealthRecordServiceServer(grpcServer, hrHandler)

	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", ":50064")
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

	// Start observability HTTP server
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

	log.Printf("[%s] gRPC server listening on :50064", serviceName)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("[%s] Failed to serve: %v", serviceName, err)
	}
	<-ctx.Done()
	log.Printf("[%s] Shutdown complete", serviceName)
}
