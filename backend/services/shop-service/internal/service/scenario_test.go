package service_test

import (
	"sync"
	"testing"

	"github.com/manpasik/backend/services/shop-service/internal/service"
)

// TestScenario_FullOrderJourney는 주문 → 결제 → 배송 → 환불 전체 흐름을 검증합니다.
func TestScenario_FullOrderJourney(t *testing.T) {
	inv := service.NewInventoryManager()
	sm := service.NewOrderStateMachine()
	rp := service.NewRefundProcessor()
	calc := service.NewDiscountCalculator()

	// 재고 설정
	_ = inv.SetStock("SKU-A", 100, 10)

	// 1. 주문 → 재고 예약
	if err := inv.Reserve("SKU-A", 5); err != nil {
		t.Fatalf("Reserve 실패: %v", err)
	}

	// 2. 할인 적용
	final, _ := calc.Calculate(50000, []service.Discount{
		{Type: service.DiscountSubscriber, Value: 10},
	}, true)
	if final != 45000 {
		t.Errorf("final = %f, want 45000", final)
	}

	// 3. 상태 전이 (정상 흐름)
	flow := []string{
		service.OrderStatusPaid,
		service.OrderStatusPacking,
		service.OrderStatusShipping,
		service.OrderStatusDelivered,
	}
	current := service.OrderStatusPending
	for _, next := range flow {
		if err := sm.Transition(current, next); err != nil {
			t.Fatalf("%s → %s 실패: %v", current, next, err)
		}
		current = next
	}

	// 4. 재고 확정
	_ = inv.Confirm("SKU-A", 5)

	// 5. 환불 흐름
	if err := sm.Transition(service.OrderStatusDelivered, service.OrderStatusRefunding); err != nil {
		t.Fatalf("환불 시작 실패: %v", err)
	}

	refundReq := &service.RefundRequest{Amount: 45000, Reason: service.RefundReasonDefective}
	refund, _ := rp.CalculateRefundAmount(refundReq)
	if refund != 45000 {
		t.Errorf("불량 전액 환불 실패: %f", refund)
	}
	_ = sm.Transition(service.OrderStatusRefunding, service.OrderStatusRefunded)
}

