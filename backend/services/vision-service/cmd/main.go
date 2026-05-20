// vision-service: 음식 인식 및 칼로리 분석 마이크로서비스
//
// 포트: gRPC :50071
// 의존: 인메모리 저장소 (향후 PostgreSQL 추가)
//
// 기능:
// - 음식 이미지 AI 분석
// - 칼로리 및 영양소 추정
// - 분석 이력 조회
// - 일일 영양 섭취 요약
//
// NOTE: Proto 정의가 추가되면 gRPC VisionService 등록 활성화
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

	"github.com/manpasik/backend/services/vision-service/internal/handler"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manpasik/backend/services/vision-service/internal/repository/memory"
	"github.com/manpasik/backend/services/vision-service/internal/repository/postgres"
	"github.com/manpasik/backend/services/vision-service/internal/service"
	"github.com/manpasik/backend/services/vision-service/internal/vision"
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

const serviceName = "vision-service"

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

	// FoodAnalysisRepository: DB 스위치 (PostgreSQL / 인메모리)
	var analysisRepo service.FoodAnalysisRepository

	if _, dbHostSet := os.LookupEnv("DB_HOST"); dbHostSet && cfg.DB.Host != "" && cfg.DB.DBName != "" {
		connCtx, connCancel := context.WithTimeout(context.Background(), 5*time.Second)
		pool, poolErr := pgxpool.New(connCtx, cfg.DB.DSN())
		connCancel()
		if poolErr != nil {
			log.Printf("[%s] DB connection failed, using memory: %v", serviceName, poolErr)
			analysisRepo = memory.NewFoodAnalysisRepository()
		} else {
			pingCtx, pingCancel := context.WithTimeout(context.Background(), 3*time.Second)
			if pingErr := pool.Ping(pingCtx); pingErr != nil {
				pingCancel()
				pool.Close()
				log.Printf("[%s] DB ping failed, using memory: %v", serviceName, pingErr)
				analysisRepo = memory.NewFoodAnalysisRepository()
			} else {
				pingCancel()
				defer pool.Close()
				log.Printf("[%s] Connected to PostgreSQL", serviceName)
				analysisRepo = postgres.NewFoodAnalysisRepository(pool)
			}
		}
	} else {
		analysisRepo = memory.NewFoodAnalysisRepository()
	}

	visionSvc := service.NewVisionService(logger, analysisRepo)

	// OpenAI Vision Analyzer 설정 (환경변수 기반)
	if apiKey := os.Getenv("VISION_API_KEY"); apiKey != "" {
		model := os.Getenv("VISION_MODEL")
		baseURL := os.Getenv("VISION_BASE_URL")
		analyzer := vision.NewOpenAIVisionAnalyzer(apiKey, model, baseURL)
		visionSvc.SetAnalyzer(&visionAnalyzerAdapter{inner: analyzer})
		log.Printf("[%s] OpenAI Vision Analyzer 활성화 (model=%s)", serviceName, model)
	} else {
		log.Printf("[%s] Vision Analyzer 미설정, fallback 모드", serviceName)
	}

	visionHandler := handler.NewVisionHandler(visionSvc, logger)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.RequestIDInterceptor(),
			observability.UnaryServerInterceptor(metrics),
		),
	)

	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus(serviceName, healthpb.HealthCheckResponse_SERVING)

	v1.RegisterVisionServiceServer(grpcServer, visionHandler)

	reflection.Register(grpcServer)

	grpcPort := cfg.GRPCPort
	if grpcPort == "" {
		grpcPort = ":50071"
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

	log.Printf("[%s] gRPC server listening on %s", serviceName, grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("[%s] Failed to serve: %v", serviceName, err)
	}
	<-ctx.Done()
	log.Printf("[%s] Shutdown complete", serviceName)
}

// visionAnalyzerAdapter는 vision.OpenAIVisionAnalyzer를 service.VisionAnalyzer로 래핑합니다.
type visionAnalyzerAdapter struct {
	inner *vision.OpenAIVisionAnalyzer
}

func (a *visionAnalyzerAdapter) AnalyzeFood(ctx context.Context, imageURL string) ([]service.FoodItem, error) {
	results, err := a.inner.AnalyzeFood(ctx, imageURL)
	if err != nil {
		return nil, err
	}
	items := make([]service.FoodItem, len(results))
	for i, r := range results {
		nutrients := make([]service.NutrientInfo, len(r.Nutrients))
		for j, n := range r.Nutrients {
			nutrients[j] = service.NutrientInfo{Name: n.Name, Amount: n.Amount, Unit: n.Unit, DV: n.DV}
		}
		items[i] = service.FoodItem{
			Name:        r.Name,
			Confidence:  r.Confidence,
			CalorieKcal: r.CalorieKcal,
			PortionG:    r.PortionG,
			Nutrients:   nutrients,
		}
	}
	return items, nil
}
