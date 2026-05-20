// Package service: 동시성 안전 재고 관리 + 주문 상태머신 + 환불/추천
package service

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"
)

// ============================================================================
// 재고 관리 (동시성 안전)
// ============================================================================

// InventoryItem은 SKU별 재고입니다.
type InventoryItem struct {
	SKU         string
	Available   int
	Reserved    int
	LowStockThreshold int
	LastUpdated time.Time
}

// InventoryManager는 동시성 안전 재고 관리자입니다.
type InventoryManager struct {
	mu       sync.RWMutex
	items    map[string]*InventoryItem
	skuLocks map[string]*sync.Mutex
}

// NewInventoryManager는 새 재고 관리자를 생성합니다.
func NewInventoryManager() *InventoryManager {
	return &InventoryManager{
		items:    make(map[string]*InventoryItem),
		skuLocks: make(map[string]*sync.Mutex),
	}
}

// SetStock은 SKU의 재고량을 설정합니다.
func (m *InventoryManager) SetStock(sku string, available int, lowStockThreshold int) error {
	if sku == "" {
		return errors.New("sku required")
	}
	if available < 0 {
		return errors.New("available cannot be negative")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	item, ok := m.items[sku]
	if !ok {
		item = &InventoryItem{SKU: sku}
		m.items[sku] = item
		m.skuLocks[sku] = &sync.Mutex{}
	}
	item.Available = available
	item.LowStockThreshold = lowStockThreshold
	item.LastUpdated = time.Now().UTC()
	return nil
}

// itemAndLock은 SKU의 item과 SKU별 lock을 짧은 글로벌 lock 하에 조회합니다.
//
// 반환된 lock과 item은 호출자가 SKU lock을 잡은 뒤에만 변경해야 합니다.
func (m *InventoryManager) itemAndLock(sku string) (*InventoryItem, *sync.Mutex, error) {
	m.mu.RLock()
	item, ok := m.items[sku]
	skuLock := m.skuLocks[sku]
	m.mu.RUnlock()
	if !ok || skuLock == nil {
		return nil, nil, fmt.Errorf("sku %s not found", sku)
	}
	return item, skuLock, nil
}

// Reserve는 재고를 예약합니다 (주문 시).
func (m *InventoryManager) Reserve(sku string, quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be positive")
	}

	item, skuLock, err := m.itemAndLock(sku)
	if err != nil {
		return err
	}

	skuLock.Lock()
	defer skuLock.Unlock()

	if item.Available < quantity {
		return fmt.Errorf("insufficient stock: available=%d, requested=%d", item.Available, quantity)
	}

	item.Available -= quantity
	item.Reserved += quantity
	item.LastUpdated = time.Now().UTC()
	return nil
}

// Confirm은 예약을 확정합니다 (배송 시).
func (m *InventoryManager) Confirm(sku string, quantity int) error {
	item, skuLock, err := m.itemAndLock(sku)
	if err != nil {
		return err
	}

	skuLock.Lock()
	defer skuLock.Unlock()

	if item.Reserved < quantity {
		return fmt.Errorf("reserved too low: reserved=%d, confirm=%d", item.Reserved, quantity)
	}
	item.Reserved -= quantity
	item.LastUpdated = time.Now().UTC()
	return nil
}

// Release는 예약을 취소합니다 (주문 취소 시).
func (m *InventoryManager) Release(sku string, quantity int) error {
	item, skuLock, err := m.itemAndLock(sku)
	if err != nil {
		return err
	}

	skuLock.Lock()
	defer skuLock.Unlock()

	if item.Reserved < quantity {
		return fmt.Errorf("reserved too low: reserved=%d, release=%d", item.Reserved, quantity)
	}
	item.Reserved -= quantity
	item.Available += quantity
	item.LastUpdated = time.Now().UTC()
	return nil
}

// GetStock은 SKU의 재고를 조회합니다.
func (m *InventoryManager) GetStock(sku string) (*InventoryItem, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	item, ok := m.items[sku]
	if !ok {
		return nil, fmt.Errorf("sku %s not found", sku)
	}
	copy := *item
	return &copy, nil
}

