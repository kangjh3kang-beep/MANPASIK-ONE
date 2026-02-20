// Package memory는 인메모리 마켓플레이스 저장소입니다 (개발/테스트용).
package memory

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/manpasik/backend/services/marketplace-service/internal/service"
)

// ============================================================================
// ProductRepository
// ============================================================================

// ProductRepository는 인메모리 상품 저장소입니다.
type ProductRepository struct {
	mu       sync.RWMutex
	products map[string]*service.PartnerProduct // key: ID
}

// NewProductRepository는 인메모리 ProductRepository를 생성합니다.
// 시드 데이터를 포함합니다.
func NewProductRepository() *ProductRepository {
	repo := &ProductRepository{
		products: make(map[string]*service.PartnerProduct),
	}

	now := time.Now().UTC()
	seeds := []*service.PartnerProduct{
		{
			ID: "prod-001", PartnerID: "partner-001",
			Name: "스마트 혈압계 X200", Description: "블루투스 연동 자동 혈압 측정기",
			Price: 89000, Category: "device", ImageURL: "https://img.example.com/bp-x200.png",
			IsActive: true, CreatedAt: now, UpdatedAt: now,
		},
		{
			ID: "prod-002", PartnerID: "partner-001",
			Name: "체성분 분석기 B100", Description: "가정용 체성분 분석 장비",
			Price: 149000, Category: "device", ImageURL: "https://img.example.com/body-b100.png",
			IsActive: true, CreatedAt: now, UpdatedAt: now,
		},
		{
			ID: "prod-003", PartnerID: "partner-002",
			Name: "프리미엄 건강검진 패키지", Description: "종합 건강검진 + AI 분석 리포트",
			Price: 350000, Category: "health_package", ImageURL: "https://img.example.com/pkg-premium.png",
			IsActive: true, CreatedAt: now, UpdatedAt: now,
		},
		{
			ID: "prod-004", PartnerID: "partner-002",
			Name: "영양제 정기구독", Description: "맞춤형 영양제 월간 배송 서비스",
			Price: 45000, Category: "supplement", ImageURL: "https://img.example.com/supp-monthly.png",
			IsActive: false, CreatedAt: now, UpdatedAt: now,
		},
	}

	for _, p := range seeds {
		repo.products[p.ID] = p
	}

	return repo
}

// FindByPartner는 파트너 ID로 상품을 조회합니다.
func (r *ProductRepository) FindByPartner(_ context.Context, partnerID string) ([]*service.PartnerProduct, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*service.PartnerProduct
	for _, p := range r.products {
		if p.PartnerID == partnerID {
			cp := *p
			result = append(result, &cp)
		}
	}
	return result, nil
}

// FindByCategory는 카테고리로 상품을 조회합니다.
func (r *ProductRepository) FindByCategory(_ context.Context, category string) ([]*service.PartnerProduct, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*service.PartnerProduct
	for _, p := range r.products {
		if p.Category == category {
			cp := *p
			result = append(result, &cp)
		}
	}
	return result, nil
}

// FindByPartnerAndCategory는 파트너 ID + 카테고리로 상품을 조회합니다.
func (r *ProductRepository) FindByPartnerAndCategory(_ context.Context, partnerID, category string) ([]*service.PartnerProduct, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*service.PartnerProduct
	for _, p := range r.products {
		if p.PartnerID == partnerID && p.Category == category {
			cp := *p
			result = append(result, &cp)
		}
	}
	return result, nil
}

// FindAll은 전체 상품을 조회합니다.
func (r *ProductRepository) FindAll(_ context.Context) ([]*service.PartnerProduct, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*service.PartnerProduct
	for _, p := range r.products {
		cp := *p
		result = append(result, &cp)
	}
	return result, nil
}

// FindByID는 상품 ID로 조회합니다.
func (r *ProductRepository) FindByID(_ context.Context, id string) (*service.PartnerProduct, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, ok := r.products[id]
	if !ok {
		return nil, nil
	}
	cp := *p
	return &cp, nil
}

