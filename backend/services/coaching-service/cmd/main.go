// coaching-service: AI 건강 코칭 마이크로서비스
//
// 포트: gRPC :50061
// 의존: PostgreSQL(선택) — 미설정 시 인메모리 저장소 사용
//
// 기능:
// - 건강 목표 설정 / 조회
// - AI 코칭 메시지 생성 (측정 피드백, 일일 팁, 목표 진행, 경고, 동기부여, 추천)
// - 일일 건강 리포트 생성
// - 주간 건강 리포트 조회
// - 개인화 추천 조회
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
	"github.com/manpasik/backend/services/coaching-service/internal/handler"
	"github.com/manpasik/backend/services/coaching-service/internal/repository/memory"
	"github.com/manpasik/backend/services/coaching-service/internal/repository/postgres"
	"github.com/manpasik/backend/services/coaching-service/internal/service"
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

const serviceName = "coaching-service"

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

	// Repository: PostgreSQL 또는 인메모리
	var goalRepo service.HealthGoalRepository
	var msgRepo service.CoachingMessageRepository
	var reportRepo service.DailyReportRepository

	if _, dbHostSet := os.LookupEnv("DB_HOST"); dbHostSet && cfg.DB.Host != "" && cfg.DB.DBName != "" {
		connCtx, connCancel := context.WithTimeout(context.Background(), 5*time.Second)
		pool, poolErr := pgxpool.New(connCtx, cfg.DB.DSN())
		connCancel()
		if poolErr != nil {
			log.Printf("[%s] DB 풀 생성 실패, 인메모리 사용: %v", serviceName, poolErr)
			goalRepo = memory.NewHealthGoalRepository()
			msgRepo = memory.NewCoachingMessageRepository()
			reportRepo = memory.NewDailyReportRepository()
		} else {
			pingCtx, pingCancel := context.WithTimeout(context.Background(), 3*time.Second)
			if pingErr := pool.Ping(pingCtx); pingErr != nil {
				pingCancel()
				pool.Close()
				log.Printf("[%s] DB Ping 실패, 인메모리 사용: %v", serviceName, pingErr)
				goalRepo = memory.NewHealthGoalRepository()
				msgRepo = memory.NewCoachingMessageRepository()
				reportRepo = memory.NewDailyReportRepository()
			} else {
				pingCancel()
				defer pool.Close()
				goalRepo = postgres.NewHealthGoalRepository(pool)
				msgRepo = postgres.NewCoachingMessageRepository(pool)
				reportRepo = postgres.NewDailyReportRepository(pool)
				log.Printf("[%s] DB 연결됨: %s", serviceName, cfg.DB.DBName)
			}
		}
	} else {
		goalRepo = memory.NewHealthGoalRepository()
		msgRepo = memory.NewCoachingMessageRepository()
		reportRepo = memory.NewDailyReportRepository()
		log.Printf("[%s] 인메모리 저장소 사용", serviceName)
	}

	coachingSvc := service.NewCoachingService(logger, goalRepo, msgRepo, reportRepo)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.RequestIDInterceptor(),
			observability.UnaryServerInterceptor(metrics),
		),
	)

	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus(serviceName, healthpb.HealthCheckResponse_SERVING)

	coachingHandler := handler.NewCoachingHandler(coachingSvc, logger)
	v1.RegisterCoachingServiceServer(grpcServer, coachingHandler)

	reflection.Register(grpcServer)

	grpcPort := cfg.GRPCPort
	if grpcPort == "" {
		grpcPort = ":50061"
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

	log.Printf("[%s] gRPC server listening on %s", serviceName, grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("[%s] Failed to serve: %v", serviceName, err)
	}
	<-ctx.Done()
	log.Printf("[%s] Shutdown complete", serviceName)
}
