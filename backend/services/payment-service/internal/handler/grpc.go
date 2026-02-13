// Package handler는 payment-service의 gRPC 핸들러입니다.
package handler

import (
	"context"

	"github.com/manpasik/backend/services/payment-service/internal/service"
	apperrors "github.com/manpasik/backend/shared/errors"
	v1 "github.com/manpasik/backend/shared/gen/go/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// PaymentHandler는 PaymentService gRPC 서버를 구현합니다.
type PaymentHandler struct {
	v1.UnimplementedPaymentServiceServer
	svc *service.PaymentService
	log *zap.Logger
}

func NewPaymentHandler(svc *service.PaymentService, log *zap.Logger) *PaymentHandler {
	return &PaymentHandler{svc: svc, log: log}
}

func (h *PaymentHandler) CreatePayment(ctx context.Context, req *v1.CreatePaymentRequest) (*v1.PaymentDetail, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id는 필수입니다")
	}

	p, err := h.svc.CreatePayment(ctx, req.UserId, req.OrderId, req.SubscriptionId,
		service.PaymentType(req.PaymentType), req.AmountKrw, req.PaymentMethod)
	if err != nil {
		return nil, toGRPC(err)
	}
	return paymentToProto(p), nil
}

func (h *PaymentHandler) ConfirmPayment(ctx context.Context, req *v1.ConfirmPaymentRequest) (*v1.PaymentDetail, error) {
	if req == nil || req.PaymentId == "" {
		return nil, status.Error(codes.InvalidArgument, "payment_id는 필수입니다")
	}

	p, err := h.svc.ConfirmPayment(ctx, req.PaymentId, req.PgTransactionId, req.PgProvider, req.GetPaymentKey())
	if err != nil {
		return nil, toGRPC(err)
	}
	return paymentToProto(p), nil
}

func (h *PaymentHandler) GetPayment(ctx context.Context, req *v1.GetPaymentRequest) (*v1.PaymentDetail, error) {
	if req == nil || req.PaymentId == "" {
		return nil, status.Error(codes.InvalidArgument, "payment_id는 필수입니다")
	}

	p, err := h.svc.GetPayment(ctx, req.PaymentId)
	if err != nil {
		return nil, toGRPC(err)
	}
	return paymentToProto(p), nil
}

func (h *PaymentHandler) ListPayments(ctx context.Context, req *v1.ListPaymentsRequest) (*v1.ListPaymentsResponse, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id는 필수입니다")
	}

	payments, total, err := h.svc.ListPayments(ctx, req.UserId, req.Limit, req.Offset)
	if err != nil {
		return nil, toGRPC(err)
	}

	protoPayments := make([]*v1.PaymentDetail, 0, len(payments))
	for _, p := range payments {
		protoPayments = append(protoPayments, paymentToProto(p))
	}

	return &v1.ListPaymentsResponse{Payments: protoPayments, TotalCount: total}, nil
}

func (h *PaymentHandler) RefundPayment(ctx context.Context, req *v1.RefundPaymentRequest) (*v1.RefundResponse, error) {
	if req == nil || req.PaymentId == "" {
		return nil, status.Error(codes.InvalidArgument, "payment_id는 필수입니다")
	}

	refund, payment, err := h.svc.RefundPayment(ctx, req.PaymentId, req.RefundAmountKrw, req.Reason)
	if err != nil {
		return nil, toGRPC(err)
	}

	return &v1.RefundResponse{
		RefundId:        refund.ID,
		PaymentId:       refund.PaymentID,
		RefundAmountKrw: refund.RefundAmountKRW,
		PaymentStatus:   v1.PaymentStatus(payment.Status),
		RefundedAt:      timestamppb.New(refund.CreatedAt),
	}, nil
}

func paymentToProto(p *service.Payment) *v1.PaymentDetail {
	detail := &v1.PaymentDetail{
		PaymentId:       p.ID,
		UserId:          p.UserID,
		OrderId:         p.OrderID,
		SubscriptionId:  p.SubscriptionID,
		PaymentType:     v1.PaymentType(p.PaymentType),
		AmountKrw:       p.AmountKRW,
		Status:          v1.PaymentStatus(p.Status),
		PaymentMethod:   p.PaymentMethod,
		PgTransactionId: p.PgTransactionID,
		PgProvider:      p.PgProvider,
		CreatedAt:       timestamppb.New(p.CreatedAt),
	}
	if p.CompletedAt != nil {
		detail.CompletedAt = timestamppb.New(*p.CompletedAt)
	}
	return detail
}

func toGRPC(err error) error {
	if err == nil {
		return nil
	}
	if ae, ok := err.(*apperrors.AppError); ok {
		return ae.ToGRPC()
	}
	if s, ok := status.FromError(err); ok {
		return s.Err()
	}
	return status.Error(codes.Internal, "내부 오류가 발생했습니다")
}