// LowStockAlerts는 임계값 미만 재고 SKU 목록을 반환합니다.
func (m *InventoryManager) LowStockAlerts() []*InventoryItem {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var alerts []*InventoryItem
	for _, item := range m.items {
		if item.LowStockThreshold > 0 && item.Available < item.LowStockThreshold {
			c := *item
			alerts = append(alerts, &c)
		}
	}
	sort.Slice(alerts, func(i, j int) bool {
		return alerts[i].Available < alerts[j].Available
	})
	return alerts
}

// IsInStock은 재고 여부를 반환합니다.
func (m *InventoryManager) IsInStock(sku string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	item, ok := m.items[sku]
	return ok && item.Available > 0
}

// RemoveSKU는 더 이상 사용하지 않는 SKU의 재고와 lock을 정리합니다.
//
// 운영 종료된 상품 정리에 사용. 호출 시 진행 중인 예약/확정이 없어야 합니다.
func (m *InventoryManager) RemoveSKU(sku string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.items, sku)
	delete(m.skuLocks, sku)
}

// ============================================================================
// 주문 상태 머신 (8단계)
// ============================================================================

// OrderStatus는 주문 상태입니다.
const (
	OrderStatusPending    = "pending"     // 1. 결제 대기
	OrderStatusPaid       = "paid"        // 2. 결제 완료
	OrderStatusPacking    = "packing"     // 3. 포장 중
	OrderStatusShipping   = "shipping"    // 4. 배송 중
	OrderStatusDelivered  = "delivered"   // 5. 배송 완료
	OrderStatusCancelled  = "cancelled"   // 6. 취소됨
	OrderStatusRefunding  = "refunding"   // 7. 환불 처리 중
	OrderStatusRefunded   = "refunded"    // 8. 환불 완료
)

// orderTransitions는 허용된 주문 상태 전이입니다.
var orderTransitions = map[string][]string{
	OrderStatusPending:   {OrderStatusPaid, OrderStatusCancelled},
	OrderStatusPaid:      {OrderStatusPacking, OrderStatusCancelled, OrderStatusRefunding},
	OrderStatusPacking:   {OrderStatusShipping, OrderStatusCancelled, OrderStatusRefunding},
	OrderStatusShipping:  {OrderStatusDelivered, OrderStatusRefunding},
	OrderStatusDelivered: {OrderStatusRefunding},
	OrderStatusCancelled: {},
	OrderStatusRefunding: {OrderStatusRefunded},
	OrderStatusRefunded:  {},
}

// OrderStateMachine은 주문 상태 전이를 관리합니다.
type OrderStateMachine struct{}

// NewOrderStateMachine은 새 주문 상태 머신을 생성합니다.
func NewOrderStateMachine() *OrderStateMachine {
	return &OrderStateMachine{}
}

// CanTransition은 상태 전이가 가능한지 검증합니다.
func (sm *OrderStateMachine) CanTransition(from, to string) bool {
	allowed, ok := orderTransitions[from]
	if !ok {
		return false
	}
	for _, s := range allowed {
		if s == to {
			return true
		}
	}
	return false
}

// Transition은 상태 전이를 시도합니다.
func (sm *OrderStateMachine) Transition(from, to string) error {
	if !sm.CanTransition(from, to) {
		return fmt.Errorf("invalid transition: %s → %s", from, to)
	}
	return nil
}

// IsTerminal은 종착 상태인지 확인합니다.
func (sm *OrderStateMachine) IsTerminal(status string) bool {
	allowed, ok := orderTransitions[status]
	return ok && len(allowed) == 0
}

// AllowedTransitions는 현재 상태에서 가능한 다음 상태 목록을 반환합니다.
func (sm *OrderStateMachine) AllowedTransitions(from string) []string {
	return orderTransitions[from]
}

// ============================================================================
// 환불 처리
// ============================================================================

// RefundReason은 환불 사유입니다.
type RefundReason string

