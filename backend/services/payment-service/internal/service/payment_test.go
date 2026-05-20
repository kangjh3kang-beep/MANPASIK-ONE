package service

import (
	"context"
	"fmt"
	"testing"

	"go.uber.org/zap"
)

type fakePayRepo struct {
	payments map[string]*Payment
	byUserID map[string][]*Payment
}

func newFakePayRepo() *fakePayRepo {
	return &fakePayRepo{
		payments: make(map[string]*Payment),
		byUserID: make(map[string][]*Payment),
	}
}

func (r *fakePayRepo) Create(_ context.Context, p *Payment) error {
	r.payments[p.ID] = p
	r.byUserID[p.UserID] = append(r.byUserID[p.UserID], p)
	return nil
}

func (r *fakePayRepo) GetByID(_ context.Context, id string) (*Payment, error) {
	return r.payments[id], nil
}

func (r *fakePayRepo) GetByOrderID(_ context.Context, orderID string) (*Payment, error) {
	for _, p := range r.payments {
		if p.OrderID == orderID {
			return p, nil
		}
	}
	return nil, nil
}

func (r *fakePayRepo) ListByUserID(_ context.Context, userID string, limit, offset int32) ([]*Payment, int32, error) {
	list := r.byUserID[userID]
	return list, int32(len(list)), nil
}

func (r *fakePayRepo) Update(_ context.Context, p *Payment) error {
	r.payments[p.ID] = p
	return nil
}

type fakeRefundRepo struct {
	refunds map[string][]*Refund
}

func newFakeRefundRepo() *fakeRefundRepo {
	return &fakeRefundRepo{refunds: make(map[string][]*Refund)}
}

func (r *fakeRefundRepo) Create(_ context.Context, ref *Refund) error {
	r.refunds[ref.PaymentID] = append(r.refunds[ref.PaymentID], ref)
	return nil
}

func (r *fakeRefundRepo) GetByPaymentID(_ context.Context, paymentID string) ([]*Refund, error) {
	return r.refunds[paymentID], nil
}

func newTestPaymentService() *PaymentService {
	return NewPaymentService(zap.NewNop(), newFakePayRepo(), newFakeRefundRepo())
}

func TestCreatePayment(t *testing.T) {
	svc := newTestPaymentService()
	ctx := context.Background()

	p, err := svc.CreatePayment(ctx, "user-1", "order-1", "", PaymentTypeOneTime, 29000, "card")
	if err != nil {
		t.Fatalf("CreatePayment 실패: %v", err)
	}
	if p.Status != PaymentStatusPending {
		t.Errorf("상태: got %d, want %d", p.Status, PaymentStatusPending)
	}
	if p.AmountKRW != 29000 {
		t.Errorf("금액: got %d, want 29000", p.AmountKRW)
	}
}

func TestCreatePayment_InvalidAmount(t *testing.T) {
	svc := newTestPaymentService()
	ctx := context.Background()

	_, err := svc.CreatePayment(ctx, "user-1", "order-1", "", PaymentTypeOneTime, 0, "card")
	if err == nil {
		t.Error("금액 0 결제 생성이 성공했습니다")
	}
}

func TestCreatePayment_NoMethod(t *testing.T) {
	svc := newTestPaymentService()
	ctx := context.Background()

	_, err := svc.CreatePayment(ctx, "user-1", "order-1", "", PaymentTypeOneTime, 10000, "")
	if err == nil {
		t.Error("결제 수단 없는 결제 생성이 성공했습니다")
	}
}

func TestConfirmPayment(t *testing.T) {
	svc := newTestPaymentService()
	ctx := context.Background()

	p, _ := svc.CreatePayment(ctx, "user-1", "order-1", "", PaymentTypeOneTime, 29000, "card")

	confirmed, err := svc.ConfirmPayment(ctx, p.ID, "pg-tx-123", "toss", "")
	if err != nil {
		t.Fatalf("ConfirmPayment 실패: %v", err)
	}
	if confirmed.Status != PaymentStatusCompleted {
		t.Errorf("상태: got %d, want %d", confirmed.Status, PaymentStatusCompleted)
	}
	if confirmed.PgTransactionID != "pg-tx-123" {
		t.Errorf("PG TX ID: got %s, want pg-tx-123", confirmed.PgTransactionID)
	}
	if confirmed.CompletedAt == nil {
		t.Error("CompletedAt이 설정되지 않았습니다")
	}
}

