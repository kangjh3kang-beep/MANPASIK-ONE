// admin-service: 관리자 및 시스템 관리 마이크로서비스
//
// 포트: gRPC :50068
// 의존: PostgreSQL(선택) — 미설정 시 인메모리 저장소 사용
//
// 기능:
// - 관리자 생성 / 조회 / 목록
// - 관리자 역할 변경 / 비활성화
// - 사용자 목록 조회 (관리자용)
// - 시스템 통계 조회
// - 감사 로그 조회
// - 시스템 설정 관리
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
	"github.com/manpasik/backend/services/admin-service/internal/crypto"
	"github.com/manpasik/backend/services/admin-service/internal/handler"
	"github.com/manpasik/backend/services/admin-service/internal/repository/memory"
	"github.com/manpasik/backend/services/admin-service/internal/repository/postgres"
	"github.com/manpasik/backend/services/admin-service/internal/service"
	"github.com/manpasik/backend/shared/config"
	"github.com/manpasik/backend/shared/events"
	v1 "github.com/manpasik/backend/shared/gen/go/v1"
	"github.com/manpasik/backend/shared/middleware"
	"github.com/manpasik/backend/shared/observability"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

const serviceName = "admin-service"

func main() {
	cfg := config.LoadFromEnv(serviceName)
	if cfg.GRPCPort == ":50051" {
		cfg.GRPCPort = ":50068"
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

	var adminRepo service.AdminRepository
	var auditRepo service.AuditLogRepository
	var configRepo service.SystemConfigRepository
	var userRepo service.UserSummaryRepository
	var metaRepo service.ConfigMetadataRepository
	var transRepo service.ConfigTranslationRepository

	if _, dbHostSet := os.LookupEnv("DB_HOST"); dbHostSet && cfg.DB.Host != "" && cfg.DB.DBName != "" {
		connCtx, connCancel := context.WithTimeout(context.Background(), 5*time.Second)
		pool, err := pgxpool.New(connCtx, cfg.DB.DSN())
		connCancel()
		if err != nil {
			log.Printf("[%s] DB 풀 생성 실패, 인메모리 사용: %v", serviceName, err)
			adminRepo = memory.NewAdminRepository()
			auditRepo = memory.NewAuditLogRepository()
			configRepo = memory.NewSystemConfigRepository()
			userRepo = memory.NewUserSummaryRepository()
			metaRepo = memory.NewConfigMetadataRepository()
			transRepo = memory.NewConfigTranslationRepository()
		} else {
			pingCtx, pingCancel := context.WithTimeout(context.Background(), 3*time.Second)
			if pingErr := pool.Ping(pingCtx); pingErr != nil {
				pingCancel()
				pool.Close()
				log.Printf("[%s] DB Ping 실패, 인메모리 사용: %v", serviceName, pingErr)
				adminRepo = memory.NewAdminRepository()
				auditRepo = memory.NewAuditLogRepository()
				configRepo = memory.NewSystemConfigRepository()
				userRepo = memory.NewUserSummaryRepository()
				metaRepo = memory.NewConfigMetadataRepository()
				transRepo = memory.NewConfigTranslationRepository()
			} else {
				pingCancel()
				defer pool.Close()
				adminRepo = postgres.NewAdminRepository(pool)
				auditRepo = postgres.NewAuditLogRepository(pool)
				configRepo = postgres.NewSystemConfigRepository(pool)
				userRepo = postgres.NewUserSummaryRepository(pool)
				metaRepo = postgres.NewConfigMetadataRepository(pool)
				transRepo = postgres.NewConfigTranslationRepository(pool)
				log.Printf("[%s] DB 연결됨: %s", serviceName, cfg.DB.DBName)
			}
		}
	} else {
		adminRepo = memory.NewAdminRepository()
		auditRepo = memory.NewAuditLogRepository()
		configRepo = memory.NewSystemConfigRepository()
		userRepo = memory.NewUserSummaryRepository()
		metaRepo = memory.NewConfigMetadataRepository()
		transRepo = memory.NewConfigTranslationRepository()
		log.Printf("[%s] 인메모리 저장소 사용", serviceName)
	}

	adminSvc := service.NewAdminService(logger, adminRepo, auditRepo, configRepo, userRepo)

	// AES-256-GCM 암호화기 (CONFIG_ENCRYPTION_KEY 환경변수)
	encryptor, encErr := crypto.NewAESEncryptor(os.Getenv("CONFIG_ENCRYPTION_KEY"))
	if encErr != nil {
		log.Printf("[%s] 암호화 키 로드 실패 (암호화 비활성): %v", serviceName, encErr)
	}

	// 이벤트 버스 (Kafka 또는 인메모리 fallback)
	var eventPublisher events.EventPublisher
	kafkaCfg := events.KafkaAdapterConfig{
		Brokers:     cfg.Kafka.Brokers,
		GroupID:     cfg.Kafka.GroupID,
		TopicPrefix: "manpasik.",
	}
	kafkaBus, kafkaErr := events.NewKafkaEventBus(kafkaCfg)
	if kafkaErr != nil {
		log.Printf("[%s] Kafka 연결 실패, 인메모리 이벤트 버스 사용: %v", serviceName, kafkaErr)
		eventPublisher = events.NewEventBus()
	} else {
		eventPublisher = kafkaBus
		defer kafkaBus.Close()
		log.Printf("[%s] Kafka 이벤트 버스 연결됨", serviceName)
	}

	// ConfigManager 생성
	cfgMgr := service.NewConfigManager(logger, configRepo, metaRepo, transRepo, auditRepo, encryptor, eventPublisher)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.RequestIDInterceptor(),
			observability.UnaryServerInterceptor(metrics),
		),
	)

	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus(serviceName, healthpb.HealthCheckResponse_SERVING)

	adminHandler := handler.NewAdminHandler(adminSvc, logger)
	adminHandler.SetConfigManager(cfgMgr)
	v1.RegisterAdminServiceServer(grpcServer, adminHandler)

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

	log.Printf("[%s] gRPC server listening on %s", serviceName, cfg.GRPCPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("[%s] Failed to serve: %v", serviceName, err)
	}
	<-ctx.Done()
	log.Printf("[%s] Shutdown complete", serviceName)
}