const (
	RefundReasonDefective       RefundReason = "defective"        // 불량
	RefundReasonWrongItem       RefundReason = "wrong_item"       // 오배송
	RefundReasonChangeOfMind    RefundReason = "change_of_mind"   // 단순 변심
	RefundReasonNotAsDescribed  RefundReason = "not_as_described" // 광고 다름
	RefundReasonDamaged         RefundReason = "damaged"          // 파손
)

// RefundRequest는 환불 요청입니다.
type RefundRequest struct {
	OrderID     string
	UserID      string
	Reason      RefundReason
	Description string
	Amount      float64
	Currency    string
	RequestedAt time.Time
}

// RefundProcessor는 환불 처리를 담당합니다.
type RefundProcessor struct {
	policies map[RefundReason]RefundPolicy
}

// RefundPolicy는 환불 사유별 정책입니다.
type RefundPolicy struct {
	FullRefundAllowed   bool
	RestockingFeeRate   float64 // 0~1, 재입고 수수료
	ReturnShippingPaidBy string  // "buyer", "seller"
	ApprovalAuto        bool
}

// NewRefundProcessor는 새 환불 처리자를 생성합니다.
func NewRefundProcessor() *RefundProcessor {
	return &RefundProcessor{
		policies: map[RefundReason]RefundPolicy{
			RefundReasonDefective:      {FullRefundAllowed: true, RestockingFeeRate: 0, ReturnShippingPaidBy: "seller", ApprovalAuto: true},
			RefundReasonWrongItem:      {FullRefundAllowed: true, RestockingFeeRate: 0, ReturnShippingPaidBy: "seller", ApprovalAuto: true},
			RefundReasonDamaged:        {FullRefundAllowed: true, RestockingFeeRate: 0, ReturnShippingPaidBy: "seller", ApprovalAuto: true},
			RefundReasonNotAsDescribed: {FullRefundAllowed: true, RestockingFeeRate: 0.05, ReturnShippingPaidBy: "buyer", ApprovalAuto: false},
			RefundReasonChangeOfMind:   {FullRefundAllowed: true, RestockingFeeRate: 0.10, ReturnShippingPaidBy: "buyer", ApprovalAuto: false},
		},
	}
}

// CalculateRefundAmount는 정책에 따른 환불 금액을 계산합니다.
func (rp *RefundProcessor) CalculateRefundAmount(req *RefundRequest) (float64, error) {
	if req == nil {
		return 0, errors.New("request is nil")
	}
	policy, ok := rp.policies[req.Reason]
	if !ok {
		return 0, fmt.Errorf("unknown refund reason: %s", req.Reason)
	}

	refund := req.Amount
	if !policy.FullRefundAllowed {
		refund = req.Amount * 0.5
	}
	// 재입고 수수료 차감
	if policy.RestockingFeeRate > 0 {
		refund -= req.Amount * policy.RestockingFeeRate
	}
	if refund < 0 {
		refund = 0
	}
	return refund, nil
}

// IsAutoApproved는 자동 승인 여부를 반환합니다.
func (rp *RefundProcessor) IsAutoApproved(reason RefundReason) bool {
	policy, ok := rp.policies[reason]
	return ok && policy.ApprovalAuto
}

// ============================================================================
// 할인 계산 (쿠폰 + 구독자 할인)
// ============================================================================

// DiscountType은 할인 유형입니다.
type DiscountType string

const (
	DiscountPercent     DiscountType = "percent"
	DiscountFixed       DiscountType = "fixed"
	DiscountSubscriber  DiscountType = "subscriber"
)

// Discount는 할인 항목입니다.
type Discount struct {
	Type      DiscountType
	Value     float64 // percent or fixed amount
	MaxAmount float64 // 최대 할인 금액 (percent 한정)
	MinOrder  float64 // 최소 주문 금액
}

// DiscountCalculator는 할인 계산기입니다.
type DiscountCalculator struct{}

// NewDiscountCalculator는 새 할인 계산기를 생성합니다.
func NewDiscountCalculator() *DiscountCalculator {
	return &DiscountCalculator{}
}

