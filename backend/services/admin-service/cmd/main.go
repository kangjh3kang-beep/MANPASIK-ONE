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
	"github.com/manpasik/backend/shared/medical/llm"
	"github.com/manpasik/backend/shared/observability"
	"github.com/manpasik/backend/shared/ops/dashboard"
	"github.com/manpasik/backend/shared/tenancy"
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

	// 운영 대시보드: 헬스 어그리게이터 (관리자 노출)
	dashAgg := dashboard.NewAggregator()
	dashAgg.Register(dashboard.NewSimpleProvider("admin-service", func(_ context.Context) dashboard.ModuleSnapshot {
		return dashboard.ModuleSnapshot{Status: dashboard.StatusHealthy}
	}))
	// /ops 하위에 모듈 헬스 노출. /metrics(Prometheus 메트릭)는 기존 엔드포인트가 우선.
	dashHandler := dashboard.NewHTTPHandler(dashAgg)
	dashHandler.SetPathPrefix("/ops")

	// 멀티테넌트 멤버십 + 초대 REST API (Phase Z + AA).
	// DB_HOST 설정 시 PostgreSQL 영속화 자동 활성, 미설정 시 인메모리.
	tenancyMemStore, tenancyInvStore, tenancyPool := buildAdminTenancyStores(cfg, serviceName)
	if tenancyPool != nil {
		defer tenancyPool.Close()
	}
	tenancyEngine := tenancy.NewPolicyEngine(tenancyMemStore)
	invSvc, _ := tenancy.NewInvitationService(tenancyInvStore, tenancyMemStore, tenancy.InvitationServiceConfig{})
	invSvc.SetPolicyEngine(tenancyEngine)

	// Webhook dispatcher (Phase AL-1) — WEBHOOK_URL 설정 시 활성
	webhookDispatcher := buildWebhookDispatcher(serviceName)
	if webhookDispatcher != nil {
		invSvc.SetWebhookDispatcher(webhookDispatcher)
		webhookDispatcher.Start(context.Background())
		defer webhookDispatcher.Stop()
	}

	tenancyHTTP := tenancy.NewHTTPHandler(invSvc, tenancyMemStore, tenancyEngine)
	tenancyHTTP.SetPathPrefix("/api/v1")

	// Tenancy 운영 메트릭 → ops/dashboard 어그리게이터 등록 (Phase AD-1).
	tenantLister := tenancy.NewMembershipBackedTenantLister(tenancyMemStore, func() []string {
		// admin-service 는 모든 사용자 ID 를 알 수 없으므로 비어있는 목록 반환.
		// 실제 운영에서는 user-service 와 통합하여 활성 사용자 ID 주입.
		return nil
	})
	statsCollector := tenancy.NewStatsCollector(tenancyMemStore, tenancyInvStore, tenantLister)
	dashAgg.Register(dashboard.AdapterFromMetricsCollector("tenancy", statsCollector))

	// Webhook 통계 dashboard 등록 (Phase AL-1)
	if webhookDispatcher != nil {
		dashAgg.Register(dashboard.AdapterFromMetricsCollector("webhook", webhookDispatcher))
	}

	// Start observability HTTP server
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/metrics", metrics.PrometheusHandler())
		mux.HandleFunc("/health", healthCheck.Handler())
		dashHandler.RegisterRoutes(mux)
		tenancyHTTP.RegisterRoutes(mux)
		// 별도 통계 엔드포인트 — JSON 직접 노출
		mux.HandleFunc("/ops/tenancy/stats", func(w http.ResponseWriter, _ *http.Request) {
			tenancy.WriteStatsJSON(w, statsCollector.Collect())
		})
		// Prometheus exposition (Phase AE-1)
		mux.HandleFunc("/ops/tenancy/metrics", func(w http.ResponseWriter, _ *http.Request) {
			tenancy.WritePrometheusMetrics(w, statsCollector.Collect())
		})
		// LLM quota REST API (Phase AJ-3 + AK-2)
		// DB_HOST 설정 시 PostgresQuotaStore 자동 활성, 미설정 시 인메모리.
		quotaStore := buildQuotaStore(tenancyPool)
		quotaHandler := llm.NewQuotaHTTPHandler(quotaStore)
		quotaHandler.SetPathPrefix("/ops/tenancy")
		quotaHandler.RegisterRoutes(mux)
		// LLM Audit log (Phase AL-2 + AM-2) — admin-service 운영 가시성
		auditLog := buildAuditLog(serviceName, tenancyPool)
		_ = auditLog
		// PostgresAuditLog 만 통계 REST API 등록 (BatchAuditLog 는 wrapper)
		if tenancyPool != nil {
			pgxAdapter := llm.NewPgxAdapter(tenancyPool)
			pgAuditLog, _ := llm.NewPostgresAuditLog(pgxAdapter, pgxAdapter)
			auditHandler := llm.NewAuditHTTPHandler(pgAuditLog)
			auditHandler.SetPathPrefix("/ops/tenancy")
			auditHandler.RegisterRoutes(mux)
		}
		// Webhook Prometheus 메트릭 (Phase AL-1)
		mux.HandleFunc("/ops/tenancy/webhook/metrics", func(w http.ResponseWriter, _ *http.Request) {
			tenancy.WriteWebhookPrometheusMetrics(w, webhookDispatcher)
		})
		// Webhook DLQ REST API (Phase AO-2) — 재시도 한계 초과 이벤트 운영 관리
		dlqHandler := tenancy.NewWebhookDLQHandler(webhookDispatcher)
		dlqHandler.SetPathPrefix("/ops/tenancy")
		dlqHandler.RegisterRoutes(mux)
		// LLM Audit 실패 조회 (Phase AN-3) — Audit handler 가 위에 등록됨
		// /ops/tenancy/audit/failures?tenant=X&limit=N
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

