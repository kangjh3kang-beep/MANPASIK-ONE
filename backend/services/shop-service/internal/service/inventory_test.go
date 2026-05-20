package service_test

import (
	"sync"
	"testing"

	"github.com/manpasik/backend/services/shop-service/internal/service"
)

// TestInventoryManager_SetAndGetStock은 재고 설정/조회를 검증합니다.
func TestInventoryManager_SetAndGetStock(t *testing.T) {
	m := service.NewInventoryManager()
	if err := m.SetStock("SKU-001", 100, 10); err != nil {
		t.Fatalf("SetStock 실패: %v", err)
	}

	stock, err := m.GetStock("SKU-001")
	if err != nil {
		t.Fatalf("GetStock 실패: %v", err)
	}
	if stock.Available != 100 {
		t.Errorf("Available = %d, want 100", stock.Available)
	}
}

// TestInventoryManager_Reserve는 재고 예약을 검증합니다.
func TestInventoryManager_Reserve(t *testing.T) {
	m := service.NewInventoryManager()
	_ = m.SetStock("SKU", 100, 10)

	if err := m.Reserve("SKU", 30); err != nil {
		t.Fatalf("Reserve 실패: %v", err)
	}

	stock, _ := m.GetStock("SKU")
	if stock.Available != 70 {
		t.Errorf("Available = %d, want 70", stock.Available)
	}
	if stock.Reserved != 30 {
		t.Errorf("Reserved = %d, want 30", stock.Reserved)
	}
}

// TestInventoryManager_InsufficientStock는 재고 부족 거부를 검증합니다.
func TestInventoryManager_InsufficientStock(t *testing.T) {
	m := service.NewInventoryManager()
	_ = m.SetStock("SKU", 10, 5)

	err := m.Reserve("SKU", 20)
	if err == nil {
		t.Error("재고 초과 예약이 허용됨")
	}
}

// TestInventoryManager_Release는 예약 취소를 검증합니다.
func TestInventoryManager_Release(t *testing.T) {
	m := service.NewInventoryManager()
	_ = m.SetStock("SKU", 100, 10)
	_ = m.Reserve("SKU", 30)

	if err := m.Release("SKU", 30); err != nil {
		t.Fatalf("Release 실패: %v", err)
	}

	stock, _ := m.GetStock("SKU")
	if stock.Available != 100 {
		t.Errorf("Available = %d, want 100", stock.Available)
	}
	if stock.Reserved != 0 {
		t.Errorf("Reserved = %d, want 0", stock.Reserved)
	}
}

// TestInventoryManager_Confirm는 예약 확정을 검증합니다.
func TestInventoryManager_Confirm(t *testing.T) {
	m := service.NewInventoryManager()
	_ = m.SetStock("SKU", 100, 10)
	_ = m.Reserve("SKU", 30)

	if err := m.Confirm("SKU", 30); err != nil {
		t.Fatalf("Confirm 실패: %v", err)
	}

	stock, _ := m.GetStock("SKU")
	if stock.Available != 70 {
		t.Errorf("Available = %d, want 70", stock.Available)
	}
	if stock.Reserved != 0 {
		t.Errorf("Reserved = %d, want 0", stock.Reserved)
	}
}

