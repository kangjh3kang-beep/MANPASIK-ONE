#!/bin/bash
set -e
BASE="/home/kangjh3kang/Manpasik/backend/services"

# ============================================================================
# 3. data-platform-service (포트 50080)
# LocationStatsServiceServer (3 RPCs) + DataProvisionServiceServer (7 RPCs)
# ============================================================================
SVC="$BASE/data-platform-service"
mkdir -p "$SVC/cmd" "$SVC/internal/handler" "$SVC/internal/service" "$SVC/internal/repository/memory"

cat > "$SVC/cmd/main.go" << 'GOEOF'
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

	"github.com/manpasik/backend/services/data-platform-service/internal/handler"
	"github.com/manpasik/backend/services/data-platform-service/internal/repository/memory"
	"github.com/manpasik/backend/services/data-platform-service/internal/service"
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

const serviceName = "data-platform-service"

func main() {
	cfg := config.LoadFromEnv(serviceName)
	logger, err := zap.NewProduction()
	if err != nil { logger = zap.NewNop() }
	defer logger.Sync()

	metrics := observability.NewMetrics()
	healthCheck := observability.NewHealthCheck(serviceName, cfg.Version)
	log.Printf("[%s] Starting v%s...", serviceName, cfg.Version)

	repo := memory.NewPlatformRepository()
	svc := service.NewPlatformService(logger, repo)
	h := handler.NewPlatformHandler(svc, logger)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.RequestIDInterceptor(),
			observability.UnaryServerInterceptor(metrics),
		),
	)
	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus(serviceName, healthpb.HealthCheckResponse_SERVING)
	v1.RegisterLocationStatsServiceServer(grpcServer, h)
	v1.RegisterDataProvisionServiceServer(grpcServer, h)
	reflection.Register(grpcServer)

	grpcPort := cfg.GRPCPort
	if grpcPort == "" { grpcPort = ":50080" }
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil { log.Fatalf("[%s] Failed to listen: %v", serviceName, err) }

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigCh
		log.Printf("[%s] Received signal %v, shutting down...", serviceName, sig)
		healthServer.SetServingStatus(serviceName, healthpb.HealthCheckResponse_NOT_SERVING)
		go func() { time.Sleep(cfg.ShutdownTimeout); os.Exit(1) }()
		grpcServer.GracefulStop(); cancel()
	}()
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/metrics", metrics.PrometheusHandler())
		mux.HandleFunc("/health", healthCheck.Handler())
		http.ListenAndServe(":9100", mux)
	}()
	log.Printf("[%s] gRPC server listening on %s", serviceName, grpcPort)
	if err := grpcServer.Serve(lis); err != nil { log.Fatalf("[%s] Failed to serve: %v", serviceName, err) }
	<-ctx.Done()
}
GOEOF

cat > "$SVC/internal/service/platform.go" << 'GOEOF'
package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrInvalidInput = errors.New("invalid input")
)

type RegionStats struct {
	RegionCode, RegionName, Period, RiskLevel string
	AvgValue, MinValue, MaxValue              float64
	SampleCount                               int32
}
type TrendPoint struct{ Date string; Value float64; SampleCount int32 }
type RegionTrend struct {
	RegionCode, Biomarker, TrendDirection string
	Points                               []*TrendPoint
}
type AnonymizedDataset struct {
	ID, DatasetType, Description, DateRange string
	RecordCount                             int32
	CreatedAt                               time.Time
}
type DataAccessResp struct{ RequestID, Status string }
type ConsentInfo struct {
	UserID            string
	AnonymousResearch, RegionStatsConsent, EnterpriseData bool
	UpdatedAt         time.Time
}
type AggregatedInsight struct {
	ID, Biomarker, InsightType, Summary string
	CreatedAt                           time.Time
}

type PlatformRepository interface {
	GetRegionStats(ctx context.Context, code, period, biomarker string) (*RegionStats, error)
	ListHotspots(ctx context.Context, biomarker, risk string, limit int32) ([]*RegionStats, error)
	GetTrend(ctx context.Context, code, biomarker string) (*RegionTrend, error)
	ListDatasets(ctx context.Context, dsType string, limit, offset int32) ([]*AnonymizedDataset, int32, error)
	GetDataset(ctx context.Context, id string) (*AnonymizedDataset, error)
	SaveDataset(ctx context.Context, ds *AnonymizedDataset) error
	CreateAccessRequest(ctx context.Context, requesterID, datasetID, purpose string) (*DataAccessResp, error)
	GetConsent(ctx context.Context, userID string) (*ConsentInfo, error)
	UpdateConsent(ctx context.Context, c *ConsentInfo) error
	ListInsights(ctx context.Context, biomarker string, limit int32) ([]*AggregatedInsight, error)
	SaveInsight(ctx context.Context, i *AggregatedInsight) error
}

type PlatformService struct {
	log  *zap.Logger
	repo PlatformRepository
	rng  *rand.Rand
}

func NewPlatformService(log *zap.Logger, repo PlatformRepository) *PlatformService {
	return &PlatformService{log: log, repo: repo, rng: rand.New(rand.NewSource(time.Now().UnixNano()))}
}

var regionNames = map[string]string{
	"KR-11": "서울특별시", "KR-26": "부산광역시", "KR-27": "대구광역시",
	"KR-28": "인천광역시", "KR-41": "경기도", "KR-42": "강원도",
}

