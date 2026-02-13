// Package handler는 shop-service의 gRPC 핸들러입니다.
package handler

import (
	"context"

	"github.com/manpasik/backend/services/shop-service/internal/service"
	apperrors "github.com/manpasik/backend/shared/errors"
	v1 "github.com/manpasik/backend/shared/gen/go/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ShopHandler는 ShopService gRPC 서버를 구현합니다.
type ShopHandler struct {
	v1.UnimplementedShopServiceServer
	svc *service.ShopService
	log *zap.Logger
}

// NewShopHandler는 ShopHandler를 생성합니다.
func NewShopHandler(svc *service.ShopService, log *zap.Logger) *ShopHandler {
	return &ShopHandler{svc: svc, log: log}
}

func (h *ShopHandler) ListProducts(ctx context.Context, req *v1.ListProductsRequest) (*v1.ListProductsResponse, error) {
	category := service.ProductCategory(req.Category)
	products, total, err := h.svc.ListProducts(ctx, category, req.Limit, req.Offset)
	if err != nil {
		return nil, toGRPC(err)
	}

	protoProducts := make([]*v1.Product, 0, len(products))
	for _, p := range products {
		protoProducts = append(protoProducts, productToProto(p))
	}

	return &v1.ListProductsResponse{Products: protoProducts, TotalCount: total}, nil
}

func (h *ShopHandler) GetProduct(ctx context.Context, req *v1.GetProductRequest) (*v1.Product, error) {
	if req == nil || req.ProductId == "" {
		return nil, status.Error(codes.InvalidArgument, "product_id는 필수입니다")
	}

	product, err := h.svc.GetProduct(ctx, req.ProductId)
	if err != nil {
		return nil, toGRPC(err)
	}
	return productToProto(product), nil
}

func (h *ShopHandler) AddToCart(ctx context.Context, req *v1.AddToCartRequest) (*v1.Cart, error) {
	if req == nil || req.UserId == "" || req.ProductId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id와 product_id는 필수입니다")
	}

	items, total, err := h.svc.AddToCart(ctx, req.UserId, req.ProductId, req.Quantity)
	if err != nil {
		return nil, toGRPC(err)
	}
	return cartToProto(req.UserId, items, total), nil
}

func (h *ShopHandler) GetCart(ctx context.Context, req *v1.GetCartRequest) (*v1.Cart, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id는 필수입니다")
	}

	items, total, err := h.svc.GetCart(ctx, req.UserId)
	if err != nil {
		return nil, toGRPC(err)
	}
	return cartToProto(req.UserId, items, total), nil
}

func (h *ShopHandler) RemoveFromCart(ctx context.Context, req *v1.RemoveFromCartRequest) (*v1.Cart, error) {
	if req == nil || req.UserId == "" || req.CartItemId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id와 cart_item_id는 필수입니다")
	}

	items, total, err := h.svc.RemoveFromCart(ctx, req.UserId, req.CartItemId)
	if err != nil {
		return nil, toGRPC(err)
	}
	return cartToProto(req.UserId, items, total), nil
}

func (h *ShopHandler) CreateOrder(ctx context.Context, req *v1.CreateOrderRequest) (*v1.Order, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id는 필수입니다")
	}

	order, err := h.svc.CreateOrder(ctx, req.UserId, req.ShippingAddress, req.PaymentMethod)
	if err != nil {
		return nil, toGRPC(err)
	}
	return orderToProto(order), nil
}

func (h *ShopHandler) GetOrder(ctx context.Context, req *v1.GetOrderRequest) (*v1.Order, error) {
	if req == nil || req.OrderId == "" {
		return nil, status.Error(codes.InvalidArgument, "order_id는 필수입니다")
	}

	order, err := h.svc.GetOrder(ctx, req.OrderId)
	if err != nil {
		return nil, toGRPC(err)
	}
	return orderToProto(order), nil
}

func (h *ShopHandler) ListOrders(ctx context.Context, req *v1.ListOrdersRequest) (*v1.ListOrdersResponse, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id는 필수입니다")
	}

	orders, total, err := h.svc.ListOrders(ctx, req.UserId, req.Limit, req.Offset)
	if err != nil {
		return nil, toGRPC(err)
	}

	protoOrders := make([]*v1.Order, 0, len(orders))
	for _, o := range orders {
		protoOrders = append(protoOrders, orderToProto(o))
	}

	return &v1.ListOrdersResponse{Orders: protoOrders, TotalCount: total}, nil
}

// --- 헬퍼 ---

func productToProto(p *service.Product) *v1.Product {
	return &v1.Product{
		ProductId:   p.ID,
		Name:        p.Name,
		Description: p.Description,
		Category:    v1.ProductCategory(p.Category),
		PriceKrw:    p.PriceKRW,
		Stock:       p.Stock,
		ImageUrl:    p.ImageURL,
		IsActive:    p.IsActive,
		CreatedAt:   timestamppb.New(p.CreatedAt),
	}
}

func cartToProto(userID string, items []*service.CartItem, total int32) *v1.Cart {
	protoItems := make([]*v1.CartItem, 0, len(items))
	for _, item := range items {
		protoItems = append(protoItems, &v1.CartItem{
			CartItemId:    item.ID,
			ProductId:     item.ProductID,
			ProductName:   item.ProductName,
			Quantity:      item.Quantity,
			UnitPriceKrw:  item.UnitPriceKRW,
			TotalPriceKrw: item.TotalPriceKRW,
		})
	}
	return &v1.Cart{UserId: userID, Items: protoItems, TotalPriceKrw: total}
}

func orderToProto(o *service.Order) *v1.Order {
	protoItems := make([]*v1.OrderItem, 0, len(o.Items))
	for _, item := range o.Items {
		protoItems = append(protoItems, &v1.OrderItem{
			ProductId:     item.ProductID,
			ProductName:   item.ProductName,
			Quantity:      item.Quantity,
			UnitPriceKrw:  item.UnitPriceKRW,
			TotalPriceKrw: item.TotalPriceKRW,
		})
	}
	return &v1.Order{
		OrderId:         o.ID,
		UserId:          o.UserID,
		Items:           protoItems,
		TotalPriceKrw:   o.TotalPriceKRW,
		Status:          v1.OrderStatus(o.Status),
		ShippingAddress: o.ShippingAddress,
		PaymentId:       o.PaymentID,
		CreatedAt:       timestamppb.New(o.CreatedAt),
		UpdatedAt:       timestamppb.New(o.UpdatedAt),
	}
}

func toGRPC(err error) error {
	if err == nil {
		return nil
	}
	if ae, ok := err.(*apperrors.AppError); ok {
		return ae.ToGRPC()
	}
	if s, ok := status.FromError(err); ok {
		return s.Err()
	}
	return status.Error(codes.Internal, "내부 오류가 발생했습니다")
}