func TestConfirmPayment_AlreadyCompleted(t *testing.T) {
	svc := newTestPaymentService()
	ctx := context.Background()

	p, _ := svc.CreatePayment(ctx, "user-1", "order-1", "", PaymentTypeOneTime, 29000, "card")
	_, _ = svc.ConfirmPayment(ctx, p.ID, "pg-tx-1", "toss", "")

	_, err := svc.ConfirmPayment(ctx, p.ID, "pg-tx-2", "toss", "")
	if err == nil {
		t.Error("이미 완료된 결제 재확인이 성공했습니다")
	}
}

func TestGetPayment(t *testing.T) {
	svc := newTestPaymentService()
	ctx := context.Background()

	created, _ := svc.CreatePayment(ctx, "user-1", "order-1", "", PaymentTypeOneTime, 29000, "card")

	p, err := svc.GetPayment(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetPayment 실패: %v", err)
	}
	if p.ID != created.ID {
		t.Errorf("ID: got %s, want %s", p.ID, created.ID)
	}
}

func TestListPayments(t *testing.T) {
	svc := newTestPaymentService()
	ctx := context.Background()

	_, _ = svc.CreatePayment(ctx, "user-list", "order-1", "", PaymentTypeOneTime, 29000, "card")
	_, _ = svc.CreatePayment(ctx, "user-list", "order-2", "", PaymentTypeOneTime, 35000, "bank")

	payments, total, err := svc.ListPayments(ctx, "user-list", 20, 0)
	if err != nil {
		t.Fatalf("ListPayments 실패: %v", err)
	}
	if total != 2 {
		t.Errorf("전체 수: got %d, want 2", total)
	}
	if len(payments) != 2 {
		t.Errorf("반환 수: got %d, want 2", len(payments))
	}
}

func TestRefundPayment_Full(t *testing.T) {
	svc := newTestPaymentService()
	ctx := context.Background()

	p, _ := svc.CreatePayment(ctx, "user-1", "order-1", "", PaymentTypeOneTime, 29000, "card")
	_, _ = svc.ConfirmPayment(ctx, p.ID, "pg-tx-1", "toss", "")

	refund, payment, err := svc.RefundPayment(ctx, p.ID, 0, "테스트 환불")
	if err != nil {
		t.Fatalf("RefundPayment 실패: %v", err)
	}
	if refund.RefundAmountKRW != 29000 {
		t.Errorf("환불 금액: got %d, want 29000", refund.RefundAmountKRW)
	}
	if payment.Status != PaymentStatusRefunded {
		t.Errorf("결제 상태: got %d, want %d", payment.Status, PaymentStatusRefunded)
	}
}

func TestRefundPayment_Partial(t *testing.T) {
	svc := newTestPaymentService()
	ctx := context.Background()

	p, _ := svc.CreatePayment(ctx, "user-1", "order-1", "", PaymentTypeOneTime, 29000, "card")
	_, _ = svc.ConfirmPayment(ctx, p.ID, "pg-tx-1", "toss", "")

	refund, payment, err := svc.RefundPayment(ctx, p.ID, 10000, "부분 환불")
	if err != nil {
		t.Fatalf("부분 환불 실패: %v", err)
	}
	if refund.RefundAmountKRW != 10000 {
		t.Errorf("환불 금액: got %d, want 10000", refund.RefundAmountKRW)
	}
	if payment.Status != PaymentStatusPartialRefund {
		t.Errorf("결제 상태: got %d, want %d", payment.Status, PaymentStatusPartialRefund)
	}
}