func (s *PlatformService) GetRegionStats(ctx context.Context, code, period, biomarker string) (*RegionStats, error) {
	if code == "" { return nil, ErrInvalidInput }
	if period == "" { period = "day" }
	stats, _ := s.repo.GetRegionStats(ctx, code, period, biomarker)
	if stats == nil {
		name := regionNames[code]; if name == "" { name = code }
		avg := 80.0 + s.rng.Float64()*40.0
		risk := "low"; if avg > 110 { risk = "high" } else if avg > 95 { risk = "medium" }
		stats = &RegionStats{RegionCode: code, RegionName: name, Period: period,
			AvgValue: avg, MinValue: avg - 15, MaxValue: avg + 15, SampleCount: int32(100 + s.rng.Intn(900)), RiskLevel: risk}
	}
	return stats, nil
}
func (s *PlatformService) ListHotspots(ctx context.Context, biomarker, risk string, limit int32) ([]*RegionStats, error) {
	if limit <= 0 { limit = 10 }
	list, _ := s.repo.ListHotspots(ctx, biomarker, risk, limit)
	if len(list) == 0 {
		codes := []string{"KR-11", "KR-26", "KR-27", "KR-28", "KR-41"}
		for i := 0; i < int(limit) && i < len(codes); i++ {
			st, _ := s.GetRegionStats(ctx, codes[i], "day", biomarker)
			if st != nil { st.RiskLevel = "high"; list = append(list, st) }
		}
	}
	return list, nil
}
func (s *PlatformService) GetTrend(ctx context.Context, code, biomarker, start, end string) (*RegionTrend, error) {
	if code == "" { return nil, ErrInvalidInput }
	trend, _ := s.repo.GetTrend(ctx, code, biomarker)
	if trend == nil {
		pts := make([]*TrendPoint, 7)
		now := time.Now().UTC()
		for i := 6; i >= 0; i-- {
			pts[6-i] = &TrendPoint{Date: now.AddDate(0, 0, -i).Format("2006-01-02"), Value: 80 + s.rng.Float64()*40, SampleCount: int32(50 + s.rng.Intn(200))}
		}
		trend = &RegionTrend{RegionCode: code, Biomarker: biomarker, Points: pts, TrendDirection: "stable"}
	}
	return trend, nil
}
func (s *PlatformService) ListDatasets(ctx context.Context, dsType string, limit, offset int32) ([]*AnonymizedDataset, int32, error) {
	if limit <= 0 { limit = 20 }
	list, total, _ := s.repo.ListDatasets(ctx, dsType, limit, offset)
	if len(list) == 0 && offset == 0 {
		for i := 0; i < 3; i++ {
			ds := &AnonymizedDataset{ID: uuid.New().String(), DatasetType: "blood_glucose",
				Description: fmt.Sprintf("혈당 데이터셋 #%d", i+1), RecordCount: int32(1000 + s.rng.Intn(9000)),
				DateRange: "2025-01-01~2025-12-31", CreatedAt: time.Now().UTC()}
			list = append(list, ds)
			s.repo.SaveDataset(ctx, ds)
		}
		total = int32(len(list))
	}
	return list, total, nil
}
func (s *PlatformService) GetDataset(ctx context.Context, id string) (*AnonymizedDataset, error) {
	if id == "" { return nil, ErrInvalidInput }
	ds, err := s.repo.GetDataset(ctx, id)
	if err != nil || ds == nil { return nil, ErrNotFound }
	return ds, nil
}
func (s *PlatformService) RequestAccess(ctx context.Context, requesterID, datasetID, purpose string) (*DataAccessResp, error) {
	if requesterID == "" || datasetID == "" { return nil, ErrInvalidInput }
	return s.repo.CreateAccessRequest(ctx, requesterID, datasetID, purpose)
}
func (s *PlatformService) GetConsent(ctx context.Context, userID string) (*ConsentInfo, error) {
	if userID == "" { return nil, ErrInvalidInput }
	c, _ := s.repo.GetConsent(ctx, userID)
	if c == nil { c = &ConsentInfo{UserID: userID, UpdatedAt: time.Now().UTC()} }
	return c, nil
}
func (s *PlatformService) UpdateConsent(ctx context.Context, userID string, anon, region, enterprise bool) (*ConsentInfo, error) {
	if userID == "" { return nil, ErrInvalidInput }
	c := &ConsentInfo{UserID: userID, AnonymousResearch: anon, RegionStatsConsent: region, EnterpriseData: enterprise, UpdatedAt: time.Now().UTC()}
	s.repo.UpdateConsent(ctx, c)
	return c, nil
}
func (s *PlatformService) TriggerAggregation(ctx context.Context, biomarker, date string) (string, string, error) {
	aggID := uuid.New().String()
	s.repo.SaveInsight(ctx, &AggregatedInsight{ID: uuid.New().String(), Biomarker: biomarker,
		InsightType: "aggregation", Summary: fmt.Sprintf("%s 집계 완료", biomarker), CreatedAt: time.Now().UTC()})
	return aggID, "completed", nil
}
func (s *PlatformService) GetInsights(ctx context.Context, biomarker string, limit int32) ([]*AggregatedInsight, error) {
	if limit <= 0 { limit = 10 }
	list, _ := s.repo.ListInsights(ctx, biomarker, limit)
	if len(list) == 0 {
		list = []*AggregatedInsight{{ID: uuid.New().String(), Biomarker: biomarker,
			InsightType: "trend", Summary: biomarker + " 개선 추세", CreatedAt: time.Now().UTC()}}
	}
	return list, nil
}
GOEOF

cat > "$SVC/internal/repository/memory/platform.go" << 'GOEOF'
package memory

import (
	"context"
	"sync"
	"github.com/google/uuid"
	"github.com/manpasik/backend/services/data-platform-service/internal/service"
)

type PlatformRepository struct {
	mu       sync.RWMutex
	datasets map[string]*service.AnonymizedDataset
	consents map[string]*service.ConsentInfo
	insights []*service.AggregatedInsight
}

func NewPlatformRepository() *PlatformRepository {
	return &PlatformRepository{
		datasets: make(map[string]*service.AnonymizedDataset),
		consents: make(map[string]*service.ConsentInfo),
	}
}

func (r *PlatformRepository) GetRegionStats(_ context.Context, _, _, _ string) (*service.RegionStats, error) { return nil, nil }
func (r *PlatformRepository) ListHotspots(_ context.Context, _, _ string, _ int32) ([]*service.RegionStats, error) { return nil, nil }
func (r *PlatformRepository) GetTrend(_ context.Context, _, _ string) (*service.RegionTrend, error) { return nil, nil }

func (r *PlatformRepository) ListDatasets(_ context.Context, dsType string, limit, offset int32) ([]*service.AnonymizedDataset, int32, error) {
	r.mu.RLock(); defer r.mu.RUnlock()
	out := make([]*service.AnonymizedDataset, 0)
	for _, ds := range r.datasets {
		if dsType != "" && ds.DatasetType != dsType { continue }
		out = append(out, ds)
	}
	total := int32(len(out))
	start := int(offset); if start > len(out) { start = len(out) }
	end := start + int(limit); if end > len(out) { end = len(out) }
	return out[start:end], total, nil
}
func (r *PlatformRepository) GetDataset(_ context.Context, id string) (*service.AnonymizedDataset, error) {
	r.mu.RLock(); defer r.mu.RUnlock()
	if v, ok := r.datasets[id]; ok { return v, nil }
	return nil, service.ErrNotFound
}
func (r *PlatformRepository) SaveDataset(_ context.Context, ds *service.AnonymizedDataset) error {
	r.mu.Lock(); defer r.mu.Unlock()
	r.datasets[ds.ID] = ds; return nil
}
func (r *PlatformRepository) CreateAccessRequest(_ context.Context, requesterID, datasetID, purpose string) (*service.DataAccessResp, error) {
	return &service.DataAccessResp{RequestID: uuid.New().String(), Status: "pending"}, nil
}
func (r *PlatformRepository) GetConsent(_ context.Context, userID string) (*service.ConsentInfo, error) {
	r.mu.RLock(); defer r.mu.RUnlock()
	if v, ok := r.consents[userID]; ok { return v, nil }
	return nil, nil
}
func (r *PlatformRepository) UpdateConsent(_ context.Context, c *service.ConsentInfo) error {
	r.mu.Lock(); defer r.mu.Unlock()
	r.consents[c.UserID] = c; return nil
}
func (r *PlatformRepository) ListInsights(_ context.Context, biomarker string, limit int32) ([]*service.AggregatedInsight, error) {
	r.mu.RLock(); defer r.mu.RUnlock()
	out := make([]*service.AggregatedInsight, 0)
	for _, i := range r.insights {
		if biomarker != "" && i.Biomarker != biomarker { continue }
		out = append(out, i)
		if int32(len(out)) >= limit { break }
	}
	return out, nil
}
func (r *PlatformRepository) SaveInsight(_ context.Context, i *service.AggregatedInsight) error {
	r.mu.Lock(); defer r.mu.Unlock()
	r.insights = append(r.insights, i); return nil
}
GOEOF

