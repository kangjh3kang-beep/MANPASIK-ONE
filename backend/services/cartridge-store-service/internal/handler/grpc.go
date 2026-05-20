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

// === DeveloperService ===

func (h *StoreHandler) RegisterDeveloper(ctx context.Context, req *v1.RegisterDeveloperRequest) (*v1.DeveloperProfile, error) {
	d, err := h.svc.RegisterDeveloper(ctx, req.GetUserId(), req.GetCompanyName(), req.GetTier())
	if err != nil {
		return nil, toGRPC(err)
	}
	return devToProto(d), nil
}

func (h *StoreHandler) GetDeveloperProfile(ctx context.Context, req *v1.GetDeveloperProfileRequest) (*v1.DeveloperProfile, error) {
	d, err := h.svc.GetDeveloper(ctx, req.GetUserId())
	if err != nil {
		return nil, toGRPC(err)
	}
	return devToProto(d), nil
}

func (h *StoreHandler) CreateApiKey(ctx context.Context, req *v1.CreateApiKeyRequest) (*v1.ApiKeyResponse, error) {
	k, err := h.svc.CreateApiKey(ctx, req.GetDeveloperId())
	if err != nil {
		return nil, toGRPC(err)
	}
	return &v1.ApiKeyResponse{
		KeyId:     k.KeyID,
		ApiKey:    k.Key,
		Name:      "default",
		CreatedAt: timestamppb.New(k.CreatedAt),
	}, nil
}

func (h *StoreHandler) SubmitCartridge(ctx context.Context, req *v1.SubmitCartridgeRequest) (*v1.SubmitCartridgeResponse, error) {
	item, err := h.svc.SubmitCartridge(ctx, req.GetDeveloperId(), req.GetCartridgeName(), req.GetDescription(), req.GetCategory(), "", int64(req.GetPriceKrw()))
	if err != nil {
		return nil, toGRPC(err)
	}
	return &v1.SubmitCartridgeResponse{SubmissionId: item.ID, Status: item.Status}, nil
}

func (h *StoreHandler) ListSubmissions(ctx context.Context, req *v1.ListSubmissionsRequest) (*v1.ListSubmissionsResponse, error) {
	list, err := h.svc.ListSubmissions(ctx, req.GetDeveloperId(), req.GetLimit())
	if err != nil {
		return nil, toGRPC(err)
	}
	out := make([]*v1.CartridgeSubmission, 0, len(list))
	for _, r := range list {
		out = append(out, &v1.CartridgeSubmission{
			Id:            r.SubmissionID,
			CartridgeName: r.CartridgeID,
			Status:        r.Status,
			ReviewComment: r.ReviewerComment,
			SubmittedAt:   timestamppb.New(r.SubmittedAt),
		})
	}
	return &v1.ListSubmissionsResponse{Submissions: out, Total: int32(len(out))}, nil
}

// === StoreService ===

func (h *StoreHandler) ListStoreItems(ctx context.Context, req *v1.ListStoreItemsRequest) (*v1.ListStoreItemsResponse, error) {
	list, total, err := h.svc.ListItems(ctx, req.GetCategory(), req.GetLimit(), req.GetOffset())
	if err != nil {
		return nil, toGRPC(err)
	}
	out := make([]*v1.StoreItem, 0, len(list))
	for _, item := range list {
		out = append(out, itemToProto(item))
	}
	return &v1.ListStoreItemsResponse{Items: out, Total: total}, nil
}

func (h *StoreHandler) SearchCartridges(ctx context.Context, req *v1.SearchCartridgesRequest) (*v1.SearchCartridgesResponse, error) {
	list, err := h.svc.SearchItems(ctx, req.GetQuery(), req.GetLimit())
	if err != nil {
		return nil, toGRPC(err)
	}
	out := make([]*v1.StoreItem, 0, len(list))
	for _, item := range list {
		out = append(out, itemToProto(item))
	}
	return &v1.SearchCartridgesResponse{Results: out, Total: int32(len(out))}, nil
}

func (h *StoreHandler) GetStoreItem(ctx context.Context, req *v1.GetStoreItemRequest) (*v1.StoreItem, error) {
	item, err := h.svc.GetItem(ctx, req.GetItemId())
	if err != nil {
		return nil, toGRPC(err)
	}
	return itemToProto(item), nil
}

func (h *StoreHandler) PurchaseCartridge(ctx context.Context, req *v1.PurchaseCartridgeRequest) (*v1.PurchaseCartridgeResponse, error) {
	p, err := h.svc.Purchase(ctx, req.GetUserId(), req.GetItemId())
	if err != nil {
		return nil, toGRPC(err)
	}
	return &v1.PurchaseCartridgeResponse{PurchaseId: p.ID, Success: true, Message: "purchased"}, nil
}