func TestRefundPayment_ExceedAmount(t *testing.T) {
	svc := newTestPaymentService()
	ctx := context.Background()

	p, _ := svc.CreatePayment(ctx, "user-1", "order-1", "", PaymentTypeOneTime, 29000, "card")
	_, _ = svc.ConfirmPayment(ctx, p.ID, "pg-tx-1", "toss", "")

	_, _, err := svc.RefundPayment(ctx, p.ID, 50000, "초과 환불")
	if err == nil {
		t.Error("결제 금액 초과 환불이 성공했습니다")
	}
}

func TestRefundPayment_NotCompleted(t *testing.T) {
	svc := newTestPaymentService()
	ctx := context.Background()

	p, _ := svc.CreatePayment(ctx, "user-1", "order-1", "", PaymentTypeOneTime, 29000, "card")

	_, _, err := svc.RefundPayment(ctx, p.ID, 0, "미완료 결제 환불")
	if err == nil {
		t.Error("미완료 결제 환불이 성공했습니다")
	}
}

// =============================================================================
// Phase F 테스트 보강
// =============================================================================

func TestGetPayment_NotFound(t *testing.T) {
	svc := newTestPaymentService()
	_, err := svc.GetPayment(context.Background(), "nonexistent-id")
	if err == nil {
		t.Error("존재하지 않는 결제에 에러가 반환되어야 합니다")
	}
}

func TestListPayments_DefaultLimit(t *testing.T) {
	svc := newTestPaymentService()
	ctx := context.Background()

	_, _ = svc.CreatePayment(ctx, "user-dl", "order-1", "", PaymentTypeOneTime, 10000, "card")

	payments, total, err := svc.ListPayments(ctx, "user-dl", 0, 0)
	if err != nil {
		t.Fatalf("ListPayments 실패: %v", err)
	}
	if total != 1 {
		t.Errorf("total: got %d, want 1", total)
	}
	if len(payments) != 1 {
		t.Errorf("len: got %d, want 1", len(payments))
	}
}

func TestCreatePayment_SubscriptionType(t *testing.T) {
	svc := newTestPaymentService()
	ctx := context.Background()

	p, err := svc.CreatePayment(ctx, "user-1", "", "sub-1", PaymentTypeSubscription, 19900, "card")
	if err != nil {
		t.Fatalf("구독 결제 생성 실패: %v", err)
	}
	if p.PaymentType != PaymentTypeSubscription {
		t.Errorf("PaymentType: got %d, want %d", p.PaymentType, PaymentTypeSubscription)
	}
	if p.SubscriptionID != "sub-1" {
		t.Errorf("SubscriptionID: got %s, want sub-1", p.SubscriptionID)
	}
}

func TestConfirmPayment_NotFound(t *testing.T) {
	svc := newTestPaymentService()
	ctx := context.Background()

	_, err := svc.ConfirmPayment(ctx, "nonexistent", "pg-tx-1", "toss", "")
	if err == nil {
		t.Error("존재하지 않는 결제 확인에 에러가 반환되어야 합니다")
	}
}

func TestRefundPayment_AlreadyRefunded(t *testing.T) {
	svc := newTestPaymentService()
	ctx := context.Background()

	p, _ := svc.CreatePayment(ctx, "user-1", "order-1", "", PaymentTypeOneTime, 29000, "card")
	_, _ = svc.ConfirmPayment(ctx, p.ID, "pg-tx-1", "toss", "")
	_, _, _ = svc.RefundPayment(ctx, p.ID, 0, "전액 환불")

	_, _, err := svc.RefundPayment(ctx, p.ID, 0, "재환불")
	if err == nil {
		t.Error("이미 환불된 결제의 재환불에 에러가 반환되어야 합니다")
	}
}

func TestRefundPayment_NotFound(t *testing.T) {
	svc := newTestPaymentService()
	ctx := context.Background()

	_, _, err := svc.RefundPayment(ctx, "nonexistent", 0, "없는 결제")
	if err == nil {
		t.Error("존재하지 않는 결제 환불에 에러가 반환되어야 합니다")
	}
}

func TestCreatePayment_NegativeAmount(t *testing.T) {
	svc := newTestPaymentService()
	ctx := context.Background()

	_, err := svc.CreatePayment(ctx, "user-1", "order-1", "", PaymentTypeOneTime, -5000, "card")
	if err == nil {
		t.Error("음수 금액 결제에 에러가 반환되어야 합니다")
	}
}

