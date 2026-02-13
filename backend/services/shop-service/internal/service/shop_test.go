package service

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// --- 테스트용 가짜 저장소 ---

type fakeProductRepo struct {
	products map[string]*Product
}

func newFakeProductRepo() *fakeProductRepo {
	r := &fakeProductRepo{products: make(map[string]*Product)}
	// 기본 테스트 상품 추가
	r.products["prod-1"] = &Product{ID: "prod-1", Name: "혈당 카트리지 10팩", Category: CategoryCartridge, PriceKRW: 29000, Stock: 100, IsActive: true, CreatedAt: time.Now()}
	r.products["prod-2"] = &Product{ID: "prod-2", Name: "리더기 V2", Category: CategoryReader, PriceKRW: 199000, Stock: 50, IsActive: true, CreatedAt: time.Now()}
	r.products["prod-3"] = &Product{ID: "prod-3", Name: "단종 카트리지", Category: CategoryCartridge, PriceKRW: 15000, Stock: 0, IsActive: false, CreatedAt: time.Now()}
	return r
}

func (r *fakeProductRepo) List(_ context.Context, category ProductCategory, limit, offset int32) ([]*Product, int32, error) {
	var result []*Product
	for _, p := range r.products {
		if category == CategoryUnknown || p.Category == category {
			result = append(result, p)
		}
	}
	total := int32(len(result))
	start := int(offset)
	if start >= len(result) {
		return nil, total, nil
	}
	end := start + int(limit)
	if end > len(result) {
		end = len(result)
	}
	return result[start:end], total, nil
}

func (r *fakeProductRepo) GetByID(_ context.Context, id string) (*Product, error) {
	p, ok := r.products[id]
	if !ok {
		return nil, nil
	}
	return p, nil
}

func (r *fakeProductRepo) Create(_ context.Context, product *Product) error {
	r.products[product.ID] = product
	return nil
}

func (r *fakeProductRepo) UpdateStock(_ context.Context, id string, delta int32) error {
	if p, ok := r.products[id]; ok {
		p.Stock += delta
	}
	return nil
}

type fakeCartRepo struct {
	carts map[string][]*CartItem // userID -> items
}

func newFakeCartRepo() *fakeCartRepo {
	return &fakeCartRepo{carts: make(map[string][]*CartItem)}
}

func (r *fakeCartRepo) GetByUserID(_ context.Context, userID string) ([]*CartItem, error) {
	return r.carts[userID], nil
}

func (r *fakeCartRepo) AddItem(_ context.Context, item *CartItem) error {
	r.carts[item.UserID] = append(r.carts[item.UserID], item)
	return nil
}

func (r *fakeCartRepo) RemoveItem(_ context.Context, userID, cartItemID string) error {
	items := r.carts[userID]
	for i, item := range items {
		if item.ID == cartItemID {
			r.carts[userID] = append(items[:i], items[i+1:]...)
			return nil
		}
	}
	return nil
}

func (r *fakeCartRepo) Clear(_ context.Context, userID string) error {
	delete(r.carts, userID)
	return nil
}

type fakeOrderRepo struct {
	orders   map[string]*Order
	byUserID map[string][]*Order
}

func newFakeOrderRepo() *fakeOrderRepo {
	return &fakeOrderRepo{
		orders:   make(map[string]*Order),
		byUserID: make(map[string][]*Order),
	}
}

func (r *fakeOrderRepo) Create(_ context.Context, order *Order) error {
	r.orders[order.ID] = order
	r.byUserID[order.UserID] = append(r.byUserID[order.UserID], order)
	return nil
}

func (r *fakeOrderRepo) GetByID(_ context.Context, id string) (*Order, error) {
	return r.orders[id], nil
}

func (r *fakeOrderRepo) ListByUserID(_ context.Context, userID string, limit, offset int32) ([]*Order, int32, error) {
	orders := r.byUserID[userID]
	total := int32(len(orders))
	return orders, total, nil
}

func (r *fakeOrderRepo) UpdateStatus(_ context.Context, id string, status OrderStatus) error {
	if o, ok := r.orders[id]; ok {
		o.Status = status
	}
	return nil
}

func newTestShopService() *ShopService {
	return NewShopService(zap.NewNop(), newFakeProductRepo(), newFakeCartRepo(), newFakeOrderRepo())
}

// --- 테스트 ---

func TestListProducts_All(t *testing.T) {
	svc := newTestShopService()
	ctx := context.Background()

	products, total, err := svc.ListProducts(ctx, CategoryUnknown, 20, 0)
	if err != nil {
		t.Fatalf("ListProducts 실패: %v", err)
	}
	if total != 3 {
		t.Errorf("전체 상품 수: got %d, want 3", total)
	}
	if len(products) != 3 {
		t.Errorf("반환 상품 수: got %d, want 3", len(products))
	}
}

func TestListProducts_ByCategory(t *testing.T) {
	svc := newTestShopService()
	ctx := context.Background()

	products, _, err := svc.ListProducts(ctx, CategoryCartridge, 20, 0)
	if err != nil {
		t.Fatalf("ListProducts 실패: %v", err)
	}
	for _, p := range products {
		if p.Category != CategoryCartridge {
			t.Errorf("카테고리 필터 실패: got %d, want %d", p.Category, CategoryCartridge)
		}
	}
}

