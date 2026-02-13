package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manpasik/backend/services/ai-inference-service/internal/handler"
	"github.com/manpasik/backend/services/ai-inference-service/internal/llm"
	"github.com/manpasik/backend/services/ai-inference-service/internal/repository/memory"
	"github.com/manpasik/backend/services/ai-inference-service/internal/repository/postgres"
	"github.com/manpasik/backend/services/ai-inference-service/internal/service"
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

const serviceName = "ai-inference-service"

func main() {
	cfg := config.LoadFromEnv(serviceName)
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	metrics := observability.NewMetrics()
	healthCheck := observability.NewHealthCheck(serviceName, "1.0.0")

	// 포트 결정
	port := cfg.GRPCPort
	if port == "" {
		port = ":50058"
	}

	// Repository: PostgreSQL 또는 인메모리
	var analysisRepo service.AnalysisRepository
	var healthScoreRepo service.HealthScoreRepository

	if _, dbHostSet := os.LookupEnv("DB_HOST"); dbHostSet && cfg.DB.Host != "" && cfg.DB.DBName != "" {
		connCtx, connCancel := context.WithTimeout(context.Background(), 5*time.Second)
		pool, poolErr := pgxpool.New(connCtx, cfg.DB.DSN())
		connCancel()
		if poolErr != nil {
			log.Printf("[%s] DB 풀 생성 실패, 인메모리 사용: %v", serviceName, poolErr)
			analysisRepo = memory.NewAnalysisRepository()
			healthScoreRepo = memory.NewHealthScoreRepository()
		} else {
			pingCtx, pingCancel := context.WithTimeout(context.Background(), 3*time.Second)
			if pingErr := pool.Ping(pingCtx); pingErr != nil {
				pingCancel()
				pool.Close()
				log.Printf("[%s] DB Ping 실패, 인메모리 사용: %v", serviceName, pingErr)
				analysisRepo = memory.NewAnalysisRepository()
				healthScoreRepo = memory.NewHealthScoreRepository()
			} else {
				pingCancel()
				defer pool.Close()
				analysisRepo = postgres.NewAnalysisRepository(pool)
				healthScoreRepo = postgres.NewHealthScoreRepository(pool)
				log.Printf("[%s] DB 연결됨: %s", serviceName, cfg.DB.DBName)
			}
		}
	} else {
		analysisRepo = memory.NewAnalysisRepository()
		healthScoreRepo = memory.NewHealthScoreRepository()
		log.Printf("[%s] 인메모리 저장소 사용", serviceName)
	}

	// LLM 클라이언트 초기화 (환경변수 LLM_API_KEY가 설정된 경우에만)
	var svcOpts []service.InferenceOption
	if apiKey, ok := os.LookupEnv("LLM_API_KEY"); ok && apiKey != "" {
		llmClient := llm.NewOpenAIClientFromEnv()
		svcOpts = append(svcOpts, service.WithLLMClient(llmClient))
		log.Printf("[%s] LLM 클라이언트 활성화 (모델: %s)", serviceName, llmClient.Model())
	} else {
		log.Printf("[%s] LLM_API_KEY 미설정 — LLM 기능 비활성화", serviceName)
	}

	// Service
	svc := service.NewInferenceService(analysisRepo, healthScoreRepo, svcOpts...)

	// gRPC Handler
	h := handler.NewInferenceHandler(svc)

	// gRPC Server
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.RequestIDInterceptor(),
			observability.UnaryServerInterceptor(metrics),
		),
	)
	v1.RegisterAiInferenceServiceServer(grpcServer, h)

	// Health check
	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("manpasik.v1.AiInferenceService", healthpb.HealthCheckResponse_SERVING)

	// Reflection (개발용)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}

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

	go func() {
		logger.Info("ai-inference-service started", zap.String("port", port))
		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatal("failed to serve", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	fmt.Printf("\nReceived signal: %v. Shutting down...\n", sig)
	grpcServer.GracefulStop()
	logger.Info("ai-inference-service stopped")
}