// TestScenario_InventoryConcurrentReservations는 다중 동시 예약을 검증합니다.
func TestScenario_InventoryConcurrentReservations(t *testing.T) {
	inv := service.NewInventoryManager()
	_ = inv.SetStock("CONCURRENT", 50, 5)

	var wg sync.WaitGroup
	successes := 0
	mu := sync.Mutex{}

	for i := 0; i < 30; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := inv.Reserve("CONCURRENT", 2); err == nil {
				mu.Lock()
				successes++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	stock, _ := inv.GetStock("CONCURRENT")
	expectedAvail := 50 - successes*2
	if stock.Available != expectedAvail {
		t.Errorf("Available = %d, want %d", stock.Available, expectedAvail)
	}
}

// TestScenario_RecommendationEnginePersonalized는 개인화 추천을 검증합니다.
func TestScenario_RecommendationEnginePersonalized(t *testing.T) {
	re := service.NewRecommendationEngine()

	// 케이스 1: 당뇨 환자 (혈당 측정 빈번 + 이상 수치)
	diabeticProfile := &service.UserMeasurementProfile{
		UserID:              "diabetic-001",
		DominantBiomarkers:  []string{"glucose", "hba1c"},
		AvgFrequencyPerWeek: 5,
		HasAbnormalReadings: true,
	}
	recs := re.Recommend(diabeticProfile, 10)

	hasGlucose := false
	hasSubscription := false
	for _, r := range recs {
		if r.ProductSKU == "BIO-GLU-001" || r.ProductSKU == "BIO-GLU-002" {
			hasGlucose = true
		}
		if r.ProductSKU == "SUB-MONTHLY-001" {
			hasSubscription = true
		}
	}
	if !hasGlucose {
		t.Error("당뇨 환자에게 혈당 카트리지 미추천")
	}
	if !hasSubscription {
		t.Error("고빈도 사용자에게 구독 미추천")
	}

	// 점수가 내림차순인지 검증
	for i := 1; i < len(recs); i++ {
		if recs[i].Score > recs[i-1].Score {
			t.Errorf("정렬 위반: [%d]=%f > [%d]=%f", i, recs[i].Score, i-1, recs[i-1].Score)
		}
	}
}

// TestScenario_DiscountStacking은 다중 할인 누적 적용을 검증합니다.
func TestScenario_DiscountStacking(t *testing.T) {
	calc := service.NewDiscountCalculator()

	final, totalDiscount := calc.Calculate(100000, []service.Discount{
		{Type: service.DiscountPercent, Value: 10},     // 10% (10000)
		{Type: service.DiscountFixed, Value: 5000},      // 5000원
		{Type: service.DiscountSubscriber, Value: 5},    // 구독자 5% (현재 85000의 5% = 4250)
	}, true)

	// 10% 할인 후: 90000, 5000 차감 후: 85000, 구독 5% 차감: 85000-4250=80750
	if final != 80750 {
		t.Errorf("final = %f, want 80750", final)
	}
	if totalDiscount != 19250 {
		t.Errorf("totalDiscount = %f, want 19250", totalDiscount)
	}
}

// TestScenario_OrderCancellationFlow는 결제 전 주문 취소를 검증합니다.
func TestScenario_OrderCancellationFlow(t *testing.T) {
	inv := service.NewInventoryManager()
	sm := service.NewOrderStateMachine()

	_ = inv.SetStock("CANCEL-A", 50, 5)
	_ = inv.Reserve("CANCEL-A", 3)

	// pending → cancelled
	if err := sm.Transition(service.OrderStatusPending, service.OrderStatusCancelled); err != nil {
		t.Fatalf("취소 실패: %v", err)
	}

	// 재고 복구
	_ = inv.Release("CANCEL-A", 3)

	stock, _ := inv.GetStock("CANCEL-A")
	if stock.Available != 50 {
		t.Errorf("Available = %d, want 50", stock.Available)
	}
}

// TestScenario_LowStockMonitoring은 재고 부족 알림 시나리오입니다.
func TestScenario_LowStockMonitoring(t *testing.T) {
	inv := service.NewInventoryManager()

	_ = inv.SetStock("LOW-A", 8, 10)
	_ = inv.SetStock("LOW-B", 3, 10)
	_ = inv.SetStock("OK-A", 50, 10)

	alerts := inv.LowStockAlerts()
	if len(alerts) != 2 {
		t.Errorf("alerts = %d, want 2", len(alerts))
	}

	// 가장 적은 재고가 첫 번째여야 함
	if alerts[0].Available != 3 {
		t.Errorf("first alert available = %d, want 3", alerts[0].Available)
	}
}

// TestScenario_RefundPolicyComparison은 환불 사유별 정책 차이를 검증합니다.
func TestScenario_RefundPolicyComparison(t *testing.T) {
	rp := service.NewRefundProcessor()
	amount := 10000.0

	cases := []struct {
		reason         service.RefundReason
		expectedRefund float64
		expectedAuto   bool
	}{
		{service.RefundReasonDefective, 10000, true},      // 100% + 자동
		{service.RefundReasonChangeOfMind, 9000, false},   // 90% + 수동
		{service.RefundReasonNotAsDescribed, 9500, false}, // 95% + 수동
	}

	for _, c := range cases {
		req := &service.RefundRequest{Amount: amount, Reason: c.reason}
		refund, _ := rp.CalculateRefundAmount(req)
		if refund != c.expectedRefund {
			t.Errorf("%s: refund = %f, want %f", c.reason, refund, c.expectedRefund)
		}
		if rp.IsAutoApproved(c.reason) != c.expectedAuto {
			t.Errorf("%s: auto = %v, want %v", c.reason, rp.IsAutoApproved(c.reason), c.expectedAuto)
		}
	}
}

// TestScenario_OrderTerminalStates는 종착 상태 전이 거부를 검증합니다.
func TestScenario_OrderTerminalStates(t *testing.T) {
	sm := service.NewOrderStateMachine()

	terminals := []string{
		service.OrderStatusCancelled,
		service.OrderStatusRefunded,
	}

	for _, terminal := range terminals {
		if !sm.IsTerminal(terminal) {
			t.Errorf("%s가 종착이 아님", terminal)
		}
		// 종착 상태에서 다른 상태로 전이 불가
		if err := sm.Transition(terminal, service.OrderStatusPaid); err == nil {
			t.Errorf("%s → paid 전이가 허용됨", terminal)
		}
	}
}
