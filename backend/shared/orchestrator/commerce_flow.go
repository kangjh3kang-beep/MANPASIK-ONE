package orchestrator

import (
	"context"
	"fmt"
	"log"
	"time"
)

// CommerceFlowOrchestrator coordinates subscription → shop → payment
type CommerceFlowOrchestrator struct {
	subscriptionChecker SubscriptionChecker
	orderCreator        OrderCreator
	paymentProcessor    PaymentProcessor
}

// SubscriptionChecker validates subscription access
type SubscriptionChecker interface {
	CheckFeatureAccess(ctx context.Context, userID, feature string) (bool, string, error)
	GetSubscriptionTier(ctx context.Context, userID string) (string, error)
}

// OrderCreator creates orders
type OrderCreator interface {
	CreateOrder(ctx context.Context, userID string, items []OrderItem) (*Order, error)
}

// OrderItem represents an item to order
type OrderItem struct {
	ProductID string
	Quantity  int32
	Price     float64
}

// Order represents a created order
type Order struct {
	OrderID     string
	UserID      string
	TotalAmount float64
	Status      string
	Items       []OrderItem
	CreatedAt   time.Time
}

// PaymentProcessor processes payments
type PaymentProcessor interface {
	ProcessPayment(ctx context.Context, orderID, userID string, amount float64, method string) (*PaymentResult, error)
}

// PaymentResult from payment processing
type PaymentResult struct {
	PaymentID string
	OrderID   string
	Status    string
	Amount    float64
	PaidAt    time.Time
}

// NewCommerceFlowOrchestrator creates a new commerce orchestrator
func NewCommerceFlowOrchestrator(sc SubscriptionChecker, oc OrderCreator, pp PaymentProcessor) *CommerceFlowOrchestrator {
	return &CommerceFlowOrchestrator{
		subscriptionChecker: sc,
		orderCreator:        oc,
		paymentProcessor:    pp,
	}
}

// PurchaseCartridge orchestrates the full cartridge purchase flow
func (o *CommerceFlowOrchestrator) PurchaseCartridge(ctx context.Context, userID, cartridgeProductID string, quantity int32, paymentMethod string) (*PurchaseResult, error) {
	log.Printf("[CommerceFlow] Starting cartridge purchase for user=%s", userID)

	// Step 1: Check subscription allows cartridge purchase
	hasAccess, tier, err := o.subscriptionChecker.CheckFeatureAccess(ctx, userID, "cartridge_purchase")
	if err != nil {
		return nil, fmt.Errorf("구독 확인 실패: %w", err)
	}
	if !hasAccess {
		return nil, fmt.Errorf("카트리지 구매 권한 없음 (현재 구독: %s)", tier)
	}
	log.Printf("[CommerceFlow] Step 1: Subscription verified (tier=%s)", tier)

	// Step 2: Apply tier discount
	basePrice := 29900.0 // default cartridge price
	discount := tierDiscount(tier)
	finalPrice := basePrice * (1 - discount)

	// Step 3: Create order
	order, err := o.orderCreator.CreateOrder(ctx, userID, []OrderItem{
		{ProductID: cartridgeProductID, Quantity: quantity, Price: finalPrice},
	})
	if err != nil {
		return nil, fmt.Errorf("주문 생성 실패: %w", err)
	}
	log.Printf("[CommerceFlow] Step 2-3: Order created (id=%s, total=%.0f)", order.OrderID, order.TotalAmount)

	// Step 4: Process payment
	payment, err := o.paymentProcessor.ProcessPayment(ctx, order.OrderID, userID, order.TotalAmount, paymentMethod)
	if err != nil {
		return nil, fmt.Errorf("결제 실패: %w", err)
	}
	log.Printf("[CommerceFlow] Step 4: Payment processed (id=%s, status=%s)", payment.PaymentID, payment.Status)

	return &PurchaseResult{
		OrderID:   order.OrderID,
		PaymentID: payment.PaymentID,
		Amount:    payment.Amount,
		Discount:  discount * 100,
		Tier:      tier,
		Status:    payment.Status,
	}, nil
}

// PurchaseResult is the final result of a purchase
type PurchaseResult struct {
	OrderID   string
	PaymentID string
	Amount    float64
	Discount  float64
	Tier      string
	Status    string
}

func tierDiscount(tier string) float64 {
	switch tier {
	case "premium":
		return 0.15
	case "professional":
		return 0.25
	case "enterprise":
		return 0.30
	default:
		return 0.0
	}
}
