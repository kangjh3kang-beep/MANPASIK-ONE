// Package memory는 인메모리 쇼핑 저장소입니다 (개발/테스트용).
package memory

import (
	"context"
	"sync"
	"time"

	"github.com/manpasik/backend/services/shop-service/internal/service"
)

// ProductRepository는 인메모리 상품 저장소입니다.
type ProductRepository struct {
	mu       sync.RWMutex
	products map[string]*service.Product
}

// NewProductRepository는 기본 상품이 포함된 인메모리 저장소를 생성합니다.
func NewProductRepository() *ProductRepository {
	r := &ProductRepository{
		products: make(map[string]*service.Product),
	}
	// 기본 상품 데이터
	now := time.Now()
	defaults := []*service.Product{
		{ID: "cart-blood-glucose-10", Name: "혈당 카트리지 10팩", Description: "혈당 측정용 카트리지 10개입", Category: service.CategoryCartridge, PriceKRW: 29000, Stock: 1000, IsActive: true, CreatedAt: now},
		{ID: "cart-cholesterol-10", Name: "콜레스테롤 카트리지 10팩", Description: "콜레스테롤 측정용 카트리지 10개입", Category: service.CategoryCartridge, PriceKRW: 35000, Stock: 800, IsActive: true, CreatedAt: now},
		{ID: "cart-hemoglobin-10", Name: "헤모글로빈 카트리지 10팩", Description: "헤모글로빈 측정용 카트리지 10개입", Category: service.CategoryCartridge, PriceKRW: 32000, Stock: 600, IsActive: true, CreatedAt: now},
		{ID: "reader-v2", Name: "ManPaSik 리더기 V2", Description: "차동측정 기반 범용 분석 리더기 2세대", Category: service.CategoryReader, PriceKRW: 199000, Stock: 200, IsActive: true, CreatedAt: now},
		{ID: "reader-v2-bundle", Name: "리더기 V2 + 카트리지 번들", Description: "리더기 V2 + 혈당·콜레스테롤 카트리지 각 10팩", Category: service.CategoryBundle, PriceKRW: 249000, Stock: 100, IsActive: true, CreatedAt: now},
		{ID: "acc-case", Name: "리더기 보호 케이스", Description: "프리미엄 실리콘 보호 케이스", Category: service.CategoryAccessory, PriceKRW: 25000, Stock: 500, IsActive: true, CreatedAt: now},
	}
	for _, p := range defaults {
		r.products[p.ID] = p
	}
	return r
}

func (r *ProductRepository) List(_ context.Context, category service.ProductCategory, limit, offset int32) ([]*service.Product, int32, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var filtered []*service.Product
	for _, p := range r.products {
		if category == service.CategoryUnknown || p.Category == category {
			if p.IsActive {
				cp := *p
				filtered = append(filtered, &cp)
			}
		}
	}

	total := int32(len(filtered))
	start := int(offset)
	if start >= len(filtered) {
		return nil, total, nil
	}
	end := start + int(limit)
	if end > len(filtered) {
		end = len(filtered)
	}
	return filtered[start:end], total, nil
}

func (r *ProductRepository) GetByID(_ context.Context, id string) (*service.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.products[id]
	if !ok {
		return nil, nil
	}
	cp := *p
	return &cp, nil
}

func (r *ProductRepository) Create(_ context.Context, product *service.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *product
	r.products[product.ID] = &cp
	return nil
}

func (r *ProductRepository) UpdateStock(_ context.Context, id string, delta int32) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if p, ok := r.products[id]; ok {
		p.Stock += delta
	}
	return nil
}

// CartRepository는 인메모리 장바구니 저장소입니다.
type CartRepository struct {
	mu    sync.RWMutex
	carts map[string][]*service.CartItem
}

func NewCartRepository() *CartRepository {
	return &CartRepository{carts: make(map[string][]*service.CartItem)}
}

func (r *CartRepository) GetByUserID(_ context.Context, userID string) ([]*service.CartItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := r.carts[userID]
	result := make([]*service.CartItem, len(items))
	for i, item := range items {
		cp := *item
		result[i] = &cp
	}
	return result, nil
}

func (r *CartRepository) AddItem(_ context.Context, item *service.CartItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *item
	r.carts[item.UserID] = append(r.carts[item.UserID], &cp)
	return nil
}

func (r *CartRepository) RemoveItem(_ context.Context, userID, cartItemID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	items := r.carts[userID]
	for i, item := range items {
		if item.ID == cartItemID {
			r.carts[userID] = append(items[:i], items[i+1:]...)
			return nil
		}
	}
	return nil
}

func (r *CartRepository) Clear(_ context.Context, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.carts, userID)
	return nil
}

// OrderRepository는 인메모리 주문 저장소입니다.
type OrderRepository struct {
	mu       sync.RWMutex
	orders   map[string]*service.Order
	byUserID map[string][]*service.Order
}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{
		orders:   make(map[string]*service.Order),
		byUserID: make(map[string][]*service.Order),
	}
}

func (r *OrderRepository) Create(_ context.Context, order *service.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *order
	r.orders[order.ID] = &cp
	r.byUserID[order.UserID] = append(r.byUserID[order.UserID], &cp)
	return nil
}

func (r *OrderRepository) GetByID(_ context.Context, id string) (*service.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	o, ok := r.orders[id]
	if !ok {
		return nil, nil
	}
	cp := *o
	return &cp, nil
}

func (r *OrderRepository) ListByUserID(_ context.Context, userID string, limit, offset int32) ([]*service.Order, int32, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	orders := r.byUserID[userID]
	total := int32(len(orders))
	return orders, total, nil
}

func (r *OrderRepository) UpdateStatus(_ context.Context, id string, s service.OrderStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if o, ok := r.orders[id]; ok {
		o.Status = s
	}
	return nil
}
