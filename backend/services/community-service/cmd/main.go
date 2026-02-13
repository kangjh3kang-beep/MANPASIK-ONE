// community-service: 건강 커뮤니티 마이크로서비스
//
// 포트: gRPC :50067
// 의존: PostgreSQL(선택) — 미설정 시 인메모리 저장소 사용
//
// 기능:
// - 게시글 작성 / 조회 / 목록 / 좋아요
// - 댓글 작성 / 목록
// - 건강 챌린지 생성 / 조회 / 참가 / 목록
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
	"github.com/manpasik/backend/services/community-service/internal/handler"
	esRepo "github.com/manpasik/backend/services/community-service/internal/repository/elasticsearch"
	"github.com/manpasik/backend/services/community-service/internal/repository/memory"
	"github.com/manpasik/backend/services/community-service/internal/repository/postgres"
	"github.com/manpasik/backend/services/community-service/internal/service"
	"github.com/manpasik/backend/shared/config"
	"github.com/manpasik/backend/shared/search"
	v1 "github.com/manpasik/backend/shared/gen/go/v1"
	"github.com/manpasik/backend/shared/middleware"
	"github.com/manpasik/backend/shared/observability"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

const serviceName = "community-service"

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

	// PostRepository, CommentRepository, ChallengeRepository: PostgreSQL 또는 인메모리
	var postRepo service.PostRepository
	var commentRepo service.CommentRepository
	var challengeRepo service.ChallengeRepository

	if _, dbHostSet := os.LookupEnv("DB_HOST"); dbHostSet && cfg.DB.Host != "" && cfg.DB.DBName != "" {
		connCtx, connCancel := context.WithTimeout(context.Background(), 5*time.Second)
		pool, poolErr := pgxpool.New(connCtx, cfg.DB.DSN())
		connCancel()
		if poolErr != nil {
			log.Printf("[%s] DB 풀 생성 실패, 인메모리 사용: %v", serviceName, poolErr)
			postRepo = memory.NewPostRepository()
			commentRepo = memory.NewCommentRepository()
			challengeRepo = memory.NewChallengeRepository()
		} else {
			pingCtx, pingCancel := context.WithTimeout(context.Background(), 3*time.Second)
			if pingErr := pool.Ping(pingCtx); pingErr != nil {
				pingCancel()
				pool.Close()
				log.Printf("[%s] DB Ping 실패, 인메모리 사용: %v", serviceName, pingErr)
				postRepo = memory.NewPostRepository()
				commentRepo = memory.NewCommentRepository()
				challengeRepo = memory.NewChallengeRepository()
			} else {
				pingCancel()
				defer pool.Close()
				postRepo = postgres.NewPostRepository(pool)
				commentRepo = postgres.NewCommentRepository(pool)
				challengeRepo = postgres.NewChallengeRepository(pool)
				log.Printf("[%s] DB 연결됨: %s", serviceName, cfg.DB.DBName)
			}
		}
	} else {
		postRepo = memory.NewPostRepository()
		commentRepo = memory.NewCommentRepository()
		challengeRepo = memory.NewChallengeRepository()
		log.Printf("[%s] 인메모리 저장소 사용", serviceName)
	}

	communitySvc := service.NewCommunityService(logger, postRepo, commentRepo, challengeRepo)

	// SearchIndexer: Elasticsearch 또는 비활성화
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
			indexer := esRepo.NewPostSearchIndexer(esClient)
			_ = indexer.EnsureIndex(context.Background())
			communitySvc.SetSearchIndexer(indexer)
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

	communityHandler := handler.NewCommunityHandler(communitySvc, logger)
	v1.RegisterCommunityServiceServer(grpcServer, communityHandler)

	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", ":50067")
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

	log.Printf("[%s] gRPC server listening on :50067", serviceName)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("[%s] Failed to serve: %v", serviceName, err)
	}
	<-ctx.Done()
	log.Printf("[%s] Shutdown complete", serviceName)
}