func TestGetProduct(t *testing.T) {
	svc := newTestShopService()
	ctx := context.Background()

	product, err := svc.GetProduct(ctx, "prod-1")
	if err != nil {
		t.Fatalf("GetProduct 실패: %v", err)
	}
	if product.Name != "혈당 카트리지 10팩" {
		t.Errorf("상품명: got %s, want 혈당 카트리지 10팩", product.Name)
	}
}

func TestGetProduct_NotFound(t *testing.T) {
	svc := newTestShopService()
	ctx := context.Background()

	_, err := svc.GetProduct(ctx, "non-existent")
	if err == nil {
		t.Error("존재하지 않는 상품 조회가 성공했습니다")
	}
}

func TestAddToCart(t *testing.T) {
	svc := newTestShopService()
	ctx := context.Background()

	userID := uuid.New().String()
	items, total, err := svc.AddToCart(ctx, userID, "prod-1", 2)
	if err != nil {
		t.Fatalf("AddToCart 실패: %v", err)
	}
	if len(items) != 1 {
		t.Errorf("장바구니 항목 수: got %d, want 1", len(items))
	}
	if total != 58000 { // 29000 * 2
		t.Errorf("장바구니 합계: got %d, want 58000", total)
	}
}

func TestAddToCart_InactiveProduct(t *testing.T) {
	svc := newTestShopService()
	ctx := context.Background()

	_, _, err := svc.AddToCart(ctx, "user-1", "prod-3", 1)
	if err == nil {
		t.Error("비활성 상품 장바구니 추가가 성공했습니다")
	}
}

func TestAddToCart_InvalidQuantity(t *testing.T) {
	svc := newTestShopService()
	ctx := context.Background()

	_, _, err := svc.AddToCart(ctx, "user-1", "prod-1", 0)
	if err == nil {
		t.Error("수량 0 장바구니 추가가 성공했습니다")
	}
}

func TestCreateOrder(t *testing.T) {
	svc := newTestShopService()
	ctx := context.Background()

	userID := uuid.New().String()
	// 장바구니에 상품 추가
	_, _, _ = svc.AddToCart(ctx, userID, "prod-1", 2) // 58000
	_, _, _ = svc.AddToCart(ctx, userID, "prod-2", 1) // 199000

	order, err := svc.CreateOrder(ctx, userID, "서울시 강남구 테헤란로 123", "card")
	if err != nil {
		t.Fatalf("CreateOrder 실패: %v", err)
	}
	if order.TotalPriceKRW != 257000 { // 58000 + 199000
		t.Errorf("주문 합계: got %d, want 257000", order.TotalPriceKRW)
	}
	if order.Status != OrderPending {
		t.Errorf("주문 상태: got %d, want %d", order.Status, OrderPending)
	}
	if len(order.Items) != 2 {
		t.Errorf("주문 항목 수: got %d, want 2", len(order.Items))
	}

	// 장바구니 비워졌는지 확인
	items, _, _ := svc.GetCart(ctx, userID)
	if len(items) != 0 {
		t.Errorf("주문 후 장바구니가 비워지지 않았습니다: %d items", len(items))
	}
}

func TestCreateOrder_EmptyCart(t *testing.T) {
	svc := newTestShopService()
	ctx := context.Background()

	_, err := svc.CreateOrder(ctx, "user-empty", "주소", "card")
	if err == nil {
		t.Error("빈 장바구니 주문이 성공했습니다")
	}
}

func TestGetOrder(t *testing.T) {
	svc := newTestShopService()
	ctx := context.Background()

	userID := uuid.New().String()
	_, _, _ = svc.AddToCart(ctx, userID, "prod-1", 1)
	created, _ := svc.CreateOrder(ctx, userID, "주소", "card")

	order, err := svc.GetOrder(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetOrder 실패: %v", err)
	}
	if order.ID != created.ID {
		t.Errorf("주문 ID: got %s, want %s", order.ID, created.ID)
	}
}

func TestListOrders(t *testing.T) {
	svc := newTestShopService()
	ctx := context.Background()

	userID := uuid.New().String()
	// 2개 주문 생성
	_, _, _ = svc.AddToCart(ctx, userID, "prod-1", 1)
	_, _ = svc.CreateOrder(ctx, userID, "주소1", "card")
	_, _, _ = svc.AddToCart(ctx, userID, "prod-2", 1)
	_, _ = svc.CreateOrder(ctx, userID, "주소2", "bank")

	orders, total, err := svc.ListOrders(ctx, userID, 20, 0)
	if err != nil {
		t.Fatalf("ListOrders 실패: %v", err)
	}
	if total != 2 {
		t.Errorf("주문 수: got %d, want 2", total)
	}
	if len(orders) != 2 {
		t.Errorf("반환 주문 수: got %d, want 2", len(orders))
	}
}

func TestRemoveFromCart(t *testing.T) {
	svc := newTestShopService()
	ctx := context.Background()

	userID := uuid.New().String()
	items, _, _ := svc.AddToCart(ctx, userID, "prod-1", 1)
	itemID := items[0].ID

	remaining, total, err := svc.RemoveFromCart(ctx, userID, itemID)
	if err != nil {
		t.Fatalf("RemoveFromCart 실패: %v", err)
	}
	if len(remaining) != 0 {
		t.Errorf("제거 후 항목 수: got %d, want 0", len(remaining))
	}
	if total != 0 {
		t.Errorf("제거 후 합계: got %d, want 0", total)
	}
}
