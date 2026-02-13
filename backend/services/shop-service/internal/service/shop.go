// Package service는 shop-service의 비즈니스 로직을 구현합니다.
package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	apperrors "github.com/manpasik/backend/shared/errors"
	"go.uber.org/zap"
)

// ProductCategory는 상품 카테고리입니다.
type ProductCategory int32

const (
	CategoryUnknown   ProductCategory = 0
	CategoryCartridge ProductCategory = 1
	CategoryReader    ProductCategory = 2
	CategoryAccessory ProductCategory = 3
	CategoryBundle    ProductCategory = 4
)

// OrderStatus는 주문 상태입니다.
type OrderStatus int32

const (
	OrderUnknown   OrderStatus = 0
	OrderPending   OrderStatus = 1
	OrderPaid      OrderStatus = 2
	OrderShipped   OrderStatus = 3
	OrderDelivered OrderStatus = 4
	OrderCancelled OrderStatus = 5
	OrderRefunded  OrderStatus = 6
)

// Product는 상품 엔티티입니다.
type Product struct {
	ID          string
	Name        string
	Description string
	Category    ProductCategory
	PriceKRW    int32
	Stock       int32
	ImageURL    string
	IsActive    bool
	CreatedAt   time.Time
}

// CartItem은 장바구니 항목입니다.
type CartItem struct {
	ID            string
	UserID        string
	ProductID     string
	ProductName   string
	Quantity      int32
	UnitPriceKRW  int32
	TotalPriceKRW int32
}

// Order는 주문 엔티티입니다.
type Order struct {
	ID              string
	UserID          string
	Items           []*OrderItem
	TotalPriceKRW   int32
	Status          OrderStatus
	ShippingAddress string
	PaymentID       string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// OrderItem은 주문 항목입니다.
type OrderItem struct {
	ProductID     string
	ProductName   string
	Quantity      int32
	UnitPriceKRW  int32
	TotalPriceKRW int32
}

// ProductRepository는 상품 저장소 인터페이스입니다.
type ProductRepository interface {
	List(ctx context.Context, category ProductCategory, limit, offset int32) ([]*Product, int32, error)
	GetByID(ctx context.Context, id string) (*Product, error)
	Create(ctx context.Context, product *Product) error
	UpdateStock(ctx context.Context, id string, delta int32) error
}

// CartRepository는 장바구니 저장소 인터페이스입니다.
type CartRepository interface {
	GetByUserID(ctx context.Context, userID string) ([]*CartItem, error)
	AddItem(ctx context.Context, item *CartItem) error
	RemoveItem(ctx context.Context, userID, cartItemID string) error
	Clear(ctx context.Context, userID string) error
}

// OrderRepository는 주문 저장소 인터페이스입니다.
type OrderRepository interface {
	Create(ctx context.Context, order *Order) error
	GetByID(ctx context.Context, id string) (*Order, error)
	ListByUserID(ctx context.Context, userID string, limit, offset int32) ([]*Order, int32, error)
	UpdateStatus(ctx context.Context, id string, status OrderStatus) error
}

// ShopService는 쇼핑 비즈니스 로직입니다.
type ShopService struct {
	logger      *zap.Logger
	productRepo ProductRepository
	cartRepo    CartRepository
	orderRepo   OrderRepository
}

// NewShopService는 새 ShopService를 생성합니다.
func NewShopService(
	logger *zap.Logger,
	productRepo ProductRepository,
	cartRepo CartRepository,
	orderRepo OrderRepository,
) *ShopService {
	return &ShopService{
		logger:      logger,
		productRepo: productRepo,
		cartRepo:    cartRepo,
		orderRepo:   orderRepo,
	}
}

// ListProducts는 상품 목록을 조회합니다.
func (s *ShopService) ListProducts(ctx context.Context, category ProductCategory, limit, offset int32) ([]*Product, int32, error) {
	if limit <= 0 {
		limit = 20
	}
	products, total, err := s.productRepo.List(ctx, category, limit, offset)
	if err != nil {
		s.logger.Error("상품 목록 조회 실패", zap.Error(err))
		return nil, 0, apperrors.New(apperrors.ErrInternal, "상품 목록 조회에 실패했습니다")
	}
	return products, total, nil
}

// GetProduct는 상품 상세를 조회합니다.
func (s *ShopService) GetProduct(ctx context.Context, productID string) (*Product, error) {
	product, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrInternal, "상품 조회에 실패했습니다")
	}
	if product == nil {
		return nil, apperrors.New(apperrors.ErrNotFound, "상품을 찾을 수 없습니다")
	}
	return product, nil
}