// Save는 상품을 저장합니다.
func (r *ProductRepository) Save(_ context.Context, product *service.PartnerProduct) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *product
	r.products[product.ID] = &cp
	return nil
}

// Update는 상품 정보를 업데이트합니다.
func (r *ProductRepository) Update(_ context.Context, product *service.PartnerProduct) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *product
	r.products[product.ID] = &cp
	return nil
}

// ============================================================================
// PartnerRepository
// ============================================================================

// PartnerRepository는 인메모리 파트너 저장소입니다.
type PartnerRepository struct {
	mu       sync.RWMutex
	partners map[string]*service.Partner // key: ID
	nextID   int
}

// NewPartnerRepository는 인메모리 PartnerRepository를 생성합니다.
// 시드 데이터를 포함합니다.
func NewPartnerRepository() *PartnerRepository {
	repo := &PartnerRepository{
		partners: make(map[string]*service.Partner),
		nextID:   100,
	}

	now := time.Now().UTC()
	seeds := []*service.Partner{
		{
			ID: "partner-001", Name: "헬스디바이스코리아",
			Description: "스마트 헬스케어 디바이스 제조사",
			ContactEmail: "biz@healthdevice.kr", Status: "active", CreatedAt: now,
		},
		{
			ID: "partner-002", Name: "웰니스메디",
			Description: "건강검진 및 영양제 전문 기업",
			ContactEmail: "partner@wellnessmed.com", Status: "active", CreatedAt: now,
		},
	}

	for _, p := range seeds {
		repo.partners[p.ID] = p
	}

	return repo
}

// Save는 파트너를 저장합니다. ID가 비어 있으면 자동 생성합니다.
func (r *PartnerRepository) Save(_ context.Context, partner *service.Partner) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if partner.ID == "" {
		r.nextID++
		partner.ID = fmt.Sprintf("partner-%03d", r.nextID)
	}

	cp := *partner
	r.partners[partner.ID] = &cp
	return nil
}

// FindByID는 파트너 ID로 조회합니다.
func (r *PartnerRepository) FindByID(_ context.Context, id string) (*service.Partner, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, ok := r.partners[id]
	if !ok {
		return nil, nil
	}
	cp := *p
	return &cp, nil
}

// FindAll은 전체 파트너를 조회합니다.
func (r *PartnerRepository) FindAll(_ context.Context) ([]*service.Partner, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*service.Partner
	for _, p := range r.partners {
		cp := *p
		result = append(result, &cp)
	}
	return result, nil
}

// ============================================================================
// StatsRepository
// ============================================================================

// StatsRepository는 인메모리 파트너 통계 저장소입니다.
type StatsRepository struct {
	mu    sync.RWMutex
	stats map[string]*service.PartnerStats // key: PartnerID
}

// NewStatsRepository는 인메모리 StatsRepository를 생성합니다.
// 시드 데이터를 포함합니다.
func NewStatsRepository() *StatsRepository {
	repo := &StatsRepository{
		stats: make(map[string]*service.PartnerStats),
	}

	seeds := []*service.PartnerStats{
		{PartnerID: "partner-001", TotalProducts: 2, TotalOrders: 150, Revenue: 12350000},
		{PartnerID: "partner-002", TotalProducts: 2, TotalOrders: 85, Revenue: 32750000},
	}

	for _, s := range seeds {
		repo.stats[s.PartnerID] = s
	}

	return repo
}

// FindByPartnerID는 파트너 ID로 통계를 조회합니다.
func (r *StatsRepository) FindByPartnerID(_ context.Context, partnerID string) (*service.PartnerStats, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	s, ok := r.stats[partnerID]
	if !ok {
		return nil, nil
	}
	cp := *s
	return &cp, nil
}

// Save는 파트너 통계를 저장합니다.
func (r *StatsRepository) Save(_ context.Context, stats *service.PartnerStats) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *stats
	r.stats[stats.PartnerID] = &cp
	return nil
}
