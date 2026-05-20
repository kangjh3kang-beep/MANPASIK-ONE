// telemedicine-service: 원격진료 마이크로서비스
//
// 포트: gRPC :50065
// 의존: PostgreSQL(선택) — 미설정 시 인메모리 저장소 사용
//
// 기능:
// - 원격진료 상담 생성/조회/목록
// - 의사 매칭
// - 비디오 세션 시작/종료
// - 상담 평점
package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manpasik/backend/services/telemedicine-service/internal/handler"
	"github.com/manpasik/backend/services/telemedicine-service/internal/repository/memory"
	"github.com/manpasik/backend/services/telemedicine-service/internal/repository/postgres"
	"github.com/manpasik/backend/services/telemedicine-service/internal/service"
	"github.com/manpasik/backend/services/telemedicine-service/internal/webrtc"
	"github.com/manpasik/backend/shared/config"
	v1 "github.com/manpasik/backend/shared/gen/go/v1"
	"github.com/manpasik/backend/shared/tenancy"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

const serviceName = "telemedicine-service"

func main() {
	cfg := config.LoadFromEnv(serviceName)

	logger, err := zap.NewProduction()
	if err != nil {
		logger = zap.NewNop()
	}
	defer logger.Sync()

	log.Printf("[%s] Starting v%s...", serviceName, cfg.Version)
	log.Printf("[%s] gRPC port: %s", serviceName, cfg.GRPCPort)

	// ConsultationRepository, DoctorRepository, VideoSessionRepository: PostgreSQL 또는 인메모리
	var consultationRepo service.ConsultationRepository
	var doctorRepo service.DoctorRepository
	var videoSessionRepo service.VideoSessionRepository

	if _, dbHostSet := os.LookupEnv("DB_HOST"); dbHostSet && cfg.DB.Host != "" && cfg.DB.DBName != "" {
		connCtx, connCancel := context.WithTimeout(context.Background(), 5*time.Second)
		pool, poolErr := pgxpool.New(connCtx, cfg.DB.DSN())
		connCancel()
		if poolErr != nil {
			log.Printf("[%s] DB 풀 생성 실패, 인메모리 사용: %v", serviceName, poolErr)
			consultationRepo = memory.NewConsultationRepository()
			doctorRepo = memory.NewDoctorRepository()
			videoSessionRepo = memory.NewVideoSessionRepository()
		} else {
			pingCtx, pingCancel := context.WithTimeout(context.Background(), 3*time.Second)
			if pingErr := pool.Ping(pingCtx); pingErr != nil {
				pingCancel()
				pool.Close()
				log.Printf("[%s] DB Ping 실패, 인메모리 사용: %v", serviceName, pingErr)
				consultationRepo = memory.NewConsultationRepository()
				doctorRepo = memory.NewDoctorRepository()
				videoSessionRepo = memory.NewVideoSessionRepository()
			} else {
				pingCancel()
				defer pool.Close()
				consultationRepo = postgres.NewConsultationRepository(pool)
				doctorRepo = postgres.NewDoctorRepository(pool)
				videoSessionRepo = postgres.NewVideoSessionRepository(pool)
				log.Printf("[%s] DB 연결됨: %s", serviceName, cfg.DB.DBName)
			}
		}
	} else {
		consultationRepo = memory.NewConsultationRepository()
		doctorRepo = memory.NewDoctorRepository()
		videoSessionRepo = memory.NewVideoSessionRepository()
		log.Printf("[%s] 인메모리 저장소 사용", serviceName)
	}

	teleSvc := service.NewTelemedicineService(logger, consultationRepo, doctorRepo, videoSessionRepo)

	// WebRTC 제공자 초기화: AGORA_APP_ID 환경변수 설정 시 Agora 사용
	if appID := os.Getenv("AGORA_APP_ID"); appID != "" {
		appCert := os.Getenv("AGORA_APP_CERT")
		customerID := os.Getenv("AGORA_CUSTOMER_ID")
		secret := os.Getenv("AGORA_SECRET")
		baseURL := os.Getenv("AGORA_BASE_URL")
		provider := webrtc.NewAgoraProvider(webrtc.AgoraConfig{
			AppID:      appID,
			AppCert:    appCert,
			CustomerID: customerID,
			Secret:     secret,
			BaseURL:    baseURL,
		})
		teleSvc.SetWebRTCProvider(provider)
		log.Printf("[%s] Agora WebRTC 제공자 활성화", serviceName)
	} else {
		log.Printf("[%s] WebRTC 제공자 미설정, 폴백 사용", serviceName)
	}

	// 멀티테넌트 RBAC: TENANCY_ENFORCED=true 시에만 테넌시 인터셉터 활성화 (점진 도입).
	tenancyEngine := tenancy.NewPolicyEngine(tenancy.NewMemoryMembershipStore())
	serverOpts := tenancy.MaybeServerOptions(tenancyEngine, nil)
	if len(serverOpts) > 0 {
		log.Printf("[%s] 멀티테넌트 RBAC 인터셉터 활성화", serviceName)
	}

	grpcServer := grpc.NewServer(serverOpts...)

	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus(serviceName, healthpb.HealthCheckResponse_SERVING)

	teleHandler := handler.NewTelemedicineHandler(teleSvc, logger)
	v1.RegisterTelemedicineServiceServer(grpcServer, teleHandler)

	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", ":50065")
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

	log.Printf("[%s] gRPC server listening on :50065", serviceName)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("[%s] Failed to serve: %v", serviceName, err)
	}
	<-ctx.Done()
	log.Printf("[%s] Shutdown complete", serviceName)
}