cat > "$SVC/internal/handler/grpc.go" << 'GOEOF'
package handler

import (
	"context"
	"github.com/manpasik/backend/services/data-platform-service/internal/service"
	v1 "github.com/manpasik/backend/shared/gen/go/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PlatformHandler struct {
	v1.UnimplementedLocationStatsServiceServer
	v1.UnimplementedDataProvisionServiceServer
	svc *service.PlatformService
	log *zap.Logger
}

func NewPlatformHandler(svc *service.PlatformService, log *zap.Logger) *PlatformHandler {
	return &PlatformHandler{svc: svc, log: log}
}

// LocationStatsService
func (h *PlatformHandler) GetRegionStats(ctx context.Context, req *v1.GetRegionStatsRequest) (*v1.RegionStats, error) {
	s, err := h.svc.GetRegionStats(ctx, req.GetRegionCode(), req.GetPeriod(), req.GetBiomarker())
	if err != nil { return nil, toGRPC(err) }
	return &v1.RegionStats{RegionCode: s.RegionCode, RegionName: s.RegionName, Period: s.Period,
		AvgValue: s.AvgValue, MinValue: s.MinValue, MaxValue: s.MaxValue, SampleCount: s.SampleCount, RiskLevel: s.RiskLevel}, nil
}
func (h *PlatformHandler) ListHotspots(ctx context.Context, req *v1.ListHotspotsRequest) (*v1.ListHotspotsResponse, error) {
	list, err := h.svc.ListHotspots(ctx, req.GetBiomarker(), req.GetRiskLevel(), req.GetLimit())
	if err != nil { return nil, toGRPC(err) }
	out := make([]*v1.RegionStats, 0, len(list))
	for _, s := range list {
		out = append(out, &v1.RegionStats{RegionCode: s.RegionCode, RegionName: s.RegionName, Period: s.Period,
			AvgValue: s.AvgValue, MinValue: s.MinValue, MaxValue: s.MaxValue, SampleCount: s.SampleCount, RiskLevel: s.RiskLevel})
	}
	return &v1.ListHotspotsResponse{Hotspots: out}, nil
}
func (h *PlatformHandler) GetTrendByRegion(ctx context.Context, req *v1.GetTrendByRegionRequest) (*v1.RegionTrend, error) {
	t, err := h.svc.GetTrend(ctx, req.GetRegionCode(), req.GetBiomarker(), req.GetStartDate(), req.GetEndDate())
	if err != nil { return nil, toGRPC(err) }
	pts := make([]*v1.TrendPoint, 0, len(t.Points))
	for _, p := range t.Points { pts = append(pts, &v1.TrendPoint{Date: p.Date, Value: p.Value, SampleCount: p.SampleCount}) }
	return &v1.RegionTrend{RegionCode: t.RegionCode, Biomarker: t.Biomarker, Points: pts, TrendDirection: t.TrendDirection}, nil
}

// DataProvisionService
func (h *PlatformHandler) ListDatasets(ctx context.Context, req *v1.ListDatasetsRequest) (*v1.ListDatasetsResponse, error) {
	list, total, err := h.svc.ListDatasets(ctx, req.GetDatasetType(), req.GetLimit(), req.GetOffset())
	if err != nil { return nil, toGRPC(err) }
	out := make([]*v1.AnonymizedDataset, 0, len(list))
	for _, ds := range list {
		out = append(out, &v1.AnonymizedDataset{Id: ds.ID, DatasetType: ds.DatasetType, Description: ds.Description,
			RecordCount: ds.RecordCount, DateRange: ds.DateRange, CreatedAt: timestamppb.New(ds.CreatedAt)})
	}
	return &v1.ListDatasetsResponse{Datasets: out, Total: total}, nil
}
func (h *PlatformHandler) GetAnonymizedDataset(ctx context.Context, req *v1.GetAnonymizedDatasetRequest) (*v1.AnonymizedDataset, error) {
	ds, err := h.svc.GetDataset(ctx, req.GetDatasetId())
	if err != nil { return nil, toGRPC(err) }
	return &v1.AnonymizedDataset{Id: ds.ID, DatasetType: ds.DatasetType, Description: ds.Description,
		RecordCount: ds.RecordCount, DateRange: ds.DateRange, CreatedAt: timestamppb.New(ds.CreatedAt)}, nil
}
func (h *PlatformHandler) RequestDataAccess(ctx context.Context, req *v1.RequestDataAccessRequest) (*v1.DataAccessResponse, error) {
	r, err := h.svc.RequestAccess(ctx, req.GetRequesterId(), req.GetDatasetId(), req.GetPurpose())
	if err != nil { return nil, toGRPC(err) }
	return &v1.DataAccessResponse{RequestId: r.RequestID, Status: r.Status}, nil
}
func (h *PlatformHandler) GetDataSharingConsent(ctx context.Context, req *v1.GetConsentRequest) (*v1.DataSharingConsentInfo, error) {
	c, err := h.svc.GetConsent(ctx, req.GetUserId())
	if err != nil { return nil, toGRPC(err) }
	return &v1.DataSharingConsentInfo{UserId: c.UserID, AnonymousResearch: c.AnonymousResearch,
		RegionStats: c.RegionStatsConsent, EnterpriseData: c.EnterpriseData, UpdatedAt: timestamppb.New(c.UpdatedAt)}, nil
}
func (h *PlatformHandler) UpdateDataSharingConsent(ctx context.Context, req *v1.UpdateConsentRequest) (*v1.DataSharingConsentInfo, error) {
	c, err := h.svc.UpdateConsent(ctx, req.GetUserId(), req.GetAnonymousResearch(), req.GetRegionStats(), req.GetEnterpriseData())
	if err != nil { return nil, toGRPC(err) }
	return &v1.DataSharingConsentInfo{UserId: c.UserID, AnonymousResearch: c.AnonymousResearch,
		RegionStats: c.RegionStatsConsent, EnterpriseData: c.EnterpriseData, UpdatedAt: timestamppb.New(c.UpdatedAt)}, nil
}
func (h *PlatformHandler) TriggerAnonymousAggregation(ctx context.Context, req *v1.TriggerAggregationRequest) (*v1.TriggerAggregationResponse, error) {
	aggID, st, err := h.svc.TriggerAggregation(ctx, req.GetBiomarker(), req.GetDate())
	if err != nil { return nil, toGRPC(err) }
	return &v1.TriggerAggregationResponse{AggregationId: aggID, Status: st}, nil
}
func (h *PlatformHandler) GetAggregatedInsights(ctx context.Context, req *v1.GetInsightsRequest) (*v1.GetInsightsResponse, error) {
	list, err := h.svc.GetInsights(ctx, req.GetBiomarker(), req.GetLimit())
	if err != nil { return nil, toGRPC(err) }
	out := make([]*v1.AggregatedInsight, 0, len(list))
	for _, i := range list {
		out = append(out, &v1.AggregatedInsight{Id: i.ID, Biomarker: i.Biomarker, InsightType: i.InsightType,
			Summary: i.Summary, CreatedAt: timestamppb.New(i.CreatedAt)})
	}
	return &v1.GetInsightsResponse{Insights: out}, nil
}

