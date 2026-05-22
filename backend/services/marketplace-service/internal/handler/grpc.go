// Package handler는 marketplace-service의 gRPC 핸들러를 구현합니다.
package handler

import (
	"context"

	"github.com/manpasik/backend/services/marketplace-service/internal/service"
	v1 "github.com/manpasik/backend/shared/gen/go/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// MarketplaceHandler는 gRPC MarketplaceService 핸들러입니다.
type MarketplaceHandler struct {
	v1.UnimplementedMarketplaceServiceServer
	svc *service.MarketplaceService
	log *zap.Logger
}

// NewMarketplaceHandler는 MarketplaceHandler를 생성합니다.
func NewMarketplaceHandler(svc *service.MarketplaceService, log *zap.Logger) *MarketplaceHandler {
	return &MarketplaceHandler{svc: svc, log: log}
}

// ListProducts는 상품 목록을 조회하는 RPC입니다.
func (h *MarketplaceHandler) ListProducts(ctx context.Context, req *v1.ListMarketplaceProductsRequest) (*v1.ListMarketplaceProductsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "요청이 비어있습니다")
	}

	products, err := h.svc.ListPartnerProducts(ctx, req.PartnerId, req.Category)
	if err != nil {
		return nil, toGRPC(err)
	}

	pbProducts := make([]*v1.MarketplaceProduct, 0, len(products))
	for _, p := range products {
		pbProducts = append(pbProducts, productToProto(p))
	}

	return &v1.ListMarketplaceProductsResponse{
		Products: pbProducts,
		Total:    int32(len(pbProducts)),
	}, nil
}

// GetProduct는 상품을 조회하는 RPC입니다.
func (h *MarketplaceHandler) GetProduct(ctx context.Context, req *v1.GetMarketplaceProductRequest) (*v1.MarketplaceProduct, error) {
	if req == nil || req.ProductId == "" {
		return nil, status.Error(codes.InvalidArgument, "product_id는 필수입니다")
	}

	product, err := h.svc.GetProduct(ctx, req.ProductId)
	if err != nil {
		return nil, toGRPC(err)
	}

	return productToProto(product), nil
}

// CreateProduct는 상품을 생성하는 RPC입니다.
func (h *MarketplaceHandler) CreateProduct(ctx context.Context, req *v1.CreateMarketplaceProductRequest) (*v1.MarketplaceProduct, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "요청이 비어있습니다")
	}
	if req.PartnerId == "" {
		return nil, status.Error(codes.InvalidArgument, "partner_id는 필수입니다")
	}
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name은 필수입니다")
	}

	product := &service.PartnerProduct{
		PartnerID:   req.PartnerId,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Category:    req.Category,
		ImageURL:    req.ImageUrl,
	}

	created, err := h.svc.CreateProduct(ctx, product)
	if err != nil {
		return nil, toGRPC(err)
	}

	return productToProto(created), nil
}

// UpdateProduct는 상품을 업데이트하는 RPC입니다.
func (h *MarketplaceHandler) UpdateProduct(ctx context.Context, req *v1.UpdateMarketplaceProductRequest) (*v1.UpdateMarketplaceProductResponse, error) {
	if req == nil || req.ProductId == "" {
		return nil, status.Error(codes.InvalidArgument, "product_id는 필수입니다")
	}

	product := &service.PartnerProduct{
		ID:          req.ProductId,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Category:    req.Category,
		ImageURL:    req.ImageUrl,
		IsActive:    req.IsActive,
	}

	ok, err := h.svc.UpdateProduct(ctx, product)
	if err != nil {
		return nil, toGRPC(err)
	}

	return &v1.UpdateMarketplaceProductResponse{Success: ok}, nil
}

// RegisterPartner는 파트너를 등록하는 RPC입니다.
func (h *MarketplaceHandler) RegisterPartner(ctx context.Context, req *v1.RegisterMarketplacePartnerRequest) (*v1.RegisterMarketplacePartnerResponse, error) {
	if req == nil || req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name은 필수입니다")
	}
	if req.ContactEmail == "" {
		return nil, status.Error(codes.InvalidArgument, "contact_email은 필수입니다")
	}

	partner := &service.Partner{
		Name:         req.Name,
		Description:  req.Description,
		ContactEmail: req.ContactEmail,
	}

	partnerID, err := h.svc.RegisterPartner(ctx, partner)
	if err != nil {
		return nil, toGRPC(err)
	}

	return &v1.RegisterMarketplacePartnerResponse{
		PartnerId: partnerID,
		Status:    "pending",
	}, nil
}

// GetPartnerStats는 파트너 통계를 조회하는 RPC입니다.
func (h *MarketplaceHandler) GetPartnerStats(ctx context.Context, req *v1.GetPartnerStatsRequest) (*v1.PartnerStatsResponse, error) {
	if req == nil || req.PartnerId == "" {
		return nil, status.Error(codes.InvalidArgument, "partner_id는 필수입니다")
	}

	stats, err := h.svc.GetPartnerStats(ctx, req.PartnerId)
	if err != nil {
		return nil, toGRPC(err)
	}

	return &v1.PartnerStatsResponse{
		PartnerId:     stats.PartnerID,
		TotalProducts: int32(stats.TotalProducts),
		TotalOrders:   int32(stats.TotalOrders),
		Revenue:       stats.Revenue,
	}, nil
}

// ============================================================================
// Helpers
// ============================================================================

func productToProto(p *service.PartnerProduct) *v1.MarketplaceProduct {
	result := &v1.MarketplaceProduct{
		Id:          p.ID,
		PartnerId:   p.PartnerID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Category:    p.Category,
		ImageUrl:    p.ImageURL,
		IsActive:    p.IsActive,
	}
	if !p.CreatedAt.IsZero() {
		result.CreatedAt = timestamppb.New(p.CreatedAt)
	}
	if !p.UpdatedAt.IsZero() {
		result.UpdatedAt = timestamppb.New(p.UpdatedAt)
	}
	return result
}

func toGRPC(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case isError(err, service.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case isError(err, service.ErrInvalidInput), isError(err, service.ErrInvalidCategory):
		return status.Error(codes.InvalidArgument, err.Error())
	case isError(err, service.ErrPartnerNotActive), isError(err, service.ErrInvalidTransition):
		return status.Error(codes.FailedPrecondition, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}

func isError(err, target error) bool {
	for e := err; e != nil; {
		if e == target {
			return true
		}
		// errors.Unwrap에 해당하는 간단한 체크
		if u, ok := e.(interface{ Unwrap() error }); ok {
			e = u.Unwrap()
		} else {
			break
		}
	}
	return false
}
