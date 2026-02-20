// Package service는 marketplace-service의 비즈니스 로직을 구현합니다.
package service

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// ============================================================================
// 도메인 모델
// ============================================================================

// PartnerProduct는 파트너 상품 엔티티입니다.
type PartnerProduct struct {
	ID          string
	PartnerID   string
	Name        string
	Description string
	Price       float64
	Category    string
	ImageURL    string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Partner는 파트너 엔티티입니다.
type Partner struct {
	ID           string
	Name         string
	Description  string
	ContactEmail string
	Status       string // "pending", "active", "suspended"
	CreatedAt    time.Time
}

// PartnerStats는 파트너 통계 엔티티입니다.
type PartnerStats struct {
	PartnerID     string
	TotalProducts int
	TotalOrders   int
	Revenue       float64
}

// ============================================================================
// Repository 인터페이스
// ============================================================================

// ProductRepository는 상품 데이터 저장소 인터페이스입니다.
type ProductRepository interface {
	FindByPartner(ctx context.Context, partnerID string) ([]*PartnerProduct, error)
	FindByCategory(ctx context.Context, category string) ([]*PartnerProduct, error)
	FindByPartnerAndCategory(ctx context.Context, partnerID, category string) ([]*PartnerProduct, error)
	FindAll(ctx context.Context) ([]*PartnerProduct, error)
	FindByID(ctx context.Context, id string) (*PartnerProduct, error)
	Save(ctx context.Context, product *PartnerProduct) error
	Update(ctx context.Context, product *PartnerProduct) error
}

// PartnerRepository는 파트너 데이터 저장소 인터페이스입니다.
type PartnerRepository interface {
	Save(ctx context.Context, partner *Partner) error
	FindByID(ctx context.Context, id string) (*Partner, error)
	FindAll(ctx context.Context) ([]*Partner, error)
}

// StatsRepository는 파트너 통계 저장소 인터페이스입니다.
type StatsRepository interface {
	FindByPartnerID(ctx context.Context, partnerID string) (*PartnerStats, error)
	Save(ctx context.Context, stats *PartnerStats) error
}

// ============================================================================
// 에러 정의
// ============================================================================

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrNotFound     = errors.New("not found")
	ErrInternal     = errors.New("internal error")
)

// ============================================================================
// MarketplaceService
// ============================================================================

// MarketplaceService는 마켓플레이스 비즈니스 로직입니다.
type MarketplaceService struct {
	productRepo ProductRepository
	partnerRepo PartnerRepository
	statsRepo   StatsRepository
}

// NewMarketplaceService는 새 MarketplaceService를 생성합니다.
func NewMarketplaceService(
	productRepo ProductRepository,
	partnerRepo PartnerRepository,
	statsRepo StatsRepository,
) *MarketplaceService {
	return &MarketplaceService{
		productRepo: productRepo,
		partnerRepo: partnerRepo,
		statsRepo:   statsRepo,
	}
}

// ListPartnerProducts는 파트너 상품 목록을 조회합니다.
// partnerID와 category 모두 비어 있으면 전체 상품을 반환합니다.
func (s *MarketplaceService) ListPartnerProducts(ctx context.Context, partnerID, category string) ([]*PartnerProduct, error) {
	switch {
	case partnerID != "" && category != "":
		return s.productRepo.FindByPartnerAndCategory(ctx, partnerID, category)
	case partnerID != "":
		return s.productRepo.FindByPartner(ctx, partnerID)
	case category != "":
		return s.productRepo.FindByCategory(ctx, category)
	default:
		return s.productRepo.FindAll(ctx)
	}
}

// RegisterPartner는 새 파트너를 등록하고 파트너 ID를 반환합니다.
func (s *MarketplaceService) RegisterPartner(ctx context.Context, partner *Partner) (string, error) {
	if partner == nil {
		return "", fmt.Errorf("%w: partner is nil", ErrInvalidInput)
	}
	if partner.Name == "" {
		return "", fmt.Errorf("%w: partner name is required", ErrInvalidInput)
	}
	if partner.ContactEmail == "" {
		return "", fmt.Errorf("%w: contact email is required", ErrInvalidInput)
	}

	if partner.Status == "" {
		partner.Status = "pending"
	}
	if partner.CreatedAt.IsZero() {
		partner.CreatedAt = time.Now().UTC()
	}

	if err := s.partnerRepo.Save(ctx, partner); err != nil {
		return "", fmt.Errorf("%w: %v", ErrInternal, err)
	}

	// 파트너 통계 초기화
	initialStats := &PartnerStats{
		PartnerID:     partner.ID,
		TotalProducts: 0,
		TotalOrders:   0,
		Revenue:       0,
	}
	_ = s.statsRepo.Save(ctx, initialStats)

	return partner.ID, nil
}

// GetPartnerStats는 파트너 통계를 조회합니다.
func (s *MarketplaceService) GetPartnerStats(ctx context.Context, partnerID string) (*PartnerStats, error) {
	if partnerID == "" {
		return nil, fmt.Errorf("%w: partner_id is required", ErrInvalidInput)
	}

	stats, err := s.statsRepo.FindByPartnerID(ctx, partnerID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	if stats == nil {
		return nil, fmt.Errorf("%w: partner stats not found for %s", ErrNotFound, partnerID)
	}
	return stats, nil
}

// UpdateProduct는 상품 정보를 업데이트합니다.
func (s *MarketplaceService) UpdateProduct(ctx context.Context, product *PartnerProduct) (bool, error) {
	if product == nil {
		return false, fmt.Errorf("%w: product is nil", ErrInvalidInput)
	}
	if product.ID == "" {
		return false, fmt.Errorf("%w: product ID is required", ErrInvalidInput)
	}

	existing, err := s.productRepo.FindByID(ctx, product.ID)
	if err != nil {
		return false, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	if existing == nil {
		return false, fmt.Errorf("%w: product %s not found", ErrNotFound, product.ID)
	}

	product.UpdatedAt = time.Now().UTC()
	if err := s.productRepo.Update(ctx, product); err != nil {
		return false, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	return true, nil
}
