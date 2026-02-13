// Package service는 payment-service의 비즈니스 로직을 구현합니다.
package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	apperrors "github.com/manpasik/backend/shared/errors"
	"go.uber.org/zap"
)

// PaymentType는 결제 유형입니다.
type PaymentType int32

const (
	PaymentTypeUnknown      PaymentType = 0
	PaymentTypeOneTime      PaymentType = 1
	PaymentTypeSubscription PaymentType = 2
)

// PaymentStatus는 결제 상태입니다.
type PaymentStatus int32

const (
	PaymentStatusUnknown       PaymentStatus = 0
	PaymentStatusPending       PaymentStatus = 1
	PaymentStatusCompleted     PaymentStatus = 2
	PaymentStatusFailed        PaymentStatus = 3
	PaymentStatusCancelled     PaymentStatus = 4
	PaymentStatusRefunded      PaymentStatus = 5
	PaymentStatusPartialRefund PaymentStatus = 6
)

// Payment는 결제 엔티티입니다.
type Payment struct {
	ID              string
	UserID          string
	OrderID         string
	SubscriptionID  string
	PaymentType     PaymentType
	AmountKRW       int32
	Status          PaymentStatus
	PaymentMethod   string
	PgTransactionID string
	PgProvider      string
	CreatedAt       time.Time
	CompletedAt     *time.Time
}

// Refund는 환불 엔티티입니다.
type Refund struct {
	ID              string
	PaymentID       string
	RefundAmountKRW int32
	Reason          string
	CreatedAt       time.Time
}

// PaymentRepository는 결제 저장소 인터페이스입니다.
type PaymentRepository interface {
	Create(ctx context.Context, payment *Payment) error
	GetByID(ctx context.Context, id string) (*Payment, error)
	ListByUserID(ctx context.Context, userID string, limit, offset int32) ([]*Payment, int32, error)
	Update(ctx context.Context, payment *Payment) error
}

// RefundRepository는 환불 저장소 인터페이스입니다.
type RefundRepository interface {
	Create(ctx context.Context, refund *Refund) error
	GetByPaymentID(ctx context.Context, paymentID string) ([]*Refund, error)
}

// ============================================================================
// 이벤트 타입 (Kafka 발행용)
// ============================================================================

// PaymentCompletedEvent는 결제 완료 이벤트입니다.
type PaymentCompletedEvent struct {
	PaymentID      string
	UserID         string
	OrderID        string
	SubscriptionID string
	PaymentType    string
	AmountKRW      int32
	PaymentMethod  string
	PgProvider     string
	PgTransactionID string
	CompletedAt    time.Time
}

// PaymentFailedEvent는 결제 실패 이벤트입니다.
type PaymentFailedEvent struct {
	PaymentID    string
	UserID       string
	OrderID      string
	AmountKRW    int32
	ErrorCode    string
	ErrorMessage string
}

// PaymentRefundedEvent는 환불 완료 이벤트입니다.
type PaymentRefundedEvent struct {
	PaymentID      string
	RefundID       string
	UserID         string
	RefundAmountKRW int32
	Reason         string
	IsFullRefund   bool
}

// EventPublisher는 이벤트 발행 인터페이스입니다 (Kafka).
type EventPublisher interface {
	PublishPaymentCompleted(ctx context.Context, event *PaymentCompletedEvent) error
	PublishPaymentFailed(ctx context.Context, event *PaymentFailedEvent) error
	PublishPaymentRefunded(ctx context.Context, event *PaymentRefundedEvent) error
}

// PaymentGateway는 PG사(Toss 등) 결제 승인·취소 연동 인터페이스입니다.
// B-3: Toss 연동 시 이 인터페이스 구현체를 주입합니다. 미설정 시 Noop 사용.
type PaymentGateway interface {
	// Confirm은 PG사에 결제 승인을 요청합니다. paymentKey는 PG에서 발급한 키, orderId는 우리 주문 ID, amountKRW는 금액(원).
	Confirm(ctx context.Context, paymentKey, orderId string, amountKRW int32) (pgTransactionID string, err error)
	// Cancel은 PG사에 결제 취소(환불)를 요청합니다.
	Cancel(ctx context.Context, paymentKey, reason string) error
}

// PaymentService는 결제 비즈니스 로직입니다.
type PaymentService struct {
	logger         *zap.Logger
	payRepo        PaymentRepository
	refundRepo     RefundRepository
	eventPublisher EventPublisher
	pgGateway      PaymentGateway // optional: Toss 등 실 PG 연동
}

