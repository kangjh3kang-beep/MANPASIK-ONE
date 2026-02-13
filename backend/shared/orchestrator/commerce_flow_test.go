package orchestrator

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// --- Mock implementations for Commerce Flow ---

type mockSubscriptionChecker struct {
	hasAccess bool
	tier      string
	err       error
}

func (m *mockSubscriptionChecker) CheckFeatureAccess(ctx context.Context, userID, feature string) (bool, string, error) {
	return m.hasAccess, m.tier, m.err
}

func (m *mockSubscriptionChecker) GetSubscriptionTier(ctx context.Context, userID string) (string, error) {
	return m.tier, m.err
}

type mockOrderCreator struct {
	order *Order
	err   error
}

func (m *mockOrderCreator) CreateOrder(ctx context.Context, userID string, items []OrderItem) (*Order, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.order, nil
}

type mockPaymentProcessor struct {
	result *PaymentResult
	err    error
}

func (m *mockPaymentProcessor) ProcessPayment(ctx context.Context, orderID, userID string, amount float64, method string) (*PaymentResult, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.result, nil
}

// --- Tests ---

func TestPurchaseCartridge_Success(t *testing.T) {
	sc := &mockSubscriptionChecker{hasAccess: true, tier: "premium"}
	oc := &mockOrderCreator{order: &Order{
		OrderID:     "order-001",
		UserID:      "user-123",
		TotalAmount: 25415, // 29900 * 0.85
		Status:      "created",
		CreatedAt:   time.Now(),
	}}
	pp := &mockPaymentProcessor{result: &PaymentResult{
		PaymentID: "pay-001",
		OrderID:   "order-001",
		Status:    "completed",
		Amount:    25415,
		PaidAt:    time.Now(),
	}}

	orchestrator := NewCommerceFlowOrchestrator(sc, oc, pp)
	result, err := orchestrator.PurchaseCartridge(context.Background(), "user-123", "cartridge-abc", 1, "card")

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if result.OrderID != "order-001" {
		t.Errorf("expected OrderID=order-001, got %s", result.OrderID)
	}
	if result.PaymentID != "pay-001" {
		t.Errorf("expected PaymentID=pay-001, got %s", result.PaymentID)
	}
	if result.Tier != "premium" {
		t.Errorf("expected Tier=premium, got %s", result.Tier)
	}
	if result.Discount != 15 {
		t.Errorf("expected Discount=15, got %.0f", result.Discount)
	}
	if result.Status != "completed" {
		t.Errorf("expected Status=completed, got %s", result.Status)
	}
}

func TestPurchaseCartridge_NoAccess(t *testing.T) {
	sc := &mockSubscriptionChecker{hasAccess: false, tier: "free"}
	oc := &mockOrderCreator{}
	pp := &mockPaymentProcessor{}

	orchestrator := NewCommerceFlowOrchestrator(sc, oc, pp)
	result, err := orchestrator.PurchaseCartridge(context.Background(), "user-456", "cartridge-abc", 1, "card")

	if err == nil {
		t.Fatal("expected error for no access, got nil")
	}
	if result != nil {
		t.Errorf("expected nil result, got %+v", result)
	}
	expected := "카트리지 구매 권한 없음 (현재 구독: free)"
	if err.Error() != expected {
		t.Errorf("expected error %q, got %q", expected, err.Error())
	}
}

func TestPurchaseCartridge_PaymentFailed(t *testing.T) {
	sc := &mockSubscriptionChecker{hasAccess: true, tier: "premium"}
	oc := &mockOrderCreator{order: &Order{
		OrderID:     "order-002",
		UserID:      "user-789",
		TotalAmount: 25415,
		Status:      "created",
		CreatedAt:   time.Now(),
	}}
	pp := &mockPaymentProcessor{err: fmt.Errorf("insufficient funds")}

	orchestrator := NewCommerceFlowOrchestrator(sc, oc, pp)
	result, err := orchestrator.PurchaseCartridge(context.Background(), "user-789", "cartridge-abc", 1, "card")

	if err == nil {
		t.Fatal("expected payment error, got nil")
	}
	if result != nil {
		t.Errorf("expected nil result, got %+v", result)
	}
	if err.Error() != "결제 실패: insufficient funds" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}
