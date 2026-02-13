// reservation-service: 예약 마이크로서비스
//
// 포트: gRPC :50066
// 의존: PostgreSQL(선택) — 미설정 시 인메모리 저장소 사용
//
// 기능:
// - 의료 시설 검색/조회
// - 예약 가능 시간대 조회
// - 예약 생성/조회/목록/취소
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
	"github.com/manpasik/backend/services/reservation-service/internal/handler"
	"github.com/manpasik/backend/services/reservation-service/internal/repository/memory"
	"github.com/manpasik/backend/services/reservation-service/internal/repository/postgres"
	"github.com/manpasik/backend/services/reservation-service/internal/service"
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

const serviceName = "reservation-service"

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

	// 저장소 초기화: PostgreSQL 또는 인메모리
	var facilityRepo service.FacilityRepository
	var slotRepo service.SlotRepository
	var reservationRepo service.ReservationRepository
	var doctorRepo service.DoctorRepository
	var regionRepo service.RegionRepository

	if _, dbHostSet := os.LookupEnv("DB_HOST"); dbHostSet && cfg.DB.Host != "" && cfg.DB.DBName != "" {
		connCtx, connCancel := context.WithTimeout(context.Background(), 5*time.Second)
		pool, poolErr := pgxpool.New(connCtx, cfg.DB.DSN())
		connCancel()
		if poolErr != nil {
			log.Printf("[%s] DB 풀 생성 실패, 인메모리 사용: %v", serviceName, poolErr)
			facilityRepo = memory.NewFacilityRepository()
			slotRepo = memory.NewSlotRepository()
			reservationRepo = memory.NewReservationRepository()
			doctorRepo = memory.NewDoctorRepository()
			regionRepo = memory.NewRegionRepository()
		} else {
			pingCtx, pingCancel := context.WithTimeout(context.Background(), 3*time.Second)
			if pingErr := pool.Ping(pingCtx); pingErr != nil {
				pingCancel()
				pool.Close()
				log.Printf("[%s] DB Ping 실패, 인메모리 사용: %v", serviceName, pingErr)
				facilityRepo = memory.NewFacilityRepository()
				slotRepo = memory.NewSlotRepository()
				reservationRepo = memory.NewReservationRepository()
				doctorRepo = memory.NewDoctorRepository()
				regionRepo = memory.NewRegionRepository()
			} else {
				pingCancel()
				defer pool.Close()
				facilityRepo = postgres.NewFacilityRepository(pool)
				slotRepo = postgres.NewSlotRepository(pool)
				reservationRepo = postgres.NewReservationRepository(pool)
				doctorRepo = postgres.NewDoctorRepository(pool)
				regionRepo = postgres.NewRegionRepository(pool)
				log.Printf("[%s] DB 연결됨: %s", serviceName, cfg.DB.DBName)
			}
		}
	} else {
		facilityRepo = memory.NewFacilityRepository()
		slotRepo = memory.NewSlotRepository()
		reservationRepo = memory.NewReservationRepository()
		doctorRepo = memory.NewDoctorRepository()
		regionRepo = memory.NewRegionRepository()
		log.Printf("[%s] 인메모리 저장소 사용", serviceName)
	}

	resSvc := service.NewReservationService(logger, facilityRepo, slotRepo, reservationRepo, doctorRepo, regionRepo)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.RequestIDInterceptor(),
			observability.UnaryServerInterceptor(metrics),
		),
	)

	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus(serviceName, healthpb.HealthCheckResponse_SERVING)

	resHandler := handler.NewReservationHandler(resSvc, logger)
	v1.RegisterReservationServiceServer(grpcServer, resHandler)

	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", ":50066")
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

	log.Printf("[%s] gRPC server listening on :50066", serviceName)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("[%s] Failed to serve: %v", serviceName, err)
	}
	<-ctx.Done()
	log.Printf("[%s] Shutdown complete", serviceName)
}