// Calculate는 주문 금액에 할인을 적용합니다.
func (c *DiscountCalculator) Calculate(orderAmount float64, discounts []Discount, isSubscriber bool) (finalAmount, totalDiscount float64) {
	finalAmount = orderAmount
	for _, d := range discounts {
		if orderAmount < d.MinOrder {
			continue
		}
		var amt float64
		switch d.Type {
		case DiscountPercent:
			amt = finalAmount * (d.Value / 100)
			if d.MaxAmount > 0 && amt > d.MaxAmount {
				amt = d.MaxAmount
			}
		case DiscountFixed:
			amt = d.Value
		case DiscountSubscriber:
			if !isSubscriber {
				continue
			}
			amt = finalAmount * (d.Value / 100)
		}
		if amt > finalAmount {
			amt = finalAmount
		}
		finalAmount -= amt
		totalDiscount += amt
	}
	return finalAmount, totalDiscount
}

// ============================================================================
// 추천 엔진 (간단 룰 기반)
// ============================================================================

// UserMeasurementProfile은 추천에 사용되는 사용자 측정 패턴입니다.
type UserMeasurementProfile struct {
	UserID                string
	DominantBiomarkers    []string // 자주 측정하는 바이오마커
	AvgFrequencyPerWeek   float64
	HasAbnormalReadings   bool
}

// ProductRecommendation은 추천 결과입니다.
type ProductRecommendation struct {
	ProductSKU  string
	Reason      string
	Score       float64 // 0~1
}

// RecommendationEngine은 사용자 측정 패턴 기반 카트리지 추천 엔진입니다.
type RecommendationEngine struct {
	biomarkerToSKU map[string][]string // biomarker -> SKU list
}

// NewRecommendationEngine은 새 추천 엔진을 생성합니다.
func NewRecommendationEngine() *RecommendationEngine {
	return &RecommendationEngine{
		biomarkerToSKU: map[string][]string{
			"glucose":       {"BIO-GLU-001", "BIO-GLU-002"},
			"hba1c":         {"BIO-HBA-001"},
			"cholesterol":   {"BIO-CHO-001", "BIO-CHO-002"},
			"blood_pressure": {"BIO-BP-001"},
			"oxygen":        {"BIO-OXY-001"},
		},
	}
}

// Recommend는 사용자 프로필 기반으로 카트리지를 추천합니다.
func (re *RecommendationEngine) Recommend(profile *UserMeasurementProfile, limit int) []ProductRecommendation {
	if profile == nil {
		return nil
	}
	var recs []ProductRecommendation

	// 주력 바이오마커 기반 추천
	for _, biomarker := range profile.DominantBiomarkers {
		skus, ok := re.biomarkerToSKU[biomarker]
		if !ok {
			continue
		}
		for i, sku := range skus {
			score := 0.9 - float64(i)*0.1 // 첫 번째 SKU가 더 높은 점수
			if profile.HasAbnormalReadings {
				score += 0.05
			}
			if score > 1.0 {
				score = 1.0
			}
			recs = append(recs, ProductRecommendation{
				ProductSKU: sku,
				Reason:     fmt.Sprintf("%s 측정에 적합", biomarker),
				Score:      score,
			})
		}
	}

	// 측정 빈도가 높으면 정기구독 추천
	if profile.AvgFrequencyPerWeek >= 3 {
		recs = append(recs, ProductRecommendation{
			ProductSKU: "SUB-MONTHLY-001",
			Reason:     "주 3회 이상 측정자에게 정기구독 권장",
			Score:      0.85,
		})
	}

	// 점수 내림차순 정렬
	sort.Slice(recs, func(i, j int) bool {
		return recs[i].Score > recs[j].Score
	})

	// 중복 SKU 제거
	seen := make(map[string]bool)
	unique := make([]ProductRecommendation, 0, len(recs))
	for _, r := range recs {
		if !seen[r.ProductSKU] {
			seen[r.ProductSKU] = true
			unique = append(unique, r)
		}
	}

	if limit > 0 && limit < len(unique) {
		return unique[:limit]
	}
	return unique
}