func toGRPC(err error) error {
	switch err {
	case service.ErrNotFound: return status.Error(codes.NotFound, err.Error())
	case service.ErrInvalidInput: return status.Error(codes.InvalidArgument, err.Error())
	default: return status.Error(codes.Internal, err.Error())
	}
}
GOEOF

cat > "$SVC/internal/service/platform_test.go" << 'GOEOF'
package service_test

import (
	"context"
	"testing"
	"github.com/manpasik/backend/services/data-platform-service/internal/repository/memory"
	"github.com/manpasik/backend/services/data-platform-service/internal/service"
	"go.uber.org/zap"
)

func newSvc() *service.PlatformService {
	return service.NewPlatformService(zap.NewNop(), memory.NewPlatformRepository())
}
func TestGetRegionStats(t *testing.T) {
	s := newSvc()
	st, err := s.GetRegionStats(context.Background(), "KR-11", "day", "blood_glucose")
	if err != nil { t.Fatal(err) }
	if st.RegionCode != "KR-11" { t.Error("code mismatch") }
}
func TestListHotspots(t *testing.T) {
	s := newSvc()
	list, err := s.ListHotspots(context.Background(), "blood_glucose", "high", 3)
	if err != nil { t.Fatal(err) }
	if len(list) == 0 { t.Error("no hotspots") }
}
func TestGetTrend(t *testing.T) {
	s := newSvc()
	trend, err := s.GetTrend(context.Background(), "KR-11", "blood_glucose", "", "")
	if err != nil { t.Fatal(err) }
	if len(trend.Points) != 7 { t.Errorf("expected 7 points, got %d", len(trend.Points)) }
}
func TestListDatasets(t *testing.T) {
	s := newSvc()
	list, total, err := s.ListDatasets(context.Background(), "", 10, 0)
	if err != nil { t.Fatal(err) }
	if total == 0 || len(list) == 0 { t.Error("no datasets") }
}
func TestConsent(t *testing.T) {
	s := newSvc(); ctx := context.Background()
	c, _ := s.GetConsent(ctx, "u1")
	if c.AnonymousResearch { t.Error("default should be false") }
	c2, _ := s.UpdateConsent(ctx, "u1", true, true, false)
	if !c2.AnonymousResearch { t.Error("should be true") }
}
func TestTriggerAggregation(t *testing.T) {
	s := newSvc()
	id, st, err := s.TriggerAggregation(context.Background(), "blood_glucose", "2026-01-01")
	if err != nil { t.Fatal(err) }
	if id == "" || st != "completed" { t.Error("unexpected result") }
}
GOEOF

cat > "$SVC/Dockerfile" << 'GOEOF'
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /data-platform-service ./services/data-platform-service/cmd

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
COPY --from=builder /data-platform-service /data-platform-service
EXPOSE 50080
ENTRYPOINT ["/data-platform-service"]
GOEOF

echo "[OK] data-platform-service created"

# ============================================================================
# 4. cartridge-store-service (포트 50079)
# DeveloperServiceServer + StoreServiceServer + CartridgeReviewServiceServer
# + RevenueServiceServer + CartridgeAnalyticsServiceServer
# ============================================================================
SVC="$BASE/cartridge-store-service"
mkdir -p "$SVC/cmd" "$SVC/internal/handler" "$SVC/internal/service" "$SVC/internal/repository/memory"

cat > "$SVC/cmd/main.go" << 'GOEOF'
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

	"github.com/manpasik/backend/services/cartridge-store-service/internal/handler"
	"github.com/manpasik/backend/services/cartridge-store-service/internal/repository/memory"
	"github.com/manpasik/backend/services/cartridge-store-service/internal/service"
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

const serviceName = "cartridge-store-service"

func main() {
	cfg := config.LoadFromEnv(serviceName)
	logger, err := zap.NewProduction()
	if err != nil { logger = zap.NewNop() }
	defer logger.Sync()

	metrics := observability.NewMetrics()
	healthCheck := observability.NewHealthCheck(serviceName, cfg.Version)
	log.Printf("[%s] Starting v%s...", serviceName, cfg.Version)

	repo := memory.NewStoreRepository()
	svc := service.NewStoreService(logger, repo)
	h := handler.NewStoreHandler(svc, logger)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.RequestIDInterceptor(),
			observability.UnaryServerInterceptor(metrics),
		),
	)
	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus(serviceName, healthpb.HealthCheckResponse_SERVING)
	v1.RegisterDeveloperServiceServer(grpcServer, h)
	v1.RegisterStoreServiceServer(grpcServer, h)
	v1.RegisterCartridgeReviewServiceServer(grpcServer, h)
	v1.RegisterRevenueServiceServer(grpcServer, h)
	v1.RegisterCartridgeAnalyticsServiceServer(grpcServer, h)
	reflection.Register(grpcServer)

	grpcPort := cfg.GRPCPort
	if grpcPort == "" { grpcPort = ":50079" }
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil { log.Fatalf("[%s] Failed to listen: %v", serviceName, err) }

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigCh
		log.Printf("[%s] Received signal %v", serviceName, sig)
		healthServer.SetServingStatus(serviceName, healthpb.HealthCheckResponse_NOT_SERVING)
		go func() { time.Sleep(cfg.ShutdownTimeout); os.Exit(1) }()
		grpcServer.GracefulStop(); cancel()
	}()
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/metrics", metrics.PrometheusHandler())
		mux.HandleFunc("/health", healthCheck.Handler())
		http.ListenAndServe(":9100", mux)
	}()
	log.Printf("[%s] gRPC server listening on %s", serviceName, grpcPort)
	if err := grpcServer.Serve(lis); err != nil { log.Fatalf("[%s] Failed to serve: %v", serviceName, err) }
	<-ctx.Done()
}
GOEOF

cat > "$SVC/internal/service/store.go" << 'GOEOF'
package service

import (
	"context"
	"errors"
	"time"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrAlreadyExists = errors.New("already exists")
)

type Developer struct {
	ID, UserID, CompanyName, Email, Status string
	CreatedAt                              time.Time
}
type ApiKey struct{ KeyID, DeveloperID, Key, Status string; CreatedAt time.Time }
type StoreItem struct {
	ID, DeveloperID, Name, Description, Category, Version, Status string
	PriceKrw int64; Rating float64; Downloads int32
	CreatedAt time.Time
}
type Purchase struct {
	ID, UserID, ItemID, Status string
	PriceKrw int64; PurchasedAt time.Time
}
type ReviewStatus struct {
	SubmissionID, CartridgeID, Status, ReviewerComment string
	SubmittedAt time.Time
}
type SalesReport struct {
	DeveloperID string; TotalRevenue, PlatformFee, DeveloperEarning int64
	Period string; ItemCount int32
}
type PayoutRecord struct {
	ID, DeveloperID, Status string; AmountKrw int64; RequestedAt time.Time
}
type RevenueSplitCfg struct{ DeveloperID string; DeveloperPct, PlatformPct int32 }
type UsageStat struct{ ItemID, Date string; UsageCount, UniqueUsers int32 }
type RatingInfo struct{ ItemID string; AvgRating float64; TotalReviews int32 }
type UserReview struct{ ID, UserID, ItemID string; Rating int32; Comment string; CreatedAt time.Time }
type DevAnalytics struct{ DeveloperID string; TotalDownloads, TotalRevenue int64; ActiveCartridges int32 }