func TestSetEventPublisher(t *testing.T) {
	svc := newTestPaymentService()
	// nil 설정 후에도 기본 동작 정상
	svc.SetEventPublisher(nil)

	ctx := context.Background()
	p, err := svc.CreatePayment(ctx, "user-1", "order-1", "", PaymentTypeOneTime, 10000, "card")
	if err != nil {
		t.Fatalf("결제 생성 실패: %v", err)
	}
	_, err = svc.ConfirmPayment(ctx, p.ID, "tx", "toss", "")
	if err != nil {
		t.Fatalf("결제 확인 실패: %v", err)
	}
}

func TestSetPaymentGateway(t *testing.T) {
	svc := newTestPaymentService()
	// nil 설정 후에도 기본 동작 정상
	svc.SetPaymentGateway(nil)

	ctx := context.Background()
	p, err := svc.CreatePayment(ctx, "user-1", "order-1", "", PaymentTypeOneTime, 10000, "card")
	if err != nil {
		t.Fatalf("결제 생성 실패: %v", err)
	}
	_, err = svc.ConfirmPayment(ctx, p.ID, "tx", "toss", "")
	if err != nil {
		t.Fatalf("결제 확인 실패: %v", err)
	}
}

// =============================================================================
// Mock PG Gateway & Event Publisher
// =============================================================================

type mockPGGateway struct {
	shouldFail bool
}

func (m *mockPGGateway) Confirm(_ context.Context, paymentKey, orderId string, amountKRW int32) (string, error) {
	if m.shouldFail {
		return "", fmt.Errorf("PG 승인 거부")
	}
	return "pg-tx-" + paymentKey, nil
}

func (m *mockPGGateway) Cancel(_ context.Context, paymentKey, reason string) error {
	if m.shouldFail {
		return fmt.Errorf("PG 취소 실패")
	}
	return nil
}

type mockEventPublisher struct {
	completedCount int
	failedCount    int
	refundedCount  int
}

func (m *mockEventPublisher) PublishPaymentCompleted(_ context.Context, event *PaymentCompletedEvent) error {
	m.completedCount++
	return nil
}

func (m *mockEventPublisher) PublishPaymentFailed(_ context.Context, event *PaymentFailedEvent) error {
	m.failedCount++
	return nil
}

func (m *mockEventPublisher) PublishPaymentRefunded(_ context.Context, event *PaymentRefundedEvent) error {
	m.refundedCount++
	return nil
}

func TestConfirmPayment_WithPGGateway(t *testing.T) {
	svc := newTestPaymentService()
	pg := &mockPGGateway{}
	svc.SetPaymentGateway(pg)
	ctx := context.Background()

	p, _ := svc.CreatePayment(ctx, "user-pg", "order-pg", "", PaymentTypeOneTime, 50000, "card")

	confirmed, err := svc.ConfirmPayment(ctx, p.ID, "", "toss", "pay-key-123")
	if err != nil {
		t.Fatalf("PG ConfirmPayment 실패: %v", err)
	}
	if confirmed.PgTransactionID != "pg-tx-pay-key-123" {
		t.Errorf("PG TX ID: got %s, want pg-tx-pay-key-123", confirmed.PgTransactionID)
	}
	if confirmed.PgProvider != "toss" {
		t.Errorf("PG Provider: got %s, want toss", confirmed.PgProvider)
	}
}

func TestConfirmPayment_PGGatewayFail(t *testing.T) {
	svc := newTestPaymentService()
	pg := &mockPGGateway{shouldFail: true}
	ep := &mockEventPublisher{}
	svc.SetPaymentGateway(pg)
	svc.SetEventPublisher(ep)
	ctx := context.Background()

	p, _ := svc.CreatePayment(ctx, "user-pgf", "order-pgf", "", PaymentTypeOneTime, 50000, "card")

	_, err := svc.ConfirmPayment(ctx, p.ID, "", "toss", "pay-key-fail")
	if err == nil {
		t.Fatal("PG 승인 실패 시 에러가 반환되어야 합니다")
	}
	if ep.failedCount != 1 {
		t.Errorf("실패 이벤트 발행 횟수: got %d, want 1", ep.failedCount)
	}
}