func (h *StoreHandler) GetPurchaseHistory(ctx context.Context, req *v1.GetPurchaseHistoryRequest) (*v1.GetPurchaseHistoryResponse, error) {
	list, err := h.svc.GetPurchaseHistory(ctx, req.GetUserId(), req.GetLimit())
	if err != nil {
		return nil, toGRPC(err)
	}
	out := make([]*v1.CartridgePurchase, 0, len(list))
	for _, p := range list {
		out = append(out, &v1.CartridgePurchase{
			Id:            p.ID,
			ItemId:        p.ItemID,
			CartridgeName: "",
			PriceKrw:      int32(p.PriceKrw),
			PurchasedAt:   timestamppb.New(p.PurchasedAt),
		})
	}
	return &v1.GetPurchaseHistoryResponse{Purchases: out, Total: int32(len(out))}, nil
}

// === CartridgeReviewService ===

func (h *StoreHandler) SubmitForReview(ctx context.Context, req *v1.SubmitForReviewRequest) (*v1.ReviewStatus, error) {
	r, err := h.svc.SubmitForReview(ctx, req.GetSubmissionId(), "")
	if err != nil {
		return nil, toGRPC(err)
	}
	return reviewToProto(r), nil
}

func (h *StoreHandler) GetReviewStatus(ctx context.Context, req *v1.GetReviewStatusRequest) (*v1.ReviewStatus, error) {
	r, err := h.svc.GetReviewStatus(ctx, req.GetSubmissionId())
	if err != nil {
		return nil, toGRPC(err)
	}
	return reviewToProto(r), nil
}

func (h *StoreHandler) ApproveCartridge(ctx context.Context, req *v1.ApproveCartridgeRequest) (*v1.ReviewStatus, error) {
	r, err := h.svc.ApproveCartridge(ctx, req.GetSubmissionId(), req.GetComment())
	if err != nil {
		return nil, toGRPC(err)
	}
	return reviewToProto(r), nil
}

func (h *StoreHandler) RejectCartridge(ctx context.Context, req *v1.RejectCartridgeRequest) (*v1.ReviewStatus, error) {
	r, err := h.svc.RejectCartridge(ctx, req.GetSubmissionId(), req.GetComment())
	if err != nil {
		return nil, toGRPC(err)
	}
	return reviewToProto(r), nil
}

// === RevenueService ===

func (h *StoreHandler) GetSalesReport(ctx context.Context, req *v1.GetSalesReportRequest) (*v1.SalesReport, error) {
	r := h.svc.GetSalesReport(ctx, req.GetDeveloperId(), req.GetPeriod())
	return &v1.SalesReport{
		DeveloperId:         r.DeveloperID,
		Period:              r.Period,
		TotalRevenueKrw:     r.TotalRevenue,
		PlatformFeeKrw:      r.PlatformFee,
		DeveloperEarningKrw: r.DeveloperEarning,
		TotalSales:          r.ItemCount,
	}, nil
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
	return &v1.PayoutResponse{PayoutId: p.ID, Status: p.Status}, nil
}

// === CartridgeAnalyticsService ===

func (h *StoreHandler) GetUsageStats(ctx context.Context, req *v1.GetCartridgeUsageStatsRequest) (*v1.CartridgeUsageStatsResponse, error) {
	u := h.svc.GetUsageStats(ctx, req.GetListingId())
	return &v1.CartridgeUsageStatsResponse{
		ListingId:   u.ItemID,
		TotalUses:   u.UsageCount,
		UniqueUsers: u.UniqueUsers,
	}, nil
}

func (h *StoreHandler) GetRatings(ctx context.Context, req *v1.GetCartridgeRatingsRequest) (*v1.CartridgeRatingsResponse, error) {
	r := h.svc.GetRatings(ctx, req.GetListingId())
	return &v1.CartridgeRatingsResponse{AvgRating: r.AvgRating, TotalReviews: r.TotalReviews}, nil
}

func (h *StoreHandler) SubmitUserReview(ctx context.Context, req *v1.SubmitUserReviewRequest) (*v1.SubmitUserReviewResponse, error) {
	r := h.svc.SubmitUserReview(ctx, req.GetUserId(), req.GetListingId(), req.GetRating(), req.GetComment())
	return &v1.SubmitUserReviewResponse{ReviewId: r.ID, Success: true}, nil
}

func (h *StoreHandler) GetDeveloperAnalytics(ctx context.Context, req *v1.GetDeveloperAnalyticsRequest) (*v1.DeveloperAnalytics, error) {
	a := h.svc.GetDevAnalytics(ctx, req.GetDeveloperId())
	return &v1.DeveloperAnalytics{
		DeveloperId:     a.DeveloperID,
		TotalCartridges: a.ActiveCartridges,
		TotalDownloads:  int32(a.TotalDownloads),
		TotalRevenueKrw: a.TotalRevenue,
	}, nil
}

// === DeveloperService Phase 5: Cartridge Type Registration ===

