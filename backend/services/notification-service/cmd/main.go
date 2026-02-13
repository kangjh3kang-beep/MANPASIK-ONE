// notification-service: 알림 마이크로서비스
//
// 포트: gRPC :50062
// 의존: PostgreSQL(선택) — 미설정 시 인메모리 저장소 사용
//
// 기능:
// - 푸시/이메일/SMS/인앱 알림 발송
// - 알림 목록 조회 / 읽음 처리
// - 알림 설정(선호도) 관리
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
	"github.com/manpasik/backend/services/notification-service/internal/handler"
	"github.com/manpasik/backend/services/notification-service/internal/push"
	"github.com/manpasik/backend/services/notification-service/internal/repository/memory"
	"github.com/manpasik/backend/services/notification-service/internal/repository/postgres"
	"github.com/manpasik/backend/services/notification-service/internal/service"
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

const serviceName = "notification-service"

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

	var notiRepo service.NotificationRepository
	var prefRepo service.PreferencesRepository
	if _, dbHostSet := os.LookupEnv("DB_HOST"); dbHostSet && cfg.DB.Host != "" && cfg.DB.DBName != "" {
		connCtx, connCancel := context.WithTimeout(context.Background(), 5*time.Second)
		pool, poolErr := pgxpool.New(connCtx, cfg.DB.DSN())
		connCancel()
		if poolErr != nil {
			log.Printf("[%s] DB 풀 생성 실패, 인메모리 사용: %v", serviceName, poolErr)
			notiRepo = memory.NewNotificationRepository()
			prefRepo = memory.NewPreferencesRepository()
		} else {
			pingCtx, pingCancel := context.WithTimeout(context.Background(), 3*time.Second)
			if pingErr := pool.Ping(pingCtx); pingErr != nil {
				pingCancel()
				pool.Close()
				log.Printf("[%s] DB Ping 실패, 인메모리 사용: %v", serviceName, pingErr)
				notiRepo = memory.NewNotificationRepository()
				prefRepo = memory.NewPreferencesRepository()
			} else {
				pingCancel()
				defer pool.Close()
				notiRepo = postgres.NewNotificationRepository(pool)
				prefRepo = postgres.NewPreferencesRepository(pool)
				log.Printf("[%s] DB 연결됨: %s", serviceName, cfg.DB.DBName)
			}
		}
	} else {
		notiRepo = memory.NewNotificationRepository()
		prefRepo = memory.NewPreferencesRepository()
		log.Printf("[%s] 인메모리 저장소 사용", serviceName)
	}

	notiSvc := service.NewNotificationService(logger, notiRepo, prefRepo)

	// FCM 푸시 알림 전송기: DB config 우선 → 환경변수 fallback
	fcmKey := config.LoadConfigWithFallback(nil, "fcm.server_key", "FCM_SERVER_KEY")
	fcmProject := config.LoadConfigWithFallback(nil, "fcm.project_id", "FIREBASE_PROJECT_ID")
	if fcmKey != "" {
		fcmClient := push.NewFCMClient(push.FCMConfig{
			ServerKey: fcmKey,
			ProjectID: fcmProject,
		})
		notiSvc.SetPushSender(fcmClient)
		log.Printf("[%s] FCM 푸시 전송기 활성화", serviceName)
	} else {
		notiSvc.SetPushSender(push.NewNoopPushSender())
		log.Printf("[%s] FCM 미설정, No-op 푸시 전송기 사용", serviceName)
	}

	// 이메일 전송기: 향후 SMTP 설정 시 활성화
	notiSvc.SetEmailSender(push.NewNoopEmailSender())

	// ConfigWatcher: FCM 설정 변경 시 핫리로드
	configWatcher := events.NewEventBusConfigWatcher(events.NewEventBus())
	_ = configWatcher.Watch(context.Background(), serviceName, func(key, newValue string) error {
		switch key {
		case "fcm.server_key":
			projectID := config.LoadConfigWithFallback(nil, "fcm.project_id", "FIREBASE_PROJECT_ID")
			notiSvc.SetPushSender(push.NewFCMClient(push.FCMConfig{
				ServerKey: newValue,
				ProjectID: projectID,
			}))
			log.Printf("[%s] FCM 핫리로드 완료 (config.changed)", serviceName)
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

	notiHandler := handler.NewNotificationHandler(notiSvc, logger)
	v1.RegisterNotificationServiceServer(grpcServer, notiHandler)

	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", ":50062")
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

	log.Printf("[%s] gRPC server listening on :50062", serviceName)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("[%s] Failed to serve: %v", serviceName, err)
	}
	<-ctx.Done()
	log.Printf("[%s] Shutdown complete", serviceName)
}
