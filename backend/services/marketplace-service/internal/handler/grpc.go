// Package handler는 marketplace-service의 gRPC/HTTP 핸들러를 구현합니다.
// proto/gRPC 연동 전 순수 Go 구조체 기반 스켈레톤입니다.
package handler

import (
	"context"
	"log"

	"github.com/manpasik/backend/services/marketplace-service/internal/service"
)

// MarketplaceHandler는 마켓플레이스 요청 핸들러입니다.
type MarketplaceHandler struct {
	svc *service.MarketplaceService
}

// NewMarketplaceHandler는 새 MarketplaceHandler를 생성합니다.
func NewMarketplaceHandler(svc *service.MarketplaceService) *MarketplaceHandler {
	return &MarketplaceHandler{svc: svc}
}

// ListPartnerProducts는 파트너 상품 목록을 반환합니다.
func (h *MarketplaceHandler) ListPartnerProducts(ctx context.Context, partnerID, category string) ([]*service.PartnerProduct, error) {
	products, err := h.svc.ListPartnerProducts(ctx, partnerID, category)
	if err != nil {
		log.Printf("[marketplace-handler] ListPartnerProducts error: %v", err)
		return nil, err
	}
	return products, nil
}

// RegisterPartner는 파트너를 등록하고 ID를 반환합니다.
func (h *MarketplaceHandler) RegisterPartner(ctx context.Context, partner *service.Partner) (string, error) {
	id, err := h.svc.RegisterPartner(ctx, partner)
	if err != nil {
		log.Printf("[marketplace-handler] RegisterPartner error: %v", err)
		return "", err
	}
	return id, nil
}

// GetPartnerStats는 파트너 통계를 반환합니다.
func (h *MarketplaceHandler) GetPartnerStats(ctx context.Context, partnerID string) (*service.PartnerStats, error) {
	stats, err := h.svc.GetPartnerStats(ctx, partnerID)
	if err != nil {
		log.Printf("[marketplace-handler] GetPartnerStats error: %v", err)
		return nil, err
	}
	return stats, nil
}

// UpdateProduct는 상품을 업데이트합니다.
func (h *MarketplaceHandler) UpdateProduct(ctx context.Context, product *service.PartnerProduct) (bool, error) {
	ok, err := h.svc.UpdateProduct(ctx, product)
	if err != nil {
		log.Printf("[marketplace-handler] UpdateProduct error: %v", err)
		return false, err
	}
	return ok, nil
}
