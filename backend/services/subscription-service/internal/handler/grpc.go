// Package handler는 subscription-service의 gRPC 핸들러입니다.
package handler

import (
	"context"

	"github.com/manpasik/backend/services/subscription-service/internal/service"
	apperrors "github.com/manpasik/backend/shared/errors"
	v1 "github.com/manpasik/backend/shared/gen/go/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// SubscriptionHandler는 SubscriptionService gRPC 서버를 구현합니다.
type SubscriptionHandler struct {
	v1.UnimplementedSubscriptionServiceServer
	svc *service.SubscriptionService
	log *zap.Logger
}

// NewSubscriptionHandler는 SubscriptionHandler를 생성합니다.
func NewSubscriptionHandler(svc *service.SubscriptionService, log *zap.Logger) *SubscriptionHandler {
	return &SubscriptionHandler{svc: svc, log: log}
}

// CreateSubscription은 구독 생성 RPC입니다.
func (h *SubscriptionHandler) CreateSubscription(ctx context.Context, req *v1.CreateSubscriptionRequest) (*v1.SubscriptionDetail, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id는 필수입니다")
	}

	tier := protoTierToService(req.Tier)
	sub, err := h.svc.CreateSubscription(ctx, req.UserId, tier)
	if err != nil {
		return nil, toGRPC(err)
	}

	return subscriptionToProto(sub), nil
}

// GetSubscription은 구독 조회 RPC입니다.
func (h *SubscriptionHandler) GetSubscription(ctx context.Context, req *v1.GetSubscriptionDetailRequest) (*v1.SubscriptionDetail, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id는 필수입니다")
	}

	sub, err := h.svc.GetSubscription(ctx, req.UserId)
	if err != nil {
		return nil, toGRPC(err)
	}

	return subscriptionToProto(sub), nil
}

// UpdateSubscription은 구독 업데이트 RPC입니다.
func (h *SubscriptionHandler) UpdateSubscription(ctx context.Context, req *v1.UpdateSubscriptionRequest) (*v1.SubscriptionDetail, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id는 필수입니다")
	}

	newTier := protoTierToService(req.NewTier)
	sub, err := h.svc.UpdateSubscription(ctx, req.UserId, newTier)
	if err != nil {
		return nil, toGRPC(err)
	}

	return subscriptionToProto(sub), nil
}

// CancelSubscription은 구독 해지 RPC입니다.
func (h *SubscriptionHandler) CancelSubscription(ctx context.Context, req *v1.CancelSubscriptionRequest) (*v1.CancelSubscriptionResponse, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id는 필수입니다")
	}

	sub, err := h.svc.CancelSubscription(ctx, req.UserId, req.Reason)
	if err != nil {
		return nil, toGRPC(err)
	}

	resp := &v1.CancelSubscriptionResponse{
		Success:        true,
		EffectiveUntil: timestamppb.New(sub.ExpiresAt),
	}
	if sub.CancelledAt != nil {
		resp.CancelledAt = timestamppb.New(*sub.CancelledAt)
	}

	return resp, nil
}

// CheckFeatureAccess는 기능 접근 권한 확인 RPC입니다.
func (h *SubscriptionHandler) CheckFeatureAccess(ctx context.Context, req *v1.CheckFeatureAccessRequest) (*v1.CheckFeatureAccessResponse, error) {
	if req == nil || req.UserId == "" || req.FeatureName == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id와 feature_name은 필수입니다")
	}

	allowed, requiredTier, currentTier, msg := h.svc.CheckFeatureAccess(ctx, req.UserId, req.FeatureName)

	return &v1.CheckFeatureAccessResponse{
		Allowed:      allowed,
		RequiredTier: serviceTierToProto(requiredTier),
		CurrentTier:  serviceTierToProto(currentTier),
		Message:      msg,
	}, nil
}