// NewPaymentService는 새 PaymentService를 생성합니다.
func NewPaymentService(
	logger *zap.Logger,
	payRepo PaymentRepository,
	refundRepo RefundRepository,
) *PaymentService {
	return &PaymentService{
		logger:     logger,
		payRepo:    payRepo,
		refundRepo: refundRepo,
	}
}

// SetEventPublisher는 이벤트 발행기를 설정합니다 (optional).
func (s *PaymentService) SetEventPublisher(ep EventPublisher) {
	s.eventPublisher = ep
}

// SetPaymentGateway는 PG 연동 클라이언트를 설정합니다 (optional). B-3 Toss 연동 시 사용.
func (s *PaymentService) SetPaymentGateway(pg PaymentGateway) {
	s.pgGateway = pg
}

// CreatePayment는 결제 요청을 생성합니다.
func (s *PaymentService) CreatePayment(ctx context.Context, userID, orderID, subscriptionID string, paymentType PaymentType, amountKRW int32, paymentMethod string) (*Payment, error) {
	if amountKRW <= 0 {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "결제 금액은 0보다 커야 합니다")
	}

	if paymentMethod == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "결제 수단은 필수입니다")
	}

	now := time.Now().UTC()
	payment := &Payment{
		ID:             uuid.New().String(),
		UserID:         userID,
		OrderID:        orderID,
		SubscriptionID: subscriptionID,
		PaymentType:    paymentType,
		AmountKRW:      amountKRW,
		Status:         PaymentStatusPending,
		PaymentMethod:  paymentMethod,
		CreatedAt:      now,
	}

	if err := s.payRepo.Create(ctx, payment); err != nil {
		s.logger.Error("결제 생성 실패", zap.Error(err))
		return nil, apperrors.New(apperrors.ErrInternal, "결제 생성에 실패했습니다")
	}

	s.logger.Info("결제 요청 생성",
		zap.String("payment_id", payment.ID),
		zap.String("user_id", userID),
		zap.Int32("amount", amountKRW),
	)
	return payment, nil
}

// ConfirmPayment는 PG 콜백으로 결제를 확인합니다.
// paymentKey가 비어 있지 않고 pgGateway가 설정된 경우 Toss 등 PG 승인 API를 호출한 뒤 DB를 갱신합니다.
func (s *PaymentService) ConfirmPayment(ctx context.Context, paymentID, pgTransactionID, pgProvider, paymentKey string) (*Payment, error) {
	payment, err := s.payRepo.GetByID(ctx, paymentID)
	if err != nil || payment == nil {
		return nil, apperrors.New(apperrors.ErrNotFound, "결제를 찾을 수 없습니다")
	}

	if payment.Status != PaymentStatusPending {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "대기 상태의 결제만 확인할 수 있습니다")
	}

	// PG 연동: paymentKey가 있고 gateway가 있으면 승인 API 호출
	if s.pgGateway != nil && paymentKey != "" {
		txID, pgErr := s.pgGateway.Confirm(ctx, paymentKey, payment.OrderID, payment.AmountKRW)
		if pgErr != nil {
			s.logger.Warn("PG 승인 실패", zap.String("payment_id", paymentID), zap.Error(pgErr))
			if s.eventPublisher != nil {
				_ = s.eventPublisher.PublishPaymentFailed(ctx, &PaymentFailedEvent{
					PaymentID: paymentID, UserID: payment.UserID, OrderID: payment.OrderID,
					AmountKRW: payment.AmountKRW, ErrorCode: "PG_CONFIRM_FAILED", ErrorMessage: pgErr.Error(),
				})
			}
			return nil, apperrors.New(apperrors.ErrInternal, "결제 승인에 실패했습니다: "+pgErr.Error())
		}
		pgTransactionID = txID
		if pgProvider == "" {
			pgProvider = "toss"
		}
	}

	now := time.Now().UTC()
	payment.Status = PaymentStatusCompleted
	payment.PgTransactionID = pgTransactionID
	payment.PgProvider = pgProvider
	payment.CompletedAt = &now

	if err := s.payRepo.Update(ctx, payment); err != nil {
		s.logger.Error("결제 확인 실패", zap.Error(err))
		return nil, apperrors.New(apperrors.ErrInternal, "결제 확인에 실패했습니다")
	}

	s.logger.Info("결제 완료",
		zap.String("payment_id", paymentID),
		zap.String("pg_tx_id", pgTransactionID),
	)

	// 결제 완료 이벤트 발행 (Kafka, optional)
	if s.eventPublisher != nil {
		ptStr := "one_time"
		if payment.PaymentType == PaymentTypeSubscription {
			ptStr = "subscription"
		}
		evt := &PaymentCompletedEvent{
			PaymentID:       payment.ID,
			UserID:          payment.UserID,
			OrderID:         payment.OrderID,
			SubscriptionID:  payment.SubscriptionID,
			PaymentType:     ptStr,
			AmountKRW:       payment.AmountKRW,
			PaymentMethod:   payment.PaymentMethod,
			PgProvider:      pgProvider,
			PgTransactionID: pgTransactionID,
			CompletedAt:     now,
		}
		if err := s.eventPublisher.PublishPaymentCompleted(ctx, evt); err != nil {
			s.logger.Warn("결제 완료 이벤트 발행 실패 (비치명적)", zap.Error(err))
		}
	}

	return payment, nil
}

