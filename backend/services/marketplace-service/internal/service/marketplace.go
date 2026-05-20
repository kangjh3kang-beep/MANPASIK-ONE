// Package service는 marketplace-service의 비즈니스 로직을 구현합니다.
package service

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// ============================================================================
// 도메인 상수
// ============================================================================

// validCategories는 허용된 상품 카테고리 목록입니다.
var validCategories = map[string]bool{
	"cartridge":    true,
	"device":       true,
	"accessory":    true,
	"subscription": true,
	"service":      true,
}

// defaultCommissionRates는 카테고리별 기본 수수료율입니다.
var defaultCommissionRates = map[string]float64{
	"cartridge":    0.15,
	"device":       0.10,
	"accessory":    0.20,
	"subscription": 0.25,
	"service":      0.20,
}

// 파트너 상태 상수
const (
	PartnerStatusPending   = "pending"
	PartnerStatusActive    = "active"
	PartnerStatusSuspended = "suspended"
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
	ErrInvalidInput      = errors.New("invalid input")
	ErrNotFound          = errors.New("not found")
	ErrInternal          = errors.New("internal error")
	ErrInvalidTransition = errors.New("invalid status transition")
	ErrPartnerNotActive  = errors.New("partner is not active")
	ErrInvalidCategory   = errors.New("invalid category")
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

	// 신규 파트너는 항상 pending 상태로 시작
	partner.Status = PartnerStatusPending
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

// ApprovePartner는 pending 파트너를 active로 승인합니다.
func (s *MarketplaceService) ApprovePartner(ctx context.Context, partnerID, approvedBy string) error {
	if partnerID == "" {
		return fmt.Errorf("%w: partner_id is required", ErrInvalidInput)
	}

	partner, err := s.partnerRepo.FindByID(ctx, partnerID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInternal, err)
	}
	if partner == nil {
		return fmt.Errorf("%w: partner %s not found", ErrNotFound, partnerID)
	}

	if partner.Status == PartnerStatusActive {
		return fmt.Errorf("%w: partner already active", ErrInvalidTransition)
	}
	if partner.Status == PartnerStatusSuspended {
		return fmt.Errorf("%w: suspended partner cannot be directly approved, use reactivate", ErrInvalidTransition)
	}

	partner.Status = PartnerStatusActive
	return s.partnerRepo.Save(ctx, partner)
}

// SuspendPartner는 active 파트너를 suspended로 중지합니다.
func (s *MarketplaceService) SuspendPartner(ctx context.Context, partnerID, reason string) error {
	if partnerID == "" {
		return fmt.Errorf("%w: partner_id is required", ErrInvalidInput)
	}

	partner, err := s.partnerRepo.FindByID(ctx, partnerID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInternal, err)
	}
	if partner == nil {
		return fmt.Errorf("%w: partner %s not found", ErrNotFound, partnerID)
	}

	if partner.Status == PartnerStatusSuspended {
		return fmt.Errorf("%w: partner already suspended", ErrInvalidTransition)
	}
	if partner.Status != PartnerStatusActive {
		return fmt.Errorf("%w: only active partners can be suspended", ErrInvalidTransition)
	}

	partner.Status = PartnerStatusSuspended
	return s.partnerRepo.Save(ctx, partner)
}

// CreateProduct는 새 상품을 생성합니다.
func (s *MarketplaceService) CreateProduct(ctx context.Context, product *PartnerProduct) (*PartnerProduct, error) {
	if product == nil {
		return nil, fmt.Errorf("%w: product is nil", ErrInvalidInput)
	}
	if product.PartnerID == "" {
		return nil, fmt.Errorf("%w: partner_id is required", ErrInvalidInput)
	}

	// 파트너 active 상태 확인
	partner, err := s.partnerRepo.FindByID(ctx, product.PartnerID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	if partner == nil {
		return nil, fmt.Errorf("%w: partner %s not found", ErrNotFound, product.PartnerID)
	}
	if partner.Status != PartnerStatusActive {
		return nil, ErrPartnerNotActive
	}

	// 상품 검증
	if len(product.Name) < 2 || len(product.Name) > 100 {
		return nil, fmt.Errorf("%w: product name must be 2-100 characters", ErrInvalidInput)
	}
	if product.Price <= 0 {
		return nil, fmt.Errorf("%w: price must be greater than 0", ErrInvalidInput)
	}
	if product.Category != "" && !validCategories[product.Category] {
		return nil, fmt.Errorf("%w: %q", ErrInvalidCategory, product.Category)
	}

	now := time.Now().UTC()
	product.CreatedAt = now
	product.UpdatedAt = now
	product.IsActive = true

	if err := s.productRepo.Save(ctx, product); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	return product, nil
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

// CalculateCommission은 상품의 수수료를 계산합니다.
func (s *MarketplaceService) CalculateCommission(product *PartnerProduct) float64 {
	if product == nil {
		return 0
	}
	rate, ok := defaultCommissionRates[product.Category]
	if !ok {
		rate = 0.20 // 기본 수수료율
	}
	return product.Price * rate
}
