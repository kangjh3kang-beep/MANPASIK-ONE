package service

import (
	"context"
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