func TestConfirmPayment_WithEventPublisher(t *testing.T) {
	svc := newTestPaymentService()
	ep := &mockEventPublisher{}
	svc.SetEventPublisher(ep)
	ctx := context.Background()

	p, _ := svc.CreatePayment(ctx, "user-ep", "order-ep", "", PaymentTypeSubscription, 19900, "card")
	_, err := svc.ConfirmPayment(ctx, p.ID, "pg-tx-ep", "toss", "")
	if err != nil {
		t.Fatalf("ConfirmPayment 실패: %v", err)
	}
	if ep.completedCount != 1 {
		t.Errorf("완료 이벤트 발행 횟수: got %d, want 1", ep.completedCount)
	}
}

func TestRefundPayment_WithEventPublisher(t *testing.T) {
	svc := newTestPaymentService()
	ep := &mockEventPublisher{}
	svc.SetEventPublisher(ep)
	ctx := context.Background()

	p, _ := svc.CreatePayment(ctx, "user-ref-ep", "order-ref", "", PaymentTypeOneTime, 30000, "card")
	_, _ = svc.ConfirmPayment(ctx, p.ID, "pg-tx-ref", "toss", "")

	_, _, err := svc.RefundPayment(ctx, p.ID, 0, "환불 이벤트 테스트")
	if err != nil {
		t.Fatalf("RefundPayment 실패: %v", err)
	}
	if ep.refundedCount != 1 {
		t.Errorf("환불 이벤트 발행 횟수: got %d, want 1", ep.refundedCount)
	}
}

func TestRefundPayment_WithPGGateway(t *testing.T) {
	svc := newTestPaymentService()
	pg := &mockPGGateway{}
	svc.SetPaymentGateway(pg)
	ctx := context.Background()

	p, _ := svc.CreatePayment(ctx, "user-refpg", "order-refpg", "", PaymentTypeOneTime, 30000, "card")
	_, _ = svc.ConfirmPayment(ctx, p.ID, "", "toss", "pay-key-ref")

	_, _, err := svc.RefundPayment(ctx, p.ID, 0, "PG 환불 테스트")
	if err != nil {
		t.Fatalf("PG RefundPayment 실패: %v", err)
	}
}

func TestRefundPayment_PGCancelFail(t *testing.T) {
	svc := newTestPaymentService()
	pg := &mockPGGateway{shouldFail: true}
	svc.SetPaymentGateway(pg)
	ctx := context.Background()

	p, _ := svc.CreatePayment(ctx, "user-refpgf", "order-refpgf", "", PaymentTypeOneTime, 30000, "card")
	// PG 없이 결제 확인 (pgGateway 설정 전에 confirm, 또는 paymentKey 없이)
	svc.SetPaymentGateway(nil)
	_, _ = svc.ConfirmPayment(ctx, p.ID, "pg-tx-cancel-test", "toss", "")
	svc.SetPaymentGateway(pg)

	_, _, err := svc.RefundPayment(ctx, p.ID, 0, "PG 취소 실패 테스트")
	if err == nil {
		t.Fatal("PG 취소 실패 시 에러가 반환되어야 합니다")
	}
}

func TestConfirmPayment_DefaultProvider(t *testing.T) {
	svc := newTestPaymentService()
	pg := &mockPGGateway{}
	svc.SetPaymentGateway(pg)
	ctx := context.Background()

	p, _ := svc.CreatePayment(ctx, "user-def-prov", "order-dp", "", PaymentTypeOneTime, 10000, "card")

	// pgProvider를 빈 문자열로 → 기본 "toss"
	confirmed, err := svc.ConfirmPayment(ctx, p.ID, "", "", "pay-key-dp")
	if err != nil {
		t.Fatalf("ConfirmPayment 실패: %v", err)
	}
	if confirmed.PgProvider != "toss" {
		t.Errorf("기본 PgProvider: got %s, want toss", confirmed.PgProvider)
	}
}
