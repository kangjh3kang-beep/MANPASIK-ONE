// measurement-service: 측정 세션·시계열 데이터·벡터 검색 마이크로서비스
//
// 포트: gRPC :50054
// 의존: PostgreSQL/TimescaleDB(선택), Milvus(선택), Kafka(선택) — 미설정 시 인메모리
//
// 기능:
// - 측정 세션 관리 (StartSession, EndSession)
// - 시계열 데이터 저장 (TimescaleDB 또는 인메모리)
// - 핑거프린트 벡터 저장 (Milvus 또는 인메모리)
// - 측정 완료 이벤트 발행 (Kafka 또는 인메모리)
// - 측정 기록 조회 (GetMeasurementHistory)
// - gRPC MeasurementService
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
	"github.com/manpasik/backend/services/measurement-service/internal/handler"
	esRepo "github.com/manpasik/backend/services/measurement-service/internal/repository/elasticsearch"
	kafkaPublisher "github.com/manpasik/backend/services/measurement-service/internal/repository/kafka"
	"github.com/manpasik/backend/services/measurement-service/internal/repository/memory"
	milvusRepo "github.com/manpasik/backend/services/measurement-service/internal/repository/milvus"
	"github.com/manpasik/backend/services/measurement-service/internal/repository/postgres"
	"github.com/manpasik/backend/services/measurement-service/internal/service"
	"github.com/manpasik/backend/shared/config"
	"github.com/manpasik/backend/shared/events"
	"github.com/manpasik/backend/shared/search"
	v1 "github.com/manpasik/backend/shared/gen/go/v1"
	"github.com/manpasik/backend/shared/middleware"
	"github.com/manpasik/backend/shared/observability"
	"github.com/manpasik/backend/shared/vectordb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

const serviceName = "measurement-service"

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

	// SessionRepository / MeasurementRepository: PostgreSQL/TimescaleDB 또는 인메모리
	var sessionRepo service.SessionRepository
	var measureRepo service.MeasurementRepository
	if _, dbHostSet := os.LookupEnv("DB_HOST"); dbHostSet && cfg.DB.Host != "" && cfg.DB.DBName != "" {
		connCtx, connCancel := context.WithTimeout(context.Background(), 5*time.Second)
		pool, err := pgxpool.New(connCtx, cfg.DB.DSN())
		connCancel()
		if err != nil {
			log.Printf("[%s] DB 풀 생성 실패, 인메모리 사용: %v", serviceName, err)
			sessionRepo = memory.NewSessionRepository()
			measureRepo = memory.NewMeasurementRepository()
		} else {
			pingCtx, pingCancel := context.WithTimeout(context.Background(), 3*time.Second)
			if pingErr := pool.Ping(pingCtx); pingErr != nil {
				pingCancel()
				pool.Close()
				log.Printf("[%s] DB Ping 실패, 인메모리 사용: %v", serviceName, pingErr)
				sessionRepo = memory.NewSessionRepository()
				measureRepo = memory.NewMeasurementRepository()
			} else {
				pingCancel()
				defer pool.Close()
				sessionRepo = postgres.NewSessionRepository(pool)
				measureRepo = postgres.NewMeasurementRepository(pool)
				log.Printf("[%s] DB 연결됨: %s", serviceName, cfg.DB.DBName)
			}
		}
	} else {
		sessionRepo = memory.NewSessionRepository()
		measureRepo = memory.NewMeasurementRepository()
		log.Printf("[%s] 인메모리 저장소 사용", serviceName)
	}

	// VectorRepository: Milvus 또는 인메모리
	var vectorRepo service.VectorRepository
	if _, milvusHostSet := os.LookupEnv("MILVUS_HOST"); milvusHostSet && cfg.Milvus.Host != "" {
		milvusClient, milvusErr := vectordb.NewMilvusClient(
			cfg.Milvus.Addr(),
			cfg.Milvus.CollectionName,
			896, // 핑거프린트 최대 차원 (MAX_CHANNELS)
		)
		if milvusErr != nil {
			log.Printf("[%s] Milvus 연결 실패, 인메모리 VectorRepo 사용: %v", serviceName, milvusErr)
			vectorRepo = memory.NewVectorRepository()
		} else {
			defer milvusClient.Close()
			vectorRepo = milvusRepo.NewVectorRepository(milvusClient)
			log.Printf("[%s] Milvus 연결됨: %s (컬렉션: %s)", serviceName, cfg.Milvus.Addr(), cfg.Milvus.CollectionName)
		}
	} else {
		vectorRepo = memory.NewVectorRepository()
		log.Printf("[%s] 인메모리 VectorRepo 사용", serviceName)
	}

	// EventPublisher: Kafka(Redpanda) 또는 인메모리
	var eventPublisher service.EventPublisher
	if _, kafkaBrokersSet := os.LookupEnv("KAFKA_BROKERS"); kafkaBrokersSet && len(cfg.Kafka.Brokers) > 0 {
		eventBus, kafkaErr := events.NewKafkaEventBus(events.KafkaAdapterConfig{
			Brokers:     cfg.Kafka.Brokers,
			GroupID:     serviceName,
			TopicPrefix: "manpasik.",
		})
		if kafkaErr != nil {
			log.Printf("[%s] Kafka 연결 실패, 인메모리 EventPublisher 사용: %v", serviceName, kafkaErr)
			eventPublisher = memory.NewEventPublisher()
		} else {
			defer eventBus.Close()
			eventPublisher = kafkaPublisher.NewEventPublisher(eventBus)
			log.Printf("[%s] Kafka 연결됨: %v", serviceName, cfg.Kafka.Brokers)
		}
	} else {
		eventPublisher = memory.NewEventPublisher()
		log.Printf("[%s] 인메모리 EventPublisher 사용", serviceName)
	}

	measureSvc := service.NewMeasurementService(logger, sessionRepo, measureRepo, vectorRepo, eventPublisher)

	// SearchIndexer: Elasticsearch 또는 인메모리 (no-op)
	if _, esURLSet := os.LookupEnv("ELASTICSEARCH_URL"); esURLSet && cfg.Elasticsearch.URL != "" {
		esClient, esErr := search.NewESClient(
			cfg.Elasticsearch.URL,
			cfg.Elasticsearch.Username,
			cfg.Elasticsearch.Password,
		)
		if esErr != nil {
			log.Printf("[%s] Elasticsearch 연결 실패, 검색 인덱싱 비활성화: %v", serviceName, esErr)
		} else {
			defer esClient.Close()
			indexer := esRepo.NewSearchIndexer(esClient)
			_ = indexer.EnsureIndex(context.Background())
			measureSvc.SetSearchIndexer(indexer)
			log.Printf("[%s] Elasticsearch 연결됨: %s", serviceName, cfg.Elasticsearch.URL)
		}
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

	measureHandler := handler.NewMeasurementHandler(measureSvc, logger)
	v1.RegisterMeasurementServiceServer(grpcServer, measureHandler)

	reflection.Register(grpcServer)

	grpcPort := cfg.GRPCPort
	if grpcPort == "" {
		grpcPort = ":50054"
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