// ListSubscriptionPlans는 구독 플랜 목록 RPC입니다.
func (h *SubscriptionHandler) ListSubscriptionPlans(_ context.Context, _ *v1.ListSubscriptionPlansRequest) (*v1.ListSubscriptionPlansResponse, error) {
	planList := h.svc.ListSubscriptionPlans()

	protoPlans := make([]*v1.SubscriptionPlan, 0, len(planList))
	for _, p := range planList {
		protoPlans = append(protoPlans, &v1.SubscriptionPlan{
			Tier:                serviceTierToProto(p.Tier),
			Name:                p.Name,
			Description:         p.Description,
			MonthlyPriceKrw:     p.MonthlyPriceKRW,
			MaxDevices:          p.MaxDevices,
			MaxFamilyMembers:    p.MaxFamilyMembers,
			AiCoachingEnabled:   p.AICoachingEnabled,
			TelemedicineEnabled: p.TelemedicineEnabled,
			Features:            p.Features,
		})
	}

	return &v1.ListSubscriptionPlansResponse{Plans: protoPlans}, nil
}

// --- 헬퍼 함수 ---

func subscriptionToProto(sub *service.Subscription) *v1.SubscriptionDetail {
	detail := &v1.SubscriptionDetail{
		SubscriptionId:      sub.ID,
		UserId:              sub.UserID,
		Tier:                serviceTierToProto(sub.Tier),
		Status:              serviceStatusToProto(sub.Status),
		StartedAt:           timestamppb.New(sub.StartedAt),
		ExpiresAt:           timestamppb.New(sub.ExpiresAt),
		MaxDevices:          sub.MaxDevices,
		MaxFamilyMembers:    sub.MaxFamilyMembers,
		AiCoachingEnabled:   sub.AICoachingEnabled,
		TelemedicineEnabled: sub.TelemedicineEnabled,
		MonthlyPriceKrw:     sub.MonthlyPriceKRW,
		AutoRenew:           sub.AutoRenew,
	}
	if sub.CancelledAt != nil {
		detail.CancelledAt = timestamppb.New(*sub.CancelledAt)
	}
	return detail
}

func protoTierToService(tier v1.SubscriptionTier) service.SubscriptionTier {
	switch tier {
	case v1.SubscriptionTier_SUBSCRIPTION_TIER_BASIC:
		return service.TierBasic
	case v1.SubscriptionTier_SUBSCRIPTION_TIER_PRO:
		return service.TierPro
	case v1.SubscriptionTier_SUBSCRIPTION_TIER_CLINICAL:
		return service.TierClinical
	default:
		return service.TierFree
	}
}

func serviceTierToProto(tier service.SubscriptionTier) v1.SubscriptionTier {
	switch tier {
	case service.TierBasic:
		return v1.SubscriptionTier_SUBSCRIPTION_TIER_BASIC
	case service.TierPro:
		return v1.SubscriptionTier_SUBSCRIPTION_TIER_PRO
	case service.TierClinical:
		return v1.SubscriptionTier_SUBSCRIPTION_TIER_CLINICAL
	default:
		return v1.SubscriptionTier_SUBSCRIPTION_TIER_FREE
	}
}

func serviceStatusToProto(s service.SubscriptionStatus) v1.SubscriptionStatus {
	switch s {
	case service.StatusActive:
		return v1.SubscriptionStatus_SUBSCRIPTION_STATUS_ACTIVE
	case service.StatusCancelled:
		return v1.SubscriptionStatus_SUBSCRIPTION_STATUS_CANCELLED
	case service.StatusExpired:
		return v1.SubscriptionStatus_SUBSCRIPTION_STATUS_EXPIRED
	case service.StatusSuspended:
		return v1.SubscriptionStatus_SUBSCRIPTION_STATUS_SUSPENDED
	case service.StatusTrial:
		return v1.SubscriptionStatus_SUBSCRIPTION_STATUS_TRIAL
	default:
		return v1.SubscriptionStatus_SUBSCRIPTION_STATUS_UNKNOWN
	}
}

// toGRPC는 AppError를 gRPC status로 변환합니다.
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