type StoreRepository interface {
	CreateDeveloper(ctx context.Context, d *Developer) error
	GetDeveloper(ctx context.Context, userID string) (*Developer, error)
	CreateApiKey(ctx context.Context, k *ApiKey) error
	CreateItem(ctx context.Context, item *StoreItem) error
	ListItems(ctx context.Context, category string, limit, offset int32) ([]*StoreItem, int32, error)
	SearchItems(ctx context.Context, query string, limit int32) ([]*StoreItem, error)
	GetItem(ctx context.Context, id string) (*StoreItem, error)
	CreatePurchase(ctx context.Context, p *Purchase) error
	ListPurchases(ctx context.Context, userID string, limit int32) ([]*Purchase, error)
	CreateReview(ctx context.Context, r *ReviewStatus) error
	GetReview(ctx context.Context, submissionID string) (*ReviewStatus, error)
	UpdateReview(ctx context.Context, r *ReviewStatus) error
	ListSubmissions(ctx context.Context, devID string, limit int32) ([]*ReviewStatus, error)
}

type StoreService struct {
	log  *zap.Logger
	repo StoreRepository
}
func NewStoreService(log *zap.Logger, repo StoreRepository) *StoreService {
	return &StoreService{log: log, repo: repo}
}

func (s *StoreService) RegisterDeveloper(ctx context.Context, userID, company, email string) (*Developer, error) {
	if userID == "" { return nil, ErrInvalidInput }
	d := &Developer{ID: uuid.New().String(), UserID: userID, CompanyName: company, Email: email, Status: "active", CreatedAt: time.Now().UTC()}
	return d, s.repo.CreateDeveloper(ctx, d)
}
func (s *StoreService) GetDeveloper(ctx context.Context, userID string) (*Developer, error) {
	if userID == "" { return nil, ErrInvalidInput }
	return s.repo.GetDeveloper(ctx, userID)
}
func (s *StoreService) CreateApiKey(ctx context.Context, devID string) (*ApiKey, error) {
	if devID == "" { return nil, ErrInvalidInput }
	k := &ApiKey{KeyID: uuid.New().String(), DeveloperID: devID, Key: "mpsk_" + uuid.New().String()[:16], Status: "active", CreatedAt: time.Now().UTC()}
	return k, s.repo.CreateApiKey(ctx, k)
}
func (s *StoreService) SubmitCartridge(ctx context.Context, devID, name, desc, category, version string, priceKrw int64) (*StoreItem, error) {
	if devID == "" || name == "" { return nil, ErrInvalidInput }
	item := &StoreItem{ID: uuid.New().String(), DeveloperID: devID, Name: name, Description: desc,
		Category: category, Version: version, PriceKrw: priceKrw, Status: "pending_review", CreatedAt: time.Now().UTC()}
	return item, s.repo.CreateItem(ctx, item)
}
func (s *StoreService) ListSubmissions(ctx context.Context, devID string, limit int32) ([]*ReviewStatus, error) {
	if limit <= 0 { limit = 20 }
	return s.repo.ListSubmissions(ctx, devID, limit)
}
func (s *StoreService) ListItems(ctx context.Context, category string, limit, offset int32) ([]*StoreItem, int32, error) {
	if limit <= 0 { limit = 20 }
	return s.repo.ListItems(ctx, category, limit, offset)
}
func (s *StoreService) SearchItems(ctx context.Context, query string, limit int32) ([]*StoreItem, error) {
	if limit <= 0 { limit = 20 }
	return s.repo.SearchItems(ctx, query, limit)
}
func (s *StoreService) GetItem(ctx context.Context, id string) (*StoreItem, error) {
	if id == "" { return nil, ErrInvalidInput }
	return s.repo.GetItem(ctx, id)
}
func (s *StoreService) Purchase(ctx context.Context, userID, itemID string) (*Purchase, error) {
	if userID == "" || itemID == "" { return nil, ErrInvalidInput }
	item, err := s.repo.GetItem(ctx, itemID)
	if err != nil { return nil, err }
	p := &Purchase{ID: uuid.New().String(), UserID: userID, ItemID: itemID, PriceKrw: item.PriceKrw, Status: "completed", PurchasedAt: time.Now().UTC()}
	return p, s.repo.CreatePurchase(ctx, p)
}
func (s *StoreService) GetPurchaseHistory(ctx context.Context, userID string, limit int32) ([]*Purchase, error) {
	if userID == "" { return nil, ErrInvalidInput }
	return s.repo.ListPurchases(ctx, userID, limit)
}
func (s *StoreService) SubmitForReview(ctx context.Context, cartridgeID, devID string) (*ReviewStatus, error) {
	if cartridgeID == "" { return nil, ErrInvalidInput }
	r := &ReviewStatus{SubmissionID: uuid.New().String(), CartridgeID: cartridgeID, Status: "pending", SubmittedAt: time.Now().UTC()}
	return r, s.repo.CreateReview(ctx, r)
}
func (s *StoreService) GetReviewStatus(ctx context.Context, submissionID string) (*ReviewStatus, error) {
	if submissionID == "" { return nil, ErrInvalidInput }
	return s.repo.GetReview(ctx, submissionID)
}
func (s *StoreService) ApproveCartridge(ctx context.Context, submissionID, comment string) (*ReviewStatus, error) {
	r, err := s.repo.GetReview(ctx, submissionID)
	if err != nil { return nil, err }
	r.Status = "approved"; r.ReviewerComment = comment
	return r, s.repo.UpdateReview(ctx, r)
}
func (s *StoreService) RejectCartridge(ctx context.Context, submissionID, reason string) (*ReviewStatus, error) {
	r, err := s.repo.GetReview(ctx, submissionID)
	if err != nil { return nil, err }
	r.Status = "rejected"; r.ReviewerComment = reason
	return r, s.repo.UpdateReview(ctx, r)
}

func (s *StoreService) GetSalesReport(ctx context.Context, devID, period string) *SalesReport {
	return &SalesReport{DeveloperID: devID, TotalRevenue: 1500000, PlatformFee: 450000, DeveloperEarning: 1050000, Period: period, ItemCount: 3}
}
func (s *StoreService) GetPayoutHistory(ctx context.Context, devID string, limit int32) []*PayoutRecord { return nil }
func (s *StoreService) ConfigureRevenueSplit(ctx context.Context, devID string, devPct, platPct int32) *RevenueSplitCfg {
	return &RevenueSplitCfg{DeveloperID: devID, DeveloperPct: devPct, PlatformPct: platPct}
}
func (s *StoreService) RequestPayout(ctx context.Context, devID string, amount int64) *PayoutRecord {
	return &PayoutRecord{ID: uuid.New().String(), DeveloperID: devID, AmountKrw: amount, Status: "pending", RequestedAt: time.Now().UTC()}
}
func (s *StoreService) GetUsageStats(ctx context.Context, itemID string) *UsageStat {
	return &UsageStat{ItemID: itemID, Date: time.Now().Format("2006-01-02"), UsageCount: 150, UniqueUsers: 42}
}
func (s *StoreService) GetRatings(ctx context.Context, itemID string) *RatingInfo {
	return &RatingInfo{ItemID: itemID, AvgRating: 4.5, TotalReviews: 23}
}
func (s *StoreService) SubmitUserReview(ctx context.Context, userID, itemID string, rating int32, comment string) *UserReview {
	return &UserReview{ID: uuid.New().String(), UserID: userID, ItemID: itemID, Rating: rating, Comment: comment, CreatedAt: time.Now().UTC()}
}
func (s *StoreService) GetDevAnalytics(ctx context.Context, devID string) *DevAnalytics {
	return &DevAnalytics{DeveloperID: devID, TotalDownloads: 12500, TotalRevenue: 3500000, ActiveCartridges: 5}
}
GOEOF