// AddToCart는 장바구니에 상품을 추가합니다.
func (s *ShopService) AddToCart(ctx context.Context, userID, productID string, quantity int32) ([]*CartItem, int32, error) {
	if quantity <= 0 {
		return nil, 0, apperrors.New(apperrors.ErrInvalidInput, "수량은 1 이상이어야 합니다")
	}

	product, err := s.productRepo.GetByID(ctx, productID)
	if err != nil || product == nil {
		return nil, 0, apperrors.New(apperrors.ErrNotFound, "상품을 찾을 수 없습니다")
	}

	if !product.IsActive {
		return nil, 0, apperrors.New(apperrors.ErrInvalidInput, "판매 중지된 상품입니다")
	}

	if product.Stock < quantity {
		return nil, 0, apperrors.New(apperrors.ErrInvalidInput, "재고가 부족합니다")
	}

	item := &CartItem{
		ID:            uuid.New().String(),
		UserID:        userID,
		ProductID:     productID,
		ProductName:   product.Name,
		Quantity:      quantity,
		UnitPriceKRW:  product.PriceKRW,
		TotalPriceKRW: product.PriceKRW * quantity,
	}

	if err := s.cartRepo.AddItem(ctx, item); err != nil {
		return nil, 0, apperrors.New(apperrors.ErrInternal, "장바구니 추가에 실패했습니다")
	}

	return s.getCartWithTotal(ctx, userID)
}

// GetCart는 장바구니를 조회합니다.
func (s *ShopService) GetCart(ctx context.Context, userID string) ([]*CartItem, int32, error) {
	return s.getCartWithTotal(ctx, userID)
}

// RemoveFromCart는 장바구니 항목을 제거합니다.
func (s *ShopService) RemoveFromCart(ctx context.Context, userID, cartItemID string) ([]*CartItem, int32, error) {
	if err := s.cartRepo.RemoveItem(ctx, userID, cartItemID); err != nil {
		return nil, 0, apperrors.New(apperrors.ErrInternal, "장바구니 항목 제거에 실패했습니다")
	}
	return s.getCartWithTotal(ctx, userID)
}

// CreateOrder는 장바구니 기반 주문을 생성합니다.
func (s *ShopService) CreateOrder(ctx context.Context, userID, shippingAddress, paymentMethod string) (*Order, error) {
	items, err := s.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrInternal, "장바구니 조회에 실패했습니다")
	}

	if len(items) == 0 {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "장바구니가 비어 있습니다")
	}

	var totalPrice int32
	orderItems := make([]*OrderItem, 0, len(items))
	for _, item := range items {
		orderItems = append(orderItems, &OrderItem{
			ProductID:     item.ProductID,
			ProductName:   item.ProductName,
			Quantity:      item.Quantity,
			UnitPriceKRW:  item.UnitPriceKRW,
			TotalPriceKRW: item.TotalPriceKRW,
		})
		totalPrice += item.TotalPriceKRW
	}

	now := time.Now().UTC()
	order := &Order{
		ID:              uuid.New().String(),
		UserID:          userID,
		Items:           orderItems,
		TotalPriceKRW:   totalPrice,
		Status:          OrderPending,
		ShippingAddress: shippingAddress,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, apperrors.New(apperrors.ErrInternal, "주문 생성에 실패했습니다")
	}

	// 장바구니 비우기
	_ = s.cartRepo.Clear(ctx, userID)

	s.logger.Info("주문 생성 완료",
		zap.String("order_id", order.ID),
		zap.String("user_id", userID),
		zap.Int32("total", totalPrice),
	)
	return order, nil
}

// GetOrder는 주문 상세를 조회합니다.
func (s *ShopService) GetOrder(ctx context.Context, orderID string) (*Order, error) {
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrInternal, "주문 조회에 실패했습니다")
	}
	if order == nil {
		return nil, apperrors.New(apperrors.ErrNotFound, "주문을 찾을 수 없습니다")
	}
	return order, nil
}

// ListOrders는 주문 이력을 조회합니다.
func (s *ShopService) ListOrders(ctx context.Context, userID string, limit, offset int32) ([]*Order, int32, error) {
	if limit <= 0 {
		limit = 20
	}
	orders, total, err := s.orderRepo.ListByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, 0, apperrors.New(apperrors.ErrInternal, "주문 이력 조회에 실패했습니다")
	}
	return orders, total, nil
}

func (s *ShopService) getCartWithTotal(ctx context.Context, userID string) ([]*CartItem, int32, error) {
	items, err := s.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, 0, apperrors.New(apperrors.ErrInternal, "장바구니 조회에 실패했습니다")
	}

	var total int32
	for _, item := range items {
		total += item.TotalPriceKRW
	}

	return items, total, nil
}
