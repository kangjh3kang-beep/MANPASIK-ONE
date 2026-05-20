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
	if err != nil {
		return nil, toGRPC(err)
	}
	return &v1.RegionStats{RegionCode: s.RegionCode, RegionName: s.RegionName, Period: s.Period,
		AvgValue: s.AvgValue, MinValue: s.MinValue, MaxValue: s.MaxValue, SampleCount: s.SampleCount, RiskLevel: s.RiskLevel}, nil
}
func (h *PlatformHandler) ListHotspots(ctx context.Context, req *v1.ListHotspotsRequest) (*v1.ListHotspotsResponse, error) {
	list, err := h.svc.ListHotspots(ctx, req.GetBiomarker(), req.GetRiskLevel(), req.GetLimit())
	if err != nil {
		return nil, toGRPC(err)
	}
	out := make([]*v1.RegionStats, 0, len(list))
	for _, s := range list {
		out = append(out, &v1.RegionStats{RegionCode: s.RegionCode, RegionName: s.RegionName, Period: s.Period,
			AvgValue: s.AvgValue, MinValue: s.MinValue, MaxValue: s.MaxValue, SampleCount: s.SampleCount, RiskLevel: s.RiskLevel})
	}
	return &v1.ListHotspotsResponse{Hotspots: out}, nil
}
func (h *PlatformHandler) GetTrendByRegion(ctx context.Context, req *v1.GetTrendByRegionRequest) (*v1.RegionTrend, error) {
	t, err := h.svc.GetTrend(ctx, req.GetRegionCode(), req.GetBiomarker(), req.GetStartDate(), req.GetEndDate())
	if err != nil {
		return nil, toGRPC(err)
	}
	pts := make([]*v1.RegionTrendPoint, 0, len(t.Points))
	for _, p := range t.Points {
		pts = append(pts, &v1.RegionTrendPoint{Date: p.Date, Value: p.Value, SampleCount: p.SampleCount})
	}
	return &v1.RegionTrend{RegionCode: t.RegionCode, Biomarker: t.Biomarker, Points: pts, TrendDirection: t.TrendDirection}, nil
}

// DataProvisionService
func (h *PlatformHandler) ListDatasets(ctx context.Context, req *v1.ListDatasetsRequest) (*v1.ListDatasetsResponse, error) {
	list, total, err := h.svc.ListDatasets(ctx, req.GetDatasetType(), req.GetLimit(), req.GetOffset())
	if err != nil {
		return nil, toGRPC(err)
	}
	out := make([]*v1.AnonymizedDataset, 0, len(list))
	for _, ds := range list {
		out = append(out, &v1.AnonymizedDataset{Id: ds.ID, DatasetType: ds.DatasetType, Description: ds.Description,
			RecordCount: ds.RecordCount, DateRange: ds.DateRange, CreatedAt: timestamppb.New(ds.CreatedAt)})
	}
	return &v1.ListDatasetsResponse{Datasets: out, Total: total}, nil
}
func (h *PlatformHandler) GetAnonymizedDataset(ctx context.Context, req *v1.GetAnonymizedDatasetRequest) (*v1.AnonymizedDataset, error) {
	ds, err := h.svc.GetDataset(ctx, req.GetDatasetId())
	if err != nil {
		return nil, toGRPC(err)
	}
	return &v1.AnonymizedDataset{Id: ds.ID, DatasetType: ds.DatasetType, Description: ds.Description,
		RecordCount: ds.RecordCount, DateRange: ds.DateRange, CreatedAt: timestamppb.New(ds.CreatedAt)}, nil
}
func (h *PlatformHandler) RequestDataAccess(ctx context.Context, req *v1.RequestDataAccessRequest) (*v1.DataAccessResponse, error) {
	r, err := h.svc.RequestAccess(ctx, req.GetRequesterId(), req.GetDatasetId(), req.GetPurpose())
	if err != nil {
		return nil, toGRPC(err)
	}
	return &v1.DataAccessResponse{RequestId: r.RequestID, Status: r.Status}, nil
}
func (h *PlatformHandler) GetDataSharingConsent(ctx context.Context, req *v1.GetConsentRequest) (*v1.DataSharingConsentInfo, error) {
	c, err := h.svc.GetConsent(ctx, req.GetUserId())
	if err != nil {
		return nil, toGRPC(err)
	}
	return &v1.DataSharingConsentInfo{UserId: c.UserID, AnonymousResearch: c.AnonymousResearch,
		RegionStats: c.RegionStatsConsent, EnterpriseData: c.EnterpriseData, UpdatedAt: timestamppb.New(c.UpdatedAt)}, nil
}
func (h *PlatformHandler) UpdateDataSharingConsent(ctx context.Context, req *v1.UpdateConsentRequest) (*v1.DataSharingConsentInfo, error) {
	c, err := h.svc.UpdateConsent(ctx, req.GetUserId(), req.GetAnonymousResearch(), req.GetRegionStats(), req.GetEnterpriseData())
	if err != nil {
		return nil, toGRPC(err)
	}
	return &v1.DataSharingConsentInfo{UserId: c.UserID, AnonymousResearch: c.AnonymousResearch,
		RegionStats: c.RegionStatsConsent, EnterpriseData: c.EnterpriseData, UpdatedAt: timestamppb.New(c.UpdatedAt)}, nil
}
func (h *PlatformHandler) TriggerAnonymousAggregation(ctx context.Context, req *v1.TriggerAggregationRequest) (*v1.TriggerAggregationResponse, error) {
	aggID, st, err := h.svc.TriggerAggregation(ctx, req.GetBiomarker(), req.GetDate())
	if err != nil {
		return nil, toGRPC(err)
	}
	return &v1.TriggerAggregationResponse{AggregationId: aggID, Status: st}, nil
}
func (h *PlatformHandler) GetAggregatedInsights(ctx context.Context, req *v1.GetInsightsRequest) (*v1.GetInsightsResponse, error) {
	list, err := h.svc.GetInsights(ctx, req.GetBiomarker(), req.GetLimit())
	if err != nil {
		return nil, toGRPC(err)
	}
	out := make([]*v1.AggregatedInsight, 0, len(list))
	for _, i := range list {
		out = append(out, &v1.AggregatedInsight{Id: i.ID, Biomarker: i.Biomarker, InsightType: i.InsightType,
			Summary: i.Summary, CreatedAt: timestamppb.New(i.CreatedAt)})
	}
	return &v1.GetInsightsResponse{Insights: out}, nil
}

func toGRPC(err error) error {
	switch err {
	case service.ErrNotFound:
		return status.Error(codes.NotFound, err.Error())
	case service.ErrInvalidInput:
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