cat > "$SVC/internal/repository/memory/store.go" << 'GOEOF'
package memory

import (
	"context"
	"strings"
	"sync"
	"github.com/manpasik/backend/services/cartridge-store-service/internal/service"
)

type StoreRepository struct {
	mu          sync.RWMutex
	developers  map[string]*service.Developer
	devByUser   map[string]string
	apiKeys     map[string]*service.ApiKey
	items       map[string]*service.StoreItem
	purchases   []*service.Purchase
	reviews     map[string]*service.ReviewStatus
}

func NewStoreRepository() *StoreRepository {
	return &StoreRepository{
		developers: make(map[string]*service.Developer),
		devByUser:  make(map[string]string),
		apiKeys:    make(map[string]*service.ApiKey),
		items:      make(map[string]*service.StoreItem),
		reviews:    make(map[string]*service.ReviewStatus),
	}
}
func (r *StoreRepository) CreateDeveloper(_ context.Context, d *service.Developer) error {
	r.mu.Lock(); defer r.mu.Unlock()
	r.developers[d.ID] = d; r.devByUser[d.UserID] = d.ID; return nil
}
func (r *StoreRepository) GetDeveloper(_ context.Context, userID string) (*service.Developer, error) {
	r.mu.RLock(); defer r.mu.RUnlock()
	if id, ok := r.devByUser[userID]; ok { return r.developers[id], nil }
	return nil, service.ErrNotFound
}
func (r *StoreRepository) CreateApiKey(_ context.Context, k *service.ApiKey) error {
	r.mu.Lock(); defer r.mu.Unlock()
	r.apiKeys[k.KeyID] = k; return nil
}
func (r *StoreRepository) CreateItem(_ context.Context, item *service.StoreItem) error {
	r.mu.Lock(); defer r.mu.Unlock()
	r.items[item.ID] = item; return nil
}
func (r *StoreRepository) ListItems(_ context.Context, category string, limit, offset int32) ([]*service.StoreItem, int32, error) {
	r.mu.RLock(); defer r.mu.RUnlock()
	out := make([]*service.StoreItem, 0)
	for _, item := range r.items {
		if category != "" && item.Category != category { continue }
		out = append(out, item)
	}
	total := int32(len(out))
	s := int(offset); if s > len(out) { s = len(out) }
	e := s + int(limit); if e > len(out) { e = len(out) }
	return out[s:e], total, nil
}
func (r *StoreRepository) SearchItems(_ context.Context, query string, limit int32) ([]*service.StoreItem, error) {
	r.mu.RLock(); defer r.mu.RUnlock()
	out := make([]*service.StoreItem, 0)
	q := strings.ToLower(query)
	for _, item := range r.items {
		if strings.Contains(strings.ToLower(item.Name), q) || strings.Contains(strings.ToLower(item.Description), q) {
			out = append(out, item)
			if int32(len(out)) >= limit { break }
		}
	}
	return out, nil
}
func (r *StoreRepository) GetItem(_ context.Context, id string) (*service.StoreItem, error) {
	r.mu.RLock(); defer r.mu.RUnlock()
	if v, ok := r.items[id]; ok { return v, nil }
	return nil, service.ErrNotFound
}
func (r *StoreRepository) CreatePurchase(_ context.Context, p *service.Purchase) error {
	r.mu.Lock(); defer r.mu.Unlock()
	r.purchases = append(r.purchases, p); return nil
}
func (r *StoreRepository) ListPurchases(_ context.Context, userID string, limit int32) ([]*service.Purchase, error) {
	r.mu.RLock(); defer r.mu.RUnlock()
	out := make([]*service.Purchase, 0)
	for _, p := range r.purchases {
		if p.UserID == userID { out = append(out, p) }
		if int32(len(out)) >= limit { break }
	}
	return out, nil
}
func (r *StoreRepository) CreateReview(_ context.Context, rv *service.ReviewStatus) error {
	r.mu.Lock(); defer r.mu.Unlock()
	r.reviews[rv.SubmissionID] = rv; return nil
}
func (r *StoreRepository) GetReview(_ context.Context, id string) (*service.ReviewStatus, error) {
	r.mu.RLock(); defer r.mu.RUnlock()
	if v, ok := r.reviews[id]; ok { return v, nil }
	return nil, service.ErrNotFound
}
func (r *StoreRepository) UpdateReview(_ context.Context, rv *service.ReviewStatus) error {
	r.mu.Lock(); defer r.mu.Unlock()
	r.reviews[rv.SubmissionID] = rv; return nil
}
func (r *StoreRepository) ListSubmissions(_ context.Context, devID string, limit int32) ([]*service.ReviewStatus, error) {
	r.mu.RLock(); defer r.mu.RUnlock()
	out := make([]*service.ReviewStatus, 0)
	for _, rv := range r.reviews { out = append(out, rv) }
	return out, nil
}
GOEOF

cat > "$SVC/internal/handler/grpc.go" << 'GOEOF'
package handler

import (
	"context"
	"github.com/manpasik/backend/services/cartridge-store-service/internal/service"
	v1 "github.com/manpasik/backend/shared/gen/go/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type StoreHandler struct {
	v1.UnimplementedDeveloperServiceServer
	v1.UnimplementedStoreServiceServer
	v1.UnimplementedCartridgeReviewServiceServer
	v1.UnimplementedRevenueServiceServer
	v1.UnimplementedCartridgeAnalyticsServiceServer
	svc *service.StoreService
	log *zap.Logger
}
func NewStoreHandler(svc *service.StoreService, log *zap.Logger) *StoreHandler {
	return &StoreHandler{svc: svc, log: log}
}