// GetPayment는 결제 상세를 조회합니다.
func (s *PaymentService) GetPayment(ctx context.Context, paymentID string) (*Payment, error) {
	payment, err := s.payRepo.GetByID(ctx, paymentID)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrInternal, "결제 조회에 실패했습니다")
	}
	if payment == nil {
		return nil, apperrors.New(apperrors.ErrNotFound, "결제를 찾을 수 없습니다")
	}
	return payment, nil
}

// ListPayments는 사용자의 결제 이력을 조회합니다.
func (s *PaymentService) ListPayments(ctx context.Context, userID string, limit, offset int32) ([]*Payment, int32, error) {
	if limit <= 0 {
		limit = 20
	}
	payments, total, err := s.payRepo.ListByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, 0, apperrors.New(apperrors.ErrInternal, "결제 이력 조회에 실패했습니다")
	}
	return payments, total, nil
}

// RefundPayment는 환불을 처리합니다.
func (s *PaymentService) RefundPayment(ctx context.Context, paymentID string, refundAmountKRW int32, reason string) (*Refund, *Payment, error) {
	payment, err := s.payRepo.GetByID(ctx, paymentID)
	if err != nil || payment == nil {
		return nil, nil, apperrors.New(apperrors.ErrNotFound, "결제를 찾을 수 없습니다")
	}

	if payment.Status != PaymentStatusCompleted {
		return nil, nil, apperrors.New(apperrors.ErrInvalidInput, "완료된 결제만 환불할 수 있습니다")
	}

	// 전액 환불인 경우
	if refundAmountKRW <= 0 {
		refundAmountKRW = payment.AmountKRW
	}

	if refundAmountKRW > payment.AmountKRW {
		return nil, nil, apperrors.New(apperrors.ErrInvalidInput, "환불 금액이 결제 금액을 초과합니다")
	}

	// PG 연동: gateway가 있고 PgTransactionID가 있으면 취소 API 호출 (전액 환불만 Toss 취소 호출)
	if s.pgGateway != nil && payment.PgTransactionID != "" && refundAmountKRW == payment.AmountKRW {
		if cancelErr := s.pgGateway.Cancel(ctx, payment.PgTransactionID, reason); cancelErr != nil {
			s.logger.Warn("PG 취소 실패", zap.String("payment_id", paymentID), zap.Error(cancelErr))
			return nil, nil, apperrors.New(apperrors.ErrInternal, "결제 취소에 실패했습니다: "+cancelErr.Error())
		}
	}

	now := time.Now().UTC()
	refund := &Refund{
		ID:              uuid.New().String(),
		PaymentID:       paymentID,
		RefundAmountKRW: refundAmountKRW,
		Reason:          reason,
		CreatedAt:       now,
	}

	if err := s.refundRepo.Create(ctx, refund); err != nil {
		return nil, nil, apperrors.New(apperrors.ErrInternal, "환불 처리에 실패했습니다")
	}

	if refundAmountKRW == payment.AmountKRW {
		payment.Status = PaymentStatusRefunded
	} else {
		payment.Status = PaymentStatusPartialRefund
	}

	if err := s.payRepo.Update(ctx, payment); err != nil {
		return nil, nil, apperrors.New(apperrors.ErrInternal, "결제 상태 업데이트에 실패했습니다")
	}

	s.logger.Info("환불 완료",
		zap.String("payment_id", paymentID),
		zap.Int32("refund_amount", refundAmountKRW),
	)

	// 환불 이벤트 발행 (Kafka, optional)
	if s.eventPublisher != nil {
		evt := &PaymentRefundedEvent{
			PaymentID:       paymentID,
			RefundID:        refund.ID,
			UserID:          payment.UserID,
			RefundAmountKRW: refundAmountKRW,
			Reason:          reason,
			IsFullRefund:    refundAmountKRW == payment.AmountKRW,
		}
		if err := s.eventPublisher.PublishPaymentRefunded(ctx, evt); err != nil {
			s.logger.Warn("환불 이벤트 발행 실패 (비치명적)", zap.Error(err))
		}
	}

	return refund, payment, nil
}