// TestInventoryManager_ConcurrentReserve는 동시 예약의 안전성을 검증합니다.
func TestInventoryManager_ConcurrentReserve(t *testing.T) {
	m := service.NewInventoryManager()
	_ = m.SetStock("SKU", 100, 10)

	var wg sync.WaitGroup
	successes := 0
	mu := sync.Mutex{}
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := m.Reserve("SKU", 3); err == nil {
				mu.Lock()
				successes++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	stock, _ := m.GetStock("SKU")
	expectedAvailable := 100 - successes*3
	if stock.Available != expectedAvailable {
		t.Errorf("Available = %d, want %d (successes=%d)", stock.Available, expectedAvailable, successes)
	}
	if stock.Reserved != successes*3 {
		t.Errorf("Reserved = %d, want %d", stock.Reserved, successes*3)
	}
}

// TestInventoryManager_LowStockAlerts는 임계값 미만 알림을 검증합니다.
func TestInventoryManager_LowStockAlerts(t *testing.T) {
	m := service.NewInventoryManager()
	_ = m.SetStock("LOW-1", 5, 10)
	_ = m.SetStock("LOW-2", 8, 10)
	_ = m.SetStock("OK", 50, 10)

	alerts := m.LowStockAlerts()
	if len(alerts) != 2 {
		t.Errorf("alerts = %d, want 2", len(alerts))
	}
}

// TestOrderStateMachine_HappyPath는 정상 주문 흐름을 검증합니다.
func TestOrderStateMachine_HappyPath(t *testing.T) {
	sm := service.NewOrderStateMachine()
	flow := []string{
		service.OrderStatusPending,
		service.OrderStatusPaid,
		service.OrderStatusPacking,
		service.OrderStatusShipping,
		service.OrderStatusDelivered,
	}

	for i := 0; i < len(flow)-1; i++ {
		if err := sm.Transition(flow[i], flow[i+1]); err != nil {
			t.Errorf("%s → %s 실패: %v", flow[i], flow[i+1], err)
		}
	}
}

// TestOrderStateMachine_InvalidTransition는 잘못된 전이를 거부합니다.
func TestOrderStateMachine_InvalidTransition(t *testing.T) {
	sm := service.NewOrderStateMachine()
	err := sm.Transition(service.OrderStatusPending, service.OrderStatusDelivered)
	if err == nil {
		t.Error("pending → delivered 직접 전이가 허용됨")
	}
}

// TestOrderStateMachine_Refund는 환불 흐름을 검증합니다.
func TestOrderStateMachine_Refund(t *testing.T) {
	sm := service.NewOrderStateMachine()
	if err := sm.Transition(service.OrderStatusDelivered, service.OrderStatusRefunding); err != nil {
		t.Errorf("delivered → refunding 실패: %v", err)
	}
	if err := sm.Transition(service.OrderStatusRefunding, service.OrderStatusRefunded); err != nil {
		t.Errorf("refunding → refunded 실패: %v", err)
	}
}

// TestOrderStateMachine_Terminal는 종착 상태를 검증합니다.
func TestOrderStateMachine_Terminal(t *testing.T) {
	sm := service.NewOrderStateMachine()
	if !sm.IsTerminal(service.OrderStatusCancelled) {
		t.Error("cancelled가 종착으로 인식되지 않음")
	}
	if !sm.IsTerminal(service.OrderStatusRefunded) {
		t.Error("refunded가 종착으로 인식되지 않음")
	}
	if sm.IsTerminal(service.OrderStatusPaid) {
		t.Error("paid가 종착으로 잘못 인식됨")
	}
}

// TestRefundProcessor_DefectiveFullRefund는 불량 전액 환불을 검증합니다.
func TestRefundProcessor_DefectiveFullRefund(t *testing.T) {
	rp := service.NewRefundProcessor()
	req := &service.RefundRequest{Amount: 10000, Reason: service.RefundReasonDefective}

	refund, err := rp.CalculateRefundAmount(req)
	if err != nil {
		t.Fatalf("CalculateRefundAmount 실패: %v", err)
	}
	if refund != 10000 {
		t.Errorf("불량 환불 = %f, want 10000", refund)
	}
	if !rp.IsAutoApproved(service.RefundReasonDefective) {
		t.Error("불량 자동 승인 아님")
	}
}

// TestRefundProcessor_ChangeOfMindRestockingFee는 단순 변심 수수료 차감을 검증합니다.
func TestRefundProcessor_ChangeOfMindRestockingFee(t *testing.T) {
	rp := service.NewRefundProcessor()
	req := &service.RefundRequest{Amount: 10000, Reason: service.RefundReasonChangeOfMind}

	refund, _ := rp.CalculateRefundAmount(req)
	// 10000 - 10000*0.10 = 9000
	if refund != 9000 {
		t.Errorf("변심 환불 = %f, want 9000", refund)
	}
	if rp.IsAutoApproved(service.RefundReasonChangeOfMind) {
		t.Error("변심은 자동 승인 아니어야 함")
	}
}

// TestDiscountCalculator_Percent는 퍼센트 할인을 검증합니다.
func TestDiscountCalculator_Percent(t *testing.T) {
	calc := service.NewDiscountCalculator()
	final, discount := calc.Calculate(10000, []service.Discount{
		{Type: service.DiscountPercent, Value: 10},
	}, false)

	if discount != 1000 {
		t.Errorf("discount = %f, want 1000", discount)
	}
	if final != 9000 {
		t.Errorf("final = %f, want 9000", final)
	}
}

// TestDiscountCalculator_PercentMaxAmount는 퍼센트 할인 최대 금액 제한을 검증합니다.
func TestDiscountCalculator_PercentMaxAmount(t *testing.T) {
	calc := service.NewDiscountCalculator()
	_, discount := calc.Calculate(100000, []service.Discount{
		{Type: service.DiscountPercent, Value: 20, MaxAmount: 5000},
	}, false)

	if discount != 5000 {
		t.Errorf("discount = %f, want 5000 (최대 제한)", discount)
	}
}

// TestDiscountCalculator_Subscriber는 구독자 할인 적용을 검증합니다.
func TestDiscountCalculator_Subscriber(t *testing.T) {
	calc := service.NewDiscountCalculator()
	_, dSub := calc.Calculate(10000, []service.Discount{
		{Type: service.DiscountSubscriber, Value: 15},
	}, true)
	_, dNon := calc.Calculate(10000, []service.Discount{
		{Type: service.DiscountSubscriber, Value: 15},
	}, false)

	if dSub != 1500 {
		t.Errorf("구독자 할인 = %f, want 1500", dSub)
	}
	if dNon != 0 {
		t.Errorf("비구독자 할인 = %f, want 0", dNon)
	}
}

// TestDiscountCalculator_MinOrder는 최소 주문 금액 미달 시 할인 미적용을 검증합니다.
func TestDiscountCalculator_MinOrder(t *testing.T) {
	calc := service.NewDiscountCalculator()
	_, discount := calc.Calculate(5000, []service.Discount{
		{Type: service.DiscountPercent, Value: 10, MinOrder: 10000},
	}, false)

	if discount != 0 {
		t.Errorf("최소 주문 미달인데 할인 적용됨: %f", discount)
	}
}

// TestRecommendationEngine_DominantBiomarkers는 주력 바이오마커 추천을 검증합니다.
func TestRecommendationEngine_DominantBiomarkers(t *testing.T) {
	re := service.NewRecommendationEngine()
	profile := &service.UserMeasurementProfile{
		UserID:              "u1",
		DominantBiomarkers:  []string{"glucose"},
		AvgFrequencyPerWeek: 1,
		HasAbnormalReadings: false,
	}

	recs := re.Recommend(profile, 5)
	if len(recs) == 0 {
		t.Error("추천이 비어 있음")
	}

	hasGlucoseSKU := false
	for _, r := range recs {
		if r.ProductSKU == "BIO-GLU-001" {
			hasGlucoseSKU = true
			break
		}
	}
	if !hasGlucoseSKU {
		t.Error("glucose 카트리지 미추천")
	}
}

// TestRecommendationEngine_HighFrequencySubscription는 고빈도 사용자 정기구독 추천을 검증합니다.
func TestRecommendationEngine_HighFrequencySubscription(t *testing.T) {
	re := service.NewRecommendationEngine()
	profile := &service.UserMeasurementProfile{
		UserID:              "u1",
		AvgFrequencyPerWeek: 5,
	}

	recs := re.Recommend(profile, 10)
	hasSubscription := false
	for _, r := range recs {
		if r.ProductSKU == "SUB-MONTHLY-001" {
			hasSubscription = true
			break
		}
	}
	if !hasSubscription {
		t.Error("정기구독 미추천")
	}
}

// TestRecommendationEngine_AbnormalBoost는 이상 측정 시 점수 부스트를 검증합니다.
func TestRecommendationEngine_AbnormalBoost(t *testing.T) {
	re := service.NewRecommendationEngine()
	normal := &service.UserMeasurementProfile{
		DominantBiomarkers: []string{"glucose"},
		HasAbnormalReadings: false,
	}
	abnormal := &service.UserMeasurementProfile{
		DominantBiomarkers: []string{"glucose"},
		HasAbnormalReadings: true,
	}

	normalRecs := re.Recommend(normal, 1)
	abnormalRecs := re.Recommend(abnormal, 1)

	if len(normalRecs) == 0 || len(abnormalRecs) == 0 {
		t.Skip("추천이 비어 있음")
	}
	if abnormalRecs[0].Score <= normalRecs[0].Score {
		t.Errorf("이상 측정 점수 = %f, 정상 점수 = %f, 부스트 적용 안됨", abnormalRecs[0].Score, normalRecs[0].Score)
	}
}
