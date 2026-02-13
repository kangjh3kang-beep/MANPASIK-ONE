// user-service: 사용자 프로필·구독·가족 그룹 관리 마이크로서비스
//
// 포트: gRPC :50052
// 의존: PostgreSQL(선택) — 미설정 시 인메모리 저장소 사용
//
// 기능:
// - 프로필 CRUD (GetProfile, UpdateProfile)
// - 구독 관리 (GetSubscription, 4티어: Free/Basic/Pro/Clinical)
// - 가족 그룹 관리
// - gRPC UserService
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
	"github.com/manpasik/backend/services/user-service/internal/handler"
	"github.com/manpasik/backend/services/user-service/internal/repository/memory"
	"github.com/manpasik/backend/services/user-service/internal/repository/postgres"
	"github.com/manpasik/backend/services/user-service/internal/service"
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

const serviceName = "user-service"

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

	// ProfileRepository / SubscriptionRepository / FamilyRepository: PostgreSQL 또는 인메모리
	var profileRepo service.ProfileRepository
	var subRepo service.SubscriptionRepository
	var familyRepo service.FamilyRepository
	if _, dbHostSet := os.LookupEnv("DB_HOST"); dbHostSet && cfg.DB.Host != "" && cfg.DB.DBName != "" {
		connCtx, connCancel := context.WithTimeout(context.Background(), 5*time.Second)
		pool, err := pgxpool.New(connCtx, cfg.DB.DSN())
		connCancel()
		if err != nil {
			log.Printf("[%s] DB 풀 생성 실패, 인메모리 사용: %v", serviceName, err)
			profileRepo = memory.NewProfileRepository()
			subRepo = memory.NewSubscriptionRepository()
			familyRepo = memory.NewFamilyRepository()
		} else {
			pingCtx, pingCancel := context.WithTimeout(context.Background(), 3*time.Second)
			if pingErr := pool.Ping(pingCtx); pingErr != nil {
				pingCancel()
				pool.Close()
				log.Printf("[%s] DB Ping 실패, 인메모리 사용: %v", serviceName, pingErr)
				profileRepo = memory.NewProfileRepository()
				subRepo = memory.NewSubscriptionRepository()
				familyRepo = memory.NewFamilyRepository()
			} else {
				pingCancel()
				defer pool.Close()
				profileRepo = postgres.NewProfileRepository(pool)
				subRepo = postgres.NewSubscriptionRepository(pool)
				familyRepo = postgres.NewFamilyRepository(pool)
				log.Printf("[%s] DB 연결됨: %s", serviceName, cfg.DB.DBName)
			}
		}
	} else {
		profileRepo = memory.NewProfileRepository()
		subRepo = memory.NewSubscriptionRepository()
		familyRepo = memory.NewFamilyRepository()
		log.Printf("[%s] 인메모리 저장소 사용", serviceName)
	}

	userSvc := service.NewUserService(logger, profileRepo, subRepo, familyRepo)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.RequestIDInterceptor(),
			observability.UnaryServerInterceptor(metrics),
		),
	)

	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus(serviceName, healthpb.HealthCheckResponse_SERVING)

	userHandler := handler.NewUserHandler(userSvc, logger)
	v1.RegisterUserServiceServer(grpcServer, userHandler)

	reflection.Register(grpcServer)

	grpcPort := cfg.GRPCPort
	if grpcPort == "" {
		grpcPort = ":50052"
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