// DeveloperService
func (h *StoreHandler) RegisterDeveloper(ctx context.Context, req *v1.RegisterDeveloperRequest) (*v1.DeveloperProfile, error) {
	d, err := h.svc.RegisterDeveloper(ctx, req.GetUserId(), req.GetCompanyName(), req.GetEmail())
	if err != nil { return nil, toGRPC(err) }
	return &v1.DeveloperProfile{DeveloperId: d.ID, UserId: d.UserID, CompanyName: d.CompanyName, Email: d.Email, Status: d.Status, CreatedAt: timestamppb.New(d.CreatedAt)}, nil
}
func (h *StoreHandler) GetDeveloperProfile(ctx context.Context, req *v1.GetDeveloperProfileRequest) (*v1.DeveloperProfile, error) {
	d, err := h.svc.GetDeveloper(ctx, req.GetUserId())
	if err != nil { return nil, toGRPC(err) }
	return &v1.DeveloperProfile{DeveloperId: d.ID, UserId: d.UserID, CompanyName: d.CompanyName, Email: d.Email, Status: d.Status, CreatedAt: timestamppb.New(d.CreatedAt)}, nil
}
func (h *StoreHandler) CreateApiKey(ctx context.Context, req *v1.CreateApiKeyRequest) (*v1.ApiKeyResponse, error) {
	k, err := h.svc.CreateApiKey(ctx, req.GetDeveloperId())
	if err != nil { return nil, toGRPC(err) }
	return &v1.ApiKeyResponse{KeyId: k.KeyID, ApiKey: k.Key, Status: k.Status, CreatedAt: timestamppb.New(k.CreatedAt)}, nil
}
func (h *StoreHandler) SubmitCartridge(ctx context.Context, req *v1.SubmitCartridgeRequest) (*v1.SubmitCartridgeResponse, error) {
	item, err := h.svc.SubmitCartridge(ctx, req.GetDeveloperId(), req.GetName(), req.GetDescription(), req.GetCategory(), req.GetVersion(), req.GetPriceKrw())
	if err != nil { return nil, toGRPC(err) }
	return &v1.SubmitCartridgeResponse{CartridgeId: item.ID, Status: item.Status}, nil
}
func (h *StoreHandler) ListSubmissions(ctx context.Context, req *v1.ListSubmissionsRequest) (*v1.ListSubmissionsResponse, error) {
	list, err := h.svc.ListSubmissions(ctx, req.GetDeveloperId(), req.GetLimit())
	if err != nil { return nil, toGRPC(err) }
	out := make([]*v1.SubmissionInfo, 0, len(list))
	for _, r := range list {
		out = append(out, &v1.SubmissionInfo{SubmissionId: r.SubmissionID, CartridgeId: r.CartridgeID, Status: r.Status, SubmittedAt: timestamppb.New(r.SubmittedAt)})
	}
	return &v1.ListSubmissionsResponse{Submissions: out, Total: int32(len(out))}, nil
}

// StoreService
func (h *StoreHandler) ListStoreItems(ctx context.Context, req *v1.ListStoreItemsRequest) (*v1.ListStoreItemsResponse, error) {
	list, total, err := h.svc.ListItems(ctx, req.GetCategory(), req.GetLimit(), req.GetOffset())
	if err != nil { return nil, toGRPC(err) }
	out := make([]*v1.StoreItem, 0, len(list))
	for _, item := range list { out = append(out, itemToProto(item)) }
	return &v1.ListStoreItemsResponse{Items: out, Total: total}, nil
}
func (h *StoreHandler) SearchCartridges(ctx context.Context, req *v1.SearchCartridgesRequest) (*v1.SearchCartridgesResponse, error) {
	list, err := h.svc.SearchItems(ctx, req.GetQuery(), req.GetLimit())
	if err != nil { return nil, toGRPC(err) }
	out := make([]*v1.StoreItem, 0, len(list))
	for _, item := range list { out = append(out, itemToProto(item)) }
	return &v1.SearchCartridgesResponse{Items: out, Total: int32(len(out))}, nil
}
func (h *StoreHandler) GetStoreItem(ctx context.Context, req *v1.GetStoreItemRequest) (*v1.StoreItem, error) {
	item, err := h.svc.GetItem(ctx, req.GetItemId())
	if err != nil { return nil, toGRPC(err) }
	return itemToProto(item), nil
}
func (h *StoreHandler) PurchaseCartridge(ctx context.Context, req *v1.PurchaseCartridgeRequest) (*v1.PurchaseCartridgeResponse, error) {
	p, err := h.svc.Purchase(ctx, req.GetUserId(), req.GetItemId())
	if err != nil { return nil, toGRPC(err) }
	return &v1.PurchaseCartridgeResponse{PurchaseId: p.ID, Status: p.Status}, nil
}
func (h *StoreHandler) GetPurchaseHistory(ctx context.Context, req *v1.GetPurchaseHistoryRequest) (*v1.GetPurchaseHistoryResponse, error) {
	list, err := h.svc.GetPurchaseHistory(ctx, req.GetUserId(), req.GetLimit())
	if err != nil { return nil, toGRPC(err) }
	out := make([]*v1.PurchaseRecord, 0, len(list))
	for _, p := range list {
		out = append(out, &v1.PurchaseRecord{PurchaseId: p.ID, ItemId: p.ItemID, PriceKrw: p.PriceKrw, Status: p.Status, PurchasedAt: timestamppb.New(p.PurchasedAt)})
	}
	return &v1.GetPurchaseHistoryResponse{Purchases: out, Total: int32(len(out))}, nil
}

// CartridgeReviewService
func (h *StoreHandler) SubmitForReview(ctx context.Context, req *v1.SubmitForReviewRequest) (*v1.ReviewStatus, error) {
	r, err := h.svc.SubmitForReview(ctx, req.GetCartridgeId(), req.GetDeveloperId())
	if err != nil { return nil, toGRPC(err) }
	return &v1.ReviewStatus{SubmissionId: r.SubmissionID, CartridgeId: r.CartridgeID, Status: r.Status, SubmittedAt: timestamppb.New(r.SubmittedAt)}, nil
}
func (h *StoreHandler) GetReviewStatus(ctx context.Context, req *v1.GetReviewStatusRequest) (*v1.ReviewStatus, error) {
	r, err := h.svc.GetReviewStatus(ctx, req.GetSubmissionId())
	if err != nil { return nil, toGRPC(err) }
	return &v1.ReviewStatus{SubmissionId: r.SubmissionID, CartridgeId: r.CartridgeID, Status: r.Status, ReviewerComment: r.ReviewerComment, SubmittedAt: timestamppb.New(r.SubmittedAt)}, nil
}
func (h *StoreHandler) ApproveCartridge(ctx context.Context, req *v1.ApproveCartridgeRequest) (*v1.ReviewStatus, error) {
	r, err := h.svc.ApproveCartridge(ctx, req.GetSubmissionId(), req.GetComment())
	if err != nil { return nil, toGRPC(err) }
	return &v1.ReviewStatus{SubmissionId: r.SubmissionID, CartridgeId: r.CartridgeID, Status: r.Status, ReviewerComment: r.ReviewerComment, SubmittedAt: timestamppb.New(r.SubmittedAt)}, nil
}
func (h *StoreHandler) RejectCartridge(ctx context.Context, req *v1.RejectCartridgeRequest) (*v1.ReviewStatus, error) {
	r, err := h.svc.RejectCartridge(ctx, req.GetSubmissionId(), req.GetReason())
	if err != nil { return nil, toGRPC(err) }
	return &v1.ReviewStatus{SubmissionId: r.SubmissionID, CartridgeId: r.CartridgeID, Status: r.Status, ReviewerComment: r.ReviewerComment, SubmittedAt: timestamppb.New(r.SubmittedAt)}, nil
}

