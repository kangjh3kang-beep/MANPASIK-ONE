// assistant-service: AI 비서 마이크로서비스
// 포트: gRPC :50077
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
	aimod "github.com/manpasik/backend/services/assistant-service/internal/ai"
	"github.com/manpasik/backend/services/assistant-service/internal/handler"
	"github.com/manpasik/backend/services/assistant-service/internal/repository/memory"
	"github.com/manpasik/backend/services/assistant-service/internal/repository/postgres"
	"github.com/manpasik/backend/services/assistant-service/internal/service"
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

const serviceName = "assistant-service"

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

	var repo service.AssistantRepository

	if _, dbHostSet := os.LookupEnv("DB_HOST"); dbHostSet && cfg.DB.Host != "" && cfg.DB.DBName != "" {
		connCtx, connCancel := context.WithTimeout(context.Background(), 5*time.Second)
		pool, poolErr := pgxpool.New(connCtx, cfg.DB.DSN())
		connCancel()
		if poolErr != nil {
			log.Printf("[%s] DB connection failed, using memory: %v", serviceName, poolErr)
			repo = memory.NewAssistantRepository()
		} else {
			pingCtx, pingCancel := context.WithTimeout(context.Background(), 3*time.Second)
			if pingErr := pool.Ping(pingCtx); pingErr != nil {
				pingCancel()
				pool.Close()
				log.Printf("[%s] DB ping failed, using memory: %v", serviceName, pingErr)
				repo = memory.NewAssistantRepository()
			} else {
				pingCancel()
				defer pool.Close()
				log.Printf("[%s] Connected to PostgreSQL", serviceName)
				repo = postgres.NewAssistantRepository(pool)
			}
		}
	} else {
		repo = memory.NewAssistantRepository()
	}
	svc := service.NewAssistantService(logger, repo)

	// OpenAI AI 응답기: OPENAI_API_KEY 환경변수 기반
	if apiKey := os.Getenv("OPENAI_API_KEY"); apiKey != "" {
		model := os.Getenv("OPENAI_MODEL")
		if model == "" {
			model = "gpt-4o"
		}
		baseURL := os.Getenv("OPENAI_BASE_URL")
		responder := aimod.NewOpenAIResponder(aimod.OpenAIConfig{
			APIKey:  apiKey,
			Model:   model,
			BaseURL: baseURL,
		})
		svc.SetAIResponder(responder)
		log.Printf("[%s] OpenAI AI 응답기 활성화 (model=%s)", serviceName, model)
	} else {
		log.Printf("[%s] OPENAI_API_KEY 미설정, 키워드 기반 응답 사용", serviceName)
	}

	grpcHandler := handler.NewAssistantHandler(svc, logger)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.RequestIDInterceptor(),
			observability.UnaryServerInterceptor(metrics),
		),
	)

	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus(serviceName, healthpb.HealthCheckResponse_SERVING)

	v1.RegisterAssistantServiceServer(grpcServer, grpcHandler)
	reflection.Register(grpcServer)

	grpcPort := cfg.GRPCPort
	if grpcPort == "" {
		grpcPort = ":50077"
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
		<-sigCh
		healthServer.SetServingStatus(serviceName, healthpb.HealthCheckResponse_NOT_SERVING)
		go func() { time.Sleep(cfg.ShutdownTimeout); os.Exit(1) }()
		grpcServer.GracefulStop()
		cancel()
	}()

	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/metrics", metrics.PrometheusHandler())
		mux.HandleFunc("/health", healthCheck.Handler())
		http.ListenAndServe(":9100", mux)
	}()

	log.Printf("[%s] gRPC server listening on %s", serviceName, grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("[%s] Failed to serve: %v", serviceName, err)
	}
	<-ctx.Done()
}