// buildAuditLog 는 DB pool 시 PostgresAuditLog → BatchAuditLog wrapping (Phase AL-2).
//
// 미설정 시 인메모리 (개발용). 영속화는 운영 환경에서 LLM 호출 cost 분리 필수.
func buildAuditLog(serviceName string, pool *pgxpool.Pool) llm.TenancyAuditLog {
	if pool == nil {
		log.Printf("[%s] LLM Audit 인메모리 (DB pool 없음)", serviceName)
		return llm.NewMemoryAuditLog()
	}
	adapter := llm.NewPgxAdapter(pool)
	pgLog, err := llm.NewPostgresAuditLog(adapter, adapter)
	if err != nil {
		log.Printf("[%s] PostgresAuditLog 생성 실패: %v", serviceName, err)
		return llm.NewMemoryAuditLog()
	}
	// BatchAuditLog 로 wrapping (성능)
	batched := llm.NewBatchAuditLog(pgLog, llm.BatchAuditConfig{
		MaxSize:       1000,
		FlushInterval: 5 * time.Second,
		OnFlushError: func(err error) {
			log.Printf("[%s] LLM Audit flush 에러: %v", serviceName, err)
		},
	})
	batched.Start(context.Background())
	log.Printf("[%s] PostgresAuditLog + Batch flush 활성화", serviceName)
	return batched
}

// buildWebhookDispatcher 는 WEBHOOK_URL 환경변수가 설정되었으면 dispatcher 생성 (Phase AL-1).
// 미설정 시 nil 반환 (webhook 비활성).
func buildWebhookDispatcher(serviceName string) *tenancy.WebhookDispatcher {
	url := os.Getenv("WEBHOOK_URL")
	if url == "" {
		log.Printf("[%s] Webhook 비활성 (WEBHOOK_URL 미설정)", serviceName)
		return nil
	}
	mode := os.Getenv("WEBHOOK_MODE")
	if mode == "" {
		mode = "generic"
	}
	cfg := tenancy.WebhookConfig{
		URL:        url,
		Mode:       mode,
		MaxRetries: 3,
	}
	cfg.OnError = func(_ tenancy.Event, err error) {
		log.Printf("[%s] webhook 에러: %v", serviceName, err)
	}
	d, err := tenancy.NewWebhookDispatcher(cfg, nil)
	if err != nil {
		log.Printf("[%s] WebhookDispatcher 생성 실패: %v", serviceName, err)
		return nil
	}
	log.Printf("[%s] Webhook 활성 (mode=%s)", serviceName, mode)
	return d
}

