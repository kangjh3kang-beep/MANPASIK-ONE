// Package postgres는 payment-service의 PostgreSQL 저장소 구현입니다.
package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manpasik/backend/services/payment-service/internal/service"
)

// ============================================================================
// PaymentRepository
// ============================================================================

// PaymentRepository는 PostgreSQL 기반 PaymentRepository 구현입니다.
type PaymentRepository struct {
	pool *pgxpool.Pool
}

// NewPaymentRepository는 PaymentRepository를 생성합니다.
func NewPaymentRepository(pool *pgxpool.Pool) *PaymentRepository {
	return &PaymentRepository{pool: pool}
}

// Create는 결제를 생성합니다.
func (r *PaymentRepository) Create(ctx context.Context, payment *service.Payment) error {
	const q = `INSERT INTO payments
		(id, user_id, order_id, amount, currency, type, status, pg_tx_id, description, created_at, updated_at)
		VALUES ($1,$2,$3,$4,'KRW',$5,$6,$7,$8,$9,$10)`
	_, err := r.pool.Exec(ctx, q,
		payment.ID,
		payment.UserID,
		payment.OrderID,
		int64(payment.AmountKRW),
		mapPaymentType(payment.PaymentType),
		mapPaymentStatus(payment.Status),
		nilIfEmpty(payment.PgTransactionID),
		nilIfEmpty(payment.PaymentMethod), // store method in description column
		payment.CreatedAt,
		payment.CreatedAt, // updated_at = created_at on insert
	)
	return err
}

// GetByID는 결제 ID로 조회합니다.
func (r *PaymentRepository) GetByID(ctx context.Context, id string) (*service.Payment, error) {
	const q = `SELECT id, user_id, order_id, amount, type, status,
		COALESCE(pg_tx_id, ''), COALESCE(description, ''), created_at, updated_at
		FROM payments WHERE id = $1`

	var p service.Payment
	var amount int64
	var payType, payStatus string
	err := r.pool.QueryRow(ctx, q, id).Scan(
		&p.ID, &p.UserID, &p.OrderID, &amount,
		&payType, &payStatus,
		&p.PgTransactionID, &p.PaymentMethod,
		&p.CreatedAt, &p.CompletedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	p.AmountKRW = int32(amount)
	p.PaymentType = reversePaymentType(payType)
	p.Status = reversePaymentStatus(payStatus)
	return &p, nil
}

// ListByUserID는 사용자의 결제 이력을 조회합니다.
func (r *PaymentRepository) ListByUserID(ctx context.Context, userID string, limit, offset int32) ([]*service.Payment, int32, error) {
	var total int32
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM payments WHERE user_id = $1`, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	const q = `SELECT id, user_id, order_id, amount, type, status,
		COALESCE(pg_tx_id, ''), COALESCE(description, ''), created_at, updated_at
		FROM payments WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.pool.Query(ctx, q, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var payments []*service.Payment
	for rows.Next() {
		var p service.Payment
		var amount int64
		var payType, payStatus string
		if err := rows.Scan(
			&p.ID, &p.UserID, &p.OrderID, &amount,
			&payType, &payStatus,
			&p.PgTransactionID, &p.PaymentMethod,
			&p.CreatedAt, &p.CompletedAt,
		); err != nil {
			return nil, 0, err
		}
		p.AmountKRW = int32(amount)
		p.PaymentType = reversePaymentType(payType)
		p.Status = reversePaymentStatus(payStatus)
		payments = append(payments, &p)
	}
	return payments, total, rows.Err()
}

// Update는 결제 정보를 업데이트합니다.
func (r *PaymentRepository) Update(ctx context.Context, payment *service.Payment) error {
	const q = `UPDATE payments SET
		status=$1, pg_tx_id=$2, updated_at=NOW()
		WHERE id=$3`
	_, err := r.pool.Exec(ctx, q,
		mapPaymentStatus(payment.Status),
		nilIfEmpty(payment.PgTransactionID),
		payment.ID,
	)
	return err
}

// ============================================================================
// RefundRepository
// ============================================================================

// RefundRepository는 PostgreSQL 기반 RefundRepository 구현입니다.
type RefundRepository struct {
	pool *pgxpool.Pool
}

// NewRefundRepository는 RefundRepository를 생성합니다.
func NewRefundRepository(pool *pgxpool.Pool) *RefundRepository {
	return &RefundRepository{pool: pool}
}

// Create는 환불을 생성합니다.
func (r *RefundRepository) Create(ctx context.Context, refund *service.Refund) error {
	const q = `INSERT INTO refunds (id, payment_id, amount, reason, status, created_at)
		VALUES ($1,$2,$3,$4,'COMPLETED',$5)`
	_, err := r.pool.Exec(ctx, q,
		refund.ID,
		refund.PaymentID,
		int64(refund.RefundAmountKRW),
		refund.Reason,
		refund.CreatedAt,
	)
	return err
}

// GetByPaymentID는 결제 ID에 해당하는 환불 목록을 조회합니다.
func (r *RefundRepository) GetByPaymentID(ctx context.Context, paymentID string) ([]*service.Refund, error) {
	const q = `SELECT id, payment_id, amount, COALESCE(reason, ''), created_at
		FROM refunds WHERE payment_id = $1 ORDER BY created_at DESC`
	rows, err := r.pool.Query(ctx, q, paymentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var refunds []*service.Refund
	for rows.Next() {
		var rf service.Refund
		var amount int64
		if err := rows.Scan(
			&rf.ID, &rf.PaymentID, &amount, &rf.Reason, &rf.CreatedAt,
		); err != nil {
			return nil, err
		}
		rf.RefundAmountKRW = int32(amount)
		refunds = append(refunds, &rf)
	}
	return refunds, rows.Err()
}

// ============================================================================
// 헬퍼: Go enum ↔ DB ENUM 매핑
// ============================================================================

func mapPaymentType(t service.PaymentType) string {
	switch t {
	case service.PaymentTypeOneTime:
		return "CARD"
	case service.PaymentTypeSubscription:
		return "BANK_TRANSFER"
	default:
		return "CARD"
	}
}

func reversePaymentType(s string) service.PaymentType {
	switch s {
	case "CARD", "MOBILE":
		return service.PaymentTypeOneTime
	case "BANK_TRANSFER", "VIRTUAL_ACCOUNT":
		return service.PaymentTypeSubscription
	default:
		return service.PaymentTypeUnknown
	}
}

func mapPaymentStatus(s service.PaymentStatus) string {
	switch s {
	case service.PaymentStatusPending:
		return "PENDING"
	case service.PaymentStatusCompleted:
		return "CONFIRMED"
	case service.PaymentStatusFailed:
		return "FAILED"
	case service.PaymentStatusCancelled:
		return "FAILED"
	case service.PaymentStatusRefunded:
		return "REFUNDED"
	case service.PaymentStatusPartialRefund:
		return "PARTIAL_REFUND"
	default:
		return "PENDING"
	}
}

func reversePaymentStatus(s string) service.PaymentStatus {
	switch s {
	case "PENDING":
		return service.PaymentStatusPending
	case "CONFIRMED":
		return service.PaymentStatusCompleted
	case "FAILED":
		return service.PaymentStatusFailed
	case "REFUNDED":
		return service.PaymentStatusRefunded
	case "PARTIAL_REFUND":
		return service.PaymentStatusPartialRefund
	default:
		return service.PaymentStatusUnknown
	}
}

func nilIfEmpty(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