func (h *StoreHandler) RegisterCartridgeType(ctx context.Context, req *v1.RegisterCartridgeTypeRequest) (*v1.RegisterCartridgeTypeResponse, error) {
	if req == nil || req.DeveloperId == "" || req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "developer_id와 name은 필수입니다")
	}

	def := &service.CartridgeDefinition{
		DeveloperID:          req.DeveloperId,
		CategoryCode:         int(req.CategoryCode),
		TypeIndex:            int(req.TypeIndex),
		Name:                 req.Name,
		Version:              req.Version,
		InterfaceVersion:     req.CsiVersion, // proto 필드명 유지 (외부 API 계약), 내부에서는 추상화된 InterfaceVersion으로 매핑
		PinConfigJSON:        req.PinConfigJson,
		AFEBlocks:            req.AfeBlocks,
		MeasurementModes:     req.MeasurementModes,
		CalibrationSchemaJSON: req.CalibrationSchemaJson,
		AlphaDefault:         req.AlphaDefault,
		TempCompensation:     req.TempCompensation,
		HumidityCompensation: req.HumidityCompensation,
		RegulatoryClass:      req.RegulatoryClass,
		RegulatoryPath:       req.RegulatoryPath,
	}

	result, err := h.svc.RegisterCartridgeType(ctx, def)
	if err != nil {
		return nil, toGRPC(err)
	}

	// 인터페이스 호환성 판정은 Phase 5 SSOT(하드웨어 영역)의 책임입니다.
	// 본 SW 게이트웨이는 메타데이터를 기록할 뿐이며, 실제 호환성 판정은
	// 디바이스/펌웨어 측 또는 별도 인증 시스템에서 수행합니다.
	return &v1.RegisterCartridgeTypeResponse{
		DefinitionId:       result.ID,
		Status:             result.ApprovalStatus,
		CsiCompatibility:   "metadata_recorded", // Phase 5에서 실제 판정
		ValidationWarnings: nil,
		CreatedAt:          timestamppb.New(result.CreatedAt),
	}, nil
}

func (h *StoreHandler) GetCartridgeDefinition(ctx context.Context, req *v1.GetCartridgeDefinitionRequest) (*v1.CartridgeDefinitionDetail, error) {
	if req == nil || req.DefinitionId == "" {
		return nil, status.Error(codes.InvalidArgument, "definition_id는 필수입니다")
	}

	def, err := h.svc.GetCartridgeDefinition(ctx, req.DefinitionId)
	if err != nil {
		return nil, toGRPC(err)
	}

	return &v1.CartridgeDefinitionDetail{
		Id:                    def.ID,
		DeveloperId:           def.DeveloperID,
		CategoryCode:          int32(def.CategoryCode),
		TypeIndex:             int32(def.TypeIndex),
		Name:                  def.Name,
		Version:               def.Version,
		CsiVersion:            def.InterfaceVersion, // proto 필드명은 외부 호환성 위해 유지
		PinConfigJson:         def.PinConfigJSON,
		AfeBlocks:             def.AFEBlocks,
		MeasurementModes:      def.MeasurementModes,
		CalibrationSchemaJson: def.CalibrationSchemaJSON,
		AlphaDefault:          def.AlphaDefault,
		TempCompensation:      def.TempCompensation,
		HumidityCompensation:  def.HumidityCompensation,
		RegulatoryClass:       def.RegulatoryClass,
		RegulatoryPath:        def.RegulatoryPath,
		ApprovalStatus:        def.ApprovalStatus,
		CreatedAt:             timestamppb.New(def.CreatedAt),
	}, nil
}

// === helpers ===

func devToProto(d *service.Developer) *v1.DeveloperProfile {
	return &v1.DeveloperProfile{
		Id: d.ID, UserId: d.UserID, CompanyName: d.CompanyName,
		Tier: d.Tier, Status: d.Status, CartridgeCount: 0,
		CreatedAt: timestamppb.New(d.CreatedAt),
	}
}

func itemToProto(item *service.StoreItem) *v1.StoreItem {
	return &v1.StoreItem{
		Id:            item.ID,
		CartridgeName: item.Name,
		DeveloperName: item.DeveloperID,
		Category:      item.Category,
		Description:   item.Description,
		PriceKrw:      int32(item.PriceKrw),
		Rating:        item.Rating,
		ReviewCount:   0,
		DownloadCount: item.Downloads,
		PublishedAt:   timestamppb.New(item.CreatedAt),
	}
}

func reviewToProto(r *service.ReviewStatus) *v1.ReviewStatus {
	return &v1.ReviewStatus{
		SubmissionId:    r.SubmissionID,
		Status:          r.Status,
		ReviewerComment: r.ReviewerComment,
		UpdatedAt:       timestamppb.New(r.SubmittedAt),
	}
}

func toGRPC(err error) error {
	switch err {
	case service.ErrNotFound:
		return status.Error(codes.NotFound, err.Error())
	case service.ErrInvalidInput:
		return status.Error(codes.InvalidArgument, err.Error())
	case service.ErrAlreadyExists:
		return status.Error(codes.AlreadyExists, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
