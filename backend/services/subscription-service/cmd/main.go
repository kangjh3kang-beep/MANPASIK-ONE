// subscription-service: SaaS 구독 관리 마이크로서비스
//
// 포트: gRPC :50055
// 의존: PostgreSQL(선택) — 미설정 시 인메모리 저장소 사용
//
// 기능:
// - 구독 생성 (회원가입 시 Free 자동)
// - 구독 조회 / 변경 (업·다운그레이드)
// - 구독 해지
// - 기능 접근 권한 확인 (티어별)
// - 구독 플랜 목록 조회
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
	"github.com/manpasik/backend/services/subscription-service/internal/handler"
	cacheRepo "github.com/manpasik/backend/services/subscription-service/internal/repository/cache"
	kafkaPublisher "github.com/manpasik/backend/services/subscription-service/internal/repository/kafka"
	"github.com/manpasik/backend/services/subscription-service/internal/repository/memory"
	"github.com/manpasik/backend/services/subscription-service/internal/repository/postgres"
	"github.com/manpasik/backend/services/subscription-service/internal/service"
	"github.com/manpasik/backend/shared/config"
	"github.com/manpasik/backend/shared/events"
	v1 "github.com/manpasik/backend/shared/gen/go/v1"
	"github.com/manpasik/backend/shared/middleware"
	"github.com/manpasik/backend/shared/observability"
	redisclient "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

const serviceName = "subscription-service"

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

	// SubscriptionRepository: PostgreSQL 또는 인메모리
	var subRepo service.SubscriptionRepository
	var historyRepo service.SubscriptionHistoryRepository

	if _, dbHostSet := os.LookupEnv("DB_HOST"); dbHostSet && cfg.DB.Host != "" && cfg.DB.DBName != "" {
		connCtx, connCancel := context.WithTimeout(context.Background(), 5*time.Second)
		pool, poolErr := pgxpool.New(connCtx, cfg.DB.DSN())
		connCancel()
		if poolErr != nil {
			log.Printf("[%s] DB 풀 생성 실패, 인메모리 사용: %v", serviceName, poolErr)
			subRepo = memory.NewSubscriptionRepository()
			historyRepo = memory.NewSubscriptionHistoryRepository()
		} else {
			pingCtx, pingCancel := context.WithTimeout(context.Background(), 3*time.Second)
			if pingErr := pool.Ping(pingCtx); pingErr != nil {
				pingCancel()
				pool.Close()
				log.Printf("[%s] DB Ping 실패, 인메모리 사용: %v", serviceName, pingErr)
				subRepo = memory.NewSubscriptionRepository()
				historyRepo = memory.NewSubscriptionHistoryRepository()
			} else {
				pingCancel()
				defer pool.Close()
				subRepo = postgres.NewSubscriptionRepository(pool)
				historyRepo = postgres.NewSubscriptionHistoryRepository(pool)
				log.Printf("[%s] DB 연결됨: %s", serviceName, cfg.DB.DBName)
			}
		}
	} else {
		subRepo = memory.NewSubscriptionRepository()
		historyRepo = memory.NewSubscriptionHistoryRepository()
		log.Printf("[%s] 인메모리 저장소 사용", serviceName)
	}

	// Redis 캐시 레이어: REDIS_HOST 설정 시 SubscriptionRepository를 캐시 데코레이터로 래핑
	if _, redisHostSet := os.LookupEnv("REDIS_HOST"); redisHostSet && cfg.Redis.Host != "" {
		rdb := redisclient.NewClient(&redisclient.Options{
			Addr:     cfg.Redis.Addr(),
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		})
		pingCtx, pingCancel := context.WithTimeout(context.Background(), 3*time.Second)
		if err := rdb.Ping(pingCtx).Err(); err != nil {
			pingCancel()
			log.Printf("[%s] Redis 연결 실패, 캐시 미사용: %v", serviceName, err)
		} else {
			pingCancel()
			defer rdb.Close()
			subRepo = cacheRepo.NewSubscriptionRepository(subRepo, rdb)
			log.Printf("[%s] Redis 캐시 SubscriptionRepository 연결됨: %s", serviceName, cfg.Redis.Addr())
		}
	}

	subSvc := service.NewSubscriptionService(logger, subRepo, historyRepo)

	// EventPublisher: Kafka(Redpanda) 또는 인메모리
	if _, kafkaBrokersSet := os.LookupEnv("KAFKA_BROKERS"); kafkaBrokersSet && len(cfg.Kafka.Brokers) > 0 {
		eventBus, kafkaErr := events.NewKafkaEventBus(events.KafkaAdapterConfig{
			Brokers:     cfg.Kafka.Brokers,
			GroupID:     serviceName,
			TopicPrefix: "manpasik.",
		})
		if kafkaErr != nil {
			log.Printf("[%s] Kafka 연결 실패, 인메모리 EventPublisher 사용: %v", serviceName, kafkaErr)
			subSvc.SetEventPublisher(memory.NewEventPublisher())
		} else {
			defer eventBus.Close()
			subSvc.SetEventPublisher(kafkaPublisher.NewEventPublisher(eventBus))
			log.Printf("[%s] Kafka 연결됨: %v", serviceName, cfg.Kafka.Brokers)
		}
	} else {
		subSvc.SetEventPublisher(memory.NewEventPublisher())
		log.Printf("[%s] 인메모리 EventPublisher 사용", serviceName)
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.RequestIDInterceptor(),
			observability.UnaryServerInterceptor(metrics),
		),
	)

	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus(serviceName, healthpb.HealthCheckResponse_SERVING)

	subHandler := handler.NewSubscriptionHandler(subSvc, logger)
	v1.RegisterSubscriptionServiceServer(grpcServer, subHandler)

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
