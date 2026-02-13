// payment-service: 결제 관리 마이크로서비스
//
// 포트: gRPC :50057
// 의존: PostgreSQL(선택) — 미설정 시 인메모리 저장소 사용
//
// 기능:
// - 결제 요청 생성 (일회성·구독)
// - PG 콜백 결제 확인
// - 결제 조회/이력
// - 환불 처리 (전액·부분)
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
	"github.com/manpasik/backend/services/payment-service/internal/handler"
	"github.com/manpasik/backend/services/payment-service/internal/pg"
	kafkaPublisher "github.com/manpasik/backend/services/payment-service/internal/repository/kafka"
	"github.com/manpasik/backend/services/payment-service/internal/repository/memory"
	"github.com/manpasik/backend/services/payment-service/internal/repository/postgres"
	"github.com/manpasik/backend/services/payment-service/internal/service"
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

const serviceName = "payment-service"

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
	var payRepo service.PaymentRepository
	var refundRepo service.RefundRepository

	if _, dbHostSet := os.LookupEnv("DB_HOST"); dbHostSet && cfg.DB.Host != "" && cfg.DB.DBName != "" {
		connCtx, connCancel := context.WithTimeout(context.Background(), 5*time.Second)
		pool, poolErr := pgxpool.New(connCtx, cfg.DB.DSN())
		connCancel()
		if poolErr != nil {
			log.Printf("[%s] DB 풀 생성 실패, 인메모리 사용: %v", serviceName, poolErr)
			payRepo = memory.NewPaymentRepository()
			refundRepo = memory.NewRefundRepository()
		} else {
			pingCtx, pingCancel := context.WithTimeout(context.Background(), 3*time.Second)
			if pingErr := pool.Ping(pingCtx); pingErr != nil {
				pingCancel()
				pool.Close()
				log.Printf("[%s] DB Ping 실패, 인메모리 사용: %v", serviceName, pingErr)
				payRepo = memory.NewPaymentRepository()
				refundRepo = memory.NewRefundRepository()
			} else {
				pingCancel()
				defer pool.Close()
				payRepo = postgres.NewPaymentRepository(pool)
				refundRepo = postgres.NewRefundRepository(pool)
				log.Printf("[%s] DB 연결됨: %s", serviceName, cfg.DB.DBName)
			}
		}
	} else {
		payRepo = memory.NewPaymentRepository()
		refundRepo = memory.NewRefundRepository()
		log.Printf("[%s] 인메모리 저장소 사용", serviceName)
	}

	paySvc := service.NewPaymentService(logger, payRepo, refundRepo)

	// EventPublisher: Kafka(Redpanda) 또는 인메모리
	if _, kafkaBrokersSet := os.LookupEnv("KAFKA_BROKERS"); kafkaBrokersSet && len(cfg.Kafka.Brokers) > 0 {
		eventBus, kafkaErr := events.NewKafkaEventBus(events.KafkaAdapterConfig{
			Brokers:     cfg.Kafka.Brokers,
			GroupID:     serviceName,
			TopicPrefix: "manpasik.",
		})
		if kafkaErr != nil {
			log.Printf("[%s] Kafka 연결 실패, 인메모리 EventPublisher 사용: %v", serviceName, kafkaErr)
			paySvc.SetEventPublisher(memory.NewEventPublisher())
		} else {
			defer eventBus.Close()
			paySvc.SetEventPublisher(kafkaPublisher.NewEventPublisher(eventBus))
			log.Printf("[%s] Kafka 연결됨: %v", serviceName, cfg.Kafka.Brokers)
		}
	} else {
		paySvc.SetEventPublisher(memory.NewEventPublisher())
		log.Printf("[%s] 인메모리 EventPublisher 사용", serviceName)
	}

	// PG 연동: DB config 우선 → 환경변수 fallback
	tossSecret := config.LoadConfigWithFallback(nil, "toss.secret_key", "TOSS_SECRET_KEY")
	tossAPIURL := config.LoadConfigWithFallback(nil, "toss.api_url", "TOSS_API_URL")
	if tossAPIURL == "" {
		tossAPIURL = cfg.Toss.APIURL
	}
	if tossSecret == "" {
		tossSecret = cfg.Toss.SecretKey
	}

	if tossSecret != "" {
		paySvc.SetPaymentGateway(pg.NewTossClient(tossSecret, tossAPIURL))
		log.Printf("[%s] Toss PG 연동 사용", serviceName)
	} else {
		paySvc.SetPaymentGateway(pg.NewNoopGateway())
		log.Printf("[%s] PG Noop 사용 (TOSS_SECRET_KEY 미설정)", serviceName)
	}

	// ConfigWatcher: 설정 변경 시 PG 게이트웨이 핫리로드
	configWatcher := events.NewEventBusConfigWatcher(func() events.EventPublisher {
		// EventPublisher가 아직 없으면 인메모리 사용
		return events.NewEventBus()
	}())
	_ = configWatcher.Watch(context.Background(), serviceName, func(key, newValue string) error {
		switch key {
		case "toss.secret_key":
			newTossURL := config.LoadConfigWithFallback(nil, "toss.api_url", "TOSS_API_URL")
			if newTossURL == "" {
				newTossURL = cfg.Toss.APIURL
			}
			paySvc.SetPaymentGateway(pg.NewTossClient(newValue, newTossURL))
			log.Printf("[%s] Toss PG 핫리로드 완료 (config.changed)", serviceName)
		case "toss.api_url":
			currentSecret := config.LoadConfigWithFallback(nil, "toss.secret_key", "TOSS_SECRET_KEY")
			if currentSecret == "" {
				currentSecret = cfg.Toss.SecretKey
			}
			if currentSecret != "" {
				paySvc.SetPaymentGateway(pg.NewTossClient(currentSecret, newValue))
				log.Printf("[%s] Toss API URL 핫리로드 완료", serviceName)
			}
		}
		return nil
	})
	defer configWatcher.Close()

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.RequestIDInterceptor(),
			observability.UnaryServerInterceptor(metrics),
		),
	)

	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus(serviceName, healthpb.HealthCheckResponse_SERVING)

	payHandler := handler.NewPaymentHandler(paySvc, logger)
	v1.RegisterPaymentServiceServer(grpcServer, payHandler)

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