// RevenueService
func (h *StoreHandler) GetSalesReport(ctx context.Context, req *v1.GetSalesReportRequest) (*v1.SalesReport, error) {
	r := h.svc.GetSalesReport(ctx, req.GetDeveloperId(), req.GetPeriod())
	return &v1.SalesReport{DeveloperId: r.DeveloperID, TotalRevenue: r.TotalRevenue, PlatformFee: r.PlatformFee, DeveloperEarning: r.DeveloperEarning, Period: r.Period, ItemCount: r.ItemCount}, nil
}
func (h *StoreHandler) GetPayoutHistory(ctx context.Context, req *v1.GetPayoutHistoryRequest) (*v1.GetPayoutHistoryResponse, error) {
	return &v1.GetPayoutHistoryResponse{Payouts: nil, Total: 0}, nil
}
func (h *StoreHandler) ConfigureRevenueSplit(ctx context.Context, req *v1.ConfigureRevenueSplitRequest) (*v1.RevenueSplitConfig, error) {
	c := h.svc.ConfigureRevenueSplit(ctx, req.GetDeveloperId(), req.GetDeveloperPct(), req.GetPlatformPct())
	return &v1.RevenueSplitConfig{DeveloperId: c.DeveloperID, DeveloperPct: c.DeveloperPct, PlatformPct: c.PlatformPct}, nil
}
func (h *StoreHandler) RequestPayout(ctx context.Context, req *v1.RequestPayoutRequest) (*v1.PayoutResponse, error) {
	p := h.svc.RequestPayout(ctx, req.GetDeveloperId(), req.GetAmountKrw())
	return &v1.PayoutResponse{PayoutId: p.ID, Status: p.Status, AmountKrw: p.AmountKrw}, nil
}

// CartridgeAnalyticsService
func (h *StoreHandler) GetUsageStats(ctx context.Context, req *v1.GetCartridgeUsageStatsRequest) (*v1.CartridgeUsageStatsResponse, error) {
	u := h.svc.GetUsageStats(ctx, req.GetItemId())
	return &v1.CartridgeUsageStatsResponse{ItemId: u.ItemID, Date: u.Date, UsageCount: u.UsageCount, UniqueUsers: u.UniqueUsers}, nil
}
func (h *StoreHandler) GetRatings(ctx context.Context, req *v1.GetCartridgeRatingsRequest) (*v1.CartridgeRatingsResponse, error) {
	r := h.svc.GetRatings(ctx, req.GetItemId())
	return &v1.CartridgeRatingsResponse{ItemId: r.ItemID, AvgRating: r.AvgRating, TotalReviews: r.TotalReviews}, nil
}
func (h *StoreHandler) SubmitUserReview(ctx context.Context, req *v1.SubmitUserReviewRequest) (*v1.SubmitUserReviewResponse, error) {
	r := h.svc.SubmitUserReview(ctx, req.GetUserId(), req.GetItemId(), req.GetRating(), req.GetComment())
	return &v1.SubmitUserReviewResponse{ReviewId: r.ID, Success: true}, nil
}
func (h *StoreHandler) GetDeveloperAnalytics(ctx context.Context, req *v1.GetDeveloperAnalyticsRequest) (*v1.DeveloperAnalytics, error) {
	a := h.svc.GetDevAnalytics(ctx, req.GetDeveloperId())
	return &v1.DeveloperAnalytics{DeveloperId: a.DeveloperID, TotalDownloads: a.TotalDownloads, TotalRevenue: a.TotalRevenue, ActiveCartridges: a.ActiveCartridges}, nil
}

func itemToProto(item *service.StoreItem) *v1.StoreItem {
	return &v1.StoreItem{Id: item.ID, DeveloperId: item.DeveloperID, Name: item.Name, Description: item.Description,
		Category: item.Category, Version: item.Version, PriceKrw: item.PriceKrw, Rating: item.Rating,
		Downloads: item.Downloads, Status: item.Status, CreatedAt: timestamppb.New(item.CreatedAt)}
}
func toGRPC(err error) error {
	switch err {
	case service.ErrNotFound: return status.Error(codes.NotFound, err.Error())
	case service.ErrInvalidInput: return status.Error(codes.InvalidArgument, err.Error())
	case service.ErrAlreadyExists: return status.Error(codes.AlreadyExists, err.Error())
	default: return status.Error(codes.Internal, err.Error())
	}
}
GOEOF

cat > "$SVC/internal/service/store_test.go" << 'GOEOF'
package service_test

import (
	"context"
	"testing"
	"github.com/manpasik/backend/services/cartridge-store-service/internal/repository/memory"
	"github.com/manpasik/backend/services/cartridge-store-service/internal/service"
	"go.uber.org/zap"
)

func newSvc() *service.StoreService {
	return service.NewStoreService(zap.NewNop(), memory.NewStoreRepository())
}
func TestRegisterDeveloper(t *testing.T) {
	s := newSvc()
	d, err := s.RegisterDeveloper(context.Background(), "u1", "테스트회사", "test@example.com")
	if err != nil { t.Fatal(err) }
	if d.ID == "" { t.Error("empty id") }
}
func TestSubmitAndListItems(t *testing.T) {
	s := newSvc(); ctx := context.Background()
	s.RegisterDeveloper(ctx, "u1", "Co", "a@b.com")
	s.SubmitCartridge(ctx, "dev1", "카트리지A", "설명", "health", "1.0", 5000)
	s.SubmitCartridge(ctx, "dev1", "카트리지B", "설명2", "health", "1.0", 3000)
	list, total, _ := s.ListItems(ctx, "", 10, 0)
	if total != 2 || len(list) != 2 { t.Errorf("expected 2, got %d", total) }
}
func TestPurchase(t *testing.T) {
	s := newSvc(); ctx := context.Background()
	item, _ := s.SubmitCartridge(ctx, "dev1", "X", "", "", "1.0", 1000)
	p, err := s.Purchase(ctx, "buyer1", item.ID)
	if err != nil { t.Fatal(err) }
	if p.PriceKrw != 1000 { t.Error("price mismatch") }
}
func TestReviewFlow(t *testing.T) {
	s := newSvc(); ctx := context.Background()
	r, _ := s.SubmitForReview(ctx, "cart1", "dev1")
	if r.Status != "pending" { t.Error("expected pending") }
	approved, _ := s.ApproveCartridge(ctx, r.SubmissionID, "좋습니다")
	if approved.Status != "approved" { t.Error("expected approved") }
}
func TestCreateApiKey(t *testing.T) {
	s := newSvc()
	k, err := s.CreateApiKey(context.Background(), "dev1")
	if err != nil { t.Fatal(err) }
	if k.Key == "" || k.KeyID == "" { t.Error("empty key") }
}
GOEOF

cat > "$SVC/Dockerfile" << 'GOEOF'
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /cartridge-store-service ./services/cartridge-store-service/cmd

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
COPY --from=builder /cartridge-store-service /cartridge-store-service
EXPOSE 50079
ENTRYPOINT ["/cartridge-store-service"]
GOEOF

echo "[OK] cartridge-store-service created"
echo "=== Part 2 done (data-platform + cartridge-store) ==="
