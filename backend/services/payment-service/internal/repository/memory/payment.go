// Package memory는 인메모리 결제 저장소입니다 (개발/테스트용).
package memory

import (
	"context"
	"sync"

	"github.com/manpasik/backend/services/payment-service/internal/service"
)

// PaymentRepository는 인메모리 결제 저장소입니다.
type PaymentRepository struct {
	mu       sync.RWMutex
	payments map[string]*service.Payment
	byUserID map[string][]*service.Payment
}

func NewPaymentRepository() *PaymentRepository {
	return &PaymentRepository{
		payments: make(map[string]*service.Payment),
		byUserID: make(map[string][]*service.Payment),
	}
}

func (r *PaymentRepository) Create(_ context.Context, p *service.Payment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *p
	r.payments[p.ID] = &cp
	r.byUserID[p.UserID] = append(r.byUserID[p.UserID], &cp)
	return nil
}

func (r *PaymentRepository) GetByID(_ context.Context, id string) (*service.Payment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.payments[id]
	if !ok {
		return nil, nil
	}
	cp := *p
	return &cp, nil
}

func (r *PaymentRepository) ListByUserID(_ context.Context, userID string, limit, offset int32) ([]*service.Payment, int32, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := r.byUserID[userID]
	total := int32(len(list))
	return list, total, nil
}

func (r *PaymentRepository) Update(_ context.Context, p *service.Payment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *p
	r.payments[p.ID] = &cp
	// byUserID도 업데이트
	list := r.byUserID[p.UserID]
	for i, existing := range list {
		if existing.ID == p.ID {
			list[i] = &cp
			break
		}
	}
	return nil
}

// RefundRepository는 인메모리 환불 저장소입니다.
type RefundRepository struct {
	mu      sync.RWMutex
	refunds map[string][]*service.Refund // paymentID -> refunds
}

func NewRefundRepository() *RefundRepository {
	return &RefundRepository{refunds: make(map[string][]*service.Refund)}
}

func (r *RefundRepository) Create(_ context.Context, ref *service.Refund) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *ref
	r.refunds[ref.PaymentID] = append(r.refunds[ref.PaymentID], &cp)
	return nil
}

func (r *RefundRepository) GetByPaymentID(_ context.Context, paymentID string) ([]*service.Refund, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.refunds[paymentID], nil
}