// buildQuotaStore 는 DB pool 이 있으면 PostgresQuotaStore, 없으면 인메모리 사용 (Phase AK-2).
func buildQuotaStore(pool *pgxpool.Pool) llm.QuotaStore {
	if pool == nil {
		log.Printf("[%s] Quota 인메모리 store 사용 (DB pool 없음)", serviceName)
		return llm.NewMemoryQuotaStore()
	}
	adapter := llm.NewPgxAdapter(pool)
	store, err := llm.NewPostgresQuotaStore(adapter, adapter)
	if err != nil {
		log.Printf("[%s] PostgresQuotaStore 생성 실패, 인메모리 fallback: %v", serviceName, err)
		return llm.NewMemoryQuotaStore()
	}
	log.Printf("[%s] PostgresQuotaStore 활성화", serviceName)
	return store
}

// buildAdminTenancyStores 는 admin-service 의 tenancy store 를 DB_HOST 설정에
// 따라 선택. gateway 의 buildTenancyStores 와 동일 로직.
func buildAdminTenancyStores(cfg *config.ServiceConfig, serviceName string) (
	tenancy.MembershipStore, tenancy.InvitationStore, *pgxpool.Pool) {
	dbHostSet := false
	if _, ok := os.LookupEnv("DB_HOST"); ok {
		dbHostSet = true
	}
	if !dbHostSet || cfg.DB.Host == "" || cfg.DB.DBName == "" {
		log.Printf("[%s] Tenancy 인메모리 store 사용 (DB_HOST 미설정)", serviceName)
		return tenancy.NewMemoryMembershipStore(),
			tenancy.NewMemoryInvitationStore(),
			nil
	}

	connCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	pool, err := pgxpool.New(connCtx, cfg.DB.DSN())
	cancel()
	if err != nil {
		log.Printf("[%s] Tenancy DB 풀 생성 실패, 인메모리: %v", serviceName, err)
		return tenancy.NewMemoryMembershipStore(),
			tenancy.NewMemoryInvitationStore(),
			nil
	}
	pingCtx, pingCancel := context.WithTimeout(context.Background(), 3*time.Second)
	if err := pool.Ping(pingCtx); err != nil {
		pingCancel()
		pool.Close()
		log.Printf("[%s] Tenancy DB Ping 실패, 인메모리: %v", serviceName, err)
		return tenancy.NewMemoryMembershipStore(),
			tenancy.NewMemoryInvitationStore(),
			nil
	}
	pingCancel()

	adapter := tenancy.NewPgxAdapter(pool)
	memStore, err := tenancy.NewPostgresMembershipStore(adapter)
	if err != nil {
		pool.Close()
		log.Printf("[%s] PostgresMembershipStore 생성 실패: %v", serviceName, err)
		return tenancy.NewMemoryMembershipStore(),
			tenancy.NewMemoryInvitationStore(),
			nil
	}
	invStore, err := tenancy.NewPostgresInvitationStore(adapter)
	if err != nil {
		pool.Close()
		log.Printf("[%s] PostgresInvitationStore 생성 실패: %v", serviceName, err)
		return memStore, tenancy.NewMemoryInvitationStore(), nil
	}
	log.Printf("[%s] Tenancy PostgreSQL store 활성화 (db=%s)", serviceName, cfg.DB.DBName)
	return memStore, invStore, pool
}
